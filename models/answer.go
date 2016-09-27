package models

import (
	"time"
)

type Answer struct {
	ID              string    `json:"answerID"`
	QuestionID      string    `json:"answerID"`
	UserID          string    `json:"answerUserID"`
	Username        string    `json:"answerUsername"`
	IsCurrentAnswer bool      `json:"answerCurrent"`
	Content         string    `json:"answerContent"`
	Upvotes         int       `json:"answerUpvotes"`
	ReqUpvotes      int       `json:"answerRequiredUpvotes"`
	LastEditedAt    time.Time `json:"answerLastEditedAt"`
}

func NewAnswer() *Answer {
	return &Answer{}
}

func (answer *Answer) GetMissingFields() string {

	if answer.Content == "" {
		return "Content\n"
	}

	return ""
}
