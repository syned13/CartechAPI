package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	us "github.com/CartechAPI/user"
	"github.com/CartechAPI/utils"
)

type apiResponse struct {
	Message string `json:"message"`
}

var (
	ErrMissingFields  = errors.New("missing_fields")
	ErrNotUniqueEmail = errors.New("email must be unique")
)

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
		return
	}
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := us.User{}

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid_request_body")
			return
		}

		if user.Email == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "missing email")
			return
		}

		if user.Password == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "missing password")
		}

	}
}

func validateSignUpFields(user us.User) error {
	if user.Email == "" || user.LastName == "" || user.Name == "" || user.Password == "" || user.PhoneNumber == "" {
		return ErrMissingFields
	}

	return nil
}

func validateUniqueCredentials(db *sql.DB, user us.User) error {
	retrievedUser, err := us.GetUserByEmail(db, user.Email)
	if err != nil {
		return err
	}

	if retrievedUser != nil {
		return ErrNotUniqueEmail
	}

	return nil
}

func SignUp(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &us.User{}

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid_request_body")
			return
		}

		err = validateSignUpFields(*user)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadGateway, err.Error())
			return
		}

		err = validateUniqueCredentials(db, *user)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadGateway, err.Error())
			return
		}

		user, err = CreateUser(db, user)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "unexpected error")
			return
		}

		user.Password = ""
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}
