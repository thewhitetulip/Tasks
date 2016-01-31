package views

import (
	"bufio"
	"github.com/thewhitetulip/Tasks/db"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

var homeTemplate *template.Template
var deletedTemplate *template.Template
var completedTemplate *template.Template
var editTemplate *template.Template
var searchTemplate *template.Template
var templates *template.Template
var message string //message will store the message to be shown as notification
var err error

//ShowAllTasksFunc is used to handle the "/" URL which is the default ons
//TODO add http404 error
func ShowAllTasksFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		context := db.GetTasks("pending") //true when you want non deleted notes
		if message != "" {
			context.Message = message
		}
		context.CSRFToken = "abcd"
		message = ""
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "csrftoken", Value: "abcd", Expires: expiration}
		http.SetCookie(w, &cookie)
		homeTemplate.Execute(w, context)
	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

//ShowTrashTaskFunc is used to handle the "/trash" URL which is used to show the deleted tasks
func ShowTrashTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		context := db.GetTasks("deleted") //false when you want deleted notes
		if message != "" {
			context.Message = message
			message = ""
		}
		deletedTemplate.Execute(w, context)
	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

//ShowCompleteTasksFunc is used to populate the "/completed/" URL
func ShowCompleteTasksFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		context := db.GetTasks("completed") //false when you want finished notes
		completedTemplate.Execute(w, context)
	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/", http.StatusFound)
	}
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
	} else if strings.HasSuffix(path, ".js") {
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
