package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func SuccessResponse(c *gin.Context, httpCode int, msg string, data interface{}) {
	switch httpCode / 100 {
	case 2:
		c.JSON(httpCode, response{
			StatusCode: httpCode,
			Status:  "success, request OK!",
			Message: msg,
			Data:    data,
		})
	default:
		c.JSON(http.StatusInternalServerError, response{
			StatusCode: httpCode,
			Status:  "error, internal server error",
			Message: msg,
			Data:    nil,
		})
	}
}

func FailureOrErrorResponse(c *gin.Context, httpCode int, msg string, err error) {
	switch httpCode {
		case 400:
			c.JSON(httpCode, response{
				StatusCode: httpCode,
				Status: "failed, bad request",
				Message: msg,
				Data: gin.H{
					"error" : err.Error(),
				},
			})
		case 401:
			c.JSON(httpCode, response{
				StatusCode: httpCode,
				Status: "failed, status unauthorized",
				Message: msg,
				Data: gin.H{
					"error" : err.Error(),
				},
			})
		case 403:
			c.JSON(httpCode, response{
				StatusCode: httpCode,
				Status: "failed, status forbidden",
				Message: msg,
				Data: gin.H{
					"error" : err.Error(),
				},
			})
		case 404:
			c.JSON(httpCode, response{
				StatusCode: httpCode,
				Status: "failed, status not found",
				Message: msg,
				Data: gin.H{
					"error" : err.Error(),
				},
			})
		case 5:
			c.JSON(httpCode, response{
				StatusCode: httpCode,
				Status: "error, internal server error",
				Message: msg,
				Data: gin.H {
					"error": err.Error(),
				},
			})
			
		default:
			c.JSON(httpCode, response{
				StatusCode: httpCode,
				Status: "error, internal server error",
				Message: msg,
				Data: gin.H {
					"error": err.Error(),
				},
			})
	}
}