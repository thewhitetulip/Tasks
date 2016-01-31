package views

import (
	"crypto/md5"
	"fmt"
	"github.com/thewhitetulip/Tasks/db"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

// UploadedFileHandler is used to handle the uploaded file related requests
func UploadedFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("into the handler")
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
		r.ParseForm()
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println(err)
		}

		taskPriority, priorityErr := strconv.Atoi(r.FormValue("priority"))
		if priorityErr != nil {
			log.Print(priorityErr)
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
		title := template.HTMLEscapeString(r.Form.Get("title"))
		content := template.HTMLEscapeString(r.Form.Get("content"))
		formToken := template.HTMLEscapeString(r.Form.Get("CSRFToken"))

		cookie, _ := r.Cookie("csrftoken")
		if formToken == cookie.Value {
			if handler != nil {
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

				filelink := "<br> <a href=/files/" + token + ">" + handler.Filename + "</a>"
				content = content + filelink

				fileTruth := db.AddFile(handler.Filename, token)
				if fileTruth != nil {
					message = "Error adding filename in db"
					log.Println("error adding task to db")
				}
			}

			taskTruth := db.AddTask(title, content, taskPriority)

			if taskTruth != nil {
				message = "Error adding task"
				log.Println("error adding task to db")
			} else {
				message = "Task added"
				log.Println("added task to db")
			}
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			log.Fatal("CSRF mismatch")
		}

	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
