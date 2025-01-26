// Author: Vinhthuy Phan, 2018
package restHandlers

import (
	"fmt"
	"github.com/GPTA/src/frontEnd"
	"github.com/GPTA/src/models"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	// "strconv"
)

type AnswersBoardMessage struct {
	Counts  map[string]int
	Content string
	Total   int
}

// -----------------------------------------------------------------------------------
func ViewAnswersHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("filename")
	passcode := r.FormValue("pc")
	if prob, ok := models.ActiveProblems[filename]; ok && passcode == models.Passcode {
		t, err := template.New("").Parse(frontEnd.VIEW_ANSWERS_TEMPLATE)
		if err == nil {
			answers := prob.Answers
			counts := make(map[string]int)
			total := 0
			for i := 0; i < len(answers); i++ {
				counts[answers[i]]++
				total++
			}
			content := prob.Info.Description
			w.Header().Set("Content-Type", "text/html")
			data := &AnswersBoardMessage{Counts: counts, Content: content, Total: total}
			err = t.Execute(w, data)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}
}

//-----------------------------------------------------------------------------------
