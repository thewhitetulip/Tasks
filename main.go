package main

/**
 * This is the main file for the Task application
 * License: MIT
 **/
import (
	"log"
	"net/http"

	"github.com/thewhitetulip/Tasks/config"
	"github.com/thewhitetulip/Tasks/views"
)

func main() {
	values := config.ReadConfig("config.json")
	views.PopulateTemplates()
	http.HandleFunc("/", views.ShowAllTasksFunc)
	http.HandleFunc("/complete/", views.CompleteTaskFunc)
	//delete permanently deletes from db
	http.HandleFunc("/delete/", views.DeleteTaskFunc)
	http.HandleFunc("/files/", views.UploadedFileHandler)
	http.HandleFunc("/deleted/", views.ShowTrashTaskFunc)
	//trash moves to recycle bin
	http.HandleFunc("/trash/", views.TrashTaskFunc)
	http.HandleFunc("/edit/", views.EditTaskFunc)
	http.HandleFunc("/completed/", views.ShowCompleteTasksFunc)
	http.HandleFunc("/restore/", views.RestoreTaskFunc)
	http.HandleFunc("/incomplete/", views.RestoreFromCompleteFunc)
	http.HandleFunc("/add/", views.AddTaskFunc)
	http.HandleFunc("/update/", views.UpdateTaskFunc)
	http.HandleFunc("/search/", views.SearchTaskFunc)
	//http.HandleFunc("/static/", ServeStaticFunc)
	http.Handle("/static/", http.FileServer(http.Dir("public")))
	log.Println("running server on ", values.ServerPort)
	log.Fatal(http.ListenAndServe(values.ServerPort, nil))
}
