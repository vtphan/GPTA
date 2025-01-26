// Author: Vinhthuy Phan, 2018
package restHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"net/http"
)

// -----------------------------------------------------------------------------------
func StudentGetsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	var js []byte
	var err error

	models.BoardsSem.Lock()
	defer models.BoardsSem.Unlock()

	if _, ok := models.Students[uid]; ok {
		js, err = json.Marshal(models.Students[uid].Boards)
		for _, b := range models.Students[uid].Boards {
			if b.Pid != 0 {
				repository.AddOrUpdateStudentStatus(uid, b.Pid, "Working", "", "", "")
				_ = repository.IncrementProblemStatActive(b.Pid)
			}
		}
		models.Students[uid].Boards = []*models.Board{}
		if err == nil {
			// fmt.Println(string(js))
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	}
	fmt.Println(err.Error())
	js, err = json.Marshal([]*models.Board{})
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//-----------------------------------------------------------------------------------
