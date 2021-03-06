package order

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/CartechAPI/shared"
	"github.com/lib/pq"
)

var (
	// ErrNoRowsAffected no rows affected
	ErrNoRowsAffected = errors.New("no rows affected")
)

func insertServiceOrder(db *sql.DB, serviceOrder ServiceOrder) (int, error) {
	query := `INSERT INTO service_order_table 
				(service_id, user_id, created_at, status, lat, lng) 
				VALUES ($1, $2, NOW(), $3, $4, $5) 
				RETURNING service_order_id`

	id := 0
	err := db.QueryRow(query, serviceOrder.ServiceID, serviceOrder.UserID, serviceOrder.Status, serviceOrder.Lat, serviceOrder.Lng).Scan(&id)
	if err != nil {
		log.Println("error inserting into service_order: " + err.Error())
		return 0, err
	}

	return id, nil
}

func getServiceOrderByUserIDAndStatus(db *sql.DB, userID int) ([]ServiceOrder, error) {
	query := `SELECT service_order_id, service_order_table.service_id, user_id, mechanic_id, created_at, started_at, status, finished_at, cancelled_at, lat, lng, display_name 
	FROM service_order_table 
	LEFT JOIN service_table ON service_order_table.service_id = service_table.service_id
	WHERE user_id = $1 AND (status = 'pending' OR status = 'in_progress')`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("error while selecting from service_order_table by user_id and status: " + err.Error())
		return nil, err
	}

	return scanServiceOrders(rows)
}

func getServiceOrderByID(db *sql.DB, serviceOrderID int) (*ServiceOrder, error) {
	query := "SELECT * FROM service_order_table WHERE service_order_id = $1"
	row := db.QueryRow(query, serviceOrderID)

	serviceOrder := ServiceOrder{}
	var mechanicID sql.NullInt64
	var startedAt, finishedAt, cancelledAt sql.NullTime
	var lat, lng sql.NullFloat64
	err := row.Scan(&serviceOrder.ServiceOrderID, &serviceOrder.ServiceID, &serviceOrder.UserID, &mechanicID, &serviceOrder.CreatedAt, &startedAt, &serviceOrder.Status, &finishedAt, &cancelledAt, &lat, &lng)
	if err == sql.ErrNoRows {
		return nil, shared.NewShowableError("not found", http.StatusNotFound) //TODO: make this error a constant in a shared package
	}

	if err != nil {
		log.Println("error scanning row on getServiceOrderByID: " + err.Error())
		return nil, err
	}

	if mechanicID.Valid {
		serviceOrder.MechanicID = (int)(mechanicID.Int64)
	}

	if startedAt.Valid {
		serviceOrder.StartedAt = &startedAt.Time
	}

	if finishedAt.Valid {
		serviceOrder.FinishedAt = &finishedAt.Time
	}

	if cancelledAt.Valid {
		serviceOrder.CancelledAt = &cancelledAt.Time
	}

	if lat.Valid && lng.Valid {
		serviceOrder.Lat = lat.Float64
		serviceOrder.Lng = lng.Float64
	}

	return &serviceOrder, nil
}

// TODO: paginate this
func selectAllOrdersFromUser(db *sql.DB, userID int) ([]ServiceOrder, error) {
	query := `SELECT service_order_id, service_order_table.service_id, user_id, mechanic_id, created_at, started_at, status, finished_at, cancelled_at, lat, lng, display_name
	FROM service_order_table 
	LEFT JOIN service_table ON service_order_table.service_id = service_table.service_id
	WHERE user_id = $1`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("error while selecting user service_order: " + err.Error())
		return nil, err
	}

	return scanServiceOrders(rows)
}

func selectAllOrdersFromMechanic(db *sql.DB, mechanicID int) ([]ServiceOrder, error) {
	query := `SELECT service_order_id, service_order_table.service_id, user_id, mechanic_id, created_at, started_at, status, finished_at, cancelled_at, lat, lng, display_name
	FROM service_order_table 
	LEFT JOIN service_table ON service_order_table.service_id = service_table.service_id
	WHERE mechanic_id = $1`

	rows, err := db.Query(query, mechanicID)
	if err != nil {
		log.Println("error while selecting mechanic service_order: " + err.Error())
		return nil, err
	}

	return scanServiceOrders(rows)
}

// TODO: paginate this
func selectAllOrders(db *sql.DB) ([]ServiceOrder, error) {
	query := `SELECT service_order_id, service_order_table.service_id, user_id, mechanic_id, created_at, started_at, status, finished_at, cancelled_at, lat, lng, display_name
	FROM service_order_table 
	LEFT JOIN service_table ON service_order_table.service_id = service_table.service_id`

	rows, err := db.Query(query)
	if err != nil {
		log.Println("error while selecting all service_order: " + err.Error())
		return nil, err
	}

	return scanServiceOrders(rows)
}

