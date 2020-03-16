package service

// Service is the representation of a mechanic service
type Service struct {
	ServiceID         int    `json:"service_id"`
	ServiceName       string `json:"service_name"`
	ServiceCategoryID int    `json:"service_category_id"`
}

// Category is the representation of a service category
type Category struct {
	ServiceCategory   string `json:"service_category"`
	ServiceCategoryID int    `json:"service_category_id"`
}
