package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/harisquqo/quqo-challenge-1/domain/repository/cache_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"go.uber.org/zap"
)

type redisRepo struct {
	p *base.Persistence
}

func (r redisRepo) SetKey(key string, value interface{}, expiration time.Duration) error {
	// Check if there is a Redis connection 
	if r.p.DbRedis == nil {
		return errors.New("REDIS NOT FOUND")
	}

	cacheEntry, err := json.Marshal(value)
	if err != nil {
		return err
	}

	ctx := context.TODO()

	err = r.p.DbRedis.Set(ctx, key, cacheEntry, expiration).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r redisRepo) GetKey(key string, src interface{}) error {
	// Check if there is a Redis connection 
	if r.p.DbRedis == nil {
		return errors.New("REDIS NOT FOUND")
	}

	ctx := context.TODO()


	result, err := r.p.DbRedis.Get(ctx, key).Result()
	if err != nil {
		zap.S().Error("1. Redis GetKey ERROR", "error", err, "key", key)
		return err
	}

	err = json.Unmarshal([]byte(result), &src)
	if err != nil {
		zap.S().Error("2. Redis GetKey ERROR", "error", err, "key", key)
		return err
	}
	
	return nil
}

func (r redisRepo) DelKey(key string) error {
	// Check if there is a Redis connection 
	if r.p.DbRedis == nil {
		return errors.New("REDIS NOT FOUND")
	}

	ctx := context.TODO()


	_, err := r.p.DbRedis.Del(ctx, key).Result()
	if err != nil {
		zap.S().Error("1. Redis GetKey ERROR", "error", err, "key", key)
		return err
	}

	
	return nil
}



func NewRedisRepository(p *base.Persistence) cache_repository.CacheRepository {
	return &redisRepo{p}
}