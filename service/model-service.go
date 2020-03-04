package service

// Service is the representation of a mechanic service
type Service struct {
	ServiceID       int             `json:"service_id"`
	ServiceName     string          `json:"service_name"`
	ServiceCategory ServiceCategory `json:"service_category"`
}

// ServiceCategory is the representation of a service category
type ServiceCategory struct {
	ServiceCategory   string `json:"service_category"`
	ServiceCategoryID int    `json:"service_category_id"`
}
