package mechanic

// Mechanic represents a mechanic
type Mechanic struct {
	MechanicID  int     `json:"mechanic_id"`
	Name        string  `json:"name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email"`
	NationalID  string  `json:"national_id"`
	Password    string  `json:"password"`
	Score       float32 `json:"score"`
	Bio         string  `json:"bio"`
	PhoneNumber string  `json:"phone_number"`
}
