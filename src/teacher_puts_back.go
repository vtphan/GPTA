// Author: Vinhthuy Phan, 2018
package main

import (
	"fmt"
	"github.com/GPTA/src/models"
	"net/http"
	"strconv"
)

// -----------------------------------------------------------------------------------
func teacher_puts_backHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	sid, _ := strconv.Atoi(r.FormValue("sid"))
	models.SubSem.Lock()
	defer models.SubSem.Unlock()
	if _, ok := models.Submissions[sid]; ok {
		models.WorkingSubs = append(models.WorkingSubs, models.Submissions[sid])
		fmt.Fprintf(w, "Submission has been put back into the queue.")
	} else {
		fmt.Fprintf(w, "Unknown submission.")
	}
}

// //-----------------------------------------------------------------------------------
