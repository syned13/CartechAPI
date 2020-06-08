package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/CartechAPI/shared"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidTokenSigningMethod = errors.New("invalid token signing method")
	ErrInvalidToken              = errors.New("invalid token")
	ErrCouldNoGetClaims          = errors.New("could not get claims")
	ErrCouldNot                  = errors.New("could not")
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
func DecodeToken(authToken string) (shared.ClientType, int, error) {
	token, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidTokenSigningMethod
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		log.Println("error_parsing_jwt: ", err.Error())
		return "", 0, err
	}

	if !token.Valid {
		return "", 0, ErrInvalidToken
	}

	tokenClaims := shared.TokenClaims{}

	claimsBytes, err := json.Marshal(token.Claims)
	if err != nil {
		return "", 0, err
	}

	err = json.Unmarshal(claimsBytes, &tokenClaims)
	if err != nil {
		return "", 0, err
	}

	return tokenClaims.ClientType, tokenClaims.ID, nil
}
