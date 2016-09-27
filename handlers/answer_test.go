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

type MockAnswerStore struct {
	AnswerSlotAvailable bool
}

func (store *MockAnswerStore) IsAnswerSlotAvailable(questionID string) (bool, error) {
	return store.AnswerSlotAvailable, nil
}

func (store *MockAnswerStore) StoreAnswer(questionID, userID, content string, reqUpvotes int) (error, int) {
	return nil, 0
}

func (store *MockAnswerStore) CastVote(answerID string, vote int) (string, error, int) {
	return "", errors.New("No answer exists with the provided answer id"), http.StatusBadRequest
}

func (store *MockAnswerStore) AssessAnswers(questionID string) (error, int) {
	return nil, 0
}

func TestServeSubmitAnswerWithNoAvailableAnswerSlots(t *testing.T) {

	mockStore := &MockAnswerStore{AnswerSlotAvailable: false}

	r, err := http.NewRequest("PUT", "api/Test/0ab2a26f-c383-45d6-a14f-448eae016641/answer", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	ServeSubmitAnswer(mockStore)(nil, w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected a status code of 401 due to the fact that there were no answer slots available, but recieved a status code of %d", w.Code)
	} else if w.Body.String() != "Maximum capacity for answers has been reached\n" {
		t.Errorf("Expected the responsewriter body to contain \"Maximum capacity for answers has been reached\", but the responsewriter body contains \"%s\"", w.Body.String())
	}
}

func TestServeSubmitAnswer(t *testing.T) {

	mockStore := &MockAnswerStore{AnswerSlotAvailable: true}

	r, err := http.NewRequest("PUT", "api/Test/0ab2a26f-c383-45d6-a14f-448eae016641/answer", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	ServeSubmitAnswer(mockStore)(&m.Context{&auth.AuthContext{UserID: ""}, &MockRepStore{}, &models.Answer{Content: ""}}, w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected a status code of 200, but recieved a status code of %d", w.Code)
	}
}

func TestServeCastAnswerVote(t *testing.T) {

	r, err := http.NewRequest("PUT", "api/Test/0ab2a26f-c383-45d6-a14f-448eae016641/answer", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	mockStore := &MockAnswerStore{}
	ServeCastAnswerVote(mockStore)(nil, w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected a status code of 400, but recieved a status code of %d", w.Code)
	} else if w.Body.String() != "No answer exists with the provided answer id\n" {
		t.Errorf("Expected the responsewriter body to contain \"No answer exists with the provided answer id\", but the responsewriter body contains \"%s\"", w.Body.String())
	}
}
