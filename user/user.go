package user

// User represents an user
type User struct {
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}
