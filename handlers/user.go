package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/patelndipen/AP1/datastores"
	m "github.com/patelndipen/AP1/middleware"
	"github.com/patelndipen/AP1/models"
	"github.com/patelndipen/AP1/services"
)

func ServeFindUser(store datastores.UserStoreServices) m.HandlerFunc {
	return func(c *m.Context, w http.ResponseWriter, r *http.Request) {

		user, err, statusCode := store.FindUser(mux.Vars(r)["filter"], mux.Vars(r)["searchVal"])
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}
		services.PrintJSON(w, user)
	}
}

func ServeRegisterUser(store datastores.UserStoreServices) m.HandlerFunc {
	return func(c *m.Context, w http.ResponseWriter, r *http.Request) {

		newUser := c.ParsedModel.(*models.UnauthUser)

		/*
			if store.IsUsernameRegistered(newUser.Username) {
				http.Error(w, "Username already exists", http.StatusBadRequest)
				return
			}
		*/

		err, statusCode := store.StoreUser(newUser.Username, newUser.HashPassword())
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		w.WriteHeader(http.StatusCreated)

	}
}

func ServeLogin(store datastores.UserStoreServices) m.HandlerFunc {
	return func(c *m.Context, w http.ResponseWriter, r *http.Request) {

		credentials := c.ParsedModel.(*models.UnauthUser)

		retrievedUser, err, statusCode := store.FindUser("username", credentials.Username)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		c.UserID = retrievedUser.ID
		token, err := c.Login(credentials.Password, retrievedUser.HashedPassword)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		services.PrintJSON(w, token)
	}
}

func ServeLogout() m.HandlerFunc {
	return func(c *m.Context, w http.ResponseWriter, r *http.Request) {
		c.Logout(r.Header.Get("Authorization")[7:]) //Sends the signed token without the "BEARER:" prefix
	}
}
