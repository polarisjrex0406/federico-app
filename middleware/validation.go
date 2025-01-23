package middleware

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/polarisjrex0406/federico-app/dto"
	"github.com/polarisjrex0406/federico-app/utils"
)

func HeaderValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		sourceType := c.GetHeader("Source-type")
		// Validate Source-type
		if sourceType != "game" && sourceType != "server" && sourceType != "payment" {
			utils.SendResponseFailure(c, http.StatusBadRequest, dto.CODE_FAILED_HEADER_SOURCE_TYPE_WRONG, dto.MESSAGE_FAILED_HEADER_SOURCE_TYPE_WRONG, nil)
			return
		}

		// Validate Content-Type
		contentType := c.GetHeader("Content-Type")
		if contentType != "application/json" {
			utils.SendResponseFailure(c, http.StatusUnsupportedMediaType, dto.CODE_FAILED_HEADER_CONTENT_TYPE_WRONG, dto.MESSAGE_FAILED_HEADER_CONTENT_TYPE_WRONG, nil)
			return
		}

		c.Next()
	}
}

func UserDoTransactionValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.UserDoTransactionRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			utils.SendResponseFailure(c, http.StatusBadRequest, dto.CODE_FAILED_REQUEST_BODY_NOT_VALID, dto.MESSAGE_FAILED_GET_REQUEST_BODY, nil)
			return
		}

		// Validate amount
		if !regexp.MustCompile(`^\d+(\.\d{1,2})?$`).MatchString(req.Amount) {
			utils.SendResponseFailure(c, http.StatusBadRequest, dto.CODE_FAILED_REQUEST_BODY_NOT_VALID, dto.MESSAGE_FAILED_GET_REQUEST_BODY, nil)
			return
		}

		// Validate state
		if req.State != "win" && req.State != "loss" {
			utils.SendResponseFailure(c, http.StatusBadRequest, dto.CODE_FAILED_REQUEST_BODY_NOT_VALID, dto.MESSAGE_FAILED_GET_REQUEST_BODY, nil)
			return
		}

		c.Set("validatedRequest", req)
		c.Next()
	}
}
