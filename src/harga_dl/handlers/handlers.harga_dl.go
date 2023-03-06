package handlers

import (
	"net/http"
	"rapsshop-project/model"
	"rapsshop-project/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type hargaDLHandler struct {
	HargaDLUsecase model.HargaDLUsecase
}

func NewHargaDLHandler(r *gin.RouterGroup, hdlh model.HargaDLUsecase, jwtMiddleware gin.HandlerFunc) {
	hargaDLHandler := &hargaDLHandler{HargaDLUsecase: hdlh}
	r.POST("/price", jwtMiddleware, hargaDLHandler.CreateNewPrice)
	r.GET("/price", hargaDLHandler.GetLatestPrice)
	r.PATCH("/price", jwtMiddleware, hargaDLHandler.UpdateLatestPrice)
	r.DELETE("/price", jwtMiddleware, hargaDLHandler.DeleteLatesPrice)
}

func (hdlh *hargaDLHandler) CreateNewPrice(c *gin.Context){
	var input model.InputHargaDL

	if err := c.BindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "must bind with json", err)
		return
	}

	if err := hdlh.HargaDLUsecase.CreateNewPrice(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when add new price", err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "success add new price", input)
}

func (hdlh *hargaDLHandler) GetLatestPrice(c *gin.Context) {
	latestPrice, err := hdlh.HargaDLUsecase.GetLatestPrice()
	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "price not found", err)
		return
	}

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed get latest price", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success add new price", latestPrice)
}

func (hdlh *hargaDLHandler) UpdateLatestPrice(c *gin.Context) {
	var input model.InputHargaDL

	if err := c.BindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "must bind with json", err)
		return
	}

	updatedPrice, err := hdlh.HargaDLUsecase.UpdateLatestPrice(&input)

	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "price not found", err)
		return
	}

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed update latest price", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success update latest price", updatedPrice)
}

func (hdlh *hargaDLHandler) DeleteLatesPrice(c *gin.Context) {
	if err := hdlh.HargaDLUsecase.DeleteLatestPrice(); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed delete latest price", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success delete latest price", nil)

}