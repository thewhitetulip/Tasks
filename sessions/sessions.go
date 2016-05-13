package sessions

import (
	"net/http"

	"github.com/gorilla/sessions"
)

//Store the cookie store which is going to store session data in the cookie
var Store = sessions.NewCookieStore([]byte("secret-password"))

//IsLoggedIn will check if the user has an active session and return True
func IsLoggedIn(r *http.Request) bool {
	session, err := Store.Get(r, "session")
	if err == nil && (session.Values["loggedin"] == "true") {
		return true
	}
	return false
}
