package views

import (
	"net/http"

	"github.com/thewhitetulip/Tasks/sessions"
)

//LogoutFunc Implements the logout functionality. WIll delete the session information from the cookie store
func LogoutFunc(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	if session.Values["loggedin"] != "false" {
		session.Values["loggedin"] = "false"
		session.Save(r, w)
		http.Redirect(w, r, "/login", 302)
		return
	}
	http.Redirect(w, r, "/login", 302)
}

//LoginFunc implements the login functionality, will add a cookie to the cookie store for managing authentication
func LoginFunc(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")

	if r.Method == "POST" && r.FormValue("password") == "secret" && r.FormValue("username") == "user" {
		session.Values["loggedin"] = "true"
		session.Save(r, w)
		http.Redirect(w, r, "/", 302)
		return
	}

	if session.Values["loggedin"] == "true" {
		http.Redirect(w, r, "/", 302)
	} else {
		loginTemplate.Execute(w, nil)
	}

}
