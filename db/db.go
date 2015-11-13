package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/thewhitetulip/task/types"
	"strings"
)

var database *sql.DB
var err error
var tag string
var name string
var tags string

func init() {
	database, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		fmt.Println(err)
	}
}

func Close() {
	database.Close()
}

func GetTasks(deleted bool) []types.Task{
	var task []types.Task
	var TaskId int
	var TaskTitle string
	var TaskContent string
	var getTasksql string
	if deleted == true{
		getTasksql = "select id, title, content from task where is_deleted!='Y' order by created_date asc"
	} else {
		getTasksql = "select id, title, content from task where is_deleted='Y' order by created_date asc"
	}

	rows, err := database.Query(getTasksql)
	if err!=nil{
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next(){
		err := rows.Scan(&TaskId, &TaskTitle, &TaskContent)
		if err!=nil{
			fmt.Println(err)
		}
		a := types.Task{Id:TaskId, Title:TaskTitle, Content:TaskContent}
		task = append(task, a)
	}

	return task
}

func GetTaskById(id int) types.Task{
	var task types.Task
	var TaskId int
	var TaskTitle string
	var TaskContent string
	getTasksql := "select id, title, content from task where id=?"

	rows, err := database.Query(getTasksql, id)
	if err!=nil{
		fmt.Println(err)
	}
	defer rows.Close()
	if rows.Next(){
		err := rows.Scan(&TaskId, &TaskTitle, &TaskContent)
		if err!=nil{
			fmt.Println(err)
		}
		task = types.Task{Id:TaskId, Title:TaskTitle, Content:TaskContent}
	}
	return task
}

func ArchiveTask(id int) error {
	stmt, err := database.Prepare("update task set is_deleted='Y', finish_date=datetime(),last_modified_at=datetime() where id=?")
	if err!=nil{
		fmt.Println(err)
	}
	tx, err := database.Begin()
	if err!=nil{
		fmt.Println(err)
	}
	_, err = tx.Stmt(stmt).Exec(id)
	if err!=nil{
		fmt.Println(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

func DeleteAll() error {
	stmt, err := database.Prepare("delete from task where is_deleted='Y'")
	if err!=nil{
		fmt.Println(err)
	}
	tx, err := database.Begin()
	if err!=nil{
		fmt.Println(err)
	}
	_, err = tx.Stmt(stmt).Exec()
	if err!=nil{
		fmt.Println("doing rollback")
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

func RestoreTask(id int) error {
	restoreSql, err := database.Prepare("update task set is_deleted='N',last_modified_at=datetime() where id=?")
	if err!=nil{
		fmt.Println(err)
	}
	tx, err := database.Begin()
	if err!=nil{
		fmt.Println(err)
	}
	_, err = tx.Stmt(restoreSql).Exec(id)
	if err!=nil{
		fmt.Println("doing rollback")
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

func DeleteTask(id int) error {
	deleteSQL, err := database.Prepare("delete from task where id = ?")
	if err!=nil{
		fmt.Println(err)
	}
	tx, err := database.Begin()
	if err!=nil{
		fmt.Println(err)
	}
	_, err = tx.Stmt(deleteSQL).Exec(id)
	if err!=nil{
		fmt.Println(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

func AddTask(title, content string) error {
	restoreSql, err := database.Prepare("insert into task(title, content, created_date, last_modified_at) values(?,?,datetime(), datetime())")
	if err!=nil{
		fmt.Println(err)
	}
	tx, err := database.Begin()
	_, err = tx.Stmt(restoreSql).Exec(title, content)
	if err!=nil{
		fmt.Println(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

func UpdateTask(id int, title string, content string) error{
	Sql, err := database.Prepare("update task set title=?, content=? where id=?")
	if err!=nil{
		fmt.Println(err)
	}
	tx, err := database.Begin()

	if err!=nil{
		fmt.Println(err)
	}
	_, err = tx.Stmt(Sql).Exec(title, content, id)
	if err!=nil{
		fmt.Println(err)
		tx.Rollback()
	} else {
		fmt.Println(tx.Commit())
	}
	return err
}

func SearchTask(query string) []types.Task{
	stmt := "select id, title, content from task where title like '%"+query+"%' or content like '%"+query+"%'"
	var task []types.Task
	var TaskId int
	var TaskTitle string
	var TaskContent string

	rows, err := database.Query(stmt, query, query)
	if err!=nil{
		fmt.Println(err)
	}
	for rows.Next(){
		err := rows.Scan(&TaskId, &TaskTitle, &TaskContent)
		if err!=nil{
			fmt.Println(err)
		}
		TaskTitle = strings.Replace(TaskTitle, query, "<span class='highlight'>"+query+"</span>", -1)
		TaskContent = strings.Replace(TaskContent, query, "<span class='highlight'>"+query+"</span>", -1)
		a := types.Task{Id:TaskId, Title:TaskTitle, Content:TaskContent}
		task = append(task, a)
	}
	return task
}
