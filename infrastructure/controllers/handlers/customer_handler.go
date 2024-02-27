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


//	@Summary		Save Customer
//	@Description	Saves a new customer to the database.
//	@Tags			Customer
//	@Accept			json
//	@Produce		json
//	@Param			customer	body		customer_entity.Customer			true	"Customer object to be saved"
//	@Success		201			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/customers [post]
func (cr Customer) SaveCustomer(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	customer := customer_entity.Customer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusUnprocessableEntity, responseContextData.ResponseData(entity.StatusFail, "Invalid JSON", ""))
		return
	}

	cr.CustomerRepo = application.NewCustomerApplication(cr.Persistence, c)

	savedCustomer, saveErr := cr.CustomerRepo.SaveCustomer(&customer)

	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, saveErr["db_error"], ""))
		return
	}

	c.JSON(http.StatusCreated, responseContextData.ResponseData(entity.StatusSuccess, "Customer saved successfully", savedCustomer))
}

//	@Summary		Get All Customers
//	@Description	Retrieves all customers from the database.
//	@Tags			Customer
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.ResponseContext	"Success"
//	@Failure		500	{object}	entity.ResponseContext	"Internal server error"
//	@Router			/customers [get]
func (cr Customer) GetAllCustomers(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	cr.CustomerRepo = application.NewCustomerApplication(cr.Persistence, c)

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

//	@Summary		Get Customer
//	@Description	Retrieves a customer by its ID.
//	@Tags			Customer
//	@Accept			json
//	@Produce		json
//	@Param			customer_id	path		int						true	"Customer ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/customers/{customer_id} [get]
func (cr Customer) GetCustomer(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	customerID, err := strconv.ParseInt(c.Param("customer_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	cr.CustomerRepo = application.NewCustomerApplication(cr.Persistence, c)

	customer, err := cr.CustomerRepo.GetCustomer(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, fmt.Sprintf("Customer %v obtained", customerID), customer))
}

//	@Summary		Delete Customer
//	@Description	Deletes a customer by its ID.
//	@Tags			Customer
//	@Accept			json
//	@Produce		json
//	@Param			customer_id	path		int						true	"Customer ID"
//	@Success		200			{object}	entity.ResponseContext	"Success"
//	@Failure		400			{object}	entity.ResponseContext	"Bad request"
//	@Failure		500			{object}	entity.ResponseContext	"Internal server error"
//	@Router			/customers/{customer_id} [delete]
func (cr Customer) DeleteCustomer(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	customerID, err := strconv.ParseInt(c.Param("customer_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}
	cr.CustomerRepo = application.NewCustomerApplication(cr.Persistence, c)

	deleteErr := cr.CustomerRepo.DeleteCustomer(customerID)
	// TODO: when deleting a customer, need to delete all the inventory in it

	if deleteErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, deleteErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Customer deleted successfully", ""))
}

//	@Summary		Update Customer
//	@Description	Updates a customer.
//	@Tags			Customer
//	@Accept			json
//	@Produce		json
//	@Param			customer_id	path		int							true	"Customer ID"
//	@Param			customer	body		customer_entity.Customer	true	"Customer object to be updated"
//	@Success		200			{object}	entity.ResponseContext		"Success"
//	@Failure		400			{object}	entity.ResponseContext		"Bad request"
//	@Failure		404			{object}	entity.ResponseContext		"Customer not found"
//	@Failure		422			{object}	entity.ResponseContext		"Unprocessable entity"
//	@Failure		500			{object}	entity.ResponseContext		"Internal server error"
func (cr Customer) UpdateCustomer(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	customerID, err := strconv.ParseInt(c.Param("customer_id"), 10, 64)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, responseContextData.ResponseData(entity.StatusFail, "Invalid Customer ID", ""))
		return
	}

	// Check if the Customer exists
	cr.CustomerRepo = application.NewCustomerApplication(cr.Persistence, c)

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

	cr.CustomerRepo = application.NewCustomerApplication(cr.Persistence, c)

	// Update the Customer
	updatedCustomer, updateErr := cr.CustomerRepo.UpdateCustomer(existingCustomer)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, responseContextData.ResponseData(entity.StatusFail, updateErr.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Customer updated successfully", updatedCustomer))
}
