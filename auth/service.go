package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	mec "github.com/CartechAPI/mechanic"
	usr "github.com/CartechAPI/user"
	"github.com/CartechAPI/utils"
	jwt "github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

// ErrMissingToken missing token
var ErrMissingToken = errors.New("missing token")

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

// GenerateToken returns the jwt for the user logged in
func GenerateToken(user *usr.User) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
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
