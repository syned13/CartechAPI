package user

import (
	"database/sql"
)

func GetUserByEmail(db *sql.DB, username string) (*User, error) {
	user := User{}

	query := "SELECT * FROM user_table WHERE email = ($1);"
	err := db.QueryRow(query, username).Scan(&user.UserID, &user.Name, &user.LastName, &user.Email, &user.Password, &user.PhoneNumber)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, nil
}

func InsertUser(db *sql.DB, user *User) (*User, error) {
	query := "INSERT INTO user_table (name, last_name, email, password, phone_number) VALUES($1, $2, $3, $4, $5) RETURNING user_id;"
	err := db.QueryRow(query, user.Name, user.LastName, user.Email, user.Password, user.PhoneNumber).Scan(&user.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
