package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	m "github.com/patelndipen/AP1/middleware"
	"github.com/patelndipen/AP1/models"
	auth "github.com/patelndipen/AP1/services"
)

type MockUserStore struct {
	FindUserErr        error
	FindUserStatusCode int
	//IsRegistered bool
}

func (store *MockUserStore) FindUser(filter, searchVal string) (*models.User, error, int) {
	return nil, store.FindUserErr, store.FindUserStatusCode
}

func (store *MockUserStore) StoreUser(username, hashedpassword string) (error, int) {
	return errors.New("Username already exists"), http.StatusConflict
}

/*
func (store *MockUserStore) IsUsernameRegistered(username string) bool {
	return store.IsRegistered
}
*/

func TestServeFindUserWithInvalidUser(t *testing.T) {

	r, err := http.NewRequest("GET", "api/username/NonExistent", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	ServeFindUser(&MockUserStore{FindUserErr: errors.New("No user exists with the provided information"), FindUserStatusCode: http.StatusBadRequest})(m.NewContext(), w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected a status code of 400, but recieved an http status code of %d", w.Code)
	} else if w.Body.String() != "No user exists with the provided information\n" {
		t.Errorf("Expected the content of the responsewriter to be \"No user exists with the provided information\", but instead the responsewriter contains %s", w.Body.String())
	}
}

func TestServeRegisterUserWithRegisteredUsername(t *testing.T) {

	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	registeredUser := &models.UnauthUser{Username: "RegisteredUsername"}
	c := &m.Context{ParsedModel: registeredUser}
	ServeRegisterUser(&MockUserStore{})(c, w, r)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected a status code of 409 Conflict, but recieved an http status code of %d", w.Code)
	} else if w.Body.String() != "Username already exists\n" {
		t.Errorf("Expected the content of the responsewriter to be \"Username already exists\", but instead the responsewriter contains %s", w.Body.String())
	}
}

func TestServeLoginWithIncorrectCredentials(t *testing.T) {

	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	ac := &auth.AuthContext{UserID: "ID"}
	unauthUser := &models.UnauthUser{Username: "Username", Password: "Wrong Password"}
	c := &m.Context{ac, nil, unauthUser}

	ServeLogin(&MockUserStore{FindUserErr: errors.New("No user exists with the provided credential"), FindUserStatusCode: http.StatusUnauthorized})(c, w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected a status code of 401 Unauthorized, but recieved an http status code of %d", w.Code)
	} else if w.Body.String() != "No user exists with the provided credential\n" {
		t.Errorf("Expected the content of the responsewriter to be \"No user exists with the provided credential\", but instead the responsewriter contains %s", w.Body.String())
	}
}
