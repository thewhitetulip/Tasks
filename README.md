# Tasks

Tasks is a simplistic golang webapp to manage tasks, I built this tool to manage tasks which I wanted to do, there is a good kanban style boards, but I felt they were a bit too much for my taste. Also I wanted to learn the golang webapp.

Features:

1. Add, update, delete task
2. Search tasks, the query is highlighted in the search results page
3. We use github flavoured markdown, which enables us for using a task list, advanced syntax highlighting and much more
4. Supports file upload, randomizes the file name, stores the user given filename in a db and works on the randomized file name for security reasons.

How you install?
==================

1. `go get github.com/thewhitetulip/Tasks`
1. change dir to the respective folder and create the db file: `cat schema.sql | sqlite3 tasks.db`
1. run `go build`
1. `./Task`
1. open [localhost:8080](http://localhost:8080)

Either this or download the latest from the release tab above and enjoy!


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
