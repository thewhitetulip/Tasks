package db

/*
Stores the database functions related to tasks like
GetTaskByID(id int)
GetTasks(status string)
DeleteAll()
*/

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
		return nil
	}
	return tx
}

func (db Database) prepare(q string) (stmt *sql.Stmt) {
	stmt, err := db.db.Prepare(q)
	if err != nil {
		log.Println(err)
		return nil
	}
	return stmt
}

func (db Database) query(q string, args ...interface{}) (rows *sql.Rows) {
	rows, err := db.db.Query(q, args...)
	if err != nil {
		log.Println(err)
		return nil
	}
	return rows
}

func init() {
	database.db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
}

//Close function closes this database connection
func Close() {
	database.db.Close()
}

//GetTasks retrieves all the tasks depending on the
//status pending or trashed or completed
func GetTasks(status, category string) (types.Context, error) {
	var task []types.Task
	var context types.Context
	var TaskID int
	var TaskTitle string
	var TaskContent string
	var TaskCreated time.Time
	var TaskPriority string
	var getTasksql string
	var rows *sql.Rows

	basicSQL := "select id, title, content, created_date, priority from task t"
	if status == "pending" && category == "" {
		getTasksql = basicSQL + " where finish_date is null and is_deleted='N' order by priority desc, created_date asc"
	} else if status == "deleted" {
		getTasksql = basicSQL + " where is_deleted='Y' order by priority desc, created_date asc"
	} else if status == "completed" {
		getTasksql = basicSQL + " where finish_date is not null order by priority desc, created_date asc"
	}

	if category != "" {
		status = category
		getTasksql = "select t.id, title, content, created_date, priority from task t, category c where c.id = t.cat_id and name = ?  and  is_deleted!='Y'  order by priority desc, created_date asc, finish_date asc"
		rows, err = database.db.Query(getTasksql, category)
		if err != nil {
			log.Println("something went wrong while getting query")
		}
	} else {
		rows = database.query(getTasksql)
	}
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
	return context, nil
}

//GetTaskByID function gets the tasks from the ID passed to the function, used to populate EditTask
func GetTaskByID(id int) (types.Context, error) {
	var tasks []types.Task
	var task types.Task

	getTasksql := "select t.id, t.title, t.content, t.priority, c.name from task t left outer join category c  where c.id = t.cat_id and t.id=?"

	rows := database.query(getTasksql, id)
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&task.Id, &task.Title, &task.Content, &task.Priority, &task.Category)
		if err != nil {
			log.Println(err)
			//send email to respective people
		}
	}
	tasks = append(tasks, task)
	context := types.Context{Tasks: tasks, Navigation: "edit"}
	return context, nil
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
func AddTask(title, content, category string, taskPriority int) error {
	var err error
	if category == "" {
		err = taskQuery("insert into task(title, content, priority, created_date, last_modified_at) values(?,?,?,datetime(), datetime())", title, content, taskPriority)
	} else {
		categoryID := GetCategoryByName(category)
		err = taskQuery("insert into task(title, content, priority, created_date, last_modified_at, cat_id) values(?,?,?,datetime(), datetime(), ?)", title, content, taskPriority, categoryID)
	}
	return err
}

//GetCategoryIdByName will return the category ID for the category, used in the edit task
//function where we need to be able to update the categoryID of the task
func GetCategoryIdByName(category string) int {
	var categoryID int
	getTasksql := "select id from category where name=?"

	rows := database.query(getTasksql, category)
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&categoryID)
		if err != nil {
			log.Println(err)
			//send email to respective people
		}
	}

	return categoryID
}

//UpdateTask is used to update the tasks in the database
func UpdateTask(id int, title, content, category string, priority int) error {
	categoryID := GetCategoryIdByName(category)
	err := taskQuery("update task set title=?, content=?, cat_id=?, priority = ? where id=?", title, content, categoryID, priority, id)
	return err
}

//taskQuery encapsulates running multiple queries which don't do much things
func taskQuery(sql string, args ...interface{}) error {
	SQL := database.prepare(sql)
	tx := database.begin()
	_, err = tx.Stmt(SQL).Exec(args...)
	if err != nil {
		log.Println(err)
		tx.Rollback()
	} else {
		tx.Commit()
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
