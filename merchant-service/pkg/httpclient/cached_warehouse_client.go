package httpclient

import (
	"context"
	"fmt"
	"time"
	"warehouse-go/merchant-service/pkg/redis"

	"github.com/gofiber/fiber/v2/log"
	goredis "github.com/redis/go-redis/v9"
)

type CachedWarehouseClient struct {
	client WarehouseClientInterface
	redis  *redis.RedisClient
	ttl    time.Duration
}

func NewCachedWarehouseClient(WarehouseClient WarehouseClientInterface, redisClient *redis.RedisClient) *CachedWarehouseClient {
	return &CachedWarehouseClient{
		client: WarehouseClient,
		redis:  redisClient,
		ttl:    1 * time.Hour,
	}
}

func (cwc *CachedWarehouseClient) generateCacheKey(prefix string, id uint) string {
	return fmt.Sprintf("warehouse:%s:%d", prefix, id)
}

func (cwc *CachedWarehouseClient) generateProductStockCacheKey(warehouseID uint, productID uint) string {
	return fmt.Sprintf("warehouse:product_stock:%d:%d", warehouseID, productID)
}

func (cwc *CachedWarehouseClient) GetWarehouseByID(ctx context.Context, warehouseID uint) (*WarehouseResponse, error) {
	cacheKey := cwc.generateCacheKey("single", warehouseID)

	var cachedWarehouse WarehouseResponse
	err := cwc.redis.Get(ctx, cacheKey, &cachedWarehouse)
	if err == nil {
		log.Infof("[CachedWarehouse] GetWarehouseByID - Cache Hit: %v", cachedWarehouse)
		return &cachedWarehouse, nil
	}

	// Log cache miss (bukan error kritis)
	if err != goredis.Nil {
		log.Warnf("[CachedWarehouse] GetWarehouseByID - Cache error (continuing): %v", err)
	}

	// Fetch dari API
	warehouse, err := cwc.client.GetWarehouseByID(ctx, warehouseID)
	if err != nil {
		log.Errorf("[CachedWarehouse] GetWarehouseByID - API error: %v", err)
		return nil, err
	}

	// Simpan ke cache
	if err := cwc.redis.Set(ctx, cacheKey, warehouse, cwc.ttl); err != nil {
		log.Warnf("[CachedWarehouse] GetWarehouseByID - Failed to cache (continuing): %v", err)
	}

	return warehouse, nil
}

func (cwc *CachedWarehouseClient) GetWarehouseProductStock(ctx context.Context, warehouseID uint, productID uint) (*WarehouseProductStockResponse, error) {
	// PERBAIKAN 1: Gunakan cache key yang spesifik untuk product stock
	cacheKey := cwc.generateProductStockCacheKey(warehouseID, productID)

	var cachedWarehouseProductStock WarehouseProductStockResponse
	err := cwc.redis.Get(ctx, cacheKey, &cachedWarehouseProductStock)
	if err == nil {
		log.Infof("[CachedWarehouse] GetWarehouseProductStock - Cache Hit: %v", cachedWarehouseProductStock)
		return &cachedWarehouseProductStock, nil
	}

	// PERBAIKAN 2: Log cache miss, tapi lanjutkan ke API call
	if err != goredis.Nil {
		log.Warnf("[CachedWarehouse] GetWarehouseProductStock - Cache error (continuing): %v", err)
	} else {
		log.Infof("[CachedWarehouse] GetWarehouseProductStock - Cache miss, fetching from API")
	}

	// Fetch dari API
	warehouseProductStock, err := cwc.client.GetWarehouseProductStock(ctx, warehouseID, productID)
	if err != nil {
		log.Errorf("[CachedWarehouse] GetWarehouseProductStock - API error: %v", err)
		return nil, err
	}

	// Simpan ke cache
	if err := cwc.redis.Set(ctx, cacheKey, warehouseProductStock, cwc.ttl); err != nil {
		log.Warnf("[CachedWarehouse] GetWarehouseProductStock - Failed to cache (continuing): %v", err)
	}

	// PERBAIKAN 3: Return data dari API call, bukan dari cached variable
	return warehouseProductStock, nil
}