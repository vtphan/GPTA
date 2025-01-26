// Author: Vinhthuy Phan, 2018
package main

import (
	"fmt"
	"github.com/GPTA/src/repository"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"sort"
	"strconv"
)

// -----------------------------------------------------------------------------------
type SubmissionData struct {
	Flag      string
	Start     int64
	At        int64
	Completed int64
}

// -----------------------------------------------------------------------------------
func analyze_submissionsHandler(w http.ResponseWriter, r *http.Request) {
	pidStr := r.FormValue("pid")
	pid, err := strconv.Atoi(pidStr) // Convert string to int
	if err != nil {
		fmt.Println("Error converting pid to int:", err)
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	// Use GetProblemUploadTime to retrieve the problem upload time
	start, err := repository.GetProblemUploadTime(pid)
	if err != nil {
		fmt.Println("Error retrieving problem upload time:", err)
		http.Error(w, "Failed to retrieve problem upload time", http.StatusInternalServerError)
		return
	}

	records := make(map[int][]*SubmissionData)

	// Use GetSubmissionsByProblemID to retrieve the submissions for the problem
	submissions, err := repository.GetSubmissionsByProblemID(pid)
	if err != nil {
		fmt.Println("Error retrieving submissions:", err)
		http.Error(w, "Failed to retrieve submissions", http.StatusInternalServerError)
		return
	}

	for _, submission := range submissions {
		sid := submission.StudentID
		priority := submission.SubmissionCategory
		at := submission.CodeSubmittedAt
		completed := submission.Completed

		if _, ok := records[sid]; !ok {
			records[sid] = make([]*SubmissionData, 0)
		}

		flag := "unknown"
		if priority == 1 {
			flag = "Got it!"
		} else if priority == 2 {
			flag = "Help!"
		}

		records[sid] = append(
			records[sid],
			&SubmissionData{
				Flag:      flag,
				Start:     start.UnixNano(),
				At:        at.UnixNano(),
				Completed: completed.UnixNano(),
			})
	}

	// Sort submissions by submission time (at)
	for sid, _ := range records {
		sort.Slice(records[sid], func(i, j int) bool {
			return records[sid][i].At < records[sid][j].At
		})
	}

	w.Header().Set("Content-Type", "text/html")
	t, _ := template.New("").Parse(ANALYZE_SUBMISSIONS_TEMPLATE)
	err = t.Execute(w, records)
	if err != nil {
		fmt.Println(err)
	}
}

// -----------------------------------------------------------------------------------
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
