package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// -----------------------------------------------------------------
func add_multiple(filename, role string) {
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
	var teacher Teacher
	var student Student
	var id int

	if role == "teacher" {
		teacher, err = FindTeacherByName(name)
		if err != nil && err.Error() != RecordNotFound {
			log.Fatal(err)
		}
	} else {
		student, err = FindStudentByName(name)
		if err != nil && err.Error() != RecordNotFound {
			log.Fatal(err)
		}
	}
	if err != nil && err.Error() == RecordNotFound {
		fmt.Printf("%s already exists. Choose a different name.\n", name)
		return
	}
	if role != "teacher" {
		password = RandStringRunes(12)
	}
	if role == "teacher" {
		teacher, err = AddTeacher(name, password)
	} else {
		student, err = AddStudent(name, password)
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
		init_teacher(id, name, password) //todo - map
	} else {
		init_student(id, name, password) //todo - map
	}
	fmt.Printf("|%s| is added. Must complete registeration.\n", name)
}

// -----------------------------------------------------------------
func complete_registrationHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	role := r.FormValue("role")
	course_id := r.FormValue("course_id")
	if course_id != Config.CourseId {
		fmt.Fprintf(w, "Failed")
		return
	}
	var err error
	var rows *sql.Rows
	if role == "teacher" {
		rows, err = Database.Query("select id, password from teacher where name=?", name)
		defer rows.Close()
	} else if role == "student" {
		rows, err = Database.Query("select id, password from student where name=?", name)
		defer rows.Close()
	} else {
		fmt.Fprintf(w, "Failed")
		return
	}
	if err != nil {
		fmt.Fprintf(w, "Failed")
		// log.Fatal(err)
	}
	var password string
	var id int
	for rows.Next() {
		rows.Scan(&id, &password)
		msg := fmt.Sprintf("%d,%s", id, password)
		// msg := ""
		// if Config.NameServer != "" {
		// 	msg = fmt.Sprintf("%d,%s,%s,%s", id, password, Config.CourseId, Config.NameServer)
		// } else {
		// 	msg = fmt.Sprintf("%d,%s,%s", id, password, Config.CourseId)
		// }
		fmt.Fprintf(w, msg)
		return
	}
	fmt.Fprintf(w, "Failed")
}

//-----------------------------------------------------------------
