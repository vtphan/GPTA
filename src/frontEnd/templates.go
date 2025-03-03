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
			<span>Dashboard for Problem: {{.ProblemName}}</span>
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
			<span>Coding Status: <strong>{{ .Status.CodingStat }} </strong></span>
			<span>Help Status: <strong>{{ .Status.HelpStat }} </strong></span>
			<span>Submission Status: <strong> {{ .Status.SubmissionStat }} </strong></span>
			<span>Tutoring Status: <strong>{{ .Status.TutoringStat }} </strong></span>
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
var PROBLEM_DASHBOARD_TEMPLATE = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>Problem Dashboard</title>
<meta http-equiv="refresh" content="10" >
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
		background-color: #AE1431;
		color: #FFFFFF;
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
	  <span>Dashboard for Problem: {{.ProblemName}}</span>
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
			<button id="deactivate-button" class="button is-danger">Deactivate!</button>
		{{end}}
	{{end}}
	<h4 class="title is-4">Exercise Statement</h4>
	<div class="accordions">
		<h3>{{.ProblemName}}</h3>
		<div>
			<textarea id="editor">{{ .Code }}</textarea>
		</div>
	</div>
	<h4 class="title is-4">Statistics for {{.ProblemName}}</h4>
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
					<th>Coding Status</th>
					<th>Help Status</th>
					<th>Submission Status</th>
					<th>Tutoring Status</th>
				</tr>
			</thead>
			<tbody>
				{{range .StudentInfo}}
				<tr>
					<td>{{if ne .CodingStat "Idle"}}<a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{$.ProblemID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}#code-snapshot">{{.StudentName}}</a>{{else}}{{.StudentName}}{{end}}</td>
					<td>{{if and (eq $.IsActive true) (ne .CodingStat "Idle") (ne .LastUpdatedAt.IsZero true) }}<a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{$.ProblemID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}">{{ formatTimeSince .LastUpdatedAt }} ago</a>{{end}}</td>
					<td>{{.CodingStat}}</td>
					<td>{{if ne .HelpStat ""}}<a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{$.ProblemID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}#ask-for-help">{{.HelpStat}}</a>{{end}}</td>
					<td>{{if ne .SubmissionStat ""}}<a href="/student_dashboard_code_snapshot?student_id={{.StudentID}}&problem_id={{$.ProblemID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}#submission">{{.SubmissionStat}}</a>{{end}}</td>
					<td>{{.TutoringStat}}</td>
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
		  $(document).ready(function(){
			$('#view-exercise-link').attr("href", "/view_exercises"+window.location.search);
			$('#deactivate-button').click(function(){
				if (confirm("Deactivate the problem?") == true) {
					$.post("/teacher_deactivates_problems", {filename: {{.ProblemName}}, uid: {{.UserID}}, role: {{.UserRole}}{{if ne .Password ""}}, password: {{.Password}}{{end}} })
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
  top: 20px; /* Adjust as needed */
  right: 20px; /* Adjust as needed */
  background-color: #2196F3;
  color: white;
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
  transition: background-color 0.3s ease;
}

.settings-button:hover {
  background-color: #1976D2;
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
	<div class="topcorner">{{.Username}}({{.UserRole}})</div>
		<h2 class="title is-2">Exercises</h2>
		{{if ne .UserRole "student"}}
		<a id="new-problem" class="button is-success" href="">
			<span class="icon is-small">
			<i class="fa-solid fa-plus"></i>
			</span>
			<span style="color: #242424;">Broadcast New Exercise</span>
		</a>
		<a id="export-button" class="button is-primary" href="">
			<span class="icon is-small">
			<i class="fa-solid fa-plus"></i>
			</span>
			<span style="color: #242424;">Export Score</span>
		</a>
		{{end}}

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
		<button class="settings-button" id="settings-button">
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
					<tr {{if eq .IsActive true}}class="is-selected"{{end}} style="color: #242424;">
						<td><a href="/problem_dashboard?problem_id={{.ID}}&uid={{$.UserID}}&role={{$.UserRole}}{{if ne $.Password ""}}&password={{$.Password}}{{end}}">{{.Filename}}</a></td>
						<td>{{ .UploadedAt.Format "Jan 02, 2006 3:04:05 PM" }}</td>
						<td>{{.Attendance}}</td>
						<td>{{.NumActive}}</td>
						<td>{{.NumHelpRequest}}</td>
						<td>{{.NumGradedCorrect}}</td>
						<td>{{.NumGradedIncorrect}}</td>
						<td>{{.NumNotGraded}}</td>
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
      background: linear-gradient(135deg, #6a11cb, #2575fc);
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
      background: linear-gradient(45deg, #6a11cb, #2575fc);
      color: white;
      border: none;
    }
    .button.is-primary:hover {
      background: linear-gradient(45deg, #2575fc, #6a11cb);
    }
    .button.is-info {
      background: linear-gradient(45deg, #1e90ff, #00bcd4);
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
  <h1 class="title">Course Settings</h1>
  
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
      <button class="button is-primary" id="toggle-peer-tutoring">
        Peer Tutoring
        <label class="switch" style="margin-left: 10px;">
          <input id="peer_tutoring_button" type="checkbox">
          <span class="slider round"></span>
        </label>
      </button>
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
        background: linear-gradient(135deg, #6a11cb, #2575fc);
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
        background: linear-gradient(135deg, #6a11cb, #2575fc);
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
      background: linear-gradient(45deg, #6a11cb, #2575fc);
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
			<span>Problem Broadcast</span>
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
					<span class="file-label" style="color: #242424;">
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
			$.post("/teacher_broadcasts", {role: "teacher", uid: uid, content: editor.getValue(), answer: answer, merit: points, effort: effort, attempts: attempt, tag: tag, filename: filename, exact_answer: "True"}, function(data, status){
				if (status == "success"){
					alert("Exercise broadcasted successfully!");
					window.location.replace("/view_exercises?role=teacher&uid="+uid);
				} else {
					alert("Failed to broadcast. Try agian!");
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
	<link rel="stylesheet" href="https://code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css" />
	<script src="https://cdn.jsdelivr.net/npm/@creativebulma/bulma-collapsible"></script>
	<style>
		#code-snapshot {
			background: darkseagreen;
			padding: 20px;
			margin: 30px;
			padding-bottom: 0px;
			border-radius: 25px;
		}
		#ask-for-help {
			background: #c1bb91;
			padding: 20px;
			margin: 30px;
			padding-bottom: 0px;
			border-radius: 25px;
		}
		#submission {
			background: #ada192;;
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
			<span>Dashboard for Problem: {{ .Feedback.ProblemName}}</span>
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
	<div class="column show" style="width: 70%;">
		<!--
		<div class="row">
			<h3 class="title is-2" style="margin-bottom: 0px;">{{ .Feedback.StudentName}}'s Dashboard for {{ .Feedback.ProblemName}}</h3>
		</div>
		-->
		<div class="row status">
			<span>Coding Status: <strong>{{ .Status.CodingStat }} </strong></span>
			<span>Help Status: <strong>{{ .Status.HelpStat }} </strong></span>
			<span>Submission Status: <strong> {{ .Status.SubmissionStat }} </strong></span>
			<span>Tutoring Status: <strong>{{ .Status.TutoringStat }} </strong></span>
		</div>

		<div class="tabs">
			<ul>
				<li class="is-active"><a>CodeSpace</a></li>
				<li><a href="/student_dashboard_feedback_provision?student_id={{.Feedback.StudentID}}&problem_id={{.Feedback.ProblemID}}&uid={{.Feedback.UserID}}&role={{.Feedback.UserRole}}{{if ne .Feedback.Password ""}}&password={{.Feedback.Password}}{{end}}">Feedback History</a></li>
			</ul>
		</div>

	</div>

	<div class="content">
		<h3>Student's latest code snapshot: </h3>
		<div id="code-snapshot">
			<div class="box" style="padding: 0px; margin-bottom: 3.5rem; border: 5px solid; border-radius: 10px;">
				<div class="message-header">
					<div class="column is-two-thirds">
						<p>Latest Code Snapshot at {{.Feedback.LastSnapshot.LastUpdated.Format "Jan 02, 2006 3:04:05 PM"}}</p>
					</div>
				</div>	
				<div id="feedback-block-99999"></div>
				<div style="background: darkseagreen;">
					<textarea id="snapshot-editor"> {{ .Feedback.LastSnapshot.Code }} </textarea>
					<div class="actions">
							<button class="button is-info" id="snapshot-check-feedback" onclick="codeSnapshotFeedback({{ .Feedback.LastSnapshot.Code }}, {{ .Feedback.UserID }})" style="margin-top:3px; margin-bottom: 3px; color: #000000;" >Check My Feedback</button>
							<button class="button is-info" id="snapshot-send-feedback" onclick="sendSnapshotFeedback({{ .Feedback.LastSnapshot.Code }}, {{ .Feedback.UserID }})" style="margin-top:3px; margin-bottom: 3px; color: #000000;" >Send Feedback</button>
							<button class="button is-info chatgpt-feedback" id="chatgpt-feedback-99999" onclick="getChatGptFeedback( 99999 ,{{ .Feedback.LastSnapshot.Code }} , {{ .Feedback.UserID }})" style="margin-top:3px; margin-bottom: 3px; color: #000000;" >ChatGPT Feedback</button>
							<button class="button is-info" id="toggle-custom-prompt" onclick="toggleCustomPromptBox()" style="margin-top:3px; color: #000000; display: flex; align-items: center;">
            Enable Custom Prompt
            <label class="switch" style="margin-left: 10px;">
                <input id="custom_prompt_toggle" type="checkbox">
                <span class="slider round"></span>
            </label>
        </button>
					</div>
					<div id="code-snapshot-feedback-block"></div>

				<div id="custom-prompt-box" class="box" style="background: #c1bb91; display: none; margin-top: 10px;">
    <h3 class="title is-4">Enter Your Custom Prompt</h3>
    <textarea id="custom-prompt-input" class="textarea" placeholder="Type your custom prompt here..."></textarea>
    <button class="button is-primary" id="submit-custom-prompt" 
        onclick="sendCustomPromptFeedback({{ .Feedback.LastSnapshot.Code }}, {{ .Feedback.UserID }}, document.getElementById('custom-prompt-input').value)" 
        style="margin-top: 10px;">
        Send Prompt
    </button>
    <div id="custom-prompt-response" style="margin-top: 15px;"></div>
</div>

			</div>
		</div>
		

		{{ if .Feedback.Messages}}
		<h3>Student's help requests: </h3>
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
								
							<div style="background: #c1bb91;">
									<div>
										<textarea class="feedback-editor" id="feedback-editor-{{ $index }}">{{ .Code }}</textarea>
									</div>

									<div class="actions">
										<button class="button is-info help-check" id="help-check-feedback-{{ $index }}" onclick="messageFeedback( {{ $index }} ,{{ .Code }} , {{ .ID }})" style="margin-top:3px; margin-bottom: 3px; color: #000000;" >Check My Feedback</button>
										<button class="button is-info chatgpt-feedback" id="chatgpt-feedback-{{ $index }}" onclick="getChatGptFeedback( {{ $index }} ,{{ .Code }} , {{ .ID }})" style="margin-top:3px; margin-bottom: 3px; color: #000000;" >ChatGPT Feedback</button>
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
		<h3>Student's submissions: </h3> 
		<div id="submission">
			{{range $index, $el := .Submission.Submissions}}
				{{ if eq .Grade "" }}
				<div class="box" style="padding: 0px; margin-bottom: 3.5rem; border: 5px solid; border-radius: 10px;">
				
					<div class="message-header">
						<div class="column is-two-thirds">
							<p>Submitted at {{.SubmittedAt.Format "Jan 02, 2006 3:04:05 PM"}}</p>
							{{if eq .Grade ""}} Not Graded {{else}} Graded {{if eq .Grade "correct"}} <span class="tag is-success">correct</span> {{else if eq .Grade "incorrect"}} <span class="tag is-danger">incorrect</span> {{else}} {{.Grade}} {{end}} {{end}}
						</div>

						<div class="column buttons" style="padding-left: 1px;">
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
					
					<div style="background: #ada192;">
						<div>
							<textarea class="submission-editor" id="editor-{{.ID}}">{{ .Code }}</textarea>
						</div>
						<div class="sub-actions">
							<button class="button is-info sub-check" id="sub-check-{{ $index }}" onclick="checkSubFeedback( {{ $index }}, {{.ID}}, {{.SnapshotID}},{{ .Code }})" style="margin-top:3px; margin-bottom: 3px; color: #000000;">Check My Feedback</button>
							<button class="button is-info sub-submit" id="sub-submit-{{ $index }}" onclick="sendGradeFeedback( {{ $index }}, {{.ID}}, {{.SnapshotID}},{{ .Code }})" style="margin-top:3px; margin-bottom: 3px; color: #000000;">Submit</button>
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
function toggleCustomPromptBox() {
        var checkbox = document.getElementById('custom_prompt_toggle');
        var promptBox = document.getElementById('custom-prompt-box');
        promptBox.style.display = checkbox.checked ? 'block' : 'none';
    }
			function getChatGptFeedback(idx, code, user_id) {
				return $.ajax({
					url: "/instructions_with_example",
					type: "POST",
					data: JSON.stringify({
						problem: {{.Feedback.ProblemName}},
						course: "{{$.CourseName}}",
						duration: 15,
						solutions: [
							{
								solution_id: 1,
								code: code,
								minute_left: 5
							}
						]
					}),
					"headers": {
						"Content-Type": "application/json",
						// 'Accept': 'application/json',
						// 'Origin': 'http://141.225.10.71:8080'
					}
					// success: function(data){
					// 	$('#feedback-block-'+i).html(data.results);
					// 	document.getElementById("help-send-feedback-"+i).classList.remove('is-hidden')
					// }
					// error: function(err) {
					// 	alert(JSON.stringify(err));
					// }
				}).done(function(data){
					if (data.feedbacks.length > 0) {
						$('#feedback-block-'+idx).html("<br/><h4>ChatGPT Feedback</h4>"+data.feedbacks[0].feedback);
						document.getElementById("feedback-block-"+idx).classList.remove('is-hidden')
					} else {
						alert("No feedback found!");
					}
					
				}).fail(function(err) {
					alert(JSON.stringify(err));
				});
			}
// Add a new function to handle sending the code and custom prompt
function sendCustomPromptFeedback(code, userID, customPrompt) {
    return $.ajax({
        url: "/process_code_with_prompt", // New endpoint
        type: "POST",
        data: JSON.stringify({
            code: code,
            custom_prompt: customPrompt,
            user_id: userID
        }),
        "headers": {
            "Content-Type": "application/json",
        }
    }).done(function(data) {
        // Process the feedback received from the server
        if (true) {
            $('#custom-prompt-response').html("<br/><h4>Feedback</h4>" + data.feedback);
        } else {
            alert("No feedback found!");
        }
    }).fail(function(err) {
        alert("Error: " + JSON.stringify(err));
    });
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
