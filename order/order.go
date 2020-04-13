package order

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/CartechAPI/shared"
)

var (
	// ErrMissingNewValue missing new value
	ErrMissingNewValue = errors.New("missing new value")
)

func updateServiceOrder(db *sql.DB, serviceOrderID int, patchRequest shared.PatchRequestBody) error {
	for _, updateOp := range patchRequest {
		if updateOp.Op == shared.PatchOpReplace {
			err := replaceOnServiceOrder(db, serviceOrderID, updateOp.Path, updateOp.Value)
			// TODO: what happens if the err occurs on the second or third updateOP? Should the message be specific on one updateOp?
			if err == ErrNoRowsAffected {
				return shared.NewShowableError("resource not found", http.StatusNotFound)
			}

			if err == ErrMissingNewValue {
				return shared.NewBadRequestError("missing new value on replace operation")
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func replaceOnServiceOrder(db *sql.DB, serviceOrderID int, toReplace string, newValue string) error {
	if newValue == "" {
		return ErrMissingNewValue
	}

	if toReplace == "status" {
		return updateServiceOrderStatus(db, serviceOrderID, ServiceOrderStatus(newValue))
	}

	return nil
}
