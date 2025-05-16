// Author: Vinhthuy Phan, 2018
package frontEnd

var STUDENT_MESSAGING_TEMPLATE = `
<html>
	<head>
  		<title>Student messaging</title>
		<meta http-equiv="refresh" content="10" />
	</head>
	<style>
		.bottom {
			position: fixed;
			bottom: 0;
			font-size: 150%;
			color: red;
		}
	</style>
	<body>
	<div class="bottom">{{.Message}}</div>
	</body>
</html>
`

var TEACHER_MESSAGING_TEMPLATE = `
<html>
	<head>
  		<title>TeacherMap messaging</title>
		<script src="https://cdn.rawgit.com/google/code-prettify/master/loader/run_prettify.js?autoload=true&skin=sons-of-obsidian"></script>
  		<script src="http://code.jquery.com/jquery-3.1.1.min.js"></script>
	    <script type="text/javascript">
			var updateInterval = 5000;		// 5 sec update interval
			var maxUpdateTime =  1800000;   // no longer update after 30 min.
			var totalUpdateTime = 0;
			function getData() {
				var url = "http://{{.Address}}/bulletin_board_data";
				$.getJSON(url, function( data ) {
					console.log(data);
					$("#p1").html(data["P1"]);
					$("#p2").html(data["P2"]);
					$("#p1u").html(data["P1Ungraded"]);
					$("#p1g").html(data["P1Graded"]);
					$("#p2u").html(data["P2Unanswered"]);
					$("#p2a").html(data["P2Answered"]);
				});
			}
			$(document).ready(function(){
				getData();
				handle = setInterval(getData, updateInterval);
			});
	    </script>
	</head>
	<style>
		.bottom {
			position: fixed;
			bottom: 0;
			text-align: center;
			width: 100%;
		}
		.bottom_left {
			position: fixed;
			bottom: 0;
			text-align: left,
			width: 50%;
		}
		.bottom_right {
			position: fixed;
			bottom: 0;
			right: 50px;
			text-align: right,
			width: 50%;
		}
		.label{ display: inline; }
		#p1, #p2, #p1g, #p1u, #p2a, #p2u, #ans, #ap, #bu, #at {
			padding: 0.75em;
			display: inline;
		}
		#p1g, #p2a { color: green; }
		#p1u, #p2u { color: red; }
		pre {
			font-family: monospace;
			font-size:120%;
			margin-top:50px;
			padding-left:2em;
			overflow-x:scroll;
			overflow-y:scroll;
			tab-size: 4;
			-moz-tab-size: 4;
		}
		.center {
		    text-align: center;
		}
		.pagination {
		    display: inline-block;
		    padding-bottom: 20px;
		}
		.pagination a {
		    color: black;
		    float: left;
		    padding: 8px 16px;
		    text-decoration: none;
		    transition: background-color .3s;
		    border: 1px solid #ddd;
		    margin: 0 4px;
		    border-radius: 5px;
		}
		.pagination a.active {
		    background-color: #4CAF50;
		    color: white;
		    border: 1px solid #4CAF50;
		    border-radius: 5px;
		}
		.pagination a:hover:not(.active) {background-color: #ddd;}
		.nav a { text-decoration: none; padding:3px;}
		.nav { display: inline-block; vertical-align: baseline;}
		#navWrap{position:absolute;top:20;right:10;}
	</style>
	<body>
	<div id="navWrap">
	{{ if .Authenticated }}
	<div class="nav"><a href="view_bulletin_board?i=0&pc={{.PC}}">First<a></div>
	<div class="nav"><a href="view_bulletin_board?i={{.PrevI}}&pc={{.PC}}">Prev<a></div>
	<div class="nav"><a href="view_bulletin_board?i={{.NextI}}&pc={{.PC}}">Next<a></div>
	<div class="nav"><a href="remove_bulletin_page?i={{.I}}&pc={{.PC}}">&#x2718;</a></div>
	{{ end }}
	</div>
	<pre class="prettyprint linenums">{{.Code}}</pre>

	<div class="bottom_left">
		<div id="p2">{{.P2}}</div> <div class="label">Help Requests:</div>
		<div id="p2u">{{.P2Unanswered}}</div> <div class="label">Pending,</div>
		<div id="p2a">{{.P2Answered}}</div> <div class="label">Answered</div>
	</div>
	<div class="bottom_right">
		<div id="p1">{{.P1}}</div> <div class="label">Submissions:</div>
		<div id="p1u">{{.P1Ungraded}}</div> <div class="label">Pending,</div>
		<div id="p1g">{{.P1Graded}}</div> <div class="label">Graded</div>
	</div>
	</body>
</html>
`
var CODESPACE_TEMPLATE = `
	<!DOCTYPE html>
	<html>
	<head>
	<title>CodeSpace</title>
	<meta http-equiv="refresh" content="120" >
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	</head>
	<body>
	<div class="container">
		<h3 class="title is-3">CodeSpace: List of Code Snapshots</h3>
		<table class="table is-striped is-fullwidth is-hoverable is-narrow">
			<thead>
				<tr>
					<th>Student</th>
					<th>Last Snapshot</th>
					<th>Time Spent</th>
					<th>Lines of Code</th>
					<th>Number of Feedback Messages</th>
					<th>Status</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
			{{ range .Snapshots }}
			<tr>
				<td>{{ .StudentName }}</td>
				<td>{{ formatTimeSince .LastUpdated }} ago</td>
				<td>{{ formatTimeSince .FirstUpdate }}</td>
				<td>{{ .LinesOfCode }}</td>
				<td>{{ .NumFeedback }}</td>
				<td>{{ if eq .Status 0 }} Not Submitted {{else if eq .Status 1}} Submitted {{else if eq .Status 2}} <span style="font-size: 1.5em; color: red;"> <i class="far fa-times-circle"></i> </span> {{else if eq .Status 3}} <span style="font-size: 1.5em; color: green;"> <i class="far fa-check-circle"></i> </span> {{end}}</td>
				<td><a href="/get_snapshot?student_id={{ .StudentID }}&problem_id={{ .ProblemID }}&uid={{$.UserID}}&role={{$.UserRole}}&password={{$.Password}}">View</a></td>
			</tr>
			{{ end }}
			</tbody>
		</table>
	</div>

	</body>
	</html>
`
var CODE_SNAPSHOT_TEMPLATE = `
<!DOCTYPE html>
	<html>
	<head>
	<title>Latest Code Snapshot from {{.Snapshot.StudentName}}</title>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	</head>
	<body>
		<div class="container">
			<section class="section">
				<h3 class="title is-3">Latest Code Snapshot from {{.Snapshot.StudentName}}</h3>
				<h4 class="title is-4">{{.Snapshot.StudentName}} ({{.Snapshot.ProblemName}} @ {{.Snapshot.LastUpdated.Format "Jan 02, 2006 3:04:05 PM"}})</h4>
				<h5 class="title is-5">If you think that this student needs help, feel free to offer a brief comment.</h5>
				{{$l := (len .HelpRequestIDs)}}
				{{if ne $l 0}}
					<b>Help requests: </b>
				{{end}}
				{{range $i, $v := .HelpRequestIDs}}
					<a href="/view_help_request?request_id={{.}}&uid={{$.UserID}}&role={{$.UserRole}}&password={{$.Password}}">Request {{add $i 1}}</a>{{if lt (add $i 1) $l}} | {{end}}
				{{end}}
				<textarea id="editor">{{ .Snapshot.Code }}</textarea>
			</section>
			{{if lt .Snapshot.Status 3}}
			<section class="section" style="margin-top: 0px !important;">
				<form action="/save_snapshot_feedback" method="POST">
					<textarea class="textarea" placeholder="Write your feedback!" name="feedback"></textarea>
					<input class="button" type="submit" value="Send Feedback">
					
					<input type="hidden" name="snapshot_id" value="{{.Snapshot.ID}}">
					<input type="hidden" name="uid" value="{{.UserID}}">
					<input type="hidden" name="role" value="{{.UserRole}}">
					<input type="hidden" name="password" value="{{.Password}}">
				</form>
			</section>
			{{end}}
			<section class="section">
				{{range .Feedbacks}}
					<article class="message">
						<div class="message-header">
						<p>{{.GivenBy}} ({{.FeedbackTime.Format "Jan 02, 2006 3:4:5 PM"}})</p>
						</div>
						<div class="message-body">
							<div class="columns">
								<div class="column is-three-quarters">{{.Feedback}}</div>
								<div class="column">
									<a onclick="autoFeedbackSubmit('yes', {{.FeedbackID}})">
										<span style="font-size: 1.5em; {{if eq .CurrentUserVote "yes"}} color: green; {{end}}">
											<i class="fas fa-thumbs-up"></i>
										</span>
									</a>
									<span>
											{{.Upvote}}
									</span>
								</div>
								<div class="column">
									<a onclick="autoFeedbackSubmit('no', {{.FeedbackID}})">
										<span style="font-size: 1.5em; {{if eq .CurrentUserVote "no"}} color: red; {{end}}">
											<i class="fas fa-thumbs-down"></i>
										</span>
									</a>
									<span>
										{{.Downvote}}
									</span>
								</div>
							</div>
							<div class="codesnapshots">
								<h3>Code Snapshot</h3>
								<div>
									<textarea class="editors">{{ .Code }}</textarea>
								</div>
							</div>
						</div>
					</article>
				{{end}}
			</section>
		</div>
		<script>
			var editor = document.getElementById("editor");
			var myCodeMirror = CodeMirror.fromTextArea(editor, {lineNumbers: true, mode: "{{getEditorMode .Snapshot.ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
			myCodeMirror.setSize("100%", 400)
			var snapshotEditors = document.getElementsByClassName("editors");
			for (i = 0;i<snapshotEditors.length; i++) {
				CodeMirror.fromTextArea(snapshotEditors[i], {lineNumbers: true, mode: "{{getEditorMode .Snapshot.ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
			}
			$( function() {
				$( ".codesnapshots" ).accordion({
					collapsible: true,
					active: false
				});
			} );
			function autoFeedbackSubmit(backFeedback, fID) {
				$.ajax({
					url: "/save_snapshot_back_feedback",
					type: "POST",
					data:  {
						feedback: backFeedback,
						feedback_id: fID,
						uid: {{.UserID}},
						role: "{{.UserRole}}",
						password: "{{.Password}}",
					},
					success: function(data){
						console.log("Success!")
					}
				});
				
				location.reload();
			}
		</script>
	</body>
	</html>
`
var STUDENT_VIEWS_FEEDBACK_TEMPLATE = `
<!DOCTYPE html>
	<html>
	<head>
	<title>Review Feedback</title>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	</head>
	<body>
		<div class="container">
		<h1 class="title">Review Feedback for Problem: {{.Filename}}</h1>
			<div class="tabs is-centered is-boxed is-medium">
				<ul>
					<li {{if eq .ViewType "forme"}}class="is-active"{{end}}>
						<a href="student_views_feedback?pid={{.CurrentPid}}&viewtype=forme&role={{.UserRole}}&uid={{.UserID}}&password={{.Password}}">
						<span class="icon is-small"><i class="fas fa-address-book" aria-hidden="true"></i></span>
						<span>For me</span>
						</a>
					</li>
					<li {{if eq .ViewType "all"}}class="is-active"{{end}}>
						<a href="student_views_feedback?pid={{.CurrentPid}}&viewtype=all&role={{.UserRole}}&uid={{.UserID}}&password={{.Password}}">
						<span class="icon is-small"><i class="fas fa-list-ul" aria-hidden="true"></i></span>
						<span>All</span>
						</a>
					</li>
				</ul>
			</div>
			<section class="section">
				{{range .Feedbacks}}
					<article class="message">
						<div class="message-header">
						<p>{{.GivenBy}} gave feedback on {{$.Filename}} at ({{.FeedbackTime.Format "Jan 02, 2006 3:04:05 PM"}})</p>
						</div>
						<div class="message-body">
							<div class="columns">
								<div class="column is-three-quarters">{{.Feedback}}</div>
								<div class="column">
									<a onclick="autoFeedbackSubmit('yes', {{.FeedbackID}})">
										<span style="font-size: 1.5em; {{if eq .CurrentUserVote "yes"}} color: green; {{end}}">
											<i class="fas fa-thumbs-up"></i>
										</span>
									</a>
									<span>
											{{.Upvote}}
									</span>
								</div>
								<div class="column">
									<a onclick="autoFeedbackSubmit('no', {{.FeedbackID}})">
										<span style="font-size: 1.5em; {{if eq .CurrentUserVote "no"}} color: red; {{end}}">
											<i class="fas fa-thumbs-down"></i>
										</span>
									</a>
									<span>
										{{.Downvote}}
									</span>
								</div>
							</div>
							<div class="codesnapshots">
								<h3>Code Snapshot</h3>
								<div>
									<textarea class="editors">{{ .Code }}</textarea>
								</div>
							</div>
						</div>
					</article>
				{{end}}
			</section>
			<nav class="pagination is-rounded" role="navigation" aria-label="pagination">
			{{if not (eq .NextPid -1)}}
				<a class="pagination-next" href="student_views_feedback?pid={{.NextPid}}&viewtype={{.ViewType}}&role={{.UserRole}}&uid={{.UserID}}&password={{.Password}}">Next</a>
			{{end}}
				<ul class="pagination-list">
				</ul>
			</nav>
		</div>
		<script>
			var snapshotEditors = document.getElementsByClassName("editors");
			
			for (let i = 0; i<snapshotEditors.length; i++){
				CodeMirror.fromTextArea(snapshotEditors[i], {lineNumbers: true, mode: "{{getEditorMode $.Filename}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
			}
			$( function() {
				$( ".codesnapshots" ).accordion({
					collapsible: true,
					active: false
				});
			} );
			function autoFeedbackSubmit(backFeedback, fID) {
				$.ajax({
					url: "/save_snapshot_back_feedback",
					type: "POST",
					data:  {
						feedback: backFeedback,
						feedback_id: fID,
						uid: {{.UserID}},
						role: "{{.UserRole}}",
						password: "{{.Password}}",
					},
					success: function(data){
						console.log("Success!")
					}
				});
				
				location.reload();
			}
		</script>
	</body>
	</html>
`
var TEACHER_VIEWS_FEEDBACK_TEMPLATE = `
<!DOCTYPE html>
	<html>
	<head>
	<title>Review Feedback</title>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	</head>
	<body>
		<div class="container">
		<<h1 class="title">Review Feedback for Problem: {{.Filename}}</h1>
			<section class="section">
				{{range .Feedbacks}}
					<article class="message">
						<div class="message-header">
						<p>{{.GivenBy}} gave feedback on {{$.Filename}} at ({{.FeedbackTime.Format "Jan 02, 2006 3:04:05 PM"}})</p>
						</div>
						<div class="message-body">
							<div class="columns">
								<div class="column is-three-quarters">{{.Feedback}}</div>
								<div class="column">
									<a onclick="autoFeedbackSubmit('yes', {{.FeedbackID}})">
										<span style="font-size: 1.5em; {{if eq .CurrentUserVote "yes"}} color: green; {{end}}">
											<i class="fas fa-thumbs-up"></i>
										</span>
									</a>
									<span>
											{{.Upvote}}
									</span>
								</div>
								<div class="column">
									<a onclick="autoFeedbackSubmit('no', {{.FeedbackID}})">
										<span style="font-size: 1.5em; {{if eq .CurrentUserVote "no"}} color: red; {{end}}">
											<i class="fas fa-thumbs-down"></i>
										</span>
									</a>
									<span>
										{{.Downvote}}
									</span>
								</div>
							</div>
							<div class="codesnapshots">
								<h3>Code Snapshot</h3>
								<div>
									<textarea class="editors">{{ .Code }}</textarea>
								</div>
							</div>
						</div>
					</article>
				{{end}}
			</section>
			<nav class="pagination is-rounded" role="navigation" aria-label="pagination">
			{{if not (eq .NextPid -1)}}
				<a class="pagination-next" href="teacher_views_feedback?pid={{.NextPid}}&role={{.UserRole}}&uid={{.UserID}}&password={{.Password}}">Next</a>
			{{end}}
				<ul class="pagination-list">
				</ul>
			</nav>
		</div>
		<script>
			var snapshotEditors = document.getElementsByClassName("editors");
			
			for (let i = 0; i<snapshotEditors.length; i++){
				CodeMirror.fromTextArea(snapshotEditors[i], {lineNumbers: true, mode: "{{getEditorMode .Filename}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
			}
			$( function() {
				$( ".codesnapshots" ).accordion({
					collapsible: true,
					active: false
				});
			} );
			function autoFeedbackSubmit(backFeedback, fID) {
				$.ajax({
					url: "/save_snapshot_back_feedback",
					type: "POST",
					data:  {
						feedback: backFeedback,
						feedback_id: fID,
						uid: {{.UserID}},
						role: "{{.UserRole}}",
						password: "{{.Password}}",
					},
					success: function(data){
						console.log("Success!")
					}
				});
				
				location.reload();
			}
		</script>
	</body>
	</html>
`
var HELP_REQUEST_LIST_TEMPLATE = `
	<!DOCTYPE html>
	<html>
	<head>
	<title>Help Hotline</title>
	<meta http-equiv="refresh" content="120" >
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	</head>
	<body>
	<div class="container">
	<h3 class="title is-3">Help Hotline/<h3>
	<h5 class="title is-5">Currently, there are {{.NumHelpNeeded}} students who need help.</h5>
		<table class="table is-striped is-fullwidth is-hoverable is-narrow">
			<thead>
				<tr>
					<th>Student</th>
					<th>Given At</th>
					<th># of Reply</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
			{{ range .HelpRequests }}
			<tr>
				<td>{{ .StudentName }}</td>
				<td>{{ formatTimeSince .GivenAt }} ago</td>
				<td>{{ .NumReply}}</td>
				<td><a href="/view_help_request?request_id={{.ID}}&uid={{$.UserID}}&role={{$.UserRole}}&password={{$.Password}}">View</a></td>
			</tr>
			{{ end }}
			</tbody>
		</table>
	</div>

	</body>
	</html>
`

var HELP_REQUEST_VIEW_TEMPLATE = `
<!DOCTYPE html>
	<html>
	<head>
	<title>Help Request</title>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	</head>
	<body>
		<div class="container">
			<h3 class="title is-3">Help Request</h3>
			<section class="section">
				<article class="message">
					<div class="message-header">
					<p>Help Request from {{.StudentName}} At ({{.GivenAt.Format "Jan 02, 2006 3:04:05 PM"}})</p>
					</div>
					<div class="message-body">
						{{.Explanation}}
					</div>
				</article>
				<textarea id="editor">{{ .Snapshot }}</textarea>
				<form action="/save_snapshot_feedback" method="POST">
					<label class="label">Feedback</label>
					<div class="field">
						<textarea class="textarea" placeholder="Write your feedback!" name="feedback"></textarea>
					</div>
					<div class="control">
						<input class="button" type="submit" value="Send Feedback">
					</div>
						
						
						<input type="hidden" name="snapshot_id" value="{{.SnapshotID}}">
						<input type="hidden" name="uid" value="{{.UserID}}">
						<input type="hidden" name="role" value="{{.UserRole}}">
						<input type="hidden" name="password" value="{{.Password}}">
					
				</form>
				<footer class="footer">
					<div class="content has-text-centered">
					<a href="/get_snapshot?snapshot_id={{.SnapshotID}}&uid={{$.UserID}}&role={{$.UserRole}}&password={{$.Password}}">View Snapshot</a>
					</div>
				</footer>
			</section>
		</div>
	</body>
	<script>
			var editor = document.getElementById("editor");
			var myCodeMirror = CodeMirror.fromTextArea(editor, {lineNumbers: true, mode: "{{getEditorMode .ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
			myCodeMirror.setSize("100%", 400)
			
		</script>
	</html>
`
var FEEDBACK_PROVISION_TEMPLATE = `
	<!DOCTYPE html>
	<html>
	<head>
	<title>Student Dashboard</title>
	<meta http-equiv="refresh" content="120" >
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />

	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	<script src="https://cdn.jsdelivr.net/npm/@creativebulma/bulma-collapsible"></script>
	<style>
		.status {
			display: flex;
			justify-content: space-between;
		}
		.menu {
			padding: 10px;
			padding-left: 100px;
			padding-right: 100px;
		}
		.show {
			top: 6%;
			position: fixed;
			z-index: 200;
			background: white;
		}
		.content {
			padding-top: 7%;
		}
		.topcorner{
			position:absolute;
			top:0;
			right:0;
		}
	</style>
	</head>
	<body>
	<div class="container">
	<nav class="navbar is-fixed-top breadcrumb menu" role="navigation" aria-label="breadcrumbs">
	<div class="navbar-start"> 
	<ul>
	  <li>
		<a id="view-exercise-link" href="#">
		  <span class="icon is-small">
			<i class="fas fa-home" aria-hidden="true"></i>
		  </span>
		  <span>Exercises</span>
		</a>
	  </li>
	  <li>
		<a id="problem-dashboard-link" href="#">
		<span class="icon is-small">
			<i class="fas fa-book" aria-hidden="true"></i>
		  </span>
			<span>Exercise Dashboard</span>
		</a>
	   </li>
	  <li class="is-active">
		<a href="#">
			<span class="icon is-small">
				<i class="fas fa-puzzle-piece" aria-hidden="true"></i>
			</span>
		  <span>{{.StudentName}}'s Dashboard</span>
		</a>
	  </li>
	</ul>
	</div>
	<div class="navbar-end"> 
		<div class="navbar-item"> <a href="#">{{.Username}} ({{.UserRole}})</a> </div>
	</div>

	</nav>
	<!--
	<nav class="breadcrumb is-right" aria-label="breadcrumbs">
		<ul>
		<li class="is-active"><a href="#">{{.Username}}({{.UserRole}})</a></li>
		</ul>
  	</nav>
	-->
	<div class="content">
	<div class="column is-two-thirds show" style="width: 70%;">
	<!--
		<div class="row">
			<h2 class="title is-2">{{.StudentName}}'s Dashboard for {{.ProblemName}}</h2>
		</div>
	-->
		<div class="row status">
			<span>Status: <strong>{{ .Status.CodingStat }} </strong></span>
			<span>Help Status: <strong>{{ .Status.HelpStat }} </strong></span>
			<span>Submission: <strong> {{ .Status.SubmissionStat }} </strong></span>
			<span>Progress: <strong>{{ .Status.Percentage }}%</strong></span>
		</div>

		<div class="tabs">
			<ul>
			<li><a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{.ProblemID}}&uid={{.UserID}}&role={{.UserRole}}{{if ne .Password ""}}&password={{.Password}}{{end}}">CodeSpace</a></li>
				<li class="is-active"><a>Feedback History</a></li>
			</ul>
		</div>

	</div>
		
	<div class="content">
		<div>
			<section class="section" style="padding: 20px">
				{{range .Messages}}
					<article class="message" style="margin-left: 25px; padding-bottom: 20px;">
						<div class="message-header">
						<p>{{if eq .Type 0}}{{.Name}} asked for help{{else if eq .Event "at_submission"}} Submission Snapshot taken {{else}} Regular Snapshot taken {{end}} at ({{.GivenAt.Format "Jan 02, 2006 3:04:05 PM"}})</p>
						</div>
						<div class="message-body">
							{{.Message}}
						</div>
						<div style="margin-left:20px;">
							{{if .Code }}

								{{range .Feedbacks}}
									<article class="message" style="margin-left: 25px;">
										<div class="message-header">
										<p>Reply from {{.Name}} given at {{.GivenAt.Format "Jan 02, 2006 3:04:05 PM"}} </p>
										</div>
										<div class="message-body">
											<div class="columns">
												<div class="column is-four-fifths">
													<textarea class="message-feedback">{{ .Feedback }}</textarea>
												</div>
												{{ if not (eq .Upvote 0) }}
												<div class="column" style="text-align: center;">
													<div style="font-size: 32px;">
														{{.Upvote}}
													</div>
													<p> Student found it helpful.</p>
													<!--
													<button class="button is-info" onclick="autoFeedbackSubmit('yes', {{.FeedbackID}})" style="margin-top:3px;" >Thank you <span style="margin:5px;"> ({{.Upvote}})</span> </button>
													-->
												</div>
												{{ end }}
												
											</div>

											
											
										</div>
									</article>
								{{end}}

							{{ end }}

						</div>
					</article>
					
				{{end}}
			</section>
		</div>
	</div>
	</div>
	<script>
		$(document).ready(function(){
			$('#view-exercise-link').attr("href", "/view_exercises"+window.location.search);
			$('#problem-dashboard-link').attr("href", "/problem_dashboard"+window.location.search+"&problem_id={{.ProblemID}}");
		

			var snapshotEditors = document.getElementsByClassName("message-feedback");
			
			for (let i = 0; i<snapshotEditors.length; i++){
				var code = CodeMirror.fromTextArea(snapshotEditors[i], {lineNumbers: true, mode: "{{getEditorMode .ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
				code.setSize("100%", 500);
			}

		});



		function autoFeedbackSubmit(backFeedback, fID) {
			$.ajax({
				url: "/save_snapshot_back_feedback",
				type: "POST",
				data:  {
					feedback: backFeedback,
					feedback_id: fID,
					uid: {{.UserID}},
					role: "{{.UserRole}}",
					{{if ne .Password ""}}password: "{{.Password}}",{{end}}
				},
				success: function(data){
					console.log("Success!")
				}
			});
			
			location.reload();
		}
	</script>

	</body>
	</html>
`
var SCAFFOLDING_TEMPLATE = `
<!DOCTYPE html>
<html lang="en">
<head>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/javascript/javascript.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" />
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/dracula.min.css" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.16/theme/monokai.min.css">
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/material.min.css" />
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/shadowfox.min.css" />
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/base16-dark.min.css" />
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/ayu-dark.min.css" />
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/tomorrow-night-bright.min.css" />
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/cobalt.min.css" />
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/lucario.min.css" />
<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.63.1/mode/markdown/markdown.js"></script>
	


    <script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" />

    <script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
    <script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js"></script>
    <link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />

    <style>
         body {
            font-family: 'Arial', sans-serif;
            margin: 20px;
        }
        .container {
            max-width: 900px;
            margin: auto;
            background: white;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
        }
        .accordions {
            margin-top: 20px;
            margin-bottom: 20px !important; /* Added margin for spacing */
        }
        .editor-container {
            margin-top: 10px;
        }
        h1 {
            text-align: center;
            color: #333;
        }
        h3 {
            background: #3498db;
            color: white;
            padding: 10px;
            border-radius: 5px;
            margin-bottom: 10px; /* Added spacing */
        }
        .CodeMirror {
            height: auto;
            min-height: 250px;
            font-size: 14px;
            padding-left: 10px; /* Prevent overlap */
        }
#scaffolding-dropdown {
    width: 100%; /* Make it span the full width of the container */
    max-width: 250px; /* Set a maximum width */
    padding: 10px; /* Add padding to the dropdown */
    border-radius: 5px; /* Rounded corners */
    border: 1px solid #ccc; /* Light border */
    background: #ffffff; /* White background */
    box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.1); 
}

    </style>
</head>
<body>
    <div class="container">
		<div style="display: flex; justify-content: center; align-items: center; position: relative; width: 100%; padding: 20px;">
    <h1 style="font-size: 2.5em; margin: 0;">Scaffoldings</h1>
	<a href="/scaffolding-guidance-html.html" 
   target="_blank" 
   style="position: absolute; left: 20px; font-size: 1em; text-decoration: none; color: #007bff;">
  Guidelines
</a>

    <select id="scaffolding-dropdown" 
        style="position: absolute; right: 20px; padding: 8px; font-size: 1em;">
        <option value="">Generate Scaffolds</option>
        <option value="1">Fill-in-the-Blanks</option>
        <option value="2">Step-by-Step Tasks</option>
        <option value="3">Guided Code with Hints</option>
        <option value="4">Debug This Code</option>
        <option value="5">Incremental Feature Implementation</option>
    </select>
</div>

        {{if or (index . 1) (index . 2) (index . 3) (index . 4) (index . 5)}}
            <!-- If at least one scaffolding data exists, render the sections -->
            
            {{if index . 1}} <!-- Check if scaffolding data exists for index 1 -->
            <div class="accordions">
                <h3>Fill-in-the-Blanks</h3>
                <div>
                    {{range $level, $scaffoldings := index . 1}}
                        <h4 style="font-weight: bold; font-size: 1.25em; margin-bottom: 10px; margin-top: 10px;">
                            Level {{$level}} 
                            {{if eq $level 1}}(Struggling){{end}}
                            {{if eq $level 2}}(Developing){{end}}
                            {{if eq $level 3}}(Nearly Proficient){{end}}
                        </h4>

                        {{range $scaffoldings}}
                        <div class="editor-container">
                            <textarea class="code-editor">{{.ScaffoldingMaterial}}</textarea>
                        </div>
                        {{end}}
                    {{end}}
                </div>
            </div>
            {{end}} <!-- End check for scaffolding data -->

            {{if index . 2}} <!-- Check if scaffolding data exists for index 2 -->
            <div class="accordions">
                <h3>Step-by-Step Tasks</h3>
                <div>
                    {{range $level, $scaffoldings := index . 2}}
                        <h4 style="font-weight: bold; font-size: 1.25em; margin-bottom: 10px; margin-top: 10px;">
                            Level {{$level}} 
                            {{if eq $level 1}}(Struggling){{end}}
                            {{if eq $level 2}}(Developing){{end}}
                            {{if eq $level 3}}(Nearly Proficient){{end}}
                        </h4>

                        {{range $scaffoldings}}
                        <div class="editor-container">
                            <textarea class="code-editor">{{.ScaffoldingMaterial}}</textarea>
                        </div>
                        {{end}}
                    {{end}}
                </div>
            </div>
            {{end}} <!-- End check for scaffolding data -->

            {{if index . 3}} <!-- Check if scaffolding data exists for index 3 -->
            <div class="accordions">
                <h3>Guided Code with Hints</h3>
                <div>
                    {{range $level, $scaffoldings := index . 3}}
                        <h4 style="font-weight: bold; font-size: 1.25em; margin-bottom: 10px; margin-top: 10px;">
                            Level {{$level}} 
                            {{if eq $level 1}}(Struggling){{end}}
                            {{if eq $level 2}}(Developing){{end}}
                            {{if eq $level 3}}(Nearly Proficient){{end}}
                        </h4>

                        {{range $scaffoldings}}
                        <div class="editor-container">
                            <textarea class="code-editor">{{.ScaffoldingMaterial}}</textarea>
                        </div>
                        {{end}}
                    {{end}}
                </div>
            </div>
            {{end}} <!-- End check for scaffolding data -->

            {{if index . 4}} <!-- Check if scaffolding data exists for index 4 -->
            <div class="accordions">
                <h3>Debug This Code</h3>
                <div>
                    {{range $level, $scaffoldings := index . 4}}
                        <h4 style="font-weight: bold; font-size: 1.25em; margin-bottom: 10px; margin-top: 10px;">
                            Level {{$level}} 
                            {{if eq $level 1}}(Struggling){{end}}
                            {{if eq $level 2}}(Developing){{end}}
                            {{if eq $level 3}}(Nearly Proficient){{end}}
                        </h4>

                        {{range $scaffoldings}}
                        <div class="editor-container">
                            <textarea class="code-editor">{{.ScaffoldingMaterial}}</textarea>
                        </div>
                        {{end}}
                    {{end}}
                </div>
            </div>
            {{end}} <!-- End check for scaffolding data -->

            {{if index . 5}} <!-- Check if scaffolding data exists for index 5 -->
            <div class="accordions">
                <h3>Incremental Feature Implementation</h3>
                <div>
                    {{range $level, $scaffoldings := index . 5}}
                        <h4 style="font-weight: bold; font-size: 1.25em; margin-bottom: 10px; margin-top: 10px;">
                            Level {{$level}} 
                            {{if eq $level 1}}(Struggling){{end}}
                            {{if eq $level 2}}(Developing){{end}}
                            {{if eq $level 3}}(Nearly Proficient){{end}}
                        </h4>

                        {{range $scaffoldings}}
                        <div class="editor-container">
                            <textarea class="code-editor">{{.ScaffoldingMaterial}}</textarea>
                        </div>
                        {{end}}
                    {{end}}
                </div>
            </div>
            {{end}} <!-- End check for scaffolding data -->

        {{else}} <!-- If no scaffolding data exists, show the fallback message -->
            <div class="no-data-message" style="text-align: center; font-weight: bold; background-color: aliceblue;">
    <p>No scaffolding data generated yet.</p>
</div>


        {{end}} <!-- End check for any scaffolding data -->

    </div>
</body>
    <script>
     $(document).ready(function () {
        $(".accordions").accordion({
            header: "h3",
            active: false,
            collapsible: true,
            heightStyle: "content",
            activate: function (event, ui) {
                if (ui.newPanel.length) {
                    ui.newPanel.css("max-height", "900px").css("overflow-y", "auto");
                    ui.newPanel.find(".CodeMirror").each(function () {
                        this.CodeMirror.refresh();
                    });
                }
            }
        });

        $(".code-editor").each(function (index, textarea) {
            let content = $(textarea).text().trim();
            let editor = CodeMirror.fromTextArea(textarea, {
                lineNumbers: true,
                mode: "markdown",
                theme: "monokai",
                matchBrackets: true,
                indentUnit: 4,
                indentWithTabs: true,
                readOnly: true
            });

            if (content) {
                editor.setValue(content);
            } else {
                editor.setValue("// No content available...");
            }

            editor.refresh();
        });
    });
document.getElementById("scaffolding-dropdown").addEventListener("change", sendScaffoldingStrategy);

function sendScaffoldingStrategy() {
    const dropdown = document.getElementById("scaffolding-dropdown");
    const selectedStrategy = dropdown.value;
    if (!selectedStrategy) return;

    // Retrieve problem_id from URL
    const urlParams = new URLSearchParams(window.location.search);
    const problemId = urlParams.get("problem_id");

    if (!problemId) {
        alert("Error: Problem ID is missing!");
        return;
    }

    // Disable the dropdown and show "Generating..."
    dropdown.disabled = true;
    const originalText = dropdown.options[dropdown.selectedIndex].text;
    dropdown.options[dropdown.selectedIndex].text = "Generating...";

    const requestData = {
        problem_id: problemId,
        scaffolding_strategy: selectedStrategy
    };

    fetch("/process_scaffolding", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(requestData)
    })
    .then(response => response.json().then(data => ({ status: response.status, body: data })))
    .then(({ status, body }) => {
        if (status === 409) {
            alert(body.message); // Show message from backend (Scaffolding already exists)
        } else if (status === 200) {
            alert("Scaffolding strategy submitted successfully!");
			window.location.reload();
        } else {
            alert("Unexpected response from server!");
        }
    })
    .catch(error => {
        console.error("Error:", error);
        alert("Error submitting scaffolding strategy!");
    })
    .finally(() => {
        // Re-enable the dropdown and restore the original text
        dropdown.disabled = false;
        dropdown.options[dropdown.selectedIndex].text = originalText;
    });
}
    </script>
</body>
</html>
`

var SC_TEMPLATE = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>RapidResponse</title>
    <!-- Included Feather Icons for UI elements -->
    <script src="https://cdnjs.cloudflare.com/ajax/libs/feather-icons/4.29.0/feather.min.js"></script>
    <link rel="stylesheet" href="style.css?v=1.1">
</head>
<style>
/* Reset and Base Styles */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
}

body {
    background-color: #f9fafb;
    color: #111827;
    height: 100vh;
    display: flex;
    flex-direction: column;
}

/* Header Styles */
.header-content {
    background-color: #1d4ed8;
    color: white;
    padding: 1rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: nowrap;
    position: relative; /* Added to ensure proper z-index context */
    z-index: 1; /* Base z-index for header */
}

