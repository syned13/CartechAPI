package user

import (
	"database/sql"
	"fmt"

	"github.com/CartechAPI/shared"
	"github.com/lib/pq"
)

var (
	// ErrNotUniqueField not unique field
	ErrNotUniqueField = shared.NewBadRequestError("not unique fields")
	// ErrNotUniqueEmail not unique email
	ErrNotUniqueEmail = shared.NewBadRequestError("not unique email")
	//ErrNotUniquePhoneNumber not unique phone number
	ErrNotUniquePhoneNumber = shared.NewBadRequestError("not unique phone number")
)

var uniqueConstraintsErrs = map[string]error{
	"user_table_email_key":        ErrNotUniqueEmail,
	"user_table_phone_number_key": ErrNotUniquePhoneNumber,
}

const uniqueViolationCode = "23505"

// GetUserByEmail searchs for an user by its email and returns it
func GetUserByEmail(db *sql.DB, username string) (*User, error) {
	user := User{}

	fmt.Println(username)
	query := "SELECT * FROM user_table WHERE email = $1;"
	err := db.QueryRow(query, username).Scan(&user.UserID, &user.Name, &user.LastName, &user.Email, &user.Password, &user.PhoneNumber)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("failed_to_get_user: " + err.Error())
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, nil
}

// InsertUser inserts a new user into the user table
func InsertUser(db *sql.DB, user *User) (*User, error) {
	query := "INSERT INTO user_table (name, last_name, email, password, phone_number) VALUES($1, $2, $3, $4, $5) RETURNING user_id;"
	err := db.QueryRow(query, user.Name, user.LastName, user.Email, user.Password, user.PhoneNumber).Scan(&user.UserID)
	if err != nil {
		pqError := err.(*pq.Error)
		if pqError.Code == uniqueViolationCode {
			if err, ok := uniqueConstraintsErrs[pqError.Constraint]; ok {
				return nil, err
			}

			return nil, ErrNotUniqueField
		}

		return nil, err
	}

	fmt.Println(user)
	return user, nil
}
