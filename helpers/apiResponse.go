package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct{}

var ResponseHelper = ApiResponse{}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (r ApiResponse) Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func (r ApiResponse) Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Status:  "error",
		Message: message,
	})
}

func (r ApiResponse) BadRequest(c *gin.Context, message string) {
	r.Error(c, http.StatusBadRequest, message)
}

func (r ApiResponse) NotFound(c *gin.Context, message string) {
	r.Error(c, http.StatusNotFound, message)
}

func (r ApiResponse) Unauthorized(c *gin.Context, message string) {
	r.Error(c, http.StatusUnauthorized, message)
}

func (r ApiResponse) InternalError(c *gin.Context, message string) {
	r.Error(c, http.StatusInternalServerError, message)
}
