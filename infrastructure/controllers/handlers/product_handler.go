package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/application"
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/product_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)




type Product struct {
	productRepo product_repository.ProductRepository
	Persistence *base.Persistence
}



func NewProduct(p *base.Persistence) *Product {
	return &Product{
		Persistence: p,
	}
}

func (pr *Product) SaveProduct(c *gin.Context) {
	var product = entity.Product{}

	if err:= c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return	
	}

	fmt.Println(product)

	pr.productRepo = application.NewProductApplication(pr.Persistence)

	savedProduct, saveErr := pr.productRepo.SaveProduct(&product)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, savedProduct)

}

func (pr *Product) SaveMultipleProducts(c *gin.Context) {
	var product = []entity.Product{}

	if err:= c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return	
	}

	fmt.Println(product)
	pr.productRepo = application.NewProductApplication(pr.Persistence)

	savedProduct, saveErr := pr.productRepo.SaveMultipleProducts(&product)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, savedProduct)

}

func (pr *Product) GetAllProducts(c *gin.Context) {
	pr.productRepo = application.NewProductApplication(pr.Persistence)

	allProduct, err := pr.productRepo.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allProduct)
}

func (pr *Product) GetProduct(c *gin.Context) {
	productId, err := strconv.ParseInt((c.Param("product_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	pr.productRepo = application.NewProductApplication(pr.Persistence)

	product, err := pr.productRepo.GetProduct(productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, product)
}

func (pr *Product) DeleteProduct(c *gin.Context) {
	productId, err := strconv.ParseInt((c.Param("product_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	pr.productRepo = application.NewProductApplication(pr.Persistence)


	deleteErr := pr.productRepo.DeleteProduct(productId)

	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, deleteErr.Error())
		return
	}

	c.JSON(http.StatusOK, "Product Deleted")

}

func (pr *Product) UpdateProduct(c *gin.Context) {
	productId, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Check if the product exists
	pr.productRepo = application.NewProductApplication(pr.Persistence)

	existingProduct, err := pr.productRepo.GetProduct(productId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Bind the JSON request body to the existing product
	if err := c.ShouldBindJSON(&existingProduct); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid JSON"})
		return
	}

	pr.productRepo = application.NewProductApplication(pr.Persistence)


	// Update the product
	updatedProduct, updateErr := pr.productRepo.UpdateProduct(existingProduct)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}


func (pr *Product) SearchProduct(c *gin.Context) {
	var productsName = c.Query("name")

    if productsName == "" {
        c.JSON(http.StatusOK, gin.H{})
        return
    }
	pr.productRepo = application.NewProductApplication(pr.Persistence)

	products, searchErr := pr.productRepo.SearchProduct(productsName)
	if searchErr != nil {
		c.JSON(http.StatusInternalServerError, searchErr.Error())
		return
	} else if len(products) == 0 {
		c.JSON(http.StatusOK, "No such product found")
		return
	}

	c.JSON(http.StatusOK, products)
}


func (pr *Product) UpdateMongo() {
	pr.productRepo = application.NewProductApplication(pr.Persistence)
	updateErr := pr.productRepo.UpdateProductsInMongo()

	if updateErr != nil {
		fmt.Println("fail to update products in mongbodb")

	}
}