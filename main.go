package main

/**
 * This is the main file for the Task application
 * License: MIT
 **/
import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/thewhitetulip/task/viewmodels"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

var homeTemplate *template.Template
var deletedTemplate *template.Template
var editTemplate *template.Template
var searchTemplate *template.Template
var err error

func main() {
	defer viewmodels.Close()
	homeTemplate, err = template.ParseFiles("./templates/home.gtpl")
	if err != nil {
		fmt.Println(err)
	}
	deletedTemplate, err = template.ParseFiles("./templates/deleted.gtpl")
	if err != nil {
		fmt.Println(err)
	}

	editTemplate, err = template.ParseFiles("./templates/edit.gtpl")
	if err != nil {
		fmt.Println(err)
	}
	searchTemplate, err = template.ParseFiles("./templates/search.gtpl")
	if err != nil {
		fmt.Println(err)
	}
	router := httprouter.New()
	router.GET("/", ShowAllTasks)
	router.GET("/archive/:id", ArchiveTask)
	router.GET("/delete/:id", DeleteTask)
	router.GET("/edit/:id", EditTask)
	router.GET("/trash/", ShowTrashTask)
	router.GET("/restore/:id", RestoreTask)
	router.POST("/add/", AddTask)
	router.POST("/update/", UpdateTask)
	router.POST("/search/", SearchTask)
	router.NotFound = http.FileServer(http.Dir("public"))
	fmt.Println("running on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func ShowAllTasks(w http.ResponseWriter, r *http.Request, parm httprouter.Params) {
	context := viewmodels.GetTasks(true) //true when you want non deleted notes
	homeTemplate.Execute(w, context)
}

func SearchTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	r.ParseForm()
	query := r.Form.Get("query")
	context := viewmodels.SearchTask(query)
	searchTemplate.Execute(w, context)
}

func AddTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	title := r.Form.Get("title")
	content := r.Form.Get("content")
	truth := viewmodels.AddTask(title, content)
	if truth == true {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func ShowTrashTask(w http.ResponseWriter, r *http.Request, parm httprouter.Params) {
	context := viewmodels.GetTasks(false) //false when you want finished notes
	deletedTemplate.Execute(w, context)
}

func EditTask(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	id, err := strconv.Atoi(param.ByName("id"))
	if err!=nil{
		fmt.Println(err)
	} else {
		task := viewmodels.GetTaskById(id)
		editTemplate.Execute(w, task)
	}
}

func ArchiveTask(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	id, err := strconv.Atoi(param.ByName("id"))
	if err != nil {
		fmt.Println(err)
	} else {
		viewmodels.ArchiveTask(id)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	id := param.ByName("id")
	if id == "all"{
		viewmodels.DeleteAll()
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println(err)
		} else {
			viewmodels.DeleteTask(id)
			http.Redirect(w, r, "/trash/", http.StatusFound)
		}
	}
}

func RestoreTask(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	id, err := strconv.Atoi(param.ByName("id"))
	if err != nil {
		fmt.Println(err)
	} else {
		viewmodels.RestoreTask(id)
		http.Redirect(w, r, "/trash/", http.StatusFound)
	}
}

func UpdateTask(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	r.ParseForm()
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err!=nil{
		fmt.Println(err)
	}
	title := r.Form.Get("title")
	content := r.Form.Get("content")
	viewmodels.UpdateTask(id, title, content)
	http.Redirect(w, r, "/", http.StatusFound)
}
