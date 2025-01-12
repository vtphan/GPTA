package main

import (
	"encoding/csv"
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

	// Get the problems and associated scores using GORM.
	var scores []struct {
		StudentID int
		Filename  string
		Score     int
	}
	if err := Database.Model(&Problem{}).
		Select("Score.StudentID, Problem.Filename, Score.Score").
		Joins("JOIN Scores Score ON Problem.ID = Score.ProblemID").
		Find(&scores).Error; err != nil {
		log.Fatal(err)
	}

	// Map data to organize scores by filename.
	data := make(map[string]map[int]int)
	for _, score := range scores {
		if data[score.Filename] == nil {
			data[score.Filename] = make(map[int]int)
		}
		data[score.Filename][score.StudentID] = score.Score
	}

	// Prepare CSV header.
	csvData := make([][]string, 1)
	csvData[0] = make([]string, len(data)+1)
	csvData[0][0] = "Name" // First column is for student names
	var i = 1
	for filename := range data {
		csvData[0][i] = filename
		i++
	}

	// Get all students.
	var students []Student
	if err := Database.Find(&students).Error; err != nil {
		log.Fatal(err)
	}

	// Prepare CSV rows.
	for _, student := range students {
		dt := make([]string, len(data)+1)
		dt[0] = student.Name // Assuming Name field in Student struct
		i = 1
		for _, v := range data {
			s, ok := v[student.ID]
			if ok {
				dt[i] = strconv.Itoa(s)
			} else {
				dt[i] = "0" // Default value if no score found.
			}
			i++
		}
		csvData = append(csvData, dt)
	}

	// Write to CSV file.
	file, err := os.Create("score.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(csvData); err != nil {
		log.Fatal(err)
	}

	// Send the file as a response.
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote("score.csv"))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, file.Name())
}
