package types

/*
Package types is used to store the context struct which
is passed while templates are executed.
*/
//Task is the struct used to identify tasks
type Task struct {
	Id        int
	Title     string
	Content   string
	Created   string
	Priority  string
	Category  string
	Referer   string
	Comments  []Comment
	IsOverdue bool
}

//Comment is the struct used to populate comments per tasks
type Comment struct {
	ID      int
	Content string
	Created string
}

//Context is the struct passed to templates
type Context struct {
	Tasks      []Task
	Navigation string
	Search     string
	Message    string
	CSRFToken  string
	Categories []CategoryCount
	Referer    string
}

//CategoryCount is the struct used to populate the sidebar
//which contains the category name and the count of the tasks
//in each category

type CategoryCount struct {
	Name  string
	Count int
}
