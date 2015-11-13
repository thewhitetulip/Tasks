<!DOCTYPE html>

<html>

	<head>

		<title>Deleted Tasks</title>

		<!-- Mobile viewport optimized -->
		<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">

		<!-- Bootstrap CSS -->
		<link href="/static/css/bootstrap.min.css" rel="stylesheet">
		<link href="/static/css/bootstrap-glyphicons.css" rel="stylesheet">

		<!-- Custom CSS -->
		<link href="/static/css/styles.css" rel="stylesheet">
		<link href="/static/css/sidebar.css" rel="stylesheet">
		<link href="/static/css/sidebar-bootstrap.css" rel="stylesheet">
		<link href="/static/css/font-awesome.min.css" rel="stylesheet" >

		<!-- Include Modernizr in the head, before any other Javascript -->
		<script src="/static/js/modernizr-2.6.2.min.js"></script>
		<!-- All Javascript at the bottom of the page for faster page loading -->
		<script src="/static/js/jquery.min.js"></script>
		<!-- If no online access, fallback to our hardcoded version of jQuery
		<script>window.jQuery || document.write('<script src="/static/js/jquery-1.8.2.min.js"><\/script>')</script>
		-->
		<!-- Bootstrap JS -->
		<script src="/static/js/bootstrap.min.js"></script>

		<!-- Custom JS -->
		<script src="/static/js/script.js"></script>
		<script src="/static/js/hammer.min.js"></script>
		<script src="/static/js/sidebar.js"></script>

	</head>
	

<body>
		<!-- The navigation bar-->
		<nav class="navbar navbar-default navbar-fixed-top mainHeader">
			<div class="container-fluid">
				<div class="navbar-header">

					<a class="navbar-brand" href="/trash/"> Trash</a>
				</div>
			</div>
		</nav>
		<!-- SIDEBAR -->
		<div data-sidebar="true" class="sidebar-trigger">

			<a class="sidebar-toggle" href="">
				<span class="glyphicon glyphicon-align-justify"></span>
			</a>


	<div class="sidebar-wrapper sidebar-default">
				<div class="sidebar-scroller">
					<ul class="sidebar-menu">
						<li class="sidebar-group"><span>Tasks</span>
							<ul class="sidebar-group-menu">
								<li class="sidebar-item">
									<a href="/" ><span class="glyphicon glyphicon-file"></span> <span class="nav-item">All</span></a>
								</li>
								<!--<li class="sidebar-item">

									<a href="" ><span class="glyphicon glyphicon-time"></span>  <span class="nav-item"> Reminders</span></a>
								</li>-->
								<li class="sidebar-item">
									<a href="/trash/" class="active"><span class="glyphicon glyphicon-trash"></span>  <span class="nav-item"> Trash</span></a>
								</li>
<!--
								<li class="sidebar-item"><a href="">
									<span class="glyphicon glyphicon-folder-open"></span> <span class="nav-item">Uncategorized</span></a>
								</li>
								<li class="sidebar-item">
									<a href=""><span class="glyphicon glyphicon-cog"></span> <span class="nav-item">Settings</span></a>
								</li>
								<li class="sidebar-item">
									<a href="#changeLogModal"  data-toggle="modal"><span class="glyphicon glyphicon-hand-up"></span> ChangeLog</a>
								</li>
-->
							</ul>
						</li>

					</ul>
				</div>
			</div>
		</div>

		<!--end mainHeader -->
		{{if .}}
		<a href="/delete/all"><button class="btn-danger btn glyphicon glyphicon-trash floating-action-icon"></button></a>
		{{end}}

		<!-- Add note modal -->
		<div class="modal fade " id="addNoteModal" tabindex="-1" role="dialog" aria-labelledby="newNoteLabel" aria-hidden="true">
			<div class="modal-dialog">
				<div class="modal-content">
					<div class="modal-header">
						<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
						<h4 class="modal-title" id="newNoteLabel"><span class="glyphicon glyphicon-pencil"></span>  New Task</h4>
					</div>
					<div class="modal-body">
						<form>
							<div class="form-group">
                <!-- <label for="note-title" class="control-label">Title:</label> -->
                <input type="text" class="form-control" id="add-note-title" placeholder="Title" style="border:none;border-bottom:1px solid gray; box-shadow:none;">
							</div>
							<div class="form-group">
                 <!-- <label for="note-content" class="control-label">Content:</label> -->
                <textarea class="form-control" id="add-note-content" placeholder="Content" rows="10" style="border:none;border-bottom:1px solid gray; box-shadow:none;"></textarea>
							</div>
						</form>
					</div>
					<div class="modal-footer">
						<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
						<button type="button" class="btn btn-primary" id="addNoteModalSaveBtn">Save</button>
					</div>
				</div>
			</div>
		</div>

		<!-- modal for opening note -->
		<div class="modal fade" id="openNoteModal" tabindex="-1" role="dialog" aria-hidden="true">
			<div class="modal-dialog">
				<div class="modal-content">
					<div class="modal-header">
						<button class="close" data-dismiss="modal"> &times;</button>
						<h4 class="modal-title"></h4>
					</div>
					<div class="modal-body">

					</div>
				</div>
			</div>
		</div>

        <div class="timeline">
		{{ if .}}
		    {{range .}}
            <div class="note">
                <p class="noteHeading">{{.Title}}</p><hr>
                <p class="noteContent">{{.Content}}</p>
                <span class="notefooter">
                    <ul class="menu">
                     <!-- <li role="presentation">
                          <a role="menuitem" tabindex="-1" href="/share/{{.Id}}">
                          <span class="glyphicon glyphicon-share"></span>  Share</a></li>
                      <li role="presentation">
                          <a role="menuitem" tabindex="-1" href="/mask/{{.Id}}">
						  <span class="glyphicon glyphicon-lock"></span> Mask</a></li> 
					  <li role="presentation"><a role="menuitem" tabindex="-1" href="/archive/{{.Id}}">
                        <span class="glyphicon glyphicon-inbox"></span>  Edit</a></li>!-->
                      <li role="presentation"><a role="menuitem" tabindex="-1" href="/delete/{{.Id}}">
                        <span class="glyphicon glyphicon-trash"></span>  Delete</a></li>
                      <li role="presentation"><a role="menuitem" tabindex="-1" href="/restore/{{.Id}}">
                           <span class="glyphicon glyphicon-inbox"></span>  Restore</a></li>
                        </ul>
                </span>
            </div>
			{{end}}
		{{else}}
		<p>No tasks here</p>
		{{end}}
        </div>
<footer class="footer" >
	Made in India with <span class="glyphicon glyphicon-heart"></span> by <a href="htp://github.com/thewhitetulip">@thewhitetulip</a>
</footer>


	</body>
</html>
