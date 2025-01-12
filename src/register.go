package main

import (
	"bufio"
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
	var id int64

	if role == "teacher" {
		var teacherExists bool
		err := Database.Model(&Teacher{}).Select("COUNT(*) > 0").Where("name = ?", name).Find(&teacherExists).Error
		if err != nil {
			log.Fatal(err)
		}
		if teacherExists {
			fmt.Printf("%s already exists. Choose a different name.\n", name)
			return
		}
	} else {
		var studentExists bool
		err := Database.Model(&Student{}).Select("COUNT(*) > 0").Where("name = ?", name).Find(&studentExists).Error
		if err != nil {
			log.Fatal(err)
		}
		if studentExists {
			fmt.Printf("%s already exists. Choose a different name.\n", name)
			return
		}
	}

	if role != "teacher" {
		password = RandStringRunes(12)
	}

	if role == "teacher" {
		teacher := Teacher{Name: name, Password: password}
		err := Database.Create(&teacher).Error
		if err != nil {
			log.Fatal(err)
		}
		id = int64(teacher.ID)
		init_teacher(int(id), name, password)
	} else {
		student := Student{Name: name, Password: password}
		err := Database.Create(&student).Error
		if err != nil {
			log.Fatal(err)
		}
		id = int64(student.ID)
		init_student(int(id), name, password)
	}

	fmt.Printf("|%s| is added. Must complete registration.\n", name)
}

func complete_registrationHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	role := r.FormValue("role")
	course_id := r.FormValue("course_id")

	if course_id != Config.CourseId {
		fmt.Fprintf(w, "Failed")
		return
	}

	var id int
	var password string

	if role == "teacher" {
		var teacher Teacher
		err := Database.Model(&Teacher{}).Select("id, password").Where("name = ?", name).First(&teacher).Error
		if err != nil {
			fmt.Fprintf(w, "Failed")
			return
		}
		id = teacher.ID
		password = teacher.Password
	} else if role == "student" {
		var student Student
		err := Database.Model(&Student{}).Select("id, password").Where("name = ?", name).First(&student).Error
		if err != nil {
			fmt.Fprintf(w, "Failed")
			return
		}
		id = student.ID
		password = student.Password
	} else {
		fmt.Fprintf(w, "Failed")
		return
	}

	msg := fmt.Sprintf("%d,%s", id, password)
	fmt.Fprintf(w, msg)
}
