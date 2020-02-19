package auth

import (
	"database/sql"

	usr "github.com/CartechAPI/user"

	"golang.org/x/crypto/bcrypt"
)

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
