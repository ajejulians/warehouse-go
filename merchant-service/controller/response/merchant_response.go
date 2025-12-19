package response

import "warehouse-go/merchant-service/pkg/pagination"

type MerchantResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	Photo        string `json:"photo"`
	Phone        string `json:"phone"`
	KeeperID     uint   `json:"keeper_id"`
	KeepersName  string `json:"keepers_name"`
	ProductCount int    `json:"product_count"`
}

type MerchantWithProductResponse struct {
	ID               uint              `json:"id"`
	Name             string            `json:"name"`
	Address          string            `json:"address"`
	Photo            string            `json:"photo"`
	Phone            string            `json:"phone"`
	KeeperID         uint              `json:"keeper_id"`
	KeepersName      string            `json:"keepers_name"`
	MerchantProducts []MerchantProduct `json:"merchant_products"`
}

type MerchantPaginationResponse struct {
	Message    string                        `json:"message"`
	Data       []MerchantResponse            `json:"data"`
	Pagination pagination.PaginationResponse `json:"pagination"`
}

type UploadResponse struct {
	URL      string `json:"url"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
}
