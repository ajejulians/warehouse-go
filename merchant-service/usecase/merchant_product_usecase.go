package usecase

import (
	"context"
	"errors"
	"time"
	"warehouse-go/merchant-service/model"
	"warehouse-go/merchant-service/pkg/httpclient"
	"warehouse-go/merchant-service/pkg/rabbitmq"
	"warehouse-go/merchant-service/repository"

	"github.com/gofiber/fiber/v2/log"
)

type MerchantProductUsecaseInterface interface {
	CreateMerchantProduct(ctx context.Context, merchantProduct *model.MerchantProduct) error
	GetMerchantProductByID(ctx context.Context, id uint) (*model.MerchantProduct, *httpclient.ProductResponse, *httpclient.WarehouseResponse, error)
	GetMerchantProducts(ctx context.Context, page, limit int, search, sortBy, sortOrder string, merchantID uint, productID uint) ([]model.MerchantProduct, []httpclient.ProductResponse, []httpclient.WarehouseResponse, int64, error)
	GetMerchantProductByBarcode(ctx context.Context, barcode string, merchantID uint) (*model.MerchantProduct, *httpclient.ProductResponse, *httpclient.WarehouseResponse, error)
	UpdateMerchantProduct(ctx context.Context, merchantProuduct *model.MerchantProduct) error
	DeleteMerchantProduct(ctx context.Context, id uint) error
	DeleteAllProductMerchantProducts(ctx context.Context, productID uint) error

	GetProductTotalStock(ctx context.Context, productID uint) (int, error)
}

type merchantProductUsecase struct {
	merchantProductRepo repository.MerchantProductRepositoryInterface
	productClient       httpclient.ProductClientInterface
	warehouseClient     httpclient.WarehouseClientInterface
	rabbitMQService     *rabbitmq.RabbitMQService
}

// GetMerchantProductByBarcode implements MerchantProductUsecaseInterface.
func (m *merchantProductUsecase) GetMerchantProductByBarcode(ctx context.Context, barcode string, merchantID uint) (*model.MerchantProduct, *httpclient.ProductResponse, *httpclient.WarehouseResponse, error) {
	product, err := m.productClient.GetProductByBarcode(ctx, barcode)
	if err != nil {
		log.Errorf("[MerchantProductUsecase] GetMerchantProductByBarcode - 1: %v", err)
		return nil, nil, nil, err
	}

	merchantProduct, err := m.merchantProductRepo.GetMerchantProductByProductIDAndMerchantID(ctx, product.ID, merchantID)
	if err != nil {
		log.Errorf("[MerchantProductUsecase] GetMerchantProductByBarcode - 2: %v", err)
		return nil, nil, nil, err 
	}

	warehouse, err := m.warehouseClient.GetWarehouseByID(ctx, merchantProduct.WarehouseID)
	if err != nil {
		log.Errorf("[MerchantProductUsecase] GetMerchantProductByBarcode - 2: %v", err)
		return nil, nil, nil, err 
	}

	return merchantProduct, product, warehouse, nil
}

// CreateMerchantProduct implements MerchantProductUsecaseInterface.
func (m *merchantProductUsecase) CreateMerchantProduct(ctx context.Context, merchantProduct *model.MerchantProduct) error {
	warehouseProductStock, err := m.warehouseClient.GetWarehouseProductStock(ctx, merchantProduct.WarehouseID, merchantProduct.ProductID)
	if err != nil {
		log.Errorf("[MerchantProductUsecase] CreateMerchantProduct - 1: %v", err)
		return err
	}

	if warehouseProductStock.Stock < merchantProduct.Stock {
		log.Errorf("[MerchantProductUsecase] CreateMerchantProduct - 2: %v", err)
		return errors.New("stock not enough")
	}

	if err := m.merchantProductRepo.CreateMerchantProduct(ctx, merchantProduct); err != nil {
		log.Errorf("[MerchantProductUsecase] CreateMerchantProduct - 3: %v", err)
		return err
	}

	stockReductionEvent := rabbitmq.StockReductionEvent{
		WarehouseID: merchantProduct.WarehouseID,
		ProductID:   merchantProduct.ProductID,
		Stock:       merchantProduct.Stock,
		MerchantID:  merchantProduct.MerchantID,
		Timestamp:   time.Now(),
	}

	if err := m.rabbitMQService.PublishStockReductionEvent(ctx, stockReductionEvent); err != nil {
		log.Errorf("[MerchantProductUsecase] CreateMerchantProduct - 3: %v", err)
	}
	return nil
}

