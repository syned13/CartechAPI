package order

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/CartechAPI/shared"
	"github.com/CartechAPI/utils"
	"github.com/streadway/amqp"
)

var (
	// ErrMissingNewValue missing new value
	ErrMissingNewValue = errors.New("missing new value")
	// ErrMissingUserID missing user id
	ErrMissingUserID = shared.NewBadRequestError("missing user id")
	// ErrMissingServiceID missing service id
	ErrMissingServiceID = shared.NewBadRequestError("missing service id")
	// ErrMissingLatitude missing latitude
	ErrMissingLatitude = shared.NewBadRequestError("missing service location latitude")
	// ErrMissingLongitude missing service location longitude
	ErrMissingLongitude = shared.NewBadRequestError("missing service location longitude")
	// ErrMultipleServiceOrders multiple service orders
	ErrMultipleServiceOrders = shared.NewBadRequestError("user is not allowed to have more than one service")
	// ErrInvalidStatus invalid status
	ErrInvalidStatus = shared.NewBadRequestError("invalid status")
)

// AssignerQueue assigner queue
const AssignerQueue = "assign-order"

func validateServiceOrderFields(serviceOrder ServiceOrder) error {
	if serviceOrder.UserID == 0 {
		return ErrMissingUserID
	}

	if serviceOrder.ServiceID == 0 {
		return ErrMissingServiceID
	}

	if serviceOrder.Lat == 0 {
		return ErrMissingLatitude
	}

	if serviceOrder.Lng == 0 {
		return ErrMissingLongitude
	}

	return nil
}

func createServiceOrder(db *sql.DB, channel *amqp.Channel, serviceOrder *ServiceOrder) error {
	err := validateServiceOrderFields(*serviceOrder)
	if err != nil {
		return err
	}

	// serviceOrders, err := getServiceOrderByUserIDAndStatus(db, serviceOrder.UserID)
	// if err != nil {
	// 	return err
	// }

	// if len(serviceOrders) > 0 {
	// 	return ErrMultipleServiceOrders
	// }

	serviceOrder.Status = ServiceOrderStatusPending
	err = insertServiceOrder(db, *serviceOrder)
	if err != nil {
		return err
	}

	err = assignOrder(channel, *serviceOrder)
	if err != nil {
		return err
	}

	return nil
}

func assignOrder(channel *amqp.Channel, order ServiceOrder) error {
	marshalledOrder, err := json.Marshal(order)
	if err != nil {
		log.Println("failed_to_marshall_order:" + err.Error())
		return err
	}

	err = channel.Publish(
		"",
		AssignerQueue,
		false,
		false,
		amqp.Publishing{
			Body: marshalledOrder,
		},
	)

	if err != nil {
		log.Println("failed_to_publish_message_assign_order: " + err.Error())
		return err
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

func isServiceOrderStatusValid(status ServiceOrderStatus) bool {
	if _, ok := ValidServiceOrderStatus[status]; ok {
		return true
	}

	return false
}

func getAllServiceOrders(db *sql.DB, token string, status string) ([]ServiceOrder, error) {
	if !isServiceOrderStatusValid(ServiceOrderStatus(status)) && status != "" {
		return nil, ErrInvalidStatus
	}

	var err error
	serviceOrders := []ServiceOrder{}

	if status != "" {
		serviceOrders, err = selectAllOrdersByStatus(db, ServiceOrderStatus(status))
		if err != nil {
			return nil, err
		}

		return serviceOrders, nil
	}

	// TODO: move this to another function to filter by user
	clientType, id, err := utils.DecodeToken(token)
	if err != nil {
		return nil, err
	}

	if clientType == shared.ClientTypeUser {
		serviceOrders, err = selectAllOrdersFromUser(db, id)
		if err != nil {
			return nil, err
		}
	}

	if clientType == shared.ClientTypeAdmin {
		serviceOrders, err = selectAllOrders(db)
		if err != nil {
			return nil, err
		}
	}

	return serviceOrders, nil
}

func getAllPastServiceOrders(db *sql.DB, clientType shared.ClientType, id int) ([]ServiceOrder, error) {
	if clientType == shared.ClientTypeUser {
		serviceOrders, err := selectAllPastServiceOrdersByUser(db, id)
		if err != nil {
			return nil, err
		}

		return serviceOrders, nil
	}

	return nil, errors.New("not yet implemented")
}

func getAllCurrentOrders(db *sql.DB, clientType shared.ClientType, id int) ([]ServiceOrder, error) {
	if clientType == shared.ClientTypeUser {
		serviceOrders, err := selectAllCurrentOrdersByUser(db, id)
		if err != nil {
			return nil, err
		}

		return serviceOrders, nil
	}

	return nil, errors.New("not yet implemented")
}
