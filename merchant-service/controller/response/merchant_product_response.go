package response

import (
	"time"
	"warehouse-go/merchant-service/pkg/pagination"
)

type MerchantProductResponse struct {
	ID          uint      `json:"id:"`
	MerchantID  uint      `json:"merchant_id"`
	ProductID   uint      `json:"product_id"`
	Stock       int       `json:"stock"`
	WarehouseID uint      `json:"warehouse_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	//Relations
	Merchant MerchantResponse `json:"merchant,omitempty"`
}

type MerchantProduct struct {
	ID						uint 	`json:"id"`
	MerchantID  			uint    `json:"merchant_id"`
	ProductID   			uint    `json:"product_id"`
	ProductName 			string	`json:"product_name"`
	ProductAbout 			string	`json:"product_about"`
	ProductPhoto			string 	`json:"product_photo"`
	ProductPrice 			int 	`json:"product_price"`
	ProductCategory 		string 	`json:"product_category"`
	ProductCategoryphoto 	string 	`json:"product_category_photo"`
	Stock 					int 	`json:"stock"`
	WarehouseID 			uint    `json:"warehouse_id"`
	WarehouseName 			string 	`json:"warehouse_name"`
	WarehousePhoto 			string 	`json:"warehouse_photo"`
	WarehousePhone 			string 	`json:"warehouse_phone"`
}

type ProductTotalStockResponse struct {
	ProductID   uint      	`json:"product_id"`
	TotalStock 	int 		`json:"total_stock"`
}

type GetAllMerchantProductResponse struct {
	MerchantProducts 	[]MerchantProduct	`json:"merchant_products"`
	Pagination 			pagination.PaginationResponse 	`json:"pagination"`
}