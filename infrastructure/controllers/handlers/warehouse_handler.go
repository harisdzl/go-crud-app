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

// SaveWarehouse saves a single warehouse to the database.
//	@Summary		Save a single warehouse
//	@Description	SaveWarehouse saves a single warehouse to the database.
//	@Tags			Warehouse
//	@Accept			json
//	@Produce		json
//	@Param			Warehouse	body		warehouse_entity.Warehouse		true	"Warehouse object to be saved"
//	@Success		201			{object}	entity.ResponseContext	"Successfully saved warehouse"
//	@Failure		400			{object}	entity.ResponseContext	"Invalid JSON"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/warehouses [post]
func (pr *Warehouse) SaveWarehouse(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	warehouse := warehouse_entity.Warehouse{}
	

	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusUnprocessableEntity, responseContextData.ResponseData(entity.StatusFail, "Invalid JSON", ""))
		return
	}

	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence, c)

	savedWarehouse, saveErr := pr.WarehouseRepo.SaveWarehouse(&warehouse)

	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, saveErr["db_error"], ""))
		return
	}

	c.JSON(http.StatusCreated, responseContextData.ResponseData(entity.StatusSuccess, "Warehouse saved successfully", savedWarehouse))
}

// GetInventoriesInWarehouse retrieves all inventories in a warehouse.
//	@Summary		Get Inventories in Warehouse
//	@Description	Retrieves all inventories in a warehouse.
//	@Tags			Warehouse
//	@Accept			json
//	@Produce		json
//	@Param			warehouse_id	path		int						true	"Warehouse ID"
//	@Success		200				{object}	entity.ResponseContext	"Success"
//	@Failure		400				{object}	entity.ResponseContext	"Bad request"
//	@Failure		500				{object}	entity.ResponseContext	"Internal server error"
//	@Router			/warehouses/{warehouse_id}/inventories [get]
func (pr *Warehouse) GetInventoriesInWarehouse(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	

	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence, c)

	warehouseID, err := strconv.ParseInt(c.Param("warehouse_id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	warehouse, warehouseErr := application.NewInventoryApplication(pr.Persistence, c).GetAllInventoryInWarehouse(warehouseID)

	if warehouseErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, warehouseErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Inventories obtained successfully", warehouse))
}

// GetAllWarehouses retrieves all warehouses.
//	@Summary		Get All Warehouses
//	@Description	Retrieves all warehouses.
//	@Tags			Warehouse
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.ResponseContext	"Success"
//	@Failure		500	{object}	entity.ResponseContext	"Internal server error"
//	@Router			/warehouses [get]
func (pr *Warehouse) GetAllWarehouses(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	

	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence, c)

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

// GetWarehouse retrieves a specific warehouse by ID.
//	@Summary		Get Warehouse
//	@Description	Retrieves a specific warehouse by ID.
//	@Tags			Warehouse
//	@Accept			json
//	@Produce		json
//	@Param			warehouse_id	path		int						true	"Warehouse ID"
//	@Success		200				{object}	entity.ResponseContext	"Success"
//	@Failure		400				{object}	entity.ResponseContext	"Bad request"
//	@Failure		500				{object}	entity.ResponseContext	"Internal server error"
//	@Router			/warehouses/{warehouse_id} [get]
func (pr *Warehouse) GetWarehouse(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	warehouseID, err := strconv.ParseInt(c.Param("warehouse_id"), 10, 64)
	

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence, c)

	warehouse, err := pr.WarehouseRepo.GetWarehouse(warehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("Warehouse %v obtained", warehouseID), warehouse))
}

// DeleteWarehouse deletes a warehouse by ID.
//	@Summary		Delete Warehouse
//	@Description	Deletes a warehouse by ID.
//	@Tags			Warehouse
//	@Accept			json
//	@Produce		json
//	@Param			warehouse_id	path		int						true	"Warehouse ID"
//	@Success		200				{object}	entity.ResponseContext	"Success"
//	@Failure		400				{object}	entity.ResponseContext	"Bad request"
//	@Failure		500				{object}	entity.ResponseContext	"Internal server error"
//	@Router			/warehouses/{warehouse_id} [delete]
func (pr *Warehouse) DeleteWarehouse(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	warehouseID, err := strconv.ParseInt(c.Param("warehouse_id"), 10, 64)
	

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence, c)

	deleteErr := pr.WarehouseRepo.DeleteWarehouse(warehouseID)
	// TODO: when deleting a warehouse, need to delete all the inventory in it

	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, deleteErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Warehouse deleted successfully", ""))
}

// UpdateWarehouse updates a warehouse.
//	@Summary		Update Warehouse
//	@Description	Updates a warehouse.
//	@Tags			Warehouse
//	@Accept			json
//	@Produce		json
//	@Param			warehouse_id	path		int						true	"Warehouse ID"
//	@Success		200				{object}	entity.ResponseContext	"Success"
//	@Failure		400				{object}	entity.ResponseContext	"Bad request"
//	@Failure		404				{object}	entity.ResponseContext	"Not found"
//	@Failure		422				{object}	entity.ResponseContext	"Unprocessable entity"
//	@Router			/warehouses/{warehouse_id} [put]
func (pr *Warehouse) UpdateWarehouse(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	warehouseID, err := strconv.ParseInt(c.Param("warehouse_id"), 10, 64)
	

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, "Invalid Warehouse ID", ""))
		return
	}

	// Check if the Warehouse exists
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence, c)

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

	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence, c)

	// Update the Warehouse
	updatedWarehouse, updateErr := pr.WarehouseRepo.UpdateWarehouse(existingWarehouse)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, updateErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Warehouse updated successfully", updatedWarehouse))
}

// SearchWarehouse searches for warehouses by name.
//	@Summary		Search Warehouse
//	@Description	Searches for warehouses by name.
//	@Tags			Warehouse
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string					false	"Warehouse name"
//	@Success		200		{object}	entity.ResponseContext	"Success"
//	@Failure		500		{object}	entity.ResponseContext	"Internal server error"
//	@Router			/warehouses/search [get]
func (pr *Warehouse) SearchWarehouse(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	var WarehousesName = c.Query("name")
	

	if WarehousesName == "" {
		c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "", gin.H{}))
		return
	}
	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence, c)

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


func (pr *Warehouse) UpdateWarehouseSearchDB(c *gin.Context) {
	

	pr.WarehouseRepo = application.NewWarehouseApplication(pr.Persistence, c)
	updateErr := pr.WarehouseRepo.UpdateWarehousesInSearchDB()

	if updateErr != nil {
		fmt.Println("fail to update Warehouses in mongbodb")

	}
}