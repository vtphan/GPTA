package restHandlers

import (
	"bufio"
	"fmt"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"log"
	"net/http"
	"os"
	"strings"
)

// -----------------------------------------------------------------
func AddMultiple(filename, role string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	password := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" {
			if role == "teacher" {
				scanner.Scan()
				password = strings.TrimSpace(scanner.Text())
				if password == "" {
					fmt.Printf("Password can not be empty!")
					log.Fatal("Password can not be empty!")
				}
			}
			add_user(name, role, password)
		}
	}
}

// -----------------------------------------------------------------
func add_user(name, role string, password string) {
	var err error
	var teacher models.Teacher
	var student models.Student
	var id int

	if role == "teacher" {
		teacher, err = repository.FindTeacherByName(name)
		if err != nil && err.Error() != models.RecordNotFound {
			log.Fatal(err)
		}
	} else {
		student, err = repository.FindStudentByName(name)
		if err != nil && err.Error() != models.RecordNotFound {
			log.Fatal(err)
		}
	}
	if err == nil {
		fmt.Printf("%s already exists. Choose a different name.\n", name)
		return
	}
	if err.Error() == models.RecordNotFound {

		if role != "teacher" {
			password = models.RandStringRunes(12)
		}
		if role == "teacher" {
			teacher, err = repository.AddTeacher(name, password)
		} else {
			student, err = repository.AddStudent(name, password)
		}
		if err != nil {
			log.Fatal(err)
		}
		if role == "teacher" {
			id = teacher.ID
		} else {
			id = student.ID
		}
		if role == "teacher" {
			repository.InitTeacher(id, name, password)
		} else {
			repository.InitStudent(id, name, password)
		}
	}
	fmt.Printf("|%s| is added. Must complete registeration.\n", name)
}

// -----------------------------------------------------------------
func CompleteRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	role := r.FormValue("role")
	course_id := r.FormValue("course_id")

	// Check course ID
	if course_id != models.Config.CourseId {
		fmt.Fprintf(w, "Failed")
		return
	}

	var msg string

	if role == "teacher" {
		teacher, err := repository.GetTeacherByName(name)
		if err != nil {
			fmt.Fprintf(w, "Failed")
			return
		}
		msg = fmt.Sprintf("%d,%s", teacher.ID, teacher.Password)
	} else if role == "student" {
		student, err := repository.GetStudentByName(name)
		if err != nil {
			fmt.Fprintf(w, "Failed")
			return
		}
		msg = fmt.Sprintf("%d,%s", student.ID, student.Password)
	} else {
		fmt.Fprintf(w, "Failed")
		return
	}

	// Respond with the appropriate message
	fmt.Fprintf(w, msg)
}

//-----------------------------------------------------------------
