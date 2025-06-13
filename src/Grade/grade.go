package Grade

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GPTA/src/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Grade struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	StudentID int       `gorm:"not null"`
	ProblemID int       `gorm:"not null"`
	Grade     string    `gorm:"type:varchar(50)"`
	GradedAt  time.Time `gorm:"autoCreateTime"`
}

type GradeRequest struct {
	StudentID int    `json:"student_id"`
	Grade     string `json:"grade"`
	ProblemID int    `json:"problem_id"`
}

func HandleGradeSubmission(w http.ResponseWriter, r *http.Request) {
	// Allow only POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the request
	var req GradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("🎯 Grade received - Student ID: %d | Grade: %s | Problem ID: %d\n",
		req.StudentID, req.Grade, req.ProblemID)

	err := GradeStudent(req)
	if err != nil {
		http.Error(w, "Failed to create grade", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("🔄 Grade updated successfully"))
}

func GradeStudent(req GradeRequest) error {
	var grade Grade
	err := models.DB.Where("student_id = ? AND problem_id = ?", req.StudentID, req.ProblemID).First(&grade).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Grade doesn't exist — create new
		grade = Grade{
			StudentID: req.StudentID,
			ProblemID: req.ProblemID,
			Grade:     req.Grade,
		}
		if err = models.DB.Create(&grade).Error; err != nil {
			return err
		}
		return err
	} else if err != nil {
		return err
	}

	// Grade exists — update it
	grade.Grade = req.Grade
	if err = models.DB.Save(&grade).Error; err != nil {
		return err
	}
	return err
}

func DeleteGradeByStudentAndProblem(studentID int, problemID int) error {
	var grade Grade

	// Try to find the grade first
	err := models.DB.Where("student_id = ? AND problem_id = ?", studentID, problemID).First(&grade).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Nothing to delete
		return nil
	} else if err != nil {
		// Some other error occurred
		return err
	}

	// Grade found — proceed with deletion
	if err := models.DB.Delete(&grade).Error; err != nil {
		return err
	}

	return nil
}
