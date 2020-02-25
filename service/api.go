package service

import (
	"database/sql"
	"net/http"

	"github.com/CartechAPI/utils"
)

func GetAllServices(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header["Authorization"]
		if len(token) == 0 || token[0] == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "mising authentication token")
			return
		}

		user := utils.DecodeToken(token[0])
		if user == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "client is unauthorized to perform the request")
			return
		}

		serviceCategories, err := GetAllServiceCategories(db)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "unexpected error")
		}

		utils.RespondJSON(w, http.StatusOK, map[string]interface{}{"services": serviceCategories})
	}
}
