package search

import (
	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/search_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/search/mongo"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/search/openSearch"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)


const (
	Mongo = "Mongo"
	ElasticSearch = "ElasticSearch"
	OpenSearch = "OpenSearch"
)

// NewSearchRepository creates a new search repository based on the specified type
func NewSearchRepository(repositoryType string, p *base.Persistence, c *gin.Context) search_repository.SearchRepository {
	switch repositoryType {
	case Mongo:
		return mongo.NewMongoRepository(p, c)
	case OpenSearch:
		return openSearch.NewOpenSearchRepository(p, c)
	default:
		return mongo.NewMongoRepository(p, c)
	}
}
