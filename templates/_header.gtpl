		<!-- The navigation bar-->
		<nav class="navbar navbar-default navbar-fixed-top mainHeader">
			<div class="container-fluid">
				<div class="navbar-header">

					<a class="navbar-brand" href=""> Tasks</a>
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
									<a href="" class="active"><span class="glyphicon glyphicon-file"></span> <span class="nav-item">Notes</span></a>
								</li>
								<li class="sidebar-item">
									<a href="" ><span class="glyphicon glyphicon-time"></span>  <span class="nav-item"> Reminders</span></a>
								</li>
								<li class="sidebar-item">
									<a href="" ><span class="glyphicon glyphicon-trash"></span>  <span class="nav-item"> Trash</span></a>
								</li>
								<li class="sidebar-item"><a href="">
									<span class="glyphicon glyphicon-folder-open"></span> <span class="nav-item">Uncategorized</span></a>
								</li>
								<li class="sidebar-item">
									<a href=""><span class="glyphicon glyphicon-cog"></span> <span class="nav-item">Settings</span></a>
								</li>
								<li class="sidebar-item">
									<a href="#changeLogModal"  data-toggle="modal"><span class="glyphicon glyphicon-hand-up"></span> ChangeLog</a>
								</li>
							</ul>
						</li>

					</ul>
				</div>
			</div>
		</div>

		<!--end mainHeader -->
		<button class=" btn-danger btn glyphicon glyphicon-plus floating-action-icon"></button>

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
