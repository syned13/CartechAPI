package order

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/CartechAPI/auth"
	"github.com/CartechAPI/shared"
	"github.com/CartechAPI/utils"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

// CreateServiceOrder receives the request to create a service request
func CreateServiceOrder(db *sql.DB, channel *amqp.Channel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceOrder := &ServiceOrder{}
		// TODO: shouldnt we use the user id from the jwt?
		_, _, err := auth.UserAuthenticationMiddleware(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "client is unauthorized to perform the request")
			return
		}

		err = json.NewDecoder(r.Body).Decode(&serviceOrder)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		serviceOrder, err = createServiceOrder(db, channel, serviceOrder)
		if err, ok := err.(shared.PublicError); ok {
			showableError := err.(shared.ShowableError)
			utils.RespondWithError(w, showableError.StatusCode, showableError.Message)
			return
		}

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		utils.RespondJSON(w, http.StatusOK, serviceOrder)
	}
}

// GetAllServiceOrders handles the request for getting all service orders
func GetAllServiceOrders(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "client is unauthorized to perform the request")
			return
		}

		status := r.URL.Query().Get("status")

		serviceOrders, err := getAllServiceOrders(db, token, status)
		if err != nil {
			log.Println(err.Error())
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		utils.RespondJSON(w, http.StatusOK, serviceOrders)
	}
}

// GetAllPastServiceOrders returns the past service orders
func GetAllPastServiceOrders(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientType, id, err := auth.UserAuthenticationMiddleware(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "client is unauthorized to perform the request")
			return
		}

		serviceOrders, err := getAllPastServiceOrders(db, clientType, id)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		utils.RespondJSON(w, http.StatusOK, serviceOrders)
	}
}

// GetAllCurrentOrders handles the request for getting all current orders
func GetAllCurrentOrders(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientType, id, err := auth.UserAuthenticationMiddleware(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "client is unauthorized to perform the request")
			return
		}

		serviceOrders, err := getAllCurrentOrders(db, clientType, id)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		utils.RespondJSON(w, http.StatusOK, serviceOrders)
	}
}

// GetServiceOrder handles the request of a GET to a specific service order
func GetServiceOrder(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		serviceOrderID, err := strconv.Atoi(params["order_id"])
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid request param")
			return
		}
		serviceOrder, err := getServiceOrderByID(db, serviceOrderID)
		if showableErr, ok := err.(shared.ShowableError); ok {
			utils.RespondWithError(w, showableErr.StatusCode, "resource not found")
			return
		}

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		utils.RespondJSON(w, http.StatusOK, serviceOrder)
	}
}

func validatePatchRequestBodyFields(patchRequest shared.PatchRequestBody) error {
	if len(patchRequest) > 1 {
		return errors.New("only one op is permitted")
	}

	for _, request := range patchRequest {
		if request.Op == "" {
			return errors.New("missing patch operation")
		}

		if request.Path == "" {
			return errors.New("missing patch path")
		}
	}

	return nil
}

// UpdateServiceOrder handles the request of updating a service order
func UpdateServiceOrder(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		serviceOrderID, err := strconv.Atoi(params["order_id"])
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid request param")
			return
		}

		patchRequest := shared.PatchRequestBody{}
		err = json.NewDecoder(r.Body).Decode(&patchRequest)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		err = validatePatchRequestBodyFields(patchRequest)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = updateServiceOrder(db, serviceOrderID, patchRequest)
		if showableError, ok := err.(shared.ShowableError); ok {
			utils.RespondWithError(w, showableError.StatusCode, showableError.Message)
			return
		}

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}
}

// AssignMechanicToOrder assings a mechanic to an order
func AssignMechanicToOrder(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mechanicIDS := r.URL.Query().Get("mechanic_id")
		if mechanicIDS == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "missing mechanic id")
			return
		}

		mechanicID, err := strconv.Atoi(mechanicIDS)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid mechanic id")
			return
		}

		params := mux.Vars(r)
		orderID, err := strconv.Atoi(params["order_id"])
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid order id")
			return
		}

		err = assignMechanicToOrder(db, mechanicID, orderID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		utils.RespondJSON(w, 200, "ok")
	}
}
