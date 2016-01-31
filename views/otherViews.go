package views

import (
	"github.com/thewhitetulip/Tasks/db"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

//PopulateTemplates is used to parse all templates present in
//the templates folder
func PopulateTemplates() {
	var allFiles []string
	templatesDir := "./public/templates/"
	files, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		log.Println(err)
		os.Exit(1) // No point in running app if templates aren't read
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, templatesDir+filename)
		}
	}

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	templates, err = template.ParseFiles(allFiles...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	homeTemplate = templates.Lookup("home.html")
	deletedTemplate = templates.Lookup("deleted.html")

	editTemplate = templates.Lookup("edit.html")
	searchTemplate = templates.Lookup("search.html")
	completedTemplate = templates.Lookup("completed.html")

}

//CompleteTaskFunc is used to show the complete tasks, handles "/completed/" url
func CompleteTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id, err := strconv.Atoi(r.URL.Path[len("/complete/"):])
		if err != nil {
			log.Println(err)
		} else {
			err = db.CompleteTask(id)
			if err != nil {
				message = "Complete task failed"
			} else {
				message = "Task marked complete"
			}
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/", http.StatusFound)
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
		message = "Method not allowed"
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

//UpdateTaskFunc is used to update a task, handes "/update/" URL
func UpdateTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		id, err := strconv.Atoi(r.Form.Get("id"))
		if err != nil {
			log.Println(err)
		}
		title := r.Form.Get("title")
		content := r.Form.Get("content")
		err = db.UpdateTask(id, title, content)
		if err != nil {
			message = "Error updating task"
		} else {
			message = "Task updated"
		}
		http.Redirect(w, r, "/", http.StatusFound)

	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
