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
	http.HandleFunc("/", views.RequiresLogin(views.ShowAllTasksFunc))
	http.HandleFunc("/login/", views.LoginFunc)
	http.HandleFunc("/logout/", views.RequiresLogin(views.LogoutFunc))
	http.HandleFunc("/add-category/", views.RequiresLogin(views.AddCategoryFunc))
	http.HandleFunc("/add-comment/", views.RequiresLogin(views.AddCommentFunc))
	http.HandleFunc("/del-comment/", views.RequiresLogin(views.DeleteCommentFunc))
	http.HandleFunc("/del-category/", views.RequiresLogin(views.DeleteCategoryFunc))
	http.HandleFunc("/upd-category/", views.RequiresLogin(views.UpdateCategoryFunc))
	http.HandleFunc("/category/", views.RequiresLogin(views.ShowCategoryFunc))
	http.HandleFunc("/complete/", views.RequiresLogin(views.CompleteTaskFunc))
	http.HandleFunc("/delete/", views.RequiresLogin(views.DeleteTaskFunc))
	http.HandleFunc("/files/", views.RequiresLogin(views.UploadedFileHandler))
	http.HandleFunc("/deleted/", views.RequiresLogin(views.ShowTrashTaskFunc))
	http.HandleFunc("/trash/", views.RequiresLogin(views.TrashTaskFunc))
	http.HandleFunc("/edit/", views.RequiresLogin(views.EditTaskFunc))
	http.HandleFunc("/completed/", views.RequiresLogin(views.ShowCompleteTasksFunc))
	http.HandleFunc("/restore/", views.RequiresLogin(views.RestoreTaskFunc))
	http.HandleFunc("/incomplete/", views.RequiresLogin(views.RestoreFromCompleteFunc))
	http.HandleFunc("/add/", views.RequiresLogin(views.AddTaskFunc))
	http.HandleFunc("/update/", views.RequiresLogin(views.UpdateTaskFunc))
	http.HandleFunc("/search/", views.RequiresLogin(views.SearchTaskFunc))
	//http.HandleFunc("/static/", ServeStaticFunc)
	http.Handle("/static/", http.FileServer(http.Dir("public")))
	log.Println("running server on ", values.ServerPort)
	log.Fatal(http.ListenAndServe(values.ServerPort, nil))
}
