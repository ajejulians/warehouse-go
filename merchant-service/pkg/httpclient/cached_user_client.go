package httpclient

import (
	"context"
	"fmt"
	"time"
	"warehouse-go/merchant-service/pkg/redis"

	"github.com/gofiber/fiber/v2/log"
)

type CachedUserClient struct {
	client UserClientInterface
	redis  *redis.RedisClient
	ttl 	time.Duration
}

func NewCachedUserClient(userClient UserClientInterface, redisClient *redis.RedisClient) *CachedUserClient {
	return &CachedUserClient{
		client: userClient,
		redis: redisClient,
		ttl : 1 * time.Hour,
	}
}

func (cuc *CachedUserClient) generateCacheKey(prefix string, id uint) string {
	return fmt.Sprintf("user:%s:%d", prefix, id)
}

func (cuc *CachedUserClient) GetUserByID(ctx context.Context, userID uint) (*UserResponse, error) {
    cacheKey := cuc.generateCacheKey("single", userID)

    var cachedUser UserResponse
    err := cuc.redis.Get(ctx, cacheKey, &cachedUser)
    if err == nil {
        // Cache HIT
        log.Infof("[CachedUserClient] Cache HIT - user: %+v", cachedUser)
        return &cachedUser, nil
    }

    // Cache MISS
    log.Warnf("[CachedUserClient] Cache MISS - calling user-service...")

    // Call real user service
    user, err := cuc.client.GetUserByID(ctx, userID)
    if err != nil {
        log.Errorf("[CachedUserClient] Error calling user service: %v", err)
        return nil, err
    }

    // Save to cache
    if err := cuc.redis.Set(ctx, cacheKey, user, cuc.ttl); err != nil {
        log.Errorf("[CachedUserClient] Failed to write cache: %v", err)
    }

    return user, nil
}

