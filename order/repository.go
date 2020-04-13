package order

import (
	"database/sql"
	"fmt"
)

func insertServiceOrder(db *sql.DB, serviceOrder ServiceOrder) error {
	query := "INSERT INTO service_order_table (service_id, user_id, created_at, status) VALUES ($1, $2, NOW(), $3)"
	_, err := db.Exec(query, serviceOrder.ServiceID, serviceOrder.UserID, serviceOrder.Status)
	if err != nil {
		fmt.Println("error inserting into service_order: " + err.Error())
		return err
	}

	return nil
}

func getServiceOrderByUserIDAndStatus(db *sql.DB, userID int) ([]ServiceOrder, error) {
	query := "SELECT * FROM service_order_table WHERE user_id = $1 AND (status = 'pending' OR status = 'in_progress')"
	rows, err := db.Query(query, userID)
	if err != nil {
		fmt.Println("error while selecting from service_order_table by user_id and status: " + err.Error())
		return nil, err
	}

	serviceOrders := []ServiceOrder{}
	for rows.Next() {
		serviceOrder := ServiceOrder{}

		// TODO: handle this null types
		var mechanicID sql.NullInt64
		var startedAt, finishedAt, cancelledAt sql.NullTime
		err = rows.Scan(&serviceOrder.ServiceOrderID, &serviceOrder.ServiceID, &serviceOrder.UserID, &mechanicID, &serviceOrder.CreatedAt, &startedAt, &serviceOrder.Status, &finishedAt, &cancelledAt)
		if err != nil {
			fmt.Println("error while scanning rows from service_order_table by user_id and status: " + err.Error())
			return nil, err
		}

		serviceOrders = append(serviceOrders, serviceOrder)
	}

	return serviceOrders, nil
}
