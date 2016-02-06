package views

/*Holds the fetch related view handlers*/

import (
	"net/http"
	"text/template"
	"time"

	"github.com/thewhitetulip/Tasks/db"
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
		context, err := db.GetTasks("pending", "")
		categories := db.GetCategories()
		if err != nil {
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
		if message != "" {
			context.Message = message
		}
		context.CSRFToken = "abcd"
		context.Categories = categories
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
		context, err := db.GetTasks("deleted", "")
		categories := db.GetCategories()
		context.Categories = categories
		if err != nil {
			http.Redirect(w, r, "/trash", http.StatusInternalServerError)
		}
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
		context, err := db.GetTasks("completed", "")
		categories := db.GetCategories()
		context.Categories = categories
		if err != nil {
			http.Redirect(w, r, "/completed", http.StatusInternalServerError)
		}
		completedTemplate.Execute(w, context)
	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

//ShowCategoryFunc will populate the /category/<id> URL which shows all the tasks related
// to that particular category
func ShowCategoryFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		category := r.URL.Path[len("/category/"):]
		context, err := db.GetTasks("", category)
		categories := db.GetCategories()

		if err != nil {
			http.Redirect(w, r, "/", http.StatusInternalServerError)
		}
		if message != "" {
			context.Message = message
		}
		context.CSRFToken = "abcd"
		context.Categories = categories
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
