package httpclient

import (
	"warehouse-go/merchant-service/controller/response"
)

func MapProductResponseToMerchantProduct(product *ProductResponse) *response.MerchantProduct{
	return &response.MerchantProduct{
		ID: product.ID,
		ProductID: product.ID,
		ProductName: product.Name,
		ProductAbout: product.About,
		ProductPhoto: product.Thumbnail,
		ProductPrice: int(product.Price),
		ProductCategory: product.Category.Name,
		ProductCategoryphoto: product.Category.Photo,
		Stock: 0,
		WarehouseID: 0,
		WarehouseName: "",
		WarehousePhoto: "",
		WarehousePhone: "",
	}
}

func MapWarehouseResponseToMerchantProduct(warehouse *WarehouseResponse) *response.MerchantProduct {
	return &response.MerchantProduct{
		WarehouseID: warehouse.ID,
		WarehouseName: warehouse.Name,
		WarehousePhoto: warehouse.Photo,
		WarehousePhone: warehouse.Phone,
	}
}