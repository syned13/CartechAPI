package order

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/CartechAPI/auth"
	"github.com/CartechAPI/shared"
	"github.com/CartechAPI/utils"
)

func validateServiceOrderFields(serviceOrder ServiceOrder) error {
	if serviceOrder.UserID == 0 {
		return shared.NewBadRequestError("missing user id")
	}

	if serviceOrder.ServiceID == 0 {
		return shared.NewBadRequestError("missing service id")
	}

	return nil
}

// CreateServiceOrder receives the request to create a service request
func CreateServiceOrder(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceOrder := ServiceOrder{}
		err := auth.UserAuthenticationMiddleware(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "client is unauthorized to perform the request")
			return
		}

		err = json.NewDecoder(r.Body).Decode(&serviceOrder)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		err = validateServiceOrderFields(serviceOrder)
		if err, ok := err.(shared.PublicError); ok {
			showableError := err.(shared.ShowableError)
			utils.RespondWithError(w, showableError.StatusCode, showableError.Message)
			return
		}

		serviceOrders, err := getServiceOrderByUserIDAndStatus(db, serviceOrder.UserID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal sever error")
			return
		}

		if len(serviceOrders) > 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "user is not allowed to have more than one service")
			return
		}

		serviceOrder.Status = ServiceOrderStatusPending
		err = insertServiceOrder(db, serviceOrder)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "internal sever error")
			return
		}

		utils.RespondJSON(w, http.StatusOK, serviceOrder)
	}
}
