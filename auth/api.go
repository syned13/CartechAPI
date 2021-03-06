package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
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
		user := &us.User{}

		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid_request_body")
			return
		}

		token, user, err := login(db, user)
		if showableErr, ok := err.(shared.ShowableError); ok {
			utils.RespondWithError(w, showableErr.StatusCode, showableErr.Message)
			return
		}
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		marshalledUser, _ := json.Marshal(user)

		responseMap := make(map[string]interface{})
		if responseMap == nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "nil!")
			return
		}

		err = json.Unmarshal(marshalledUser, &responseMap)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		responseMap["token"] = token

		utils.RespondJSON(w, http.StatusOK, responseMap)
		return
	}
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

		user, err = signUp(db, user)
		if showableErr, ok := err.(shared.ShowableError); ok {
			utils.RespondWithError(w, showableErr.StatusCode, showableErr.Message)
			return
		}
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		user.Password = ""
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func validateMechanicSignUpFields(mechanic mec.Mechanic) error {
	if mechanic.Name == "" {
		return ErrMissingName
	}
	if mechanic.LastName == "" {
		return ErrMissingLastName
	}
	if mechanic.Email == "" {
		return ErrMissingEmail
	}
	if mechanic.PhoneNumber == "" {
		return ErrMissingPhoneNumber
	}
	if mechanic.Password == "" {
		return ErrMissingPassword
	}

	return nil
}

// MechanicLogin is the login function for the mechanic
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

		if !isPasswordCorrect(mechanic.Password, retrievedMechanic.Password) {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid password")
			return
		}

		retrievedMechanic.Password = ""

		token, err := GenerateToken(retrievedMechanic.MechanicID, shared.ClientTypeMechanic)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		responseMap := map[string]interface{}{}
		responseMap["token"] = token
		responseMap["mechanic"] = retrievedMechanic

		utils.RespondJSON(w, http.StatusOK, responseMap)
		return
	}
}

// MechanichSignUp the sign up for the mechanic
func MechanichSignUp(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mechanic := mec.Mechanic{}

		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&mechanic)
		if err != nil {
			log.Println(err)
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

func validateSessionFields(session Session) error {
	if session.UserID == 0 {
		return errors.New("missing user id")
	}

	if session.Token == "" {
		return errors.New("missing token")
	}

	if session.UserType == "" {
		return errors.New("missing user type")
	}

	return nil
}

// StoreSession handles the request for storing a session
func StoreSession(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := &Session{}

		err := json.NewDecoder(r.Body).Decode(session)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid user body")
			return
		}

		err = validateSessionFields(*session)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		session, err = saveSession(db, *session)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		utils.RespondJSON(w, http.StatusCreated, session)
	}
}
