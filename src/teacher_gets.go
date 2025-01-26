// Author: Vinhthuy Phan, 2018
package main

import (
	"encoding/json"
	"fmt"
	"github.com/GPTA/src/models"
	"net/http"
	"strconv"
)

// -----------------------------------------------------------------------------------
// Return a submission by index or priority
// -----------------------------------------------------------------------------------
func teacher_getsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	index, _ := strconv.Atoi(r.FormValue("index"))
	priority, _ := strconv.Atoi(r.FormValue("priority"))

	models.BoardsSem.Lock()
	defer models.BoardsSem.Unlock()

	selected := &models.Submission{}

	if index >= 0 {
		// Try to select by index first
		selected = models.WorkingSubs[index]
		models.WorkingSubs = append(models.WorkingSubs[:index], models.WorkingSubs[index+1:]...)
	} else if priority > 0 {
		// Try to select by priority: 1 (got it), 2 (help me)
		for i := 0; i < len(models.WorkingSubs); i++ {
			if models.WorkingSubs[i].Priority == priority {
				selected = models.WorkingSubs[i]
				models.WorkingSubs = append(models.WorkingSubs[:i], models.WorkingSubs[i+1:]...)
				// Students[selected.Uid].SubmissionStatus = 1
				subStat := &models.StudentSubmissionStatus{
					Filename:      selected.Filename,
					AttemptNumber: selected.AttemptNumber,
					Status:        1,
				}
				models.Students[selected.Uid].SubmissionStatus = append(models.Students[selected.Uid].SubmissionStatus, subStat)
				break
			}
		}
	} else {
		// Try to select the first highest priority
		first_sub_w_priority := []int{-1, -1, -1}
		for i := 0; i < len(models.WorkingSubs); i++ {
			p := models.WorkingSubs[i].Priority
			if first_sub_w_priority[p] == -1 {
				first_sub_w_priority[p] = i
			}
		}
		for i := len(first_sub_w_priority) - 1; i > 0; i-- {
			if first_sub_w_priority[i] != -1 {
				j := first_sub_w_priority[i]
				selected = models.WorkingSubs[j]
				models.WorkingSubs = append(models.WorkingSubs[:j], models.WorkingSubs[j+1:]...)
				// Students[selected.Uid].SubmissionStatus = 1
				subStat := &models.StudentSubmissionStatus{
					Filename:      selected.Filename,
					AttemptNumber: selected.AttemptNumber,
					Status:        1,
				}
				models.Students[selected.Uid].SubmissionStatus = append(models.Students[selected.Uid].SubmissionStatus, subStat)
				break
			}
		}
	}
	js, err := json.Marshal(selected)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// -----------------------------------------------------------------------------------
func teacher_gets_queueHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	js, err := json.Marshal(models.WorkingSubs)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//-----------------------------------------------------------------------------------
