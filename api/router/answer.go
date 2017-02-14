package router

import (
	"github.com/gorilla/mux"
)

const (
	CreatePendingAnswer = "put:pending_answer"
	UpdateAnswerVote    = "put:answer_vote"
)

func InitAnswerRoutes(r *mux.Router) *mux.Router {

	//PUT
	r.Path("/{category:[a-z]+}/{questionID:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/answer}").Methods("PUT").Name(CreatePendingAnswer)

	r.Path("/{category:[a-z]+}/answer/{answerID:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/vote/{vote:-1|1}}").Methods("PUT").Name(UpdateAnswerVote)
	return r
}
