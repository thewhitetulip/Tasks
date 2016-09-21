package views

import (
	"log"
	"net/http"

	"github.com/thewhitetulip/Tasks/db"
	"github.com/thewhitetulip/Tasks/sessions"
)

//RequiresLogin is a middleware which will be used for each httpHandler to check if there is any active session
func RequiresLogin(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessions.IsLoggedIn(r) {
			http.Redirect(w, r, "/login/", 302)
			return
		}
		handler(w, r)
	}
}

//LogoutFunc Implements the logout functionality. WIll delete the session information from the cookie store
func LogoutFunc(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Store.Get(r, "session")
	if err == nil { //If there is no error, then remove session
		if session.Values["loggedin"] != "false" {
			session.Values["loggedin"] = "false"
			session.Save(r, w)
		}
	}
	http.Redirect(w, r, "/login", 302) //redirect to login irrespective of error or not
}

//LoginFunc implements the login functionality, will add a cookie to the cookie store for managing authentication
func LoginFunc(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Store.Get(r, "session")

	if err != nil {
		log.Println("error identifying session")
		loginTemplate.Execute(w, nil)
		return
	}

	switch r.Method {
	case "GET":
		loginTemplate.Execute(w, nil)
	case "POST":
		log.Print("Inside POST")
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if (username != "" && password != "") && db.ValidUser(username, password) {
			session.Values["loggedin"] = "true"
			session.Values["username"] = username
			session.Save(r, w)
			log.Print("user ", username, " is authenticated")
			http.Redirect(w, r, "/", 302)
			return
		}
		log.Print("Invalid user " + username)
		loginTemplate.Execute(w, nil)
	default:
		http.Redirect(w, r, "/login/", http.StatusUnauthorized)
	}
}
