package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/patelndipen/AP1/datastores"
	m "github.com/patelndipen/AP1/middleware"
	"github.com/patelndipen/AP1/models"
	"github.com/patelndipen/AP1/services"
)

const (
	MAX_REP             = 25
	QUESTION_ASKING_FEE = -2
)

func ServePostByID(store datastores.QuestionStoreServices) m.HandlerFunc {
	return func(c *m.Context, w http.ResponseWriter, r *http.Request) {
		var post []models.ModelServices
		question, answer, err, statusCode := store.FindPostByID(mux.Vars(r)["questionId"])
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		post = append(post, question)
		if answer != nil {
			post = append(post, answer)
		}
		services.PrintJSON(w, post)
	}
}

func ServeQuestionsByFilter(store datastores.QuestionStoreServices) m.HandlerFunc {
	return func(c *m.Context, w http.ResponseWriter, r *http.Request) {

		questions, err, statusCode := store.FindQuestionsByFilter(mux.Vars(r)["filter"], mux.Vars(r)["val"])
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}
		services.PrintJSON(w, questions)
	}

}

func ServeSortedQuestions(store datastores.QuestionStoreServices) m.HandlerFunc {
	return func(c *m.Context, w http.ResponseWriter, r *http.Request) {
		routeVars := mux.Vars(r)
		questions, err, statusCode := store.SortQuestions(routeVars["postComponent"], routeVars["sortedBy"], routeVars["order"], routeVars["offset"])
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}
		services.PrintJSON(w, questions)
	}
}

func ServeSubmitQuestion(store datastores.QuestionStoreServices) m.HandlerFunc {
	return func(c *m.Context, w http.ResponseWriter, r *http.Request) {

		newQuestion := c.ParsedModel.(*models.Question)
		category := mux.Vars(r)["category"]

		err := c.RepStore.UpdateRep(category, c.UserID, QUESTION_ASKING_FEE)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		err, statusCode := store.StoreQuestion(newQuestion.UserID, category, newQuestion.Title, newQuestion.Content)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		w.WriteHeader(http.StatusCreated)

	}
}

func ServeCastQuestionVote(store datastores.AnswerStoreServices) m.HandlerFunc {
	return func(c *m.Context, w http.ResponseWriter, r *http.Request) {
		vote := 1
		urlParams := mux.Vars(r)

		if urlParams["vote"] == "downvote" {
			vote = -1
		}

		voteRecipient, err, statusCode := store.CastVote(urlParams["questionID"], vote)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
		}

		rep, err := c.RepStore.FindRep(urlParams["category"], voteRecipient)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if rep <= MAX_REP {
			c.RepStore.UpdateRep(urlParams["category"], voteRecipient, 1)
		}

		err, statusCode = store.AssessAnswers(urlParams["questionID"])
		if err != nil {
			http.Error(w, err.Error(), statusCode)
		}
	}
}

/*
		newCurrentAnsRecipient := store.AssessAnswers(urlParams["questionID"])

		if voteRecipient != newCurrentRecipient && c.RepStore.FindRep(urlParams["category"], voteRecipient) <= maxRep {
			c.RepStore.UpdateRep(urlParams["category"], voteRecipient, 1)
		} else if newCurrentAnsRecipient != "" {
			c.RepStore.UpdateRep(urlParams["category"], newCurrentAnsRecipient, )
		}
	}
*/
