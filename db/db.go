package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" //we want to use sqlite natively
	"github.com/thewhitetulip/Tasks/types"
	"log"
	"strings"
	"time"
)

var database *sql.DB
var err error

func init() {
	database, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Println(err)
	}
}

//Close function closes this database connection
func Close() {
	database.Close()
}

//GetTasks retrieves all the tasks depending on the
//status pending or trashed or completed
func GetTasks(status string) types.Context {
	var task []types.Task
	var context types.Context
	var TaskID int
	var TaskTitle string
	var TaskContent string
	var TaskCreated time.Time
	var TaskPriority string
	var getTasksql string

	basicSQL := "select id, title, content, created_date, priority from task "
	if status == "pending" {
		getTasksql = basicSQL + " where finish_date is null and is_deleted='N' order by priority desc, created_date asc"
	} else if status == "deleted" {
		getTasksql = basicSQL + " where is_deleted='Y' order by priority desc, created_date asc"
	} else if status == "completed" {
		getTasksql = basicSQL + " where finish_date is not null order by priority desc, created_date asc"
	}

	rows, err := database.Query(getTasksql)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&TaskID, &TaskTitle, &TaskContent, &TaskCreated, &TaskPriority)
		TaskContent = strings.Replace(TaskContent, "\n", "<br>", -1)
		if err != nil {
			log.Println(err)
		}
		TaskCreated = TaskCreated.Local()
		a := types.Task{Id: TaskID, Title: TaskTitle, Content: TaskContent, Created: TaskCreated.Format(time.UnixDate)[0:20], Priority: TaskPriority}
		task = append(task, a)
	}
	context = types.Context{Tasks: task, Navigation: status}
	return context
}

//GetTaskByID function gets the tasks from the ID passed to the function, used to populate EditTask
func GetTaskByID(id int) types.Context {
	var tasks []types.Task
	var task types.Task

	getTasksql := "select id, title, content, priority from task where id=?"

	rows, err := database.Query(getTasksql, id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&task.Id, &task.Title, &task.Content, &task.Priority)
		if err != nil {
			log.Println(err)
			//send email to respective people
		}
	}
	tasks = append(tasks, task)
	context := types.Context{Tasks: tasks, Navigation: "edit"}
	return context
}

//TrashTask is used to delete the task
func TrashTask(id int) error {
	trashSQL, err := database.Prepare("update task set is_deleted='Y',last_modified_at=datetime() where id=?")
	if err != nil {
		log.Println(err)
	}
	tx, err := database.Begin()
	if err != nil {
		log.Println(err)
	}
	_, err = tx.Stmt(trashSQL).Exec(id)
	if err != nil {
		log.Println("doing rollback")
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

//CompleteTask  is used to mark tasks as complete
func CompleteTask(id int) error {
	stmt, err := database.Prepare("update task set is_deleted='Y', finish_date=datetime(),last_modified_at=datetime() where id=?")
	if err != nil {
		log.Println(err)
	}
	tx, err := database.Begin()
	if err != nil {
		log.Println(err)
	}
	_, err = tx.Stmt(stmt).Exec(id)
	if err != nil {
		log.Println(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

//DeleteAll is used to empty the trash
func DeleteAll() error {
	stmt, err := database.Prepare("delete from task where is_deleted='Y'")
	if err != nil {
		log.Println(err)
	}
	tx, err := database.Begin()
	if err != nil {
		log.Println(err)
	}
	_, err = tx.Stmt(stmt).Exec()
	if err != nil {
		log.Println("doing rollback")
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

//RestoreTask is used to restore tasks from the Trash
func RestoreTask(id int) error {
	restoreSQL, err := database.Prepare("update task set is_deleted='N',last_modified_at=datetime() where id=?")
	if err != nil {
		log.Println(err)
	}
	tx, err := database.Begin()
	if err != nil {
		log.Println(err)
	}
	_, err = tx.Stmt(restoreSQL).Exec(id)
	if err != nil {
		log.Println("doing rollback")
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

//RestoreTask is used to restore tasks from the Trash
func RestoreTaskFromComplete(id int) error {
	restoreSQL, err := database.Prepare("update task set finish_date=null,last_modified_at=datetime() where id=?")
	if err != nil {
		log.Println(err)
	}
	tx, err := database.Begin()
	if err != nil {
		log.Println(err)
	}
	_, err = tx.Stmt(restoreSQL).Exec(id)
	if err != nil {
		log.Println("doing rollback")
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

//DeleteTask is used to delete the task from the database
func DeleteTask(id int) error {
	deleteSQL, err := database.Prepare("delete from task where id = ?")
	if err != nil {
		log.Println(err)
	}
	tx, err := database.Begin()
	if err != nil {
		log.Println(err)
	}
	_, err = tx.Stmt(deleteSQL).Exec(id)
	if err != nil {
		log.Println(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

//AddTask is used to add the task in the database
func AddTask(title, content string, taskPriority int) error {
	restoreSQL, err := database.Prepare("insert into task(title, content, priority, created_date, last_modified_at) values(?,?,?,datetime(), datetime())")
	if err != nil {
		log.Println(err)
	}
	tx, err := database.Begin()
	_, err = tx.Stmt(restoreSQL).Exec(title, content, taskPriority)
	if err != nil {
		log.Println(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

//UpdateTask is used to update the tasks in the database
func UpdateTask(id int, title string, content string) error {
	SQL, err := database.Prepare("update task set title=?, content=? where id=?")
	if err != nil {
		log.Println(err)
	}
	tx, err := database.Begin()

	if err != nil {
		log.Println(err)
	}
	_, err = tx.Stmt(SQL).Exec(title, content, id)
	if err != nil {
		log.Println(err)
		tx.Rollback()
	} else {
		log.Println(tx.Commit())
	}
	return err
}

//SearchTask is used to return the search results depending on the query
func SearchTask(query string) types.Context {
	stmt := "select id, title, content, created_date from task where title like '%" + query + "%' or content like '%" + query + "%'"
	var task []types.Task
	var TaskID int
	var TaskTitle string
	var TaskContent string
	var TaskCreated time.Time
	var context types.Context

	rows, err := database.Query(stmt, query, query)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		err := rows.Scan(&TaskID, &TaskTitle, &TaskContent, &TaskCreated)
		if err != nil {
			log.Println(err)
		}
		TaskTitle = strings.Replace(TaskTitle, query, "<span class='highlight'>"+query+"</span>", -1)
		TaskContent = strings.Replace(TaskContent, query, "<span class='highlight'>"+query+"</span>", -1)
		a := types.Task{Id: TaskID, Title: TaskTitle, Content: TaskContent, Created: TaskCreated.Format(time.UnixDate)[0:20]}
		task = append(task, a)
	}
	context = types.Context{Tasks: task, Search: query}
	return context
}
