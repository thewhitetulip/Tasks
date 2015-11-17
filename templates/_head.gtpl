<!DOCTYPE html>

<html>

<head>

  <title>Tasks</title>

  <!-- Mobile viewport optimized -->
  <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">

  <!-- Bootstrap CSS -->
  <link href="/static/css/bootstrap.min.css" rel="stylesheet">
  <link href="/static/css/bootstrap-glyphicons.css" rel="stylesheet">

  <!-- Custom CSS -->
  <link href="/static/css/styles.css" rel="stylesheet">
  <link href="/static/css/sidebar.css" rel="stylesheet">
  <link href="/static/css/sidebar-bootstrap.css" rel="stylesheet">
  <link href="/static/css/font-awesome.min.css" rel="stylesheet">

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

				<a class="navbar-brand" href="/"> Tasks</a>
				<span id="icons">
						<form action="/search/" method="POST">
							<input type="text" name="query" placeholder="Search" style="border:none;border-bottom:1px solid gray; box-shadow:none;">
							<input type="submit" text="submit" class="btn btn-default"/>
						</form>
					</span>
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
								<a href="/"  ><span class="glyphicon glyphicon-tasks"></span> <span class="nav-item">Pending</span></a>
								<!-- class="active" -->
							</li>
							<li class="sidebar-item">
								<a href="/complete/"><span class="glyphicon glyphicon-check"></span>  <span class="nav-item"> Completed</span></a>
							</li>
							<li class="sidebar-item">
								<a href="/deleted/"><span class="glyphicon glyphicon-trash"></span>  <span class="nav-item"> Deleted</span></a>
							</li>
						</ul>
					</li>

				</ul>
			</div>
		</div>
	</div>
