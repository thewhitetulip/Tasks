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
func AddFile(fileName, token string) error {
	SQL := database.prepare("insert into files values(?,?)")
	tx := database.begin()
	_, err = tx.Stmt(SQL).Exec(fileName, token)
	if err != nil {
		log.Println(err)
		tx.Rollback()
	} else {
		log.Println(tx.Commit())
	}
	return err
}

// GetFileName is used to fetch the name according to the md5 checksum from the db
func GetFileName(token string) (string, error) {
	sql := "select name from files where autoName=?"
	var fileName string
	rows := database.query(sql, fileName)
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
func GetCategories() []types.CategoryCount {
	stmt := "select c.name, count(*) from  category c left outer join task t  where c.id = t.cat_id and t.is_deleted='N' and t.finish_date is null   group by name    union     select name, 0  from category where name not in (select distinct name from task t join category c on t.cat_id = c.id and is_deleted!='Y'and t.finish_date is null)"
	rows := database.query(stmt)
	var categories []types.CategoryCount
	var category types.CategoryCount

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
func AddCategory(category string) error {
	err := taskQuery("insert into category(name) values(?)", category)
	return err
}

// GetCategoryByName will return the ID of that category passed as args
// used while inserting tasks into the table
func GetCategoryByName(category string) int {
	stmt := "select id from category where name=?"
	rows := database.query(stmt, category)
	var categoryID int

	for rows.Next() {
		err := rows.Scan(&categoryID)
		if err != nil {
			log.Println(err)
		}
	}
	return categoryID
}

//DeleteCategoryByName will be used to delete a category from the category page
func DeleteCategoryByName(category string) error {
	//first we delete entries from task and then from category
	categoryID := GetCategoryByName(category)
	query := "update task set cat_id = null where id =?"
	err := taskQuery(query, categoryID)
	if err == nil {
		err = taskQuery("delete from category where id=?", categoryID)
		if err != nil {
			return err
		}
	}
	return err
}

//UpdateCategoryByName will be used to delete a category from the category page
func UpdateCategoryByName(oldName, newName string) error {
	query := "update category set name = ? where name=?"
	log.Println(query)
	err := taskQuery(query, newName, oldName)
	return err
}

//DeleteCommentByID will actually delete the comment from db
func DeleteCommentByID(id int) error {
	query := "delete from comments where id=?"
	err := taskQuery(query, id)
	return err
}
