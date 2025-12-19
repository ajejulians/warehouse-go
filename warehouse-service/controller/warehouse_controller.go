package controller

import (
	"warehouse-go/warehouse-service/controller/request"
	"warehouse-go/warehouse-service/controller/response"
	"warehouse-go/warehouse-service/model"
	"warehouse-go/warehouse-service/pkg/conv"
	"warehouse-go/warehouse-service/pkg/pagination"
	"warehouse-go/warehouse-service/pkg/validator"
	"warehouse-go/warehouse-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type WarehouseControllerInterface interface {
	CreateWarehouse(ctx *fiber.Ctx) error
	GetAllWarehouses(ctx *fiber.Ctx) error
	GetWarehouseByID(ctx *fiber.Ctx) error
	UpdateWarehouse(ctx *fiber.Ctx) error
	DeleteWarehouse(ctx *fiber.Ctx) error
}

type warehouseController struct {
	warehouseUsecase usecase.WarehouseUsecaseInterface
}

// CreateWarehouse implements WarehouseControllerInterface.
func (w *warehouseController) CreateWarehouse(ctx *fiber.Ctx) error {
	var req request.CreateWarehouseRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid body request",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[WarehousController] CreateWarehouse - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Warehouse{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Photo:   req.Photo,
	}

	if err := w.warehouseUsecase.CreateWarehouse(ctx.Context(), &reqModel); err != nil {
		log.Errorf("[WarehouseController] CreateWarehouse - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create warehouse",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouse created successfully",
	})
}

// DeleteWarehoouse implements WarehouseControllerInterface.
func (w *warehouseController) DeleteWarehouse(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	warehouseID := conv.StringToUint(id)

	if err := w.warehouseUsecase.DeleteWarehouse(ctx.Context(), warehouseID); err != nil {
		log.Errorf("[WarehouseController] DeleteWarehouse - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete warehouse",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouse deleted successfully",
	})
}

// GetAllWarehouses implements WarehouseControllerInterface.
func (w *warehouseController) GetAllWarehouses(ctx *fiber.Ctx) error {
	var req request.GetAllWarehouseRequest
	if err := ctx.QueryParser(&req); err != nil {
		log.Errorf("[WarehouseController] GetAllWarehouses - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[WarehouseController] GetAllWarehoues - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	warehouses, total, err := w.warehouseUsecase.GetAllWarehouse(ctx.Context(), req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[WarehouseController] GetAllWarehouses - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get all warehouses",
		})
	}

	pagination := pagination.CalculatePagination(req.Page, req.Limit, int(total))
	var warehousesResponse []response.WarehouseResponse
	for _, warehouse := range warehouses {
		warehousesResponse = append(warehousesResponse, response.WarehouseResponse{
			ID:           warehouse.ID,
			Name:         warehouse.Name,
			Address:      warehouse.Address,
			Photo:        warehouse.Photo,
			Phone:        warehouse.Phone,
			CountProduct: len(warehouse.WarehouseProducts),
		})
	}

	response := response.GetAllWarehouseResponse{
		Warehouse:  warehousesResponse,
		Pagination: pagination,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouses fetched successfully",
		"data":    response,
	})
}

// GetWarehouseByID implements WarehouseControllerInterface.
func (w *warehouseController) GetWarehouseByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	warehouseID := conv.StringToUint(id)

	warehouse, err := w.warehouseUsecase.GetWarehouseByID(ctx.Context(), warehouseID)
	if err != nil {
		log.Errorf("[WarehouseController] GetWarehouseByID - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get warehouse",
		})
	}

	respWarehouses := response.DetailWarehouseResponse{
		ID:      warehouse.ID,
		Name:    warehouse.Name,
		Address: warehouse.Address,
		Photo:   warehouse.Photo,
		Phone:   warehouse.Phone,
	}

	for _, warehouseProduct := range warehouse.WarehouseProducts {
		respWarehouses.WarehouseProducts = append(respWarehouses.WarehouseProducts, response.WarehouseProductResponse{
			ID:          warehouseProduct.ID,
			WarehouseID: warehouseProduct.WarehouseID,
			ProductID:   warehouseProduct.ProductID,
			Stock:       int(warehouseProduct.Stock),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    respWarehouses,
		"message": "warehouses fetched successfully",
	})
}

// UpdateWarehouse implements WarehouseControllerInterface.
func (w *warehouseController) UpdateWarehouse(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	warehouseID := conv.StringToUint(id)

	var req request.CreateWarehouseRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[WarehouseController] UpdateWarehouse - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[WarehouseController] UpdateWarehouses - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Warehouse{
		ID:      warehouseID,
		Name:    req.Name,
		Address: req.Address,
		Photo:   req.Photo,
		Phone:   req.Phone,
	}

	if err := w.warehouseUsecase.UpdateWarehouse(ctx.Context(), &reqModel); err != nil {
		log.Errorf("[WarehouseController] UpdateWarehouse - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update warehouse",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouse updated successfully",
	})
}

func NewWarehouseController(warehouseUsecase usecase.WarehouseUsecaseInterface) WarehouseControllerInterface {
	return &warehouseController{warehouseUsecase: warehouseUsecase}
}
