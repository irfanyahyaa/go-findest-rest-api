package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	response := gin.H{
		"status":  statusCode,
		"message": message,
		"data":    data,
	}

	c.JSON(statusCode, response)
}

func Created(c *gin.Context, message string, data interface{}) {
	SendResponse(c, http.StatusCreated, data, message)
}

func NotFound(c *gin.Context, message string, data interface{}) {
	SendResponse(c, http.StatusNotFound, data, message)
}

func InternalServerError(c *gin.Context, message string, data interface{}) {
	SendResponse(c, http.StatusInternalServerError, data, message)
}
