package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	mec "github.com/CartechAPI/mechanic"
	"github.com/CartechAPI/shared"
	us "github.com/CartechAPI/user"
	"github.com/CartechAPI/utils"
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

// Index returns handler of GET / endpoint
func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
		return
	}
}

// Login returns handler of POST /login endpoint
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

		if err == sql.ErrNoRows || userRetrieved == nil {
			utils.RespondWithError(w, http.StatusBadRequest, "incorrect email or password")
			return
		}

		if !isPasswordCorrect(user.Password, userRetrieved.Password) {
			utils.RespondWithError(w, http.StatusBadRequest, "incorrect email or password")
			return
		}

		user = *userRetrieved
		token, err := GenerateToken(&user)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		user.Password = ""
		marshalledUser, _ := json.Marshal(user)

		var responseMap map[string]interface{}

		_ = json.Unmarshal(marshalledUser, &responseMap)
		responseMap["token"] = token

		utils.RespondJSON(w, http.StatusOK, responseMap)
		return
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

// SignUp returns the handler of the POST /signup endpoint
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
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err = CreateUser(db, user)
		if err, ok := err.(shared.PublicError); ok {
			showableError := err.(shared.ShowableError)
			utils.RespondWithError(w, showableError.StatusCode, showableError.Message)
			return
		}

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "unexpected error")
			return
		}

		user.Password = ""
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func validateMechanicSignUpFields(mechanic mec.Mechanic) error {
	if mechanic.Name == "" || mechanic.LastName == "" || mechanic.Email == "" || mechanic.NationalID == "" || mechanic.PhoneNumber == "" || mechanic.Password == "" {
		return ErrMissingFields
	}

	return nil
}

func MechanicLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mechanic := mec.Mechanic{}
		err := json.NewDecoder(r.Body).Decode(&mechanic)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if mechanic.Email == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "missing email")
			return
		}

		if mechanic.Password == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "missing password")
			return
		}

		retrievedMechanic, err := mec.GetMechanicByEmail(db, mechanic.Email)
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid email")
			return
		}

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(retrievedMechanic.Password), []byte(mechanic.Password))
		if !isPasswordCorrect(mechanic.Password, retrievedMechanic.Password) {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid password")
			return
		}

		retrievedMechanic.Password = ""

		token, err := GenerateMechanicToken(*retrievedMechanic)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		responseMap := map[string]interface{}{}
		responseMap["token"] = token
		responseMap["mechanic"] = retrievedMechanic

		utils.RespondJSON(w, http.StatusCreated, responseMap)
		return
	}
}

func MechanichSignUp(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mechanic := mec.Mechanic{}

		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&mechanic)
		if err != nil {
			fmt.Println(err)
			utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		err = validateMechanicSignUpFields(mechanic)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(mechanic.Password), 10)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		mechanic.Password = string(hashedPassword)
		err = mec.InsertMechanic(db, mechanic)
		if err, ok := err.(shared.PublicError); ok {
			showableError := err.(shared.ShowableError)
			utils.RespondWithError(w, showableError.StatusCode, showableError.Message)
			return
		}

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}
}
