// Author: Vinhthuy Phan, 2018
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type StudentReport struct {
	Points   int
	Filename string
	Date     int64
}

// -----------------------------------------------------------------------------------
func student_gets_reportHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	// Fetch report data
	report, err := FetchStudentReport(uid)
	if err != nil {
		log.Printf("Error fetching student report: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert to JSON
	js, err := json.Marshal(report)
	if err != nil {
		log.Printf("Error marshaling JSON: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//-----------------------------------------------------------------------------------
