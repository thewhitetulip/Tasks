package db

import (
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" //we want to use sqlite natively
	md "github.com/shurcooL/github_flavored_markdown"
	"github.com/thewhitetulip/Tasks/types"
)

var database Database
var err error

//Database encapsulates database
type Database struct {
	db *sql.DB
}

func (db Database) begin() (tx *sql.Tx) {
	tx, err := db.db.Begin()
	if err != nil {
		log.Println(err)
	}
	return tx
}

func (db Database) prepare(q string) (stmt *sql.Stmt) {
	stmt, err := db.db.Prepare(q)
	if err != nil {
		log.Println(err)
	}
	return stmt
}

func (db Database) query(q string, args ...interface{}) (rows *sql.Rows) {
	rows, err := db.db.Query(q, args...)
	if err != nil {
		log.Println(err)
	}
	return rows
}

func init() {
	database.db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Println(err)
	}
}

//Close function closes this database connection
func Close() {
	database.db.Close()
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

	rows := database.query(getTasksql)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&TaskID, &TaskTitle, &TaskContent, &TaskCreated, &TaskPriority)
		TaskContent = string(md.Markdown([]byte(TaskContent)))
		// TaskContent = strings.Replace(TaskContent, "\n", "<br>", -1)
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

	rows := database.query(getTasksql, id)
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
	err := taskQuery("update task set is_deleted='Y',last_modified_at=datetime() where id=?", id)
	return err
}

//CompleteTask  is used to mark tasks as complete
func CompleteTask(id int) error {
	err := taskQuery("update task set is_deleted='Y', finish_date=datetime(),last_modified_at=datetime() where id=?", id)
	return err
}

//DeleteAll is used to empty the trash
func DeleteAll() error {
	err := taskQuery("delete from task where is_deleted='Y'")
	return err
}

//RestoreTask is used to restore tasks from the Trash
func RestoreTask(id int) error {
	err := taskQuery("update task set is_deleted='N',last_modified_at=datetime() where id=?", id)
	return err
}

//RestoreTaskFromComplete is used to restore tasks from the Trash
func RestoreTaskFromComplete(id int) error {
	err := taskQuery("update task set finish_date=null,last_modified_at=datetime() where id=?", id)
	return err
}

//DeleteTask is used to delete the task from the database
func DeleteTask(id int) error {
	err := taskQuery("delete from task where id = ?", id)
	return err
}

//AddTask is used to add the task in the database
func AddTask(title, content string, taskPriority int) error {
	err := taskQuery("insert into task(title, content, priority, created_date, last_modified_at) values(?,?,?,datetime(), datetime())", title, content, taskPriority)
	return err
}

//UpdateTask is used to update the tasks in the database
func UpdateTask(id int, title string, content string) error {
	err := taskQuery("update task set title=?, content=? where id=?", title, content)
	return err
}

func taskQuery(sql string, args ...interface{}) error {
	SQL := database.prepare("update task set title=?, content=? where id=?")
	tx := database.begin()
	_, err = tx.Stmt(SQL).Exec(args...)
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

	rows := database.query(stmt, query, query)

	for rows.Next() {
		err := rows.Scan(&TaskID, &TaskTitle, &TaskContent, &TaskCreated)
		if err != nil {
			log.Println(err)
		}
		TaskTitle = strings.Replace(TaskTitle, query, "<span class='highlight'>"+query+"</span>", -1)
		TaskContent = strings.Replace(TaskContent, query, "<span class='highlight'>"+query+"</span>", -1)
		TaskContent = string(md.Markdown([]byte(TaskContent)))
		a := types.Task{Id: TaskID, Title: TaskTitle, Content: TaskContent, Created: TaskCreated.Format(time.UnixDate)[0:20]}
		task = append(task, a)
	}
	context = types.Context{Tasks: task, Search: query}
	return context
}
