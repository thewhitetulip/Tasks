{{template "_head.gtpl"}}

  <!--end mainHeader -->
  <button class=" btn-danger btn glyphicon glyphicon-plus floating-action-icon floating-action-icon-add"></button>

  <!-- Add note modal -->
  <div class="modal fade " id="addNoteModal" tabindex="-1" role="dialog" aria-labelledby="newNoteLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
          <h4 class="modal-title" id="newNoteLabel"><span class="glyphicon glyphicon-pencil"></span>  New Task</h4>
        </div>
        <div class="modal-body">
          <form action="/add/" method="POST">
            <div class="form-group">
              <!-- <label for="note-title" class="control-label">Title:</label> -->
              <input type="text" name="title" class="form-control" id="add-note-title" placeholder="Title" style="border:none;border-bottom:1px solid gray; box-shadow:none;">
            </div>
            <div class="form-group">
              <!-- <label for="note-content" class="control-label">Content:</label> -->
              <textarea class="form-control" name="content" id="add-note-content" placeholder="Content" rows="10" style="border:none;border-bottom:1px solid gray; box-shadow:none;"></textarea>
            </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
          <input type="submit" text="submit" class="btn btn-default" />
        </div>
        </form>
      </div>
    </div>
  </div>

  <div class="timeline">
    {{ if .}} {{range .}}
    <div class="note">
      <p class="noteHeading">{{.Title}}</p>
      <hr>
      <p class="noteContent">{{.Content}}</p>
      <span class="notefooter">
                    <ul class="menu">
                      <li role="presentation">
                          <!-- <a role="menuitem" tabindex="-1" href="/share/{{.Id}}"> </a>-->
                          <span class="glyphicon glyphicon-time"></span> {{.Created}}</li>
      <!--  <li role="presentation">
                          <a role="menuitem" tabindex="-1" href="/mask/{{.Id}}">
						  <span class="glyphicon glyphicon-lock"></span> Mask</a></li> !-->
      <li role="presentation">
        <a role="menuitem" tabindex="-1" href="/trash/{{.Id}}">
          <span class="glyphicon glyphicon-trash"></span> Trash</a>
      </li>
      <li role="presentation">
        <a role="menuitem" tabindex="-1" href="/complete/{{.Id}}">
          <span class="glyphicon glyphicon-check"></span> Complete</a>
      </li>
      <li role="presentation">
        <a role="menuitem" tabindex="-1" href="/edit/{{.Id}}">
          <span class="glyphicon glyphicon-pencil"></span> Edit</a>
      </li>

      </ul>
      </span>
    </div>
    {{end}} {{else}}
    <p>No tasks here</p>
    {{end}}
  </div>
	{{template "footer.gtpl"}}

</body>

</html>
