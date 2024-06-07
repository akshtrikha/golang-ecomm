package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akshtrikha/golang-ecomm/config"
	"github.com/akshtrikha/golang-ecomm/services/auth"
	"github.com/akshtrikha/golang-ecomm/types"
	"github.com/akshtrikha/golang-ecomm/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/maolinc/copier"
)

// Handler struct
type Handler struct {
	store types.UserStore
}

// NewHandler constructor takes UserStore as a dependency
// This will allow the Handler to manage user data in the database
func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

// RegisterRoutes func
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/list-users", h.handleListUsers).Methods("GET")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("handle /login endpoint hit")

	// get the json payload
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		log.Println("Payload parsing met with an error")
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	log.Println("Payload parsing completed")

	// validate the json payload
	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload: %v", validationErrors))
		return
	}

	// find the user
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	log.Printf("User found: %+v\n", u)

	// check password
	if ok := auth.ComparePassword(u.Password, payload.Password); !ok {
		log.Println("Invalid Password")
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Invalid Password"))
		return
	}

	// generate jwt
	secret := config.Envs.JWTSecret
	token, err := auth.GenerateJWT(secret, u.ID)
	if err != nil {
		log.Println("Error creating jwt")
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// respond with jwt and user in the response payload
	var response types.LoginUserResponse
	copier.Copy(&response, &u)
	response.Token = token

	log.Printf("Response: %+v\n", response)

	utils.WriteJSON(w, http.StatusFound, response)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("handle /register endpoint hit")

	// get the json payload
	// var payload types.RegisterUserPayload = types.RegisterUserPayload{}
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		log.Println("Payload Parsing met with an error")
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	log.Println("Payload parsing completed")
	// validating the payload
	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload: %v", validationErrors))
		return
	}

	// check if the user exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// create the new user if it doesn't exist
	id, err := h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	var response types.RegisterUserResponse
	response.Message = "User Created Successfully"
	response.Email = payload.Email
	response.ID = id

	log.Printf("Response: %+v", response)

	utils.WriteJSON(w, http.StatusCreated, response)
}

func (h *Handler) handleListUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("handle /list-users endpoint hit")

	// check the token is present or not
	token := r.Header.Get("Authorization")
	if token == "" {
		log.Println("jwt token not found in the headers")
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Please register/login first"))
		return
	}

	// get the token and verify
	secret := config.Envs.JWTSecret
	claims, err := auth.VerifyTokenAndClaims(token, secret)
	if err != nil {
		log.Printf("token verification failed, error: %+v", err)
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Invalid Token"))
		return
	}

	// get all the users from the database
	userID := claims["userID"]
	log.Printf("Request from the user: %v", userID)

	result, err := h.store.GetAllUsers()
	if err != nil {
		log.Println("Error encountered while getting all the users")
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	var users types.ListUsersResponse
	copier.Copy(&users.Users, &result)

	// respond
	utils.WriteJSON(w, http.StatusOK, users.Users)
}
