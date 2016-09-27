package router

import "github.com/gorilla/mux"

const (
	ReadUser   = "get:user"
	CreateUser = "post:user"
	Login      = "post:login"
	Logout     = "post:logout"
)

func InitUserRoutes(r *mux.Router) *mux.Router {

	//GET
	r.Path("/{filter:id|username}/{searchVal:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}|[a-z0-9]}").Methods("GET").Name(ReadUser)

	//POST
	r.Path("/register").Methods("POST").Name(CreateUser)
	r.Path("/login").Methods("POST").Name(Login)
	r.Path("/logout").Methods("POST").Name(Logout)

	return r

}
