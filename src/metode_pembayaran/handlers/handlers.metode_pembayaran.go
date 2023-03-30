package handlers

import (
	"net/http"
	"rapsshop-project/entities"
	"rapsshop-project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type metodePembayaranHandler struct {
	MetodePembayaranUsecase entities.MetodePembayaranUsecase
}
func NewMetodePembayaranHandler(r *gin.RouterGroup, mpu entities.MetodePembayaranUsecase, jwtMiddleware gin.HandlerFunc) {
	handlerMetodeBayar := &metodePembayaranHandler{MetodePembayaranUsecase: mpu}
	r.POST("/payment", jwtMiddleware, handlerMetodeBayar.CreateNewPayment)
	r.GET("/payment/:id", handlerMetodeBayar.GetDetailPaymentByID)
	r.GET("/payments", handlerMetodeBayar.GetAllPayment)
	r.PATCH("/payment/:id", jwtMiddleware, handlerMetodeBayar.PatchDetailPaymentByID)
	r.DELETE("/payment/:id", jwtMiddleware, handlerMetodeBayar.DeletePaymentByID)
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

func (mph *metodePembayaranHandler) GetDetailPaymentByID(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "error convert string to uint", err)
		return
	}

	detail, err := mph.MetodePembayaranUsecase.GetDetailPembayaranByID(uint(idUint))

	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "payment method not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "payment method found", detail)
}

func (mph *metodePembayaranHandler) PatchDetailPaymentByID(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "error convert string to uint", err)
		return
	}

	var pacthMethod entities.InputMetodePembayaran

	if err := c.BindJSON(&pacthMethod); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "must bind with json", err)
		return
	}

	if err := mph.MetodePembayaranUsecase.PatchDetailPembayaranByID(uint(idUint), &pacthMethod); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when patch detail payment method", err)
		return
	}

	detail, err := mph.MetodePembayaranUsecase.GetDetailPembayaranByID(uint(idUint))
	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "payment id not found", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "success patch detail payment method", detail)
}

func (mph *metodePembayaranHandler) DeletePaymentByID(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "error convert string to uint", err)
		return
	}

	if err := mph.MetodePembayaranUsecase.DeletePembayaranByID(uint(idUint)); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when delete payment method", err)
		return		
	}

	utils.SuccessResponse(c, http.StatusOK, "success delete payment method", nil)
}