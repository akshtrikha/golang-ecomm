package utils

import (
	// "bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Validate is a singleton of the validator struct
// Singleton is helpful as this package uses caching
var Validate = validator.New()

// ParseJSON function to parse the json body received in request
func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		log.Fatal("Request Body is empty")
	}

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(r.Body)

	log.Println("Parsing the JSON Body")
	// log.Println(buf.String())
	return json.NewDecoder(r.Body).Decode(payload)
}

// WriteJSON function
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

// WriteError function
func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}
