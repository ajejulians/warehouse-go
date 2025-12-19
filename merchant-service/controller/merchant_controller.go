package controller

import (
	"warehouse-go/merchant-service/controller/request"
	"warehouse-go/merchant-service/controller/response"
	"warehouse-go/merchant-service/model"
	"warehouse-go/merchant-service/pkg/conv"
	"warehouse-go/merchant-service/pkg/httpclient"
	"warehouse-go/merchant-service/pkg/pagination"
	"warehouse-go/merchant-service/pkg/validator"
	"warehouse-go/merchant-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type MerchantControllerInterface interface {
	CreateMerchant(c *fiber.Ctx) error
	GetAllMerchant(c *fiber.Ctx) error
	GetMerchantByID(c *fiber.Ctx) error
	UpdateMerchant(c *fiber.Ctx) error
	DeleteMerchant(c *fiber.Ctx) error
}

type merchantController struct {
	merchantUsecase usecase.MerchantUsecaseInterface
}

// CreateMerchant implements MerchantControllerInterface.
func (m *merchantController) CreateMerchant(c *fiber.Ctx) error {
	var req request.CreateMerchantRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[MerchantController] CreateMerchant - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : "Invalid request body",
		})
	}

		if err := validator.Validate(req); err != nil {
			log.Errorf("[MerchantController] CreateMerchant - 2: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message" : err.Error(),
			})
		}

		reqModel := model.Merchant{
			Name: req.Name,
			KeeperID: req.KeeperID,
			Address: req.Address,
			Phone: req.Phone,
			Photo: req.Photo,
		}

		if err := m.merchantUsecase.CreateMerchant(c.Context(), &reqModel); err != nil {
			log.Errorf("[MerchantController] CreateMerchant - 2: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message" : "Failed to create Merchant",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message" : "Merchant created successfully",
		})
} 

// DeleteMerchant implements MerchantControllerInterface.
func (m *merchantController) DeleteMerchant(c *fiber.Ctx) error {
	id := c.Params("id")

	merchantID := conv.StringToUint(id)

	if err := m.merchantUsecase.DeleteMerchant(c.Context(), merchantID); err != nil {
		log.Errorf("[MerchantController] DeleteMerchant - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "Merchant deleted successfully",
	})
}

