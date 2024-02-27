package application

import (
	"context"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/category_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/category_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/categories"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type CategoryApp struct {
	p *base.Persistence
	c context.Context
}

func NewCategoryApplication(p *base.Persistence, c context.Context) category_repository.CategoryRepository {
	return &CategoryApp{p, c}
}


func (c *CategoryApp) SaveCategory(category *category_entity.Category) (*category_entity.Category, map[string]string) {
	repoCategory := categories.NewCategoryRepository(c.p, c.c)
	return repoCategory.SaveCategory(category)
}


func (c *CategoryApp) GetCategory(categoryId int64) (*category_entity.Category, error) {
	repoCategory := categories.NewCategoryRepository(c.p, c.c)

	return repoCategory.GetCategory(categoryId)
}


func (c *CategoryApp) GetAllCategories() ([]category_entity.Category, error) {
	repoCategory := categories.NewCategoryRepository(c.p, c.c)
	return repoCategory.GetAllCategories()
}

func (c *CategoryApp) GetParentCategories(categoryId int64) ([]category_entity.Category, error) {
	repoCategory := categories.NewCategoryRepository(c.p, c.c)
	return repoCategory.GetParentCategories(categoryId)
}

func (c *CategoryApp) UpdateCategory(image *category_entity.Category) (*category_entity.Category, error) {
	repoCategory := categories.NewCategoryRepository(c.p, c.c)

	return repoCategory.UpdateCategory(image)
}


func (c *CategoryApp) DeleteCategory(categoryId int64) (error) {
	repoCategory := categories.NewCategoryRepository(c.p, c.c)

	return repoCategory.DeleteCategory(categoryId)
}
