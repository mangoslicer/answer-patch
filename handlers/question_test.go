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

type MockQuestionStore struct {
	ExistingID string
}

type MockRepStore struct {
}

func (store *MockQuestionStore) FindPostByID(id string) (*models.Question, *models.Answer, error, int) {
	return nil, nil, errors.New("No question exists with the provided id"), http.StatusBadRequest

}

func (store *MockQuestionStore) FindQuestionsByFilter(filter, val string) ([]*models.Question, error, int) {
	return nil, errors.New("No question(s) found with the provided query"), http.StatusBadRequest

}

func (store *MockQuestionStore) SortQuestions(postComponent, filter, order, offset string) ([]*models.Question, error, int) {
	return nil, errors.New("No questions match the specifications in the url"), http.StatusBadRequest
}

func (store *MockQuestionStore) StoreQuestion(user_id, title, content, category string) (error, int) {
	return errors.New("The provided title is not unique"), http.StatusBadRequest
}

func (store *MockRepStore) FindRep(category, userID string) (int, error) {
	return 0, nil
}

func (store *MockRepStore) UpdateRep(category, userID string, rep int) error {
	return nil
}

func TestServePostByIDWithInvalidID(t *testing.T) {

	//Creates a request with an invalid ID
	r, err := http.NewRequest("GET", "api/posts/5975ea52-2a91-483f-b52f-2b0257886773", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	ServePostByID(new(MockQuestionStore))(m.NewContext(), w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected a status code of 400, but recieved an http status code of %d", w.Code)
	} else if w.Body.String() != "No question exists with the provided id\n" {
		t.Errorf("Expected the content of the responsewriter to be \"No question exists with the provided id\", but instead the responsewriter contains %s", w.Body.String())
	}
}

func TestServeQuestionsByFilterWithInvalidAuthor(t *testing.T) {

	//Creates a request with an invalid Author
	r, err := http.NewRequest("GET", "api/questions/posted-by/NonExistent", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	ServeQuestionsByFilter(new(MockQuestionStore))(m.NewContext(), w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected a status code of 400, because the MockQuestionStore's FindQuestionsByAuthor method always returns nil as a result, but recieved an http status code of %d", w.Code)
	} else if w.Body.String() != "No question(s) found with the provided query\n" {
		t.Errorf("Expected the content of the responsewriter to be \"No question(s) found with the provided query\", but instead the responsewriter contains %s", w.Body.String())
	}
}

func TestServeSortedQuestions(t *testing.T) {

	//Creates a request with filters that no questions satisfy
	r, err := http.NewRequest("GET", "api/questions/upvotes/desc/10", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	ServeSortedQuestions(new(MockQuestionStore))(m.NewContext(), w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected a status code of 400, because the MockQuestionStore's FindQuestionsByFilter method always returns nil as a result, but recieved an http status code of %d", w.Code)
	} else if w.Body.String() != "No questions match the specifications in the url\n" {
		t.Errorf("Expected the content of the responsewriter to be \"No questions match the specifications in the url\", but instead the responsewriter contains %s", w.Body.String())
	}
}

func TestServeSubmitQuestionWithExistingQuestion(t *testing.T) {

	existingQuestion := &models.Question{UserID: "0c1b2b91-9164-4d52-87b0-9c4b444ee62d", Username: "Tester1", Title: "Where is the best sushi place?", Content: "I have cravings"}

	c := &m.Context{auth.NewAuthContext(nil), &MockRepStore{}, existingQuestion}

	r, err := http.NewRequest("POST", "api/question/TestCategory", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	ServeSubmitQuestion(&MockQuestionStore{ExistingID: "526c4576-0e49-4e90-b760-e6976c698574"})(c, w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected a status code of 400 due to the existence of a question with the same title as that of the question recieved in the request body, recieved a status code of %d", w.Code)

	} else if w.Body.String() != "The provided title is not unique\n" {
		t.Errorf("Expected the content of the responsewriter to be \"The provided title is not unique\", but instead the responsewriter contains %s", w.Body.String())
	}
}
