package handlers

import (
	"net/http"
	"rapsshop-project/model"
	"rapsshop-project/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type stockDLHandler struct {
	StockDLUsecase model.StockDLUsecase
}

func NewAdminHandler(r *gin.RouterGroup, sdlh model.StockDLUsecase, jwtMiddleware gin.HandlerFunc) {
	stockDLHandler := &stockDLHandler{StockDLUsecase: sdlh}
	r.POST("/stock", jwtMiddleware, stockDLHandler.CreateNewStock)
	r.GET("/stocks", stockDLHandler.GetAllStockData)
	r.GET("/stock", stockDLHandler.GetLatestStockData)
	r.PATCH("/stock", jwtMiddleware, stockDLHandler.UpdateStockData)
	r.DELETE("/stock", jwtMiddleware, stockDLHandler.DeleteStockData)
}

func (sdlh *stockDLHandler) CreateNewStock(c *gin.Context) {
	var input model.InputStockDL

	if err := c.BindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "must bind with json", err)
		return
	}

	if err := sdlh.StockDLUsecase.CreateNewStock(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when add new stock data", err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "success add new stock data", input.StockDL)
}

func (sdlh *stockDLHandler) GetAllStockData(c *gin.Context) {
	allStock, err := sdlh.StockDLUsecase.GetAllStock()

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when query all data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success fetch all stock data", allStock)
}

func (sdlh *stockDLHandler) GetLatestStockData(c *gin.Context) {
	stockDL, err := sdlh.StockDLUsecase.GetLatestDataStock()

	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "stock not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "stock found", stockDL)
}

func(sdlh *stockDLHandler) UpdateStockData(c *gin.Context) {
	var updateStock model.InputStockDL

	if err := c.BindJSON(&updateStock); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "must bind with json", err)
		return
	}

	result, err := sdlh.StockDLUsecase.UpdateTambahStock(&updateStock)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed to update stock data", err)
		return
	}

	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "stock not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success update data medsos", result)
}

func(sdlh *stockDLHandler) DeleteStockData(c *gin.Context) {
	if err := sdlh.StockDLUsecase.DeleteStock(); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed to delete data stock", err)
		return	
	}

	utils.SuccessResponse(c, http.StatusOK, "success delete data stock", nil)
}