.header-title {
    font-size: 1.5rem;
    font-weight: bold;
    margin-right: 10px;
    white-space: nowrap;
}

.header-info {
    display: flex;
    gap: 1rem;
    white-space: nowrap;
}

/* Timer Styles - Simplified for one row */
.timer-container {
    display: flex;
    align-items: center;
    background-color: rgba(255, 255, 255, 0.1);
    border-radius: 4px;
    padding: 4px 8px;
    margin: 0 10px;
}

.timer-display {
    font-size: 1.1rem;
    font-weight: bold;
    color: white;
    margin-right: 8px;
    min-width: 50px;
}

.timer-buttons {
    display: flex;
    gap: 4px;
}

.timer-btn {
    background-color: rgba(255, 255, 255, 0.15);
    color: white;
    border: none;
    border-radius: 4px;
    width: 28px;
    height: 28px;
    font-size: 0.75rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
    padding: 0;
}

.timer-btn:hover {
    background-color: rgba(255, 255, 255, 0.25);
}

.timer-btn.reset {
    background-color: rgba(239, 68, 68, 0.5);
}

.timer-btn.reset:hover {
    background-color: rgba(239, 68, 68, 0.7);
}

.timer-btn.pause {
    background-color: rgba(245, 158, 11, 0.5);
}

.timer-btn.pause:hover {
    background-color: rgba(245, 158, 11, 0.7);
}

.timer-warning {
    color: #fef08a;
    animation: pulse 1.5s infinite;
}

/* Add this for responsive design */
@media (max-width: 960px) {
    .header-content {
        flex-wrap: wrap;
        gap: 10px;
    }
    
    .header-title {
        order: 1;
    }
    
    .timer-container {
        order: 2;
        margin-left: auto;
    }
    
    .header-info {
        order: 3;
        flex-basis: 100%;
        justify-content: center;
    }
}

@keyframes pulse {
    0% { opacity: 1; }
    50% { opacity: 0.5; }
    100% { opacity: 1; }
}


/* Main Layout */
.main-content {
    display: flex;
    flex: 1;
    overflow: hidden;
    position: relative; /* Added to ensure proper z-index context */
}

/* Collapsible elements */
.collapsible {
    cursor: pointer;
    position: relative;
}

.collapsible:after {
    content: '\25BC'; /* Down arrow */
    position: absolute;
    right: 10px;
    top: 10px;
    color: rgba(176, 177, 180, 0.4); /* Very light gray with more transparency */
    transition: transform 0.3s, opacity 0.2s;
}

.collapsible.collapsed:after {
    transform: rotate(-90deg);
}

.collapsible-content {
    max-height: 1000px;
    overflow: hidden;
    transition: max-height 0.3s ease-out;
}

.collapsible-content.collapsed {
    max-height: 0;
    padding-top: 0;
    padding-bottom: 0;
    margin: 0;
    overflow: hidden;
}

/* Left Sidebar */
.error-analysis {
    width: 33.333%;
    padding: 1rem;
    background-color: white;
    border-right: 1px solid #e5e7eb;
    overflow-y: auto;
}

/* Exercise Description */
.exercise-box {
    background-color: #f0f9ff;
    padding: 1rem;
    border-radius: 0.375rem;
    border: 1px solid #bae6fd;
    margin-bottom: 1.5rem;
    position: relative;
    padding-top: 25px; /* Increased to make room for the title */
}

.exercise-box-header {
    font-weight: bold;
    font-size: 0.75rem;
    position: absolute;
    top: 10px;
    left: 10px;
    color: #0369a1;
}

/* Remove the pseudo-element since we now have an explicit header */
.exercise-box::before {
    content: none;
}

.section-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 0.75rem;
    display: flex;
    align-items: center;
}

.section-title i {
    margin-right: 0.5rem;
}

.error-box {
    background-color: #fee2e2;
    padding: 0.75rem;
    border-radius: 0.375rem;
    border: 1px solid #fecaca;
    margin-bottom: 1.5rem;
    position: relative;
}

.error-box::before {
    content: 'ERROR PATTERNS';
    font-size: 0.75rem;
    font-weight: bold;
    position: absolute;
    top: -10px;
    left: 10px;
    background-color: #fee2e2;
    padding: 0 6px;
    color: #b91c1c;
    border-radius: 4px;
}

/* ==== OUTER FRAME ==================================================== */
.code-box{
    background:#f3f4f6;
    border:1px solid #d1d5db;
    border-radius:.375rem;
    position:relative;

    padding:.75rem;          /* normal inner padding               */
    overflow:visible;        /* so the label can hang over the edge*/
    margin-bottom: 1.5rem;
}

/* label pill */
.code-box::before{
    content:"STUDENT CODE";
    position:absolute; 
    top:-10px; 
    left:10px;
    font:600 .75rem/1 sans-serif;
    padding:0 6px;
    border-radius:4px;
    background:#f3f4f6;
    color:#4b5563;
    pointer-events:none;
}

/* ==== SCROLLABLE AREA =============================================== */
.code-scroll{
    /* scrolling */
    overflow:auto;           /* both axes, only when needed         */
    max-height:24rem;        /* optional vertical cap               */

    /* code look‑&‑feel */
    font-family:monospace;
    white-space:pre;         /* preserve indentation, no wrapping   */
    line-height:1.35;
    /* optional extra styles … */
}

.reasoning-box {
    background-color: #eff6ff;
    padding: 0.75rem;
    border-radius: 0.375rem;
    border: 1px solid #dbeafe;
    position: relative;
    padding-top: 25px; /* Increased to make room for the title */
}

.reasoning-box-header {
    font-weight: bold;
    font-size: 0.75rem;
    position: absolute;
    top: 10px;
    left: 10px;
    color: #1d4ed8;
}

/* Remove the pseudo-element since we now have an explicit header */
.reasoning-box::before {
    content: none;
}

/* Tooltips - Position below icon */
.tooltip {
    position: relative;
    display: inline-block;
    cursor: help;
    margin-left: 8px;
}

.tooltip .tooltiptext {
    visibility: hidden;
    width: 300px;
    background-color: #605f5f;
    color: #fff;
    text-align: center;
    border-radius: 6px;
    padding: 8px;
    position: absolute;
    z-index: 9999; /* High z-index to be above other elements */
    top: 125%; /* Position below instead of above */
    left: 50%;
    margin-left: -100px;
    opacity: 0;
    transition: opacity 0.3s;
    font-size: 0.875rem;
    font-weight: normal;
}

.tooltip .tooltiptext::after {
    content: "";
    position: absolute;
    bottom: 100%; /* Changed from top to bottom */
    left: 50%;
    /* margin-left: -5px;
    border-width: 5px; */
    border-style: solid;
    border-color: transparent transparent #333 transparent; /* Changed arrow direction */
}

.tooltip:hover .tooltiptext {
    visibility: visible;
    opacity: 1;
}

/* Scaffolding Section */
.scaffolding-section {
    width: 66.667%;
    padding: 1rem;
    overflow-y: auto;
    position: relative; /* Added to ensure proper z-index context */
}

.scaffolds-title {
    display: flex;
    align-items: center;
    margin-bottom: 1rem;
    font-size: 1.25rem;
    font-weight: 600;
    position: relative; /* Added to create stacking context */
}

.scaffolds-subtitle {
    margin-left: 0.5rem;
    font-size: 0.875rem;
    color: #6b7280;
    font-weight: normal;
}

.scaffold-grid {
    display: grid;
    gap: 1rem;
}

.scaffold-card {
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    transition: all 0.2s ease-in-out;
    background-color: white;
}

.scaffold-card:hover {
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
}

.scaffold-card.selected {
    border: 2px solid #3b82f6;
}

.scaffold-header {
    padding: 1rem;
    cursor: pointer;
    border-radius: 0.5rem;
    position: relative;
}

.scaffold-title-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
}

.scaffold-name {
    font-weight: 600;
    font-size: 1.125rem;
    display: flex;
    align-items: center;
}

.scaffold-badges {
    display: flex;
    gap: 0.5rem;
}

.badge {
    padding: 0.25rem 0.5rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    display: flex;
    align-items: center;
    gap: 4px;
}

/* Badge Colors */
.badge-effort-low {
    background-color: #d1fae5;
    color: #065f46;
}

.badge-effort-moderate {
    background-color: #fef3c7;
    color: #92400e;
}

.badge-effort-high {
    background-color: #ffedd5;
    color: #9a3412;
}

.badge-action-low {
    background-color: #dbeafe;
    color: #1e40af;
}

.badge-action-moderate {
    background-color: #f3e8ff;
    color: #6b21a8;
}

.badge-action-high {
    background-color: #e0e7ff;
    color: #3730a3;
}

.scaffold-support {
    font-size: 0.875rem;
    margin-bottom: 0.75rem;
    line-height: 1.4;
}

.scaffold-metrics {
    font-size: 0.875rem;
    color: #4b5563;
    line-height: 1.4;
}

.scaffold-detail {
    border-top: 1px solid #e5e7eb;
    padding: 1rem;
    background-color: #f9fafb;
    display: none;
    animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
}

.detail-title {
    font-weight: 500;
    margin-bottom: 0.5rem;
}

.content-box {
    background-color: white;
    padding: 0.75rem;
    border-radius: 0.25rem;
    border: 1px solid #e5e7eb;
    margin-bottom: 1rem;
    line-height: 1.5;
}

.content-box pre {
    font-family: monospace;
    white-space: pre-wrap;
    background-color: #f3f4f6;
    padding: 0.5rem;
    border-radius: 0.25rem;
    margin: 0.5rem 0;
    overflow-x: auto;
}

.content-box strong {
    font-weight: 600;
}

.scaffold-reasoning {
    font-size: 0.875rem;
    margin-bottom: 1rem;
    line-height: 1.5;
}

.action-row {
    display: flex;
    justify-content: space-between;
    margin-top: 1rem;
}

.apply-button {
    background-color: #2563eb;
    color: white;
    padding: 0.5rem 1rem;
    border-radius: 0.375rem;
    border: none;
    display: flex;
    align-items: center;
    cursor: pointer;
    font-size: 0.875rem;
    font-weight: 500;
    transition: background-color 0.2s;
}

.apply-button:hover {
    background-color: #1d4ed8;
}

.apply-button i {
    margin-right: 0.5rem;
}

.copy-button {
    background-color: #4b5563;
    color: white;
    padding: 0.5rem 1rem;
    border-radius: 0.375rem;
    border: none;
    display: flex;
    align-items: center;
    cursor: pointer;
    font-size: 0.875rem;
    font-weight: 500;
    transition: background-color 0.2s;
}

.copy-button:hover {
    background-color: #374151;
}

.copy-button i {
    margin-right: 0.5rem;
}

/* Quick Nav */
.quick-nav {
    position: fixed;
    bottom: 50px;
    right: 20px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    z-index: 1000;
}

.quick-nav-btn {
    background-color: #3b82f6;
    color: white;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s, background-color 0.2s;
}

.quick-nav-btn:hover {
    background-color: #2563eb;
    transform: scale(1.05);
}

/* Footer */
footer {
    background-color: #e5e7eb;
    padding: 0.75rem;
    text-align: center;
    color: #4b5563;
    font-size: 0.875rem;
    border-top: 1px solid #d1d5db;
}

/* Responsive styles */
@media (max-width: 768px) {
    .main-content {
        flex-direction: column;
    }
    
    .error-analysis, .scaffolding-section {
        width: 100%;
    }
    
    .timer-container {
        position: static;
        margin: 10px auto;
        width: 90%;
    }
}

/* Notification */
.notification {
    position: fixed;
    top: 20px;
    left: 50%;
    transform: translateX(-50%);
    background-color: #10b981;
    color: white;
    padding: 10px 20px;
    border-radius: 4px;
    z-index: 1000;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    opacity: 0;
    transition: opacity 0.3s;
}

.notification.show {
    opacity: 1;
}

/* Previous Scaffold Box */
.previous-scaffold-box {
    background-color: #ede9fe;
    padding: 0.75rem;
    border-radius: 0.375rem;
    border: 1px solid #cbd5e1;
    margin-bottom: 1.5rem;
    position: relative;
    display: none; /* Hide by default, will show only if there's content */
}

.previous-scaffold-box::before {
    content: 'PREVIOUS SCAFFOLD';
    font-size: 0.75rem;
    font-weight: bold;
    position: absolute;
    top: -10px;
    left: 10px;
    background-color: #ede9fe;
    padding: 0 6px;
    color: #7c3aed;
    border-radius: 4px;
}

.previous-scaffold-content {
    font-size: 0.875rem;
    line-height: 1.5;
}
</style>
<body>
    <!-- Header -->
    <header>
        <div class="header-content">
            <h1 class="header-title">RapidResponse</h1>
            <div class="header-info">
            </div>
			<button id="regenerate-btn" style="margin: 0 20px; margin-right: 280px; padding: 8px 14px; font-size: 14px;">
      🔄 Regenerate Response
    </button>
            <div class="timer-container">
                <div class="timer-display" id="timerDisplay">00:00</div>
            </div>
        </div>
    </header>

    <!-- Notification -->
    <div class="notification" id="notification"></div>

    <!-- Main Content -->
    <div class="main-content">
        <!-- Left Sidebar - Error Analysis -->
        <div class="error-analysis">
            <!-- Collapsible Exercise Box -->
            <div class="exercise-box collapsible collapsed" id="exerciseBox">
                <div class="exercise-box-header">EXERCISE</div>
                <div class="collapsible-content collapsed" id="exerciseDescription"></div>
            </div>
            
            <div class="code-box" id="studentCode"><pre class="code-scroll"></pre></div>
            <div class="error-box" id="errorPatterns"></div>
            <div class="previous-scaffold-box" id="previousScaffoldBox"></div>
            
            <!-- Collapsible Reasoning Box -->
            <div class="reasoning-box collapsible collapsed" id="reasoningBox">
                <div class="reasoning-box-header">PEDAGOGICAL REASONING</div>
                <div class="collapsible-content collapsed" id="pedagogicalReasoning"></div>
            </div>
        </div>

        <!-- Center - Scaffolding Options -->
        <div class="scaffolding-section">
            <div>
                <h2 class="scaffolds-title">
                    Scaffolding Strategies
                    <span class="scaffolds-subtitle">(ordered from lowest to highest effort)</span>
                    <div class="tooltip">
                        <i data-feather="info" size="16"></i>
                        <span class="tooltiptext">Select the appropriate scaffold based on student needs and time constraints</span>
                    </div>
                </h2>
            </div>

            <div class="scaffold-grid" id="scaffoldGrid">
                <!-- Scaffold cards will be inserted here by JavaScript -->
            </div>
        </div>
    </div>

    <!-- Quick Nav Buttons -->
    <div class="quick-nav">
        <div class="quick-nav-btn" id="scrollToTop" title="Scroll to Top">
            <i data-feather="arrow-up"></i>
        </div>
    </div>

    <!-- Footer -->
    <footer>
        Crafting Scaffolds to Bridge Learning Gaps
    </footer>

    <script>
document.getElementById('regenerate-btn').addEventListener('click', function () {
	const userConfirmed = window.confirm('Are you sure you want to regenerate the Response?');
    if (!userConfirmed) return;

    const params = new URLSearchParams(window.location.search);
    const studentID = params.get('student_id');
    const problemID = params.get('problem_id');
    const snapshotId = params.get('snapshot_id');
    const userId = params.get('uid');
    const userRole = params.get('role');

	const button = this;
    button.disabled = true;
    const originalText = button.textContent;
    button.textContent = '⏳ Regenerating...';


    const actionUrl = '/sc_view?student_id=' + encodeURIComponent(studentID) +
                      '&problem_id=' + encodeURIComponent(problemID) +
                      '&snapshot_id=' + encodeURIComponent(snapshotId) +
                      '&uid=' + encodeURIComponent(userId) +
                      '&role=' + encodeURIComponent(userRole)

    const tempForm = document.createElement('form');
    tempForm.method = 'POST';
    tempForm.action = actionUrl;

    const hiddenFields = {
      student_id: studentID,
      problem_id: problemID,
      snapshot_id: snapshotId,
      uid: userId,
      role: userRole,
      isNew: 'true'
    };

    for (const key in hiddenFields) {
      const input = document.createElement('input');
      input.type = 'hidden';
      input.name = key;
      input.value = hiddenFields[key];
      tempForm.appendChild(input);
    }

    document.body.appendChild(tempForm);
    tempForm.submit();
  });

// Initialize Feather Icons
document.addEventListener('DOMContentLoaded', () => {
    feather.replace();
});

// Format content with Markdown-like syntax
function formatContent(content) {
    const backtick = String.fromCharCode(96);
    const tripleBacktick = backtick + backtick + backtick;
    const codeBlockRegex = new RegExp(tripleBacktick + '(.*?)' + tripleBacktick, 'gs');

    return content
        .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
        .replace(/\*(.*?)\*/g, '<em>$1</em>')
        .replace(codeBlockRegex, '<pre>$1</pre>')
        .replace(/\n/g, '<br>');
}


// Stopwatch functionality (replaced timer)
function initializeTimer() {
    let timerDisplay = document.getElementById('timerDisplay');
    let startButton = document.getElementById('startTimer');
    let pauseButton = document.getElementById('pauseTimer');
    let resetButton = document.getElementById('resetTimer');

    let totalSeconds = 0; // Start from 0 for a stopwatch
    let timerInterval;
    let isRunning = false;

    function updateDisplay() {
        const minutes = Math.floor(totalSeconds / 60);
        const seconds = totalSeconds % 60;
        timerDisplay.textContent = minutes.toString().padStart(2, '0') + ':' + seconds.toString().padStart(2, '0');

        // Add highlight class when more than 10 minutes have passed
        if (totalSeconds >= 600) { // 10 minutes = 600 seconds
            timerDisplay.classList.add('timer-warning');
        } else {
            timerDisplay.classList.remove('timer-warning');
        }
    }

    function startTimer() {
        if (!isRunning) {
            isRunning = true;
            if (startButton) startButton.disabled = true;
            if (pauseButton) pauseButton.disabled = false;

            timerInterval = setInterval(() => {
                totalSeconds++;
                updateDisplay();

                // Optional: Notify when reaching 10 minutes
                if (totalSeconds === 600) {
                    showNotification("10 minutes have passed in this session.");
                }
            }, 1000);
        }
    }

    function pauseTimer() {
        clearInterval(timerInterval);
        isRunning = false;
        if (startButton) startButton.disabled = false;
        if (pauseButton) pauseButton.disabled = true;
    }

    function resetTimer() {
        clearInterval(timerInterval);
        isRunning = false;
        totalSeconds = 0;
        updateDisplay();
        if (startButton) startButton.disabled = false;
        if (pauseButton) pauseButton.disabled = true;
        timerDisplay.classList.remove('timer-warning');
    }

    // Event listeners - only add if buttons exist
    if (startButton) startButton.addEventListener('click', startTimer);
    if (pauseButton) pauseButton.addEventListener('click', pauseTimer);
    if (resetButton) resetButton.addEventListener('click', resetTimer);

    // Initialize display
    updateDisplay();

    // Start automatically regardless of button presence
    setTimeout(() => {
        startTimer();
    }, 1500);
}

// Notification function
function showNotification(message) {
    const notification = document.getElementById('notification');
    notification.textContent = message;
    notification.classList.add('show');

    setTimeout(() => {
        notification.classList.remove('show');
    }, 3000);
}

// Populate the dashboard with data
function populateDashboard() {
    // Fill error analysis section
    document.getElementById('exerciseDescription').textContent = studentData.exerciseDescription;
    document.getElementById('errorPatterns').textContent = studentData.errorPatterns;
    // document.getElementById('studentCode').textContent = studentData.studentCode;
    const codeBox  = document.getElementById('studentCode');
    let codePre  = codeBox.querySelector('.code-scroll');
    if (!codePre) {                         // first time the page runs
        codePre           = document.createElement('pre');
        codePre.className = 'code-scroll';
        codeBox.appendChild(codePre);
    }

codePre.textContent = studentData.studentCode;

    document.getElementById('pedagogicalReasoning').textContent = studentData.reasoning;
    const previousScaffoldBox = document.getElementById('previousScaffoldBox');
    if (studentData.previousScaffold && studentData.previousScaffold.trim() !== '') {
        previousScaffoldBox.style.display = 'block';

        // Clear any existing content first
        previousScaffoldBox.innerHTML = '';

        // Create and append the content div
        const contentDiv = document.createElement('div');
        contentDiv.className = 'previous-scaffold-content';
        contentDiv.innerHTML = formatContent(studentData.previousScaffold);
        previousScaffoldBox.appendChild(contentDiv);
    } else {
        previousScaffoldBox.style.display = 'none';
    }

    // Create scaffold cards
    const scaffoldGrid = document.getElementById('scaffoldGrid');

    studentData.scaffolds.forEach((scaffold, index) => {
        // Create main card element
        const card = document.createElement('div');
        card.className = 'scaffold-card';
        card.dataset.index = index;

        // Create header section
        const header = document.createElement('div');
        header.className = 'scaffold-header';

        // Title row with badges
        const titleRow = document.createElement('div');
        titleRow.className = 'scaffold-title-row';

        const name = document.createElement('h3');
        name.className = 'scaffold-name';
        name.innerHTML = scaffold.strategy + ' <span style="font-size: 0.8rem; color: #6b7280; margin-left: 8px;">(~' + scaffold.estimatedTime + ')</span>';

        const badges = document.createElement('div');
        badges.className = 'scaffold-badges';

        const effortBadge = document.createElement('span');
        effortBadge.className = 'badge badge-effort-' + scaffold.effortLevel.toLowerCase();
        effortBadge.innerHTML = '<i data-feather="activity" size="14"></i> ' + scaffold.effortLevel + ' Effort';

        const actionBadge = document.createElement('span');
        actionBadge.className = 'badge badge-action-' + scaffold.actionabilityLevel.toLowerCase();
		actionBadge.innerHTML = '<i data-feather="zap" size="14"></i> ' + scaffold.actionabilityLevel + ' Actionability';

        badges.appendChild(effortBadge);
        badges.appendChild(actionBadge);
        titleRow.appendChild(name);
        titleRow.appendChild(badges);

        // Support and metrics sections
        const support = document.createElement('div');
        support.className = 'scaffold-support';
        support.innerHTML = '<strong>Learning Support:</strong> ' + scaffold.learningSupport;

        const metrics = document.createElement('div');
        metrics.className = 'scaffold-metrics';
        metrics.innerHTML = '<strong>Success Metrics:</strong> ' + scaffold.implementationSuccessMetrics;

        // Assemble header
        header.appendChild(titleRow);
        header.appendChild(support);
        header.appendChild(metrics);

        // Create detail section (hidden by default)
        const detail = document.createElement('div');
        detail.className = 'scaffold-detail';
        detail.id = 'scaffold-detail-' + index;

        // Content to share
        const contentTitle = document.createElement('h4');
        contentTitle.className = 'detail-title';
        contentTitle.textContent = 'Content to Share with Student:';

        const contentBox = document.createElement('div');
        contentBox.className = 'content-box';
        contentBox.innerHTML = formatContent(scaffold.content);

        // Reasoning
        const reasoningTitle = document.createElement('h4');
        reasoningTitle.className = 'detail-title';
        reasoningTitle.textContent = 'Recommendation Reasoning:';

        const reasoning = document.createElement('p');
        reasoning.className = 'scaffold-reasoning';
        reasoning.textContent = scaffold.recommendationReasoning;

        // Action buttons
        const actionRow = document.createElement('div');
        actionRow.className = 'action-row';

        const copyButton = document.createElement('button');
        copyButton.className = 'copy-button';
        copyButton.innerHTML = '<i data-feather="copy"></i> Copy to Clipboard';
        copyButton.addEventListener('click', (e) => {
            e.stopPropagation();

            // Replace deprecated document.execCommand with Clipboard API
            const textToCopy = scaffold.content.replace(/\*\*/g, '').replaceAll(String.fromCharCode(96) + String.fromCharCode(96) + String.fromCharCode(96), '');

// Use the modern Clipboard API
navigator.clipboard.writeText(textToCopy)
.then(() => {
showNotification('Scaffold content copied to clipboard!');
})
.catch(err => {
console.error('Could not copy text: ', err);
showNotification('Failed to copy to clipboard. Please try again.');
});
});

const applyButton = document.createElement('button');
applyButton.className = 'apply-button';
applyButton.innerHTML = '<i data-feather="arrow-right-circle"></i> Use This Scaffold';
applyButton.addEventListener('click', (e) => {
  e.stopPropagation();
  showNotification('Preparing scaffold for editing...');

  const editableContent = formatContent(scaffold.content);

  // Create popup overlay
  const overlay = document.createElement('div');
  overlay.className = 'modal-overlay';
  overlay.style.cssText = 
  'position: fixed; top: 0; left: 0; width: 100%; height: 100%;' +
  'background: rgba(0,0,0,0.5); display: flex; justify-content: center; align-items: center; z-index: 1000;';


  // Create popup modal
  const modal = document.createElement('div');
  modal.className = 'modal-box';
  modal.style.cssText = 
  'background: white; padding: 20px; border-radius: 10px; width: 600px; max-width: 90%;' +
  'box-shadow: 0 0 10px rgba(0,0,0,0.3);';


  const title = document.createElement('h4');
  title.textContent = 'Edit Scaffold Content';
  title.style.marginBottom = '10px';

  const textarea = document.createElement('textarea');
  textarea.style.cssText = 'width: 100%; height: 200px; margin-bottom: 10px;';
  textarea.value = editableContent.replace(/<br\s*\/?>/g, '\n');

  const confirmButton = document.createElement('button');
  confirmButton.textContent = 'Confirm & Send';
  confirmButton.style.cssText = 'margin-right: 10px;';

  const cancelButton = document.createElement('button');
  cancelButton.textContent = 'Cancel';

  // Confirm button event
  confirmButton.addEventListener('click', async () => {
    const editedContent = textarea.value.replace(/\n/g, '<br>'); // convert newlines to HTML
    const urlParams = new URLSearchParams(window.location.search);
    const studentId = parseInt(urlParams.get('student_id'));
    const problemId = parseInt(urlParams.get('problem_id'));
    const snapshotId = parseInt(urlParams.get('snapshot_id'));
    const authorId = parseInt(urlParams.get('uid'));
    const authorRole = urlParams.get('role');

    const timerText = document.getElementById('timerDisplay')?.textContent || "00:00";
    const [minutes, seconds] = timerText.split(':').map(Number);
    const duration = (minutes * 60) + seconds;

    try {
      const response = await fetch('/assign_scaffold', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          student_id: studentId,
          problem_id: problemId,
          duration: duration,
          scaffolding: editedContent,
          snapshot_id: snapshotId,
          uid: authorId,
          role: authorRole
        }),
      });

      if (!response.ok) throw new Error('Failed to assign scaffold');
      showNotification('Scaffold assigned successfully!');
    } catch (error) {
      console.error(error);
      showNotification('Error assigning scaffold.');
    }

    document.body.removeChild(overlay);
  });

  cancelButton.addEventListener('click', () => {
    document.body.removeChild(overlay);
  });

  modal.appendChild(title);
  modal.appendChild(textarea);
  modal.appendChild(confirmButton);
  modal.appendChild(cancelButton);
  overlay.appendChild(modal);
  document.body.appendChild(overlay);
});




actionRow.appendChild(copyButton);
actionRow.appendChild(applyButton);

// Assemble detail section
detail.appendChild(contentTitle);
detail.appendChild(contentBox);
detail.appendChild(reasoningTitle);
detail.appendChild(reasoning);
detail.appendChild(actionRow);

// Assemble full card
card.appendChild(header);
card.appendChild(detail);
scaffoldGrid.appendChild(card);

// Add click event to header
header.addEventListener('click', () => {
// Toggle selection class
document.querySelectorAll('.scaffold-card').forEach(c => {
c.classList.remove('selected');
});
card.classList.add('selected');

// Hide all details
document.querySelectorAll('.scaffold-detail').forEach(d => {
d.style.display = 'none';
});

// Show this detail
detail.style.display = 'block';

// Re-initialize feather icons for any new elements
feather.replace();
});
});

// Re-initialize feather icons
feather.replace();
}

// Initialize quick nav functionality
function initializeQuickNav() {
const scrollToTopBtn = document.getElementById('scrollToTop');

scrollToTopBtn.addEventListener('click', () => {
window.scrollTo({
top: 0,
behavior: 'smooth'
});
});

// Show/hide button based on scroll position
window.addEventListener('scroll', () => {
if (window.pageYOffset > 300) {
scrollToTopBtn.style.display = 'flex';
} else {
scrollToTopBtn.style.display = 'none';
}
});

// Initial check
if (window.pageYOffset <= 300) {
scrollToTopBtn.style.display = 'none';
}
}

function initializeCollapsibles() {
document.querySelectorAll('.collapsible').forEach(collapsible => {
collapsible.addEventListener('click', function() {
this.classList.toggle('collapsed');

// Find the content element
const content = this.querySelector('.collapsible-content');
if (content) {
content.classList.toggle('collapsed');
}
});
});
}

// Initialize everything
document.addEventListener('DOMContentLoaded', () => {
populateDashboard();
initializeTimer();
initializeQuickNav();
initializeCollapsibles();

// Pre-select the first scaffold
setTimeout(() => {
const firstScaffold = document.querySelector('.scaffold-header');
if (firstScaffold) {
firstScaffold.click();
}
}, 500);
});
</script>
    <script>
