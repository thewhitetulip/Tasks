package views

import (
	"bufio"
	"fmt"
	"github.com/thewhitetulip/Tasks/db"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

var homeTemplate *template.Template
var deletedTemplate *template.Template
var completedTemplate *template.Template
var editTemplate *template.Template
var searchTemplate *template.Template
var templates *template.Template
var err error

//PopulateTemplates is used to parse all templates present in
//the templates folder
func PopulateTemplates() {
	var allFiles []string
	templatesDir := "./public/templates/"
	files, err := ioutil.ReadDir(templatesDir)
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, templatesDir+filename)
		}
	}

	if err != nil {
		fmt.Println(err)
	}
	templates, err = template.ParseFiles(allFiles...)
	homeTemplate = templates.Lookup("home.html")
	deletedTemplate = templates.Lookup("deleted.html")

	editTemplate = templates.Lookup("edit.html")
	searchTemplate = templates.Lookup("search.html")
	completedTemplate = templates.Lookup("completed.html")

}

//ShowAllTasksFunc is used to handle the "/" URL which is the default ons
//TODO add http404 error
func ShowAllTasksFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		context := db.GetTasks("pending") //true when you want non deleted notes
		homeTemplate.Execute(w, context)
	}
}

//ShowTrashTaskFunc is used to handle the "/trash" URL which is used to show the deleted tasks
func ShowTrashTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		context := db.GetTasks("trashed") //false when you want deleted notes
		deletedTemplate.Execute(w, context)
	}
}

//SearchTaskFunc is used to handle the /search/ url, handles the search function
func SearchTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		query := r.Form.Get("query")
		context := db.SearchTask(query)
		searchTemplate.Execute(w, context)
	} else {

	}

}

//AddTaskFunc is used to handle the addition of new task, "/add" URL
func AddTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		title := r.Form.Get("title")
		content := r.Form.Get("content")
		truth := db.AddTask(title, content)
		if truth != nil {
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			fmt.Println(err)
		}
	}
}

//ShowCompleteTasksFunc is used to populate the "/completed/" URL
func ShowCompleteTasksFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		context := db.GetTasks("complete") //false when you want finished notes
		completedTemplate.Execute(w, context)
	}
}

//EditTaskFunc is used to edit tasks, handles "/edit/" URL
func EditTaskFunc(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/edit/"):])
	if err != nil {
		fmt.Println(err)
	} else {
		task := db.GetTaskById(id)
		editTemplate.Execute(w, task)
	}
}

//CompleteTaskFunc is used to show the complete tasks, handles "/completed/" url
func CompleteTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id, err := strconv.Atoi(r.URL.Path[len("/complete/"):])
		if err != nil {
			fmt.Println(err)
		} else {
			err := db.CompleteTask(id)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("redirecting to home")
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

//DeleteTaskFunc is used to
func DeleteTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Path[len("/delete/"):]
		if id == "all" {
			db.DeleteAll()
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			id, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println(err)
			} else {
				db.DeleteTask(id)
				http.Redirect(w, r, "/deleted/", http.StatusFound)
			}
		}
	}
}

//TrashTaskFunc is used to populate the "/trash/" URL
func TrashTaskFunc(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/trash/"):])
	if err != nil {
		fmt.Println(err)
	} else {
		db.TrashTask(id)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

//RestoreTaskFunc is used to restore task from trash, handles "/restore/" URL
func RestoreTaskFunc(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/restore/"):])
	if err != nil {
		fmt.Println(err)
	} else {
		db.RestoreTask(id)
		http.Redirect(w, r, "/deleted/", http.StatusFound)
	}
}

//UpdateTaskFunc is used to update a task, handes "/update/" URL
func UpdateTaskFunc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		fmt.Println(err)
	}
	title := r.Form.Get("title")
	content := r.Form.Get("content")
	db.UpdateTask(id, title, content)
	http.Redirect(w, r, "/", http.StatusFound)
}

//ServeStaticFunc is used to serve static files
//TODO: replace this with the http.FileServer
func ServeStaticFunc(w http.ResponseWriter, r *http.Request) {
	path := "./public" + r.URL.Path
	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "application/javascript"
	} else {
		contentType = "plain/text"
	}

	f, err := os.Open(path)

	if err == nil {
		defer f.Close()
		w.Header().Add("Content Type", contentType)

		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}
