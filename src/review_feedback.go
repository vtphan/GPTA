package main

import (
	"log"
)

func countVotes(feedbackID int, voteType string) int {
	var count int64
	err := Database.Model(&SnapshotBackFeedback{}).
		Where("snapshot_feedback_id = ? AND is_helpful = ?", feedbackID, voteType).
		Count(&count).Error
	if err != nil {
		log.Println("Error counting votes:", err)
		return 0
	}
	return int(count)
}

func getFeedbackData(uid int, role string) []*FeedbackData {
	var feedbacks []*FeedbackData

	// Fetch all feedbacks
	var snapshotFeedbacks []SnapshotBackFeedback
	err := Database.Find(&snapshotFeedbacks).Error
	if err != nil {
		log.Fatal("Error fetching snapshot feedbacks:", err)
	}

	for _, feedback := range snapshotFeedbacks {
		// Fetch upvote and downvote counts
		upvote := countVotes(feedback.ID, "yes")
		downvote := countVotes(feedback.ID, "no")

		// Fetch current user's vote
		var userVote SnapshotBackFeedback
		err = Database.Where("snapshot_feedback_id = ? AND author_id = ? AND author_role = ?", feedback.ID, uid, role).
			First(&userVote).Error
		currentUserVote := ""
		if err == nil {
			currentUserVote = userVote.IsHelpful
		}

		// Fetch author name based on role
		authorName := ""
		if feedback.AuthorRole == "teacher" {
			var teacher Teacher
			err = Database.First(&teacher, feedback.AuthorID).Error
			if err == nil {
				authorName = teacher.Name
			}
		} else {
			var student Student
			err = Database.First(&student, feedback.AuthorID).Error
			if err == nil {
				authorName = student.Name
			}
		}

		// Append the feedback to the result list
		feedbacks = append(feedbacks, &FeedbackData{
			FeedbackID:      feedback.ID,
			Feedback:        feedback.IsHelpful, // Or any other relevant field
			FeedbackTime:    feedback.GivenAt,
			Upvote:          upvote,
			Downvote:        downvote,
			CurrentUserVote: currentUserVote,
			GivenBy:         authorName,
			Code:            "", // Adjust as necessary
		})
	}

	return feedbacks
}

//func studentViewsFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
//	spid := r.FormValue("pid")
//	pid := -1
//	nextPid := -1
//	filename := ""
//	if spid == "" {
//		rows, err := Database.Query("select id, filename from problems order by id desc limit 2")
//		defer rows.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//		if rows.Next() {
//			rows.Scan(&pid, &filename)
//		}
//		tmp := ""
//		if rows.Next() {
//			rows.Scan(&nextPid, &tmp)
//		}
//		rows.Close()
//	} else {
//		pid, _ = strconv.Atoi(spid)
//		rows, err := Database.Query("select id from problems where id<? order by id desc limit 1", pid)
//		defer rows.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//		if rows.Next() {
//			rows.Scan(&nextPid)
//		}
//		rows.Close()
//
//		rows2, err2 := Database.Query("select filename from problems where id = ?", pid)
//		defer rows2.Close()
//		if err2 != nil {
//			log.Fatal(err2)
//		}
//		if rows2.Next() {
//			rows2.Scan(&filename)
//		}
//		rows2.Close()
//	}
//	viewType := r.FormValue("viewtype")
//	role := r.FormValue("role")
//	var rows *sql.Rows
//	var err error
//	if viewType == "forme" {
//		rows, err = Database.Query("select F.id, F.feedback, F.author_id, F.author_role, F.given_at, C.code from code_snapshots C, snapshot_feedbacks F where C.id=F.snapshot_id and C.student_id=? and C.problem_id = ? order by F.given_at desc", uid, pid)
//		defer rows.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//	} else if viewType == "all" {
//		active := false
//		for _, prob := range ActiveProblems {
//			if prob.Info.Pid == pid {
//				if prob.Active {
//					active = true
//				}
//				break
//			}
//		}
//		if active == true {
//			if _, ok := HelpEligibleStudents[pid][uid]; ok {
//				rows, err = Database.Query("select F.id, F.feedback, F.author_id, F.author_role, F.given_at, C.code from code_snapshots C, snapshot_feedbacks F where C.id=F.snapshot_id and C.problem_id = ? order by F.given_at desc", pid)
//			} else {
//				rows, err = Database.Query("select F.id, F.feedback, F.author_id, F.author_role, F.given_at, C.code from code_snapshots C, snapshot_feedbacks F where C.id=F.snapshot_id and C.student_id=? and C.problem_id = ? order by F.given_at desc", uid, pid)
//			}
//		} else {
//			rows, err = Database.Query("select F.id, F.feedback, F.author_id, F.author_role, F.given_at, C.code from code_snapshots C, snapshot_feedbacks F where C.id=F.snapshot_id and C.problem_id = ? order by F.given_at desc", pid)
//		}
//		defer rows.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//	} else {
//		log.Fatal("Invalid parameter!")
//	}
//	feedbacks := getFeedbackData(rows, uid, role)
//	data := struct {
//		Feedbacks  []*FeedbackData
//		ViewType   string
//		UserRole   string
//		UserID     int
//		Password   string
//		Filename   string
//		CurrentPid int
//		NextPid    int
//	}{
//		Feedbacks:  feedbacks,
//		ViewType:   viewType,
//		UserRole:   role,
//		UserID:     uid,
//		Password:   r.FormValue("password"),
//		Filename:   filename,
//		CurrentPid: pid,
//		NextPid:    nextPid,
//	}
//	temp := template.New("")
//	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
//	t, err := temp.Funcs(ownFuncs).Parse(STUDENT_VIEWS_FEEDBACK_TEMPLATE)
//	if err != nil {
//		log.Fatal(err)
//	}
//	w.Header().Set("Content-Type", "text/html")
//	err = t.Execute(w, data)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		log.Fatal(err)
//	}
//}