// Student data (from paste.txt)
const studentData = {{.Text}};
</script>
</body>
</html>
`

var VIS_TEMPLATE = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Student Coding Performance Dashboard</title>
    <link rel="stylesheet" href="assets/css/styles.css">
    <link rel="stylesheet" href="assets/css/concept_map.css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/3.9.1/chart.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/prism.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-python.min.js"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-tomorrow.min.css">
</head>
<style>
/* Concept Map Container */
.concept-map-container {
    height: 1000px;
    width: 100%;
    position: relative;
    margin-top: 1rem;
    border: 1px solid #e3e6f0;
    border-radius: 4px;
    overflow: hidden;
    background: rgba(248, 249, 252, 0.5);
  }
  
  /* Simple Concept Map */
  .simple-concept-map {
    position: relative;
    width: 100%;
    height: 100%;
    overflow: auto;
  }
  
  /* Hub in center */
  .concept-hub {
    position: absolute;
    width: 120px;
    height: 120px;
    background-color: #4e73df;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    left: 400px;
    top: 250px;
    transform: translate(-50%, -50%);
    z-index: 2;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    border: 2px solid #3a5ccc;
  }
  
  .concept-hub span {
    color: white;
    font-weight: bold;
    text-align: center;
    font-size: 14px;
    padding: 10px;
  }
  
  /* Nodes */
  .concept-node {
    position: absolute;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transform: translate(-50%, -50%);
    z-index: 2;
    cursor: pointer;
    transition: box-shadow 0.2s ease;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  }
  
  .concept-node:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  }
  
  /* Misconception nodes */
  .misconception-node {
    width: 170px;
    height: 170px;
    border: 2px solid rgba(0, 0, 0, 0.2);
  }
  
  /* Error nodes */
  .error-node {
    width: 100px;
    height: 100px;
    background-color: rgb(100, 130, 220);
    /* background-color: hsl(15, 85%, 55%); */
    border: 2px solid #3a5ccc;
  }
  
  /* Node content */
  .node-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    padding: 5px;
    text-align: center;
  }
  
  .node-title {
    font-size: 12px;
    font-weight: bold;
    line-height: 1.2;
    max-height: 70%;
    overflow: hidden;
    display: -webkit-box;
    display: -ms-flexbox;
    display: box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    text-overflow: ellipsis;
}
  
  .error-node .node-title {
    color: white;
  }
  
  .node-percentage {
    font-size: 12px;
    font-weight: bold;
    margin-top: 5px;
  }
  
  /* Connections */
  .connection {
    position: absolute;
    height: 2px;
    background-color: rgba(128, 128, 128, 0.3);
    transform-origin: left center;
    z-index: 1;
  }
  
  .hub-connection {
    background-color: rgba(78, 115, 223, 0.4);
  }
  
  .node-connection {
    background-color: rgba(128, 128, 128, 0.3);
  }
  
  .connection-highlight {
    height: 3px;
    background-color: rgba(28, 200, 138, 0.8);
    z-index: 2;
  }
  
  /* Legend */
  .concept-map-legend {
    position: absolute;
    top: 20px;
    left: 20px;
    background-color: rgba(255, 255, 255, 0.9);
    padding: 10px;
    border-radius: 4px;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  }
  
  .legend-title {
    font-weight: bold;
    margin-bottom: 10px;
    font-size: 14px;
  }
  
  .legend-item {
    display: flex;
    align-items: center;
    margin-bottom: 5px;
  }
  
  .legend-icon {
    width: 16px;
    height: 16px;
    border-radius: 50%;
    margin-right: 10px;
  }
  
  .misconception-icon {
    background-color: hsl(185, 100%, 60%); /* Medium brightness cyan */
    border: 1px solid rgba(0, 0, 0, 0.2);
}

  
  .error-icon {
    background-color: rgba(78, 115, 223, 0.7);
    border: 1px solid #3a5ccc;
  }
  
  .legend-label {
    font-size: 12px;
  }
  
  /* Tooltip */
  .concept-tooltip {
    position: absolute;
    background-color: rgba(0, 0, 0, 0.8);
    color: white;
    padding: 10px;
    border-radius: 4px;
    font-size: 12px;
    max-width: 250px;
    z-index: 1000;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
    pointer-events: none;
  }
  
  /* Dark mode support */
  .dark-mode .concept-map-container {
    background: rgba(26, 32, 44, 0.3);
    border-color: #2d3748;
  }
  
  .dark-mode .concept-hub {
    border-color: #2c4687;
  }
  
  .dark-mode .error-node {
    border-color: #2c4687;
  }
  
  .dark-mode .concept-map-legend {
    background-color: rgba(26, 32, 44, 0.8);
    color: #e2e8f0;
  }
  
  .dark-mode .legend-label {
    color: #e2e8f0;
  }
  
  .dark-mode .node-title {
    color: #333;
  }
  
  .dark-mode .error-node .node-title {
    color: #fff;
  }
  
  .dark-mode .concept-tooltip {
    background-color: rgba(26, 32, 44, 0.9);
  }

  /* Concept Map Controls */
.concept-map-controls {
  background-color: rgba(248, 249, 252, 0.9);
  border: 1px solid #e3e6f0;
  border-radius: 4px;
  padding: 15px;
  margin-bottom: 15px;
}

.control-section h4 {
  margin-top: 0;
  margin-bottom: 10px;
  font-size: 16px;
  color: #5a5c69;
}

.control-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
  margin-bottom: 15px;
}

.control-item {
  display: flex;
  flex-direction: column;
}

.control-item label {
  margin-bottom: 5px;
  font-size: 14px;
}

.control-item input[type="range"] {
  width: 100%;
}

.reset-button {
  background-color: #4e73df;
  color: white;
  border: none;
  padding: 8px 15px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.reset-button:hover {
  background-color: #2e59d9;
}

/* Dark mode support */
.dark-mode .concept-map-controls {
  background-color: rgba(26, 32, 44, 0.8);
  border-color: #2d3748;
  color: #e2e8f0;
}

.dark-mode .control-section h4 {
  color: #e2e8f0;
}

.dark-mode .reset-button {
  background-color: #3151b7;
}

.dark-mode .reset-button:hover {
  background-color: #263d8f;
}

:root {
    --primary-color: #4e73df;
    --secondary-color: #1cc88a;
    --warning-color: #f6c23e;
    --danger-color: #e74a3b;
    --info-color: #36b9cc;
    --dark-color: #5a5c69;
    --light-color: #f8f9fc;
    --gray-100: #f8f9fc;
    --gray-200: #eaecf4;
    --gray-300: #dddfeb;
    --gray-400: #d1d3e2;
    --gray-500: #b7b9cc;
    --gray-600: #858796;
    --gray-700: #6e707e;
    --gray-800: #5a5c69;
    --gray-900: #3a3b45;
    --card-border-radius: 0.5rem;
    --transition-speed: 0.3s;
    --sidebar-width: 16rem;
    --sidebar-collapsed-width: 4.5rem;
    --font-primary: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
    font-family: var(--font-primary);
}

body {
    background-color: var(--gray-100);
    color: var(--gray-800);
    transition: background-color var(--transition-speed);
    line-height: 1.5;
}

/* Dark Mode Styles */
.dark-mode {
    background-color: #1a1a2e;
    color: #e1e1e1;
}

.dark-mode .section {
    background-color: #16213e;
    box-shadow: 0 0.15rem 1.75rem 0 rgba(0, 0, 0, 0.2);
}

.dark-mode .section-header {
    border-bottom: 1px solid #2a2d45;
}

.dark-mode .section-header h2 {
    color: #e1e1e1;
}

.dark-mode .collapse-btn {
    color: #e1e1e1;
}

.dark-mode .metric-card {
    background-color: #16213e;
    box-shadow: 0 0.15rem 1.75rem 0 rgba(0, 0, 0, 0.2);
}

.dark-mode .metric-card h3 {
    color: #e1e1e1;
}

.dark-mode .metric-card .value {
    color: #e1e1e1;
}

.dark-mode .metric-card .percentage {
    color: #b1b1b3;
}

.dark-mode .info-card {
    background-color: #16213e;
    box-shadow: 0 0.15rem 1.75rem 0 rgba(0, 0, 0, 0.2);
}

.dark-mode .info-card h3 {
    color: #e1e1e1;
    border-bottom: 1px solid #2a2d45;
}

.dark-mode .tabs {
    border-bottom: 1px solid #2a2d45;
}

.dark-mode .tab.active {
    border-bottom: 2px solid #36b9cc;
    color: #36b9cc;
}

.dark-mode .tab:hover:not(.active) {
    background-color: #2a2d45;
    border-bottom: 2px solid #16213e;
}

.dark-mode .code-preview {
    background-color: #2a2d45;
}

.dark-mode .tag {
    background-color: #2a2d45;
    color: #e1e1e1;
}

.dark-mode .correlation-matrix th, .dark-mode .correlation-matrix td {
    border: 1px solid #2a2d45;
}

.dark-mode .correlation-matrix th {
    background-color: #2a2d45;
}

.dark-mode .sidebar {
    background-image: linear-gradient(180deg, #1a1a2e 10%, #0f3460 100%);
}

.dark-mode .highlight-item {
    background-image: linear-gradient(180deg, #1a1a2e 10%, #0f3460 100%);
    box-shadow: 0 0.15rem 1.75rem 0 rgba(0, 0, 0, 0.2);
}

.dark-mode .strategyList li:before {
    color: #36b9cc;
}

/* Layout Styles */
.dashboard-container {
    display: flex;
    min-height: 100vh;
}

.sidebar {
    width: var(--sidebar-width);
    background-color: var(--primary-color);
    background-image: linear-gradient(180deg, var(--primary-color) 10%, #224abe 100%);
    color: white;
    transition: all var(--transition-speed);
    box-shadow: 0 0.15rem 1.75rem 0 rgba(58, 59, 69, 0.15);
    z-index: 10;
    position: fixed;
    left: 0;
    top: 0;
    height: 100vh;
    overflow-y: auto;
    padding: 1rem 0;
}

.sidebar h1 {
    font-size: 1.2rem;
    text-align: center;
    padding: 1rem;
    margin-bottom: 2rem;
    text-transform: uppercase;
    letter-spacing: 0.05rem;
}

.nav-item {
    padding: 0.85rem 1.5rem;
    display: flex;
    align-items: center;
    cursor: pointer;
    transition: background-color 0.2s;
    position: relative;
}

.nav-item.active {
    background-color: rgba(255, 255, 255, 0.15);
}

.nav-item:hover {
    background-color: rgba(255, 255, 255, 0.1);
}

.nav-item i {
    margin-right: 0.75rem;
    width: 1.5rem;
    text-align: center;
}

.main-content {
    flex: 1;
    margin-left: var(--sidebar-width);
    padding: 1.5rem;
    transition: margin-left var(--transition-speed);
}

/* Toggle Switch Styles */
.toggle-wrapper {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    margin-bottom: 1.5rem;
}

.toggle-label {
    margin-right: 0.5rem;
    font-size: 0.875rem;
}

.toggle-switch {
    position: relative;
    display: inline-block;
    width: 3rem;
    height: 1.5rem;
}

.toggle-switch input {
    opacity: 0;
    width: 0;
    height: 0;
}

.toggle-slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: #ccc;
    transition: .4s;
    border-radius: 34px;
}

.toggle-slider:before {
    position: absolute;
    content: "";
    height: 1rem;
    width: 1rem;
    left: 0.25rem;
    bottom: 0.25rem;
    background-color: white;
    transition: .4s;
    border-radius: 50%;
}

input:checked + .toggle-slider {
    background-color: #2196F3;
}

input:checked + .toggle-slider:before {
    transform: translateX(1.5rem);
}

/* Section Styles */
.section {
    background-color: white;
    border-radius: var(--card-border-radius);
    box-shadow: 0 0.15rem 1.75rem 0 rgba(58, 59, 69, 0.15);
    margin-bottom: 1.5rem;
    transition: all var(--transition-speed);
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.25rem;
    border-bottom: 1px solid var(--gray-300);
}

.section-header h2 {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--gray-800);
}

.collapse-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    color: var(--gray-800);
    cursor: pointer;
}

.section-body {
    padding: 1.25rem;
    overflow: hidden;
    transition: max-height 0.3s;
}

.collapsed {
    max-height: 0;
    padding: 0 1.25rem;
}

/* Metric Cards Styles */
.metrics-cards {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
    margin-bottom: 1.5rem;
}

.metric-card {
    background-color: white;
    border-radius: var(--card-border-radius);
    padding: 1.25rem;
    text-align: center;
    box-shadow: 0 0.15rem 1.75rem 0 rgba(58, 59, 69, 0.15);
    border-left: 0.25rem solid var(--primary-color);
    transition: transform 0.3s;
}

.metric-card:hover {
    transform: translateY(-5px);
}

.metric-card h3 {
    font-size: 0.85rem;
    font-weight: 700;
    text-transform: uppercase;
    margin-bottom: 0.5rem;
    color: var(--gray-800);
}

.metric-card .value {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--gray-800);
    margin-bottom: 0.25rem;
}

.metric-card .percentage {
    font-size: 0.875rem;
    color: var(--gray-600);
}

.metric-card.success {
    border-left-color: var(--secondary-color);
}

.metric-card.success .value {
    color: var(--secondary-color);
}

.metric-card.warning {
    border-left-color: var(--warning-color);
}

.metric-card.warning .value {
    color: var(--warning-color);
}

.metric-card.danger {
    border-left-color: var(--danger-color);
}

.metric-card.danger .value {
    color: var(--danger-color);
}

.metric-card.info {
    border-left-color: var(--info-color);
}

.metric-card.info .value {
    color: var(--info-color);
}

/* Chart Container Styles */
.chart-container {
    display: flex;
    flex-wrap: wrap;
    gap: 1.5rem;
    margin-bottom: 1.5rem;
}

.chart-wrapper {
    flex: 1;
    min-width: 300px;
    height: 300px;
    position: relative;
}

/* Info Card Styles */
.card-wrapper {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: 1rem;
    margin-bottom: 1.5rem;
}

.info-card {
    background-color: white;
    border-radius: var(--card-border-radius);
    padding: 1.25rem;
    box-shadow: 0 0.15rem 1.75rem 0 rgba(58, 59, 69, 0.15);
    transition: all var(--transition-speed);
}

.info-card h3 {
    font-size: 1.1rem;
    font-weight: 600;
    margin-bottom: 0.75rem;
    color: var(--gray-800);
    border-bottom: 1px solid var(--gray-300);
    padding-bottom: 0.5rem;
}

.info-card h4 {
    font-size: 1rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
    color: var(--gray-800);
}

.info-card .description {
    margin-bottom: 1rem;
    font-size: 0.925rem;
    line-height: 1.5;
}

.code-preview {
    background-color: var(--gray-100);
    border-radius: 0.35rem;
    padding: 1rem;
    overflow-x: auto;
    margin-bottom: 1rem;
    font-family: monospace;
    font-size: 0.875rem;
    line-height: 1.5;
}

/* Tag Styles */
.tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 0.75rem;
}

.tag {
    background-color: var(--gray-300);
    color: var(--gray-800);
    border-radius: 1rem;
    padding: 0.25rem 0.75rem;
    font-size: 0.75rem;
    font-weight: 600;
}

.tag.primary {
    background-color: var(--primary-color);
    color: white;
}

.tag.success {
    background-color: var(--secondary-color);
    color: white;
}

.tag.warning {
    background-color: var(--warning-color);
    color: white;
}

.tag.danger {
    background-color: var(--danger-color);
    color: white;
}

.tag.info {
    background-color: var(--info-color);
    color: white;
}

/* Correlation Matrix Styles */
.heatmap-container {
    width: 100%;
    overflow-x: auto;
    margin-bottom: 1.5rem;
}

.correlation-matrix {
    width: 100%;
    border-collapse: collapse;
    margin-bottom: 1.5rem;
}

.correlation-matrix th, .correlation-matrix td {
    padding: 0.75rem;
    text-align: center;
    border: 1px solid var(--gray-300);
}

.correlation-matrix th {
    background-color: var(--gray-100);
    font-weight: 600;
}

.correlation-value {
    display: block;
    width: 100%;
    height: 100%;
    border-radius: 0.25rem;
    color: white;
    font-weight: 600;
    padding: 0.5rem;
}

/* Tab Styles */
.tabs {
    display: flex;
    margin-bottom: 1rem;
    border-bottom: 1px solid var(--gray-300);
}

.tab {
    padding: 0.75rem 1.25rem;
    cursor: pointer;
    transition: all 0.2s;
    border-bottom: 2px solid transparent;
}

.tab.active {
    border-bottom: 2px solid var(--primary-color);
    color: var(--primary-color);
    font-weight: 600;
}

.tab:hover:not(.active) {
    background-color: var(--gray-100);
    border-bottom: 2px solid var(--gray-300);
}

.tab-content {
    display: none;
}

.tab-content.active {
    display: block;
}

/* Priority Indicator Styles */
.priority-indicator {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    border-radius: 1rem;
    font-size: 0.75rem;
    font-weight: 600;
    color: white;
    margin-bottom: 0.5rem;
}

.priority-high {
    background-color: var(--danger-color);
}

.priority-medium {
    background-color: var(--warning-color);
}

.priority-low {
    background-color: var(--secondary-color);
}

/* Strategy List Styles */
.strategyList {
    list-style-type: none;
    margin-left: 0;
    padding-left: 0;
}

.strategyList li {
    position: relative;
    padding-left: 1.5rem;
    margin-bottom: 0.5rem;
    line-height: 1.5;
}

.strategyList li:before {
    content: "→";
    position: absolute;
    left: 0;
    color: var(--primary-color);
    font-weight: bold;
}

/* Highlights Container Styles */
.highlights-container {
    padding: 0.75rem 0;
}

.highlight-items {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1rem;
    margin-bottom: 1rem;
}

.secondary-highlights {
    margin-top: 1rem;
}

.highlight-item {
    background-color: var(--primary-color);
    background-image: linear-gradient(180deg, var(--primary-color) 10%, #224abe 100%);
    color: white;
    border-radius: var(--card-border-radius);
    padding: 1.5rem;
    box-shadow: 0 0.15rem 1.75rem 0 rgba(58, 59, 69, 0.15);
    transition: transform 0.3s;
}

.highlight-item:hover {
    transform: translateY(-5px);
}

.highlight-item h3 {
    font-size: 1.1rem;
    font-weight: 600;
    margin-bottom: 0.75rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
    padding-bottom: 0.5rem;
}

.highlight-item p {
    font-size: 0.925rem;
    line-height: 1.5;
}

.progress-container {
    background-color: rgba(255, 255, 255, 0.1);
    border-radius: 0.5rem;
    height: 0.5rem;
    overflow: hidden;
    margin: 0.5rem 0 1rem;
}

.progress-bar {
    height: 100%;
    background-color: var(--secondary-color);
    border-radius: 0.5rem;
}

.progress-label {
    display: flex;
    justify-content: space-between;
    font-size: 0.75rem;
    margin-bottom: 0.25rem;
}

/* Code Comparison Styles */
.code-comparison-container {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
}

@media (max-width: 768px) {
    .code-comparison-container {
        grid-template-columns: 1fr;
    }
}

.code-comparison {
    border: 1px solid var(--gray-300);
    border-radius: var(--card-border-radius);
    overflow: hidden;
}

.code-comparison-header {
    padding: 0.75rem;
    background-color: var(--gray-200);
    border-bottom: 1px solid var(--gray-300);
    font-weight: 600;
    display: flex;
    justify-content: space-between;
}

.code-comparison-body {
    padding: 1rem;
}

/* Timeline Styles */
.timeline-container {
    margin-top: 1.5rem;
}

.timeline-header {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1rem;
    margin-bottom: 1rem;
}

.timeline-phase {
    background-color: var(--gray-200);
    padding: 0.75rem;
    border-radius: var(--card-border-radius);
    text-align: center;
    font-weight: 600;
}

.timeline-content {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1rem;
}

.timeline-item {
    background-color: white;
    border-radius: var(--card-border-radius);
    padding: 1rem;
    box-shadow: 0 0.15rem 1.75rem 0 rgba(58, 59, 69, 0.1);
    border-left: 3px solid var(--primary-color);
}

.timeline-item.priority-high {
    border-left-color: var(--danger-color);
}

.timeline-item.priority-medium {
    border-left-color: var(--warning-color);
}

.timeline-item.priority-low {
    border-left-color: var(--secondary-color);
}


/* Comparative Analysis Styles */
.comparative-analysis {
    margin-top: 1.5rem;
}

.factor-list {
    list-style-type: none;
    margin: 0;
    padding: 0;
}

.factor-list li {
    padding: 0.5rem 0;
    border-bottom: 1px solid var(--gray-200);
    display: flex;
    align-items: center;
}

.factor-list li:last-child {
    border-bottom: none;
}

.factor-icon {
    width: 2rem;
    height: 2rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    margin-right: 1rem;
    flex-shrink: 0;
}

.factor-icon.success {
    background-color: rgba(28, 200, 138, 0.1);
    color: var(--secondary-color);
}

.factor-icon.danger {
    background-color: rgba(231, 74, 59, 0.1);
    color: var(--danger-color);
}

/* Action Plan Overview Styles */
.action-plan-overview {
    display: flex;
    gap: 1.5rem;
    margin-bottom: 1.5rem;
}

.action-summary {
    flex: 1;
    min-width: 300px;
}

/* Add to the CSS file */
.effort-indicator {
    display: inline-block;
    padding: 0.25rem 0.5rem;
    border-radius: 0.25rem;
    margin: 0.5rem 0;
    font-size: 0.875rem;
    background-color: #f8f9fc;
    border: 1px solid #d1d3e2;
}

.implementation-tag {
    display: inline-block;
    padding: 0.25rem 0.5rem;
    border-radius: 0.25rem;
    margin: 0.5rem 0;
    font-size: 0.875rem;
}

.implementation-tag.correct {
    background-color: #1cc88a;
    color: white;
}

/* Modal Styles - Add to style.css */
.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.85);
    z-index: 1000;
    display: flex;
    justify-content: center;
    align-items: center;
    opacity: 0;
    visibility: hidden;
    transition: opacity 0.3s ease, visibility 0.3s ease;
  }
  
  .modal-overlay.active {
    opacity: 1;
    visibility: visible;
  }
  
  .modal-content {
    background-color: var(--bg-color, #fff);
    border-radius: 8px;
    padding: 20px;
    max-width: 80%;
    max-height: 80vh;
    overflow-y: auto;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
    transform: scale(0.8);
    transition: transform 0.3s ease;
  }
  
  .modal-overlay.active .modal-content {
    transform: scale(1);
  }
  
  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
    border-bottom: 1px solid var(--border-color, #e3e6f0);
    padding-bottom: 10px;
  }
  
  .modal-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin: 0;
  }
  
  .modal-close {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1.5rem;
    color: var(--text-color, #5a5c69);
  }
  
  .modal-body {
    margin-bottom: 15px;
  }
  
  .code-preview {
    cursor: pointer;
    position: relative;
  }
  
  .code-preview:after {
    content: "🔍 Click to expand";
    position: absolute;
    top: 5px;
    right: 10px;
    font-size: 0.75rem;
    color: #fff;
    background-color: rgba(0, 0, 0, 0.5);
    padding: 2px 8px;
    border-radius: 4px;
    opacity: 0;
    transition: opacity 0.2s ease;
  }
  
  .code-preview:hover:after {
    opacity: 1;
  }
  
  .dark-mode .modal-content {
    background-color: #2c3136;
    color: #e3e6f0;
  }
  
  .dark-mode .modal-close {
    color: #e3e6f0;
  }

/* Responsive Styles */
@media (max-width: 992px) {
    .action-plan-overview {
        flex-direction: column;
    }
    
    .chart-wrapper {
        min-width: 100%;
    }
}

@media (max-width: 768px) {
    .sidebar {
        width: var(--sidebar-collapsed-width);
    }

    .sidebar h1, .nav-item span {
        display: none;
    }

    .nav-item i {
        margin-right: 0;
        font-size: 1.25rem;
    }

    .main-content {
        margin-left: var(--sidebar-collapsed-width);
    }
    
    .metrics-cards {
        grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    }
}

@media (max-width: 576px) {
    .metrics-cards {
        grid-template-columns: 1fr;
    }

    .chart-wrapper {
        min-width: 100%;
    }

    .card-wrapper {
        grid-template-columns: 1fr;
    }
    
    .timeline-header, .timeline-content {
        grid-template-columns: 1fr;
    }
    
    .main-content {
        padding: 1rem;
    }
}

.nav-item.disabled {
  pointer-events: none;
  opacity: 0.6;
}


</style>
<body>
    <div class="dashboard-container">
        <aside class="sidebar">
            <h1>Student Performance</h1>
            <div class="nav-item active" data-section="summary">
                <i>📊</i> <span>Summary Dashboard</span>
            </div>
            <div class="nav-item" data-section="performance">
                <i>📈</i> <span>Stage 1: Performance</span>
            </div>
            <div class="nav-item" data-section="errors">
                <i>❌</i> <span>Stage 2: Errors</span>
            </div>
            <div class="nav-item" data-section="correlations">
                <i>🔄</i> <span>Stage 3: Correlations</span>
            </div>
            <div class="nav-item" data-section="misconceptions">
                <i>💡</i> <span>Stage 4: Misconceptions</span>
            </div>
            <div class="nav-item" data-section="action-plan">
                <i>📝</i> <span>Stage 5: Action Plan</span>
            </div>
            <div class="nav-item" data-section="additional">
                <i>➕</i> <span>Stage 6: Insights</span>
            </div>
					<div class="nav-item" data-section="regenerate-btn" id="regenerate-sidebar">
  <i>🔄</i> <span>Regenerate Response</span>
</div>
        </aside>

        <main class="main-content">
            <div class="toggle-wrapper">
                <span class="toggle-label">Dark Mode</span>
                <label class="toggle-switch">
                    <input type="checkbox" id="darkModeToggle">
                    <span class="toggle-slider"></span>
                </label>
            </div>

            <!-- Summary Dashboard Section -->
            <section id="summary" class="section">
                <div class="section-header">
                    <h2>Key Findings Overview </h2>
                    <button class="collapse-btn" data-target="summary-body">−</button>
                </div>
                <div id="summary-body" class="section-body">
                    <div class="highlights-container">
                        <div class="highlight-items">
                            <div class="highlight-item">
                                <h3>Performance Summary</h3>
                                <p id="highlight-performance">Loading...</p>
                                <div class="progress-container">
                                    <div class="progress-bar" id="passing-progress" style="width: 0%"></div>
                                </div>
                                <div class="progress-label">
                                    <span>Satisfactory Performance</span>
                                    <span id="passing-percentage">0%</span>
                                </div>
                            </div>
                            <div class="highlight-item">
                                <h3>Top Error</h3>
                                <p id="highlight-top-error">Loading...</p>
                            </div>
                            <div class="highlight-item">
                                <h3>Primary Misconception</h3>
                                <p id="highlight-misconception">Loading...</p>
                            </div>
                        </div>
                        <div class="highlight-items secondary-highlights">
                            <div class="highlight-item">
                                <h3>Priority Intervention</h3>
                                <p id="highlight-intervention">Loading...</p>
                            </div>
                            <div class="highlight-item">
                                <h3>Student Success Gap</h3>
                                <p id="highlight-gap">Loading...</p>
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            <!-- Performance Analysis Section -->
            <section id="performance" class="section">
                <div class="section-header">
                    <h2>Stage 1: Performance Analysis</h2>
                    <button class="collapse-btn" data-target="performance-body">−</button>
                </div>
                <div id="performance-body" class="section-body">
                    <div class="metrics-cards">
                        <div class="metric-card">
                            <h3>Total Submissions</h3>
                            <div class="value" id="totalSubmissions">-</div>
                        </div>
                        <div class="metric-card success">
                            <h3>Strong</h3>
                            <div class="value" id="strongSubmissions">-</div>
                            <div class="percentage" id="strongPercentageText">-</div>
                        </div>
                        <div class="metric-card info">
                            <h3>Good Progress</h3>
                            <div class="value" id="goodProgressSubmissions">-</div>
                            <div class="percentage" id="goodProgressPercentageText">-</div>
                        </div>
                        <div class="metric-card warning">
                            <h3>Struggling</h3>
                            <div class="value" id="strugglingSubmissions">-</div>
                            <div class="percentage" id="strugglingPercentageText">-</div>
                        </div>
                        <div class="metric-card danger">
                            <h3>Poor</h3>
                            <div class="value" id="poorSubmissions">-</div>
                            <div class="percentage" id="poorPercentageText">-</div>
                        </div>
                    </div>
                    <div class="chart-container">
                        <div class="chart-wrapper">
                            <canvas id="performanceDistributionChart"></canvas>
                        </div>
                        <div class="chart-wrapper">
                            <canvas id="performanceGaugeChart"></canvas>
                        </div>
                    </div>
                    <div class="info-card">
                        <h3>Performance Overview</h3>
                        <p id="performanceSummary" class="description">Loading...</p>
                    </div>
                    <div class="comparative-analysis">
                        <h3>Comparative Analysis</h3>
                        <div class="tabs">
                            <div class="tab active" data-tab="success-factors">Success Factors</div>
                            <div class="tab" data-tab="challenge-areas">Challenge Areas</div>
                        </div>
                        <div id="success-factors" class="tab-content active">
                            <div class="info-card">
                                <h3>What Successful Students Do Differently</h3>
                                <div id="successFactorsList"></div>
                            </div>
                        </div>
                        <div id="challenge-areas" class="tab-content">
                            <div class="info-card">
                                <h3>Common Challenges for Struggling Students</h3>
                                <div id="challengeAreasList"></div>
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            <!-- Error Analysis Section -->
            <section id="errors" class="section">
                <div class="section-header">
                    <h2>Stage 2: Error Analysis</h2>
                    <button class="collapse-btn" data-target="errors-body">−</button>
                </div>
                <div id="errors-body" class="section-body">
                    <div class="chart-container">
                        <div class="chart-wrapper">
                            <canvas id="topErrorsChart"></canvas>
                        </div>
                    </div>
                    <div id="errorCards" class="card-wrapper">
                        <!-- Error cards will be dynamically inserted here -->
                    </div>
                </div>
            </section>

            <!-- Error Correlations Section -->
            <section id="correlations" class="section">
                <div class="section-header">
                    <h2>Stage 3: Error Correlations</h2>
                    <button class="collapse-btn" data-target="correlations-body">−</button>
                </div>
                <div id="correlations-body" class="section-body">
                    <p>The matrix below shows which errors tend to occur together in the same submissions:</p>
                    <div class="heatmap-container">
                        <div class="correlation-matrix-container">
                            <table class="correlation-matrix" id="correlationMatrix">
                                <!-- Correlation matrix will be dynamically inserted here -->
                            </table>
                        </div>
                    </div>
                    <div id="correlationCards" class="card-wrapper">
                        <!-- Correlation cards will be dynamically inserted here -->
                    </div>
                </div>
            </section>

            <!-- Misconceptions Section -->
            <section id="misconceptions" class="section">
                <div class="section-header">
                    <h2>Stage 4: Potential Misconceptions</h2>
                    <button class="collapse-btn" data-target="misconceptions-body">−</button>
                </div>
                <div id="misconceptions-body" class="section-body">
                    <div class="tabs">
                        <div class="tab active" data-tab="all-misconceptions">All Misconceptions</div>
                        <div class="tab" data-tab="by-category">By Error Category</div>
                        <div class="tab" data-tab="concept-map">Misconception-Error Diagram</div>
                    </div>
                    <div id="all-misconceptions" class="tab-content active">
                        <div id="misconceptionCards" class="card-wrapper">
                            <!-- Misconception cards will be dynamically inserted here -->
                        </div>
                    </div>
                    <div id="by-category" class="tab-content">
                        <p>Select an error category to see related misconceptions:</p>
                        <div id="categoryMisconceptions" class="card-wrapper">
                            <!-- Category-based misconception cards will be dynamically inserted here -->
                        </div>
                    </div>
                    <div id="concept-map" class="tab-content">
                      <div id="concept-map-container" class="concept-map-container">
                          <!-- The map will be dynamically inserted here -->
                      </div>
                    </div>
                </div>
            </section>

            <!-- Action Plan Section -->
            <section id="action-plan" class="section">
                <div class="section-header">
                    <h2>Stage 5: Instructional Action Plan</h2>
                    <button class="collapse-btn" data-target="action-plan-body">−</button>
                </div>
                <div id="action-plan-body" class="section-body">
                    <div class="action-plan-overview">
                        <div class="chart-wrapper">
                            <canvas id="priorityChart"></canvas>
                        </div>
                        <div class="chart-wrapper">
                          <canvas id="effortDistributionChart"></canvas>
                        </div>
                        <div class="action-summary">
                            <h3>Priority Overview</h3>
                            <p>Interventions are categorized by priority level based on the number of affected students and the severity of the misconception.</p>
                        </div>
                    </div>
                    <div class="tabs">
                        <div class="tab active" data-tab="priority-view">By Priority</div>
                        <div class="tab" data-tab="misconception-view">By Misconception</div>
                        <div class="tab" data-tab="timeline-view">Implementation Timeline</div>
                    </div>
                    <div id="priority-view" class="tab-content active">
                        <div id="actionPlanItems" class="card-wrapper">
                            <!-- Action plan cards will be dynamically inserted here -->
                        </div>
                    </div>
                    <div id="misconception-view" class="tab-content">
                        <div id="actionByMisconception" class="card-wrapper">
                            <!-- Misconception-based action plan will be dynamically inserted here -->
                        </div>
                    </div>
                    <div id="timeline-view" class="tab-content">
                        <div class="timeline-container">
                            <div class="timeline-header">
                                <div class="timeline-phase">Immediate (Week 1)</div>
                                <div class="timeline-phase">Short-term (Weeks 2-3)</div>
                                <div class="timeline-phase">Long-term (Weeks 4+)</div>
                            </div>
                            <div id="timelineContent" class="timeline-content">
                                <!-- Timeline content will be dynamically inserted here -->
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            <!-- Additional Insights Section -->
            <section id="additional" class="section">
                <div class="section-header">
                    <h2>Stage 6: Additional Insights</h2>
                    <button class="collapse-btn" data-target="additional-body">−</button>
                </div>
                <div id="additional-body" class="section-body">
                    <div class="tabs">
                        <div class="tab active" data-tab="good-practices">Good Practices</div>
                        <div class="tab" data-tab="submission-patterns">Submission Patterns</div>
                        <div class="tab" data-tab="error-categories">Error Categories</div>
                    </div>
                    <div id="good-practices" class="tab-content active">
                        <div id="goodPracticeCards" class="card-wrapper">
                            <!-- Good practice cards will be dynamically inserted here -->
                        </div>
                        <div class="info-card">
                            <h3>Correct vs. Incorrect Implementations</h3>
                            <div class="code-comparison-container" id="codeComparisonContainer">
                                <!-- Code comparison will be dynamically inserted here -->
                            </div>
                        </div>
                    </div>
                    <div id="submission-patterns" class="tab-content">
                        <div id="submissionPatternsContent" class="info-card">
                            <!-- Submission patterns will be dynamically inserted here -->
                        </div>
                        <div class="chart-wrapper">
                            <canvas id="submissionTimelineChart"></canvas>
                        </div>
                    </div>
                    <div id="error-categories" class="tab-content">
                        <div class="info-card">
                            <h3>Error Categories Distribution</h3>
                            <div class="error-categories-content" id="errorCategoriesContent">
                                <!-- Error categories content will be dynamically inserted here -->
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        </main>
    </div>

    <!-- Embed data directly in the main HTML file -->
    <script>
       const analysisData = {{.Text}};
  </script>

    <script>

document.addEventListener('DOMContentLoaded', () => {
  const sidebarItem = document.querySelector('[data-section="regenerate-btn"]');

  if (sidebarItem) {
    sidebarItem.addEventListener('click', function (e) {
      // 🛑 Stop navigation AND bubbling
      e.preventDefault();
      e.stopImmediatePropagation(); // this is stronger than stopPropagation

      const userConfirmed = window.confirm('Are you sure you want to regenerate the summary?');
      if (!userConfirmed) return;

      const span = this.querySelector('span');
      span.textContent = '⏳ Regenerating...';
      this.classList.add('disabled');

      setTimeout(() => {
        handleSummarySubmission(true);
      }, 50);
    });
  }
});

function handleSummarySubmission(isNew, onComplete) {
  const urlParams = new URLSearchParams(window.location.search);
  const problemId = urlParams.get('problem_id');

  const tempForm = document.createElement('form');
  tempForm.method = 'POST';
  tempForm.action = '/vi_view?problem_id=' + encodeURIComponent(problemId);

  const inputs = {
    problem_id: problemId,
    new: isNew
  };

  for (const [name, value] of Object.entries(inputs)) {
    const input = document.createElement('input');
    input.type = 'hidden';
    input.name = name;
    input.value = value;
    tempForm.appendChild(input);
  }

  document.body.appendChild(tempForm);
  tempForm.submit();

  if (typeof onComplete === 'function') {
    onComplete();
  }
}

// Main JavaScript for the dashboard
document.addEventListener('DOMContentLoaded', () => {
    // Toggle dark mode
    const darkModeToggle = document.getElementById('darkModeToggle');
    darkModeToggle.addEventListener('change', () => {
        document.body.classList.toggle('dark-mode');
        // Redraw charts when toggling dark mode
        initializeCharts();
    });
    const conceptMapTab = document.querySelector('[data-tab="concept-map"]');
    if (conceptMapTab) {
        conceptMapTab.addEventListener('click', () => {
            // Allow time for the tab to become active
            setTimeout(createConceptMap, 10);
        });
    }
    
    // Navigation
    const navItems = document.querySelectorAll('.nav-item');
    const sections = document.querySelectorAll('.section');

    navItems.forEach(item => {
        item.addEventListener('click', () => {
            // Remove active class from all nav items
            navItems.forEach(navItem => navItem.classList.remove('active'));
            
            // Add active class to clicked nav item
            item.classList.add('active');
            
            // Hide all sections
            sections.forEach(section => section.style.display = 'none');
            
            // Show the corresponding section
            const sectionId = item.getAttribute('data-section');
            document.getElementById(sectionId).style.display = 'block';
            
            // If the misconceptions section is displayed, ensure concept map is drawn
            if (sectionId === 'misconceptions') {
                // Use setTimeout to ensure the section is visible before drawing
                setTimeout(() => {
                    const conceptMapTab = document.querySelector('[data-tab="concept-map"]');
                    if (conceptMapTab && conceptMapTab.classList.contains('active')) {
                        createConceptMap();
                    }
                }, 10);
            }
        });
    });

    // Initially show only the summary section
    sections.forEach(section => {
        if (section.id !== 'summary') {
            section.style.display = 'none';
        }
    });

    // Collapse section functionality
    const collapseBtns = document.querySelectorAll('.collapse-btn');
    collapseBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            const targetId = btn.getAttribute('data-target');
            const targetElement = document.getElementById(targetId);
            
            if (targetElement.classList.contains('collapsed')) {
                targetElement.classList.remove('collapsed');
                btn.textContent = '−';
            } else {
                targetElement.classList.add('collapsed');
                btn.textContent = '+';
            }
        });
    });

    // Tab functionality with concept map handling
    const tabs = document.querySelectorAll('.tab');
    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            // Get parent tabs container
            const tabsContainer = tab.parentElement;
            
            // Remove active class from all tabs in this container
            tabsContainer.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            
            // Add active class to clicked tab
            tab.classList.add('active');
            
            // Get the associated tab content id
            const tabContentId = tab.getAttribute('data-tab');
            
            // Hide all tab contents in the same section
            const section = tabsContainer.closest('.section-body');
            section.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));
            
            // Show the corresponding tab content
            const tabContent = section.querySelector('#' + tabContentId);
            tabContent.classList.add('active');

            // If switching to concept map tab, redraw the concept map
            if (tabContentId === 'concept-map') {
                // Use setTimeout to ensure the tab is visible before drawing
                setTimeout(createConceptMap, 10);
            }
        });
    });

    // Add resize event listener to redraw concept map on window resize
    let resizeTimeout;
    window.addEventListener('resize', () => {
        // Use a debounce pattern to avoid excessive redraws during resizing
        clearTimeout(resizeTimeout);
        resizeTimeout = setTimeout(() => {
            // Only redraw if the concept map tab is active
            const conceptMapTab = document.querySelector('[data-tab="concept-map"]');
            const misconceptionsSection = document.getElementById('misconceptions');

            if (misconceptionsSection &&
                misconceptionsSection.style.display !== 'none' &&
                conceptMapTab &&
                conceptMapTab.classList.contains('active')) {
                createConceptMap();
            }
        }, 250); // Wait 250ms after resize ends before redrawing
    });

    // Load data into the dashboard
    loadDashboardData();
});

// Function to get a value with fallback
function getValueWithFallback(obj, property, fallback) {
    return obj && obj[property] !== undefined ? obj[property] : fallback;
}

// Function to load data into the dashboard
function loadDashboardData() {
    // Summary Dashboard
    document.getElementById('highlight-performance').textContent = analysisData.stage_1_performance_analysis.overall_summary;

    // Calculate passing percentage (Good Progress + Strong)
    const goodProgressCount = analysisData.stage_1_performance_analysis.performance_distribution.Good_Progress.count;
    const strongCount = analysisData.stage_1_performance_analysis.performance_distribution.Strong.count;
    const totalCount = analysisData.stage_1_performance_analysis.total_submissions;
    const passingPercentage = ((goodProgressCount + strongCount) / totalCount * 100).toFixed(2);

    document.getElementById('passing-progress').style.width = passingPercentage + '%';
    document.getElementById('passing-percentage').textContent = passingPercentage + '%';


    // Top Error
    const topError = analysisData.stage_2_error_analysis.top_errors[0];
    document.getElementById('highlight-top-error').textContent =
  topError.category + ' (' + topError.occurrence_percentage + '): ' + topError.description;


    // Primary Misconception
    const primaryMisconception = analysisData.stage_4_misconception_analysis.potential_misconceptions[0];
    document.getElementById('highlight-misconception').textContent =
  primaryMisconception.misconception + ' (' + primaryMisconception.occurrence_percentage + '): ' + primaryMisconception.explanation;

    // Priority Intervention
    const topIntervention = analysisData.stage_5_instructional_plan.interventions.find(i => i.priority_level === "High");
    document.getElementById('highlight-intervention').textContent = topIntervention.target_misconception + ': ' + topIntervention.instructional_strategies[0];

    // Student Success Gap
    const strongPercentage = parseFloat(analysisData.stage_1_performance_analysis.performance_distribution.Strong.percentage);
    const poorPercentage = parseFloat(analysisData.stage_1_performance_analysis.performance_distribution.Poor.percentage);
    const gap = poorPercentage - strongPercentage;
    document.getElementById('highlight-gap').textContent = gap.toFixed(2) + "% more students are in the 'Poor' category than in the 'Strong' category. Focus on moving students from 'Struggling' to 'Good Progress' for biggest impact.";

    // Performance Analysis
    document.getElementById('totalSubmissions').textContent = analysisData.stage_1_performance_analysis.total_submissions;
    document.getElementById('strongSubmissions').textContent = analysisData.stage_1_performance_analysis.performance_distribution.Strong.count;
    document.getElementById('strongPercentageText').textContent = analysisData.stage_1_performance_analysis.performance_distribution.Strong.percentage;
    document.getElementById('goodProgressSubmissions').textContent = analysisData.stage_1_performance_analysis.performance_distribution.Good_Progress.count;
    document.getElementById('goodProgressPercentageText').textContent = analysisData.stage_1_performance_analysis.performance_distribution.Good_Progress.percentage;
    document.getElementById('strugglingSubmissions').textContent = analysisData.stage_1_performance_analysis.performance_distribution.Struggling.count;
    document.getElementById('strugglingPercentageText').textContent = analysisData.stage_1_performance_analysis.performance_distribution.Struggling.percentage;
    document.getElementById('poorSubmissions').textContent = analysisData.stage_1_performance_analysis.performance_distribution.Poor.count;
    document.getElementById('poorPercentageText').textContent = analysisData.stage_1_performance_analysis.performance_distribution.Poor.percentage;

    document.getElementById('performanceSummary').textContent = analysisData.stage_1_performance_analysis.overall_summary;

    // Success Factors and Challenge Areas
    const successFactorsContainer = document.getElementById('successFactorsList');
    const challengeAreasContainer = document.getElementById('challengeAreasList');

    if (analysisData.stage_6_additional_insights.comparative_analysis) {
        // Success Factors
        const successFactors = analysisData.stage_6_additional_insights.comparative_analysis.success_factors;
        let successHtml = '<ul class="factor-list">';

        successFactors.forEach(factor => {
            successHtml +=
  '<li>' +
    '<div class="factor-icon success">✓</div>' +
    '<div>' + factor + '</div>' +
  '</li>';

        });

        successHtml += '</ul>';
        successFactorsContainer.innerHTML = successHtml;

        // Challenge Areas
        const challengeFactors = analysisData.stage_6_additional_insights.comparative_analysis.challenge_factors;
        let challengeHtml = '<ul class="factor-list">';

        challengeFactors.forEach(factor => {
            challengeHtml +=
  '<li>' +
    '<div class="factor-icon danger">✗</div>' +
    '<div>' + factor + '</div>' +
  '</li>';

        });

        challengeHtml += '</ul>';
        challengeAreasContainer.innerHTML = challengeHtml;
    }

    // Error Analysis
    // Generate Error Cards
    const errorCardsContainer = document.getElementById('errorCards');
    errorCardsContainer.innerHTML = '';

    analysisData.stage_2_error_analysis.top_errors.forEach(error => {
        const card = document.createElement('div');
        card.className = 'info-card';

        card.innerHTML =
  '<h3>' + error.category + '</h3>' +
  '<div class="tags">' +
    '<div class="tag danger">' + error.occurrence_percentage + '</div>' +
    '<div class="tag primary">' + error.occurrence_count + ' Students</div>' +
  '</div>' +
  '<p class="description">' + error.description + '</p>' +
  '<div class="code-preview">' +
    '<pre><code class="language-python">' + error.example_code.join('\n') + '</code></pre>' +
  '</div>';


        errorCardsContainer.appendChild(card);
    });

    // Correlation Analysis
    // Create Correlation Matrix
    const correlationMatrix = document.getElementById('correlationMatrix');
    correlationMatrix.innerHTML = '';

    // Get unique error categories from correlations
    const correlatedErrorPairs = analysisData.stage_3_correlation_analysis.error_correlations.map(corr => corr.correlated_errors);
    const uniqueCorrelatedErrors = [...new Set(correlatedErrorPairs.flat())];

    // Create header row
    const headerRow = document.createElement('tr');
    headerRow.innerHTML = '<th></th>';
    uniqueCorrelatedErrors.forEach(error => {
        headerRow.innerHTML += '<th>' + error + '</th>';

    });
    correlationMatrix.appendChild(headerRow);

    // Create matrix rows
    uniqueCorrelatedErrors.forEach(rowError => {
        const row = document.createElement('tr');
        row.innerHTML = '<th>' + rowError + '</th>';


        uniqueCorrelatedErrors.forEach(colError => {
            let cell = document.createElement('td');

            if (rowError === colError) {
                // Diagonal cells
                cell.innerHTML = '<span class="correlation-value" style="background-color: #4e73df;">-</span>';
            } else {
                // Find if there's a correlation between these errors
                const correlation = analysisData.stage_3_correlation_analysis.error_correlations.find(corr =>
                    corr.correlated_errors.includes(rowError) && corr.correlated_errors.includes(colError)
                );

                if (correlation) {
                    cell.innerHTML = '<span class="correlation-value" style="background-color: #e74a3b;">' + correlation.correlation_percentage + '</span>';
                } else {
                    cell.innerHTML = '<span class="correlation-value" style="background-color: #cccccc;">0%</span>';
                }
            }

            row.appendChild(cell);
        });

        correlationMatrix.appendChild(row);
    });

    // Generate Correlation Cards
    const correlationCardsContainer = document.getElementById('correlationCards');
    correlationCardsContainer.innerHTML = '';

    analysisData.stage_3_correlation_analysis.error_correlations.forEach(correlation => {
        const card = document.createElement('div');
        card.className = 'info-card';

        card.innerHTML =
  '<h3>' + correlation.correlated_errors.join(' + ') + '</h3>' +
  '<div class="tags">' +
    '<div class="tag warning">' + correlation.correlation_percentage + '</div>' +
    '<div class="tag primary">' + correlation.correlation_count + ' Students</div>' +
  '</div>' +
  '<p class="description">' + correlation.hypothesis + '</p>' +
  '<div class="code-preview">' +
    '<pre><code class="language-python">' + correlation.example_code.join('\n') + '</code></pre>' +
  '</div>';


        correlationCardsContainer.appendChild(card);
    });

    // Misconceptions Analysis
    // Generate Misconception Cards
    const misconceptionCardsContainer = document.getElementById('misconceptionCards');
    misconceptionCardsContainer.innerHTML = '';

    analysisData.stage_4_misconception_analysis.potential_misconceptions.forEach(misconception => {
        const card = document.createElement('div');
        card.className = 'info-card';

        card.innerHTML =
  '<h3>' + misconception.misconception + '</h3>' +
  '<div class="tags">' +
    '<div class="tag danger">' + misconception.occurrence_percentage + '</div>' +
    '<div class="tag primary">' + misconception.occurrence_count + ' Students</div>' +
  '</div>' +
  '<div class="tags">' +
    misconception.related_error_categories.map(function(category) {
      return '<div class="tag info">' + category + '</div>';
    }).join('') +
  '</div>' +
  '<p class="description">' + misconception.explanation + '</p>' +
  '<div class="code-preview">' +
    '<pre><code class="language-python">' + misconception.example_code.join('\n') + '</code></pre>' +
  '</div>';


        misconceptionCardsContainer.appendChild(card);
    });

    // Generate category-based misconception view
    const categoryMisconceptionsContainer = document.getElementById('categoryMisconceptions');
    categoryMisconceptionsContainer.innerHTML = '';

    // Get all unique error categories from misconceptions
    const allErrorCategories = [...new Set(analysisData.stage_4_misconception_analysis.potential_misconceptions.flatMap(m => m.related_error_categories))];

    allErrorCategories.forEach(category => {
        const card = document.createElement('div');
        card.className = 'info-card';

        // Find misconceptions related to this category
        const relatedMisconceptions = analysisData.stage_4_misconception_analysis.potential_misconceptions.filter(m =>
            m.related_error_categories.includes(category)
        );

        let misconceptionsList = '';
        relatedMisconceptions.forEach(m => {
            misconceptionsList +=
  '<div style="margin-bottom: 1rem; padding-bottom: 1rem; border-bottom: 1px solid #e3e6f0;">' +
    '<h4 style="margin-bottom: 0.5rem;">' + m.misconception + ' (' + m.occurrence_percentage + ')</h4>' +
    '<p>' + m.explanation + '</p>' +
  '</div>';

        });

        card.innerHTML =
  '<h3>' + category + '</h3>' +
  '<p class="description">Related misconceptions:</p>' +
  misconceptionsList;


        categoryMisconceptionsContainer.appendChild(card);
    });

    // Action Plan
    // Generate Action Plan Cards
    const actionPlanContainer = document.getElementById('actionPlanItems');
    actionPlanContainer.innerHTML = '';

    // Sort interventions by priority
    const sortedInterventions = [...analysisData.stage_5_instructional_plan.interventions].sort((a, b) => {
        const priorityOrder = { "High": 0, "Medium": 1, "Low": 2 };
        return priorityOrder[a.priority_level] - priorityOrder[b.priority_level];
    });

    sortedInterventions.forEach(intervention => {
        const card = document.createElement('div');
        card.className = 'info-card';

        let priorityClass = '';
        if (intervention.priority_level === 'High') {
            priorityClass = 'priority-high';
        } else if (intervention.priority_level === 'Medium') {
            priorityClass = 'priority-medium';
        } else {
            priorityClass = 'priority-low';
        }

        // Get estimated effort with fallback to 'Medium'
        const effort = getValueWithFallback(intervention, 'estimated_effort', 'Medium');
        let effortClass = '';
        if (effort === 'High') {
            effortClass = 'effort-high';
        } else if (effort === 'Medium') {
            effortClass = 'effort-medium';
        } else {
            effortClass = 'effort-low';
        }

        let strategiesList = '<ul class="strategyList">';
        intervention.instructional_strategies.forEach(strategy => {
            strategiesList += '<li>' + strategy + '</li>';

        });
        strategiesList += '</ul>';

        card.innerHTML =
  '<h3>' + intervention.target_misconception + '</h3>' +
  '<div class="indicators-row">' +
    '<div class="priority-indicator ' + priorityClass + '">' + intervention.priority_level + ' Priority</div>' +
    '<div class="effort-indicator ' + effortClass + '">' + effort + ' Effort</div>' +
  '</div>' +
  '<div class="tags">' +
    intervention.related_errors.map(function(error) {
      return '<div class="tag info">' + error + '</div>';
    }).join('') +
  '</div>' +
  '<p class="description">Instructional strategies:</p>' +
  strategiesList;


        actionPlanContainer.appendChild(card);
    });

    // Action Plan by Misconception
    const actionByMisconceptionContainer = document.getElementById('actionByMisconception');
    if (actionByMisconceptionContainer) {
        actionByMisconceptionContainer.innerHTML = '';

        // Group interventions by target misconception
        const misconceptionMap = new Map();
        analysisData.stage_4_misconception_analysis.potential_misconceptions.forEach(misconception => {
            const relatedInterventions = analysisData.stage_5_instructional_plan.interventions.filter(i =>
                i.target_misconception === misconception.misconception
            );

            misconceptionMap.set(misconception, relatedInterventions);
        });

        // Create cards for each misconception with its interventions
        misconceptionMap.forEach((interventions, misconception) => {
            if (interventions.length > 0) {
                const card = document.createElement('div');
                card.className = 'info-card';

                let interventionsHtml = '';
                interventions.forEach(intervention => {
                    let priorityClass = intervention.priority_level === 'High' ? 'priority-high' :
                                       intervention.priority_level === 'Medium' ? 'priority-medium' : 'priority-low';

                    // Get estimated effort with fallback
                    const effort = getValueWithFallback(intervention, 'estimated_effort', 'Medium');

                    interventionsHtml +=
  '<div style="margin-bottom: 1rem; padding-bottom: 1rem; border-bottom: 1px solid #e3e6f0;">' +
    '<div class="indicators-row">' +
      '<div class="priority-indicator ' + priorityClass + '">' + intervention.priority_level + ' Priority</div>' +
      '<div class="effort-indicator">' + effort + ' Effort</div>' +
    '</div>' +
    '<ul class="strategyList">' +
      intervention.instructional_strategies.map(function(strategy) {
        return '<li>' + strategy + '</li>';
      }).join('') +
    '</ul>' +
  '</div>';

                });

                card.innerHTML =
  '<h3>' + misconception.misconception + '</h3>' +
  '<div class="tags">' +
    '<div class="tag danger">' + misconception.occurrence_percentage + '</div>' +
  '</div>' +
  '<p class="description">' + misconception.explanation + '</p>' +
  '<h4>Intervention Strategies:</h4>' +
  interventionsHtml;


                actionByMisconceptionContainer.appendChild(card);
            }
        });
    }

    // Implementation Timeline
    const timelineContentContainer = document.getElementById('timelineContent');
    if (timelineContentContainer) {
        timelineContentContainer.innerHTML = '';

        // Group interventions by implementation phase
        const phases = {
            "Immediate": [],
            "Short-term": [],
            "Long-term": []
        };

        analysisData.stage_5_instructional_plan.interventions.forEach(intervention => {
            if (intervention.implementation_phase) {
                phases[intervention.implementation_phase].push(intervention);
            } else if (intervention.priority_level === "High") {
                phases["Immediate"].push(intervention);
            } else if (intervention.priority_level === "Medium") {
                phases["Short-term"].push(intervention);
            } else {
                phases["Long-term"].push(intervention);
            }
        });

        // Create timeline columns
        const immediateColumn = document.createElement('div');
        const shortTermColumn = document.createElement('div');
        const longTermColumn = document.createElement('div');

        // Fill immediate column
        phases["Immediate"].forEach(intervention => {
            const timelineItem = document.createElement('div');
            timelineItem.className = 'timeline-item priority-' + intervention.priority_level.toLowerCase();

            // Get estimated effort with fallback
            const effort = getValueWithFallback(intervention, 'estimated_effort', 'Medium');

            timelineItem.innerHTML =
  '<h4>' + intervention.target_misconception + '</h4>' +
  '<div class="timeline-effort">' + effort + ' Effort</div>' +
  '<p>' + intervention.instructional_strategies[0] + '</p>';


            immediateColumn.appendChild(timelineItem);
        });

        // Fill short-term column
        phases["Short-term"].forEach(intervention => {
            const timelineItem = document.createElement('div');
            timelineItem.className = 'timeline-item priority-' + intervention.priority_level.toLowerCase();

            // Get estimated effort with fallback
            const effort = getValueWithFallback(intervention, 'estimated_effort', 'Medium');

            timelineItem.innerHTML =
  '<h4>' + intervention.target_misconception + '</h4>' +
  '<div class="timeline-effort">' + effort + ' Effort</div>' +
  '<p>' + intervention.instructional_strategies[0] + '</p>';


            shortTermColumn.appendChild(timelineItem);
        });

        // Fill long-term column
        phases["Long-term"].forEach(intervention => {
            const timelineItem = document.createElement('div');
            timelineItem.className = 'timeline-item priority-' + intervention.priority_level.toLowerCase();

            // Get estimated effort with fallback
            const effort = getValueWithFallback(intervention, 'estimated_effort', 'Medium');

            timelineItem.innerHTML =
  '<h4>' + intervention.target_misconception + '</h4>' +
  '<div class="timeline-effort">' + effort + ' Effort</div>' +
  '<p>' + intervention.instructional_strategies[0] + '</p>';

            longTermColumn.appendChild(timelineItem);
        });

        timelineContentContainer.appendChild(immediateColumn);
        timelineContentContainer.appendChild(shortTermColumn);
        timelineContentContainer.appendChild(longTermColumn);
    }

    // Additional Insights
    // Generate Good Practice Cards
    const goodPracticeContainer = document.getElementById('goodPracticeCards');
    goodPracticeContainer.innerHTML = '';

    analysisData.stage_6_additional_insights.common_good_practices.forEach(practice => {
        const card = document.createElement('div');
        card.className = 'info-card';

        // Use correct_implementation flag if available
        const isCorrectImplementation = getValueWithFallback(practice, 'correct_implementation', false);
        const implementationTag = isCorrectImplementation ?
            '<div class="implementation-tag correct">✓ Correct Implementation</div>' :
            '<div class="implementation-tag">Implementation Example</div>';

        card.innerHTML =
  '<h3>Good Practice</h3>' +
  implementationTag +
  '<p class="description">' + practice.description + '</p>' +
  '<div class="code-preview">' +
    '<pre><code class="language-python">' + practice.example_code.join('\n') + '</code></pre>' +
  '</div>';


        goodPracticeContainer.appendChild(card);
    });

    // Submission Patterns
    const submissionPatternsContainer = document.getElementById('submissionPatternsContent');

    submissionPatternsContainer.innerHTML = 
  '<h3>Submission Patterns</h3>' +
  '<p class="description"><strong>Observation:</strong> ' + analysisData.stage_6_additional_insights.submission_patterns.observation + '</p>' +
  '<p class="description"><strong>Potential Impact:</strong> ' + analysisData.stage_6_additional_insights.submission_patterns.potential_impact + '</p>';


    // Code Comparison
    const codeComparisonContainer = document.getElementById('codeComparisonContainer');
    if (codeComparisonContainer) {
        // Find an error example and a good practice to compare
        const errorExample = analysisData.stage_2_error_analysis.top_errors[0];

        // Find a good practice that's marked as correct implementation if possible
        let goodPractice = analysisData.stage_6_additional_insights.common_good_practices.find(p =>
            p.correct_implementation === true
        );

        // Fallback to first good practice if none are marked as correct
        if (!goodPractice) {
            goodPractice = analysisData.stage_6_additional_insights.common_good_practices[0];
        }

        codeComparisonContainer.innerHTML =
  '<div class="code-comparison-container">' +
    '<div class="code-comparison">' +
      '<div class="code-comparison-header">' +
        '<span>❌ Incorrect Implementation</span>' +
        '<span class="tag danger">' + errorExample.category + '</span>' +
      '</div>' +
      '<div class="code-comparison-body">' +
        '<pre><code class="language-python">' + errorExample.example_code.join('\n') + '</code></pre>' +
      '</div>' +
    '</div>' +
    '<div class="code-comparison">' +
      '<div class="code-comparison-header">' +
        '<span>✓ Correct Implementation</span>' +
        '<span class="tag success">Good Practice</span>' +
      '</div>' +
      '<div class="code-comparison-body">' +
        '<pre><code class="language-python">' + goodPractice.example_code.join('\n') + '</code></pre>' +
      '</div>' +
    '</div>' +
  '</div>';

    }

    // Error Categories
    const errorCategoriesContainer = document.getElementById('errorCategoriesContent');
    if (errorCategoriesContainer) {
        errorCategoriesContainer.innerHTML = 
  '<div>' +
    '<h4>Error Categories Used in Analysis</h4>' +
    '<div class="tags" style="margin-top: 0.5rem; margin-bottom: 1.5rem;">' +
      analysisData.stage_6_additional_insights.error_category_distribution.all_categories_used.map(function(category) {
        return '<div class="tag primary">' + category + '</div>';
      }).join('') +
    '</div>' +
  '</div>' +
  '<div>' +
    '<h4>Additional Categories Identified</h4>' +
    '<div class="tags" style="margin-top: 0.5rem;">' +
      analysisData.stage_6_additional_insights.error_category_distribution.additional_categories_identified.map(function(category) {
        return '<div class="tag info">' + category + '</div>';
      }).join('') +
    '</div>' +
  '</div>';

    }

    // Apply syntax highlighting to code elements
    if (typeof Prism !== 'undefined') {
        Prism.highlightAll();
    }

    // Initialize charts
    initializeCharts();

    if (typeof setupCodePreviewListeners === 'function') {
        setupCodePreviewListeners();
    }
}
</script>
    <script>
// Chart initialization for the dashboard
function initializeCharts() {
    // Performance Distribution Chart
    createPerformanceDistributionChart();

    // Performance Gauge Chart
    createPerformanceGaugeChart();

    // Top Errors Chart
    createTopErrorsChart();

    // Priority Chart
    createPriorityChart();

    // Effort Distribution Chart (NEW)
    createEffortDistributionChart();

    // Concept Map
    const conceptMapTab = document.getElementById('concept-map');
    if (conceptMapTab && conceptMapTab.classList.contains('active')) {
        createConceptMap();
    }

    // Submission Timeline Chart
    createSubmissionTimelineChart();
}

function createPerformanceDistributionChart() {
    const ctx = document.getElementById('performanceDistributionChart');
    if (!ctx) return;

    new Chart(ctx, {
        type: 'pie',
        data: {
            labels: ['Strong', 'Good Progress', 'Struggling', 'Poor'],
            datasets: [{
                data: [
                    analysisData.stage_1_performance_analysis.performance_distribution.Strong.count,
                    analysisData.stage_1_performance_analysis.performance_distribution.Good_Progress.count,
                    analysisData.stage_1_performance_analysis.performance_distribution.Struggling.count,
                    analysisData.stage_1_performance_analysis.performance_distribution.Poor.count
                ],
                backgroundColor: [
                    '#1cc88a',
                    '#36b9cc',
                    '#f6c23e',
                    '#e74a3b'
                ]
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'bottom'
                },
                title: {
                    display: true,
                    text: 'Performance Distribution'
                },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            const label = context.label || '';
                            const value = context.raw || 0;
                            const total = context.dataset.data.reduce((a, b) => a + b, 0);
                            const percentage = Math.round((value / total) * 100);
                            return label + ': ' + value + ' (' + percentage + '%)';
                        }
                    }
                }
            }
        }
    });
}

function createPerformanceGaugeChart() {
    const ctx = document.getElementById('performanceGaugeChart');
    if (!ctx) return;

    // Calculate passing percentage (Good Progress + Strong)
    const goodProgressCount = analysisData.stage_1_performance_analysis.performance_distribution.Good_Progress.count;
    const strongCount = analysisData.stage_1_performance_analysis.performance_distribution.Strong.count;
    const totalCount = analysisData.stage_1_performance_analysis.total_submissions;
    const passingPercentage = ((goodProgressCount + strongCount) / totalCount * 100).toFixed(2);

    new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: ['Satisfactory', 'Needs Improvement'],
            datasets: [{
                data: [
                    goodProgressCount + strongCount,
                    totalCount - (goodProgressCount + strongCount)
                ],
                backgroundColor: [
                    '#1cc88a',
                    '#e74a3b'
                ],
                borderWidth: 0
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            circumference: 180,
            rotation: -90,
            cutout: '75%',
            plugins: {
                legend: {
                    display: false
                },
                title: {
                    display: true,
                    text: passingPercentage + '% Satisfactory',
                    position: 'bottom',
                    padding: {
                        top: 10
                    }
                },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            const label = context.label || '';
                            const value = context.raw || 0;
                            const total = context.dataset.data.reduce((a, b) => a + b, 0);
                            const percentage = Math.round((value / total) * 100);
                            return label + ': ' + value + ' (' + percentage + '%)';
                        }
                    }
                }
            }
        }
    });
}

function createEffortDistributionChart() {
    const ctx = document.getElementById('effortDistributionChart');
    if (!ctx) return;

    // Count interventions by effort level with fallback handling
    const efforts = {'Low': 0, 'Medium': 0, 'High': 0};

    analysisData.stage_5_instructional_plan.interventions.forEach(intervention => {
        const effort = intervention.estimated_effort || 'Medium'; // Default to Medium if missing
        efforts[effort]++;
    });

    new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: ['Low Effort', 'Medium Effort', 'High Effort'],
            datasets: [{
                data: [efforts.Low, efforts.Medium, efforts.High],
                backgroundColor: [
                    '#1cc88a', // Green for Low
                    '#f6c23e', // Yellow for Medium
                    '#e74a3b'  // Red for High
                ]
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'bottom'
                },
                title: {
                    display: true,
                    text: 'Interventions by Required Effort'
                },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            const label = context.label || '';
                            const value = context.raw || 0;
                            const total = context.dataset.data.reduce((a, b) => a + b, 0);
                            const percentage = Math.round((value / total) * 100);
                            return label + ': ' + value + ' (' + percentage + '%)';
                        }
                    }
                }
            }
        }
    });
}

function createTopErrorsChart() {
    const ctx = document.getElementById('topErrorsChart');
    if (!ctx) return;
    const errorCategories = analysisData.stage_2_error_analysis.top_errors.map(error => error.category);
    const errorCounts = analysisData.stage_2_error_analysis.top_errors.map(error => error.occurrence_count);
    const errorPercentages = analysisData.stage_2_error_analysis.top_errors.map(error =>
      parseFloat(error.occurrence_percentage.replace('%', ''))
    );
    new Chart(ctx, {
      type: 'bar',
      data: {
        labels: errorCategories,
        datasets: [
          {
            label: 'Occurrence Count',
            data: errorCounts,
            backgroundColor: '#4e73df',
            borderColor: '#4e73df',
            borderWidth: 1,
            yAxisID: 'y',
            order: 2  // Higher order means render first (underneath)
          },
          {
            label: 'Percentage',
            data: errorPercentages,
            backgroundColor: 'transparent',  // Make background transparent
            borderColor: '#1cc88a',
            borderWidth: 3,  // Slightly thicker line for visibility
            type: 'line',
            yAxisID: 'y1',
            order: 1,  // Lower order means render last (on top)
            pointBackgroundColor: '#1cc88a',
            pointRadius: 4,
            tension: 0.1  // Slight curve to the line
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            position: 'top'
          },
          title: {
            display: true,
            text: 'Top Error Categories'
          }
        },
        scales: {
          y: {
            beginAtZero: true,
            title: {
              display: true,
              text: 'Occurrence Count'
            },
            position: 'left'
          },
          y1: {
            beginAtZero: true,
            title: {
              display: true,
              text: 'Percentage (%)'
            },
            position: 'right',
            grid: {
              drawOnChartArea: false
            },
            ticks: {
              callback: function(value) {
                return value + '%';
              }
            },
            max: 100
          },
          x: {
            title: {
              display: true,
              text: 'Error Category'
            }
          }
        }
      }
    });
  }

function createPriorityChart() {
    const ctx = document.getElementById('priorityChart');
    if (!ctx) return;

    // Count interventions by priority level
    const priorities = {'High': 0, 'Medium': 0, 'Low': 0};

    analysisData.stage_5_instructional_plan.interventions.forEach(intervention => {
        priorities[intervention.priority_level]++;
    });

    new Chart(ctx, {
        type: 'pie',
        data: {
            labels: ['High Priority', 'Medium Priority', 'Low Priority'],
            datasets: [{
                data: [priorities.High, priorities.Medium, priorities.Low],
                backgroundColor: [
                    '#e74a3b',
                    '#f6c23e',
                    '#1cc88a'
                ]
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'bottom'
                },
                title: {
                    display: true,
                    text: 'Interventions by Priority'
                }
            }
        }
    });
}

function createSubmissionTimelineChart() {
    const ctx = document.getElementById('submissionTimelineChart');
    // Add a null check for both the chart element and the timeline data
    if (!ctx) return;
    if (!analysisData.stage_6_additional_insights.submission_patterns ||
        !analysisData.stage_6_additional_insights.submission_patterns.submission_timeline) {
        return;
    }

    const timeline = analysisData.stage_6_additional_insights.submission_patterns.submission_timeline;
    const dates = Object.keys(timeline);
    const submissionCounts = dates.map(date => timeline[date].count);
    // Handle both string percentages with % and numeric values
    const passingRates = dates.map(date => {
        const rate = timeline[date].passing_rate;
        return typeof rate === 'string' ? parseFloat(rate.replace('%', '')) : rate;
    });

    new Chart(ctx, {
        type: 'bar',
        data: {
            labels: dates,
            datasets: [
                {
                    label: 'Submission Count',
                    data: submissionCounts,
                    backgroundColor: '#4e73df',
                    borderColor: '#4e73df',
                    borderWidth: 1,
                    order: 2,  // Lower order means render last (on top)
                    yAxisID: 'y'
                },
                {
                    label: 'Passing Rate (%)',
                    data: passingRates,
                    backgroundColor: '#1cc88a',
                    borderColor: '#1cc88a',
                    borderWidth: 3,
                    order: 1,  // Lower order means render last (on top)
                    type: 'line',
                    yAxisID: 'y1'
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'top'
                },
                title: {
                    display: true,
                    text: 'Submission Timeline and Passing Rate'
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: 'Submission Count'
                    },
                    position: 'left'
                },
                y1: {
                    beginAtZero: true,
                    title: {
                        display: true,
                        text: 'Passing Rate (%)'
                    },
                    position: 'right',
                    grid: {
                        drawOnChartArea: false
                    },
                    ticks: {
                        callback: function(value) {
                            return value + '%';
                        }
                    },
                    max: 100
                },
                x: {
                    title: {
                        display: true,
                        text: 'Submission Date'
                    }
                }
            }
        }
    });
}


function createConceptMap() {
    // Get the container
    const container = document.getElementById('concept-map-container');
    if (!container) return;

    // Get values from controls or use defaults
    const radius1Control = document.getElementById('radius1-control');
    const radius2Control = document.getElementById('radius2-control');
    const centerXControl = document.getElementById('centerX-control');
    const centerYControl = document.getElementById('centerY-control');

    const radius1 = radius1Control ? parseInt(radius1Control.value) : 120;
    const radius2 = radius2Control ? parseInt(radius2Control.value) : 300;
    const centerX = centerXControl ? parseInt(centerXControl.value) : 600;
    const centerY = centerYControl ? parseInt(centerYControl.value) : 350;

    // Update value displays
    if (document.getElementById('radius1-value')) {
        document.getElementById('radius1-value').textContent = radius1;
    }
    if (document.getElementById('radius2-value')) {
        document.getElementById('radius2-value').textContent = radius2;
    }
    if (document.getElementById('centerX-value')) {
        document.getElementById('centerX-value').textContent = centerX;
    }
    if (document.getElementById('centerY-value')) {
        document.getElementById('centerY-value').textContent = centerY;
    }

    // Clear previous content
    container.innerHTML = '';

    // Create control panel if it doesn't exist
    if (!document.getElementById('concept-map-controls')) {
        createControlPanel(container, radius1, radius2, centerX, centerY);
    }

    // Create the concept map wrapper
    const conceptMap = document.createElement('div');
    conceptMap.className = 'simple-concept-map';

    // Process data - get misconceptions and error categories
    const misconceptions = analysisData.stage_4_misconception_analysis.potential_misconceptions;

    // Extract all unique error categories
    const errorCategories = new Set();
    misconceptions.forEach(m => {
        m.related_error_categories.forEach(e => errorCategories.add(e));
    });
    const errorCategoriesArray = Array.from(errorCategories);

    // Calculate positions in a circle for misconceptions
    const misconceptionNodes = [];

    misconceptions.forEach((misconception, i) => {
        // Calculate position in a circle
        const angle = (i / misconceptions.length) * Math.PI * 2;
        const x = centerX + Math.cos(angle) * radius1;
        const y = centerY + Math.sin(angle) * radius1;

        // Create misconception node
        const node = document.createElement('div');
        node.className = 'concept-node misconception-node';
        node.style.left = x + 'px';
		node.style.top = y + 'px';


        // Parse percentage to use for color intensity (pink with varying intensity)
        const percentage = parseFloat(misconception.occurrence_percentage);
        const brightness = 100 - (percentage * 0.5);
        node.style.backgroundColor = 'hsl(185, 100%, ' + brightness + '%)';


        // Add content
        node.innerHTML =
  '<div class="node-content">' +
    '<div class="node-title">' + misconception.misconception + '</div>' +
    '<div class="node-percentage">' + misconception.occurrence_percentage + '</div>' +
  '</div>';


        // Add tooltip
        node.setAttribute('data-tooltip', 
  '<strong>' + misconception.misconception + '</strong><br>' +
  'Affects ' + misconception.occurrence_percentage + ' of students<br>' +
  'Related to: ' + misconception.related_error_categories.join(', ')
);


        // Store node info for connections
        misconceptionNodes.push({
            element: node,
            x: x,
            y: y,
            misconception: misconception
        });

        // Add to concept map
        conceptMap.appendChild(node);

        // Add connection to center
        addConnection(conceptMap, centerX, centerY, x, y, 'hub-connection');
    });

    // Calculate positions in a larger circle for error categories
    const errorNodes = [];

    errorCategoriesArray.forEach((category, i) => {
        // Calculate position in a circle
        const angle = (i / errorCategoriesArray.length) * Math.PI * 2;
        const x = centerX + Math.cos(angle) * radius2;
        const y = centerY + Math.sin(angle) * radius2;

        // Create error category node
        const node = document.createElement('div');
        node.className = 'concept-node error-node';
        node.style.left = x + 'px';
		node.style.top = y + 'px';

        // Add content
        node.innerHTML = 
  '<div class="node-content">' +
    '<div class="node-title">' + category + '</div>' +
  '</div>';


        // Find related misconceptions for tooltip
        const relatedMisconceptions = misconceptions.filter(m =>
            m.related_error_categories.includes(category)
        );

        // Add tooltip
        node.setAttribute('data-tooltip', 
  '<strong>' + category + '</strong><br>' +
  'Related misconceptions:<br>' +
  relatedMisconceptions.map(function(m) {
    return m.misconception + ' (' + m.occurrence_percentage + ')';
  }).join('<br>')
);


        // Store node info for connections
        errorNodes.push({
            element: node,
            x: x,
            y: y,
            category: category
        });

        // Add to concept map
        conceptMap.appendChild(node);
    });

    // Add connections between misconceptions and related error categories
    misconceptionNodes.forEach(misconceptionNode => {
        const misconception = misconceptionNode.misconception;

        // For each related error category, find the node and connect
        misconception.related_error_categories.forEach(category => {
            const errorNode = errorNodes.find(n => n.category === category);
            if (errorNode) {
                addConnection(
                    conceptMap,
                    misconceptionNode.x,
                    misconceptionNode.y,
                    errorNode.x,
                    errorNode.y,
                    'node-connection'
                );
            }
        });
    });

    // Add legend
    const legend = document.createElement('div');
    legend.className = 'concept-map-legend';
    legend.innerHTML = 
  '<div class="legend-title"></div>' +
  '<div class="legend-item">' +
    '<div class="legend-icon misconception-icon"></div>' +
    '<div class="legend-label">Misconception (% of students)</div>' +
  '</div>' +
  '<div class="legend-item">' +
    '<div class="legend-icon error-icon"></div>' +
    '<div class="legend-label">Error Category</div>' +
  '</div>';

    conceptMap.appendChild(legend);

    // Add the concept map to the container
    container.appendChild(conceptMap);

    // Setup tooltips
    setupTooltips();
}

// Function to create the control panel
function createControlPanel(container, currentRadius1, currentRadius2, currentCenterX, currentCenterY) {
    const controlPanel = document.createElement('div');
    controlPanel.id = 'concept-map-controls';
    controlPanel.className = 'concept-map-controls';

    controlPanel.innerHTML = 
  '<div class="control-section">' +
    '<h4>Adjust Concept Map Layout</h4>' +
    '<div class="control-grid">' +
      '<div class="control-item">' +
        '<label for="radius1-control">Inner Circle: <span id="radius1-value">' + currentRadius1 + '</span>px</label>' +
        '<input type="range" id="radius1-control" min="80" max="200" value="' + currentRadius1 + '">' +
      '</div>' +
      '<div class="control-item">' +
        '<label for="radius2-control">Outer Circle: <span id="radius2-value">' + currentRadius2 + '</span>px</label>' +
        '<input type="range" id="radius2-control" min="150" max="400" value="' + currentRadius2 + '">' +
      '</div>' +
      '<div class="control-item">' +
        '<label for="centerX-control">Center X: <span id="centerX-value">' + currentCenterX + '</span>px</label>' +
        '<input type="range" id="centerX-control" min="300" max="900" value="' + currentCenterX + '">' +
      '</div>' +
      '<div class="control-item">' +
        '<label for="centerY-control">Center Y: <span id="centerY-value">' + currentCenterY + '</span>px</label>' +
        '<input type="range" id="centerY-control" min="200" max="500" value="' + currentCenterY + '">' +
      '</div>' +
    '</div>' +
    '<button id="reset-map-controls" class="reset-button">Reset to Default</button>' +
  '</div>';


    // Insert at the top of the container
    if (container.firstChild) {
        container.insertBefore(controlPanel, container.firstChild);
    } else {
        container.appendChild(controlPanel);
    }

    // Add event listeners for controls
    setupControlListeners();
}

// Function to set up control listeners
function setupControlListeners() {
    const radius1Control = document.getElementById('radius1-control');
    const radius2Control = document.getElementById('radius2-control');
    const centerXControl = document.getElementById('centerX-control');
    const centerYControl = document.getElementById('centerY-control');
    const resetButton = document.getElementById('reset-map-controls');

    // Update value displays during slider movement
    if (radius1Control) {
        radius1Control.addEventListener('input', function() {
            document.getElementById('radius1-value').textContent = this.value;
        });

        radius1Control.addEventListener('change', createConceptMap);
    }

    if (radius2Control) {
        radius2Control.addEventListener('input', function() {
            document.getElementById('radius2-value').textContent = this.value;
        });

        radius2Control.addEventListener('change', createConceptMap);
    }

    if (centerXControl) {
        centerXControl.addEventListener('input', function() {
            document.getElementById('centerX-value').textContent = this.value;
        });

        centerXControl.addEventListener('change', createConceptMap);
    }

    if (centerYControl) {
        centerYControl.addEventListener('input', function() {
            document.getElementById('centerY-value').textContent = this.value;
        });

        centerYControl.addEventListener('change', createConceptMap);
    }

    // Reset button
    if (resetButton) {
        resetButton.addEventListener('click', function() {
            // Set controls to default values
            if (radius1Control) radius1Control.value = 120;
            if (radius2Control) radius2Control.value = 300;
            if (centerXControl) centerXControl.value = 600;
            if (centerYControl) centerYControl.value = 350;

            // Update displayed values
            if (document.getElementById('radius1-value')) {
                document.getElementById('radius1-value').textContent = 120;
            }
            if (document.getElementById('radius2-value')) {
                document.getElementById('radius2-value').textContent = 300;
            }
            if (document.getElementById('centerX-value')) {
                document.getElementById('centerX-value').textContent = 600;
            }
            if (document.getElementById('centerY-value')) {
                document.getElementById('centerY-value').textContent = 350;
            }

            // Recreate the map
            createConceptMap();
        });
    }
}

// Helper function to add a connection line between two points
function addConnection(parent, x1, y1, x2, y2, className) {
    const length = Math.sqrt(Math.pow(x2 - x1, 2) + Math.pow(y2 - y1, 2));
    const angle = Math.atan2(y2 - y1, x2 - x1) * 180 / Math.PI;

    const connection = document.createElement('div');
connection.className = 'connection ' + className;
connection.style.width = length + 'px';
connection.style.left = x1 + 'px';
connection.style.top = y1 + 'px';
connection.style.transform = 'rotate(' + angle + 'deg)';


    parent.appendChild(connection);
}

// Setup tooltips for nodes
function setupTooltips() {
    // Create tooltip element
    const tooltip = document.createElement('div');
    tooltip.className = 'concept-tooltip';
    tooltip.style.display = 'none';
    document.body.appendChild(tooltip);

    // Add mouseover and mouseout event listeners to nodes
    const nodes = document.querySelectorAll('.concept-node');
    nodes.forEach(node => {
        node.addEventListener('mouseover', function(e) {
            const tooltipContent = this.getAttribute('data-tooltip');
            tooltip.innerHTML = tooltipContent;
            tooltip.style.display = 'block';

            // Position tooltip near mouse but not too close to avoid flickering
            tooltip.style.left = (e.pageX + 15) + 'px';
tooltip.style.top = (e.pageY - 15) + 'px';


            // Highlight connections related to this node
            if (this.classList.contains('misconception-node')) {
                // Highlight center connection
                document.querySelectorAll('.hub-connection').forEach(conn => {
                    if (parseFloat(conn.style.left) === parseFloat(this.style.left) &&
                        parseFloat(conn.style.top) === parseFloat(this.style.top)) {
                        conn.classList.add('connection-highlight');
                    }
                });

                // Highlight connections to error nodes
                document.querySelectorAll('.node-connection').forEach(conn => {
                    if (parseFloat(conn.style.left) === parseFloat(this.style.left) &&
                        parseFloat(conn.style.top) === parseFloat(this.style.top)) {
                        conn.classList.add('connection-highlight');
                    }
                });
            } else if (this.classList.contains('error-node')) {
                // Highlight connections from misconception nodes
                document.querySelectorAll('.node-connection').forEach(conn => {
                    const connAngle = parseFloat(conn.style.transform.replace('rotate(', '').replace('deg)', ''));
                    const connLength = parseFloat(conn.style.width);

                    // Calculate end point of connection
                    const startX = parseFloat(conn.style.left);
                    const startY = parseFloat(conn.style.top);
                    const angleRad = connAngle * Math.PI / 180;
                    const endX = startX + connLength * Math.cos(angleRad);
                    const endY = startY + connLength * Math.sin(angleRad);

                    // Check if this connection ends at this error node
                    if (Math.abs(endX - parseFloat(this.style.left)) < 10 &&
                        Math.abs(endY - parseFloat(this.style.top)) < 10) {
                        conn.classList.add('connection-highlight');
                    }
                });
            }
        });

        node.addEventListener('mousemove', function(e) {
            // Update tooltip position when mouse moves
            tooltip.style.left = (e.pageX + 15) + 'px';
tooltip.style.top = (e.pageY - 15) + 'px';

        });

        node.addEventListener('mouseout', function() {
            // Hide tooltip
            tooltip.style.display = 'none';

            // Remove highlight from all connections
            document.querySelectorAll('.connection-highlight').forEach(conn => {
                conn.classList.remove('connection-highlight');
            });
        });
    });
}

// Add resizable functionality
function makeConceptMapResizable() {
    let resizeTimeout;
    window.addEventListener('resize', () => {
        // Use debounce to avoid excessive redraws
        clearTimeout(resizeTimeout);
        resizeTimeout = setTimeout(() => {
            // Check if concept map tab is active
            const conceptMapTab = document.querySelector('[data-tab="concept-map"]');
            if (conceptMapTab && conceptMapTab.classList.contains('active')) {
                createConceptMap();
            }
        }, 250);
    });
}

// Call function when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    // Keep existing event handlers, add resizable functionality
    makeConceptMapResizable();
});
</script>
    <script>
// modal.js - Create this new file for the modal functionality

// Initialize modal functionality when the document is loaded
document.addEventListener('DOMContentLoaded', () => {
    // Create modal overlay
    setupModal();

    // Listen for DOM changes to catch dynamically added code previews
    setupMutationObserver();
});

// Create the modal elements and add them to the DOM
function setupModal() {
    // Create modal overlay element
    const modalOverlay = document.createElement('div');
    modalOverlay.className = 'modal-overlay';
    modalOverlay.innerHTML = 
  '<div class="modal-content">' +
    '<div class="modal-header">' +
      '<h3 class="modal-title">Code Preview</h3>' +
      '<button class="modal-close">&times;</button>' +
    '</div>' +
    '<div class="modal-body"></div>' +
  '</div>';

    document.body.appendChild(modalOverlay);

    // Close modal when clicking the close button or outside the modal
    const modalClose = modalOverlay.querySelector('.modal-close');
    modalClose.addEventListener('click', () => {
        modalOverlay.classList.remove('active');
    });

    modalOverlay.addEventListener('click', (e) => {
        if (e.target === modalOverlay) {
            modalOverlay.classList.remove('active');
        }
    });

    // Add an escape key listener to close the modal
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && modalOverlay.classList.contains('active')) {
            modalOverlay.classList.remove('active');
        }
    });

    return modalOverlay;
}

// Set up click event listeners for all code preview elements
function setupCodePreviewListeners() {
    const modalOverlay = document.querySelector('.modal-overlay');
    const modalBody = modalOverlay.querySelector('.modal-body');
    const modalTitle = modalOverlay.querySelector('.modal-title');

    const codePreviewElements = document.querySelectorAll('.code-preview');

    codePreviewElements.forEach(preview => {
        // Skip if already has a click listener
        if (preview.dataset.hasListener === 'true') {
            return;
        }

        preview.dataset.hasListener = 'true';

        preview.addEventListener('click', () => {
            // Find the nearest heading to get the context
            let heading = preview.closest('.info-card')?.querySelector('h3')?.textContent || 'Code Preview';

            // Clone the code element to show in modal
            const codeElement = preview.querySelector('pre').cloneNode(true);

            // Update modal title and content
            modalTitle.textContent = 'Code Preview: ' + heading;
            modalBody.innerHTML = '';
            modalBody.appendChild(codeElement);

            // Apply syntax highlighting if needed
            if (typeof Prism !== 'undefined') {
                Prism.highlightElement(codeElement.querySelector('code'));
            }

            // Show the modal
            modalOverlay.classList.add('active');
        });
    });
}

// Set up a mutation observer to watch for dynamically added code previews
function setupMutationObserver() {
    // Options for the observer (which mutations to observe)
    const config = { childList: true, subtree: true };

    // Callback function to execute when mutations are observed
    const callback = function(mutationsList, observer) {
        for (const mutation of mutationsList) {
            if (mutation.type === 'childList' && mutation.addedNodes.length > 0) {
                // Check if we need to set up code preview listeners
                setupCodePreviewListeners();
            }
        }
    };

    // Create an observer instance linked to the callback function
    const observer = new MutationObserver(callback);

    // Start observing the target node for configured mutations
    observer.observe(document.body, config);

    // Initial setup for already existing code previews
    setupCodePreviewListeners();
}
</script>
  </body>
</html>`

