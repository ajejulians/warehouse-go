package controller

import (
	"warehouse-go/warehouse-service/controller/request"
	"warehouse-go/warehouse-service/controller/response"
	"warehouse-go/warehouse-service/model"
	"warehouse-go/warehouse-service/pkg/conv"
	"warehouse-go/warehouse-service/pkg/httpclient"
	"warehouse-go/warehouse-service/pkg/validator"
	"warehouse-go/warehouse-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type WarehouseProductControllerInterface interface {
	GetDetailWarehouse(c *fiber.Ctx) error
	GetDetailWarehouseProductByID(c *fiber.Ctx) error
	CreateWarehouseProduct(c *fiber.Ctx) error
	GetWarehouseProductByWarehouseIDAndProductID(c *fiber.Ctx) error
	UpdateWarehouseProduct(c *fiber.Ctx) error
	DeleteWarehouseProduct(c *fiber.Ctx) error
	DeleteAllWarehouseProductByProductID(c *fiber.Ctx) error
	GetWarehouseProductByProductID(c *fiber.Ctx) error
	GetProductTotalStock(c *fiber.Ctx) error
}

type warehouseProductController struct {
	warehouseProductUsecase usecase.WarehouseProductUsecaseInterface
}

// CreateWarehouseProduct implements WarehouseProductControllerInterface.
func (w *warehouseProductController) CreateWarehouseProduct(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.CreateWarehouseProductRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseProductController] CreateWarehouseProduct - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[WarehouseProductController] CreateWarehouseProduct - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : err.Error(),
		})
	}

	warehouseID := c.Params("warehouse_id")
	warehouseIDUint := conv.StringToUint(warehouseID)

	reqModel := model.WarehouseProduct{
		WarehouseID:  warehouseIDUint,
		ProductID: req.ProductID,
		Stock: req.Stock,
	}

	if err := w.warehouseProductUsecase.CreateWarehouseProduct(ctx, &reqModel); err != nil {
		log.Errorf("[WarehouseProductController] CreateWarehouseProduct - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "Failed to create warehouse product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "warehouse product created successfully",
	})
}

// DeleteAllWarehouseProductByProductID implements WarehouseProductControllerInterface.
func (w *warehouseProductController) DeleteAllWarehouseProductByProductID(c *fiber.Ctx) error {
	ctx := c.Context()

	productID := c.Params("product_id")
	productIDUint := conv.StringToUint(productID)

	if err := w.warehouseProductUsecase.DeleteAllWarehouseProductByProductID(ctx, productIDUint); err != nil {
		log.Errorf("[WarehouseProductController] DeleteAllWarehouseProductByProductID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "Failed to delete all warehouse product by product ID",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "warehouse product deleted successfully",
	})
}

// DeleteWarehouseProduct implements WarehouseProductControllerInterface.
func (w *warehouseProductController) DeleteWarehouseProduct(c *fiber.Ctx) error {
	ctx := c.Context()

	warehouseProductID := c.Params("warehouse_product_id")
	warehouseProductIDUint := conv.StringToUint(warehouseProductID)

	if err := w.warehouseProductUsecase.DeleteWarehouseProduct(ctx, warehouseProductIDUint); err != nil {
		log.Errorf("[WarehouseProductController] DeleteAllWarehouseProductByProductID - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "Failed to delete warehouse product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "warehouse product deleted successfully",
	})
}

// GetDetailWarehouse implements WarehouseProductControllerInterface.
func (w *warehouseProductController) GetDetailWarehouse(c *fiber.Ctx) error {
	ctx := c.Context()
	warehouseID := c.Params("warehouse_id")
	warehouseIDUint := conv.StringToUint(warehouseID)

	warehouse, products, err := w.warehouseProductUsecase.GetDetailWarehouse(ctx, warehouseIDUint)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetDetailWarehouse - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "failed to get detail warehouse",
		})
	}

	respWarehouseProducts := response.DetailWarehouseResponse{
		ID:			warehouse.ID,
		Name: 		warehouse.Name,
		Address: 	warehouse.Address,
		Photo: 		warehouse.Photo,
		Phone: 		warehouse.Phone,
	}

	productMap := make(map[uint]*httpclient.ProductResponse)
	for i := range products{
		productMap[products[i].ID] = &products[i]
	}

	for _, wp := range warehouse.WarehouseProducts {
		warehouseProduct := response.WarehouseProductResponse{
			ID: wp.ID,
			WarehouseID: wp.WarehouseID,
			ProductID: wp.ProductID,
			Stock: int(wp.Stock),
		}

		if product, exists := productMap[wp.ProductID]; exists {
			warehouseProduct.ProductName = product.Name
			warehouseProduct.ProductAbout = product.About
			warehouseProduct.ProductPhoto = product.Thumbnail
			warehouseProduct.ProductPrice = int(product.Price)
			warehouseProduct.ProductCategory = product.Category.Name
			warehouseProduct.ProductCategoryPhoto = product.Category.Photo
		}

		respWarehouseProducts.WarehouseProducts = append(respWarehouseProducts.WarehouseProducts, warehouseProduct)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data" : respWarehouseProducts,
		"message" : "warehouse products fetched successfully",
	})
}

