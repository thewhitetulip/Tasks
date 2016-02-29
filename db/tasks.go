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
	var tasks []types.Task
	var task types.Task
	var TaskCreated time.Time
	var context types.Context
	var getTasksql string
	var rows *sql.Rows

	comments, err := GetComments()

	if err != nil {
		return context, err
	}

	basicSQL := "select t.id, title, content, created_date, priority, c.name from task t, category c where c.id = t.cat_id"
	if status == "pending" && category == "" {
		getTasksql = basicSQL + " and finish_date is null and is_deleted='N' order by priority desc, created_date asc"
	} else if status == "deleted" {
		getTasksql = basicSQL + " and is_deleted='Y' order by priority desc, created_date asc"
	} else if status == "completed" {
		getTasksql = basicSQL + " and finish_date is not null order by priority desc, created_date asc"
	}

	if category != "" {
		status = category
		getTasksql = basicSQL + " and name = ?  and  t.is_deleted!='Y' and t.finish_date is null  order by priority desc, created_date asc, finish_date asc"
		rows, err = database.db.Query(getTasksql, category)

		if err != nil {
			log.Println("something went wrong while getting query")
		}
	} else {
		rows = database.query(getTasksql)
	}
	defer rows.Close()
	for rows.Next() {
		task = types.Task{}

		err = rows.Scan(&task.Id, &task.Title, &task.Content, &TaskCreated, &task.Priority, &task.Category)

		task.Content = string(md.Markdown([]byte(task.Content)))
		// TaskContent = strings.Replace(TaskContent, "\n", "<br>", -1)
		if err != nil {
			log.Println(err)
		}

		if comments[task.Id] != nil {
			task.Comments = comments[task.Id]
		}

		TaskCreated = TaskCreated.Local()
		if task.Priority != "1" { // if priority is not 1 then calculate, else why bother?
			CurrentTime := time.Now().Local()
			diff := CurrentTime.Sub(TaskCreated).Hours()
			if diff > 168 {
				task.IsOverdue = true // If one week then overdue by default
			}
		}
		task.Created = TaskCreated.Format("Jan 01 2006")

		tasks = append(tasks, task)
	}
	context = types.Context{Tasks: tasks, Navigation: status}
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
	var tasks []types.Task
	var task types.Task
	var TaskCreated time.Time
	var context types.Context

	comments, err := GetComments()
	if err != nil {
		log.Println("SearchTask: something went wrong in finding comments")
	}

	stmt := "select t.id, title, content, created_date, priority, c.name from task t, category c where c.id = t.cat_id and (title like '%" + query + "%' or content like '%" + query + "%') order by created_date desc"

	rows := database.query(stmt, query, query)

	for rows.Next() {
		err := rows.Scan(&task.Id, &task.Title, &task.Content, &TaskCreated, &task.Priority, &task.Category)
		if err != nil {
			log.Println(err)
		}

		if comments[task.Id] != nil {
			task.Comments = comments[task.Id]
		}

		task.Title = strings.Replace(task.Title, query, "<span class='highlight'>"+query+"</span>", -1)
		task.Content = strings.Replace(task.Content, query, "<span class='highlight'>"+query+"</span>", -1)
		task.Content = string(md.Markdown([]byte(task.Content)))

		TaskCreated = TaskCreated.Local()
		CurrentTime := time.Now().Local()
		week := TaskCreated.AddDate(0, 0, 7)

		if (week.String() < CurrentTime.String()) && (task.Priority != "1") {
			task.IsOverdue = true // If one week then overdue by default
		}
		task.Created = TaskCreated.Format("Jan 01 2006")

		tasks = append(tasks, task)
	}
	context = types.Context{Tasks: tasks, Search: query, Navigation: "search"}
	return context
}

//GetComments is used to get comments, all of them.
//We do not want 100 different pages to show tasks, we want to use as few pages as possible
//so we are going to populate everything on the damn home pages
func GetComments() (map[int][]types.Comment, error) {
	commentMap := make(map[int][]types.Comment)

	var taskID int
	var comment types.Comment
	var created time.Time

	stmt := "select id, taskID, content, created from comments;"
	rows := database.query(stmt)

	for rows.Next() {
		err := rows.Scan(&comment.ID, &taskID, &comment.Content, &created)
		if err != nil {
			return commentMap, err
		}
		// comment.Content = string(md.Markdown([]byte(comment.Content))) ## have to fix the <p> issue markdown support
		created = created.Local()
		comment.Created = created.Format("02 Jan 2006 15:04:05")
		commentMap[taskID] = append(commentMap[taskID], comment)
	}
	return commentMap, nil
}

//AddComments will be used to add comments in the database
func AddComments(id int, comment string) error {
	stmt := "insert into comments(taskID, content, created) values (?,?,datetime())"
	err := taskQuery(stmt, id, comment)

	if err != nil {
		return err
	}

	log.Println("added comment to task ID ", id)

	return nil
}
