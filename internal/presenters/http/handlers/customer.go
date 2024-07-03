package handlers

import (
	"encoding/json"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomerHandler struct {
	*BaseHandler
	customerUsecase CustomerUsecase
}

func NewCustomerHandler(
	baseHandler *BaseHandler,
	customerUsecase CustomerUsecase,
) *CustomerHandler {
	return &CustomerHandler{
		BaseHandler:     baseHandler,
		customerUsecase: customerUsecase,
	}
}

type CreateCustomerRequest struct {
	CreateCustomerRequestBody
}

type CreateCustomerRequestBody struct {
	Message struct {
		Data []byte `json:"data"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type Customer struct {
	Uuid string `json:"uuid"`
}

func ModelToCustomerResponse(user *models.Customer) *Customer {
	return &Customer{
		Uuid: user.Uuid.String(),
	}
}

// CreateCustomer godoc
// @Summary Create user
// @Description Create user
// @Tags users
// @Accept  json
// @Produce application/json
// @Security BearerAuth
// @Param raw body handlers.CreateCustomerRequestBody true "Create user"
// @Success 201 {object} handlers.Response{data=handlers.Customer}
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 401 {object} handlers.ErrorResponse
// @Failure 403 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /v1/customers [post]
func (ch CustomerHandler) CreateCustomer(c echo.Context) error {
	var req CreateCustomerRequest
	if err := c.Bind(&req); err != nil {
		return ch.ErrorResponse(c, http.StatusBadRequest, "could not bind request", err)
	}
	var customerReq Customer
	if err := json.Unmarshal(req.Message.Data, &customerReq); err != nil {
		return ch.ErrorResponse(c, http.StatusBadRequest, "could not unmarshal request", err)
	}
	customer, err := ch.customerUsecase.CreateCustomer(customerReq.Uuid)
	if err != nil {
		return err
	}
	return ch.SuccessResponse(c, http.StatusCreated, "customerReq was successfully created", ModelToCustomerResponse(customer))
}
