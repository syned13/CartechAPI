package service

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/CartechAPI/auth"
	"github.com/CartechAPI/utils"
	"github.com/gorilla/mux"
)

// GetAllServices returns all the mechanic services
func GetAllServices(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.UserAuthenticationMiddleware(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "client is unauthorized to perform the request")
			return
		}

		serviceCategories, err := getAllServiceCategories(db)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "unexpected error")
			return
		}

		utils.RespondJSON(w, http.StatusOK, map[string]interface{}{"services": serviceCategories})
	}
}

// GetAllServiceCategories returns all the services categories
func GetAllServiceCategories(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.UserAuthenticationMiddleware(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "client is unauthorized to perform the request")
			return
		}

		serviceCategories, err := getAllServiceCategories(db)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "unexpected error")
		}

		utils.RespondJSON(w, http.StatusOK, map[string]interface{}{"categories": serviceCategories})
	}
}

// GetServicesByCategoryID returns all services within a category
func GetServicesByCategoryID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.UserAuthenticationMiddleware(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "client is unauthorized to perform the request")
			return
		}

		params := mux.Vars(r)
		categoryID, err := strconv.Atoi(params["category_id"])
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "invalid request param")
			return
		}

		category, err := getCategoryByID(db, categoryID)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.RespondWithError(w, http.StatusNotFound, "resource not found")
				return
			}

			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		services, err := getServicesByCategoryID(db, categoryID)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.RespondWithError(w, http.StatusNotFound, "resource not found")
				return
			}

			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		responseMap := map[string]interface{}{}
		responseMap["service_category_id"] = category.ServiceCategoryID
		responseMap["service_category"] = category.ServiceCategory
		responseMap["services"] = services

		utils.RespondJSON(w, http.StatusOK, responseMap)
	}
}
