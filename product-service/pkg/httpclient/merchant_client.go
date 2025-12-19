package httpclient

import (
	"net/http"
	"time"
	"warehouse-go/product-service/configs"
)

type MerchantClient struct {
	urlMerchantService string
	httpClient         *http.Client
}

type MerchantResponse struct {
	ProductID uint `json:"product_id"`
	TotalStock int `json:"total_stock"`	
}

type MerchantProductServiceResponse struct {
	Message string `json:"message"`
	Data 	MerchantResponse `json:"data"`
	Error	string `json:"error,omitempty"`
}

func NewMerchantClient(cfg configs.Config) *MerchantClient {
	return &MerchantClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		urlMerchantService: cfg.App.UrlMerchantService,
	}
}


