// Author: Vinhthuy Phan, 2018
package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"sort"
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
	// if r.FormValue("pc") != Passcode {
	//     fmt.Fprintf(w, "Unauthorized")
	//     return
	// }

	pid := r.FormValue("pid")
	if pid == "" {
		http.Error(w, "Problem ID is required", http.StatusBadRequest)
		return
	}

	var problem Problem
	if err := DB.First(&problem, pid).Error; err != nil {
		http.Error(w, "Problem not found", http.StatusNotFound)
		fmt.Println(err)
		return
	}

	var submissions []Submission
	if err := DB.Where("problem_id = ?", pid).Find(&submissions).Error; err != nil {
		http.Error(w, "Failed to fetch submissions", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	records := make(map[int][]*SubmissionData)
	for _, submission := range submissions {
		if _, ok := records[submission.StudentID]; !ok {
			records[submission.StudentID] = make([]*SubmissionData, 0)
		}

		flag := "unknown"
		if submission.SubmissionCategory == 1 {
			flag = "Got it!"
		} else if submission.SubmissionCategory == 2 {
			flag = "Help!"
		}

		records[submission.StudentID] = append(records[submission.StudentID], &SubmissionData{
			Flag:      flag,
			Start:     problem.ProblemUploadedAt.UnixNano(),
			At:        submission.CodeSubmittedAt.UnixNano(),
			Completed: submission.Completed.UnixNano(),
		})
	}

	// Sort submissions for each student by submission time
	for sid := range records {
		sort.Slice(records[sid], func(i, j int) bool {
			return records[sid][i].At < records[sid][j].At
		})
	}

	// Render the HTML template
	w.Header().Set("Content-Type", "text/html")
	t, err := template.New("").Parse(ANALYZE_SUBMISSIONS_TEMPLATE)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := t.Execute(w, records); err != nil {
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
