package router

import (
	"github.com/gorilla/mux"
)

const (
	ReadPost              = "get:post"
	ReadQuestionsByFilter = "get:questions_by_filter"
	ReadSortedQuestions   = "get:sorted_questions"
	CreateQuestion        = "post:question"
)

func InitQuestionRoutes(r *mux.Router) *mux.Router {

	//GET
	r.Path("/post/{questionId:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}").Methods("GET").Name(ReadPost)
	r.Path("/questions/{filter:posted-by|answered-by|category}/{val:[A-Za-z0-9]+}").Methods("GET").Name(ReadQuestionsByFilter)
	r.Path("/{postComponent:questions|answers}/{sortedBy:upvotes|edits|date}/{order:desc|asc}/{offset:[0-9]+}").Methods("GET").Name(ReadSortedQuestions)

	//POST
	r.Path("/question/{category:[a-z]+}").Methods("POST").Name(CreateQuestion)

	return r

}
