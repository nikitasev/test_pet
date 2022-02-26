package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"test_pet/internal/config"
	"test_pet/internal/domain/entity"
	"test_pet/internal/domain/service"
	"time"
)

type RedisUserListCache struct {
	client *redis.Client
	config config.Cache
}

func NewRedisUserListCache(client *redis.Client, cfg config.Cache) *RedisUserListCache {
	return &RedisUserListCache{client: client, config: cfg}
}

func (r *RedisUserListCache) SaveList(list []entity.User, limit, offset int32) error {
	key := generateCacheKey(limit, offset)
	value, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return r.client.Set(key, value, time.Second*r.config.CacheExpireInSeconds).Err()
}

func (r *RedisUserListCache) GetListByParams(limit, offset int32) ([]entity.User, error) {
	key := generateCacheKey(limit, offset)
	value, err := r.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, service.MissCacheError
		}

		return nil, err
	}

	var list []entity.User
	err = json.Unmarshal([]byte(value), &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func generateCacheKey(limit, offset int32) string {
	return fmt.Sprintf("userList_%d_%d", limit, offset)
}
