package views

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/thewhitetulip/Tasks/db"
	"github.com/thewhitetulip/Tasks/types"
)

var mySigningKey = []byte("secret")

//GetTokenHandler will get a token for the username and password
func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		token := jwt.New(jwt.SigningMethodHS256)

		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		log.Println(username, " ", password)
		if db.ValidUser(username, password) {
			/* Set token claims */
			token.Claims["username"] = username
			token.Claims["exp"] = time.Now().Add(time.Hour * 5).Unix()

			/* Sign the token with our secret */
			tokenString, _ := token.SignedString(mySigningKey)

			/* Finally, write the token to the browser window */
			w.Write([]byte(tokenString))
		} else {
			w.Write([]byte("Authentication failed"))
		}
	}
}

//ValidateToken will validate the token
func ValidateToken(myToken string) (bool, string) {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	if (err != nil) || (!token.Valid) {
		return false, ""
	}
	return token.Valid, token.Claims["username"].(string)
}

//GetTasksFuncAPI fetches tasks depending on the request, the authorization will be taken care by our middleare
//in this function we will return all the tasks to the user or tasks per category
//GET /api/get-tasks/
func GetTasksFuncAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var strTaskID string
		var err error
		var message string
		var task types.Task
		var tasks types.Tasks
		var status types.Status

		token := r.Header["Token"][0]

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		if !IsTokenValid {
			status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}

		log.Println("token is valid " + username + " is logged in")

		strTaskID = r.URL.Path[len("/api/get-task/"):]
		//this is when we get a request for all the tasks for that user
		if strTaskID == "" {
			context, err := db.GetTasks(username, "pending", "")
			if err != nil {
				message = "GetTasksFuncAPI: api.go: Server error"
				log.Println(message)
				status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
				w.WriteHeader(http.StatusInternalServerError)
				err = json.NewEncoder(w).Encode(tasks)

				if err != nil {
					panic(err)
				}
				return
			}

			tasks = context.Tasks
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(tasks)

			if err != nil {
				panic(err)
			}
			return
		}
		//this is when we get a request for a particular task

		taskID, err := strconv.Atoi(strTaskID)
		if err != nil {
			message = "GetTasksFuncAPI: api.go: Invalid taskID " + strTaskID
			log.Println(message)

			status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}
		ctx, err := db.GetTaskByID(username, taskID)
		task = ctx.Tasks[0]

		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			panic(err)
		}

	}
}

//AddTaskFuncAPI will add the tasks for the user
func AddTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		token := r.Header["Token"][0]
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		if !IsTokenValid {
			status := types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}

		log.Println("token is valid " + username + " is logged in")

		r.ParseForm()
		category := r.FormValue("category")
		title := template.HTMLEscapeString(r.Form.Get("title"))
		content := template.HTMLEscapeString(r.Form.Get("content"))
		taskPriority, priorityErr := strconv.Atoi(r.FormValue("priority"))

		if priorityErr != nil {
			log.Print(priorityErr)
			message = "Bad task priority"
		}
		priorityList := []int{1, 2, 3}
		found := false
		for _, priority := range priorityList {
			if taskPriority == priority {
				found = true
			}
		}
		//If someone gives us incorrect priority number, we give the priority
		//to that task as 1 i.e. Low
		if !found {
			taskPriority = 1
		}
		var taskErr bool

		if title != "" && content != "" {
			taskTruth := db.AddTask(title, content, category, taskPriority, username, "")
			if taskTruth != nil {
				taskErr = true
			}
		}

		var statusCode int
		var message string

		if !taskErr {
			statusCode = http.StatusInternalServerError
			message = "Error adding task to db"
		} else {
			statusCode = http.StatusOK
			message = "Task added to db"
		}
		status := types.Status{StatusCode: statusCode, Message: message}
		json.NewEncoder(w).Encode(status)
	} else {
		var statusCode int
		var message string

		statusCode = http.StatusBadRequest
		message = "Invalid request"
		status := types.Status{StatusCode: statusCode, Message: message}
		json.NewEncoder(w).Encode(status)

	}
}

//UpdateTaskFuncAPI will add the tasks for the user
func UpdateTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var taskErr bool
		token := r.Header["Token"][0]
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		if !IsTokenValid {
			status := types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}

		log.Println("token is valid " + username + " is logged in")

		r.ParseForm()
		strID := r.Form.Get("id")
		id, err := strconv.Atoi(strID)
		if err != nil {
			log.Println(err)
			taskErr = true
		}
		category := r.Form.Get("category")
		title := r.Form.Get("title")
		content := r.Form.Get("content")
		priority, err := strconv.Atoi(r.Form.Get("priority"))

		if err != nil {
			log.Println(err)
			priority = 1
		}

		if strID != "" && title != "" && content != "" {
			err = db.UpdateTask(id, title, content, category, priority, username)
			if err != nil {
				taskErr = true
			}
			taskErr = false
		} else {
			taskErr = true
		}

		var statusCode int
		var message string

		if taskErr {
			statusCode = http.StatusBadRequest
			message = "unable to update task id "
		} else {
			statusCode = http.StatusOK
			message = "updated task id "
		}

		status := types.Status{StatusCode: statusCode, Message: message}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
	}
}

