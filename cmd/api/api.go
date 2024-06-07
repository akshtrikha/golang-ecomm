package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/akshtrikha/golang-ecomm/services/product"
	"github.com/akshtrikha/golang-ecomm/services/user"
	"github.com/gorilla/mux"
)

// APIServer struct
type APIServer struct {
	addr string
	db   *sql.DB
}

// NewAPIServer constructor to create and return a new APIServer
func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

// Run func
func (s *APIServer) Run() error {
	// Create a mux router
	router := mux.NewRouter()

	// create a subrouter out of the router
	// the main router routes all the /api/v1 apis.
	// this helps use to do versioning of the api endpoints.
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// create a dao for user
	userStore := user.NewStore(s.db)
	productStore := product.NewStore(s.db)

	// this is used to create a handler of the user service.
	// the user handler will help us handle routes related to the user.
	// here we are injecting the userStore dependency to the handler.
	// this will allow the handler to do everything with the user.
	// from routing to handing user data
	userHandler := user.NewHandler(userStore)
	productHandler := product.NewHandler(productStore)

	// pass the subrouter to this function
	// to delegeate the route management
	// for user service
	userHandler.RegisterRoutes(subrouter)
	productHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	// start the http server on s.addr
	// pass the router to handle routing of this server
	return http.ListenAndServe(s.addr, router)
}
