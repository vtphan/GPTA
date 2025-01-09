// Author: Vinhthuy Phan, 2018
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// -----------------------------------------------------------------------------------
func student_getsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	var js []byte
	var err error

	BoardsSem.Lock()
	defer BoardsSem.Unlock()

	if _, ok := Students[uid]; ok {
		js, err = json.Marshal(Students[uid].Boards)
		for _, b := range Students[uid].Boards {
			if b.Pid != 0 {
				// Update the student status to "Working"
				addOrUpdateStudentStatus(uid, b.Pid, "Working", "", "", "")

				// Increment the active count for the problem
				err = IncProblemStatActive(b.Pid)
				if err != nil {
					log.Printf("Error incrementing active problem count for problem ID %d: %v", b.Pid, err)
				}
			}
		}
		Students[uid].Boards = []*Board{}
		if err == nil {
			// fmt.Println(string(js))
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	}
	fmt.Println(err.Error())
	js, err = json.Marshal([]*Board{})
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//-----------------------------------------------------------------------------------
