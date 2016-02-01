package views

/*
Holds the delete related view handlers
*/

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/thewhitetulip/Tasks/db"
)

//TrashTaskFunc is used to populate the trash tasks
func TrashTaskFunc(w http.ResponseWriter, r *http.Request) {
	//for best UX we want the user to be returned to the page making
	//the delete transaction, we use the r.Referer() function to get the link
	var redirectUrl string
	redirect := strings.Split(r.Referer(), "/")
	index := len(redirect) - 1
	if len(redirect) == 4 {
		redirectUrl = "/"
	} else {
		redirectUrl = redirect[index]
	}

	if r.Method == "GET" {
		id, err := strconv.Atoi(r.URL.Path[len("/trash/"):])
		if err != nil {
			log.Println("TrashTaskFunc", err)
			http.Redirect(w, r, redirectUrl, http.StatusBadRequest)
		} else {
			err = db.TrashTask(id)
			if err != nil {
				message = "Error trashing task"
			} else {
				message = "Task trashed"
			}
			http.Redirect(w, r, redirectUrl, http.StatusFound)
		}
	} else {
		message = "Method not allowed"
		http.Redirect(w, r, redirectUrl, http.StatusFound)
	}
}

//RestoreTaskFunc is used to restore task from trash, handles "/restore/" URL
func RestoreTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id, err := strconv.Atoi(r.URL.Path[len("/restore/"):])
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/deleted", http.StatusBadRequest)
		} else {
			err = db.RestoreTask(id)
			if err != nil {
				message = "Restore failed"
			} else {
				message = "Task restored"
			}
			http.Redirect(w, r, "/deleted/", http.StatusFound)
		}
	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/", http.StatusFound)
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
			task, err := db.GetTaskByID(id)
			if err != nil {
				task.Message = "Error fetching Tasks"
			}
			editTemplate.Execute(w, task)
		}
	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

//DeleteTaskFunc is used to delete a task, trash = move to recycle bin, delete = permanent delete
func DeleteTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id := r.URL.Path[len("/delete/"):]
		if id == "all" {
			err := db.DeleteAll()
			if err != nil {
				message = "Error deleting tasks"
				http.Redirect(w, r, "/", http.StatusInternalServerError)
			}
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			id, err := strconv.Atoi(id)
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/", http.StatusBadRequest)
			} else {
				err = db.DeleteTask(id)
				if err != nil {
					message = "Error deleting task"
				} else {
					message = "Task deleted"
				}
				http.Redirect(w, r, "/deleted", http.StatusFound)
			}
		}
	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

//RestoreFromCompleteFunc restores the task from complete to pending
func RestoreFromCompleteFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id, err := strconv.Atoi(r.URL.Path[len("/incomplete/"):])
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/completed", http.StatusBadRequest)
		} else {
			err = db.RestoreTaskFromComplete(id)
			if err != nil {
				message = "Restore failed"
			} else {
				message = "Task restored"
			}
			http.Redirect(w, r, "/completed", http.StatusFound)
		}
	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/completed", http.StatusFound)
	}
}
