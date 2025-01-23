// Author: Vinhthuy Phan, 2018
package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"strconv"
)

// -----------------------------------------------------------------------------------
type ScoreEntry struct {
	Name     string
	Points   int
	Attempts int
	Count    int
}

type TagsViewData struct {
	Tags            map[int]string
	SubmissionCount map[string]int
	Scores          map[int]*ScoreEntry
	PC              string
}

// -----------------------------------------------------------------------------------
func reportHandler(w http.ResponseWriter, r *http.Request) {
	// Check passcode
	if r.FormValue("pc") != Passcode {
		fmt.Fprintf(w, "Unauthorized")
		return
	}

	// Fetch tags
	tags, err := GetTags()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch student scores
	scores, err := GetStudentScores()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Prepare the data for rendering
	record := &TagsViewData{
		Tags:            tags,
		Scores:          scores,
		SubmissionCount: make(map[string]int), // You can add logic to fill this if needed
		PC:              Passcode,
	}

	// Render the view
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.New("").Parse(TAGS_VIEW_TEMPLATE)
	err = t.Execute(w, record)
	if err != nil {
		fmt.Println(err)
	}
}

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

// -----------------------------------------------------------------------------------
type ProblemPerformance struct {
	Pid       int
	Timestamp int64
	Correct   int
	Incorrect int
	Activity  float32
	Success   float32
	PC        string
}

type TagData struct {
	Description string
	Performance map[int]*ProblemPerformance
}

// -----------------------------------------------------------------------------------
func report_tagHandler(w http.ResponseWriter, r *http.Request) {
	// Check passcode
	if r.FormValue("pc") != Passcode {
		fmt.Fprintf(w, "Unauthorized")
		return
	}

	// Get tag_id from request
	tag_id := r.FormValue("tag_id")
	tagID, err := strconv.Atoi(tag_id)
	if err != nil {
		fmt.Println("Error converting tag_id to int:", err)
		return
	}

	// Get the tag description using GORM
	tagDescription, err := GetTagDescriptionByID(tagID)
	if err != nil {
		fmt.Println("Error retrieving tag description:", err)
		return
	}
	fmt.Println("Tag Description:", tagDescription)

	// Get problem performance data by tag ID
	record, err := GetProblemPerformanceByTagID(tagID)
	if err != nil {
		fmt.Println("Error retrieving problem performance:", err)
		return
	}

	// Get the student count
	studentCount, err := GetStudentCount()
	if err != nil {
		fmt.Println("Error retrieving student count:", err)
		return
	}

	// Calculate success and activity for each problem
	for pid, _ := range record {
		record[pid].Success = float32(record[pid].Correct) / float32(record[pid].Correct+record[pid].Incorrect)
		record[pid].Activity = record[pid].Activity / studentCount
	}

	// Render the template with the tag description and performance data
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.New("").Parse(TAG_REPORT_TEMPLATE)
	err = t.Execute(w, &TagData{Description: tagDescription, Performance: record})
	if err != nil {
		fmt.Println(err)
	}
}

//-----------------------------------------------------------------------------------

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
