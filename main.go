package main

/**
 * This is the main file for the Task application
 * License: MIT
 **/
import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/thewhitetulip/Tasks/config"
	"github.com/thewhitetulip/Tasks/views"
)

func main() {
	values, err := config.ReadConfig("config.json")
	var port *string

	if err != nil {
		port = flag.String("port", "", "IP address")
		flag.Parse()

		//User is expected to give :8080 like input, if they give 8080
		//we'll append the required ':'
		if !strings.HasPrefix(*port, ":") {
			*port = ":" + *port
			log.Println("port is " + *port)
		}

		values.ServerPort = *port
	}

	views.PopulateTemplates()

	//Login logout
	http.HandleFunc("/login/", views.LoginFunc)
	http.HandleFunc("/logout/", views.RequiresLogin(views.LogoutFunc))
	http.HandleFunc("/signup/", views.SignUpFunc)

	http.HandleFunc("/add-category/", views.RequiresLogin(views.AddCategoryFunc))
	http.HandleFunc("/add-comment/", views.RequiresLogin(views.AddCommentFunc))
	http.HandleFunc("/add/", views.RequiresLogin(views.AddTaskFunc))

	//these handlers are used to delete
	http.HandleFunc("/del-comment/", views.RequiresLogin(views.DeleteCommentFunc))
	http.HandleFunc("/del-category/", views.RequiresLogin(views.DeleteCategoryFunc))
	http.HandleFunc("/delete/", views.RequiresLogin(views.DeleteTaskFunc))

	//these handlers update
	http.HandleFunc("/upd-category/", views.RequiresLogin(views.UpdateCategoryFunc))
	http.HandleFunc("/update/", views.RequiresLogin(views.UpdateTaskFunc))

	//these handlers are used for restoring tasks
	http.HandleFunc("/incomplete/", views.RequiresLogin(views.RestoreFromCompleteFunc))
	http.HandleFunc("/restore/", views.RequiresLogin(views.RestoreTaskFunc))

	//these handlers fetch set of tasks
	http.HandleFunc("/", views.RequiresLogin(views.ShowAllTasksFunc))
	http.HandleFunc("/category/", views.RequiresLogin(views.ShowCategoryFunc))
	http.HandleFunc("/deleted/", views.RequiresLogin(views.ShowTrashTaskFunc))
	http.HandleFunc("/completed/", views.RequiresLogin(views.ShowCompleteTasksFunc))

	//these handlers perform action like delete, mark as complete etc
	http.HandleFunc("/complete/", views.RequiresLogin(views.CompleteTaskFunc))
	http.HandleFunc("/files/", views.RequiresLogin(views.UploadedFileHandler))
	http.HandleFunc("/trash/", views.RequiresLogin(views.TrashTaskFunc))
	http.HandleFunc("/edit/", views.RequiresLogin(views.EditTaskFunc))
	http.HandleFunc("/search/", views.RequiresLogin(views.SearchTaskFunc))

	http.Handle("/static/", http.FileServer(http.Dir("public")))

	// http.HandleFunc("/api/get-task/", views.GetTasksFuncAPI)
	// http.HandleFunc("/api/get-deleted-task/", views.GetDeletedTaskFuncAPI)
	// http.HandleFunc("/api/add-task/", views.AddTaskFuncAPI)
	// http.HandleFunc("/api/update-task/", views.UpdateTaskFuncAPI)
	// http.HandleFunc("/api/delete-task/", views.DeleteTaskFuncAPI)

	// http.HandleFunc("/api/get-token/", views.GetTokenHandler)
	// http.HandleFunc("/api/get-category/", views.GetCategoryFuncAPI)
	// http.HandleFunc("/api/add-category/", views.AddCategoryFuncAPI)
	// http.HandleFunc("/api/update-category/", views.UpdateCategoryFuncAPI)
	// http.HandleFunc("/api/delete-category/", views.DeleteCategoryFuncAPI)

	log.Println("running server on ", values.ServerPort)
	log.Fatal(http.ListenAndServe(values.ServerPort, nil))
}
