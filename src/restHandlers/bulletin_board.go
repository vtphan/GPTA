// Author: Vinhthuy Phan, 2018
package restHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/GPTA/src/frontEnd"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------------
type BulletinBoardMessage struct {
	Code           string
	I              int
	NextI          int
	PrevI          int
	PC             string
	P1             int
	P2             int
	P1Graded       int
	P1Ungraded     int
	P2Answered     int
	P2Unanswered   int
	ActiveProblems string
	// ActiveProblems int
	BulletinItems int
	AnswerCount   int
	Attendance    int
	Address       string
	Authenticated bool
}

// -----------------------------------------------------------------------------------
func TeacherAddsBulletinPageHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	models.BulletinSem.Lock()
	defer models.BulletinSem.Unlock()
	models.BulletinBoard = append(models.BulletinBoard, r.FormValue("content"))
	fmt.Fprintf(w, "Content added to bulletin board")
}

// -----------------------------------------------------------------
func RemoveBulletinPageHandler(w http.ResponseWriter, r *http.Request) {
	models.BulletinSem.Lock()
	defer models.BulletinSem.Unlock()
	i, _ := strconv.Atoi(r.FormValue("i"))
	passcode := r.FormValue("pc")
	if passcode == models.Passcode && i >= 0 && i < len(models.BulletinBoard) {
		models.BulletinBoard = append(models.BulletinBoard[:i], models.BulletinBoard[i+1:]...)
		http.Redirect(w, r, "view_bulletin_board?i=0&pc="+passcode, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "view_bulletin_board?i="+r.FormValue("i")+"&pc="+passcode, http.StatusSeeOther)
	}
}

// -----------------------------------------------------------------------------------
func GetBulletinBoardData(i int, passcode string) *BulletinBoardMessage {
	models.BulletinSem.Lock()
	defer models.BulletinSem.Unlock()

	if i >= len(models.BulletinBoard) {
		i = 0
	}
	// Get code and build page links
	code := ""
	if i >= 0 && i < len(models.BulletinBoard) {
		code = models.BulletinBoard[i]
	}

	// Get priority counts
	priority := []int{0, 0, 0}
	for j := 0; j < len(models.WorkingSubs); j++ {
		priority[models.WorkingSubs[j].Priority]++
	}
	next_i, prev_i := 0, 0
	if len(models.BulletinBoard) > 0 {
		next_i = (i + 1 + len(models.BulletinBoard)) % len(models.BulletinBoard)
		prev_i = (i - 1 + len(models.BulletinBoard)) % len(models.BulletinBoard)
	}
	answers := 0
	keys := make([]string, 0)
	for key := range models.ActiveProblems {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	submissions := make([]string, 0)
	for i, key := range keys {
		p := models.ActiveProblems[key]
		if p.Active {
			starting_time, err := repository.GetProblemUploadTime(p.Info.Pid)
			if err != nil {
				fmt.Println("Error retrieving problem starting time", err)
				return &BulletinBoardMessage{}
			}
			duration := time.Since(starting_time).Minutes()
			subs := len(p.Attempts)
			label := fmt.Sprintf("P%d: %d subs after %.0fm", i+1, subs, duration)
			submissions = append(submissions, label)
			answers += len(p.Answers)
		}
	}
	active_problems := strings.Join(submissions, ". ")
	fmt.Println(">", active_problems)

	data := &BulletinBoardMessage{
		Code:           code,
		I:              i,
		NextI:          next_i,
		PrevI:          prev_i,
		PC:             passcode,
		P1:             len(models.Submissions),
		P1Graded:       len(models.Submissions) - len(models.WorkingSubs),
		P1Ungraded:     len(models.WorkingSubs),
		P2:             len(models.HelpSubmissions),
		P2Answered:     len(models.HelpSubmissions) - len(models.WorkingHelpSubs),
		P2Unanswered:   len(models.WorkingHelpSubs),
		ActiveProblems: active_problems,
		BulletinItems:  len(models.BulletinBoard),
		AnswerCount:    answers,
		Attendance:     len(models.Students),
		Address:        models.Config.Address,
		Authenticated:  passcode == models.Passcode,
	}
	return data
}

// -----------------------------------------------------------------------------------
func BulletinBoardDataHandler(w http.ResponseWriter, r *http.Request) {
	data := GetBulletinBoardData(0, "")
	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

// -----------------------------------------------------------------------------------
func ViewBulletinBoardHandler(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.FormValue("i"))
	passcode := r.FormValue("pc")
	if err != nil {
		i = 0
	}

	temp := template.New("")
	t, err2 := temp.Parse(frontEnd.TEACHER_MESSAGING_TEMPLATE)
	if err2 != nil {
		log.Fatal(err2)
	}
	data := GetBulletinBoardData(i, passcode)
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

//-----------------------------------------------------------------------------------
// func student_messagesHandler(w http.ResponseWriter, r *http.Request) {
// 	stid, err := strconv.Atoi(r.FormValue("stid"))
// 	if err != nil {
// 		fmt.Fprintf(w, "Error")
// 	}
// 	_, ok := Students[stid]
// 	if ok {
// 		t := template.New("")
// 		t, err := t.Parse(STUDENT_MESSAGING_TEMPLATE)
// 		if err == nil {
// 			data := struct{ Message string }{Students[stid].Status}
// 			w.Header().Set("Content-Type", "text/html")
// 			t.Execute(w, data)
// 		} else {
// 			fmt.Println(err)
// 		}
// 	} else {
// 		fmt.Fprint(w, "Error")
// 	}
// }