// GetAllMerchant implements MerchantControllerInterface.
func (m *merchantController) GetAllMerchant(c *fiber.Ctx) error {
	var req request.GetMerchantProductRequest
	if err := c.QueryParser(&req); err != nil {
		log.Errorf("[MerchantController] GetAllMerchant - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : "Invalid request body",
		})
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

if req.KeeperID != 0 {
	productMap := make(map[uint]*httpclient.ProductResponse)
	warehouseMap := make(map[uint]*httpclient.WarehouseResponse)

	merchant, products, warehouses, err := m.merchantUsecase.GetMerchantByKeeperID(c.Context(), req.KeeperID)
	if err != nil {
		log.Errorf("[MerchantController] GetAllMerchant - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "failed to get merchant by keeper id",
		})
	}

	if merchant.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message" : "No merchant found",
		})
	} 

		keeperName, err := m.merchantUsecase.GetKeeperName(c.Context(), merchant.KeeperID)
		if err != nil {
			log.Errorf("[MerchantController] GetAllMerchant - 3: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message" : "Failed to get keeper name",
			})
		}

		merchantResponse := response.MerchantWithProductResponse{
			ID: merchant.ID,
			Name: merchant.Name,
			Address: merchant.Address,
			Photo: merchant.Photo,
			Phone: merchant.Phone,
			KeeperID: merchant.KeeperID,
			KeepersName: keeperName,
		}

		
		if len(merchant.MerchantProducts) > 0 {
			for i := range products {
				productMap[products[i].ID] = &products[i]
			}

			for i := range warehouses {
				warehouseMap[warehouses[i].ID] = &warehouses[i]
			}

			for _, mp := range merchant.MerchantProducts {
				productResponse := response.MerchantProduct{
					ID: mp.ID,
					MerchantID: mp.MerchantID,
					ProductID: mp.ProductID,
					Stock: mp.Stock,
					WarehouseID: mp.WarehouseID,
				}

				if product, exists := productMap[mp.ProductID]; exists {
					productResponse.ProductName = product.Name
					productResponse.ProductAbout = product.About
					productResponse.ProductPhoto = product.Thumbnail
					productResponse.ProductPrice = int(product.Price)
					productResponse.ProductCategory = product.Category.Name	
					productResponse.ProductCategoryphoto = product.Category.Photo
				}

				if warehouse, exists := warehouseMap[mp.WarehouseID]; exists {
					productResponse.WarehouseName = warehouse.Name
					productResponse.WarehousePhoto = warehouse.Photo
					productResponse.WarehousePhone = warehouse.Phone
				}

				merchantResponse.MerchantProducts = append(merchantResponse.MerchantProducts, productResponse)
			}
		}
		
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data" : merchantResponse,
			"message" : "Merchant products fetched successfully",
		})
	}

	merchants, total, err := m.merchantUsecase.GetAllMerchants(c.Context(), req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[MerchantController] GetAllMerchant - 4: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "Failed to get all merchants",
		})
	}

	paginationInfo := pagination.CalculatePagination(req.Page, req.Limit, int(total))
	if len(merchants) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message" : "No merchants found",
			"data" : []response.MerchantResponse{},
			"pagination" : paginationInfo,
		})
	}

	var merchantsResponse []response.MerchantResponse
	for _, merchant := range merchants {
		keeperName, err := m.merchantUsecase.GetKeeperName(c.Context(), merchant.KeeperID)
		if err != nil {
			log.Errorf("[MerchantController] GetAllMerchant - 5: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message" : "Failed to get keeper name",
			})
		}
		merchantsResponse = append(merchantsResponse, response.MerchantResponse{
			ID: merchant.ID,
			Name: merchant.Name,
			Address: merchant.Address,
			Photo: merchant.Photo,
			Phone: merchant.Phone,
			KeeperID: merchant.KeeperID,
			KeepersName: keeperName,
			ProductCount: len(merchant.MerchantProducts),
		})
	}

	resp := response.MerchantPaginationResponse{
		Message: "Merchants fetched successfully",
		Data: merchantsResponse,
		Pagination: paginationInfo,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data" : resp,
		"pagination" : paginationInfo,
		"message" : "Merchants fetched successfully",
	})
}

// GetMerchantByID implements MerchantControllerInterface.
func (m *merchantController) GetMerchantByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id := conv.StringToUint(idStr)

	merchant, err := m.merchantUsecase.GetMerchantByID(c.Context(), id)
	if err != nil {
		log.Errorf("[MerchantController] GetMerchantByID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "Failed to get merchant by id",
		})
	}
	var keeperName string
	if merchant.KeeperID != 0 {
		keeperName, err = m.merchantUsecase.GetKeeperName(c.Context(), merchant.KeeperID)
		if err != nil {
			log.Errorf("[MerchantController] GetMerchantByID - 2: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message" : "Failed to get keeper name", 
			})
		}
	}
		response := response.MerchantResponse{
			ID: merchant.ID,
			Name: merchant.Name,
			Address: merchant.Address,
			Photo: merchant.Photo,
			Phone: merchant.Phone,
			KeeperID: merchant.KeeperID,
			KeepersName: keeperName,
			ProductCount: len(merchant.MerchantProducts),
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data" : response,
			"message" : "Merchant fetched successfully",
		})
}


// UpdateMerchant implements MerchantControllerInterface.
func (m *merchantController) UpdateMerchant(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id := conv.StringToUint(idStr)

	var req request.CreateMerchantRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[MerchantController] UpdateMerchant - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[MerchantController] UpdateMerchant - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : err.Error(),
		})
	}

	reqModel := model.Merchant{
		ID: id,
		Name: req.Name,
		KeeperID: req.KeeperID,
		Address: req.Address,
		Phone: req.Phone,
		Photo: req.Photo,
	}

	if err := m.merchantUsecase.UpdateMerchant(c.Context(), &reqModel); err != nil {
		log.Errorf("[MerchantController] UpdateMerchant - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "Failed to update merchant", 
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "Merchant updated successfully",
	})

}

func NewMerchantController(merchantUsecase usecase.MerchantUsecaseInterface) MerchantControllerInterface {
	return &merchantController{
		merchantUsecase: merchantUsecase,
	}
}


//CONTROLLER MERCHANT DONE lanjut ke MerCHANT PRODUCT CONTROLLER