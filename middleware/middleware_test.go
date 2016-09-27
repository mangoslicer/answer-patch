package middleware

import (
	"bytes"
	"encoding/json"
	//	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	auth "github.com/patelndipen/AP1/services"
	"github.com/patelndipen/AP1/settings"
)

type MockTokenStore struct {
	IsStored bool
}

type MockRepStore struct {
	Rep int
}

type MockModel struct {
	Field string
}

func (store *MockTokenStore) StoreToken(key, val string, exp int) error {
	return nil
}

func (store *MockRepStore) FindRep(category, userId string) (int, error) {
	return store.Rep, nil
}

func (store *MockRepStore) UpdateRep(category, userID string, rep int) error {
	return nil
}

func (store *MockTokenStore) IsTokenStored(key string) (bool, error) {
	return store.IsStored, nil
}

func (model *MockModel) GetMissingFields() string {
	if model.Field == "" {
		return "Field"
	}
	return ""
}

func init() {
	settings.SetPreproductionEnv()
}

func TestParseRequestBody(t *testing.T) {

	model := &MockModel{Field: "value"}

	body, err := json.Marshal(model)
	if err != nil {
		t.Error(err)
	}

	r, err := http.NewRequest("", "", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	ac := &auth.AuthContext{UserID: "0", Exp: time.Now()}
	context := &Context{ac, nil, nil}

	ParseRequestBody(new(MockModel), func(c *Context, w http.ResponseWriter, r *http.Request) {

		parsedModel, ok := c.ParsedModel.(*MockModel)
		if !ok {
			http.Error(w, "context.ParsedModel is not of type*MockModel", http.StatusInternalServerError)
		}

		w.Write([]byte(parsedModel.Field))
	})(context, w, r)

	if parsedField := w.Body.String(); parsedField != model.Field {
		t.Errorf("Expected parsedModel.Field to equal %s, but instead %s was retrieved by parsing the request body ", model.Field, parsedField)
	}
}

func TestParseRequestBodyOnEmptyBody(t *testing.T) {

	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	ac := &auth.AuthContext{UserID: "0", Exp: time.Now()}
	context := &Context{ac, nil, nil}

	ParseRequestBody(new(MockModel), func(c *Context, w http.ResponseWriter, r *http.Request) {})(context, w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected a status code of 400 due to the absence of a request body, recieved a status code of %d", w.Code)
	} else if errMessage := w.Body.String(); errMessage != "No data recieved through the request\n" {
		t.Errorf("Expected \"No data recieved through the request\" to be written to the responsewriter body, but the responsewriter body contains: %s", errMessage)
	}
}

func TestCheckRepWithInsufficientRep(t *testing.T) {

	r, err := http.NewRequest("POST", "api/question/testing", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	c := &Context{&auth.AuthContext{UserID: ""}, &MockRepStore{Rep: 1}, nil}

	CheckRep(func(c *Context, w http.ResponseWriter, r *http.Request) {})(c, w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected a http status code of 403 Forbidden, because the rep requirement was not met, but recieved a status code of %d", w.Code)
	} else if w.Body.String() != "Not enough reputation in order to complete the request\n" {
		t.Errorf("Expected the response writer body to contain \"Not enough reputation in order to complete the request\", but instead the response writer body contains %s", w.Body.String())
	}
}

func TestRefreshToken(t *testing.T) {

	w := httptest.NewRecorder()

	RefreshExpiringToken(func(c *Context, w http.ResponseWriter, r *http.Request) {})(NewContext(), w, nil)

	if w.Body == nil {
		t.Errorf("Expected RefreshToken to print a token to the responsewriter body")
	}

}

func TestAuthenticateTokenWithInvalidToken(t *testing.T) {

	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Error(err)
	}

	r.Header.Set("Authorization", "BEARER:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ")

	w := httptest.NewRecorder()

	ac := auth.NewAuthContext(&MockTokenStore{IsStored: false})
	AuthenticateToken(&Context{ac, nil, nil}, func(c *Context, w http.ResponseWriter, r *http.Request) {})(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected the status code to be 401, because of the request contained an invalid JWT, but instead recieved a status code of %d", w.Code)
	} else if w.Body.String() != "Unrecognized signing method: HS256\n" {
		t.Errorf("Expected the responsewriter body to be set to \"Unrecognized signing method: HS256\", but instead the responsewriter body is set to \"%s\"", w.Body.String())
	}
}

func TestAuthenticateTokenWithNoToken(t *testing.T) {

	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	ac := auth.NewAuthContext(&MockTokenStore{IsStored: false})
	AuthenticateToken(&Context{ac, nil, nil}, func(c *Context, w http.ResponseWriter, r *http.Request) {
		if (c.Exp == time.Time{}) && (c.UserID == "") {
			w.Write([]byte("Context has a nil value for both the UserID and Exp fields"))
		}
	})(w, r)

	if w.Body.String() != "Context has a nil value for both the UserID and Exp fields" {
		t.Errorf("Expected the responsewriter's body to contain a message of \"Context has a nil value for both the UserID and Exp fields\" because there was no token in the request")
	}
}

func TestAuthenticateTokenParseOperations(t *testing.T) {

	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Error(err)
	}
	ac := &auth.AuthContext{UserID: "0", TokenStore: &MockTokenStore{IsStored: false}}
	c := &Context{ac, nil, nil}

	//JWT token with a "sub" claim set to "0"

	refreshedToken, err := ac.RefreshToken()
	if err != nil {
		t.Error(err)
	}

	r.Header.Set("Authorization", "BEARER:"+refreshedToken.SignedToken)

	w := httptest.NewRecorder()

	AuthenticateToken(c, func(c *Context, w http.ResponseWriter, r *http.Request) { w.Write([]byte(c.UserID)) })(w, r)

	if w.Body.String() != "0" {
		t.Errorf("Expected the UserID that AuthenticateToken is supposed to determine by parsing the JWT to be \"0\", but the UserID retrieved from the context struct was %s", w.Body.String())
	}
}

func TestAuthenticateTokenWithExpiredToken(t *testing.T) {

	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Error(err)
	}

	c := &Context{auth.NewAuthContext(&MockTokenStore{IsStored: true}), nil, nil}

	refreshedToken, err := c.RefreshToken()
	if err != nil {
		t.Error(err)
	}

	r.Header.Set("Authorization", "BEARER:"+refreshedToken.SignedToken)
	w := httptest.NewRecorder()

	AuthenticateToken(c, func(c *Context, w http.ResponseWriter, r *http.Request) {})(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected the status code to be a 401, because AuthenticateToken recognized that the token is stored in Redis due to the mock IsTokenStored method always returning true, but recieved a status code of %d", w.Code)
	} else if w.Body.String() != "Token is no longer valid\n" {
		t.Errorf("Expected the responsewriter body to contain \"Token is no longer valid\", because AuthenticateToken recognized that the token is stored in Redis due to the mock IsTokenStored method always returning true, but the responsewriter contained %s", w.Body.String())
	}
}
