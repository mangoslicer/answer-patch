package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/patelndipen/AP1/datastores"
	"github.com/patelndipen/AP1/handlers"
	m "github.com/patelndipen/AP1/middleware"
	auth "github.com/patelndipen/AP1/services"
	"github.com/patelndipen/AP1/settings"
)

type Server struct {
	r *mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding")
	}

	if r.Method == "OPTIONS" {
		return
	}

	s.r.ServeHTTP(w, r)
}

func main() {

	settings.SetPreproductionEnv() // Set GO_ENV to "preproduction"

	db := datastores.ConnectToPostgres()

	ac := auth.NewAuthContext(&datastores.JWTStore{datastores.ConnectToRedis()})
	c := &m.Context{ac, &datastores.RepStore{datastores.ConnectToMongoCol()}, nil}

	r := handlers.AssignHandlersToRoutes(c, db)
	http.Handle("/", &Server{r})

	fmt.Println("Listening on port 3030")
	http.ListenAndServe(":3030", nil)
}
