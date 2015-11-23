# Tasks

Tasks is a simplistic golang webapp to manage tasks, I built this tool to manage tasks which I wanted to do, there is a good kanban style boards, but I felt they were a bit too much for my taste. Also I wanted to learn the golang webapp.

Features:
1. Add, update, delete note
2. Search notes, the query is highlighted in the search results page

How you install?
==================

1. `go get github.com/thewhitetulip/Tasks`
1. change dir to the respective folder and create the db file: `cat schema.sql | sqlite3 tasks.db`
1. run `go build`
1. ./Task
1. open [localhost:8080](http://localhost:8080)

Either this or download the latest from the release tab above and enjoy!
