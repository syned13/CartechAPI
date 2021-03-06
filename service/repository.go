package service

import (
	"database/sql"
	"log"
)

func getServicesByCategoryID(db *sql.DB, categoryID int) ([]Service, error) {
	query := "SELECT * FROM service_table WHERE service_category_id = $1"
	rows, err := db.Query(query, categoryID)
	if err != nil {
		log.Println("error_while_executing_query: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	services := []Service{}

	for rows.Next() {
		service := Service{}
		err = rows.Scan(&service.ServiceID, &service.ServiceName, &service.ServiceCategoryID)
		if err != nil {
			log.Println("error_while_scanning_row_service_table: ", err.Error())
			return nil, err
		}

		services = append(services, service)
	}

	return services, nil
}

func getCategoryByID(db *sql.DB, categoryID int) (*Category, error) {
	query := "SELECT * FROM service_category_table WHERE service_category_id = $1"
	rows, err := db.Query(query, categoryID)
	if err != nil {
		log.Println("error_while_executing_query: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	category := Category{}

	rows.Next()
	err = rows.Scan(&category.ServiceCategoryID, &category.ServiceCategory)
	if err != nil {
		log.Println("error_while_scanning_row_service_category_table: ", err.Error())
		return nil, err
	}

	return &category, nil
}

// GetAllServiceCategories returns all the categories of services
func getAllServiceCategories(db *sql.DB) ([]Category, error) {
	query := "SELECT * FROM service_category_table"
	rows, err := db.Query(query)
	if err != nil {
		log.Println("error_while_executing_query: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	serviceCategories := []Category{}

	for rows.Next() {
		serviceCategory := Category{}
		err = rows.Scan(&serviceCategory.ServiceCategoryID, &serviceCategory.ServiceCategory)
		if err != nil {
			log.Println("error_while_scanning_row_service_category_table: ", err.Error())
			return nil, err
		}

		serviceCategories = append(serviceCategories, serviceCategory)
	}

	return serviceCategories, nil
}

func scanServices(rows *sql.Rows) ([]Service, error) {
	services := []Service{}

	var err error
	for rows.Next() {
		service := Service{}
		err = rows.Scan(&service.ServiceID, &service.ServiceName, &service.ServiceCategoryID)
		if err != nil {
			log.Println("error_while_scanning_row_service_table: ", err.Error())
			return nil, err
		}

		services = append(services, service)
	}

	return services, nil
}
