package search

import (
	"github.com/harisquqo/quqo-challenge-1/domain/repository/search_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/search/elasticSearch"
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
func NewSearchRepository(repositoryType string, p *base.Persistence) search_repository.SearchRepository {
	switch repositoryType {
	case Mongo:
		return mongo.NewMongoRepository(p)
	// Add cases for other search repository types if needed
	case ElasticSearch:
		return elasticSearch.NewElasticSearchRepository(p)
	case OpenSearch:
		return openSearch.NewOpenSearchRepository(p)
	default:
		return mongo.NewMongoRepository(p)
	}
}
