package db

import (
	"log"
)

// AddFile is used to add the md5 of a file name which is uploaded to our application
// this will enable us to randomize the URL without worrying about the file names
func AddFile(fileName, token string) error {
	SQL, err := database.Prepare("insert into files values(?,?)")
	if err != nil {
		log.Println(err)
	}
	tx, err := database.Begin()

	if err != nil {
		log.Println(err)
	}
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
	rows, err := database.Query(sql, fileName)
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
