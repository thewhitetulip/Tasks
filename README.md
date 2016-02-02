# Tasks

Tasks is a simplistic Go webapp to manage tasks, I built this tool to manage tasks which I wanted to do, there are many good kanban style boards, but I felt they were a bit too heavyweight for my taste. Also I wanted to learn the Go webapp development.

##### Book
I am learning writing webapps with Go as I build this application, I took to writing an introductory book about [building webapps in Go] (https://github.com/thewhitetulip/web-dev-golang-anti-textbook) because I faced a lot of problems while learning how to write webapps in Go, it, the book strives to teach by practical examples.

##Features

1. Add, update, delete task
2. Search tasks, the query is highlighted in the search results page
3. We use github flavoured markdown, which enables us for using a task list, advanced syntax highlighting and much more
4. Supports file upload, randomizes the file name, stores the user given filename in a db and works on the randomized file name for security reasons.
5. Priorities are assigned, High = 3, medium = 2 and low = 1, sorting is done on priority descending and created date ascending
6. Categories are supported, you can add tasks to different categories. 

How you install?
==================

1. `go get github.com/thewhitetulip/Tasks`
1. change dir to the respective folder and create the db file: `cat schema.sql | sqlite3 tasks.db`
1. run `go build`
1. `./Task`
1. open [localhost:8081](http://localhost:8081)

You can change the port in the [config] (https://github.com/thewhitetulip/Tasks/blob/master/config.json) file

#Screenshots
The Home Page

![Home Page] (https://github.com/thewhitetulip/Tasks/blob/master/screenshots/FrontEnd.png)

Add Task dialog

![Add Task] (https://github.com/thewhitetulip/Tasks/blob/master/screenshots/FrontEnd-Add%20task.png)

Navigation drawer

![Navigation Drawer] (https://github.com/thewhitetulip/Tasks/blob/master/screenshots/FrontEnd%20Navigation%20Drawer.png)

#License

The MIT License (MIT)

Copyright (c) 2015 Suraj Patil

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
