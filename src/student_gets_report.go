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
	// Define a slice to hold the reports
	var report []StudentReport

	// Use GORM to find the scores for the given student, including the related problem data
	err := DB.Table("scores").
		Joins("join problems on problems.id = scores.problem_id").
		Where("student_id = ?", uid).
		Select("scores.score, scores.score_given_at, problems.filename").
		Find(&report).Error

	if err != nil {
		log.Fatal(err)
	}

	// Marshal the report slice into JSON
	js, _ := json.Marshal(report)

	// Set the content type and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//-----------------------------------------------------------------------------------
