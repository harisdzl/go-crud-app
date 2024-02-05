package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/application"
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/customer_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/customer_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)




type Customer struct {
	CustomerRepo customer_repository.CustomerRepository
	Persistence *base.Persistence
}



func NewCustomer(p *base.Persistence) *Customer {
	return &Customer{
		Persistence: p,
	}
}

// SaveCustomer saves a single Customer to the database.
// @Summary Save a single Customer
// @Description SaveCustomer saves a single Customer to the database.
// @Tags Customer
// @Accept json
// @Produce json
// @Param Customer body entity.Customer true "Customer object to be saved"
// @Success 201 {object} entity.Customer "Successfully saved Customer"
// @Failure 400 {object} map[string]string "Invalid JSON"
// @Failure 422 {object} map[string]string "Unprocessable entity"
// @Router /Customers [post]
func (cr Customer) SaveCustomer(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	customer := customer_entity.Customer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusUnprocessableEntity, responseContextData.ResponseData(entity.StatusFail, "Invalid JSON", ""))
		return
	}

	cr.CustomerRepo = application.NewCustomerApplication(cr.Persistence)

	savedCustomer, saveErr := cr.CustomerRepo.SaveCustomer(&customer)

	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, saveErr["db_error"], ""))
		return
	}

	c.JSON(http.StatusCreated, responseContextData.ResponseData(entity.StatusSuccess, "Customer saved successfully", savedCustomer))
}


func (cr Customer) GetAllCustomers(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	cr.CustomerRepo = application.NewCustomerApplication(cr.Persistence)

	allCustomers, err := cr.CustomerRepo.GetAllCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	results := map[string]interface{}{
		"results" : allCustomers,
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "All customers obtained successfully", results))
}

func (cr Customer) GetCustomer(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	customerID, err := strconv.ParseInt(c.Param("customer_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	cr.CustomerRepo = application.NewCustomerApplication(cr.Persistence)

	customer, err := cr.CustomerRepo.GetCustomer(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("Customer %v obtained", customerID), customer))
}

func (cr Customer) DeleteCustomer(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	customerID, err := strconv.ParseInt(c.Param("customer_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	cr.CustomerRepo = application.NewCustomerApplication(cr.Persistence)

	deleteErr := cr.CustomerRepo.DeleteCustomer(customerID)
	// TODO: when deleting a customer, need to delete all the inventory in it

	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, deleteErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Customer deleted successfully", ""))
}

func (cr Customer) UpdateCustomer(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	customerID, err := strconv.ParseInt(c.Param("customer_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, "Invalid Customer ID", ""))
		return
	}

	// Check if the Customer exists
	cr.CustomerRepo = application.NewCustomerApplication(cr.Persistence)

	existingCustomer, err := cr.CustomerRepo.GetCustomer(customerID)
	if err != nil {
		c.JSON(http.StatusNotFound, responseContextData.ResponseData(entity.StatusFail, "Customer not found", ""))
		return
	}

	// Bind the JSON request body to the existing Customer
	if err := c.ShouldBindJSON(&existingCustomer); err != nil {
		c.JSON(http.StatusUnprocessableEntity, responseContextData.ResponseData(entity.StatusFail, "Invalid JSON", ""))
		return
	}

	cr.CustomerRepo = application.NewCustomerApplication(cr.Persistence)

	// Update the Customer
	updatedCustomer, updateErr := cr.CustomerRepo.UpdateCustomer(existingCustomer)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, updateErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Customer updated successfully", updatedCustomer))
}
