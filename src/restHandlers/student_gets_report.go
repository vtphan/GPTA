// Author: Vinhthuy Phan, 2018
package restHandlers

import (
	"encoding/json"
	"github.com/GPTA/src/repository"
	"log"
	"net/http"
)

// -----------------------------------------------------------------------------------
func StudentGetsReportHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	// Fetch report data
	report, err := repository.FetchStudentReport(uid)
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
