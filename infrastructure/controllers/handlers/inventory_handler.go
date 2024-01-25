package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/application"
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


func (in *Inventory) GetInventory(c *gin.Context) {
	productID, err := strconv.ParseInt((c.Param("product_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	in.inventoryHandlerRepo = application.NewInventoryApplication(in.Persistence)

	inventory, err := in.inventoryHandlerRepo.GetInventory(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, inventory)
}



func (in *Inventory) UpdateInventory(c *gin.Context) {
	productIDofInventory, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Check if the inventory exists
	in.inventoryHandlerRepo = application.NewInventoryApplication(in.Persistence)

	existinginventory, err := in.inventoryHandlerRepo.GetInventory(productIDofInventory)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "inventory not found"})
		return
	}

	// Bind the JSON request body to the existing inventory
	if err := c.ShouldBindJSON(&existinginventory); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid JSON"})
		return
	}

	in.inventoryHandlerRepo = application.NewInventoryApplication(in.Persistence)


	// Update the inventory
	updatedinventory, updateErr := in.inventoryHandlerRepo.UpdateInventory(existinginventory)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedinventory)
}
