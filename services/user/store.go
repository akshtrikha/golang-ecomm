package user

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/akshtrikha/golang-ecomm/types"
)

// Store struct to hold the database object
// This will be used to handle the database queries
type Store struct {
	db *sql.DB
}

// NewStore function to return a reference to the Store struct
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetUserByEmail function to run a SQL query and find the user by email
func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	// log.Printf("Inside store.GetUserByEmail with email: %s", email)
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)

		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user with email: %s not found", email)
	}

	log.Println("Returning the user")
	return u, nil
}

// GetUserByID function to run a SQL query and find the user by id
func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)

		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user with id: %v not found", id)
	}

	return u, nil
}

// CreateUser function to run a SQL query and create the user
func (s *Store) CreateUser(u types.User) (int, error) {
	result, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", u.FirstName, u.LastName, u.Email, u.Password)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetAllUsers to get all the users from the database
func (s *Store) GetAllUsers() ([]types.User, error) {
	result, err := s.db.Query("SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	var users []types.User

	for result.Next() {
		user, err := scanRowIntoUser(result)
		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return users, nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
