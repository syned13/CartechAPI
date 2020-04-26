package order

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/CartechAPI/shared"
	"github.com/streadway/amqp"
)

var (
	// ErrMissingNewValue missing new value
	ErrMissingNewValue = errors.New("missing new value")
)

const ASSIGNER_QUEUE = "assign-order"

func assignOrder(channel *amqp.Channel, order ServiceOrder) error {
	marshalledOrder, err := json.Marshal(order)
	if err != nil {
		fmt.Println("failed_to_marshall_order:" + err.Error())
		return err
	}

	err = channel.Publish(
		"",
		ASSIGNER_QUEUE,
		false,
		false,
		amqp.Publishing{
			Body: marshalledOrder,
		},
	)

	if err != nil {
		fmt.Println("failed_to_publish_message_assign_order:" + err.Error())
	}

	return nil
}

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
