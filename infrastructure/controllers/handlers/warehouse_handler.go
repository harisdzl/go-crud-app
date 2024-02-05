package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/application"
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
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
	responseContextData := entity.ResponseContext{Ctx: c}
	warehouse := warehouse_entity.Warehouse{}

	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusUnprocessableEntity, responseContextData.ResponseData(entity.StatusFail, "Invalid JSON", ""))
		return
	}

	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	savedWarehouse, saveErr := pr.WarehouseRepo.SaveWarehouse(&warehouse)

	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, saveErr["db_error"], ""))
		return
	}

	c.JSON(http.StatusCreated, responseContextData.ResponseData(entity.StatusSuccess, "Warehouse saved successfully", savedWarehouse))
}

func (pr *Warehouse) GetInventoriesInWarehouse(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	warehouseID, err := strconv.ParseInt(c.Param("warehouse_id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	warehouse, warehouseErr := application.NewInventoryApplication(pr.Persistence).GetAllInventoryInWarehouse(warehouseID)

	if warehouseErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, warehouseErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Inventories obtained successfully", warehouse))
}

func (pr *Warehouse) GetAllWarehouses(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	allWarehouses, err := pr.WarehouseRepo.GetAllWarehouses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	results := map[string]interface{}{
		"results" : allWarehouses,
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "All warehouses obtained successfully", results))
}

func (pr *Warehouse) GetWarehouse(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	warehouseID, err := strconv.ParseInt(c.Param("warehouse_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	warehouse, err := pr.WarehouseRepo.GetWarehouse(warehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("Warehouse %v obtained", warehouseID), warehouse))
}

func (pr *Warehouse) DeleteWarehouse(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	warehouseID, err := strconv.ParseInt(c.Param("warehouse_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	deleteErr := pr.WarehouseRepo.DeleteWarehouse(warehouseID)
	// TODO: when deleting a warehouse, need to delete all the inventory in it

	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, deleteErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Warehouse deleted successfully", ""))
}

func (pr *Warehouse) UpdateWarehouse(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	warehouseID, err := strconv.ParseInt(c.Param("warehouse_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, "Invalid Warehouse ID", ""))
		return
	}

	// Check if the Warehouse exists
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	existingWarehouse, err := pr.WarehouseRepo.GetWarehouse(warehouseID)
	if err != nil {
		c.JSON(http.StatusNotFound, responseContextData.ResponseData(entity.StatusFail, "Warehouse not found", ""))
		return
	}

	// Bind the JSON request body to the existing Warehouse
	if err := c.ShouldBindJSON(&existingWarehouse); err != nil {
		c.JSON(http.StatusUnprocessableEntity, responseContextData.ResponseData(entity.StatusFail, "Invalid JSON", ""))
		return
	}

	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	// Update the Warehouse
	updatedWarehouse, updateErr := pr.WarehouseRepo.UpdateWarehouse(existingWarehouse)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, updateErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Warehouse updated successfully", updatedWarehouse))
}

func (pr *Warehouse) SearchWarehouse(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	var WarehousesName = c.Query("name")

	if WarehousesName == "" {
		c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "", gin.H{}))
		return
	}
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)

	warehouses, searchErr := pr.WarehouseRepo.SearchWarehouse(WarehousesName)
	if searchErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, searchErr.Error(), ""))
		return
	} else if len(warehouses) == 0 {
		c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "No such Warehouse found", ""))
		return
	}

	results := map[string]interface{}{
		"results" : warehouses,
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "", results))
}


func (pr *Warehouse) UpdateWarehouseSearchDB() {
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence)
	updateErr := pr.WarehouseRepo.UpdateWarehousesInSearchDB()

	if updateErr != nil {
		fmt.Println("fail to update Warehouses in mongbodb")

	}
}