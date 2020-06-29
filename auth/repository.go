package auth

import "database/sql"

func saveSession(db *sql.DB, session Session) (*Session, error) {
	query := `INSERT INTO sessions 
	(created_at, user_id, user_type, device_token)
	VALUES (NOW(), $1, $2, $3) RETURNING session_id, created_at`

	err := db.QueryRow(query, session.UserID, session.UserType, session.Token).Scan(&session.SessionID, &session.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func GetLastSession(db *sql.DB, id int) (*Session, error) {
	query := `SELECT * FROM sessions WHERE user_id = $1 ORDER BY created_at LIMIT 1`

	session := Session{}
	err := db.QueryRow(query, id).Scan(&session.SessionID, &session.CreatedAt, &session.UserID, &session.UserType, &session.Token)

	return &session, err
}
