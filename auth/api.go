package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	us "github.com/CartechAPI/user"
	"github.com/CartechAPI/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type apiResponse struct {
	Message string `json:"message"`
}

var tokenSigningKey string

var (
	ErrMissingFields  = errors.New("missing_fields")
	ErrNotUniqueEmail = errors.New("email must be unique")
)

func init() {
	tokenSigningKey = os.Getenv("SECRET")
}

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
		return
	}
}

func Login(db *sql.DB) http.HandlerFunc {
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
			return
		}

		userRetrieved, err := us.GetUserByEmail(db, user.Email)
		if err != nil && err != sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusInternalServerError, "unexpected error")
			return
		}

		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusBadRequest, "incorrect email or password")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(userRetrieved.Password), []byte(user.Password))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			utils.RespondWithError(w, http.StatusBadRequest, "incorrect email or password")
			return
		}

		user = *userRetrieved
		now := time.Now()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.UserID,
			"iat":     now.String(),
		})

		signedToken, err := token.SignedString([]byte(os.Getenv("SECRET")))
		if err != nil {
			fmt.Println("error_signing_token: " + err.Error())
			utils.RespondWithError(w, http.StatusInternalServerError, "unexpected error")
			return
		}

		user.Password = ""
		marshalledUser, _ := json.Marshal(user)

		var responseMap map[string]interface{}

		_ = json.Unmarshal(marshalledUser, &responseMap)
		responseMap["token"] = signedToken

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseMap)
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
