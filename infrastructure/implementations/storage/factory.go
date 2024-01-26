package storage

import (
	"github.com/harisquqo/quqo-challenge-1/domain/repository/storage_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/storage/supabase"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)


const (
	Supabase = "Supabase"
)

// NewSearchRepository creates a new search repository based on the specified type
func NewStorageRepository(repositoryType string, p *base.Persistence) storage_repository.StorageRepository {
	switch repositoryType {
		case Supabase:
			return supabase.NewStorageRepository(p)
		// Add cases for other search repository types if needed
		default:
			return supabase.NewStorageRepository(p)
		}
}