//func teacherViewsFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
//	spid := r.FormValue("pid")
//	pid := -1
//	nextPid := -1
//	filename := ""
//	if spid == "" {
//		rows, err := Database.Query("select id, filename from problems order by id desc limit 2")
//		defer rows.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//		if rows.Next() {
//			rows.Scan(&pid, &filename)
//		}
//		tmp := ""
//		if rows.Next() {
//			rows.Scan(&nextPid, &tmp)
//		}
//		rows.Close()
//	} else {
//		pid, _ = strconv.Atoi(spid)
//		rows, err := Database.Query("select id from problems where id<? order by id desc limit 1", pid)
//		defer rows.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//		if rows.Next() {
//			rows.Scan(&nextPid)
//		}
//		rows.Close()
//
//		rows2, err2 := Database.Query("select filename from problems where id = ?", pid)
//		defer rows2.Close()
//		if err2 != nil {
//			log.Fatal(err2)
//		}
//		if rows2.Next() {
//			rows2.Scan(&filename)
//		}
//		rows2.Close()
//	}
//	role := r.FormValue("role")
//	rows, err := Database.Query("select F.id, F.feedback, F.author_id, F.author_role, F.given_at, C.code from code_snapshots C, snapshot_feedbacks F where C.id=F.snapshot_id and C.problem_id = ? order by F.given_at desc", pid)
//	defer rows.Close()
//	if err != nil {
//		log.Fatal(err)
//	}
//	feedbacks := getFeedbackData(rows, uid, role)
//	data := struct {
//		Feedbacks  []*FeedbackData
//		UserRole   string
//		UserID     int
//		Password   string
//		Filename   string
//		CurrentPid int
//		NextPid    int
//	}{
//		Feedbacks:  feedbacks,
//		UserRole:   role,
//		UserID:     uid,
//		Password:   r.FormValue("password"),
//		Filename:   filename,
//		CurrentPid: pid,
//		NextPid:    nextPid,
//	}
//	temp := template.New("")
//	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
//	t, err := temp.Funcs(ownFuncs).Parse(TEACHER_VIEWS_FEEDBACK_TEMPLATE)
//	if err != nil {
//		log.Fatal(err)
//	}
//	w.Header().Set("Content-Type", "text/html")
//	err = t.Execute(w, data)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		log.Fatal(err)
//	}
//}
