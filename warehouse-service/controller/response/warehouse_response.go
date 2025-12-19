package response

import "warehouse-go/warehouse-service/pkg/pagination"

type WarehouseResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	Photo        string `json:"photo"`
	Phone        string `json:"phone"`
	CountProduct int    `json:"count_product"`
}

type GetAllWarehouseResponse struct {
	Warehouse  []WarehouseResponse 				`json:"warehouse"`
	Pagination pagination.PaginationResponse	`json:"pagination"`
}

type DetailWarehouseResponse struct {
	ID		uint 	`json:"id"`
	Name	string	`json:"name"`
	Address string 	`json:"address"`
	Photo	string 	`json:"photo"`
	Phone 	string  `json:"phone"`
	CountProduct int `json:"count_product"`
	WarehouseProducts 	[]WarehouseProductResponse	`json:"products"`
}