package mechanic

import (
	"database/sql"
	"log"

	"github.com/CartechAPI/shared"
	"github.com/lib/pq"
)

const uniqueViolationCode = "23505"

var (
	// ErrNotUniqueField not unique field
	ErrNotUniqueField = shared.NewBadRequestError("not unique fields")
	// ErrNotUniqueEmail not unique email
	ErrNotUniqueEmail = shared.NewBadRequestError("not unique email")
	//ErrNotUniqueNationalID not unique national id
	ErrNotUniqueNationalID = shared.NewBadRequestError("not unique national id")
	//ErrNotUniquePhoneNumber not unique phone number
	ErrNotUniquePhoneNumber = shared.NewBadRequestError("not unique phone number")
)

var uniqueConstraintsErrs = map[string]error{
	"mechanic_table_email_key":        ErrNotUniqueEmail,
	"mechanic_table_national_id_key":  ErrNotUniqueNationalID,
	"mechanic_table_phone_number_key": ErrNotUniquePhoneNumber,
}

// InsertMechanic creates a new mechanic on the database
func InsertMechanic(db *sql.DB, mechanic Mechanic) error {
	query := "INSERT INTO mechanic_table (name, last_name, email, national_id, password, phone_number) VALUES($1,$2,$3,$4,$5,$6) RETURNING mechanic_id"

	err := db.QueryRow(query, mechanic.Name, mechanic.LastName, mechanic.Email, mechanic.NationalID, mechanic.Password, mechanic.PhoneNumber).Scan(&mechanic.MechanicID)
	if err != nil {
		pqError := err.(*pq.Error)
		if pqError.Code == uniqueViolationCode {
			if err, ok := uniqueConstraintsErrs[pqError.Constraint]; ok {
				return err
			}

			return ErrNotUniqueField
		}

		return err
	}

	return nil
}

// GetMechanicByEmail returns a mechanig given its email
func GetMechanicByEmail(db *sql.DB, email string) (*Mechanic, error) {
	query := "SELECT * FROM mechanic_table WHERE email = $1"

	mechanic := Mechanic{}
	err := db.QueryRow(query, email).Scan(&mechanic.MechanicID, &mechanic.Name, &mechanic.LastName, &mechanic.Email, &mechanic.NationalID, &mechanic.Password, &mechanic.Score, &mechanic.Bio, &mechanic.PhoneNumber)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &mechanic, nil
}
