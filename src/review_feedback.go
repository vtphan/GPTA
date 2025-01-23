package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getFeedbackData(rows *sql.Rows, uid int, role string) []*FeedbackData {
	var feedbacks []*FeedbackData
	feedbackID, feedback, authorID, authorRole, givenAt, code := 0, "", 0, "", time.Now(), ""

	upvote, downvote := int64(0), int64(0)
	currentUserVote := ""

	for rows.Next() {
		rows.Scan(&feedbackID, &feedback, &authorID, &authorRole, &givenAt, &code)

		upvote, _ = GetVoteCount(feedbackID, "yes")
		downvote, _ = GetVoteCount(feedbackID, "no")
		rows2, err := Database.Query("select is_helpful from snapshot_back_feedback where snapshot_feedback_id=? and author_id=? and author_role=?", feedbackID, uid, role)
		defer rows2.Close()
		if err != nil {
			log.Fatal(err)
		}
		for rows2.Next() {
			rows2.Scan(&currentUserVote)
		}
		rows2.Close()
		if authorRole == "teacher" {
			rows2, err = Database.Query("select name from teacher where id=?", authorID)
		} else {
			rows2, err = Database.Query("select name from student where id=?", authorID)
		}
		defer rows2.Close()
		if err != nil {
			log.Fatal(err)
		}
		authorName := ""
		if rows2.Next() {
			rows2.Scan(&authorName)
		}
		rows2.Close()
		feedbacks = append(feedbacks, &FeedbackData{
			FeedbackID:      feedbackID,
			Feedback:        feedback,
			FeedbackTime:    givenAt,
			Upvote:          upvote,
			Downvote:        downvote,
			CurrentUserVote: currentUserVote,
			GivenBy:         authorName,
			Code:            code,
		})
		currentUserVote = ""

	}
	return feedbacks
}
