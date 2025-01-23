package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/polarisjrex0406/federico-app/dto"
	"github.com/polarisjrex0406/federico-app/services"
	"github.com/polarisjrex0406/federico-app/utils"
)

type (
	UserHandler interface {
		DoTransaction(c *gin.Context)
		GetBalance(c *gin.Context)
	}

	userHandler struct {
		userService services.UserService
	}
)

func NewUserHandler(us services.UserService) UserHandler {
	return &userHandler{
		userService: us,
	}
}

// @Summary Updates user balance
// @Description Updates user balance
// @Tags User
// @Accept  json
// @Produce  json
// @Param Source-type header string true "Source type of the request"
// @Param body body dto.UserDoTransactionRequest true "Data for transaction"
// @Param userId path uint true "User identifier"
// @Success 200 {object} utils.Response "Successful response"
// @Router /user/{userId}/transaction [post]
func (h *userHandler) DoTransaction(c *gin.Context) { // Read the request body
	// Get path from the context
	userIdAsString := c.Param("userId")
	userId, err := utils.StringToUint(userIdAsString)
	if err != nil {
		utils.SendResponseFailure(c, http.StatusBadRequest, dto.CODE_FAILED_REQUEST_PATH_NOT_VALID, dto.MESSAGE_FAILED_GET_PATH_PARAM, nil)
		return
	}

	// Retrieve the validated request from the context
	validatedReq, exists := c.Get("validatedRequest")
	if !exists {
		utils.SendResponseFailure(c, http.StatusBadRequest, dto.CODE_FAILED_REQUEST_BODY_NOT_VALID, dto.MESSAGE_FAILED_GET_REQUEST_BODY, nil)
		return
	}
	// Type assert the request to the expected type
	assertedReq, ok := validatedReq.(dto.UserDoTransactionRequest)
	if !ok {
		utils.SendResponseFailure(c, http.StatusBadRequest, dto.CODE_FAILED_REQUEST_BODY_NOT_VALID, dto.MESSAGE_FAILED_GET_REQUEST_BODY, nil)
		return
	}

	// Update user balance
	res, err := h.userService.DoTransaction(userId, assertedReq)
	if err != nil {
		if err == dto.ErrTransactionAlreadyExists {
			utils.SendResponseFailure(c, http.StatusBadRequest, dto.CODE_FAILED_TRANSACTION_ALREADY_EXISTS, dto.MESSAGE_FAILED_USER_DO_TRANSACTION, nil)
			return
		} else if err == dto.ErrNotEnoughBalance {
			utils.SendResponseFailure(c, http.StatusBadRequest, dto.CODE_FAILED_BALANCE_NOT_ENOUGH, dto.MESSAGE_FAILED_USER_DO_TRANSACTION, nil)
			return
		}
		utils.SendResponseFailure(c, http.StatusInternalServerError, dto.CODE_FAILED_INTERNAL_PROCESS, dto.MESSAGE_FAILED_USER_DO_TRANSACTION, nil)
		return
	}
	// Return success response
	utils.SendResponseSuccess(c, http.StatusOK, dto.CODE_SUCCESS, dto.MESSAGE_SUCCESS_USER_DO_TRANSACTION, *res)
}

// @Summary Gets current user balance
// @Description Gets current user balance
// @Tags User
// @Accept  json
// @Produce  json
// @Param userId path uint true "User identifier"
// @Success 200 {object} utils.Response "Successful response"
// @Router /user/{userId}/balance [get]
func (h *userHandler) GetBalance(c *gin.Context) {
	// Get path from the context
	userIdAsString := c.Param("userId")
	userId, err := utils.StringToUint(userIdAsString)
	if err != nil {
		utils.SendResponseFailure(c, http.StatusBadRequest, dto.CODE_FAILED_REQUEST_PATH_NOT_VALID, dto.MESSAGE_FAILED_GET_PATH_PARAM, nil)
		return
	}
	// Get current user balance
	res, err := h.userService.GetBalance(userId)
	if err != nil {
		utils.SendResponseFailure(c, http.StatusInternalServerError, dto.CODE_FAILED_INTERNAL_PROCESS, dto.MESSAGE_FAILED_USER_GET_BALANCE, nil)
		return
	}
	// Return success response
	utils.SendResponseSuccess(c, http.StatusOK, dto.CODE_SUCCESS, dto.MESSAGE_SUCCESS_USER_GET_BALANCE, *res)
}
