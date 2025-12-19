package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	"warehouse-go/merchant-service/configs"

	"github.com/gofiber/fiber/v2/log"
)

type WarehouseClientInterface interface {
	GetWarehouseByID(ctx context.Context, warehouseID uint) (*WarehouseResponse, error)
	GetWarehouseProductStock(ctx context.Context, warehouseID uint, productID uint) (*WarehouseProductStockResponse, error)
}

type WarehouseClient struct {
	urlWarehouseService string
	httpClient          *http.Client
}

// GetWarehouseByID implements WarehouseClientInterface.
func (w *WarehouseClient) GetWarehouseByID(ctx context.Context, warehouseID uint) (*WarehouseResponse, error) {
	url := fmt.Sprintf("%s/api/v1/warehouses/%d", w.urlWarehouseService, warehouseID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[WarehouseClient] GetWarehouseByID - 1: %v", err)
		return nil, err
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		log.Errorf("[WarehouseClient] GetWarehouseByID - 2: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[WarehouseClient] GetWarehouseByID - 3: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[WarehouseClient] GetWarehouseByID - 4: Status=%d, Body=%s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("failed to get warehouse by id: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var warehouseResponse WarehouseServiceResponse
	if err := json.Unmarshal(body, &warehouseResponse); err != nil {
		log.Errorf("[WarehouseClient] GetWarehouseByID - 5: %v", err)
		return nil, err
	}

	return &warehouseResponse.Data, nil
}

// GetWarehouseProductsStock implements WarehouseClientInterface.
func (w *WarehouseClient) GetWarehouseProductStock(ctx context.Context, warehouseID uint, productID uint) (*WarehouseProductStockResponse, error) {
    url := fmt.Sprintf("%s/api/v1/warehouse-products/%d/detail/%d", w.urlWarehouseService, warehouseID, productID)
	var ErrWarehouseProductNotFound = errors.New("warehouse product not found")
    log.Infof("[WarehouseClient] GetWarehouseProductStock - URL: %s", url)

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        log.Errorf("[WarehouseClient] GetWarehouseProductStock - 1: %v", err)
        return nil, err
    }

    resp, err := w.httpClient.Do(req)
    if err != nil {
        log.Errorf("[WarehouseClient] GetWarehouseProductStock - 2: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Errorf("[WarehouseClient] GetWarehouseProductStock - 3: %v", err)
        return nil, err
    }

    // ✅ Bedakan 404 vs error lain
	
    if resp.StatusCode == http.StatusNotFound {
        log.Warnf("[WarehouseClient] GetWarehouseProductStock - NotFound: Status=%d, Body=%s", resp.StatusCode, string(body))
        return nil, ErrWarehouseProductNotFound
    }

    if resp.StatusCode != http.StatusOK {
        log.Errorf("[WarehouseClient] GetWarehouseProductStock - 4: Status=%d, Body=%s", resp.StatusCode, string(body))
        return nil, fmt.Errorf("failed to get warehouse product stock: status=%d, body=%s", resp.StatusCode, string(body))
    }

    var warehouseProductStockResponse WarehouseProductStockServiceResponse
    if err := json.Unmarshal(body, &warehouseProductStockResponse); err != nil {
        log.Errorf("[WarehouseClient] GetWarehouseProductStock - 5: %v, Body=%s", err, string(body))
        return nil, err
    }

    log.Infof("[WarehouseClient] GetWarehouseProductStock - Success: %+v", warehouseProductStockResponse.Data)

    return &warehouseProductStockResponse.Data, nil
}
type WarehouseResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Photo   string `json:"photo"`
	Phone   string `json:"phone"`
}

type WarehouseServiceResponse struct {
	Message string            `json:"message"`
	Data    WarehouseResponse `json:"data"`
	Error   string            `json:"error,omitempty"`
}

type WarehouseProductStockResponse struct {
	ID          uint `json:"id"`
	ProductID   uint `json:"product_id"`
	Stock       int  `json:"stock"`
	WarehouseID uint `json:"warehouse_id"`
}

type WarehouseProductStockServiceResponse struct {
	Message string                        `json:"message"`
	Data    WarehouseProductStockResponse `json:"data"`
	Error   string                        `json:"error,omitempty"`
}

func NewWarehouseClient(cfg configs.Config) WarehouseClientInterface {
	return &WarehouseClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		urlWarehouseService: cfg.App.UrlWarehouseService,
	}
}