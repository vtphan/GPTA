// Author: Vinhthuy Phan, 2018
package restHandlers

import (
	"fmt"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"log"
	"net/http"
	"time"
)

// -----------------------------------------------------------------------------------
func StudentAskHelpHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, filename := r.FormValue("content"), r.FormValue("filename")
	need_help_with := r.FormValue("need_help_with")
	sid := int64(0)
	if need_help_with == "" {
		need_help_with = "None."
	}

	msg := "your help message has been sent"

	pid := 0
	prob, ok := models.ActiveProblems[filename]
	snapshotID := 0
	if ok {
		if !prob.Active {
			msg = "Problem is no longer active. But the teacher will look at your submission."
		} else {
			pid = prob.Info.Pid
			if _, ok := prob.Attempts[uid]; !ok {
				models.ActiveProblems[filename].Attempts[uid] = prob.Info.Attempts
			}
			now := time.Now()
			snapshotID = AddCodeSnapshot(uid, pid, content, 0, now, "at_ask_for_help")

			// result, err = AddHelpSubmissionSQL.Exec(pid, uid, snapshotID, "", need_help_with, now)
			result, err := repository.AddMessage(snapshotID, need_help_with, uid, "student", now, 0)
			if err != nil {
				log.Fatal(err)
			}
			sid = int64(result)
			err = repository.IncrementProblemStatHelp(pid)
			if err != nil {
				log.Fatal(err)
			}
			repository.AddOrUpdateStudentStatus(uid, pid, "", "Asked for help", "")
		}
	} else {
		msg = "Invalid filename"
	}
	if ok && prob.Active {
		models.HelpSubSem.Lock()
		defer models.HelpSubSem.Unlock()
		sub := &models.HelpSubmission{
			Sid:        int(sid),
			Uid:        uid,
			Pid:        pid,
			Content:    need_help_with,
			Filename:   filename,
			At:         time.Now(),
			SnapshotID: snapshotID,
			Snapshot:   content,
		}
		models.WorkingHelpSubs = append(models.WorkingHelpSubs, sub)
		models.HelpSubmissions[int(sid)] = sub
	}

	fmt.Fprintf(w, msg)

}

//-----------------------------------------------------------------------------------
