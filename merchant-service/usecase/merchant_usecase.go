package usecase

import (
	"context"
	"warehouse-go/merchant-service/model"
	"warehouse-go/merchant-service/pkg/httpclient"
	"warehouse-go/merchant-service/repository"

	"github.com/gofiber/fiber/v2/log"
)

type MerchantUsecaseInterface interface {
	CreateMerchant(ctx context.Context, merchant *model.Merchant) error
	GetAllMerchants(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.Merchant, int64, error)
	GetMerchantByID(ctx context.Context, id uint) (*model.Merchant, error)
	UpdateMerchant(ctx context.Context, merchant *model.Merchant) error
	DeleteMerchant(ctx context.Context, id uint) error
	GetMerchantByKeeperID(ctx context.Context, keeperID uint) (*model.Merchant, []httpclient.ProductResponse, []httpclient.WarehouseResponse, error)
	GetKeeperName(ctx context.Context, keeperID uint) (string, error)
}

type merchantUsecase struct {
	merchantRepo    repository.MerchantRepositoryInterface
	userClient      httpclient.UserClientInterface
	warehouseClient httpclient.WarehouseClientInterface
	productClient   httpclient.ProductClientInterface
}

// CreateMerchant implements MerchantUsecaseInterface.
func (m *merchantUsecase) CreateMerchant(ctx context.Context, merchant *model.Merchant) error {
	return m.merchantRepo.CreateMerchant(ctx, merchant)
}

// DeleteMerchant implements MerchantUsecaseInterface.
func (m *merchantUsecase) DeleteMerchant(ctx context.Context, id uint) error {
	return m.merchantRepo.DeleteMerchant(ctx, id)
}

// GetAllMerchants implements MerchantUsecaseInterface.
func (m *merchantUsecase) GetAllMerchants(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.Merchant, int64, error) {
	return m.merchantRepo.GetAllMerchants(ctx, page, limit, search, sortBy, sortOrder)
}

func (m *merchantUsecase) GetKeeperName(ctx context.Context, keeperID uint) (string, error) {
	log.Infof("[MerchantUsecase] GetKeeperName called with keeperID: %d", keeperID)
	
	if keeperID != 0 {
		user, err := m.userClient.GetUserByID(ctx, keeperID)
		if err != nil {
			log.Errorf("[MerchantUsecase] GetKeeperName - 1: %v", err)
			return "", err
		}
		
		log.Infof("[MerchantUsecase] GetKeeperName - Got user name: %s", user.Name)
		return user.Name, nil
	}
	
	log.Warnf("[MerchantUsecase] GetKeeperName - keeperID is 0")
	return "", nil
}
// // GetKeeperName implements MerchantUsecaseInterface.
// func (m *merchantUsecase) GetKeeperName(ctx context.Context, keeperID uint) (string, error) {
// 	if keeperID != 0 {
// 		user, err := m.userClient.GetUserByID(ctx, keeperID)
// 		if err != nil {
// 			log.Errorf("[MerchantUsecase] GetKeeperName - 1: %v", err)
// 			return "", err
// 		}
// 		return user.Name, nil
// 	}
// 	return "", nil

// }

// GetMerchantByID implements MerchantUsecaseInterface.
func (m *merchantUsecase) GetMerchantByID(ctx context.Context, id uint) (*model.Merchant, error) {
	return m.merchantRepo.GetMerchantByID(ctx, id)
}

// GetMerchantByKeeperID implements MerchantUsecaseInterface.
func (m *merchantUsecase) GetMerchantByKeeperID(ctx context.Context, keeperID uint) (*model.Merchant, []httpclient.ProductResponse, []httpclient.WarehouseResponse, error) {
	merchant, err := m.merchantRepo.GetMerchantByKeeperID(ctx, keeperID)
	if err != nil {
		log.Errorf("[MerchantUsecase] GetMerchantByKeeperID - 1: %v", err)
		return nil, nil, nil, err
	}

	var products []httpclient.ProductResponse
	var warehouses []httpclient.WarehouseResponse

	if len(merchant.MerchantProducts) > 0 {
		for _, mp := range merchant.MerchantProducts {
			product, err := m.productClient.GetProductByID(ctx, mp.ProductID)
			if err != nil {
				log.Errorf("[MerchantUsecase] GetMerchantByID - 2: %v", err)
				return nil, nil, nil, err
			}
			products = append(products, *product)

			warehouse, err := m.warehouseClient.GetWarehouseByID(ctx, mp.WarehouseID)
			if err != nil {
				log.Errorf("[MerchantUsecase] GetMerchantByID - 3: %v", err)
				return nil, nil, nil, err
			}
			warehouses = append(warehouses, *warehouse)
		}
	}
	return merchant, products, warehouses, nil
}

// UpdateMerchant implements MerchantUsecaseInterface.
func (m *merchantUsecase) UpdateMerchant(ctx context.Context, merchant *model.Merchant) error {
	return m.merchantRepo.UpdateMerchant(ctx, merchant)
}

func NewMerchantUsecase(merchantRepo repository.MerchantRepositoryInterface, userClient httpclient.UserClientInterface, warehouseClient httpclient.WarehouseClientInterface, productClient httpclient.ProductClientInterface) MerchantUsecaseInterface {
	return &merchantUsecase{
		merchantRepo:    merchantRepo,
		userClient:      userClient,
		warehouseClient: warehouseClient,
		productClient:   productClient,
	}
}

//Lanjut Ke MerchantProduct Usecase
