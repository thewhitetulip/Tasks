package db

import "log"

//CreateUser will create a new user, take as input the parameters and
//insert it into database
func CreateUser(username, password, email string) error {
	err := taskQuery("insert into user(username, password, email) values(?,?,?)", username, password, email)
	return err
}

//ValidUser will check if the user exists in db and if exists if the username password
//combination is valid
func ValidUser(username, password string) bool {
	var passwordFromDB string
	userSQL := "select password from user where username=?"
	log.Print("validating user ", username)
	rows := database.query(userSQL, username)

	if rows.Next() {
		err := rows.Scan(&passwordFromDB)
		if err != nil {
			return false
		}
	}
	//If the password matches, return true
	if password == passwordFromDB {
		return true
	}
	//by default return false
	return false
}

//GetUserID will get the user's ID from the database
func GetUserID(username string) (int, error) {
	var userID int
	userSQL := "select id from user where username=?"
	rows := database.query(userSQL, username)

	if rows.Next() {
		err := rows.Scan(&userID)
		if err != nil {
			return -1, err
		}
	}
	rows.Close()
	return userID, nil
}
