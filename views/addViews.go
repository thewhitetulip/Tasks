package views

/*
Holds the insert task related view handlers, includes the one for file upload
*/
import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/thewhitetulip/Tasks/db"
	"github.com/thewhitetulip/Tasks/utils"
)

// UploadedFileHandler is used to handle the uploaded file related requests
func UploadedFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		token := r.URL.Path[len("/files/"):]

		//file, err := db.GetFileName(token)
		//if err != nil {
		log.Println("serving file ./files/" + token)
		http.ServeFile(w, r, "./files/"+token)
		//}
	}
}

//AddTaskFunc is used to handle the addition of new task, "/add" URL
func AddTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" { // Will work only for POST requests, will redirect to home
		var filelink string // will store the html when we have files to be uploaded, appened to the note content
		r.ParseForm()
		file, handler, err := r.FormFile("uploadfile")
		if err != nil && handler != nil {
			//Case executed when file is uploaded and yet an error occurs
			log.Println(err)
			message = "Error uploading file"
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}

		taskPriority, priorityErr := strconv.Atoi(r.FormValue("priority"))

		if priorityErr != nil {
			log.Print(priorityErr)
			message = "Bad task priority"
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
		priorityList := []int{1, 2, 3}
		found := false
		for _, priority := range priorityList {
			if taskPriority == priority {
				found = true
			}
		}
		//If someone gives us incorrect priority number, we give the priority
		//to that task as 1 i.e. Low
		if !found {
			taskPriority = 1
		}

		category := r.FormValue("category")
		title := template.HTMLEscapeString(r.Form.Get("title"))
		content := template.HTMLEscapeString(r.Form.Get("content"))
		formToken := template.HTMLEscapeString(r.Form.Get("CSRFToken"))

		cookie, _ := r.Cookie("csrftoken")
		if formToken == cookie.Value {
			if handler != nil {
				// this will be executed whenever a file is uploaded
				r.ParseMultipartForm(32 << 20) //defined maximum size of file
				defer file.Close()
				randomFileName := md5.New()
				io.WriteString(randomFileName, strconv.FormatInt(time.Now().Unix(), 10))
				io.WriteString(randomFileName, handler.Filename)
				token := fmt.Sprintf("%x", randomFileName.Sum(nil))
				f, err := os.OpenFile("./files/"+token, os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					log.Println(err)
					return
				}
				defer f.Close()
				io.Copy(f, file)

				if strings.HasSuffix(handler.Filename, ".png") || strings.HasSuffix(handler.Filename, ".jpg") {
					filelink = "<br> <img src='/files/" + token + "'/>"
				} else {
					filelink = "<br> <a href=/files/" + token + ">" + handler.Filename + "</a>"
				}
				content = content + filelink

				fileTruth := db.AddFile(handler.Filename, token)
				if fileTruth != nil {
					message = "Error adding filename in db"
					log.Println("error adding task to db")
				}
			}

			taskTruth := db.AddTask(title, content, category, taskPriority)

			if taskTruth != nil {
				message = "Error adding task"
				log.Println("error adding task to db")
				http.Redirect(w, r, "/", http.StatusInternalServerError)
			} else {
				message = "Task added"
				log.Println("added task to db")
				http.Redirect(w, r, "/", http.StatusFound)
			}
		} else {
			log.Println("CSRF mismatch")
			message = "Server Error"
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}

	}

}

//AddCategoryFunc used to add new categories to the database
func AddCategoryFunc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	category := r.Form.Get("category")
	if strings.Trim(category, " ") != "" {
		err := db.AddCategory(category)
		if err != nil {
			message = "Error adding category"
			http.Redirect(w, r, "/", http.StatusBadRequest)
		} else {
			message = "Added category"
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

//EditTaskFunc is used to edit tasks, handles "/edit/" URL
func EditTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id, err := strconv.Atoi(r.URL.Path[len("/edit/"):])
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusBadRequest)
		} else {
			redirectURL := utils.GetRedirectUrl(r.Referer())
			task, err := db.GetTaskByID(id)
			categories := db.GetCategories()
			task.Categories = categories
			task.Referer = redirectURL

			if err != nil {
				task.Message = "Error fetching Tasks"
			}
			editTemplate.Execute(w, task)
		}
	}
}

//AddCommentFunc will be used
func AddCommentFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		text := r.Form.Get("commentText")
		id := r.Form.Get("taskID")

		idInt, err := strconv.Atoi(id)

		if (err != nil) || (text == "") {
			log.Println("unable to convert into integer")
			message = "Error adding comment"
		} else {
			err = db.AddComments(idInt, text)

			if err != nil {
				log.Println("unable to insert into db")
				message = "Comment not added"
			} else {
				message = "Comment added"
			}
		}

		http.Redirect(w, r, "/", http.StatusFound)

	}
}
