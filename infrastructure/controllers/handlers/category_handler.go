package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/application"
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/category_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/category_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)


type Category struct {
	CategoryRepo category_repository.CategoryRepository
	Persistence *base.Persistence
}



func NewCategory(p *base.Persistence) *Category {
	return &Category{
		Persistence: p,
	}
}

//	@Summary		Save Category
//	@Description	Saves a new category.
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			category	body		category_entity.Category	true	"Category object to be saved"
//	@Success		200			{object}	entity.ResponseContext		"Success"
//	@Failure		400			{object}	entity.ResponseContext		"Bad request"
//	@Failure		500			{object}	entity.ResponseContext		"Internal server error"
//	@Router			/category [post]
func (ca *Category) SaveCategory(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	category := category_entity.Category{}


    if err := c.ShouldBindJSON(&category); err != nil {
        c.JSON(http.StatusBadRequest,
            responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
        return
    }

	
	ca.CategoryRepo = application.NewCategoryApplication(ca.Persistence, c)
	savedCategory, savedCategoryErr := ca.CategoryRepo.SaveCategory(&category)

	if savedCategoryErr != nil {
		log.Println(savedCategoryErr)
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, savedCategoryErr["db_error"], ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Category saved successfully", savedCategory))
}

//	@Summary		Get Category
//	@Description	Retrieves a category by its ID.
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path		int						true	"Category ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/category/{category_id} [get]
func (ca *Category) GetCategory(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	categoryID, err := strconv.ParseInt((c.Param("category_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	

	ca.CategoryRepo = application.NewCategoryApplication(ca.Persistence, c)

	category, err := ca.CategoryRepo.GetCategory(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("Category %v obtained", categoryID), category))
}

//	@Summary		Get All Categories
//	@Description	Retrieves all categories.
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.ResponseContext	"Success"
//	@Failure		500	{object}	entity.ResponseContext	"Internal server error"
//	@Router			/categories [get]
func (ca *Category) GetAllCategories(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	
	
	ca.CategoryRepo = application.NewCategoryApplication(ca.Persistence, c)

	allCategories, err := ca.CategoryRepo.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	results := map[string]interface{}{
		"results" : allCategories,
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "All categories obtained", results))
}

//	@Summary		Get Parent Categories
//	@Description	Retrieves parent categories of a given category.
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path		int						true	"Category ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/category/{category_id}/parents [get]
func (ca *Category) GetParentCategories(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	

	ca.CategoryRepo = application.NewCategoryApplication(ca.Persistence, c)
	childCategoryId, err := strconv.ParseInt((c.Param("category_id")), 10, 64)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	parentCategories, parentCategoryErr := ca.CategoryRepo.GetParentCategories(childCategoryId)

	if parentCategoryErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}


	results := map[string]interface{}{
		"results" : parentCategories,
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "All parent categories obtained", results))
}

//	@Summary		Delete Category
//	@Description	Deletes a category by its ID.
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path		int						true	"Category ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/category/{category_id} [delete]
func (ca *Category) DeleteCategory(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	

	categoryID, err := strconv.ParseInt((c.Param("category_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	ca.CategoryRepo = application.NewCategoryApplication(ca.Persistence, c)

	deleteErr := ca.CategoryRepo.DeleteCategory(categoryID)
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusError, deleteErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess,fmt.Sprintf("Category %v has been deleted", categoryID), ""))
}

//	@Summary		Update Category
//	@Description	Updates a category.
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path		int							true	"Category ID"
//	@Param			category	body		category_entity.Category	true	"Category object to be updated"
//	@Success		200			{object}	entity.ResponseContext		"Success"
//	@Failure		400			{object}	entity.ResponseContext		"Bad request"
//	@Failure		404			{object}	entity.ResponseContext		"Category not found"
//	@Failure		422			{object}	entity.ResponseContext		"Unprocessable entity"
//	@Failure		500			{object}	entity.ResponseContext		"Internal server error"
//	@Router			/category/{category_id} [put]
func (ca *Category) UpdateCategory(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	categoryID, err := strconv.ParseInt(c.Param("category_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, "Invalid category ID", ""))
		return
	}

	


	// Check if the category exists
	ca.CategoryRepo = application.NewCategoryApplication(ca.Persistence, c)

	existingCategory, err := ca.CategoryRepo.GetCategory(categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, responseContextData.ResponseData(entity.StatusFail, "Category not found", ""))
		return
	}

	// Bind the JSON request body to the existing category
	if err := c.ShouldBindJSON(&existingCategory); err != nil {
		c.JSON(http.StatusUnprocessableEntity, responseContextData.ResponseData(entity.StatusFail, "Invalid JSON", ""))
		return
	}

	ca.CategoryRepo = application.NewCategoryApplication(ca.Persistence, c)

	// Update the category
	updatedCategory, updateErr := ca.CategoryRepo.UpdateCategory(existingCategory)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, updateErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Product updated succesfully", updatedCategory))
}