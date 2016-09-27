package handlers

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/patelndipen/AP1/datastores"
	m "github.com/patelndipen/AP1/middleware"
	"github.com/patelndipen/AP1/models"
	"github.com/patelndipen/AP1/router"
)

func AssignHandlersToRoutes(c *m.Context, db *sql.DB) *mux.Router {

	r := router.InitRouter()
	r = AssignHandlersToQuestionRoutes(r, c, db)
	r = AssignHandlersToUserRoutes(r, c, db)

	return r
}

func AssignHandlersToQuestionRoutes(r *mux.Router, c *m.Context, db *sql.DB) *mux.Router {

	questionStore := &datastores.QuestionStore{db}
	//	answerStore := &datastores.AnswerStore{db}

	r.Get(router.ReadPost).Handler(m.AuthenticateToken(c, m.RefreshExpiringToken(ServePostByID(questionStore))))

	r.Get(router.ReadQuestionsByFilter).Handler(m.AuthenticateToken(c, m.RefreshExpiringToken(ServeQuestionsByFilter(questionStore))))

	r.Get(router.ReadSortedQuestions).Handler(m.AuthenticateToken(c, m.RefreshExpiringToken(ServeSortedQuestions(questionStore))))

	r.Get(router.CreateQuestion).Handler(m.AuthenticateToken(c, m.RefreshExpiringToken(ServeSubmitQuestion(questionStore))))

	return r
}

func AssignHandlersToUserRoutes(r *mux.Router, c *m.Context, db *sql.DB) *mux.Router {

	userStore := &datastores.UserStore{db}

	r.Get(router.ReadUser).Handler(m.AuthenticateToken(c, m.RefreshExpiringToken(ServeFindUser(userStore))))

	r.Get(router.CreateUser).Handler(m.ServeHTTP(m.ParseRequestBody(new(models.UnauthUser), ServeRegisterUser(userStore))))

	r.Get(router.Login).Handler(m.ServeHTTP(m.ParseRequestBody(new(models.UnauthUser), ServeLogin(userStore))))

	r.Get(router.Logout).Handler(m.AuthenticateToken(c, ServeLogout()))

	return r
}
