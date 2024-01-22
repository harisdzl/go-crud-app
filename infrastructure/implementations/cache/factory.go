package cache

import (
	"github.com/harisquqo/quqo-challenge-1/domain/repository/cache_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/cache/redis"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)


const (
	Redis = "Redis"
)

// NewCacheRepository creates a new cache repository based on the specified type
func NewCacheRepository(repositoryType string, p *base.Persistence) cache_repository.CacheRepository {
	switch repositoryType {
	case Redis:
		return redis.NewRedisRepository(p)
	// Add cases for other cache repository types if needed
	default:
		return redis.NewRedisRepository(p)
	}
}
