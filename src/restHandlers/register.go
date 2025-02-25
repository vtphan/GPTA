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
	"strconv"
	"strings"
)

// -----------------------------------------------------------------
// Function to create a course

func AddCourse(courseID string) error {
	// Check if the course already exists
	var existing models.Course
	err := models.DB.Where("course_id = ?", courseID).First(&existing).Error
	if err == nil {
		fmt.Printf("Course %s (%s) already exists.\n", courseID)
		return errors.New("Course already exists")
	}

	// If not found, insert the new course
	course := models.Course{
		CourseID: courseID,
	}
	if err := models.DB.Create(&course).Error; err != nil {
		return fmt.Errorf("failed to insert course: %w", err)
	}
	fmt.Printf("Course %s (%s) added successfully.\n", courseID)
	return nil
}

func AddMultipleCourses(courseFile string) {
	file, err := os.Open(courseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		courseID := strings.TrimSpace(scanner.Text())
		if courseID != "" {
			_ = AddCourse(courseID) // Now only passing courseID
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
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
	var currentCourseID string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Detect course ID using "COURSE:" prefix
		if strings.HasPrefix(line, "COURSE:") {
			currentCourseID = strings.TrimSpace(strings.TrimPrefix(line, "COURSE:"))
			continue
		}

		// Process based on role
		if role == "student" {
			AddStudent(line, currentCourseID)
		} else if role == "teacher" {
			parts := strings.Split(line, ",")
			if len(parts) < 2 {
				log.Fatal("Invalid teacher data format")
			}
			name := strings.TrimSpace(parts[0])
			password := strings.TrimSpace(parts[1])
			AddTeacher(name, password, currentCourseID)
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
	courseIDStr := strings.TrimSpace(courseID)
	if err != nil {
		log.Fatal("Invalid course ID")
	}
	err = AddTeacherToCourse(teacher.ID, courseIDStr)
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
	teacherID := r.URL.Query().Get("teacher_id")
	if teacherID == "" {
		http.Error(w, "Missing teacher_id", http.StatusBadRequest)
		return
	}

	var courses []models.Course
	if err := models.DB.
		Joins("JOIN teacher_classes ON teacher_classes.course_id = courses.course_id").
		Where("teacher_classes.teacher_id = ?", teacherID).
		Find(&courses).Error; err != nil {
		http.Error(w, "Failed to fetch courses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(courses); err != nil {
		http.Error(w, "Failed to encode courses", http.StatusInternalServerError)
	}
}

func AddStudentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Names    []string `json:"student_names"`
		CourseID string   `json:"course_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Trim input values and validate
	req.CourseID = strings.TrimSpace(req.CourseID)
	if req.CourseID == "" || len(req.Names) == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	for i := range req.Names {
		req.Names[i] = strings.TrimSpace(req.Names[i])
	}

	// Add each student to the course
	for _, name := range req.Names {
		if name == "" {
			continue
		}
		if err := AddStudent(name, req.CourseID); err != nil {
			http.Error(w, fmt.Sprintf("Failed to add student %s", name), http.StatusInternalServerError)
			return
		}
	}

	// Refresh student-course map
	studentClasses, err := repository.GetAllStudentClasses()
	if err != nil {
		http.Error(w, "Failed to refresh student-course map", http.StatusInternalServerError)
		return
	}
	models.StudentClassesMap = make(map[int][]string)
	for _, class := range studentClasses {
		models.StudentClassesMap[class.StudentID] = append(models.StudentClassesMap[class.StudentID], class.CourseID)
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "All students added successfully"}
	json.NewEncoder(w).Encode(response)
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
	if req.Name == "" || req.Password == "" {
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
	if req.CourseID != "" {
		models.TeacherClassesMap[id] = append(models.TeacherClassesMap[id], req.CourseID)
	}
	// Success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Teacher added to course successfully"})
}

type AddCourseRequest struct {
	CourseID  string `json:"course_id"`
	TeacherID string `json:"teacher_id"`
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

	req.CourseID = strings.TrimSpace(req.CourseID)
	req.TeacherID = strings.TrimSpace(req.TeacherID)

	if req.CourseID == "" || req.TeacherID == "" {
		http.Error(w, "Course ID and Teacher ID are required", http.StatusBadRequest)
		return
	}

	err := AddCourse(req.CourseID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding course: %v", err), http.StatusInternalServerError)
		return
	}

	teacherID, err := strconv.Atoi(req.TeacherID)
	if err != nil {
		http.Error(w, "Invalid Teacher ID format, must be an integer", http.StatusBadRequest)
		return
	}

	err = AddTeacherToCourse(teacherID, req.CourseID) // Use TeacherID from request
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding teacher to course: %v", err), http.StatusInternalServerError)
		return
	}

	// Set content type and send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Course added successfully"})
}

//-----------------------------------------------------------------
