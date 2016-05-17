package db

/*
stores the functions related to file IO and category
*/
import (
	"log"

	"github.com/thewhitetulip/Tasks/types"
)

// AddFile is used to add the md5 of a file name which is uploaded to our application
// this will enable us to randomize the URL without worrying about the file names
func AddFile(fileName, token, username string) error {
	userID, err := GetUserID(username)
	if err != nil {
		return err
	}
	err = taskQuery("insert into files values(?,?,?,datetime())", fileName, token, userID)
	return err
}

// GetFileName is used to fetch the name according to the md5 checksum from the db
func GetFileName(token string) (string, error) {
	sql := "select name from files where autoName=?"
	var fileName string
	rows := database.query(sql, fileName)
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&fileName)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}
	if err != nil {
		return "", err
	}

	return fileName, nil
}

//GetCategories will return the list of categories to be
//rendered in the template
func GetCategories(username string) []types.CategoryCount {
	userID, err := GetUserID(username)
	if err != nil {
		return nil
	}
	stmt := "select 'UNCATEGORIZED' as name, count(1) from task where cat_id=0 union  select c.name, count(*) from   category c left outer join task t  join status s on  c.id = t.cat_id and t.task_status_id=s.id where s.status!='DELETED' and c.user_id=?   group by name    union     select name, 0  from category c, user u where c.user_id=? and name not in (select distinct name from task t join category c join status s on s.id = t.task_status_id and t.cat_id = c.id and s.status!='DELETED' and c.user_id=?)"
	rows := database.query(stmt, userID, userID, userID)
	var categories []types.CategoryCount
	var category types.CategoryCount

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&category.Name, &category.Count)
		if err != nil {
			log.Println(err)
		}
		categories = append(categories, category)
	}
	return categories
}

//AddCategory is used to add the task in the database
func AddCategory(username, category string) error {
	userID, err := GetUserID(username)
	if err != nil {
		return nil
	}
	log.Println("executing query to add category")
	err = taskQuery("insert into category(name, user_id) values(?,?)", category, userID)
	return err
}

// GetCategoryByName will return the ID of that category passed as args
// used while inserting tasks into the table
func GetCategoryByName(username, category string) int {
	stmt := "select id from category where name=? and user_id = (select id from user where username=?)"
	rows := database.query(stmt, category, username)
	var categoryID int
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&categoryID)
		if err != nil {
			log.Println(err)
		}
	}
	return categoryID
}

//DeleteCategoryByName will be used to delete a category from the category page
func DeleteCategoryByName(username, category string) error {
	//first we delete entries from task and then from category
	categoryID := GetCategoryByName(username, category)
	userID, err := GetUserID(username)
	if err != nil {
		return err
	}
	query := "update task set cat_id = null where id =? and user_id = ?"
	err = taskQuery(query, categoryID, userID)
	if err == nil {
		err = taskQuery("delete from category where id=? and user_id=?", categoryID, userID)
		if err != nil {
			return err
		}
	}
	return err
}

//UpdateCategoryByName will be used to delete a category from the category page
func UpdateCategoryByName(username, oldName, newName string) error {
	userID, err := GetUserID(username)
	if err != nil {
		return err
	}
	query := "update category set name = ? where name=? and user_id=?"
	log.Println(query)
	err = taskQuery(query, newName, oldName, userID)
	return err
}

//DeleteCommentByID will actually delete the comment from db
func DeleteCommentByID(username string, id int) error {
	userID, err := GetUserID(username)
	if err != nil {
		return err
	}
	query := "delete from comments where id=? and user_id = ?"
	err = taskQuery(query, id, userID)
	return err
}
