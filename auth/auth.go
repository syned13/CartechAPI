package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	mec "github.com/CartechAPI/mechanic"
	"github.com/CartechAPI/shared"
	us "github.com/CartechAPI/user"
	usr "github.com/CartechAPI/user"
	"github.com/CartechAPI/utils"
	jwt "github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrMissingToken missing token
	ErrMissingToken = shared.NewBadRequestError("missing token")
	// ErrMissingEmail missing email
	ErrMissingEmail = shared.NewBadRequestError("missing email")
	// ErrMissingName missing name
	ErrMissingName = shared.NewBadRequestError("missing name")
	// ErrMissingLastName missing last name
	ErrMissingLastName = shared.NewBadRequestError("missing last name")
	// ErrMissingPassword missing password
	ErrMissingPassword = shared.NewBadRequestError("missing password")
	// ErrMissingPhoneNumber missing phone number
	ErrMissingPhoneNumber = shared.NewBadRequestError("missing phone number")
	// ErrInvalidCredentials invalid credentials
	ErrInvalidCredentials = shared.NewBadRequestError("incorrect email or password")
)

func login(db *sql.DB, user us.User) (string, *usr.User, error) {
	if user.Email == "" {
		return "", nil, ErrMissingEmail
	}

	if user.Password == "" {
		return "", nil, ErrMissingPassword
	}

	userRetrieved, err := us.GetUserByEmail(db, user.Email)
	if err != nil && err != sql.ErrNoRows {
		return "", nil, err
	}

	// user does not exist
	if err == sql.ErrNoRows || userRetrieved == nil {
		return "", nil, ErrInvalidCredentials
	}

	if !isPasswordCorrect(user.Password, userRetrieved.Password) {
		return "", nil, ErrInvalidCredentials
	}

	token, err := GenerateTokenForUser(userRetrieved)
	if err != nil {
		return "", nil, err
	}

	userRetrieved.Password = ""

	return token, userRetrieved, nil
}

func validateSignUpFields(user us.User) error {
	if user.Email == "" {
		return ErrMissingEmail
	}

	if user.Name == "" {
		return ErrMissingPassword
	}

	if user.LastName == "" {
		return ErrMissingLastName
	}

	if user.Password == "" {
		return ErrMissingPassword
	}

	if user.PhoneNumber == "" {
		return ErrMissingPassword
	}

	return nil
}

func signUp(db *sql.DB, user *usr.User) (*us.User, error) {
	err := validateSignUpFields(*user)
	if err != nil {
		return nil, err
	}

	user, err = CreateUser(db, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateUser creates a new user and adds it to the databse
func CreateUser(db *sql.DB, user *usr.User) (*usr.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	user, err = usr.InsertUser(db, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GenerateTokenForUser returns the jwt for the user logged in
func GenerateTokenForUser(user *usr.User) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"type":    shared.ClientTypeUser,
		"user_id": user.UserID,
		"iat":     now.String(),
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		fmt.Println("error_signing_token: " + err.Error())
		return "", err
	}

	return signedToken, nil
}

// GenerateMechanicToken returns the jwt for the mechanic logged in
func GenerateMechanicToken(mechanic mec.Mechanic) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"mechanic_id": mechanic.MechanicID,
		"iat":         now.String(),
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		fmt.Println("error_signing_token: " + err.Error())
		return "", err
	}

	return signedToken, nil
}

// UserAuthenticationMiddleware middleware for the user's restricted endpoints
func UserAuthenticationMiddleware(r *http.Request) error {
	token := r.Header["Authorization"]
	if len(token) == 0 || token[0] == "" {
		return ErrMissingToken
	}

	user, err := utils.DecodeToken(token[0])
	if user == nil {
		return err
	}

	return nil
}

func isPasswordCorrect(enteredPassword string, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(enteredPassword), []byte(storedPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	}

	return true
}
