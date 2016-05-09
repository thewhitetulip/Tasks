package views

import (
	"net/http"

	"github.com/thewhitetulip/Tasks/sessions"
)

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
		loginTemplate.Execute(w, nil) // in case of error during fetching session info, execute login template
	} else {
		isLoggedIn := session.Values["loggedin"]
		if isLoggedIn != "true" {
			if r.Method == "POST" {
				if r.FormValue("password") == "secret" && r.FormValue("username") == "user" {
					session.Values["loggedin"] = "true"
					session.Save(r, w)
					http.Redirect(w, r, "/", 302)
					return
				}
			} else if r.Method == "GET" {
				loginTemplate.Execute(w, nil)
			}
		} else {
			http.Redirect(w, r, "/", 302)
		}
	}
}
