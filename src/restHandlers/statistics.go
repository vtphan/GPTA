// Author: Vinhthuy Phan, 2018
package restHandlers

import (
	"fmt"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"html/template"

	_ "github.com/mattn/go-sqlite3"

	// "math"
	"net/http"
	"strconv"
)

type StatsData struct {
	Performance        map[string]int
	ProblemDescription string
	Durations          map[string][]float64
	NextPid            int
	PrevPid            int
	Date               string
	PC                 string
	// Durations          map[string]float64
}

// -----------------------------------------------------------------------------------
func StatisticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("pc") != models.Passcode {
		fmt.Fprintf(w, "Unauthorized")
		return
	}
	pid, err := strconv.Atoi(r.FormValue("pid"))
	if err != nil {
		fmt.Println("Unknown problem")
		fmt.Fprintf(w, "Unknown problem")
		return
	}
	if pid <= 0 { // Select the last problem (max id)
		pid, err = repository.GetLatestProblemID()
		if err != nil {
			fmt.Println("Error retrieving latest problem:", err)
			return
		}
	}
	data := &StatsData{
		Performance: make(map[string]int),
		Durations:   make(map[string][]float64),
		PC:          models.Passcode,
		NextPid:     pid + 1,
		PrevPid:     pid - 1,
	}
	if pid > 0 {
		participants, probContent, probAt, performance, durations, err := repository.GetProblemStatistics(pid)
		if err != nil {
			fmt.Println("Error retrieving problem statistics", pid, err)
			return
		}

		// Extracting the date for attendance query
		theDate := probAt.Format("2006-01-02")

		// Fetching attendance data using the Attendance model
		attendants, err := repository.GetAttendanceByDate(theDate)
		if err != nil {
			fmt.Println("Error retrieving attendance for date", theDate, err)
			return
		}

		// Calculating inactive participants
		performance["Inactive"] = len(attendants) - len(participants)

		// Populating `data` (assuming it is a struct instance in the outer scope)
		data.ProblemDescription = probContent
		data.Performance = performance
		data.Durations = durations
		data.Date = theDate
	}

	w.Header().Set("Content-Type", "Text/html")
	t, err := template.New("").Parse(STATS_TEMPLATE)
	if err != nil {
		fmt.Println(err)
	} else {
		err = t.Execute(w, data)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// -----------------------------------------------------------------------------------
var STATS_TEMPLATE = `
<html>
  <head>
  <script src="https://cdn.plot.ly/plotly-latest.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/numeric/1.2.6/numeric.min.js"></script>
    <style>
    #main { width:1000px; margin: 0 auto;}
    #performance {
      width: 500px;
      height: 500px;
      float: left;
    }
    #durations {
      width: 500px;
      height: 500px;
      margin-left: 500px;
    }
    #pre{ width:100%; display:block;}
    .spacer{ width:100%; height:40px; margin: 0 auto;}
    .pager{ font-size:120%; Text-align: center; }
    .pager a{padding:25px; Text-decoration: none;}
    .pager a:visited{color:blue}
    </style>
  </head>
  <body>
    <div id="main">
    <div id="performance"></div>
	<div id="durations"></div>
	<script>
  	var perf = [];
	{{ range $key, $val := .Performance }}
		perf.push([{{$key}}, {{$val}}]);
	{{ end }}
	perf.sort();
	var values = [];
	var labels = [];
	for (i=0; i<perf.length; i++ ){
		values.push(perf[i][1]);
		labels.push(perf[i][0]);
	}
	var data = [{
	  values: values,
	  labels: labels,
	  type: 'pie'
	}];
	Plotly.newPlot('performance', data, {'title':'Points'});

	var data2 = [];
	{{ range $key, $val := .Durations }}
		data2.push({
			type:'violin',
			name: {{$key}},
			y: {{$val}},
			box: { visible: true },
			line: { color: 'blue' },
			meanline: { visible: true }
		});
	{{ end }}
	Plotly.newPlot('durations', data2, {
		title:'Time (min)',
		yaxis: {zeroline: false},
	});
    </script>

    <div class="spacer"></div>
    <pre style="padding-left:100px;">{{.Date}}
{{.ProblemDescription}}</pre>
    <div class="spacer"></div>
    <div class="pager">
    <a href="statistics?pc={{.PC}}&pid={{.PrevPid}}">Previous</a>
    <a href="statistics?pc={{.PC}}&pid={{.NextPid}}">Next</a>
    </div>
    </div>
  </body>
</html>
`

var STATS_TEMPLATE_OLD = `
<html>
  <head>
    <!--Load the AJAX API-->
    <script type="Text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="Text/javascript">
      google.charts.load('current', {'packages':['corechart']});
      google.charts.setOnLoadCallback(performancePieChart);
      google.charts.setOnLoadCallback(durationHistogram);

      function performancePieChart() {
      	var perf = [];
		{{ range $key, $val := .Performance }}
			perf.push([{{$key}}, {{$val}}]);
		{{ end }}
		perf.sort();
	    perf.unshift(['Category', 'Count']);
	    console.log(perf);
        var data = google.visualization.arrayToDataTable(perf);
        var options = {
          title: 'Performance'
        };

        var chart = new google.visualization.PieChart(document.getElementById('performance'));
        chart.draw(data, options);
      }

      function durationHistogram() {
        var data = google.visualization.arrayToDataTable([
          ['STID', 'Duration'],
          {{ range $key, $val := .Durations }}
          	[ {{$key}}, {{$val}} ],
          {{ end }}
        ]);

        var options = {
          title: 'Durations (minutes)',
          legend: { position: 'none' },
        };

        var chart = new google.visualization.Histogram(document.getElementById('durations'));
        chart.draw(data, options);
       }
    </script>
    <style>
    #main { width:1000px; margin: 0 auto;}
    #performance {
      width: 500px;
      height: 500px;
      float: left;
    }
    #durations {
      width: 500px;
      height: 500px;
      margin-left: 400px;
    }
    #pre{ width:100%; display:block;}
    .spacer{ width:100%; height:40px; margin: 0 auto;}
    .pager{ font-size:120%; Text-align: center; }
    .pager a{padding:25px; Text-decoration: none;}
    .pager a:visited{color:blue}
    </style>
  </head>
  <body>
    <div id="main">
    <div id="performance"></div>
    <div id="durations"></div>
    <div class="spacer"></div>
    <pre style="padding-left:100px;">{{.Date}}
{{.ProblemDescription}}</pre>
    <div class="spacer"></div>
    <div class="pager">
    <a href="statistics?pc={{.PC}}&pid={{.PrevPid}}">Previous</a>
    <a href="statistics?pc={{.PC}}&pid={{.NextPid}}">Next</a>
    </div>
    </div>
  </body>
</html>
`
