package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Code    string `json:"code" example:"200.0"`
	Message string `json:"message" example:"success"`
	Data    any    `json:"data,omitempty"`
}

type EmptyObj struct{}

func SendResponseSuccess(c *gin.Context, httpStatus int, code string, message string, data any) {
	res := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	c.JSON(httpStatus, res)
}

func SendResponseFailure(c *gin.Context, httpStatus int, code string, message string, data any) {
	res := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	c.AbortWithStatusJSON(httpStatus, res)
}
