package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/application"
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/inventory_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/product_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/product_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)




type Product struct {
	productRepo product_repository.ProductHandlerRepository
	Persistence *base.Persistence
}

type SaveProductResponse struct {
    Product   product_entity.Product   `json:"product"`
    Inventory inventory_entity.Inventory `json:"inventory"`
}




func NewProduct(p *base.Persistence) *Product {
	return &Product{
		Persistence: p,
	}
}

// SaveProduct saves a single product to the database.
//	@Summary		Save a single product
//	@Description	SaveProduct saves a single product to the database.
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			product	body		product_entity.Product		true	"Product object to be saved"
//	@Success		200		{object}	product_entity.Product		"Successfully saved product"
//	@Failure		400		{object}	map[string]string	"Invalid JSON"
//	@Failure		422		{object}	map[string]string	"Unprocessable entity"
//	@Router			/products [post]
func (pr *Product) SaveProduct(c *gin.Context) {
    responseContextData := entity.ResponseContext{Ctx: c}
    productForInventory := product_entity.ProductForInventory{}
	

    if err := c.ShouldBindJSON(&productForInventory); err != nil {
        c.JSON(http.StatusBadRequest,
            responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
        return
    }

    pr.productRepo = application.NewProductApplication(pr.Persistence, c)

    // Call the application layer method to save the product
	savedProduct, savedInventory, saveErr := pr.productRepo.SaveProductAndInventory(productForInventory)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, "Fail to save product", ""))
		return
	}

	response := SaveProductResponse{
        Product:   *savedProduct,
        Inventory: *savedInventory,
    }

	
    // Send the saved product as JSON response
    c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "", response))
}
// func (pr *Product) SaveMultipleProducts(c *gin.Context) {
// 	var product = []product_entity.Product{}

// 	if err:= c.ShouldBindJSON(&product); err != nil {
// 		c.JSON(http.StatusUnprocessableEntity, gin.H{
// 			"invalid_json": "invalid json",
// 		})
// 		return	
// 	}

// 	fmt.Println(product)
// 	pr.productRepo = application.NewProductApplication(pr.Persistence, c)

// 	savedProduct, saveErr := pr.productRepo.SaveMultipleProducts(&product)
// 	if saveErr != nil {
// 		c.JSON(http.StatusInternalServerError, saveErr)
// 		return
// 	}
// 	c.JSON(http.StatusCreated, savedProduct)

// }

// GetAllProducts retrieves all products.
//	@Summary		Get All Products
//	@Description	Retrieves all products.
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.ResponseContext	"Success"
//	@Failure		500	{object}	entity.ResponseContext	"Internal server error"
//	@Router			/products [get]
func (pr *Product) GetAllProducts(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	

	pr.productRepo = application.NewProductApplication(pr.Persistence, c)

	allProduct, err := pr.productRepo.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusError, err.Error(), ""))
		return
	}
	results := map[string]interface{}{
		"results" : allProduct,
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "All products obtained", results))
}

// GetProduct retrieves a specific product by ID.
//	@Summary		Get Product
//	@Description	Retrieves a specific product by ID.
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int						true	"Product ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/products/{product_id} [get]
func (pr *Product) GetProduct(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	

	productId, err := strconv.ParseInt((c.Param("product_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	pr.productRepo = application.NewProductApplication(pr.Persistence, c)

	product, getErr := pr.productRepo.GetProduct(productId)
	if getErr != nil {
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, getErr.Error(), ""))
		return
	}
	
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("Product %v obtained", productId), product))
}

// DeleteProduct deletes a product by ID.
//	@Summary		Delete Product
//	@Description	Deletes a product by ID.
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int						true	"Product ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/products/{product_id} [delete]
func (pr *Product) DeleteProduct(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	productId, err := strconv.ParseInt((c.Param("product_id")), 10, 64)
	

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	pr.productRepo = application.NewProductApplication(pr.Persistence, c)

	deleteErr := pr.productRepo.DeleteProduct(productId)
	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusError, deleteErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess,fmt.Sprintf("Product %v has been deleted", productId), ""))
}

// UpdateProduct updates a product.
//	@Summary		Update Product
//	@Description	Updates a product.
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int						true	"Product ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		404			{object}	entity.ResponseContext	"Not found"
//	@Failure		422			{object}	entity.ResponseContext	"Unprocessable entity"
//	@Router			/products/{product_id} [put]
func (pr *Product) UpdateProduct(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	productId, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, "Invalid product ID", ""))
		return
	}

	// Check if the product exists
	pr.productRepo = application.NewProductApplication(pr.Persistence, c)

	existingProduct, err := pr.productRepo.GetProduct(productId)
	if err != nil {
		c.JSON(http.StatusNotFound, responseContextData.ResponseData(entity.StatusFail, "Product not found", ""))
		return
	}

	// Bind the JSON request body to the existing product
	if err := c.ShouldBindJSON(&existingProduct); err != nil {
		c.JSON(http.StatusUnprocessableEntity, responseContextData.ResponseData(entity.StatusFail, "Invalid JSON", ""))
		return
	}

	pr.productRepo = application.NewProductApplication(pr.Persistence, c)

	// Update the product
	updatedProduct, updateErr := pr.productRepo.UpdateProduct(existingProduct)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, updateErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Product updated succesfully", updatedProduct))
}

// SearchProduct searches for products by name.
//	@Summary		Search Product
//	@Description	Searches for products by name.
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string					false	"Product name"
//	@Success		200		{object}	entity.ResponseContext	"Success"
//	@Failure		500		{object}	entity.ResponseContext	"Internal server error"
//	@Router			/products/search [get]
func (pr *Product) SearchProduct(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	var productsName = c.Query("name")
	

	if productsName == "" {
		c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "", gin.H{}))
		return
	}
	pr.productRepo = application.NewProductApplication(pr.Persistence, c)

	products, searchErr := pr.productRepo.SearchProduct(productsName)
	if searchErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, searchErr.Error(), ""))
		return
	} else if len(products) == 0 {
		c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "No such product found", ""))
		return
	}
	results := map[string]interface{}{
		"results" : products,
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "", results))
}

func (pr *Product) UpdateProductSearchDB(c *gin.Context) {
	

	pr.productRepo = application.NewProductApplication(pr.Persistence, c)
	updateErr := pr.productRepo.UpdateProductsInSearchDB()

	if updateErr != nil {
		log.Println(updateErr)
	}
}