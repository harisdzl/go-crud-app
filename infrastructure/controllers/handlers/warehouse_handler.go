package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/application"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/warehouse_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/warehouse_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)




type Warehouse struct {
	WarehouseRepo warehouse_repository.WarehouseRepository
	Persistence *base.Persistence
}



func NewWarehouse(p *base.Persistence) *Warehouse {
	return &Warehouse{
		Persistence: p,
	}
}

// SaveWarehouse saves a single Warehouse to the database.
// @Summary Save a single Warehouse
// @Description SaveWarehouse saves a single Warehouse to the database.
// @Tags Warehouse
// @Accept json
// @Produce json
// @Param Warehouse body entity.Warehouse true "Warehouse object to be saved"
// @Success 201 {object} entity.Warehouse "Successfully saved Warehouse"
// @Failure 400 {object} map[string]string "Invalid JSON"
// @Failure 422 {object} map[string]string "Unprocessable entity"
// @Router /Warehouses [post]
func (pr *Warehouse) SaveWarehouse(c *gin.Context) {
	warehouse := warehouse_entity.Warehouse{}


	if err:= c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return	
	}

	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	savedWarehouse, saveErr := pr.WarehouseRepo.SaveWarehouse(&warehouse)

	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}

	c.JSON(http.StatusCreated, savedWarehouse)

}

func (pr *Warehouse) GetInventoriesInWarehouse(c *gin.Context) {
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	warehouseID, err := strconv.ParseInt((c.Param("warehouse_id")), 10, 64)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	warehouse, warehouseErr := application.NewInventoryApplication(pr.Persistence).GetAllInventoryInWarehouse(warehouseID)

	if warehouseErr != nil {
		c.JSON(http.StatusInternalServerError, warehouseErr.Error())
		return
	}

	c.JSON(http.StatusOK, warehouse)
}
func (pr *Warehouse) GetAllWarehouses(c *gin.Context) {
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	allWarehouse, err := pr.WarehouseRepo.GetAllWarehouses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allWarehouse)
}

func (pr *Warehouse) GetWarehouse(c *gin.Context) {
	warehouseID, err := strconv.ParseInt((c.Param("warehouse_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	Warehouse, err := pr.WarehouseRepo.GetWarehouse(warehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, Warehouse)
}

func (pr *Warehouse) DeleteWarehouse(c *gin.Context) {
	warehouseID, err := strconv.ParseInt((c.Param("warehouse_id")), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)


	deleteErr := pr.WarehouseRepo.DeleteWarehouse(warehouseID)
	// TODO: when deleting a warehouse, need to delete all the inventory in it

	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, deleteErr.Error())
		return
	}
	
	c.JSON(http.StatusOK, "Warehouse deleted")

}

func (pr *Warehouse) UpdateWarehouse(c *gin.Context) {
	warehouseID, err := strconv.ParseInt(c.Param("warehouse_id"), 10, 64)
	
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Warehouse ID"})
		return
	}

	// Check if the Warehouse exists
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	existingWarehouse, err := pr.WarehouseRepo.GetWarehouse(warehouseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
		return
	}

	// Bind the JSON request body to the existing Warehouse
	if err := c.ShouldBindJSON(&existingWarehouse); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid JSON"})
		return
	}

	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)


	// Update the Warehouse
	updatedWarehouse, updateErr := pr.WarehouseRepo.UpdateWarehouse(existingWarehouse)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedWarehouse)
}


func (pr *Warehouse) SearchWarehouse(c *gin.Context) {
	var WarehousesName = c.Query("name")

    if WarehousesName == "" {
        c.JSON(http.StatusOK, gin.H{})
        return
    }
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	Warehouses, searchErr := pr.WarehouseRepo.SearchWarehouse(WarehousesName)
	if searchErr != nil {
		c.JSON(http.StatusInternalServerError, searchErr.Error())
		return
	} else if len(Warehouses) == 0 {
		c.JSON(http.StatusOK, "No such Warehouse found")
		return
	}

	c.JSON(http.StatusOK, Warehouses)
}


func (pr *Warehouse) UpdateWarehouseSearchDB() {
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)
	updateErr := pr.WarehouseRepo.UpdateWarehousesInSearchDB()

	if updateErr != nil {
		fmt.Println("fail to update Warehouses in mongbodb")

	}
}