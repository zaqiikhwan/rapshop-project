package handlers

import (
	"net/http"
	"rapsshop-project/database/mysql"
	"rapsshop-project/entities"
	"rapsshop-project/model"
	"rapsshop-project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type penjualanDLHandler struct {
	PenjualanDLUsecase model.PenjualanDLUsecase
}

func NewPenjualanDLHandler(r *gin.RouterGroup, pdlh model.PenjualanDLUsecase, jwtMiddleware gin.HandlerFunc) {
	jualDLHandler := &penjualanDLHandler{PenjualanDLUsecase: pdlh}
	r.POST("/penjualan", jualDLHandler.CreateNewPenjualan)
	r.GET("/penjualans", jwtMiddleware, jualDLHandler.GetAllPenjualan)
	r.GET("/penjualan/:id", jwtMiddleware, jualDLHandler.GetDetailPenjualan)
	r.PATCH("/penjualan/:id", jwtMiddleware, jualDLHandler.UpdateStatusPenjualan)
	r.DELETE("/penjualan/:id", jwtMiddleware, jualDLHandler.DeletePenjualan)
}

func (pdlh *penjualanDLHandler) CreateNewPenjualan(c *gin.Context) {
	image, _ := c.FormFile("image")
	jumlahDLstr := c.PostForm("jumlah_dl")
	wa := c.PostForm("whatsapp")
	transfer := c.PostForm("transfer")
	nomorTransfer := c.PostForm("nomor_transfer")

	jumlahDL, err := strconv.Atoi(jumlahDLstr)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "bad format from jumlh_dl", err)
		return
	}

	var harga entities.HargaDL
	if err := mysql.InitDatabase().Order("id desc").First(&harga).Error; err != nil {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "price not found", err)
		return
	}

	// total := harga.HargaJualDL * jumlahDL

	if err := pdlh.PenjualanDLUsecase.Create(image, jumlahDL, (harga.HargaJualDL*jumlahDL), wa, transfer, nomorTransfer); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed make new penjualan_dl", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "success make new penjualan_dl", jumlahDL)
}

func (pdlh *penjualanDLHandler) GetAllPenjualan(c *gin.Context) {
	allData, err := pdlh.PenjualanDLUsecase.GetAll()
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed fetch all data", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "success fetch all data", allData)
}

func (pdlh *penjualanDLHandler) GetDetailPenjualan(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "error convert string to uint", err)
		return
	}

	detailData, err := pdlh.PenjualanDLUsecase.GetByID(uint(idUint))
	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "data penjualan not found", err)
		return
	}
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed fetch detail data", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "success fetch detail data", detailData)
}

func (pdlh *penjualanDLHandler) UpdateStatusPenjualan(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "error convert string to uint", err)
		return
	}

	var input entities.PenjualanDL

	if err := c.BindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "must bind with json", err)
		return
	}

	updatedData, err := pdlh.PenjualanDLUsecase.UpdateByID(uint(idUint), input)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed update status data", err)
		return
	}

	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "data penjualan not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success update status data", updatedData)
}

func (pdlh *penjualanDLHandler) DeletePenjualan(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "error convert string to uint", err)
		return
	}

	if err := pdlh.PenjualanDLUsecase.DeleteByID(uint(idUint)); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed delete data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success delete data", nil)
}