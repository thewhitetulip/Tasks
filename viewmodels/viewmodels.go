package viewmodels

import (
	"github.com/thewhitetulip/Tasks/db"
	"github.com/thewhitetulip/Tasks/types"
)

func GetTasks(status string) []types.Task {
	return db.GetTasks(status)
}

func SearchTask(query string) []types.Task {
	return db.SearchTask(query)
}

func AddTask(title, content string) bool {
	err := db.AddTask(title, content)
	if err != nil {
		return false
	}
	return true
}

func TrashTask(id int) bool {
	err := db.TrashTask(id)
	if err != nil {
		return false
	}
	return true
}

func RestoreTask(id int) bool {
	err := db.RestoreTask(id)
	if err != nil {
		return false
	}
	return true
}

func DeleteTask(id int) bool {
	err := db.DeleteTask(id)
	if err != nil {
		return false
	}
	return true
}

func DeleteAll() bool {
	err := db.DeleteAll()
	if err != nil {
		return false
	}
	return true
}

func CompleteTask(id int) bool {
	err := db.CompleteTask(id)
	if err != nil {
		return false
	}
	return true
}

func GetTaskById(id int) types.Task {
	return db.GetTaskById(id)
}

func UpdateTask(id int, title string, content string) bool {
	err := db.UpdateTask(id, title, content)
	if err != nil {
		return false
	}
	return true
}

func Close() {
	db.Close()
}
