package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akshtrikha/golang-ecomm/types"
	"github.com/gorilla/mux"
)

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("User not found")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(u types.User) (int, error) {
	return 0, nil
}

func (m *mockUserStore) GetAllUsers() ([]types.User, error)

// TestUserServiceHandlers functinoon to implement testing
func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("Should fail if the user payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "test",
			LastName:  "test",
			Email:     "invalid-email",
			Password:  "test-password",
		}

		marshalled, _ := json.Marshal(payload)
		log.Println(string(marshalled))

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("Should correctly register the user", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "test",
			LastName:  "test",
			Email:     "valid@email.com",
			Password:  "test-password",
		}

		marshalled, _ := json.Marshal(payload)
		log.Println(string(marshalled))

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}