var PROBLEM_DASHBOARD_TEMPLATE = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>Problem Dashboard</title>
<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
<script src="https://www.kryogenix.org/code/browser/sorttable/sorttable.js"></script>

<style>
	.menu {
		padding: 10px;
		padding-left: 100px;
	}
	.topcorner{
		position:absolute;
		top:0;
		right:0;
	}
	#deactivate-button {
		background-color: linear-gradient(45deg, rgb(255, 0, 0), rgb(255, 102, 102));
		color: #FFFFFF;
	}
#deactivate-button:hover {
  background: linear-gradient(45deg, rgb(255, 102, 102), rgb(255, 0, 0));
  color: white; /* Ensures text remains white on hover */
}
#feedback-dropdown {
    width: 100%; /* Make it span the full width of the container */
    max-width: 250px; /* Set a maximum width */
    padding: 10px; /* Add padding to the dropdown */
    border-radius: 5px; /* Rounded corners */
    border: 1px solid #ccc; /* Light border */
    background: #ffffff; /* White background */
    box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.1); 
}
.content {
    margin-top: 30px; 
}

#scaffolding-dropdown {
    width: 100%; /* Make it span the full width of the container */
    max-width: 250px; /* Set a maximum width */
    padding: 10px; /* Add padding to the dropdown */
    border-radius: 5px; /* Rounded corners */
    border: 1px solid #ccc; /* Light border */
    background: #ffffff; /* White background */
    box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.1); 
	display: none;
}

