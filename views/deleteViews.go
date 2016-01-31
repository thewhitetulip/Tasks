package views

import (
	"github.com/thewhitetulip/Tasks/db"
	"log"
	"net/http"
	"strconv"
)


//TrashTaskFunc is used to populate the trash tasks
func TrashTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id, err := strconv.Atoi(r.URL.Path[len("/trash/"):])
		if err != nil {
			log.Println(err)
		} else {
			err = db.TrashTask(id)
			if err != nil {
				message = "Error trashing task"
			} else {
				message = "Task trashed"
			}
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/trash", http.StatusFound)
	}
}

//RestoreTaskFunc is used to restore task from trash, handles "/restore/" URL
func RestoreTaskFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id, err := strconv.Atoi(r.URL.Path[len("/restore/"):])
		if err != nil {
			log.Println(err)
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
		} else {
			task := db.GetTaskByID(id)
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
			db.DeleteAll()
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			id, err := strconv.Atoi(id)
			if err != nil {
				log.Println(err)
			} else {
				err = db.DeleteTask(id)
				if err != nil {
					message = "Error deleting task"
				} else {
					message = "Task deleted"
				}
				http.Redirect(w, r, "/", http.StatusFound)
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
		} else {
			err = db.RestoreTaskFromComplete(id)
			if err != nil {
				message = "Restore failed"
			} else {
				message = "Task restored"
			}
			http.Redirect(w, r, "/pending/", http.StatusFound)
		}
	} else {
		message = "Method not allowed"
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
