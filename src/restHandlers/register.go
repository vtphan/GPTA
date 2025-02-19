package restHandlers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strings"
)

// -----------------------------------------------------------------
// Function to create a course

func AddCourse(courseID, courseName string) error {
	// Check if the course already exists
	var existing models.Course
	err := models.DB.Where("course_id = ?", courseID).First(&existing).Error
	if err == nil {
		fmt.Printf("Course %s (%s) already exists.\n", courseID, courseName)
		return errors.New("Course already exists")
	}

	// If not found, insert the new course
	course := models.Course{
		CourseID:   courseID,
		CourseName: courseName,
	}
	if err := models.DB.Create(&course).Error; err != nil {
		return fmt.Errorf("failed to insert course: %w", err)
	}
	fmt.Printf("Course %s (%s) added successfully.\n", courseName, courseID)
	return nil
}

func AddMultipleCourses(courseFile string) {
	file, err := os.Open(courseFile)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			parts := strings.Split(line, ",")
			courseID := parts[0]
			courseName := parts[1]
			_ = AddCourse(courseID, courseName)
		}
	}
}

// Function to associate students with courses
func AddStudentToCourse(studentID int, courseID string) error {
	// Ensure that the student is not added to the same course more than once
	var course models.Course
	err := models.DB.Where("course_id = ?", courseID).First(&course).Error
	if err != nil {
		fmt.Printf("Course %d does not exist. Skipping enrollment for student %d.\n", courseID, studentID)
		return nil // Skip if the course does not exist
	}

	var existing models.StudentClass
	err = models.DB.Where("student_id = ? AND course_id = ?", studentID, courseID).First(&existing).Error
	if err == nil {
		fmt.Printf("Student %d is already assigned to course %s.\n", studentID, courseID)
		return nil
	}

	studentClass := models.StudentClass{
		StudentID: studentID,
		CourseID:  courseID,
	}
	if err := models.DB.Create(&studentClass).Error; err != nil {
		return fmt.Errorf("failed to enroll student in course: %w", err)
	}
	return nil
}

// Function to associate teachers with courses
func AddTeacherToCourse(teacherID int, courseID string) error {

	var course models.Course
	err := models.DB.Where("course_id = ?", courseID).First(&course).Error
	if err != nil {
		fmt.Printf("Course %d does not exist. Skipping enrollment for student %d.\n", courseID, teacherID)
		return nil // Skip if the course does not exist
	}
	// Ensure that the teacher is not added to the same course more than once
	var existing models.TeacherClass
	err = models.DB.Where("teacher_id = ? AND course_id = ?", teacherID, courseID).First(&existing).Error
	if err == nil {
		fmt.Printf("Teacher %d is already assigned to course %s.\n", teacherID, courseID)
		return nil
	}

	teacherClass := models.TeacherClass{
		TeacherID: teacherID,
		CourseID:  courseID,
	}
	if err := models.DB.Create(&teacherClass).Error; err != nil {
		return fmt.Errorf("failed to assign teacher to course: %w", err)
	}
	return nil
}

