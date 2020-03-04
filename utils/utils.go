package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	usr "github.com/CartechAPI/user"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidTokenSigningMethod = errors.New("invalid token signing method")
)

type errorResponse struct {
	Message string `json:"message"`
}

// RespondWithError responds with a json with the given status code and message
func RespondWithError(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{errorMessage})
}

// RespondJSON responds with a json with the given status code and data
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// DecodeToken decodes the token and returns the claims
func DecodeToken(authToken string) *usr.User {
	token, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidTokenSigningMethod
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		fmt.Println("error_parsing_jwt: ", err.Error())
		return nil
	}

	if !token.Valid {
		return nil
	}

	var user usr.User

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		claimsBytes, _ := json.Marshal(claims)
		json.Unmarshal(claimsBytes, &user)
		return &user
	}

	return nil
}
