package main

import (
	"encoding/csv"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"log"
	"net/http"
	"os"
	"strconv"
)

func exportPointsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	if role == "student" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Fetch data from the database
	data, err := repository.FetchScores()
	if err != nil {
		log.Printf("Error fetching scores: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create CSV data
	csvData := make([][]string, 1)
	header := []string{"Student Name"}
	for filename := range data {
		header = append(header, filename)
	}
	csvData[0] = header

	for studentID := range models.Students {
		row := make([]string, len(data)+1)
		row[0] = getName(studentID, "student")
		i := 1
		for _, scores := range data {
			if score, ok := scores[studentID]; ok {
				row[i] = strconv.Itoa(score)
			} else {
				row[i] = ""
			}
			i++
		}
		csvData = append(csvData, row)
	}

	// Write CSV to a file
	file, err := os.Create("score.csv")
	if err != nil {
		log.Printf("Failed to create CSV file: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(csvData); err != nil {
		log.Printf("Failed to write CSV data: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Serve the CSV file as a response
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote("score.csv"))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, file.Name())
}