// DeleteAllProductMerchantProducts implements MerchantProductUsecaseInterface.
func (m *merchantProductUsecase) DeleteAllProductMerchantProducts(ctx context.Context, productID uint) error {
	return m.merchantProductRepo.DeleteAllProductMerchantProducts(ctx, productID)
}

// DeleteMerchantProduct implements MerchantProductUsecaseInterface.
func (m *merchantProductUsecase) DeleteMerchantProduct(ctx context.Context, id uint) error {
	return m.merchantProductRepo.DeleteMerchantProduct(ctx, id)
}

// GetMerchantProductByID implements MerchantProductUsecaseInterface.
func (m *merchantProductUsecase) GetMerchantProductByID(ctx context.Context, id uint) (*model.MerchantProduct, *httpclient.ProductResponse, *httpclient.WarehouseResponse, error) {
	merchantProduct, err := m.merchantProductRepo.GetMerchantProductByID(ctx, id)
	if err != nil {
		log.Errorf("[MerchantProductUsecase] GetMerchantProductByID - 1: %v", err)
		return nil, nil, nil, err
	}

	product, err := m.productClient.GetProductByID(ctx, merchantProduct.ProductID)
	if err != nil {
		log.Errorf("[MerchantProductUsecase] GetMerchantProductByID - 2: %v", err)
		return nil, nil, nil, err
	}

	warehouse, err := m.warehouseClient.GetWarehouseByID(ctx, merchantProduct.WarehouseID)
	if err != nil {
		log.Errorf("[MerchantProductUsecase] GetMerchantProductByID - 2: %v", err)
		return nil, nil, nil, err
	}

	return merchantProduct, product, warehouse, nil

}

// GetMerchantProducts implements MerchantProductUsecaseInterface.
func (m *merchantProductUsecase) GetMerchantProducts(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string, merchantID uint, productID uint) ([]model.MerchantProduct, []httpclient.ProductResponse, []httpclient.WarehouseResponse, int64, error) {
	merchantProducts, total, err := m.merchantProductRepo.GetMerchantProducts(ctx, page, limit, search, sortBy, sortOrder, merchantID, productID)
	if err != nil {
		log.Errorf("[MerchantProductUsecase] GetMerchantProducts - 1: %v", err)
		return nil, nil, nil, 0, err
	}

	var products []httpclient.ProductResponse
	var warehouses []httpclient.WarehouseResponse

	if len(merchantProducts) > 0 {
		for _, mp := range merchantProducts {
			product, err := m.productClient.GetProductByID(ctx, mp.ProductID)
			if err != nil {
				log.Errorf("[MerchantProductUsecase] GetMerchantProducts - 2: %v", err)
				return nil, nil, nil, 0, err
			}
			products = append(products, *product)

			warehouse, err := m.warehouseClient.GetWarehouseByID(ctx, mp.WarehouseID)
			if err != nil {
				log.Errorf("[MerchantProductUsecase] GetMerchantProducts - 3: %v", err)
				return nil, nil, nil, 0, err
			}
			warehouses = append(warehouses, *warehouse)
		}
	}

	return merchantProducts, products, warehouses, total, nil
}

// GetProductTotalStock implements MerchantProductUsecaseInterface.
func (m *merchantProductUsecase) GetProductTotalStock(ctx context.Context, productID uint) (int, error) {
	return m.merchantProductRepo.GetProductTotalStock(ctx, productID)
}

// UpdateMerchantProduct implements MerchantProductUsecaseInterface.
func (m *merchantProductUsecase) UpdateMerchantProduct(ctx context.Context, merchantProuduct *model.MerchantProduct) error {
	warehouseProductStock, err := m.warehouseClient.GetWarehouseProductStock(ctx, merchantProuduct.WarehouseID, merchantProuduct.ProductID)
		if err != nil {
		log.Errorf("[MerchantProductUsecase] UpdateMerchantProduct - 1: %v", err)
		return err
	}

	if warehouseProductStock.Stock < merchantProuduct.Stock {
		log.Errorf("[MerchantProductUsecase] UpdateMerchantProduct - 2: %v", err)
		return errors.New("stock not enough")
	}

	return m.merchantProductRepo.UpdateMerchantProduct(ctx, merchantProuduct)
}


func NewMerchantProductUsecase(merchantProductRepo repository.MerchantProductRepositoryInterface, productClient httpclient.ProductClientInterface, warehouseClient httpclient.WarehouseClientInterface, rabbitMQService *rabbitmq.RabbitMQService) MerchantProductUsecaseInterface {
	return &merchantProductUsecase{
		merchantProductRepo: merchantProductRepo,
		productClient:       productClient,
		warehouseClient:     warehouseClient,
		rabbitMQService:     rabbitMQService,
	}
}

//USECASE DONE lanjut controller