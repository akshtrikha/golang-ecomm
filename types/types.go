package types

import "time"

// UserStore interface to hold all the methods required
// for handling User operations with the database(store)
type UserStore interface {
	GetUserByEmail(string) (*User, error)
	GetUserByID(int) (*User, error)
	CreateUser(User) (int, error)
	GetAllUsers() ([]User, error)
}

// User struct to hold the data regarding the user
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

// RegisterUserPayload struct to hold the payload for /register user endpoint
type RegisterUserPayload struct {
	FirstName string `json:"firstName"  validate:"required"`
	LastName  string `json:"lastName"   validate:"required"`
	Email     string `json:"email"      validate:"required,email"`
	Password  string `json:"password"   validate:"required,min=3,max=60"`
}

// RegisterUserResponse holds the response sent for /register endpoint
type RegisterUserResponse struct {
	Message string `json:"message"`
	ID      int    `json:"id"`
	Email   string `json:"email"`
}

// LoginUserPayload struct to hold the payload for /login user endpoint
type LoginUserPayload struct {
	Email    string `json:"email"       validate:"required,email"`
	Password string `json:"password"    validate:"required,min=3,max=60"`
}

// LoginUserResponse struct to hold the response for /login user endpoint
type LoginUserResponse struct {
	Token     string    `json:"jwt"`
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// ListUsersResponse struct to hold the response fro /list-users api endpoint
type ListUsersResponse struct {
	Users []struct {
		ID        int    `json:"id"`
		FirstName string `json:"lastName"`
		LastName  string `json:"firstName"`
		Email     string `json:"email"`
	}
}

// ProductStore interface to hold all the methods required
// for handling Product operations with the database(store)
type ProductStore interface {
	AddProduct(AddProductPayload) (int, error)
	GetProducts() ([]Product, error)
}

// Product struct is used to hold the info regarding the product
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
}

// AddProductPayload Payload for add-product api endpoint
type AddProductPayload struct {
	Name        string  `json:"name"        validate:"required"`
	Description string  `json:"description" validate:"required"`
	Image       string  `json:"image"       validate:"required"`
	Price       float64 `json:"price"       validate:"required"`
	Quantity    int     `json:"quantity"    validate:"required"`
}