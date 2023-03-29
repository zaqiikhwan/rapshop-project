package handlers

import (
	"net/http"
	"rapsshop-project/entities"
	"rapsshop-project/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type metodePembayaranHandler struct {
	MetodePembayaranUsecase entities.MetodePembayaranUsecase
}
func NewMetodePembayaranHandler(r *gin.RouterGroup, mpu entities.MetodePembayaranUsecase, jwtMiddleware gin.HandlerFunc) {
	handlerMetodeBayar := &metodePembayaranHandler{MetodePembayaranUsecase: mpu}
	r.POST("/payment", jwtMiddleware, handlerMetodeBayar.CreateNewPayment)
	r.GET("/payment/:method", handlerMetodeBayar.GetDetailPayment)
	r.GET("/payments", handlerMetodeBayar.GetAllPayment)
	r.PATCH("/payment/:method", jwtMiddleware, handlerMetodeBayar.PatchDetailPayment)
}

func (mph *metodePembayaranHandler) CreateNewPayment(c *gin.Context) {
	var input entities.InputMetodePembayaran

	if err := c.BindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "must bind with json", err)
		return
	}

	if err := mph.MetodePembayaranUsecase.CreateNewPembayaran(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when add new payment method", err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "success add new payment method", input)
}

func (mph *metodePembayaranHandler) GetAllPayment(c *gin.Context) {
	allMethod, err := mph.MetodePembayaranUsecase.GetAllPembayaran()

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when fetch data payment", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success fetch payment method data", allMethod)
}

func (mph *metodePembayaranHandler) GetDetailPayment(c *gin.Context) {
	method := c.Param("method")

	detail, err := mph.MetodePembayaranUsecase.GetDetailPembayaran(method)

	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "payment method not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "payment method found", detail)
}

func (mph *metodePembayaranHandler) PatchDetailPayment(c *gin.Context) {
	method := c.Param("method")
	var pacthMethod entities.InputMetodePembayaran

	if err := c.BindJSON(&pacthMethod); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "must bind with json", err)
		return
	}

	if err := mph.MetodePembayaranUsecase.PatchDetailPembayaran(method, &pacthMethod); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when patch detail payment method", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "success patch detail payment method", pacthMethod.JenisPembayaran)
}