// Author: Vinhthuy Phan, 2018
package restHandlers

import (
	"fmt"
	"github.com/GPTA/src/frontEnd"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"strconv"
)

// -----------------------------------------------------------------------------------

// -----------------------------------------------------------------------------------
func ReportHandler(w http.ResponseWriter, r *http.Request) {
	// Check passcode
	if r.FormValue("pc") != models.Passcode {
		fmt.Fprintf(w, "Unauthorized")
		return
	}

	// Fetch tags
	tags, err := repository.GetTags()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch student scores
	scores, err := repository.GetStudentScores()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Prepare the data for rendering
	record := &models.TagsViewData{
		Tags:            tags,
		Scores:          scores,
		SubmissionCount: make(map[string]int), // You can add logic to fill this if needed
		PC:              models.Passcode,
	}

	// Render the view
	w.Header().Set("Content-Type", "Text/html")
	t, _ := template.New("").Parse(frontEnd.TAGS_VIEW_TEMPLATE)
	err = t.Execute(w, record)
	if err != nil {
		fmt.Println(err)
	}
}

// -----------------------------------------------------------------------------------

// -----------------------------------------------------------------------------------
func ReportTagHandler(w http.ResponseWriter, r *http.Request) {
	// Check passcode
	if r.FormValue("pc") != models.Passcode {
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
	tagDescription, err := repository.GetTagDescriptionByID(tagID)
	if err != nil {
		fmt.Println("Error retrieving tag description:", err)
		return
	}
	fmt.Println("Tag Description:", tagDescription)

	// Get problem performance data by tag ID
	record, err := repository.GetProblemPerformanceByTagID(tagID)
	if err != nil {
		fmt.Println("Error retrieving problem performance:", err)
		return
	}

	// Get the student count
	studentCount, err := repository.GetStudentCount()
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
	w.Header().Set("Content-Type", "Text/html")
	t, _ := template.New("").Parse(frontEnd.TAG_REPORT_TEMPLATE)
	err = t.Execute(w, &models.TagData{Description: tagDescription, Performance: record})
	if err != nil {
		fmt.Println(err)
	}
}

//-----------------------------------------------------------------------------------
