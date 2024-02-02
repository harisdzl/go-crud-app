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

func (inv *Inventory) GetInventory(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	productID, err := strconv.ParseInt((c.Param("product_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	inv.inventoryHandlerRepo = application.NewInventoryApplication(inv.Persistence)

	inventory, err := inv.inventoryHandlerRepo.GetInventory(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("Inventory for product %v obtained", productID), inventory))
}

func (inv *Inventory) UpdateInventory(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	productIDofInventory, err := strconv.ParseInt(c.Param("product_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, "Invalid product ID", ""))
		return
	}

	// Check if the inventory exists
	inv.inventoryHandlerRepo = application.NewInventoryApplication(inv.Persistence)

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

	inv.inventoryHandlerRepo = application.NewInventoryApplication(inv.Persistence)

	// Update the inventory
	updatedInventory, updateErr := inv.inventoryHandlerRepo.UpdateInventory(existingInventory)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, updateErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Inventory updated successfully", updatedInventory))
}