// Updated AddMultiple function to process courses
func AddMultiple(filename, role string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			parts := strings.Split(line, ",")
			name := parts[0]
			name = strings.TrimSpace(name)

			// Handle teacher
			if role == "teacher" {
				if len(parts) < 3 { // Ensure at least name, password, and one course
					log.Fatal("Invalid teacher data format")
				}
				password := parts[1]
				courses := parts[2:] // Courses start from index 2

				for _, courseID := range courses {
					AddTeacher(name, password, courseID)
				}
			}

			// Handle student
			if role == "student" {
				courses := parts[1:] // Course IDs start from index 1
				for _, courseID := range courses {
					AddStudent(name, courseID)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Updated AddStudent function to include courses
func AddStudent(name, courseID string) error {
	var student models.Student
	err := models.DB.Where("name = ?", name).First(&student).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		// Create student
		password := models.RandStringRunes(12)
		student, err = repository.AddStudent(name, password)
		if err != nil {
			return err
		}
	}

	// Associate student with courses
	courseIDInt := strings.TrimSpace(courseID)
	if err != nil {
		return fmt.Errorf("Invalid course ID")
	}
	err = AddStudentToCourse(student.ID, courseIDInt)
	if err != nil {
		return err
	}
	return nil
}

func AddTeacher(name, password, courseID string) (int, error) {
	var teacher models.Teacher
	password = strings.TrimSpace(password)
	err := models.DB.Where("name = ?", name).First(&teacher).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	if err == gorm.ErrRecordNotFound {
		// Create teacher
		teacher, err = repository.AddTeacher(name, password)
		if err != nil {
			return 0, err
		}
	}

	// Associate teacher with courses
	courseIDInt := strings.TrimSpace(courseID)
	if err != nil {
		log.Fatal("Invalid course ID")
	}
	err = AddTeacherToCourse(teacher.ID, courseIDInt)
	if err != nil {
		return 0, err
	}
	return teacher.ID, nil
}

// -----------------------------------------------------------------
func CompleteRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	role := r.FormValue("role")
	course_id := r.FormValue("course_id")

	var msg string

	if role == "teacher" {
		teacher, err := repository.GetTeacherByName(name)
		if err != nil {
			fmt.Fprintf(w, "Failed")
			return
		}
		// Check if the teacher is mapped to the given course ID
		if teacherCourses, exists := models.TeacherClassesMap[teacher.ID]; !exists {
			// Teacher is not mapped to any courses
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else {
			// Check if the course_id exists in the teacherCourses slice
			isAuthorized := false
			for _, course := range teacherCourses {
				if course == course_id {
					isAuthorized = true
					break
				}
			}

			if !isAuthorized {
				// Teacher is not authorized for this course
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}
		msg = fmt.Sprintf("%d,%s", teacher.ID, teacher.Password)
	} else if role == "student" {
		student, err := repository.GetStudentByName(name)
		if err != nil {
			fmt.Fprintf(w, "Failed")
			return
		}
		if studentCourses, exists := models.StudentClassesMap[student.ID]; !exists {
			// Student is not mapped to any courses
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else {
			// Check if the course_id exists in the studentCourses slice
			isAuthorized := false
			for _, course := range studentCourses {
				if course == course_id {
					isAuthorized = true
					break
				}
			}

			if !isAuthorized {
				// Student is not authorized for this course
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}
		msg = fmt.Sprintf("%d,%s", student.ID, student.Password)
	} else {
		fmt.Fprintf(w, "Failed")
		return
	}

	// Respond with the appropriate message
	fmt.Fprintf(w, msg)
}

func GetCoursesHandler(w http.ResponseWriter, r *http.Request) {
	var courses []models.Course
	if err := models.DB.Find(&courses).Error; err != nil {
		http.Error(w, "Failed to fetch courses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(courses); err != nil {
		http.Error(w, "Failed to encode courses", http.StatusInternalServerError)
	}
}

func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name     string `json:"student_name"`
		CourseID string `json:"course_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Trim input values
	req.Name = strings.TrimSpace(req.Name)
	req.CourseID = strings.TrimSpace(req.CourseID)

	// Validate required fields
	if req.Name == "" || req.CourseID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Call function to add student
	err := AddStudent(req.Name, req.CourseID)
	if err != nil {
		http.Error(w, "Failed to add student to course", http.StatusInternalServerError)
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Student added successfully"}
	json.NewEncoder(w).Encode(response) // Ensure JSON encoding
}

type AddTeacherRequest struct {
	Name     string `json:"teacher_name"`
	Password string `json:"teacher_pass"`
	CourseID string `json:"course_id"`
}

func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req AddTeacherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Trim input values
	req.Name = strings.TrimSpace(req.Name)
	req.Password = strings.TrimSpace(req.Password)
	req.CourseID = strings.TrimSpace(req.CourseID)

	// Validate required fields
	if req.Name == "" || req.Password == "" || req.CourseID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	id, err := AddTeacher(req.Name, req.Password, req.CourseID)
	if err != nil {
		http.Error(w, "Failed to add teacher to course", http.StatusInternalServerError)
		return
	}

	models.TeacherMap[id] = req.Password
	models.TeacherPass[req.Name] = req.Password
	models.TeacherNameToId[req.Name] = id
	models.TeacherIdToName[id] = req.Name
	models.TeacherClassesMap[id] = append(models.TeacherClassesMap[id], req.CourseID)
	// Success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Teacher added to course successfully"})
}

type AddCourseRequest struct {
	CourseID   string `json:"course_id"`
	CourseName string `json:"course_name"`
}

func AddCourseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req AddCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.CourseID == "" || req.CourseName == "" {
		http.Error(w, "Course ID and Course Name are required", http.StatusBadRequest)
		return
	}

	err := AddCourse(req.CourseID, req.CourseName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding course: %v", err), http.StatusInternalServerError)
		return
	}

	// Set content type and send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Course added successfully"})
}

//-----------------------------------------------------------------
