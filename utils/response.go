package utils

import (
	"github.com/gin-gonic/gin"
)

type response struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func SuccessResponse(c *gin.Context, httpCode int, msg string, data interface{}) {
	c.JSON(httpCode, response{
		StatusCode: httpCode,
		Status:     "success, request OK!",
		Message:    msg,
		Data:       data,
	})
}

func FailureOrErrorResponse(c *gin.Context, httpCode int, msg string, err error) {
	var status string
	switch httpCode {
	case 400:
		status = "failed, bad request"
	case 401:
		status = "failed, status unauthorized"
	case 403:
		status = "failed, status forbidden"
	case 404:
		status = "failed, status not found"
	case 409:
		status = "failed, status conflict"
	case 502:
		status = "error, bad gateway"
	default:
		status = "error, internal server error"
	}

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	c.JSON(httpCode, response{
		StatusCode: httpCode,
		Status:     status,
		Message:    msg,
		Data:       gin.H{"error": errMsg},
	})
}
