package httpclient

import (
	"context"
	"fmt"
	"time"
	"warehouse-go/merchant-service/pkg/redis"

	"github.com/gofiber/fiber/v2/log"
)

type CachedProductClient struct {
	client ProductClientInterface
	redis  *redis.RedisClient
	ttl 	time.Duration
}

func NewCachedProductClient(productClient ProductClientInterface, redisClient *redis.RedisClient) *CachedProductClient {
	return &CachedProductClient{
		client: productClient,
		redis : redisClient,
		ttl : 1 * time.Hour,
	}
}

func (cpc *CachedProductClient) generateCacheKey(prefix string, id uint) string {
	return fmt.Sprintf("product:%s:%d", prefix, id)
}

func (cpc *CachedProductClient) generateCacheKeyMultiple(prefix string, ids []uint) string {
	key := fmt.Sprintf("product:%s:%d", prefix)
	for _, id := range ids {
		key += fmt.Sprintf(":id", id)
	}

	return key[:len(key)-1]
}

func (cpc *CachedProductClient) GetProductByID(ctx context.Context, productID uint) (*ProductResponse, error) {
	cacheKey := cpc.generateCacheKey("single", productID)

	var cachedProduct ProductResponse
	if err := cpc.redis.Get(ctx, cacheKey, cachedProduct); err == nil {
		log.Infof("[CachedProductClient] GetProductByID - 1: %v", cachedProduct)
		return &cachedProduct, nil
	}

	product, err := cpc.client.GetProductByID(ctx, productID)
	if err != nil {
		log.Errorf("[CachedProductClient] GetProductByID - 1: %v", err)
		return nil, err
	}

	err = cpc.redis.Set(ctx, cacheKey, product, cpc.ttl)
	if err != nil {
		log.Errorf("[CachedProductClient] GetProductByID - 2: %v", err)
		return nil, err
	}

	return product, nil
}

func (cpc *CachedProductClient) GetProductByBarcode(ctx context.Context, barcode string) (*ProductResponse, error) {
	cacheKey := fmt.Sprintf("product:barcode:%s", barcode)

	var cachedProduct ProductResponse
	if err := cpc.redis.Get(ctx, cacheKey, &cachedProduct); err == nil{
		log.Infof("[CachedProductClient] GetProductByBarcode - 1: %v", cachedProduct) 
		return &cachedProduct, nil
	}

	product, err := cpc.client.GetProductByBarcode(ctx, barcode)
	if err != nil {
		log.Errorf("[CachedProductClient] GetProductByBarcode - 2: %v", err)
		return nil, err
	}

	err = cpc.redis.Set(ctx, cacheKey, product, cpc.ttl)
	if err != nil {
		log.Errorf("[CachedProductClient] GetProductByBarcode - 3: %v", err)
		return nil, err
	}

	return product, nil
}


func (cpc *CachedProductClient) GetProducts(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]ProductResponse, error) {
	return cpc.client.GetProducts(ctx, page, limit, search, sortBy, sortOrder)
}

func (cpc *CachedProductClient) HealthCheck(ctx context.Context) error {
	return cpc.client.HealthCheck(ctx)
}
