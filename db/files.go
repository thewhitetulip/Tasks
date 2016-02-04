package db

/*
stores the functions related to file IO
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
	stmt := "select c.name, count(*) from  category c left outer join task t  where c.id = t.cat_id   group by name    union     select name, 0  from category where name not in (select distinct name from task t join category c on t.cat_id = c.id)"
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

// GetCategoryById will return the ID of that category passed as args
// used while inserting tasks into the table
func GetCategoryById(category string) int {
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