.button.is-primary {
  background: cornflowerblue;
  color: white; /* Ensures text is white */
  border: none;
}

.button.is-primary:hover {
  background: linear-gradient(45deg, #2575fc, #6a11cb);
  color: white; /* Ensures text remains white when hovered */
}

.button.is-info {
  background: cornflowerblue;
  color: white; /* Ensures text is white */
  border: none;
}

.button.is-info:hover {
  background: linear-gradient(45deg, #00bcd4, #1e90ff);
  color: white; /* Ensures text remains white when hovered */
}

.button.is-success {
  background: cornflowerblue;
  color: white; /* Ensures text is white */
  border: none;
}

.button.is-success:hover {
  background: linear-gradient(45deg, #2575fc, #6a11cb);
  color: white; /* Ensures text remains white when hovered */
}

.button.is-info {
  background: cornflowerblue;
  color: white; /* Ensures text is white */
  border: none;
}

.button.is-info:hover {
  background: linear-gradient(45deg, #00bcd4, #1e90ff);
  color: white; /* Ensures text remains white when hovered */
}

td {
  border: none; /* Removes border if you don't want one */
  vertical-align: middle; /* Centers text vertically */
  padding: 10px; /* Adds some spacing for better appearance */
}


.progress-container {
  width: 100px; /* Adjust as needed */
  height: 10px;
  background: #e0e0e0;
  border-radius: 6px;
  overflow: hidden;
  display: flex;
  align-items: center; /* Centers progress bar inside container */
  margin-top: 5px; 
}

.progress-bar {
  height: 100%;
  transition: width 0.5s ease-in-out;
}


.modal {
    display: none; 
    align-items: center;  /* Vertically center */
    justify-content: center; /* Horizontally center */
  }

.modal.is-active {
    display: flex !important;
}

</style>
</head>
<body>
<div class="container">
<nav class="navbar is-fixed-top breadcrumb menu" role="navigation" aria-label="breadcrumbs">
<ul>
  <li>
	<a id="view-exercise-link" href="#">
	  <span class="icon is-small">
		<i class="fas fa-home" aria-hidden="true"></i>
	  </span>
	  <span>Exercises</span>
	</a>
  </li>
  <li class="is-active">
	<a href="#">
	  <span class="icon is-small">
		<i class="fas fa-book" aria-hidden="true"></i>
	  </span>
	  <span>Exercise Dashboard</span>
	</a>
  </li>
</ul>
</nav>
<nav class="breadcrumb is-right" aria-label="breadcrumbs">
		<ul>
		<li class="is-active"><a href="#">{{.Username}}({{.UserRole}})</a></li>
		</ul>
  	</nav>
<div class="content">
	{{if eq .UserRole "teacher"}} 
		{{if eq .IsActive true}}
			<button id="deactivate-button" class="button is-danger">Deactivate</button>
		{{end}}
	{{end}}
	{{if eq .UserRole "teacher"}} 
	<button id="api-call-button" class="button is-primary">Generate Summary</button>
{{end}}
{{if eq .UserRole "teacher"}} 
    <button id="student-progress-button" class="button is-primary">Estimate Student Progress</button>
{{end}}
{{if eq .UserRole "teacher"}}
    <select id="scaffolding-dropdown">
        <option value="">Generate Scaffolds</option>
        <option value="1">Fill-in-the-Blanks</option>
        <option value="2">Step-by-Step Tasks</option>
        <option value="3">Guided Code with Hints</option>
        <option value="4">Debug This Code</option>
        <option value="5">Incremental Feature Implementation</option>
    </select>
{{end}}
{{if eq .UserRole "teacher"}} 
    <button id="scaffolding-view-button" class="button is-primary">Generate Scaffolds</button>
{{end}}
	<div class="accordions" style="margin-top: 10px;">
		<h3>{{.ProblemName}}</h3>
		<div>
			<textarea id="editor">{{ .Code }}</textarea>
		</div>
	</div>
 <!-- Feedback Box -->
	<table class="table">
			<thead>
				<tr>
					<th>Active Students</th>
					<th>Help Requests</th>
					<th>Not Graded</th>
					<th>Correct</th>
					<th>Incorrect</th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<td>{{.NumActive}}</td>
					<td>{{.NumHelpRequest}}</td>
					<td>{{.NumNotGraded}}</td>
					<td>{{.NumGradedCorrect}}</td>
					<td>{{.NumGradedIncorrect}}</td>
				</tr>
			</tbody>
	</table>
	{{if gt (len .AnswerStats) 0}}
		<table>
			<thead>
				<tr>
					<th>Answer</th>
					<th>Student submitted</th>
				</tr>
			</thead>
			<tbody>
				{{range .AnswerStats}}
				<tr>
				<td>{{.Answer}}</td>
				<td>{{.Percent}}%</td>
				</tr>
				{{end}}
			</tobdy>
		</table>
	{{end}}
	<table class="table sortable">
			<thead>
				<tr>
					<th>Student</th>
					<th>Active</th>
					<th>Status</th>
					<th>Help Status</th>
					<th>Submission</th>
					<th>Progress</th>
				</tr>
			</thead>
			<tbody>
				{{range .StudentInfo}}
				<tr>
					<td class="{{if ne .CodingStat "Idle"}}active{{end}}">{{if ne .CodingStat "Idle"}}<a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{$.ProblemID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}#code-snapshot">{{.StudentName}}</a>{{else}}{{.StudentName}}{{end}}</td>
					<td class="{{if ne .CodingStat "Idle"}}active{{end}}">{{if and (eq $.IsActive true) (ne .CodingStat "Idle") (ne .LastUpdatedAt.IsZero true) }}<a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{$.ProblemID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}">{{ formatTimeSince .LastUpdatedAt }} ago</a>{{end}}</td>
					<td class="{{if ne .CodingStat "Idle"}}active{{end}}">{{.CodingStat}}</td>
					<td class="{{if ne .CodingStat "Idle"}}active{{end}}">{{if ne .HelpStat ""}}<a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{$.ProblemID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}#ask-for-help">{{.HelpStat}}</a>{{end}}</td>
					<td class="{{if ne .CodingStat "Idle"}}active{{end}}">{{if ne .SubmissionStat ""}}<a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{$.ProblemID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}#submission">{{.SubmissionStat}}</a>{{end}}</td>
					<td class="{{if ne .CodingStat "Idle"}}active{{end}}">
            <div class="progress-container">
                <div class="progress-bar" style="width: {{.Percentage}}%; background-color: 
                    {{if lt .Percentage 50}}#ff4d4d{{else if lt .Percentage 60}}#ffcc00{{else if lt .Percentage 80}}#ffeb3b{{else}}#008000{{end}};">
                </div>
            </div>
        </td>
				</tr>
				{{end}}
			</tbody>
	</table>
	</div>
	<script>
		var editor = document.getElementById("editor");
		var myCodeMirror = CodeMirror.fromTextArea(editor, {lineNumbers: true, mode: get_editor_mode({{.ProblemName}}), theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
		myCodeMirror.setSize("100%", 400)
		function get_editor_mode(filename) {
			filename = filename.toLowerCase();
			if (filename.endsWith('.py')) {
				return "python";
			}
			if (filename.endsWith('.java')) {
				return "text/x-java";
			}
			if (filename.endsWith('.cpp') || filename.endsWith('.c++') || filename.endsWith('.c')) {
				return "text/x-c++src";
			}
			return "text";
		  }

$('#api-call-button').on('click', function () {
  handleSummarySubmission(false);
});


function handleSummarySubmission(isNew) {

  const $button = $('#api-call-button');
  $button.text('Generating...');

  const urlParams = new URLSearchParams(window.location.search);
  const problemId = urlParams.get('problem_id');

  const tempForm = $('<form>', {
    method: 'POST',
    action: '/vi_view?problem_id=' + encodeURIComponent(problemId)
  });

  tempForm.append($('<input>', {
    type: 'hidden',
    name: 'problem_id',
    value: problemId
  }));

  tempForm.append($('<input>', {
    type: 'hidden',
    name: 'new',
    value: isNew
  }));

  $('body').append(tempForm);
  tempForm.submit();
}



		  $(document).ready(function(){
			$('#view-exercise-link').attr("href", "/view_exercises"+window.location.search);
			var problemId = new URLSearchParams(window.location.search).get('problem_id');
			$('#deactivate-button').click(function(){
				if (confirm("After you deactivate this exercise, students can no longer submit their work to this exercise.  Further, you cannot undo this action.  If you want to deactivate this exercise, click OK.") == true) {
					$.post("/teacher_deactivates_problems", {filename: {{.ProblemName}}, problem_id: problemId, uid: {{.UserID}}, role: {{.UserRole}}{{if ne .Password ""}}, password: {{.Password}}{{end}} })
					.done(function(data){
						if (data == "-1"){
							alert("Couldn't deactivate the problem! Please try again!");
						} else {
							alert("Problem deactiavated!");
							window.location.reload();
						}
					})
					.fail(function(){
						alert("Couldn't deactivate the problem! Please try again!");
					});
				}
			});
		  });
		  $(".accordions").accordion({ header: "h3", active: false, collapsible: true });
		  $(".accordions").show();
document.addEventListener("DOMContentLoaded", function() {
    const urlParams = new URLSearchParams(window.location.search);
    const problemId = urlParams.get('problem_id');

    // Fetch the feedback list when the page loads
    fetch('/get_feedback_list', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ problem_id: problemId })
    })
    .then(response => response.json())
    .then(data => {
        console.log('API Response:', data);

        // Populate the dropdown
        const dropdown = document.getElementById("feedback-dropdown");
        data.feedbacks.forEach(feedback => {
            const option = document.createElement("option");
            option.value = feedback.feedback_id;  // Store feedback_id as the value
            option.textContent = new Date(feedback.feedback_time).toLocaleString();
            dropdown.appendChild(option);
        });
    })
    .catch(error => {
        console.error('Error:', error);
    });

    // Event listener for dropdown selection
    document.getElementById("feedback-dropdown").addEventListener("change", function(event) {
        const selectedFeedbackID = event.target.value;
        const feedbackBox = document.getElementById("custom-prompt-box");
        const feedbackTextArea = document.getElementById("custom-prompt-response");
        const loadingText = document.getElementById("loading-text");

        if (selectedFeedbackID) {
            // Show the box and loading text
            feedbackBox.style.display = "block";
            loadingText.style.display = "block";
            feedbackTextArea.style.display = "none"; // Hide textarea while fetching

            fetch('/get_feedback_by_id', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ feedback_id: selectedFeedbackID })
            })
            .then(response => response.json())
            .then(data => {
                // Hide loading text and show the textarea with the response
                loadingText.style.display = "none";
                feedbackTextArea.style.display = "block";
                feedbackTextArea.value = data.feedback; // Set response text
            })
            .catch(error => {
                console.error('Error fetching feedback:', error);
                loadingText.textContent = "Error fetching feedback.";
            });
        } else {
            // Hide the box if no feedback is selected
            feedbackBox.style.display = "none";
        }
    });
});
document.getElementById('student-progress-button').addEventListener('click', function() {
    // Get the button element
    const button = document.getElementById('student-progress-button');
    
    // Change button text to "Generating..."
    button.textContent = 'Generating...';
    
    // Prepare the request data
    const urlParams = new URLSearchParams(window.location.search);
    const problemId = urlParams.get('problem_id');
    
    const requestData = {
        problem_id: problemId,
    };

    // Send the API request using fetch
    fetch('/summarize_student_progress', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData),
    })
    .then(response => response.json())
    .then(data => {
        console.log('API Response:', data);

        // Restore button text to normal after the request is completed
        button.textContent = 'Generate Students Progress';

        // Check for the correct property in the response (message)
        if (data.message && data.message !== "") {
            // Show a success prompt if the message is valid
            alert(data.message);
			location.reload();
        } else {
            // Show an alert if no message is found
            alert("No progress data found!");
        }
    })
    .catch(error => {
        // Restore button text to normal in case of error
        button.textContent = 'Generate Students Progress';
        
        // Show an error alert if the request fails
        alert("Error: " + JSON.stringify(error));
    });
});

document.getElementById("scaffolding-dropdown").addEventListener("change", sendScaffoldingStrategy);

function sendScaffoldingStrategy() {
    const dropdown = document.getElementById("scaffolding-dropdown");
    const selectedStrategy = dropdown.value;
    if (!selectedStrategy) return;

    // Retrieve problem_id from URL
    const urlParams = new URLSearchParams(window.location.search);
    const problemId = urlParams.get("problem_id");

    if (!problemId) {
        alert("Error: Problem ID is missing!");
        return;
    }

    // Disable the dropdown and show "Generating..."
    dropdown.disabled = true;
    const originalText = dropdown.options[dropdown.selectedIndex].text;
    dropdown.options[dropdown.selectedIndex].text = "Generating...";

    const requestData = {
        problem_id: problemId,
        scaffolding_strategy: selectedStrategy
    };

    fetch("/process_scaffolding", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(requestData)
    })
    .then(response => response.json().then(data => ({ status: response.status, body: data })))
    .then(({ status, body }) => {
        if (status === 409) {
            alert(body.message); // Show message from backend (Scaffolding already exists)
        } else if (status === 200) {
            alert("Scaffolding strategy submitted successfully!");
        } else {
            alert("Unexpected response from server!");
        }
    })
    .catch(error => {
        console.error("Error:", error);
        alert("Error submitting scaffolding strategy!");
    })
    .finally(() => {
        // Re-enable the dropdown and restore the original text
        dropdown.disabled = false;
        dropdown.options[dropdown.selectedIndex].text = originalText;
    });
}
 document.addEventListener("DOMContentLoaded", function() {
        var button = document.getElementById("scaffolding-view-button");
        if (button) {
            button.addEventListener("click", function() {
                // Get the current URL parameters
                var urlParams = new URLSearchParams(window.location.search);
                var problemId = urlParams.get("problem_id");
                var uid = urlParams.get("uid");
                var role = urlParams.get("role");

                if (problemId && uid && role) {
                    window.location.href = "/scaffolding_dashboard?problem_id=" + problemId + "&uid=" + uid + "&role=" + role;
                } else {
                    alert("Missing required parameters in the URL.");
                }
            });
        }
    });
	</script>
</body>
</html>
`
var PROBLEM_LIST_TEMPLATE = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>Exercises</title>
<meta http-equiv="refresh" content="10000000" >
<style>
.switch {
  position: relative;
  display: inline-block;
  width: 60px;
  height: 34px;
}

.switch input { 
  opacity: 0;
  width: 0;
  height: 0;
}

.button.is-primary {
  background: cornflowerblue !important;
  color: white; /* Ensures text is white */
  border: none;
}

.button.is-primary:hover {
  background: linear-gradient(45deg, #2575fc, #6a11cb) !important;
  color: white; /* Ensures text remains white when hovered */
}

.button.is-info {
  background: cornflowerblue !important;
  color: white; /* Ensures text is white */
  border: none;
}

.button.is-info:hover {
  background: linear-gradient(45deg, #00bcd4, #1e90ff);
  color: white; /* Ensures text remains white when hovered */
}

.button.is-success {
  background: cornflowerblue !important;
  color: white; /* Ensures text is white */
  border: none;
}

.button.is-success:hover {
  background: linear-gradient(45deg, #2575fc, #6a11cb) !important;
  color: white; /* Ensures text remains white when hovered */
}

.button.is-info {
  background: cornflowerblue !important;
  color: white; /* Ensures text is white */
  border: none;
}

.button.is-info:hover {
  background: linear-gradient(45deg, #00bcd4, #1e90ff) !important;
  color: white; /* Ensures text remains white when hovered */
}

td {
  border: none; /* Removes border if you don't want one */
  text-align: center !important; /* Centers text horizontally */
  vertical-align: middle; /* Centers text vertically */
  padding: 10px; /* Adds some spacing for better appearance */
}

tr {
  text-align: center !important; /* Centers text horizontally */
  vertical-align: middle; /* Centers text vertically */
}

td.active {
  background: white; /* Gradient for active problem */
  color: black !important;
}

td.inactive {
  background: #f6f6f6; /* White background for inactive problem */
}



.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  -webkit-transition: .4s;
  transition: .4s;
}

.slider:before {
  position: absolute;
  content: "";
  height: 26px;
  width: 26px;
  left: 4px;
  bottom: 4px;
  background-color: white;
  -webkit-transition: .4s;
  transition: .4s;
}

input:checked + .slider {
  background-color: #2196F3;
}

input:focus + .slider {
  box-shadow: 0 0 1px #2196F3;
}

input:checked + .slider:before {
  -webkit-transform: translateX(26px);
  -ms-transform: translateX(26px);
  transform: translateX(26px);
}

/* Rounded sliders */
.slider.round {
  border-radius: 34px;
}

.slider.round:before {
  border-radius: 50%;
}

.menu {
	padding: 10px;
	padding-left: 100px;
}

.topcorner{
	position:absolute;
	top:0;
	right:0;
}

.drawer {
    display: none;
    position: fixed;
    top: 0;
    right: 0;
    width: 350px;
    height: 100%;
    background-color: #f4f4f4;
    box-shadow: -2px 0 5px rgba(0, 0, 0, 0.5);
    padding: 20px;
    z-index: 1000;
    transition: transform 0.3s ease;
    overflow-y: auto;
}

.drawer.open {
    display: block;
    transform: translateX(0);
}

h4 {
    font-size: 16px;
    font-weight: bold;
    margin-bottom: 5px;
}

.input {
    width: 100%;
    padding: 8px;
    margin-bottom: 10px;
    border: 1px solid #ccc;
    border-radius: 4px;
}

button {
    width: 100%;
    margin-bottom: 15px;
}

.logout-container {
    position: absolute;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    width: 90%;
}

.drawer-content {
  padding-top: 50px;
}

.settings-button {
  position: fixed;
  top: 10px; /* Adjust as needed */
  left: 20px; /* Move to the left side */
  background: cornflowerblue; /* Apply gradient */
  color: white; /* Ensures text is white */
  border: none;
  padding: 12px;
  border-radius: 50%;
  cursor: pointer;
  font-size: 18px;
  z-index: 1100;
  width: 45px; /* Ensuring button size is appropriate */
  height: 45px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
  transition: background 0.3s ease; /* Smooth transition for background */
}

.settings-button:hover {
  background: linear-gradient(45deg, #2575fc, #6a11cb); /* Inverted gradient on hover */
}

.settings-button.is-disabled {
  opacity: 0.5;
  cursor: not-allowed;
  pointer-events: none;
}

a.is-disabled {
  pointer-events: none;
  opacity: 0.5;
}


</style>
<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
<script src="https://www.kryogenix.org/code/browser/sorttable/sorttable.js"></script>

</head>
<body>
<div class="container">
	
	
	<nav class="navbar is-fixed-top breadcrumb menu" role="navigation" aria-label="breadcrumbs">
		<ul>
			<li class="is-active">
			<a href="#">
				<span class="icon is-small">
				<i class="fas fa-home" aria-hidden="true"></i>
				</span>
				<span>Exercises</span>
			</a>
			</li>
		</ul>

	</nav>
	<div class="content">
	<div class="topcorner" style="margin-top: 78px; margin-bottom: 10px;">{{.Username}}({{.UserRole}})</div>
		<a id="new-problem" 
   class="button is-success {{if eq .UserRole "student"}}is-disabled{{end}}" 
   href="{{if ne .UserRole "student"}}/your-link{{end}}" 
   style="margin-top: 80px; margin-bottom: 10px;">
	<span class="icon is-small">
		<i class="fa-solid fa-plus"></i>
	</span>
	<span>Add a New Exercise</span>
</a>

<a id="export-button" 
   class="button is-primary {{if eq .UserRole "student"}}is-disabled{{end}}" 
   href="{{if ne .UserRole "student"}}/your-export-link{{end}}" 
   style="margin-top: 80px; margin-bottom: 10px;">
	<span class="icon is-small">
		<i class="fa-solid fa-plus"></i>
	</span>
	<span>Export Performance Data</span>
</a>


		<div class="drawer" id="settings-drawer" style="font-family: Arial, sans-serif; padding: 20px; background: white; border-radius: 10px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);">

    <h3 style="font-size: 20px; font-weight: bold; margin-bottom: 20px; color: #333;">Settings</h3>

    <hr>

	<div>
        <h4>Add TA to Course</h4>
        <input type="text" id="teacher-name" class="input" placeholder="Enter Teacher Name">
		<input type="text" id="teacher-pass" class="input" placeholder="Enter Teacher Password">
        <button class="button is-primary" id="add-teacher-btn">Add Teacher</button>
    </div>

	 <hr>

<div>
    <h4>Add Students to Course</h4>
    <input id="student-names" class="input" placeholder="Enter Comma Seperated Names">
    <button class="button is-info" id="add-students-btn">Add Students</button>
</div>

<hr style="margin: 15px 0;">
	<button class="button is-primary" id="toggle-peer-tutoring" style="margin-top: 3px; display: flex; align-items: center;">
    Peer Tutoring
    <label class="switch" style="margin-left: 10px;">
        <input id="peer_tutoring_button" type="checkbox">
        <span class="slider round"></span>
    </label>
</button>



    <hr>


    <button class="button is-danger" id="logout-button" style="width: 100%; display: flex; align-items: center; justify-content: center;">
        <i class="fas fa-sign-out-alt" style="margin-right: 10px;"></i> Logout
    </button>

</div>
		<button 
	class="settings-button {{if eq .UserRole "student"}}is-disabled{{end}}" 
	id="settings-button"
	{{if eq .UserRole "student"}}disabled{{end}}>
	<i class="fas fa-cogs"></i>
</button>


		<table class="table sortable">
				<thead>
					<tr>
						<th>Filename</th>
						<th>Posted At</th>
						<th>Attendance</th>
						<th>Active Students</th>
						<th>Help Requests</th>
						<th>Correct</th>
						<th>Incorrect</th>
						<th>Not Graded</th>
					</tr>
				</thead>
				<tbody>
					{{range .Problems}}
  <tr {{if eq .IsActive true}}class="is-selected"{{end}}>
    <td class="{{if eq .IsActive false}}inactive{{else}}active{{end}}">
      <a href="/analyse_view?problem_id={{.ID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}">
        {{.Filename}}
      </a>
    </td>
    <td class="{{if eq .IsActive false}}inactive{{else}}active{{end}}">
      {{ .UploadedAt.Format "Jan 02, 2006 3:04:05 PM" }}
    </td>
    <td class="{{if eq .IsActive false}}inactive{{else}}active{{end}}">{{.Attendance}}</td>
    <td class="{{if eq .IsActive false}}inactive{{else}}active{{end}}">{{.NumActive}}</td>
    <td class="{{if eq .IsActive false}}inactive{{else}}active{{end}}">{{.NumHelpRequest}}</td>
    <td class="{{if eq .IsActive false}}inactive{{else}}active{{end}}">{{.NumGradedCorrect}}</td>
    <td class="{{if eq .IsActive false}}inactive{{else}}active{{end}}">{{.NumGradedIncorrect}}</td>
    <td class="{{if eq .IsActive false}}inactive{{else}}active{{end}}">{{.NumNotGraded}}</td>
  </tr>
  {{end}}
				</tbody>
		</table>
	</div>
</div>
<script>
$(document).ready(function(){
	// Toggling the drawer visibility
	$('#settings-button').click(function(){
			let urlParams = new URLSearchParams(window.location.search);
			let courseID = urlParams.get("course_id");
			let teacherID = urlParams.get("uid");
		  window.location.href = "/settings_view?course_id=" + courseID + "&teacher_id=" + teacherID;
	});
	$('#logout-button').click(function(){
        if (confirm("Are you sure you want to logout?")) {
        window.location.href = "/logout"; 
    }
    });
	$('#add-course-btn').click(function(){
    let courseID = $('#course-id').val().trim();

    if (courseID === "") {
        alert("Please enter Course ID.");
        return;
    }

    $.ajax({
        url: "/add_course",
        type: "POST",
        contentType: "application/json",
        dataType: "json",
        data: JSON.stringify({ course_id: courseID }),
        success: function(response) {
            if (response && response.message) {
                alert(response.message); // Show success message
                
                // Clear input field
                $('#course-id').val('');
            } else {
                alert("Unexpected response format.");
            }
        },
        error: function(xhr) {
            alert("Error: " + xhr.responseText);
        }
    });
});


	// Add TA to Course
    $('#add-teacher-btn').click(function(){
    let teacherName = $('#teacher-name').val().trim();
    let teacherPass = $('#teacher-pass').val().trim();
    let urlParams = new URLSearchParams(window.location.search);
	let courseId = urlParams.get("course_id");

    if (teacherName === "" || teacherPass === "" || !courseId) {
        alert("Please enter Teacher Name and Password");
        return;
    }

    $.ajax({
        url: "/add_teacher",
        type: "POST",
        contentType: "application/json",
        dataType: "json",  // Ensure response is treated as JSON
        data: JSON.stringify({ 
            teacher_name: teacherName, 
            teacher_pass: teacherPass, 
            course_id: courseId
        }),
        success: function(response) {
            if (response && response.message) {
                alert(response.message); // Ensure proper handling of response
            } else {
                alert("Teacher added successfully, but response is missing data.");
            }

            // Clear input fields after successful addition
            $('#teacher-name').val('');
            $('#teacher-pass').val('');
        },
        error: function(xhr, status, error) {
            alert("Failed to add teacher: " + xhr.responseText);
        }
    });
});



	// Add Student to Course
   $('#add-students-btn').click(function() {
        let studentNames = $('#student-names').val().trim();
        let urlParams = new URLSearchParams(window.location.search);
        let courseID = urlParams.get("course_id");

        if (studentNames === "" || courseID === "") {
            alert("Please enter student names and ensure Course ID is available.");
            return;
        }

        let studentsArray = studentNames.split(',').map(name => name.trim()).filter(name => name !== "");

        if (studentsArray.length === 0) {
            alert("Please enter valid student names.");
            return;
        }

        $.ajax({
            url: "/add_students",
            type: "POST",
            contentType: "application/json",
            dataType: "json",
            data: JSON.stringify({ 
                student_names: studentsArray, 
                course_id: courseID 
            }),
            success: function(response) {
                alert(response.message || "Students added successfully");
                $('#student-names').val(''); // Clear input field
            },
            error: function(xhr, status, error) {
                alert("Failed to add students: " + xhr.responseText);
            }
        });
    });

	// Peer tutoring functionality (unchanged)
	{{if eq .PeerTutorAllowed true}}$('#peer_tutoring_button').prop('checked', true);{{end}}
	$('#new-problem').attr("href", "/teacher_web_broadcast"+window.location.search);
	$('#peer_tutoring_button').change(function(){
		var val = document.getElementById('peer_tutoring_button').checked;
		var valInt = val ? 1 : 0;
		$.post("/set_peer_tutor", {turn_on: valInt, uid: {{.UserID}}, role: {{.UserRole}}{{if ne .Password ""}}, password: {{.Password}}{{end}} }, function(data, status){
		});
	});
	$('#export-button').click(function(){

		$.ajax({
			type: "POST",
			url: "/teacher_exports_point",
			data: {uid: {{.UserID}}, role: {{.UserRole}}{{if ne .Password ""}}, password: {{.Password}}{{end}} },
			success: function(response, status, xhr) {
				// check for a filename
				var filename = "";
				var disposition = xhr.getResponseHeader('Content-Disposition');
				if (disposition && disposition.indexOf('attachment') !== -1) {
					var filenameRegex = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/;
					var matches = filenameRegex.exec(disposition);
					if (matches != null && matches[1]) filename = matches[1].replace(/['"]/g, '');
				}

				var type = xhr.getResponseHeader('Content-Type');
				var blob = new Blob([response], { type: type });

				if (typeof window.navigator.msSaveBlob !== 'undefined') {
					window.navigator.msSaveBlob(blob, filename);
				} else {
					var URL = window.URL || window.webkitURL;
					var downloadUrl = URL.createObjectURL(blob);

					if (filename) {
						var a = document.createElement("a");
						if (typeof a.download === 'undefined') {
							window.location = downloadUrl;
						} else {
							a.href = downloadUrl;
							a.download = filename;
							document.body.appendChild(a);
							a.click();
						}
					} else {
						window.location = downloadUrl;
					}

					setTimeout(function () { URL.revokeObjectURL(downloadUrl); }, 100);
				}
			}
		});
	});
	
});
</script>
</body>
</html>
`

