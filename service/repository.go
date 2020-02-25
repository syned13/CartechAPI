package service

import (
	"database/sql"
	"fmt"
)

func GetAllServiceCategories(db *sql.DB) ([]ServiceCategory, error) {
	query := "SELECT * FROM service_category_table"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("error_while_executing_query: ", err.Error())
		return nil, err
	}

	serviceCategories := []ServiceCategory{}

	for rows.Next() {
		serviceCategory := ServiceCategory{}
		err := rows.Scan(&serviceCategory.ServiceCategoryID, &serviceCategory.ServiceCategory)
		if err != nil {
			fmt.Println("error_while_scanning_row: ", err.Error())
			return nil, err
		}

		serviceCategories = append(serviceCategories, serviceCategory)
	}

	return serviceCategories, nil
}
