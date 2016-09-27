package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/patelndipen/AP1/datastores"
	m "github.com/patelndipen/AP1/middleware"
	"github.com/patelndipen/AP1/models"
	"github.com/patelndipen/AP1/services"
)

func ServeSubmitAnswer(store datastores.AnswerStoreServices) m.HandlerFunc {
	return func(c *m.Context, w http.ResponseWriter, r *http.Request) {

		questionID := mux.Vars(r)["questionID"]
		isSlotAvailable, err := store.IsAnswerSlotAvailable(questionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if !isSlotAvailable {
			http.Error(w, "Maximum capacity for answers has been reached", http.StatusForbidden)
			return
		}

		newAnswer := c.ParsedModel.(*models.Answer)
		requiredRep, err := c.RepStore.FindRep(mux.Vars(r)["category"], c.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err, statusCode := store.StoreAnswer(questionID, c.UserID, newAnswer.Content, services.CalculateCurrentAnswerEligibilityRep(requiredRep))
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func ServeCastAnswerVote(store datastores.AnswerStoreServices) m.HandlerFunc {
	return func(c *m.Context, w http.ResponseWriter, r *http.Request) {

		routeVars := mux.Vars(r)

		/*
			if !store.IsCategoryRegistered(routeVars["category"]) {
				http.Error(w, "The provided category does not exist", http.StatusBadRequest)
				return
			}
		*/
		vote := 1

		if routeVars["vote"] == "downvote" {
			vote = -1
		}

		voteRecipient, err, statusCode := store.CastVote(routeVars["answerID"], vote)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		rep, err := c.RepStore.FindRep(routeVars["category"], voteRecipient)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rep <= MAX_REP {
			err = c.RepStore.UpdateRep(routeVars["cateogry"], c.UserID, vote)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		err, statusCode = store.AssessAnswers(routeVars["answerID"])
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}
	}
}
