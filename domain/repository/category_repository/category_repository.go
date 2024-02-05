package category_repository

import "github.com/harisquqo/quqo-challenge-1/domain/entity/category_entity"


type CategoryRepository interface {
	SaveCategory(*category_entity.Category) (*category_entity.Category, map[string]string)
	GetCategory(int64) (*category_entity.Category, error)
	GetAllCategories() ([]category_entity.Category, error)
	GetParentCategories(int64) ([]category_entity.Category, error)
	UpdateCategory(*category_entity.Category) (*category_entity.Category, error)
	DeleteCategory(int64) error
}