// GetDetailWarehouseProductByID implements WarehouseProductControllerInterface.
func (w *warehouseProductController) GetDetailWarehouseProductByID(c *fiber.Ctx) error {
	ctx := c.Context()

	warehouseProductID := c.Params("warehouse_product_id")
	warehouseProductIDUint := conv.StringToUint(warehouseProductID)

	warehouseProduct, product, err := w.warehouseProductUsecase.GetDetailWarehouseProductByID(ctx, warehouseProductIDUint)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetDetailWarehouseProductByID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "Failed to get detail warehouse product by id",
		})
	}

	respWarehouseProduct := response.GetDetailWarehouseProductByIDResponse{
		ID: warehouseProduct.ID,
		WarehouseID: warehouseProduct.WarehouseID,
		ProductID: warehouseProduct.ProductID,
		Stock: int(warehouseProduct.Stock),
		WarehouseName: warehouseProduct.Warehouse.Name,
		WarehousePhoto: warehouseProduct.Warehouse.Photo,
		WarehousePhone: warehouseProduct.Warehouse.Phone,
		ProductName: product.Name,
		ProductBarcode: product.Barcode,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data" : respWarehouseProduct,
		"message" : "Warehouse product fetched successfully",
	})
}

// GetProductTotalStock implements WarehouseProductControllerInterface.
func (w *warehouseProductController) GetProductTotalStock(c *fiber.Ctx) error {
	ctx := c.Context()

	productID := c.Params("product_id")
	productIDUint := conv.StringToUint(productID)

	totalStock, err := w.warehouseProductUsecase.GetProductTotalStock(ctx, productIDUint)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetProductTotalStock - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "failed to get product total stock",
		})
	}

	resp := response.ProductTotalStockResponse{
		ProductID: productIDUint,
		TotalStock: totalStock,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data" : resp,
		"message" : "Product total stock fetched successfully",
	})
}

// GetWarehouseProductByProductID implements WarehouseProductControllerInterface.
func (w *warehouseProductController) GetWarehouseProductByProductID(c *fiber.Ctx) error {
	ctx := c.Context()

	productID := c.Params("product_id")
	productIDUint := conv.StringToUint(productID)

	warehouseProducts, err := w.warehouseProductUsecase.GetWarehouseProductByProductID(ctx, productIDUint)
	if err != nil {
		log.Errorf("[WarehouseProductController] GetWarehouseProductByProductID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "Failed to get warehouse product by product ID", 
		})
	}

	resps := []response.WarehouseResponse {}
	for _, wp := range warehouseProducts{
		resps = append(resps, response.WarehouseResponse{
			ID: wp.WarehouseID,
			Name: wp.Warehouse.Name,
			Address: wp.Warehouse.Address,
			Photo: wp.Warehouse.Photo,
			Phone: wp.Warehouse.Phone,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data" : resps,	
		"message" : "Warehouse products fetched successfully",
	})
}

// GetWarehouseProductByWarehouseIDAndProductID implements WarehouseProductControllerInterface.
func (w *warehouseProductController) GetWarehouseProductByWarehouseIDAndProductID(c *fiber.Ctx) error {
	ctx := c.Context()
	warehouseID := c.Params("warehouse_id")
	warehouseIDUint := conv.StringToUint(warehouseID)
	productID := c.Params("product_id")
	productIDUint := conv.StringToUint(productID)

	var warehouseProduct *model.WarehouseProduct
	warehouseProduct, err := w.warehouseProductUsecase.GetWarehouseProductByWarehouseIDAndProductID(ctx, warehouseIDUint, productIDUint)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "warehouse product not found",
			})
		}
		log.Errorf("[WarehouseProductController] GetWarehouseProductByWarehouseIDAndProductID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get warehouse product by warehouse id and product id",
		})
	}

	// ✅ Tambahkan return untuk response sukses di sini
	respWarehouseProduct := response.GetDetailWarehouseProductByIDResponse{
		ID:          warehouseProduct.ID,
		WarehouseID: warehouseProduct.WarehouseID,
		ProductID:   warehouseProduct.ProductID,
		Stock:       int(warehouseProduct.Stock),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    respWarehouseProduct,
		"message": "warehouse product fetched successfully",
	})
}

// UpdateWarehouseProduct implements WarehouseProductControllerInterface.
func (w *warehouseProductController) UpdateWarehouseProduct(c *fiber.Ctx) error {
	ctx := c.Context()

	warehouseProductID := c.Params("warehouse_product_id")
	warehouseProductIDUint := conv.StringToUint(warehouseProductID)
	warehouseID := c.Params("warehouse_id")
	warehouseIDUint := conv.StringToUint(warehouseID)

	var req request.CreateWarehouseProductRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseProductController] UpdateWarehouseProduct - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message" : "Invalid request body",
		})
	}

		if err := validator.Validate(req); err != nil {
		log.Errorf("[WarehouseProductController] UpdateWarehouseProduct - 2: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{

		})
	}

	reqModel := model.WarehouseProduct{
		ID: 			warehouseProductIDUint,
		WarehouseID: 	warehouseIDUint,
		ProductID: 		req.ProductID,
		Stock: 			req.Stock,
	}

	if err := w.warehouseProductUsecase.UpdateWarehouseProduct(ctx, &reqModel); err != nil {
		log.Errorf("[WarehouseProductController] UpdateWarehouseProduct - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message" : "Failed to update warehouse ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message" : "warehouse product updated successfully ",
	})
}

func NewWarehouseProductController(warehouseProductUsecase usecase.WarehouseProductUsecaseInterface) WarehouseProductControllerInterface {
	return &warehouseProductController{warehouseProductUsecase: warehouseProductUsecase}
}
