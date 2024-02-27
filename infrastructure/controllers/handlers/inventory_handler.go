package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/application"
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/inventory_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)




type Inventory struct {
	inventoryHandlerRepo inventory_repository.InventoryHandlerRepository
	Persistence *base.Persistence
}



func NewInventory(p *base.Persistence) *Inventory {
	return &Inventory{
		Persistence: p,
	}
}

//	@Summary		Get Inventory
//	@Description	Retrieves inventory information for a specific product.
//	@Tags			Inventory
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int						true	"Product ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/inventory/{product_id} [get]
func (inv *Inventory) GetInventory(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	productID, err := strconv.ParseInt((c.Param("product_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	inv.inventoryHandlerRepo = application.NewInventoryApplication(inv.Persistence, c)

	inventory, err := inv.inventoryHandlerRepo.GetInventory(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("Inventory for product %v obtained", productID), inventory))
}

//	@Summary		Update Inventory
//	@Description	Updates inventory information for a specific product.
//	@Tags			Inventory
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int						true	"Product ID"
//	@Param			inventory	body		inventory_entity.Inventory		true	"Inventory object to be updated"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		404			{object}	entity.ResponseContext	"Inventory not found"
//	@Failure		422			{object}	entity.ResponseContext	"Unprocessable entity"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/inventory/{product_id} [put]
func (inv *Inventory) UpdateInventory(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	productIDofInventory, err := strconv.ParseInt(c.Param("product_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, "Invalid product ID", ""))
		return
	}

	// Check if the inventory exists
	inv.inventoryHandlerRepo = application.NewInventoryApplication(inv.Persistence, c)

	existingInventory, err := inv.inventoryHandlerRepo.GetInventory(productIDofInventory)
	if err != nil {
		c.JSON(http.StatusNotFound, responseContextData.ResponseData(entity.StatusFail, "Inventory not found", ""))
		return
	}

	// Bind the JSON request body to the existing inventory
	if err := c.ShouldBindJSON(&existingInventory); err != nil {
		c.JSON(http.StatusUnprocessableEntity, responseContextData.ResponseData(entity.StatusFail, "Invalid JSON", ""))
		return
	}

	inv.inventoryHandlerRepo = application.NewInventoryApplication(inv.Persistence, c)

	// Update the inventory
	updatedInventory, updateErr := inv.inventoryHandlerRepo.UpdateInventory(existingInventory)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, updateErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Inventory updated successfully", updatedInventory))
}