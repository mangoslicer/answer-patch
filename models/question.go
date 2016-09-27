package models

import (
	"time"
)

type Question struct {
	ID           string    `json:"questionID"`
	UserID       string    `json:"questionUserID"`
	Username     string    `json:"questionUsername"`
	Category     string    `json:"questionCategory"`
	Title        string    `json:"questionTitle"`
	Content      string    `json:"questionContent"`
	Upvotes      int       `json:"questionUpvotes"`
	EditCount    int       `json:"answerEditCount"`
	PendingCount int       `json:"pendingAnswerCount"`
	SubmittedAt  time.Time `json:"questionSubmittedAt"`
}

func (question *Question) GetMissingFields() string {

	var missing string

	switch {
	case question.UserID == "":
		missing = "Author's user ID\n"
	case question.Username == "":
		missing += "Author's username\n"
	case question.Title == "":
		missing += "Title\n"
	case question.Category == "":
		missing += "Category\n"
	}

	return missing
}
