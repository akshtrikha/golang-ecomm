package product

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
)

// API Endpoints to create products and get all the products

// Handler to the product store which will deal
// with the database regarding products
type Handler struct {
	store types.ProductStore
}

// NewHandler constructor
func NewHandler(s *Store) *Handler {
	return &Handler{store: s}
}

// RegisterRoutes func for products
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/get-products", h.handleGetProducts).Methods("GET")
	router.HandleFunc("/add-product", h.handleAddProduct).Methods("POST")
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("handle /get-products hit")
	// get token from the headers
	token := r.Header.Get("Authorization")
	if token == "" {
		log.Println("Token missing")
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Missing token"))
		return
	}

	// verify the token
	claims, err := auth.VerifyTokenAndClaims(token, config.Envs.JWTSecret)
	if err != nil {
		log.Println("Invalid token")
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	userID := claims["userID"]
	log.Printf("Request from the user: %v", userID)

	// get the products from the database
	products, err := h.store.GetProducts()
	if err != nil {
		log.Println("Error fetching the products from the database")
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// return the products
	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handleAddProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("handler /add-products hit")

	// get the token
	token := r.Header.Get("Authorization")
	if token == "" {
		log.Println("Token Missing")
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Missing token"))
		return
	}

	// verify the token
	claims, err := auth.VerifyTokenAndClaims(token, config.Envs.JWTSecret)
	if err != nil {
		log.Println("Invalid token")
		utils.WriteError(w, http.StatusUnauthorized, err)
	}

	userID := claims["userID"]
	log.Printf("Request from the user: %v", userID)

	// get the payload
	var payload types.AddProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		log.Printf("Error parsing the payload\n")
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Parsing the payload completed: %+v", payload)

	// verify the payload
	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid Payload: %v", validationErrors))
		return
	}

	// add the product
	productID, err := h.store.AddProduct(payload)
	log.Printf("Product Added %v", productID)

	// return the product id
	utils.WriteJSON(w, http.StatusCreated, map[string]int{"id": productID})
}
