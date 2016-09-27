package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/patelndipen/AP1/models"
	"github.com/patelndipen/AP1/services"
	"github.com/patelndipen/AP1/settings"
)

const (
	MIN_REP_FOR_ASKING_QUESTION = 10
)

type HandlerFunc func(*Context, http.ResponseWriter, *http.Request)

func ServeHTTP(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := NewContext()
		fn(c, w, r)
	}
}

func ParseRequestBody(model models.ModelServices, fn HandlerFunc) HandlerFunc {
	return func(c *Context, w http.ResponseWriter, r *http.Request) {

		url := r.URL.RequestURI()
		if (url != "/login") && (url != "/register") && (c.UserID == "") && (c.Exp == time.Time{}) {
			http.Error(w, "JWT authentication required in order to complete this request", http.StatusUnauthorized)
			return
		}
		//Checks whether the request body is in JSON format
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			http.Error(w, "This api only accepts JSON payloads. Be sure to specify the \"Content-Type\" of the payload in the request header.", http.StatusBadRequest)
			return
		} else if r.Body == nil {
			http.Error(w, "No data recieved through the request", http.StatusBadRequest)
			return
		}

		//Parses request body

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			log.Fatal(err)
		}

		defer r.Body.Close()

		err = json.Unmarshal(body, model)
		if err != nil {
			// 422 -unprocessable entity
			http.Error(w, err.Error()+"\n", 422)
			return
		}

		missing := model.GetMissingFields()
		if missing != "" {
			http.Error(w, "The following fields were not recieved:\n"+missing, http.StatusBadRequest)
			return
		}

		c.ParsedModel = model

		fn(c, w, r)
	}
}

func CheckRep(fn HandlerFunc) HandlerFunc {
	return func(c *Context, w http.ResponseWriter, r *http.Request) {
		category := mux.Vars(r)["category"]
		rep, err := c.RepStore.FindRep(category, c.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rep < MIN_REP_FOR_ASKING_QUESTION {
			http.Error(w, "Not enough reputation in order to complete the request", http.StatusForbidden)
			return
		}

		fn(c, w, r)
	}
}

func RefreshExpiringToken(fn HandlerFunc) HandlerFunc {
	return func(c *Context, w http.ResponseWriter, r *http.Request) {

		//Refreshes token, if the token expires in less than 24 hours
		if (c.Exp != time.Time{}) && (c.Exp.Sub(time.Now()) < (time.Duration(24) * time.Hour)) {

			refreshedToken, err := c.RefreshToken()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			services.PrintJSON(w, refreshedToken)
		}
		fn(c, w, r)
	}
}

func AuthenticateToken(c *Context, fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token, err := jwt.ParseFromRequest(r, func(parsedToken *jwt.Token) (interface{}, error) {
			if _, ok := parsedToken.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unrecognized signing method: %v", parsedToken.Header["alg"])

			} else {
				return settings.GetPublicKey(), nil
			}
		})

		if err == jwt.ErrNoTokenInRequest {
			fn(c, w, r)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid JWT", http.StatusUnauthorized)
			return
		}

		var ok bool

		c.UserID, ok = token.Claims["sub"].(string)
		if !ok {
			log.Fatal("The underlying type of sub is not string")
		}

		isStored, err := c.TokenStore.IsTokenStored(c.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if isStored {
			http.Error(w, "Token is no longer valid", http.StatusUnauthorized)
			return
		}

		_, ok = token.Claims["exp"].(float64)
		if !ok {
			log.Fatal("The underlying type of exp is not float64")
		}

		fn(c, w, r)
	}
}