func selectAllOrdersByStatus(db *sql.DB, status ServiceOrderStatus) ([]ServiceOrder, error) {
	query := `SELECT service_order_id, service_order_table.service_id, user_id, mechanic_id, created_at, started_at, status, finished_at, cancelled_at, lat, lng, display_name
	FROM service_order_table 
	LEFT JOIN service_table ON service_order_table.service_id = service_table.service_id
	WHERE status = $1`

	rows, err := db.Query(query, status)
	if err != nil {
		log.Println("error while selecting all service_order by status: " + err.Error())
		return nil, err
	}

	return scanServiceOrders(rows)
}

func selectAllPastServiceOrdersByUser(db *sql.DB, userID int) ([]ServiceOrder, error) {
	query := fmt.Sprintf(`SELECT service_order_id, service_order_table.service_id, user_id, mechanic_id, created_at, started_at, status, finished_at, cancelled_at, lat, lng, display_name
	FROM service_order_table 
	LEFT JOIN service_table ON service_order_table.service_id = service_table.service_id
	WHERE user_id = $1 AND status = '%s';`, ServiceOrderStatusFinished)

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("error while selecting from service_order_table by user_id and status: " + err.Error())
		return nil, err
	}

	return scanServiceOrders(rows)
}

func selectAllCurrentOrdersByUser(db *sql.DB, userID int) ([]ServiceOrder, error) {
	query := fmt.Sprintf(`SELECT service_order_id, service_order_table.service_id, user_id, mechanic_id, created_at, started_at, status, finished_at, cancelled_at, lat, lng, display_name
	FROM service_order_table 
	LEFT JOIN service_table ON service_order_table.service_id = service_table.service_id 
	WHERE user_id = $1 AND (status = '%s' OR status = '%s');`, ServiceOrderStatusInProgress, ServiceOrderStatusPending)

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("error while selecting from service_order_table by user_id and status: " + err.Error())
		return nil, err
	}

	return scanServiceOrders(rows)
}

func scanServiceOrders(rows *sql.Rows) ([]ServiceOrder, error) {
	serviceOrders := []ServiceOrder{}
	for rows.Next() {
		serviceOrder := ServiceOrder{}

		var mechanicID sql.NullInt64
		var startedAt, finishedAt, cancelledAt sql.NullTime
		var lat, lng sql.NullFloat64
		var serviceName sql.NullString

		err := rows.Scan(&serviceOrder.ServiceOrderID, &serviceOrder.ServiceID, &serviceOrder.UserID, &mechanicID, &serviceOrder.CreatedAt, &startedAt, &serviceOrder.Status, &finishedAt, &cancelledAt, &lat, &lng, &serviceName)
		if err != nil {
			log.Println("error while scanning service_orders: " + err.Error())
			return nil, err
		}

		if mechanicID.Valid {
			serviceOrder.MechanicID = (int)(mechanicID.Int64)
		}

		if startedAt.Valid {
			serviceOrder.StartedAt = &startedAt.Time
		}

		if finishedAt.Valid {
			serviceOrder.FinishedAt = &finishedAt.Time
		}

		if cancelledAt.Valid {
			serviceOrder.CancelledAt = &cancelledAt.Time
		}

		if lat.Valid && lng.Valid {
			serviceOrder.Lat = lat.Float64
			serviceOrder.Lng = lng.Float64
		}

		if serviceName.Valid {
			serviceOrder.ServiceName = serviceName.String
		}

		serviceOrders = append(serviceOrders, serviceOrder)
	}

	return serviceOrders, nil
}

func updateServiceOrderStatus(db *sql.DB, serviceOrderID int, status ServiceOrderStatus) error {
	query := "UPDATE service_order_table SET status = $1 WHERE service_order_id = $2"
	result, err := db.Exec(query, string(status), serviceOrderID)
	if err != nil {
		log.Println("error updating service order status: " + err.Error())
		if pqErr, ok := err.(pq.Error); ok {
			log.Println(pqErr.Error())
		}

		return err
	}

	rowsAffectes, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	if rowsAffectes == 0 {
		log.Println("now rows affected")
		return ErrNoRowsAffected
	}

	return nil
}

func setOrderMechanic(db *sql.DB, orderID int, mechanicID int) error {
	query := `UPDATE service_order_table
			SET mechanic_id = $1, status = $2
			WHERE service_order_id = $3`

	result, err := db.Exec(query, mechanicID, ServiceOrderStatusInProgress, orderID)
	if err != nil {
		log.Println("assigning_mechanic_to_order_failed: " + err.Error())
		return err
	}

	rowsAffectes, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	if rowsAffectes == 0 {
		log.Println("now rows affected")
		return ErrNoRowsAffected
	}

	return nil
}