//DeleteTaskFuncAPI will add the tasks for the user
func DeleteTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		token := r.Header["Token"][0]
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		if !IsTokenValid {
			status := types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}

		log.Println("token is valid " + username + " is logged in")
		var strtaskID string
		strtaskID = r.URL.Path[len("/api/delete-task/"):]
		var statusCode int
		var message string

		taskID, err := strconv.Atoi(strtaskID)
		if err != nil {
			log.Println("invalid task id")
			statusCode = http.StatusBadRequest
			message = "invalid TaskID"
		} else {
			db.TrashTask(username, taskID)
			if err != nil {
				statusCode = http.StatusOK
				message = "deleted task id " + strtaskID
			} else {
				statusCode = http.StatusOK
				message = "deleted task id " + strtaskID
			}
		}

		status := types.Status{StatusCode: statusCode, Message: message}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
	}
}

//GetDeletedTaskFuncAPI will get the deleted tasks for the user
func GetDeletedTaskFuncAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var err error
		var message string
		var tasks types.Tasks
		var status types.Status

		token := r.Header["Token"][0]

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		if !IsTokenValid {
			status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}

		log.Println("token is valid " + username + " is logged in")

		//this is when we get a request for all the deleted tasks for that user
		context, err := db.GetTasks(username, "deleted", "")
		if err != nil {
			message = "GetTasksFuncAPI: api.go: Server error"
			log.Println(message)
			status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(tasks)

			if err != nil {
				panic(err)
			}
			return
		}

		tasks = context.Tasks
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(tasks)

		if err != nil {
			panic(err)
		}
		return
	}
}

//GetCategoryFuncAPI will return the categories for the user
//depends on the ID that we get, if we get all, then return all categories of the user
func GetCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var err error
		var message string
		var status types.Status

		token := r.Header["Token"][0]

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		if !IsTokenValid {
			status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}

		log.Println("token is valid " + username + " is logged in")
		categories := db.GetCategories(username)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(categories)
		if err != nil {
			panic(err)
		}
	}

}

//AddCategoryFuncAPI will add the category for the user
func AddCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var err error
		var message string
		var status types.Status
		var statusCode int
		var categoryErr bool

		token := r.Header["Token"][0]

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		if !IsTokenValid {
			status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}

		log.Println("token is valid " + username + " is logged in")
		r.ParseForm()

		category := r.Form.Get("category")
		if strings.Trim(category, " ") != "" {
			log.Println("adding category")
			err := db.AddCategory(username, category)
			if err != nil {
				categoryErr = true
			} else {
				categoryErr = false
			}
		} else {
			categoryErr = true
		}

		if categoryErr {
			statusCode = http.StatusInternalServerError
			message = "error adding category" + category
		} else {
			statusCode = http.StatusOK
			message = "added category " + category
		}

		status = types.Status{StatusCode: statusCode, Message: message}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
	}
}

//UpdateCategoryFuncAPI will update the category for the user
func UpdateCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var statusCode int
		var err error
		var message string
		var catErr bool
		var status types.Status

		token := r.Header["Token"][0]

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		if !IsTokenValid {
			status = types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}

		log.Println("token is valid " + username + " is logged in")

		r.ParseForm()

		oldName := r.URL.Path[len("/api/update-category/"):]
		oldName = strings.Trim(oldName, "/")

		newName := r.Form.Get("catname")

		if strings.Trim(newName, " ") != "" {
			err = db.UpdateCategoryByName(username, oldName, newName)
			if err != nil {
				catErr = true
			}
			catErr = false
		} else {
			catErr = true
		}

		if catErr {
			statusCode = http.StatusInternalServerError
			message = "unable to update category from " + oldName + " to " + newName
		} else {
			statusCode = http.StatusOK
			message = "updated category from " + oldName + " to " + newName

		}

		status = types.Status{StatusCode: statusCode, Message: message}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
	}
}

//DeleteCategoryFuncAPI will delete the category for the user
func DeleteCategoryFuncAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		categoryName := r.URL.Path[len("/delete-category/"):]
		categoryName = strings.Trim(categoryName, "/")
		var statusCode int
		var err error
		var message string
		var catErr bool

		token := r.Header["Token"][0]

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		IsTokenValid, username := ValidateToken(token)
		//When the token is not valid show the default error JSON document
		if !IsTokenValid {
			status := types.Status{StatusCode: http.StatusInternalServerError, Message: message}
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(status)

			if err != nil {
				panic(err)
			}
			return
		}

		log.Println("token is valid " + username + " is logged in")

		categoryName = strings.Trim(categoryName, " ")

		if categoryName != "" {
			catErr = true
		}

		err = db.DeleteCategoryByName(username, categoryName)
		if err != nil {
			catErr = true
		} else {
			catErr = false
		}

		if catErr {
			statusCode = http.StatusBadRequest
			message = "error deleting category" + categoryName
		} else {
			statusCode = http.StatusOK
			message = "deleted category " + categoryName
		}

		status := types.Status{StatusCode: statusCode, Message: message}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(status)
		if err != nil {
			panic(err)
		}
	}
}