var EXERCISE_LIST_TEMPLATE = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>Exercises</title>
<meta http-equiv="refresh" content="10000000">
<style>
body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f9fafb;
  color: #333;
  margin: 0;
  padding: 0;
}

.container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.user-info {
  position: absolute;
  top: 20px;
  right: 20px;
  font-size: 14px;
  background-color: cornflowerblue;
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
}

.page-title {
  text-align: center;
  margin-top: 70px;
  margin-bottom: 30px;
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.exercise-list {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  overflow: hidden;
}

.exercise-item {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f4f8;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.exercise-item:last-child {
  border-bottom: none;
}

.exercise-item:hover {
  background-color: #f7fafc;
}

.exercise-item.active {
  background-color: #ebf5ff;
}

.exercise-name {
  flex-grow: 1;
}

.exercise-link {
  text-decoration: none;
  color: #4a5568;
  font-weight: 500;
  display: block;
  width: 100%;
}

.exercise-link:hover {
  color: cornflowerblue;
}

.exercise-date {
  color: #718096;
  font-size: 14px;
  min-width: 180px;
  text-align: right;
}

.action-buttons {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}
</style>
<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
</head>
<body>
<div class="container">
  <div class="user-info">
    {{.Username}} ({{.UserRole}})
  </div>
  
  <div class="content">
    <div class="page-title">
      Available Exercises
    </div>
    
    {{if ne .UserRole "student"}}
    <div class="action-buttons">
      <a id="new-problem" class="button is-success" href="">
        <span>Add Exercise</span>
      </a>
      <a id="export-button" class="button is-primary" href="">
        <span>Export Data</span>
      </a>
    </div>
    {{end}}
    
    <div class="exercise-list">
      {{range .Problems}}
      <div class="exercise-item {{if eq .IsActive true}}active{{end}}">
        <div class="exercise-name">
          <a href="/view_user_feedback?student_id={{$.UserID}}&problem_id={{.ID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}" class="exercise-link">
            {{.Filename}}
          </a>
        </div>
        <div class="exercise-date">
          {{ .UploadedAt.Format "Jan 02, 2006 3:04 PM" }}
        </div>
      </div>
      {{end}}
    </div>
  </div>
</div>
</body>
</html>
`

var EXERCISE_FEEDBACK_TEMPLATE = `
	<!DOCTYPE html>
	<html>
	<head>
	<title>Student Dashboard</title>
	<meta http-equiv="refresh" content="120" >
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />

	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	<script src="https://cdn.jsdelivr.net/npm/@creativebulma/bulma-collapsible"></script>
	<style>
		.status {
			display: flex;
			justify-content: space-between;
		}
		.menu {
			padding: 10px;
			padding-left: 100px;
			padding-right: 100px;
		}
		.show {
			top: 6%;
			position: fixed;
			z-index: 200;
			background: white;
		}
		.content {
			padding-top: 3%;
		}
		.topcorner{
			position:absolute;
			top:0;
			right:0;
		}
	</style>
	</head>
	<body>
	<div class="container">
	<nav class="navbar is-fixed-top breadcrumb menu" role="navigation" aria-label="breadcrumbs">
	<div class="navbar-start"> 
	<ul>
	  <li>
		<a id="view-exercise-link" href="#">
		  <span class="icon is-small">
			<i class="fas fa-home" aria-hidden="true"></i>
		  </span>
		  <span>Exercises</span>
		</a>
	  </li>
	</ul>
	</div>
	<div class="navbar-end"> 
		<div class="navbar-item"> <a href="#">{{.Username}} ({{.UserRole}})</a> </div>
	</div>
	</nav>
	<div class="content">
	<div class="column is-two-thirds show" style="width: 70%;">
	</div>
		
	<div class="content">
    <div>
        <section class="section" style="padding: 20px">
            <h2 class="title is-3" style="padding-left: 20px">Feedback History</h2>
            {{if .Messages}}
                {{range .Messages}}
                    <article class="message" style="margin-left: 25px; padding-bottom: 20px;">
                        <div class="message-header">
                            <p>{{if eq .Type 0}}{{.Name}} asked for help{{else if eq .Event "at_submission"}} Submission Snapshot taken {{else}} Regular Snapshot taken {{end}} at ({{.GivenAt.Format "Jan 02, 2006 3:04:05 PM"}})</p>
                        </div>
                        <div class="message-body">
                            {{.Message}}
                        </div>
                        <div style="margin-left:20px;">
                            {{if .Code }}
                                {{range .Feedbacks}}
                                    <article class="message" style="margin-left: 25px;">
                                        <div class="message-header">
                                            <p>Reply from {{.Name}} given at {{.GivenAt.Format "Jan 02, 2006 3:04:05 PM"}} </p>
                                        </div>
                                        <div class="message-body">
                                            <div class="columns">
                                                <div class="column is-four-fifths">
                                                    <textarea class="message-feedback">{{ .Feedback }}</textarea>
                                                </div>
                                                {{ if not (eq .Upvote 0) }}
                                                <div class="column" style="text-align: center;">
                                                    <div style="font-size: 32px;">
                                                        {{.Upvote}}
                                                    </div>
                                                    <p> Student found it helpful.</p>
                                                </div>
                                                {{ end }}
                                            </div>
                                        </div>
                                    </article>
                                {{end}}
                            {{ end }}
                        </div>
                    </article>
                {{end}}
            {{else}}
                <div class="notification is-info is-light" style="margin-left: 25px; margin-right: 25px;">
                    <p>No feedback available</p>
                </div>
            {{end}}
        </section>
    </div>
</div>
	</div>
	<script>
		$(document).ready(function(){
			$('#view-exercise-link').attr("href", "/view_feedback"+window.location.search);
			$('#problem-dashboard-link').attr("href", "/problem_dashboard"+window.location.search+"&problem_id={{.ProblemID}}");
		

			var snapshotEditors = document.getElementsByClassName("message-feedback");
			
			for (let i = 0; i<snapshotEditors.length; i++){
				var code = CodeMirror.fromTextArea(snapshotEditors[i], {lineNumbers: true, mode: "{{getEditorMode .ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
				code.setSize("100%", 500);
			}
		});

		function autoFeedbackSubmit(backFeedback, fID) {
			$.ajax({
				url: "/save_snapshot_back_feedback",
				type: "POST",
				data:  {
					feedback: backFeedback,
					feedback_id: fID,
					uid: {{.UserID}},
					role: "{{.UserRole}}",
					{{if ne .Password ""}}password: "{{.Password}}",{{end}}
				},
				success: function(data){
					console.log("Success!")
				}
			});
			
			location.reload();
		}
	</script>
	</body>
	</html>
`

var STUDENT_PROBLEM_LIST_TEMPLATE = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>My Feedback</title>
  <meta http-equiv="refresh" content="10000000">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />
  <script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
  <style>
    /* Basic styles for a pleasing accordion */
    .collapsible {
      background-color: #f6f6f6;
      color: #333;
      cursor: pointer;
      padding: 18px;
      width: 100%;
      border: none;
      text-align: left;
      outline: none;
      font-size: 18px;
      transition: background-color 0.3s ease;
      border-radius: 5px;
      margin-bottom: 10px;
    }
    .collapsible:hover, .collapsible.active {
      background-color: #e2e2e2;
    }
    .content {
      padding: 0 18px;
      max-height: 0;
      overflow: hidden;
      transition: max-height 0.2s ease-out;
      background-color: #fff;
      border: 1px solid #ddd;
      border-radius: 5px;
      margin-bottom: 10px;
    }
    .feedback-item {
      border-bottom: 1px solid #eee;
      padding: 10px 0;
    }
    .feedback-item:last-child {
      border-bottom: none;
    }
    .feedback-header {
      font-weight: bold;
    }
  </style>
</head>
<body>
  <section class="section">
    <div class="container">
      <h1 class="title">My Feedback</h1>
      {{range .Problems}}
        <button class="collapsible">
          {{.Filename}} <span style="font-size: 14px; color: #666;">(Posted at {{.UploadedAt.Format "Jan 02, 2006 3:04:05 PM"}})</span>
        </button>
        <div class="content">
          <div class="content">
            <p><strong>Attendance:</strong> {{.Attendance}}</p>
            <p><strong>Active Students:</strong> {{.NumActive}}</p>
            <p><strong>Help Requests:</strong> {{.NumHelpRequest}}</p>
            <p><strong>Graded Correct:</strong> {{.NumGradedCorrect}}</p>
            <p><strong>Graded Incorrect:</strong> {{.NumGradedIncorrect}}</p>
            <p><strong>Not Graded:</strong> {{.NumNotGraded}}</p>
            <hr>
            <h2 class="subtitle">Feedbacks</h2>
            {{if .Feedbacks}}
              {{range .Feedbacks}}
                <div class="feedback-item">
                  <p class="feedback-header">{{.GivenBy}} <small>({{.FeedbackTime.Format "Jan 02, 2006 3:04:05 PM"}})</small></p>
                  <p>{{.Feedback}}</p>
                </div>
              {{end}}
            {{else}}
              <p>No feedback available for this exercise.</p>
            {{end}}
          </div>
        </div>
      {{end}}
    </div>
  </section>
  <script>
    // Collapsible accordion functionality
    $(document).ready(function(){
      var coll = document.getElementsByClassName("collapsible");
      for (var i = 0; i < coll.length; i++) {
        coll[i].addEventListener("click", function() {
          this.classList.toggle("active");
          var content = this.nextElementSibling;
          if (content.style.maxHeight){
            content.style.maxHeight = null;
          } else {
            content.style.maxHeight = content.scrollHeight + "px";
          } 
        });
      }
    });
  </script>
</body>
</html>
`

var SUBMISSION_VIEW_TEMPLATE = `
	<!DOCTYPE html>
	<html>
	<head>
	<title>Student Dashboard</title>
	<meta http-equiv="refresh" content="120" >
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />

	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	<style>
	.topcorner{
		position:absolute;
		top:0;
		right:0;
	}
	</style>
	</head>
	<body>
	<div class="container">
	<nav class="breadcrumb" aria-label="breadcrumbs">
		<ul>
		<li>
			<a id="view-exercise-link" href="#">
			<span class="icon is-small">
				<i class="fas fa-home" aria-hidden="true"></i>
			</span>
			<span>Exercises</span>
			</a>
		</li>
		<li>
			<a id="problem-dashboard-link" href="#">
			<span class="icon is-small">
				<i class="fas fa-book" aria-hidden="true"></i>
			</span>
				<span>Problem Dashboard ({{.ProblemName}})</span>
			</a>
		</li>
		<li class="is-active">
			<a href="#">
				<span class="icon is-small">
					<i class="fas fa-puzzle-piece" aria-hidden="true"></i>
				</span>
			<span>{{.StudentName}}'s Dashboard</span>
			</a>
		</li>
		</ul>
	</nav>
	<div class="topcorner">{{.Username}}({{.UserRole}})</div>
		<!-- <h2 class="title is-2">{{.StudentName}}'s Submissions for {{.ProblemName}}</h2> -->
		<div class="tabs">
			<ul>
				<li><a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{.ProblemID}}&uid={{.UserID}}&role={{.UserRole}}{{if ne .Password ""}}&password={{.Password}}{{end}}">Code Snapshot</a></li>
				<li><a href="/student_dashboard_feedback_provision?student_id={{.StudentID}}&problem_id={{.ProblemID}}&uid={{.UserID}}&role={{.UserRole}}{{if ne .Password ""}}&password={{.Password}}{{end}}" >Feedback</a></li>
				<li class="is-active"><a>Submissions</a></li>
			</ul>
		</div>
		<div>
			{{range .Submissions}}
			<div class="box">
			<h4 class="title is-4">Submitted at {{.SubmittedAt.Format "Jan 02, 2006 3:04:05 PM"}}</h4>
			{{if eq .Grade ""}} Not Graded {{else}} Graded {{if eq .Grade "correct"}} <span class="tag is-success">correct</span> {{else if eq .Grade "incorrect"}} <span class="tag is-danger">incorrect</span> {{else}} {{.Grade}} {{end}} {{end}}
			<div class="accordions">
				<h3>Code</h3>
				<div>
					<textarea class="editor" id="editor-{{.ID}}">{{ .Code }}</textarea>
				</div>
			</div>
			
			{{if eq .Grade ""}}
			<div class="columns">
				<div class="column is-three-quarters"><input  class="input is-info" id="{{.ID}}" type="text" placeholder="Provide your feedback!"></div>
				<div class="column"><button  class="button is-success" onclick="sendGrade({{.ID}}, {{.SnapshotID}}, 'correct')">Correct</button></div>
				<div class="column"><button  class="button is-danger" onclick="sendGrade({{.ID}}, {{.SnapshotID}}, 'incorrect')">Incorrect</button></div>
			</div>
			{{end}}
			</div>
			{{end}}
		</div>
		<script>
			$(document).ready(function(){
				$('#view-exercise-link').attr("href", "/view_exercises"+window.location.search);
				$('#problem-dashboard-link').attr("href", "/problem_dashboard"+window.location.search+"&problem_id={{.ProblemID}}");
			});
			var snapshotEditors = document.getElementsByClassName("editor");
				
			for (let i = 0; i<snapshotEditors.length; i++){
				var code = CodeMirror.fromTextArea(snapshotEditors[i], {lineNumbers: true, mode: "{{getEditorMode .ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
				code.setSize("100%", "auto");
			}
			

			function sendGrade(submission_id, snapshot_id, grade) {
				var feedback = $('#'+submission_id).val().trim();
				var code = $('#editor-'+submission_id).val();
				$.post("/teacher_grades", {content: code, changed: "", decision: grade, sid: submission_id, uid: {{.UserID}}, role: {{.UserRole}}{{if ne .Password ""}}, password: {{.Password}}{{end}}  }, function(data, status){
					if (status == "success"){
						if (feedback != "") {
							$.post("/save_snapshot_feedback", {snapshot_id: snapshot_id, feedback: feedback, uid: {{.UserID}}, role: {{.UserRole}}{{if ne .Password ""}}, password: {{.Password}}{{end}} }, function(data1, status1){
							});
						}
						alert("Graded successfully!");
						window.location.reload();
					} else {
						alert("Could not grade the submission. Please try again!");
					}
				});
			}
			$(".accordions").accordion({ header: "h3", active: false, collapsible: true });
		</script>

	</body>
	</html>
`

var AI_SETTINGS = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>AI Settings</title>
  <script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
  <style>
    body {
      background: cornflowerblue;
      font-family: "Arial", sans-serif;
      min-height: 100vh;
      display: flex;
      justify-content: center;
      align-items: center;
      margin: 0;
    }

    .dashboard-box {
      background: #fff;
      border-radius: 10px;
      padding: 30px;
      max-width: 700px;
      width: 100%;
      box-shadow: 0px 10px 30px rgba(0, 0, 0, 0.2);
      text-align: center;
    }

    .dashboard-box .title {
      font-size: 1.8rem;
      font-weight: 700;
      color: #333;
      margin-bottom: 30px;
    }

    .input, select {
      width: 100%;
      padding: 10px;
      font-size: 1rem;
      margin: 10px 0;
      border-radius: 5px;
      border: 1px solid #ccc;
    }

    .button {
      padding: 12px;
      font-size: 1rem;
      border: none;
      border-radius: 5px;
      cursor: pointer;
      width: 100%;
      margin-top: 15px;
    }

    .button-primary {
      background: cornflowerblue;
      color: white;
    }

    .button-primary:hover {
      background: linear-gradient(45deg, #2575fc, #6a11cb);
    }

    .button-danger {
      background: linear-gradient(45deg, #f56a79, #ff4757);
      color: white;
    }

    .button-danger:hover {
      background: linear-gradient(45deg, #ff4757, #f56a79);
      transform: translateY(-3px);
      box-shadow: 0 8px 16px rgba(0, 0, 0, 0.4);
    }

    .api-section {
      display: none;
      margin-top: 10px;
    }

    label {
      font-weight: bold;
      color: #444;
      display: block;
      margin-top: 10px;
    }
  </style>
</head>
<body>
  <div class="dashboard-box">
    <h1 class="title">AI Settings</h1>

    <label for="ai-select">Select AI Provider</label>
    <select id="ai-select">
      <option value="">-- Choose --</option>
      <option value="claude">Claude AI</option>
      <option value="openai">OpenAI</option>
      <option value="deepseek">DeepSeek</option>
    </select>

    <div id="claude" class="api-section">
      <label for="claude-key">Claude API Key</label>
      <input type="password" id="claude-key" class="input" placeholder="Enter Claude API Key">
    </div>

    <div id="openai" class="api-section">
      <label for="openai-key">OpenAI API Key</label>
      <input type="password" id="openai-key" class="input" placeholder="Enter OpenAI API Key">
    </div>

    <div id="deepseek" class="api-section">
      <label for="deepseek-key">DeepSeek API Key</label>
      <input type="password" id="deepseek-key" class="input" placeholder="Enter DeepSeek API Key">
    </div>

    <button class="button button-danger" id="submit-settings-btn">Submit Settings</button>
    <button class="button button-primary" id="display-prompts-btn">Display Prompts</button>
  </div>

  <script>
    $(document).ready(function () {
    $('#ai-select').on('change', function () {
      $('.api-section').hide();
      const selected = $(this).val();
      if (selected) {
        $('#' + selected).fadeIn();
      }
    });

    $('#submit-settings-btn').on('click', function () {
      const selectedAI = $('#ai-select').val();
      const apiKey = $('#' + selectedAI + '-key').val();

      if (!selectedAI || !apiKey) {
        alert("Please select an AI provider and enter its API key.");
        return;
      }

      const aiProviderIDs = {
        'claude': 1,
        'openai': 2,
        'deepseek': 3
      };

      const requestData = {
        id: aiProviderIDs[selectedAI],
        api_key: apiKey
      };

      $.ajax({
        url: '/add_api_key',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(requestData),
        success: function (response) {
          alert("API key updated successfully!");
          $('#' + selectedAI + '-key').val(''); // Clear the input
        },
        error: function (xhr) {
          alert("Error: " + xhr.responseText);
        }
      });
    });

    $('#display-prompts-btn').on('click', function () {
  window.location.href = '/prompt_view';
});

  });
  </script>
</body>
</html>
`

var PROMPT_DISPLAY = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>AI Prompt Dashboard</title>
  <script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.5/codemirror.min.css">
<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.5/codemirror.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.5/mode/javascript/javascript.min.js"></script>

  <style>
    body {
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      background: cornflowerblue;
      margin: 0;
      padding: 0;
    }

    .dashboard {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 50px 20px;
    }

    .prompt-container {
      width: 90%;
      max-width: 1100px;
      background: #ffffff;
      margin-bottom: 25px;
      border-radius: 12px;
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
      overflow: hidden;
    }

    .prompt-title {
      font-size: 20px;
      font-weight: 600;
      color: #2c3e50;
      padding: 18px 24px;
      cursor: pointer;
      background: #f0f0f0;
      transition: background 0.2s ease-in-out;
    }

    .prompt-title:hover {
      background: #e1eaff;
    }

    .prompt-content {
      display: none;
      padding: 20px 24px 30px;
    }

    pre {
      white-space: pre-wrap;
      word-wrap: break-word;
      font-size: 16px;
      line-height: 1.6;
      color: #34495e;
      background: #f8f9fa;
      padding: 15px;
      border-radius: 8px;
    }

    @media (max-width: 768px) {
      .prompt-title {
        font-size: 18px;
      }

      pre {
        font-size: 15px;
      }
    }
  </style>
</head>
<body>

  <div class="dashboard">

    <div class="prompt-container">
      <div class="prompt-title">Prompt 1: Individual AI Feedback</div>
      <div class="prompt-content">
        <textarea class="code-block">

You are assisting instructors in real-time classroom coding exercises, each lasting only 10-15 minutes. Your role is to analyze a student's current code submission, identify relevant error patterns, select appropriate scaffolding strategies, and generate multiple scaffolds clearly varying in effort and actionability. Your scaffolds must explicitly support student autonomy by prompting students to actively engage in problem-solving, reflect upon their coding decisions, and maintain responsibility for their learning outcomes. Scaffolds should facilitate productive struggle, fostering independence and deeper conceptual understanding rather than providing overly explicit or prescriptive solutions. The scaffolds you produce will be reviewed by instructors, who will select the most suitable scaffolds based on student needs and classroom context.

<exercise_description>
{Detailed description of the exercise would be placed here}
</exercise_description>

<student_code>
{The student's current code would be placed here}
</student_code>

<previous_scaffold>
{Previous scaffold content placed here, or "None" if this is the first scaffold}
</previous_scaffold>

# Programming Language
The student is coding in Python. Tailor all scaffolding approaches to Python-specific concepts, common Python errors, and Pythonic solutions.

# Time Constraint Guidance
Each scaffold MUST be implementable within the 10-15 minute exercise window. Consider:
- For Low-Effort scaffolds: Student should need <3 minutes to understand and apply
- For Moderate-Effort scaffolds: Student should need 3-7 minutes to process and implement
- For High-Effort scaffolds: Student should need 7-12 minutes to understand and integrate the concept
Avoid scaffolds requiring extensive code rewriting or introducing entirely new approaches that cannot be reasonably implemented in the remaining time.

# Scaffolding Strategies
1. **Gap-Fill Prompts:** Templates with blanks requiring student input.
2. **Incremental Hints:** Short, progressively explicit guidance.
3. **Minimal Correction Prompts:** Highlight specific errors with minimal suggested fixes.
4. **Guided Questions:** Reflective questions prompting conceptual thinking.
5. **Targeted Code Comments:** Specific, embedded comments guiding improvements without explicit solutions.
6. **Partial Worked Examples:** Concise examples aligned with the student's current approach.
7. **Reference Examples:** Adaptable examples illustrating relevant coding patterns.
8. **Error Message Translation:** Clear, student-friendly interpretations of errors.

# Effort Levels (Positive Struggle)
- **Low (Explicit):** Direct instructions; minimal cognitive challenge.
- **Moderate (Reflective):** Prompts thoughtful reflection and moderate reasoning.
- **High (Conceptual):** Requires deeper conceptual understanding and integrative thinking.

# Actionability Levels
- **High:** Clearly specifies immediate steps for student action.
- **Moderate:** Provides directional guidance, requiring some interpretation.
- **Low:** Offers conceptual insights without explicit instructions, significant interpretation required.

# Instructions
Analyze the student's submitted code and clearly describe identified error patterns. Generate multiple scaffold suggestions spanning a range of effort (low to high) and actionability (high to low). Explicitly support student autonomy, productive struggle, and continuity.

Scaffolds must incorporate the student's current code structure and approach whenever possible. Scaffolds should build upon what the student has already created rather than suggesting entirely new implementations. This ensures continuity in the student's thought process, honors their existing work, and makes scaffolds immediately actionable. Only suggest significant restructuring when the current approach contains fundamental conceptual errors that cannot be remedied through incremental improvements.

For each scaffold clearly state:
- **Strategy used**
- **Effort level (Low, Moderate, High)**
- **Actionability level (High, Moderate, Low)**
- **Content:** Scaffold content must be brief, clear, and immediately actionable—ideally presented in 1–3 concise sentences or a succinct code snippet, easily understandable within a short time frame.
- **Learning support:** Explicitly describe how the scaffold addresses misconceptions and supports student learning.
- **Implementation success metrics:** Observable indicators demonstrating the student’s successful application of scaffolds (e.g., "Student correctly implements loop structures instead of repetitive statements," "Student independently resolves index errors by accurately adjusting loop boundaries.")
- **Recommendation reasoning:** Clearly justify why this scaffold occupies its specific position in your ordered list.

**IMPORTANT:**  
Order scaffolds clearly by recommendation level. The FIRST scaffold listed should be your HIGHEST recommendation, informed by:
1. Clearly identified common errors in the student's code.
2. Previous scaffolds provided (if applicable).
3. Optimal progression in cognitive complexity and student support.
4. Feasibility of implementation within the time constraint.

        </textarea>
      </div>
    </div>

    <div class="prompt-container">
      <div class="prompt-title">Prompt 2: Class Summary</div>
      <div class="prompt-content">
        <textarea class="code-block">

You are an expert in analyzing student code submissions for introductory computer science exercises (CS1/CS2 level). Given a collection of student code submissions for a specific programming problem and their corresponding assessment results, your task is to generate a structured summary that provides instructors with actionable insights into common student difficulties.

Your analysis should follow a staged approach, where each stage builds upon insights gained in previous stages. This will ensure a systematic, consistent, and comprehensive analysis of student performance patterns.

## Staged Analysis Process

### Stage 1: Data Processing and Performance Classification
First, process the raw submission data:
1. Calculate the total number of submissions.
2. Classify each submission into one of these performance levels using these specific criteria:
   - **Poor (Failing with major issues)**:
     * Submissions that fail most or all test cases
     * Code contains 3+ fundamental errors from different error categories
     * Code structure shows significant misunderstanding of core concepts
     * Code may not compile/run at all
   - **Struggling (Failing but showing potential)**:
     * Submissions that pass some test cases but fail others
     * Code contains 1-2 fundamental errors, but shows understanding of basic concepts
     * Core algorithm demonstrates partial correctness
     * Code runs but produces incorrect results for certain inputs
   - **Good Progress (Passing with minor mistakes)**:
     * Submissions that pass most test cases
     * Code has working core functionality with minor issues
     * Errors are limited to edge cases or efficiency concerns
     * Code produces correct results for most inputs but may have minor issues
   - **Strong (Fully passing)**:
     * Submissions that pass all test cases
     * Code is efficient, well-structured, and handles all edge cases
     * Solution demonstrates thorough understanding of concepts
     * Code is well-documented and follows best practices
3. Calculate counts and percentages for each performance level.
4. Write a concise summary of overall class performance based on these categories.

### Stage 2: Error Identification and Categorization
Using the performance classifications from Stage 1, particularly focusing on the "Poor" and "Struggling" submissions:
1. Identify all errors in each submission. Start with these common error categories:
   - **Logic Errors**: Fundamental mistakes in the algorithm or program logic.
   - **Syntax Errors**: Issues with the syntax of the programming language.
   - **Data Type Errors**: Incorrect use or manipulation of data types.
   - **Control Flow Errors**: Problems with the sequence of execution (e.g., incorrect loop conditions, branching).
   - **Function/Method Errors**: Incorrect use or implementation of functions/methods.
   - **Data Structure Errors**: Incorrect use or implementation of data structures.
   - **Edge Case Handling**: Failure to handle unusual or boundary input values.
   - **Runtime Errors**: Errors that occur during the execution of the program.
   - **Algorithm Selection Errors**: Choosing inappropriate algorithms or approaches for solving problems.
   - **Variable Scope Issues**: Problems with variable accessibility (local vs. global).
   - **Off-by-One Errors**: Errors related to boundary conditions in loops or indexing.

   You may identify additional error categories if you observe patterns that don't fit well into the above categories. Include any additional categories in the all_categories_used array in Stage 6.

2. Count the frequency of each error category.
3. Select the top 5 most common error categories.
4. For each common error category:
   a. Calculate occurrence count and percentage among failing submissions.
   b. Write a clear description of the error pattern.
   c. Select a representative code example that clearly demonstrates the issue, using no more than 5-7 lines of code.
   d. Include the student_ids array listing all students who exhibited this error.
   e. Choose examples that are:
      - Clear: Demonstrate the error with minimal extraneous code
      - Typical: Represent the most common manifestation of the error
      - Concise: Short but with enough context to understand the issue
      - Educational: Show cases where the fix would be instructive

### Stage 3: Correlation and Pattern Analysis
Building on the error categorizations from Stage 2:
1. Analyze which error categories frequently co-occur in the same submissions.
2. Calculate correlation percentages for pairs of errors.
3. Identify the 3-5 strongest error correlations.
4. For each correlated pair:
   a. Report the correlation strength (percentage of failing submissions containing both errors).
   b. Report the correlation count (number of students showing both errors).
   c. Formulate a hypothesis explaining why these errors might be related.
   d. Select an example that clearly demonstrates both errors together, using no more than 8 lines of code.
   e. Include the student_ids array of students who exhibited this correlation.
   f. Choose examples where:
      - Both errors are clearly visible in proximity
      - The relationship between the errors is evident
      - The example supports your hypothesis about why they correlate

### Stage 4: Misconception Inference
Based on the error patterns and correlations identified in Stages 2 and 3:
1. Infer 3-5 potential misconceptions that might explain the observed error patterns.
2. For each misconception:
   a. Link it to specific error categories from earlier stages using the related_error_categories array.
   b. Calculate how many submissions exhibit error patterns suggesting this misconception (occurrence_count).
   c. Calculate the percentage of submissions showing this misconception (occurrence_percentage).
   d. Write a clear explanation of the likely misunderstanding.
   e. Select a representative code example that demonstrates this misconception, using no more than 5-7 lines of code.
   f. Include the student_ids array of students who exhibited this misconception.
   g. Choose examples that:
      - Clearly illustrate the conceptual misunderstanding
      - Represent common patterns across multiple students
      - Show the root cause rather than just symptoms

### Stage 5: Intervention Development
For each misconception identified in Stage 4:
1. Develop 2-3 specific instructional interventions that address the misconception.
2. For strongly correlated error pairs from Stage 3, suggest interventions that address both issues simultaneously.
3. Prioritize interventions based on the frequency and severity of associated errors.
4. For each intervention, specify:
   a. The target_misconception it addresses.
   b. An array of related_errors categories it targets.
   c. An array of instructional_strategies (2-3 specific approaches).
   d. A priority_level of "High", "Medium", or "Low".
   e. An estimated_effort of "Low", "Medium", or "High".
   f. An implementation_phase of "Immediate", "Short-term", or "Long-term".
5. Make recommendations concrete and actionable, such as:
   - Specific in-class activities
   - Targeted assignments
   - Code examples to discuss
   - Conceptual explanations to emphasize

### Stage 6: Additional Insights
After completing Stages 1-5, include:
1. Common good practices observed in successful submissions:
   a. For each good practice, provide:
      - A description of the practice.
      - Example code demonstrating this good practice.
      - A correct_implementation flag set to true.
2. Error category distribution:
   a. List all error categories used in the analysis in all_categories_used array.
   b. List any additional categories identified but not included in the main analysis in additional_categories_identified array.
3. Submission patterns:
   a. Provide an observation about when/how students submitted their work.
   b. Include a potential_impact analysis of these submission patterns.
   c. Include a submission_timeline object with dynamically determined time intervals as keys and objects containing:
      - count: number of submissions on that date
      - passing_rate: percentage of passing submissions on that date (e.g., "75%")
4. Comparative analysis:
   a. Include a success_factors array listing what successful students do differently.
   b. Include a challenge_factors array listing common challenges for struggling students.

        </textarea>
      </div>
    </div>

    <div class="prompt-container">
      <div class="prompt-title">Prompt 3: Progress Generation</div>
      <div class="prompt-content">
        <textarea class="code-block">

You are given multiple code submissions from students. Your task is to:
1. Analyze each student's submission and compare it with the correct solution.
2. Provide a one-line explanation of the student's feedback.
3. Calculate a percentage (0-100) indicating how close their submission is to the correct solution.
       </textarea>
      </div>
    </div>

  </div>

  <script>
    $(document).ready(function () {
    $('.prompt-title').on('click', function () {
      const content = $(this).next('.prompt-content');
      $('.prompt-content').not(content).slideUp();
      content.slideToggle();
    });

    // Render markdown
    $('.prompt-content').each(function () {
      const raw = $(this).text();
      $(this).html(marked.parse(raw));
    });
  });
  </script>
</body>
</html>
`

var SETTINGS_VIEW = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
  <script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
  <style>
    body {
      background: cornflowerblue;
      font-family: "Arial", sans-serif;
      min-height: 100vh;
      display: flex;
      justify-content: center;
      align-items: center;
      margin: 0;
      position: relative;
    }
    .dashboard-box {
      background: #fff;
      border-radius: 10px;
      padding: 30px;
      max-width: 900px;
      width: 100%;
      box-shadow: 0px 10px 30px rgba(0, 0, 0, 0.2);
      text-align: center;
    }
    .dashboard-box .title {
      font-size: 1.8rem;
      font-weight: 700;
      color: #333;
      margin-bottom: 30px;
    }
    .logout-container {
      position: fixed;
      top: 20px;
      right: 20px;
      width: auto;
      z-index: 10;
    }
    #logout-button {
      padding: 15px;
      background: linear-gradient(45deg, #f56a79, #ff4757);
      color: white;
      border-radius: 30px;
      text-align: center;
      font-size: 1.2rem;
      font-weight: 600;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
      transition: all 0.3s ease;
      width: 100%;
    }
    #logout-button i {
      margin-right: 10px;
    }
    #logout-button:hover {
      background: linear-gradient(45deg, #ff4757, #f56a79);
      transform: translateY(-3px);
      box-shadow: 0 8px 16px rgba(0, 0, 0, 0.4);
    }
    #logout-button:active {
      transform: translateY(1px);
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    }
    .container {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 20px;
    }
    .drawer {
      font-family: Arial, sans-serif;
      padding: 20px;
      background: white;
      border-radius: 10px;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    }
    .drawer h3 {
      font-size: 20px;
      font-weight: bold;
      margin-bottom: 20px;
      color: #333;
    }
    .switch {
      position: relative;
      display: inline-block;
      width: 34px;
      height: 20px;
    }
    .switch input {
      opacity: 0;
      width: 0;
      height: 0;
    }
    .slider {
      position: absolute;
      cursor: pointer;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background-color: #ccc;
      transition: 0.4s;
      border-radius: 50px;
    }
    .slider:before {
      position: absolute;
      content: "";
      height: 12px;
      width: 12px;
      border-radius: 50px;
      left: 4px;
      bottom: 4px;
      background-color: white;
      transition: 0.4s;
    }
    input:checked + .slider {
      background-color: #2196F3;
    }
    input:checked + .slider:before {
      transform: translateX(14px);
    }
    .input {
      width: 100%;
      padding: 10px;
      font-size: 1rem;
      margin: 10px 0;
    }
    .button {
      margin-top: 10px;
    }
    .button.is-primary, .button.is-info, .button.is-danger {
      width: 100%;
    }
    .button.is-primary {
      background: cornflowerblue;
      color: white;
      border: none;
    }
    .button.is-primary:hover {
      background: linear-gradient(45deg, #2575fc, #6a11cb);
    }
    .button.is-info {
      background: cornflowerblue;
      color: white;
      border: none;
    }
    .button.is-info:hover {
      background: linear-gradient(45deg, #00bcd4, #1e90ff);
    }
  .grid-container {
    display: grid;
    grid-template-columns: 1fr 1fr; /* Two equal columns */
    gap: 20px;
  }

  .grid-item {
    background: #fff;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  }

  /* Make it responsive */
  @media (max-width: 768px) {
    .grid-container {
      grid-template-columns: 1fr; /* Single column on small screens */
    }
  }
  </style>
</head>
<body>
  <div class="dashboard-box">
  <h1 class="title">Welcome {{.Username}}</h1>
  
  <div class="grid-container">
    
    <!-- Select Course -->
    <div class="grid-item">
      <h4>Select Course</h4>
      <div class="select">
        <select id="course-dropdown">
          <option value="">Loading courses...</option>
        </select>
      </div>
      <button class="button is-primary" id="view-exercises-btn" style="margin-top: 89px;">
    View Exercises
  </button>
    </div>

    <!-- Add TA to Course -->
    <div class="grid-item">
      <h4>Add TA to Course</h4>
      <input type="text" id="teacher-name" class="input" placeholder="Enter Teacher Name">
      <input type="text" id="teacher-pass" class="input" placeholder="Enter Teacher Password">
      <button class="button is-primary" id="add-teacher-btn">Add Teacher</button>
    </div>

    <!-- Add Students to Course -->
    <div class="grid-item">
      <h4>Add Students to Course</h4>
      <input id="student-names" class="input" placeholder="Enter Comma Separated Names">
      <button class="button is-primary" id="add-students-btn">Add Students</button>
    </div>

    <!-- Add Course -->
    <div class="grid-item">
      <h4>Add Course</h4>
      <input type="text" id="course-id" class="input" placeholder="Enter Course ID">
      <p class="help is-info">Example: <strong>S2025_COMP7712_01</strong></p>
      <button class="button is-primary" id="add-course-btn">Add Course</button>
    </div>

    <!-- Peer Tutoring -->
    

	<div class="grid-item">
  <button class="button is-primary" id="ai-settings-btn">View AI Settings</button>
</div>

    <!-- Logout -->
    <div class="grid-item">
      <button class="button is-danger" id="logout-button">
        <i class="fas fa-sign-out-alt"></i> Logout
      </button>
    </div>

  </div>
</div>

  <script>
$('#add-course-btn').click(function () {
    let courseID = $('#course-id').val().trim();
	const urlParams = new URLSearchParams(window.location.search);
    const teacherId = urlParams.get('teacher_id');

    if (courseID === "") {
        alert("Please enter Course ID.");
        return;
    }

    if (!teacherId) {
        alert("Teacher ID is missing.");
        return;
    }

    $.ajax({
        url: "/add_course",
        type: "POST",
        contentType: "application/json",
        dataType: "json",
        data: JSON.stringify({ course_id: courseID, teacher_id: teacherId }), // Include teacher_id
        success: function (response) {
            if (response && response.message) {
                alert(response.message); // Show success message
                $('#course-id').val(''); // Clear input field
				location.reload();
            } else {
                alert("Unexpected response format.");
            }
        },
        error: function (xhr) {
            alert("Error: " + xhr.responseText);
        }
    });
});

$(document).ready(function () {
    const urlParams = new URLSearchParams(window.location.search);
    const teacherId = urlParams.get('teacher_id');

    if (!teacherId) {
        alert("Teacher ID missing, redirecting to login.");
        window.location.replace('/');
        return;
    }

    // Fetch courses for the teacher
    $.get('/get_courses?teacher_id=' + teacherId, function (courses) {
        let dropdown = $('#course-dropdown');
        dropdown.empty();
        dropdown.append('<option value="">Select a course</option>');
        courses.forEach(course => {
            dropdown.append('<option value="' + course.CourseID + '">' + course.CourseID + '</option>');
        });
    }).fail(function () {
        alert("Failed to load courses.");
    });

    // Handle course selection (change URL without navigation)
    $('#course-dropdown').change(function () {
        let courseID = $(this).val();
        if (courseID) {
            // Change the URL without navigating
            const currentUrl = new URL(window.location.href);
            currentUrl.searchParams.set('course_id', courseID);
            history.pushState({}, '', currentUrl); // This updates the URL in the browser
        }
    });

    // Handle "View Exercises" button click (trigger navigation)
    $('#view-exercises-btn').click(function () {
        const courseID = $('#course-dropdown').val();
        if (courseID) {
            window.location.replace('/view_exercises?role=teacher&uid=' + teacherId + '&course_id=' + courseID);
        } else {
            alert("Please select a course.");
        }
    });
});

    $(document).ready(function () {
      // Logout functionality
      $('#logout-button').click(function () {
        if (confirm("Are you sure you want to logout?")) {
          window.location.href = "/logout";
        }
      });

      // Add TA to Course
      $('#add-teacher-btn').click(function(){
        let teacherName = $('#teacher-name').val().trim();
        let teacherPass = $('#teacher-pass').val().trim();
        let urlParams = new URLSearchParams(window.location.search);
        let courseId = urlParams.get("course_id");

        if (teacherName === "" || teacherPass === "") {
            alert("Please enter Teacher Name and Password");
            return;
        }

		if ( !courseId){
alert("Please Select a course");
return;
}

        $.ajax({
            url: "/add_teacher",
            type: "POST",
            contentType: "application/json",
            dataType: "json",
            data: JSON.stringify({ 
                teacher_name: teacherName, 
                teacher_pass: teacherPass, 
                course_id: courseId
            }),
            success: function(response) {
                alert(response.message || "Teacher added successfully");
                $('#teacher-name').val('');
                $('#teacher-pass').val('');
            },
            error: function(xhr, status, error) {
                alert("Failed to add teacher: " + xhr.responseText);
            }
        });
      });

      // Add Student to Course
      $('#add-students-btn').click(function() {
        let studentNames = $('#student-names').val().trim();
        let urlParams = new URLSearchParams(window.location.search);
        let courseID = urlParams.get("course_id");

        if (studentNames === "" ) {
            alert("Please enter student names and ensure Course ID is available.");
            return;
        }
		if ( !courseID){
alert("Please Select a course");
return;
}

        let studentsArray = studentNames.split(',').map(name => name.trim()).filter(name => name !== "");

        if (studentsArray.length === 0) {
            alert("Please enter valid student names.");
            return;
        }

        $.ajax({
            url: "/add_students",
            type: "POST",
            contentType: "application/json",
            dataType: "json",
            data: JSON.stringify({ 
                student_names: studentsArray, 
                course_id: courseID 
            }),
            success: function(response) {
                alert(response.message || "Students added successfully");
                $('#student-names').val(''); // Clear input field
            },
            error: function(xhr, status, error) {
                alert("Failed to add students: " + xhr.responseText);
            }
        });
      });

      // Toggle Peer Tutoring functionality
      $('#peer_tutoring_button').change(function () {
        if ($(this).prop('checked')) {
          alert("Peer Tutoring Enabled.");
        } else {
          alert("Peer Tutoring Disabled.");
        }
      });
    });
$('#ai-settings-btn').on('click', function () {
  window.location.href = '/ai_settings_view'; 
});
  </script>
</body>
</html>
`
var TEACHER_LOGIN = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Login</title>
    <link
      href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css"
      rel="stylesheet"
    />
    <link
      href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css"
      rel="stylesheet"
    />
    <style>
      body {
        background: cornflowerblue;
        font-family: "Arial", sans-serif;
        min-height: 100vh;
        display: flex;
        justify-content: center;
        align-items: center;
        margin: 0;
      }
      .login-box {
        background: #fff;
        border-radius: 10px;
        padding: 30px;
        max-width: 400px;
        width: 100%;
        box-shadow: 0px 10px 30px rgba(0, 0, 0, 0.2);
        text-align: center;
      }
      .login-box .title {
        font-size: 1.8rem;
        font-weight: 700;
        color: #333;
        margin-bottom: 20px;
      }
      .login-box .subtitle {
        font-size: 1rem;
        color: #777;
        margin-bottom: 30px;
      }
      .login-box .field {
        margin-bottom: 20px;
      }
      .login-box input {
        border-radius: 25px;
        padding: 10px 20px;
      }
      .login-box .button {
        border-radius: 25px;
        padding: 10px 20px;
        background: #6a11cb;
        color: white;
        font-weight: bold;
        transition: all 0.3s;
      }
      .login-box .button:hover {
        background: #2575fc;
        transform: translateY(-2px);
      }
    </style>
    <script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
  </head>
  <body>
    <div class="login-box">
      <h1 class="title">Welcome Back!</h1>
      <p class="subtitle">Sign in to manage your system</p>
      <div class="field">
        <label class="checkbox">
          <input type="checkbox" id="admin-toggle" /> Login as Admin
        </label>
      </div>

      <!-- Teacher Login Fields -->
      <div id="teacher-login">
        <div class="field">
          <p class="control has-icons-left">
            <input id="name" class="input is-medium" type="email" placeholder="Enter your username" />
            <span class="icon is-left"><i class="fas fa-user"></i></span>
          </p>
        </div>
        <div class="field">
          <p class="control has-icons-left">
            <input id="password" class="input is-medium" type="password" placeholder="Enter your password" />
            <span class="icon is-left"><i class="fas fa-lock"></i></span>
          </p>
        </div>
      </div>

      <!-- Admin Login Fields -->
      <div id="admin-login" style="display: none;">
        <div class="field">
          <p class="control has-icons-left">
            <input id="admin-username" class="input is-medium" type="text" placeholder="Enter your username" />
            <span class="icon is-left"><i class="fas fa-user"></i></span>
          </p>
        </div>
        <div class="field">
          <p class="control has-icons-left">
            <input id="admin-password" class="input is-medium" type="password" placeholder="Enter your password" />
            <span class="icon is-left"><i class="fas fa-lock"></i></span>
          </p>
        </div>
      </div>

      <button id="login" class="button is-medium">Login</button>
    </div>

    <script>
      $(document).ready(function () {
        // Toggle between Teacher and Admin login
        $('#admin-toggle').change(function () {
          if (this.checked) {
            $('#teacher-login').hide();
            $('#admin-login').show();
          } else {
            $('#teacher-login').show();
            $('#admin-login').hide();
          }
        });

        $('#login').click(function () {
          if ($('#admin-toggle').is(':checked')) {
            var adminUsername = $('#admin-username').val().trim();
            var adminPassword = $('#admin-password').val().trim();
            if (adminUsername == '' || adminPassword == '') {
              alert('Please enter username and password!');
            } else {
              $.post('/admin_signin', { username: adminUsername, password: adminPassword })
                .done(function (data) {
                  window.location.replace('/admin_dashboard');
                })
                .fail(function (xhr) {
                  alert("Login failed. Please try again.");
                });
            }
          } else {
            var name = $('#name').val().trim();
            var pass = $('#password').val().trim();
			var courseID = "F2025_COMP7712_02";
            if (name == '' || pass == '') {
              alert('Please enter email and password!');
            } else {
              $.post('/teacher_signin_complete', { username: name, password: pass })
                .done(function (data) {
					window.location.href = "/settings_view?teacher_id=" + data;
                })
                .fail(function (xhr) {
                  alert("Login failed. Please try again.");
                });
            }
          }
        });
      });
    </script>
  </body>
</html>
`

var ADMIN_DASHBOARD = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Admin Dashboard</title>
    <link
      href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css"
      rel="stylesheet"
    />
    <link
      href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css"
      rel="stylesheet"
    />
    <script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
    <style>
      body {
        background: cornflowerblue;
        font-family: "Arial", sans-serif;
        min-height: 100vh;
        display: flex;
        justify-content: center;
        align-items: center;
        margin: 0;
        position: relative;
      }
      .dashboard-box {
        background: #fff;
        border-radius: 10px;
        padding: 30px;
        max-width: 400px;
        width: 100%;
        box-shadow: 0px 10px 30px rgba(0, 0, 0, 0.2);
        text-align: center;
      }
      .dashboard-box .title {
        font-size: 1.8rem;
        font-weight: 700;
        color: #333;
        margin-bottom: 20px;
      }
      .logout-container {
        position: absolute;
        top: 20px;
        right: 20px;
        width: auto;
        z-index: 10;
      }
      #logout-button {
        padding: 15px;
        background: linear-gradient(45deg, #e74c3c, #c0392b);
        color: white;
        border-radius: 30px;
        text-align: center;
        font-size: 1.2rem;
        font-weight: 600;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
        transition: all 0.3s ease;
        width: 100%;
      }
      #logout-button i {
        margin-right: 10px;
      }
      #logout-button:hover {
        background: linear-gradient(45deg, #c0392b, #e74c3c);
        transform: translateY(-3px);
        box-shadow: 0 8px 16px rgba(0, 0, 0, 0.4);
      }
      #logout-button:active {
        transform: translateY(1px);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
      }
      .add-teacher-container {
        margin-top: 30px;
      }
      .add-teacher-container input {
        margin-bottom: 10px;
      }
.button.is-primary, .button.is-info, .button.is-danger {
      width: 100%;
    }
    .button.is-primary {
      background: cornflowerblue;
      color: white;
      border: none;
    }
    .button.is-primary:hover {
      background: linear-gradient(45deg, #2575fc, #6a11cb);
    }
    </style>
  </head>
  <body>
    <div class="dashboard-box">
      <h1 class="title">Add Teacher</h1>

      <!-- Add Teacher Form -->
      <div class="add-teacher-container">
        <input type="text" id="teacher-name" class="input" placeholder="Enter Teacher Name">
        <input type="text" id="teacher-pass" class="input" placeholder="Enter Teacher Password">
        <button class="button is-primary" id="add-teacher-btn">Add Teacher</button>
      </div>

      <!-- Logout Button -->
      <div class="logout-container">
        <button id="logout-button">
          <i class="fas fa-sign-out-alt" style="margin-right: 10px;"></i> Logout
        </button>
      </div>
    </div>

    <script>
      $(document).ready(function () {
        // Logout button functionality
        $('#logout-button').click(function () {
          if (confirm("Are you sure you want to logout?")) {
            window.location.href = "/logout";
          }
        });

        // Add Teacher button functionality
        $('#add-teacher-btn').click(function () {
          let teacherName = $('#teacher-name').val().trim();
          let teacherPass = $('#teacher-pass').val().trim();

          if (teacherName === "" || teacherPass === "") {
            alert("Please enter Teacher Name and Password.");
            return;
          }

          $.ajax({
            url: "/add_teacher",
            type: "POST",
            contentType: "application/json",
            dataType: "json",
            data: JSON.stringify({ 
              teacher_name: teacherName, 
              teacher_pass: teacherPass
            }),
            success: function (response) {
              if (response && response.message) {
                alert(response.message);
              } else {
                alert("Teacher added successfully, but response is missing data.");
              }

              // Clear input fields after successful addition
              $('#teacher-name').val('');
              $('#teacher-pass').val('');
            },
            error: function (xhr, status, error) {
              alert("Failed to add teacher: " + xhr.responseText);
            }
          });
        });
      });
    </script>
  </body>
</html>
`

var PROBLEM_FILE_UPLOAD_VIEW = `
<!DOCTYPE html>
<html lang="en">
   <head>
      <meta charset="utf-8">
      <meta name="viewport" content="width=device-width, initial-scale=1">
      <title>Broadcast Problem</title>
      <script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />

	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	  <script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<style>
		.menu {
			padding: 10px;
			padding-left: 100px;
		}
		.content {
			padding-top: 115px;
		}
		.topcorner{
			position:absolute;
			top:0;
			right:0;
		}
	</style>
	</head>
   <body>
   <div class="container">
   <nav class="navbar is-fixed-top breadcrumb menu" role="navigation" aria-label="breadcrumbs">
   <ul>
		<li>
			<a id="view-exercise-link" href="#">
			<span class="icon is-small">
				<i class="fas fa-home" aria-hidden="true"></i>
			</span>
			<span>Exercises</span>
			</a>
		</li>
		<li class="is-active">
			<a href="#">
			<span class="icon is-small">
				<i class="fas fa-book" aria-hidden="true"></i>
			</span>
			<span>Add a new exercise</span>
			</a>
		</li>
		</ul>
	</nav>
	<div class="content">
		<div id="problem" class="file is-centered is-boxed is-success has-name">
				<label class="file-label">
					<input class="file-input" type="file" name="resume">
					<span class="file-cta">
					<span class="file-icon">
						<i class="fas fa-upload"></i>
					</span>
					<span class="file-label">
						Select Exercise File
					</span>
					</span>
				</label>
			</div>
		<div style="visibility:hidden;" id="editor-area">
			<article class="message">
					<div class="message-header">
						<p><span id="filename"></span></p>
					</div>
					<div class="message-body">
						<div>
							<textarea id="editor"></textarea>
						</div>
					</div>
				</article>
		</div>
		
		<div id="answer" class="file is-centered is-info has-name">
			<label class="file-label">
			  <input class="file-input" type="file" name="resume">
			  <span class="file-cta">
				<span class="file-icon">
				  <i class="fas fa-upload"></i>
				</span>
				<span class="file-label">
				Select Answer File (if any)
				</span>
			  </span>
			  <span id="answer_filename" class="file-name">
				No file selected
			  </span>
			</label>
		  </div>
		<button style="visibility:hidden" id="submit" class="button is-success is-rounded">
				<span class="icon is-small">
					<i class="fas fa-check"></i>
				</span>
				<span>Broadcast</span>
			</button>
			<input type="hidden" id="points" value="">
			<input type="hidden" id="effort" value="">
			<input type="hidden" id="attempt" value="">
			<input type="hidden" id="tag" value="">
			<input type="hidden" id="exact_answer" value="">
		</div>
	</div>
	<script>
	$(document).ready(function(){
		$('#view-exercise-link').attr("href", "/view_exercises"+window.location.search);
	  });
	document.querySelector('#problem input[type=file]').onchange = function(){
		document.querySelector('#problem').style.visibility = "hidden";
		var file = this.files[0];
		document.querySelector('#filename').textContent = file.name;
		var reader = new FileReader();
		reader.onload = function(progressEvent){
	  
		  // By lines
		  var lines = this.result.split('\n');
		  var firstLine = lines[0];
		  lines.splice(0, 1);
		  var content = lines.join('\n');
		  if (firstLine.length == 0 || (firstLine[0]!='#' && !firstLine.startsWith('//') )){
			  alert("Invalid problem header!");
			return;
		  }
		  var prefix = '';
		  if (firstLine[0] == '#') {
			  prefix = '#';
			firstLine.replace("#", '');
		  } else {
			  prefix = '//';
			firstLine.replace("//", '');
		  }
		  params = get_problem_info(firstLine);
		//   alert(params[1]+" Points, "+params[2]+" for effort. Maximum attempts: " + params[3]);
		 $('#editor-area').css('visibility', 'visible');
		  document.querySelector('#editor').textContent = prefix + ' ' + params[1]+" Points, "+params[2]+" for effort. Maximum attempts: " + params[3] + "\n" + content;
		  var editor = document.getElementById("editor");
		  var myCodeMirror = CodeMirror.fromTextArea(editor, {lineNumbers: true, mode: get_editor_mode(file.name), theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: "nocursor"});
		  myCodeMirror.setSize("100%", 400)
		  $('#submit').css('visibility', 'visible');
		  $('#points').val(params[1]);
		  $('#effort').val(params[2]);
		  $('#attempt').val(params[3]);
		  $('#tag').val(params[4]);
		};
		reader.readAsText(file);
	  };
	  
	  document.querySelector('#answer input[type=file]').onchange = function(){
		// document.querySelector('#answer').style.visibility = "hidden";
		var file = this.files[0];
		document.querySelector('#answer_filename').textContent = file.name;
		var reader = new FileReader();
		reader.onload = function(progressEvent){
	  
		  $('#exact_answer').val(this.result);
		};
		reader.readAsText(file);
	  };

	  function get_problem_info(content) {
		  let regexpNames =  /\s*(\d+)\s+(\d+)\s+(\d+)(?:\s+(\w.*))?/mg;
		let match = regexpNames.exec(content);
		return match;
	  }
	  function get_editor_mode(filename) {
		filename = filename.toLowerCase();
		if (filename.endsWith('.py')) {
			return "python";
		}
		if (filename.endsWith('.java')) {
			return "text/x-java";
		}
		if (filename.endsWith('.cpp') || filename.endsWith('.c++') || filename.endsWith('.c')) {
			return "text/x-c++src";
		}
		return "text";
	  }
	$(document).ready(function() {
		$.ajaxSetup({
			xhrFields: {
			  withCredentials: true
			}
		});
		$('#submit').click(function() {
    var editor = document.querySelector('.CodeMirror').CodeMirror;
    var uid = new URLSearchParams(window.location.search).get('uid');
    var points = $('#points').val();
    var effort = $('#effort').val();
    var attempt = $('#attempt').val();
    var tag = $('#tag').val();
    var filename = $('#filename').text();
    var answer = $('#exact_answer').val().trim();

    $.post("/teacher_broadcasts", {
        role: "teacher", 
        uid: uid, 
        content: editor.getValue(), 
        answer: answer, 
        merit: points, 
        effort: effort, 
        attempts: attempt, 
        tag: tag, 
        filename: filename, 
        exact_answer: "True"
    })
    .done(function(data, status, xhr) {
        // Check if the request was successful (HTTP status 200)
        if (xhr.status == 200) {
            alert("Exercise broadcasted successfully!");
            window.location.replace("/view_exercises?role=teacher&uid=" + uid);
        }
    })
    .fail(function(xhr, status, error) {
        // Handle failure, display the error message from the server
        if (xhr.status === 400) {
            alert(xhr.responseText); // This will show the error message sent by the server
        } else {
            alert("Failed to broadcast. Try again!");
        }
    });
});

	});
	
	</script>
   </body>
</html>
`
var CODE_SNAPSHOT_TAB_TEMPLATE = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<title>Student Dashboard</title>
	<script src="https://kit.fontawesome.com/923539b4ee.js" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css" integrity="sha512-IgmDkwzs96t4SrChW29No3NXBIBv8baW490zk5aXvhCD8vuZM3yUSkbyTBcXohkySecyzIrUwiF/qV0cuPcL3Q==" crossorigin="anonymous" referrerpolicy="no-referrer" />

	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.js" integrity="sha512-hGVnilhYD74EGnPbzyvje74/Urjrg5LSNGx0ARG1Ucqyiaz+lFvtsXk/1jCwT9/giXP0qoXSlVDjxNxjLvmqAw==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/python/python.min.js" integrity="sha512-/mavDpedrvPG/0Grj2Ughxte/fsm42ZmZWWpHz1jCbzd5ECv8CB7PomGtw0NAnhHmE/lkDFkRMupjoohbKNA1Q==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/mode/clike/clike.min.js" integrity="sha512-GAled7oA9WlRkBaUQlUEgxm37hf43V2KEMaEiWlvBO/ueP2BLvBLKN5tIJu4VZOTwo6Z4XvrojYngoN9dJw2ug==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/codemirror.min.css" integrity="sha512-6sALqOPMrNSc+1p5xOhPwGIzs6kIlST+9oGWlI4Wwcbj1saaX9J3uzO3Vub016dmHV7hM+bMi/rfXLiF5DNIZg==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.62.3/theme/monokai.min.css" integrity="sha512-R6PH4vSzF2Yxjdvb2p2FA06yWul+U0PDDav4b/od/oXf9Iw37zl10plvwOXelrjV2Ai7Eo3vyHeyFUjhXdBCVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<script src="https://code.jquery.com/jquery-3.6.0.min.js" integrity="sha256-/xUj+3OJU5yExlq6GSYGSHk7tPXikynS7ogEvDej/m4=" crossorigin="anonymous"></script>
	<script src="https://code.jquery.com/ui/1.12.1/jquery-ui.min.js" integrity="sha256-VazP97ZCwtekAsvgPBSUwPFKdrwD3unUfSGVYrahUqU=" crossorigin="anonymous"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.13/theme/abcdef.min.css">
	<script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.13/mode/markdown/markdown.min.js"></script>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.13/theme/eclipse.min.css">
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.13/theme/twilight.min.css">

	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	<script src="https://cdn.jsdelivr.net/npm/@creativebulma/bulma-collapsible"></script>
	<style>
		#code-snapshot {
			padding: 20px;
			margin: 30px;
			padding-bottom: 0px;
			border-radius: 25px;
		}
		#ask-for-help {
			background: #4a4a4a;
			padding: 20px;
			margin: 30px;
			padding-bottom: 0px;
			border-radius: 25px;
		}
		#submission {
			background: #4a4a4a;;
			padding: 20px;
			margin: 30px;
			padding-bottom: 0px;
			border-radius: 25px;
		}
.switch {
        position: relative;
        display: inline-block;
        width: 34px;
        height: 20px;
    }

    .switch input {
        opacity: 0;
        width: 0;
        height: 0;
    }

    .slider {
        position: absolute;
        cursor: pointer;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: #ccc;
        transition: 0.4s;
        border-radius: 20px;
    }
	#custom-prompt-response {
    display: none; /* Initially hidden */
	 margin-top: 10px; 
}

    .slider::before {
        position: absolute;
        content: "";
        height: 14px;
        width: 14px;
        left: 3px;
        bottom: 3px;
        background-color: white;
        transition: 0.4s;
        border-radius: 50%;
    }

    input:checked + .slider {
        background-color: #4CAF50;
    }

    input:checked + .slider::before {
        transform: translateX(14px);
    }
		.wrapper {
			display: grid;
			grid-template-columns: repeat(2, 1fr);
			gap: 20px;
			padding: 10px;
		}
		.status {
			display: flex;
			justify-content: space-between;
		}
		input[type="radio"] {
			margin-right: 5px;
		}
		.menu {
			padding: 10px;
			padding-left: 100px;
			padding-right: 100px;
		}
		.show {
			top: 6%;
			position: fixed;
			z-index: 200;
			background: white;
		}
		.content {
			padding-top: 7%;
		}
		.actions {
			// float: right;
			margin-left: 78%
		}
		.sub-actions {
			margin-left: 73%;
		}
		.topcorner{
			position:absolute;
			top:0;
			right:0;
		}
/* Dropdown Styling */
#scaffolding-dropdown {
    padding: 10px;
    font-size: 16px;
    border-radius: 5px;
    border: 2px solid #007bff;
    background-color: #ffffff !important; /* Ensures White Background */
    color: #333;
    width: 280px;
    cursor: pointer;
    transition: all 0.3s ease-in-out;
    appearance: none; /* Removes default arrow */
}

/* Hover Effect */
#scaffolding-dropdown:hover {
    border-color: #0056b3;
}

/* Focus Effect */
#scaffolding-dropdown:focus {
    outline: none;
    border-color: #0056b3;
    box-shadow: 0px 0px 8px rgba(0, 123, 255, 0.5);
}

/* Optgroup Styling (Category Headings) */
#scaffolding-dropdown optgroup {
    font-weight: bold;
    color: #007bff !important; /* Blue Text */
    font-size: 18px;
    padding: 5px;
    text-transform: uppercase;
    font-family: Arial, sans-serif;
}

/* Individual Option Styling */
#scaffolding-dropdown option {
    padding: 10px;
    color: #333;
    background-color: #ffffff;
}

/* Workaround: Make optgroup labels more visible */
#scaffolding-dropdown optgroup::before {
    content: "⬤ "; /* Add bullet point */
    font-size: 14px;
    color: #0056b3;
}
	</style>
	</head>
	<body>
	<div class="container">
	<nav class="navbar is-fixed-top breadcrumb menu" role="navigation" aria-label="breadcrumbs">
	<div class="navbar-start"> 
	<ul>
	  <li>
		<a id="view-exercise-link" href="#">
		  <span class="icon is-small">
			<i class="fas fa-home" aria-hidden="true"></i>
		  </span>
		  <span>Exercises</span>
		</a>
	  </li>
	  <li>
		<a id="problem-dashboard-link" href="#">
		<span class="icon is-small">
			<i class="fas fa-book" aria-hidden="true"></i>
		  </span>
			<span>Exercise Dashboard</span>
		</a>
	   </li>
	  <li class="is-active">
		<a href="#">
			<span class="icon is-small">
				<i class="fas fa-puzzle-piece" aria-hidden="true"></i>
			</span>
		  <span>{{ .Feedback.StudentName}}'s Dashboard</span>
		</a>
	  </li>
	</ul>
	</div>
	<div class="navbar-end"> 
		<div class="navbar-item"> <a href="#">{{.Username}} ({{.UserRole}})</a> </div>
	</div>

	</nav>
	<!--
		<nav class="breadcrumb is-right" aria-label="breadcrumbs">
			<ul>
			<li class="is-active"><a href="#">{{.Username}}({{.UserRole}})</a></li>
			</ul>
		</nav>
	-->
	<div class="content">
	<div class="column show" style="width: 93%;">
		<!--
		<div class="row">
			<h3 class="title is-2" style="margin-bottom: 0px;">{{ .Feedback.StudentName}}'s Dashboard for {{ .Feedback.ProblemName}}</h3>
		</div>
		-->
		<div class="row status">
			<span>Status: <strong>{{ .Status.CodingStat }} </strong></span>
			<span>Help Status: <strong>{{ .Status.HelpStat }} </strong></span>
			<span>Submission: <strong> {{ .Status.SubmissionStat }} </strong></span>
			<span>Progress: <strong>{{ .Status.Percentage }}%</strong></span>
		</div>

		<div class="tabs">
			<ul>
				<li class="is-active"><a>CodeSpace</a></li>
				<li><a href="/student_dashboard_feedback_provision?student_id={{.Feedback.StudentID}}&problem_id={{.Feedback.ProblemID}}&uid={{.Feedback.UserID}}&role={{.Feedback.UserRole}}{{if ne .Feedback.Password ""}}&password={{.Feedback.Password}}{{end}}">Feedback History</a></li>
			</ul>
		</div>
<div class="feedback-box">
			<p><strong>Brief AI Assessment:</strong> {{ .Status.Explanation }}</p>
		</div>
	</div>

	<div class="content">
		<div id="code-snapshot">
			<div class="box" style="padding: 0px; margin-bottom: 3.5rem; border: 5px solid; border-radius: 10px;">
				<div class="message-header">
					<div class="column is-two-thirds">
						<p>Latest Code Snapshot at {{.Feedback.LastSnapshot.LastUpdated.Format "Jan 02, 2006 3:04:05 PM"}}</p>
					</div>
{{if eq .UserRole "teacher"}}
{{end}}
{{if eq .UserRole "teacher"}} 
<button class="button is-info" id="snapshot-send-feedback" onclick="sendSnapshotFeedback({{ .Feedback.LastSnapshot.Code }}, {{ .Feedback.UserID }})" style="margin-top:3px; margin-bottom: 3px; color: #ffff;" >Send Inline Feedback</button>
<button class="button is-info chatgpt-feedback" id="submit-custom-prompt" 
            onclick="sendCustomPromptFeedback()" 
            style="color: #ffffff;">
            Get Feedback from AI
        </button>
	{{end}}
				</div>	
				<div id="feedback-block-99999"></div>
				<div style="background: #4a4a4a;">
					<textarea id="snapshot-editor"> {{ .Feedback.LastSnapshot.Code }} </textarea>
					<div class="actions">
    <button class="button is-info" id="snapshot-check-feedback" onclick="codeSnapshotFeedback({{ .Feedback.LastSnapshot.Code }}, {{ .Feedback.UserID }})" style="margin-top:3px; margin-bottom: 3px; color: #ffff;" >Check My Feedback</button>
</div>
					
		</div>

			</div>

<div id="code-snapshot-feedback-block"></div>
{{if ne .UserRole "student"}}
<div id="custom-prompt-box" class="box" style="background: #4a4a4a; margin-top: 10px; padding: 10px;">
    <div style="display: flex; align-items: center; justify-content: space-between;">
        <h3 id="feedback-heading" style="color: #ffffff; margin: 0;">AI & Scaffolding Feedback</h3>

        <!-- Dropdown for scaffolding selection -->
        <select id="scaffolding-dropdown" style="margin-right: 10px; padding: 5px;">
            <option value="">Loading scaffolding...</option>
        </select>


        <!-- New button to send scaffolding feedback -->
        <button class="button is-info" id="send-scaffolding-feedback" 
            onclick="sendScaffoldingFeedback()" 
            style="color: #ffffff; margin-left: 10px;">
            Send Feedback
        </button>
    </div>

    <!-- Hidden initially, shown only after response is received -->
    <textarea id="custom-prompt-response"></textarea>
</div>
{{end}}

		{{ if .Feedback.Messages}}
		<h3 style="color: black;">Student's help requests:</h3>
		<div id="ask-for-help">
			<section class="section" style="padding: 0px;">
				{{range $index, $el := .Feedback.Messages}}
					{{ if eq $el.Type 0 }} <!-- Don't show regular snapshots in this block -->
					{{ if eq (len .Feedbacks) 0 }}
						<div class="box" style="padding: 0px; margin-bottom: 3.5rem; border: 5px solid; border-radius: 10px;">
							<div class="message-header">
								<div class="column is-two-thirds">
									<p>{{if eq .Type 0}}{{.Name}} asked for help{{else if eq .Event "at_submission"}} Submission Snapshot taken {{else}} Regular Snapshot taken {{end}} at {{.GivenAt.Format "Jan 02, 2006 3:04:05 PM"}}</p>
									{{ if not (eq (len .Feedbacks) 0) }}
										<span class="tag is-success">Responded ({{ len .Feedbacks}})</span>
									{{ end }}
								</div>
							</div>

							<div class="message-body">
								<h3>Student says: </h3> {{.Message}}
								<div id="feedback-block-{{ $index }}"></div>
							</div>
								
							<div style="background: #4a4a4a;">
									<div>
										<textarea class="feedback-editor" id="feedback-editor-{{ $index }}">{{ .Code }}</textarea>
									</div>

									<div class="actions">
										<button class="button is-info help-check" id="help-check-feedback-{{ $index }}" onclick="messageFeedback( {{ $index }} ,{{ .Code }} , {{ .ID }})" style="margin-top:3px; margin-bottom: 3px; color: #000000;" >Check My Feedback</button>
										<button class="button is-info help-send" id="help-send-feedback-{{ $index }}" onclick="sendMessageFeedback( {{ $index }} ,{{ .Code }} , {{ .ID }})" style="margin-top:3px; margin-bottom: 3px; color: #000000;" >Send Feedback</button>
										
									</div>
							</div>
						</div>
					{{ end }}
					{{ end }}
					
				{{end}}
			</section>
		</div>
		{{ end }}
		{{ if ne .UserRole "student"}}
		{{ if .Submission.Submissions}}
		<h3 style="color: black;">Student's submissions: </h3> 
		<div id="submission">
			{{range $index, $el := .Submission.Submissions}}
				{{ if eq .Grade "" }}
				<div class="box" style="padding: 0px; margin-bottom: 3.5rem; border: 5px solid; border-radius: 10px;">
				
					<div class="message-header">
						<div class="column is-two-thirds">
							<p>Submitted at {{.SubmittedAt.Format "Jan 02, 2006 3:04:05 PM"}}</p>
							{{if eq .Grade ""}} Not Graded {{else}} Graded {{if eq .Grade "correct"}} <span class="tag is-success">correct</span> {{else if eq .Grade "incorrect"}} <span class="tag is-danger">incorrect</span> {{else}} {{.Grade}} {{end}} {{end}}
						</div>

						<div class="column buttons" style="padding-left: 1px; display: contents;">
							{{if eq .Grade ""}}
								<button class="button"><label><input type="radio" name="grade-{{$index}}"  value="correct" onchange="setGrade( {{ $index }}, 'correct')" />Correct </label></button>
								<button class="button"><label><input type="radio" name="grade-{{$index}}"  value="incorrect" onchange="setGrade( {{ $index }}, 'incorrect')" />Incorrect </label></button>
								<button class="button"><label><input type="radio" name="grade-{{$index}}" value="0" checked onchange="removeGrade( {{ $index }})" />Not Graded</label></button>
							{{end}}
						</div>
					</div>
					<!--
					<div class="message-body">
						<div class="columns">
							{{if eq .Grade ""}} Not Graded {{else}} Graded {{if eq .Grade "correct"}} <span class="tag is-success">correct</span> {{else if eq .Grade "incorrect"}} <span class="tag is-danger">incorrect</span> {{else}} {{.Grade}} {{end}} {{end}}
						</div>
					</div>
					-->
					
					<div style="background: #4a4a4a;">
						<div>
							<textarea class="submission-editor" id="editor-{{.ID}}">{{ .Code }}</textarea>
						</div>
						<div class="sub-actions">
							<button class="button is-info sub-check" id="sub-check-{{ $index }}" onclick="checkSubFeedback( {{ $index }}, {{.ID}}, {{.SnapshotID}},{{ .Code }})" style="margin-top:3px; margin-bottom: 3px; color: #ffff;">Check My Feedback</button>
							<button class="button is-info sub-submit" id="sub-submit-{{ $index }}" onclick="sendGradeFeedback( {{ $index }}, {{.ID}}, {{.SnapshotID}},{{ .Code }})" style="margin-top:3px; margin-bottom: 3px; color: #ffff;">Submit</button>
						</div>
						<div id="sub-feedback-block-{{ $index }}"></div>
					</div>
				</div>
				{{ end }}
			{{end}}
		</div>
		{{ end }}
		{{ end }}
	</div>
</div>

		<script>
			$(document).ready(function(){
				$('#view-exercise-link').attr("href", "/view_exercises"+window.location.search);
				$('#problem-dashboard-link').attr("href", "/problem_dashboard"+window.location.search+"&problem_id={{.Feedback.ProblemID}}");
				
				// Hids all the Check Feedback buttons
				// Snapshot feedback button
				document.getElementById("snapshot-check-feedback").classList.add("is-hidden");
				// Ask for help feedback button
				document.querySelectorAll('.help-check').forEach(function(button) {
					button.classList.add("is-hidden");
				});
				// Submission feedback button
				document.querySelectorAll('.sub-check').forEach(function(button) {
					button.classList.add("is-hidden")
				});
			});

			var snapshotCodeChanged = "";
			var snapshotcode = CodeMirror.fromTextArea(document.getElementById("snapshot-editor"), {
				lineNumbers: true, 
				mode: "{{getEditorMode .Feedback.ProblemName}}", 
				theme: "monokai", 
				matchBrackets: true, 
				indentUnit: 4, 
				indentWithTabs: true, 
				readOnly: false
			});
			snapshotcode.setSize("100%", "100%");
			snapshotcode.on('change', (snapshotcode) => {
				snapshotCodeChanged = snapshotcode.doc.getValue()
			});
			var snapshotCounter = 0;
			var feedbackCounter = 0;
			var subfeedbackCounter = 0; // should be 0 with NLP
			
			// var feedbackEditors = document.getElementsByClassName("feedback-editor");
			var feedbackChangedCode = []
			for (let i = 0; i<{{ len .Feedback.Messages }}; i++){
				var askHelpEditor = document.getElementById("feedback-editor-"+i)
				if (askHelpEditor != null ) {
					var code = CodeMirror.fromTextArea(askHelpEditor, {lineNumbers: true, mode: "{{getEditorMode .Feedback.ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: false});
					code.setSize("100%", "100%");
					code.on('change', (code) => {
						feedbackChangedCode[i] = code.doc.getValue()
					});
				}

			}

			var submissionEditors = document.getElementsByClassName("submission-editor");
			var submissionChangedCode = [];
			var submissionsGrade = [];
			for (let i = 0; i<submissionEditors.length; i++){
				var code = CodeMirror.fromTextArea(submissionEditors[i], {lineNumbers: true, mode: "{{getEditorMode .Feedback.ProblemName}}", theme: "monokai", matchBrackets: true, indentUnit: 4, indentWithTabs: true, readOnly: false});
				code.setSize("100%", "100%");
				code.on('change', (code) => {
					submissionChangedCode[i] = code.doc.getValue()
					// This changes the Send Feedback into Check my feedback button
					if (subfeedbackCounter === 0 ){
						document.getElementById("sub-submit-"+i).classList.add('is-hidden')
						document.getElementById("sub-check-"+i).classList.remove('is-hidden')
					}
				});
			}
function sendScaffoldingFeedback() {
    let feedback = window.editor ? window.editor.getValue().trim() : $("#custom-prompt-response").val().trim();
    
    if (feedback === "") {
        alert("Scaffolding feedback cannot be empty!");
        return;
    }

    $.post("/save_snapshot_feedback", {
        feedback: feedback,
        snapshot_id: {{ .Feedback.LastSnapshot.ID }},
        uid: {{ .Feedback.UserID }},
        role: "{{ .Feedback.UserRole }}"
    })
    .done(function(response) {
        alert("Scaffolding feedback posted successfully!");
        
        // Redirection after posting feedback
        window.location.replace("/student_dashboard_feedback_provision?student_id={{ .Feedback.StudentID }}&problem_id={{ .Feedback.ProblemID }}&uid={{ .Feedback.UserID }}&role={{ .Feedback.UserRole }}{{ if ne .Feedback.Password "" }}&password={{ .Feedback.Password }}{{ end }}");
    })
    .fail(function() {
        alert("Could not post scaffolding feedback. Please try again!");
    });
}


$(document).ready(function () {
    const urlParams = new URLSearchParams(window.location.search);
    const problem_id = urlParams.get('problem_id');

    if (problem_id) {
        fetchScaffoldingData(problem_id);
    } else {
        console.warn("No problem_id found in URL.");
        $("#scaffolding-dropdown").html('<option value="">No problem ID provided</option>');
    }

    // Event listener for dropdown change
    $("#scaffolding-dropdown").change(function () {
    const selectedOption = $(this).find("option:selected");
    console.log("Selected Option:", selectedOption.text()); // Debugging

    // Get the parent <optgroup> of the selected <option>
    const optgroupElement = selectedOption.parent("optgroup");
    console.log("Optgroup Element:", optgroupElement.length ? optgroupElement[0] : "Not found"); // Debugging

    const selectedCategory = optgroupElement.attr("label");
    console.log("Selected Category:", selectedCategory || "Undefined"); // Debugging

    if (selectedCategory) {
        const selectedLevel = selectedOption.val();
        const levelNumber = selectedLevel.split(" ")[1];

        const scaffoldingStrategy = getScaffoldingStrategyFromCategory(selectedCategory);
        fetchScaffoldingMaterial(problem_id, levelNumber, scaffoldingStrategy);
    } else {
        console.warn("No category found for selected option.");
    }
});

});

function fetchScaffoldingData(problem_id) {
    $.ajax({
        url: "/get_scaffolding",
        type: "POST",
        data: JSON.stringify({ problem_id: parseInt(problem_id) }),
        contentType: "application/json",
        success: function (data) {
            console.log("Received scaffolding data:", data); // Debugging
            let dropdown = $("#scaffolding-dropdown");
            dropdown.empty();

            // Define categories
            const categories = {
                1: "Fill in the Blanks",
                2: "Step-by-Step Tasks",
                3: "Guided Code With Hints",
                4: "Debug This Code",
                5: "Incremental Feature Implementation"
            };

            // Group scaffolding data by strategy
            let groupedData = {};
            data.forEach(function(item) {
                let category = categories[item.ScaffoldingStrategy] || "Other";
                if (!groupedData[category]) groupedData[category] = [];
                groupedData[category].push("Level " + item.ScaffoldingLevel);
            });

            // Populate dropdown correctly with optgroups
            Object.keys(groupedData).forEach(function(category) {
                let optgroup = $("<optgroup>").attr("label", category);

                groupedData[category].forEach(function(level) {
                    // Combining the level and the category name
                    let optionText = level + " (" + category + ")";
                    optgroup.append($("<option>").val(level).text(optionText));
                });

                dropdown.append(optgroup);
            });
        },
        error: function (err) {
            console.error("Error fetching scaffolding:", err);
            $("#scaffolding-dropdown").html("<option value=''>Failed to load</option>");
        }
    });
}

function fetchScaffoldingMaterial(problem_id, scaffolding_level, scaffolding_strategy) {
    $.ajax({
        url: "/get_scaffolding_material",
        type: "POST",
        data: JSON.stringify({
            problem_id: parseInt(problem_id),
            scaffolding_level: parseInt(scaffolding_level),
            scaffolding_strategy: scaffolding_strategy
        }),
        contentType: "application/json",
        success: function (data) {
            console.log("Scaffolding material response:", data); // Debugging log

            if (data && data.scaffolding_material) {
                // Destroy previous CodeMirror instance before updating the textarea
                if (window.editor) {
                    window.editor.toTextArea(); // Convert back to textarea
                    window.editor = null; // Reset the instance
                }

                // Update the textarea's value before initializing CodeMirror
                $("#custom-prompt-response").val(data.scaffolding_material).show();
                
                // Ensure the feedback box is visible
                $("#custom-prompt-box").show();

                // Initialize CodeMirror once
                window.editor = CodeMirror.fromTextArea(document.getElementById("custom-prompt-response"), {
                    lineNumbers: true,
                    mode: "markdown",
                    theme: "twilight",
                    matchBrackets: true,
                    indentUnit: 4,
                    indentWithTabs: true,
                    readOnly: false
                });
                window.editor.setSize(null, "650px");
                $(window.editor.getWrapperElement()).css("margin-top", "10px");

                // ✅ Ensure CodeMirror gets the updated value
                window.editor.setValue(data.scaffolding_material);

            } else {
                console.warn("No scaffolding material found.");

                // Reset the editor if no data is available
                if (window.editor) {
                    window.editor.toTextArea();
                    window.editor = null;
                }

                $("#custom-prompt-response").val("No scaffolding material available").show();

                // Reinitialize CodeMirror with the new message
                window.editor = CodeMirror.fromTextArea(document.getElementById("custom-prompt-response"), {
                    lineNumbers: true,
                    mode: "python",
                    theme: "monokai",
                    matchBrackets: true,
                    indentUnit: 4,
                    indentWithTabs: true,
                    readOnly: false
                });
                window.editor.setSize(null, "650px");
                window.editor.setValue("No scaffolding material available");
            }
        },
        error: function (err) {
            console.error("Error fetching scaffolding material:", err);

            // Reset the editor on error
            if (window.editor) {
                window.editor.toTextArea();
                window.editor = null;
            }

            $("#custom-prompt-response").val("Failed to load scaffolding material.").show();

            // Reinitialize CodeMirror with error message
            window.editor = CodeMirror.fromTextArea(document.getElementById("custom-prompt-response"), {
                lineNumbers: true,
                mode: "python",
                theme: "monokai",
                matchBrackets: true,
                indentUnit: 4,
                indentWithTabs: true,
                readOnly: false
            });
            window.editor.setSize(null, "650px");
            window.editor.setValue("Failed to load scaffolding material.");
        }
    });
}


// Helper function to map the category to the correct strategy
function getScaffoldingStrategyFromCategory(category) {
    const strategyMapping = {
        "Fill in the Blanks": "1",
        "Step-by-Step Tasks": "2",
        "Guided Code With Hints": "3",
        "Debug This Code": "4",
        "Incremental Feature Implementation": "5"
    };

    return strategyMapping[category] || "";
}

const snapshotId = {{ .Feedback.LastSnapshot.ID }};
const userId = {{ .Feedback.UserID }};
const userRole = "{{ .Feedback.UserRole }}";

function sendCustomPromptFeedback() {
    const button = document.getElementById("submit-custom-prompt");
    button.textContent = "Generating...";
    button.disabled = true;

    const params = new URLSearchParams(window.location.search);
    const studentID = params.get("student_id");
    const problemID = params.get("problem_id");

    if (studentID && problemID && typeof snapshotId !== 'undefined' && typeof userId !== 'undefined' && typeof userRole !== 'undefined') {
        // Construct full URL with all required query params
        const actionUrl = '/sc_view?student_id=' + encodeURIComponent(studentID) +
                          '&problem_id=' + encodeURIComponent(problemID) +
                          '&snapshot_id=' + encodeURIComponent(snapshotId) +
                          '&uid=' + encodeURIComponent(userId) +
                          '&role=' + encodeURIComponent(userRole);

        const tempForm = document.createElement('form');
        tempForm.method = 'POST';
        tempForm.action = actionUrl;

        // Optional: Add them as hidden inputs too
        const hiddenFields = {
            student_id: studentID,
            problem_id: problemID,
            snapshot_id: snapshotId,
            uid: userId,
            role: userRole,
			isNew: 'false'
        };

        for (const key in hiddenFields) {
            const input = document.createElement('input');
            input.type = 'hidden';
            input.name = key;
            input.value = hiddenFields[key];
            tempForm.appendChild(input);
        }

        document.body.appendChild(tempForm);
        tempForm.submit();
    } else {
        alert("Missing required parameters.");
        button.textContent = "Get Feedback from AI";
        button.disabled = false;
    }
}





			function runNLP(student_code, ta_code, user_id, write) {
				return $.ajax({
					// url: "http://127.0.0.1:5000/feedback_classify",
					url: "http://141.225.10.71:5000/feedback_classify",
					type: "POST",
					data: JSON.stringify({
							student_code: student_code,
							ta_code: ta_code,
							user_id: user_id
					}),
					"headers": {
						"Content-Type": "application/json"
					},
					success: function(data){
						var pre = ''
						var str = '';
						data.results.forEach(function(item){
							pre += '<div class="card" style="margin-top: 15px; padding: 5px;"><div class="card-content"><ul style="margin:5px">'

							pre += '<li><strong>' + "Feedback: " + '</strong>' + item.feedback + '</li>'

							str = '<li><strong>Suggestions: </strong> <ul style="list-style-type:square">'
							item.output.forEach(function(suggestion){
								str += '<li>' + suggestion + '</li>'
							});
							str += '</ul></li>'
							pre += str
							pre += '</ul></div></div>'
						})
						// Do not show NLP classifier result
						$(write).html("");
						$('<div class="wrapper">' + pre +  '</div>').appendTo( write )
					},
					error: function(err) {
						// alert(JSON.stringify(err));
					}
				});
			}

			function codeSnapshotFeedback(code, user_id) {
				// Check if the code is changed.
				if (snapshotCodeChanged == "" ) {
					alert("Please provide in-line feedback!");
					return
				}
				runNLP(code, snapshotCodeChanged, user_id, '#code-snapshot-feedback-block');
				// Hide the Check button & show the Send button
				document.getElementById("snapshot-send-feedback").classList.remove('is-hidden')
				document.getElementById("snapshot-check-feedback").classList.add('is-hidden')

			}

			function sendSnapshotFeedback(code, user_id) {
				// Check if the code is changed.
				if (snapshotCodeChanged == "" ) {
					alert("Please provide in-line feedback!");
					return
				}
				runNLP(code, snapshotCodeChanged, user_id, '#code-snapshot-feedback-block');
				$.post("/save_snapshot_feedback", {feedback: snapshotCodeChanged, snapshot_id: {{.Feedback.LastSnapshot.ID}}, uid: {{ .Feedback.UserID}}, role: {{ .Feedback.UserRole}}{{if ne .Feedback.Password ""}}, password: {{ .Feedback.Password}}{{end}}  }, function(data, status){
					alert("Feedback posted successfully!");
					window.location.replace("/student_dashboard_feedback_provision?student_id={{ .Feedback.StudentID}}&problem_id={{ .Feedback.ProblemID}}&uid={{ .Feedback.UserID}}&role={{.Feedback.UserRole}}{{if ne .Feedback.Password ""}}&password={{.Feedback.Password}}{{end}}");
				})
				.fail(function() {
					alert("Could not post the feedback. Please try again!");
				});
			}

			function messageFeedback(i,code,message_id){
				// Check if the code is changed.
				if (feedbackChangedCode[i] === undefined ) {
					alert("Please provide in-line feedback!");
					return
				}

				runNLP(code, feedbackChangedCode[i], {{ .Feedback.UserID }}, '#feedback-block-'+i );
				// Hide the Check button & show the Send button
				document.getElementById("help-send-feedback-"+i).classList.remove('is-hidden')
				document.getElementById("help-check-feedback-"+i).classList.add('is-hidden')
			}

			function sendMessageFeedback(i,code,message_id) {
				// Check if the code is changed.
				if (feedbackChangedCode[i] === undefined ) {
					alert("Please provide in-line feedback!");
					return
				}
				runNLP(code, feedbackChangedCode[i], {{ .Feedback.UserID }}, '#feedback-block-'+i );
				$.post("/save_message_feedback", {feedback: feedbackChangedCode[i], message_id: message_id, uid: {{ .Feedback.UserID}}, role: {{ .Feedback.UserRole}}{{if ne .Feedback.Password ""}}, password: {{ .Feedback.Password}}{{end}}  }, function(data, status){
					alert("Feedback posted successfully!");
					window.location.replace("/student_dashboard_feedback_provision?student_id={{ .Feedback.StudentID}}&problem_id={{ .Feedback.ProblemID}}&uid={{ .Feedback.UserID}}&role={{.Feedback.UserRole}}{{if ne .Feedback.Password ""}}&password={{.Feedback.Password}}{{end}}");
				})
				.fail(function() {
					alert("Could not post the feedback. Please try again!");
				});

			}

			function setGrade(i, grade) {
				submissionsGrade[i] = grade;
				// document.getElementById("sub-submit-"+i).removeAttribute("disabled");
				// Hide the Check button & show the Submit button
				document.getElementById("sub-submit-"+i).innerHTML = "Submit Grade";
				document.getElementById("sub-submit-"+i).classList.remove('is-hidden')
				document.getElementById("sub-check-"+i).classList.add('is-hidden')

				// check if code changed and already run through nlp
				if (subfeedbackCounter != 0 && submissionChangedCode[i] != undefined ) {
					document.getElementById("sub-submit-"+i).innerHTML = "Submit Grade and Feedback";
				}

			}
			function removeGrade(index) {
				submissionsGrade[index] = undefined;
				// check if code changed and already run through nlp
				if (subfeedbackCounter != 0) {
					document.getElementById("sub-submit-"+index).innerHTML = "Submit Feedback";
				}

			}
			function checkSubFeedback (i, submission_id, snapshot_id, submittedCode) {
				// Check if the code is changed.
				var code = submissionChangedCode[i]
				if (code === undefined ) {
					alert("Please provide in-line feedback!");
					return
				}

				runNLP(submittedCode, code, {{ .Submission.UserID }}, '#sub-feedback-block-'+i );
				subfeedbackCounter++;
				// document.getElementById("sub-submit-"+i).removeAttribute("disabled");
				document.getElementById("sub-check-"+i).classList.add('is-hidden')
				document.getElementById("sub-submit-"+i).classList.remove('is-hidden')

				if (submissionsGrade[i] !== undefined && subfeedbackCounter != 0 ){
					document.getElementById("sub-submit-"+i).innerHTML = "Submit Grade and Feedback";
				} else {
					document.getElementById("sub-submit-"+i).innerHTML = "Submit Feedback";
				}
				
			}
			
			function sendGradeFeedback(i, submission_id, snapshot_id, submittedCode) {
				var code = submissionChangedCode[i]

				var grade = submissionsGrade[i]

				// if no grade, no code change, disable submit
				// if grade but no code change, send grades only.
				// if grade and  code change, disable submit to run it through nlp.
				// if grade and code change and already run through nlp, enable submit. This will send both grade and feedback.
				
				if (subfeedbackCounter != 0 && code !== undefined ) {
					if (grade !== undefined ) {
						$.post("/teacher_grades", {content: submittedCode, changed: "", decision: grade, sid: submission_id, uid: {{ .Submission.UserID}}, role: {{ .Submission.UserRole}}{{if ne .Submission.Password ""}}, password: {{ .Submission.Password}}{{end}}  }, function(data, status){
								// Save and send feedback if the code is changed
								if (code !== undefined) {
									runNLP(submittedCode, code, {{ .Submission.UserID }}, '#sub-feedback-block-'+i );
									$.post("/save_snapshot_feedback", {snapshot_id: snapshot_id, feedback: code, uid: {{ .Submission.UserID}}, role: {{ .Submission.UserRole}}{{if ne .Submission.Password ""}}, password: {{ .Submission.Password}}{{end}} }, function(data1, status1){
										alert("Graded successfully! Feedback posted successfully! ");
										window.location.reload();
									})
									.fail(function() {
										alert("Could not post the feedback. Please try again!");
									});
								} 
								// else {
								// 	alert("Graded successfully! ");
								// 	window.location.reload();
								// }
						})
						.fail(function() {
							alert("Could not grade the submission. Please try again!");
						});
					} else {
						runNLP(submittedCode, code, {{ .Submission.UserID }}, '#sub-feedback-block-'+i );
						$.post("/save_snapshot_feedback", {snapshot_id: snapshot_id, feedback: code, uid: {{ .Submission.UserID}}, role: {{ .Submission.UserRole}}{{if ne .Submission.Password ""}}, password: {{ .Submission.Password}}{{end}} }, function(data1, status1){
							alert("Feedback posted successfully! ");
							window.location.reload();
						})
						.fail(function() {
							alert("Could not grade the submission. Please try again!");
						});
					}
				}

				if ( code === undefined ){
					if (grade !== undefined ) {
						$.post("/teacher_grades", {content: submittedCode, changed: "", decision: grade, sid: submission_id, uid: {{ .Submission.UserID}}, role: {{ .Submission.UserRole}}{{if ne .Submission.Password ""}}, password: {{ .Submission.Password}}{{end}}  }, function(data, status){
							alert("Graded successfully!");
							window.location.reload();
						})
						.fail(function() {
							alert("Could not grade the submission. Please try again!");
						});

					} else {
						alert("Please provide in-line feedback!");
						return;
					}
				
				}
				
			}
			
		</script>

	</body>
	</html>
`

var ACTIVITY_VIEW_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['bar','line']});
      google.charts.setOnLoadCallback(drawStudentCount);
      google.charts.setOnLoadCallback(drawExerciseCount);

      function drawStudentCount() {
	      var data = google.visualization.arrayToDataTable([
	        ['Time', 'Student Count'],
			{{ range $day, $val := . }}
				[  new Date({{$day}} / 1000000), {{$val.SidCount}} ],
			{{ end }}
	      ]);
	      var options = {
	        title: 'Daily student participation',
        	height: 350,
	        vAxis: { textStyle: { fontSize: 20} },
            hAxis: { title: '', textStyle: { fontSize: 20} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Bar(document.getElementById('chart1_div'));
        chart.draw(data, google.charts.Bar.convertOptions(options));
      }

      function drawExerciseCount() {
	      var data = google.visualization.arrayToDataTable([
	        ['Time', 'Excercise Count'],
			{{ range $day, $val := . }}
				[  new Date({{$day}} / 1000000), {{$val.PidCount}} ],
			{{ end }}
	      ]);
	      var options = {
	        title: 'Daily exercise',
        	height: 350,
	        vAxis: { textStyle: { fontSize: 20} },
            hAxis: { title: '', textStyle: { fontSize: 20} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Bar(document.getElementById('chart2_div'));
        chart.draw(data, google.charts.Bar.convertOptions(options));
      }
    </script>
    <style>
    #chart1_div,#chart2_div{ margin: auto; width:75%; }
    .spacer{ width:100%; height:40px; }
    </style>
  </head>

  <body>
    <div id="chart1_div"></div>
    <div class="spacer"></div>
    <div id="chart2_div"></div>
  </body>
</html>
`

var VIEW_ANSWERS_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['bar']});
      google.charts.setOnLoadCallback(drawChart);
      function drawChart() {
        var data = google.visualization.arrayToDataTable([
          ["Answer", "Count", {role: 'annotation'}],
          	{{$total := .Total}}
			{{ range $key, $value := .Counts }}
				[{{ $key }}, {{ $value }}, Math.round(100 * {{$value}} / {{ $total }})  + '%'],
			{{ end }}
		]);
        var options = {
        	'title':'Total votes: {{ .Total }}',
        	'height':300,
        	'legend': {position: "none"},
        	'fontSize': 24,
            vAxis: { textStyle: { fontSize: 20} },
            hAxis: { title: "", textStyle: { fontSize: 20} },
        };
        var chart = new google.charts.Bar(document.getElementById('chart_div'));
        chart.draw(data, google.charts.Bar.convertOptions(options));
      }
    </script>
    <style>
    #chart_div{ margin: auto; width:70%; }
    pre{ margin: auto; width:60%}
    .spacer{ width:100%; height:40px; }
    </style>
  </head>

  <body>
    <div id="chart_div"></div>
    <div class="spacer"></div>
    <pre id="content">{{ .Content }}</pre>
    <div class="spacer"></div>
  </body>
</html>
`

var ANALYZE_SUBMISSIONS_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['corechart','bar']});
      google.charts.setOnLoadCallback(drawWaitingTime);
      google.charts.setOnLoadCallback(drawResponseTime);
      google.charts.setOnLoadCallback(drawAttempts);

      function drawResponseTime() {
        var data = google.visualization.arrayToDataTable([
			['Student', 'Duration'],
			{{ range $sid, $rec := . }}
				{{ range $rec }}
				[  String({{$sid}}), ({{.At}} - {{.Start}})/1e9],
				{{ end }}
			{{ end }}
        ]);

        var options = {
          title: 'Student response time',
          legend: { position: 'none' },
        };

        var chart = new google.visualization.Histogram(document.getElementById('response'));
        chart.draw(data, options);
      }

      //---------------------------------------------------
      function drawWaitingTime() {
        var data = google.visualization.arrayToDataTable([
			['Student', 'Waiting time'],
			{{ range $sid, $rec := . }}
				{{ range $rec }}
				[  String({{$sid}}), ({{.Completed}} - {{.At}})/1e9],
				{{ end }}
			{{ end }}
        ]);

        var options = {
          title: 'Waiting for teacher',
          legend: { position: 'none' },
        };

        var chart = new google.visualization.Histogram(document.getElementById('waiting'));
        chart.draw(data, options);
      }

      //---------------------------------------------------
      function drawAttempts() {
        var data = google.visualization.arrayToDataTable([
			['Student', 'Attempts'],
			{{ range $sid, $rec := . }}
				{{ $length := len $rec }}
				[  String({{$sid}}), {{$length}} ],
			{{ end }}
        ]);

        var options = {
			title: 'Solution attempts',
			legend: { position: 'none' },
	        hAxis: {
	            viewWindowMode:'explicit',
		        viewWindow: { min:1 }
	        },
        };

        var chart = new google.visualization.Histogram(document.getElementById('attempts'));
        chart.draw(data, options);
      }


      //---------------------------------------------------
    </script>
    <style>
    .spacer{ width:100%; height:40px; }
    .row{
	    display:flex;
	    flex-direction:row;
	    justify-content: space-around;
    }
    #attempts,#response,#waiting{
	    width:420px; height:400px;
	    display:flex;
	    flex-direction:column;
    }
    </style>
  </head>
  <body>
  	<div class="row">
    <div id="attempts"></div>
    <div id="response"></div>
    <div id="waiting"></div>
	</div>
  </body>
</html>
`

var TAGS_VIEW_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['corechart', 'table']});
      google.charts.setOnLoadCallback(drawActivity);
      google.charts.setOnLoadCallback(drawScores);

      function drawActivity() {
	      var data = google.visualization.arrayToDataTable([
	        ['Date', 'Total submissions'],
			{{ range $day, $val := .SubmissionCount }}
				[  {{$day}}, {{$val}} ],
			{{ end }}
	      ]);
	      var options = {
	        title: '',
        	height: 350,
            hAxis: { title: 'Total submissions' },
	        legend: { position: 'none' },
	      };
        var chart = new google.visualization.SteppedAreaChart(document.getElementById('chart_div'));
        chart.draw(data, options);
      }

      function drawScores() {
        var data = new google.visualization.DataTable();
        data.addColumn('string', 'Name');
        data.addColumn('number', 'Points');
        data.addColumn('number', 'Attempts');
        data.addRows([
			{{ range $student_id, $entry := .Scores }}
				[{{$entry.Name}},{{$entry.Points}},{{$entry.Attempts}}],
			{{ end }}
        ]);
        var table = new google.visualization.Table(document.getElementById('scores_div'));
        table.draw(data, {showRowNumber: true, width: '400px'});
      }
    </script>
    <style>
    body { margin:auto; width:90%; font-size:16pt;}
    #chart_div,#scores_div{ margin:auto; }
    .spacer{ width:100%; height:30px; }
    </style>
  </head>
  <body>
  	<div class="spacer"></div>
  	<h4>Activities</h4>
  	<div id="chart_div"></div>
  	<div class="spacer"></div>
  	<h4>Learning objectives</h4>
  	<ul>
	{{$pc := .PC}}
	{{ range $tag_id, $tag_des := .Tags }}
	<li><a href="report_tag?pc={{$pc}}&tag_id={{$tag_id}}" target="_blank">{{$tag_des}}</a></li>
	{{ end }}
	</ul>
	<h4>Points</h4>
	<div id="scores_div"></div>
  	<div class="spacer"></div>
  </body>
</html>
`

var TAG_REPORT_TEMPLATE = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['bar','scatter']});
      google.charts.setOnLoadCallback(draw_success);
      google.charts.setOnLoadCallback(draw_participation);

      function draw_success() {
	      var data = google.visualization.arrayToDataTable([
	        ['Time', 'Success'],
			{{ range $pid, $rec := .Performance }}
				[ new Date({{$rec.Timestamp}} / 1000000), {{$rec.Success}} ],
			{{ end }}
	      ]);
	      var options = {
	        title: 'Success', legend: {position: 'none'},
	        vAxis: {
	        	textStyle: { fontSize: 18},
	            viewWindowMode:'explicit',
    	        viewWindow: { min:0, max:1.05 }
	        },
            hAxis: { title: '', textStyle: {fontSize: 18} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Scatter(document.getElementById('success'));
        chart.draw(data, google.charts.Scatter.convertOptions(options));
      }

      function draw_participation() {
	      var data = google.visualization.arrayToDataTable([
	        ['Time', 'Participation'],
			{{ range $pid, $rec := .Performance }}
				[ new Date({{$rec.Timestamp}} / 1000000), {{$rec.Activity}} ],
			{{ end }}
	      ]);
	      var options = {
	        title: 'Participation', legend: {position: 'none'},
	        vAxis: {
	        	textStyle: {fontSize: 18},
	            viewWindowMode:'explicit',
    	        viewWindow: { min:0, max:1.05 }
	        },
            hAxis: { title: '', textStyle: {fontSize: 18} },
        	fontSize: 24,
	      };
        var chart = new google.charts.Scatter(document.getElementById('participation'));
        chart.draw(data, google.charts.Scatter.convertOptions(options));
      }
      </script>
    <style>
    body{ margin: auto; width:75%; }
    .row{
	    display:flex;
	    flex-direction:row;
	    justify-content: space-around;
    }
    #success,#participation{
	    width:450px; height:400px;
	    display:flex;
	    flex-direction:column;
    }
    .spacer{ width:100%; height:40px; }
    #problem_ids{
    	margin:auto;
    	height:100px;
    	padding-top:20px;
		overflow-x: scroll;
	    white-space: nowrap;
    	text-align: center;
		vertical-align: middle;
    }
    .problem_id{
		padding: 15px 20px 15px 20px;
    	text-align: center;
    	font-size: 110%;
    	border:2px solid #dedede;
		display: inline-block;
    }
    .problem_id a{
    	text-decoration: none;
    }
    </style>
  </head>

  <body>
  	<div class="spacer"><h2>{{.Description}}</h2></div>
	<div class="row">
	    <div id="success"></div>
	    <div id="participation"></div>
    </div>
	<div class="spacer"></div>
    <div id="problem_ids">
	{{ range $pid, $rec := .Performance }}
		<div class="problem_id"><a href="analyze_submissions?pid={{$pid}}&pc={{$rec.PC}}" target="_blank">{{$pid}}</a></div>
	{{ end }}
	</div>
  </body>
</html>
`

var ANALYSIS_TEMPLATE = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="./vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>CodeInsight</title>
    <script type="module" crossorigin src="./assets/index-TFIfn1AY.js"></script>
    <link rel="stylesheet" crossorigin href="./assets/index-BC89AxIr.css">
  </head>
  <body>
    <div id="root"></div>
  </body>

</html>
`
