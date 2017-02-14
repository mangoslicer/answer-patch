package router

import (
	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r = InitQuestionRoutes(r)
	r = InitAnswerRoutes(r)
	r = InitUserRoutes(r)

	return r
}
