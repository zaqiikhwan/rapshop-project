package handlers

import (
	"errors"
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
	AdminRepository model.AdminRepository
}

func NewPenjualanDLHandler(r *gin.RouterGroup, pdlh model.PenjualanDLUsecase, adminRepo model.AdminRepository,jwtMiddleware gin.HandlerFunc) {
	jualDLHandler := &penjualanDLHandler{PenjualanDLUsecase: pdlh,AdminRepository: adminRepo}
	r.POST("/penjualan", jualDLHandler.CreateNewPenjualan)
	r.GET("/penjualans", jwtMiddleware, jualDLHandler.GetAllPenjualan)
	r.GET("/penjualan/:id", jwtMiddleware, jualDLHandler.GetDetailPenjualan)
	r.GET("/rekapitulasi", jwtMiddleware, jualDLHandler.GetRekapByDate)
	r.GET("/profit", jwtMiddleware, jualDLHandler.GetProfit)
	r.GET("/penjualan/total", jwtMiddleware, jualDLHandler.GetTotalPenjualan)
	r.PATCH("/penjualan/:id", jwtMiddleware, jualDLHandler.UpdateStatusPenjualan)
	r.DELETE("/penjualan/:id", jwtMiddleware, jualDLHandler.DeletePenjualan)
}

func (pdlh *penjualanDLHandler) GetRekapByDate(c *gin.Context) {
	_date := c.Query("_date")

	rekapJual, rekapBeli, err := pdlh.PenjualanDLUsecase.GetByDate(_date)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed fetch rekap all data by date", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success fetch rekap all data by date", map[string]any{"tanggal" : _date, "penjualan": rekapJual, "pembelian": rekapBeli})
}

func (pdlh *penjualanDLHandler) GetTotalPenjualan(c *gin.Context) {
	_date := c.Query("_date")

	rekapJual, err := pdlh.PenjualanDLUsecase.GetTotal(_date)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed fetch total penjualan data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success fetch total penjualan data", rekapJual)
}

func (pdlh *penjualanDLHandler) CreateNewPenjualan(c *gin.Context) {
	image, _ := c.FormFile("image")
	jumlahDLstr := c.PostForm("jumlah_dl")
	wa := c.PostForm("whatsapp")
	transfer := c.PostForm("transfer")
	nomorTransfer := c.PostForm("nomor_transfer")
	nama := c.PostForm("nama")

	jumlahDL, err := strconv.Atoi(jumlahDLstr)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "bad format from jumlh_dl", err)
		return
	}

	var harga entities.StockDL
	if err := mysql.InitDatabase().Order("id desc").First(&harga).Error; err != nil {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "price not found", err)
		return
	}

	// total := harga.HargaJualDL * jumlahDL

	if err := pdlh.PenjualanDLUsecase.Create(image, jumlahDL, (harga.HargaJualDL*jumlahDL), wa, transfer, nomorTransfer, nama, harga.HargaJualDL); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed make new penjualan_dl", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "success make new penjualan_dl", jumlahDL)
}

func (pdlh *penjualanDLHandler) GetAllPenjualan(c *gin.Context) {
	_start := c.Query("_start")
	_end := c.Query("_end")

	_startInt, err := strconv.Atoi(_start)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed when convert str to int", err)
		return
	}

	_endInt, err := strconv.Atoi(_end)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed when convert str to int", err)
		return
	}
	allData, len, err := pdlh.PenjualanDLUsecase.GetAll(_startInt, _endInt)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed fetch all data", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "success fetch all data", gin.H{"data": allData,"total": len})
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

func (pdlh *penjualanDLHandler) GetProfit(c *gin.Context) {
	_date := c.Query("_date")

	RekapProfit, err := pdlh.PenjualanDLUsecase.GetProfit(_date)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed fetch profit data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success fetch profit data", RekapProfit)
}

func (pdlh *penjualanDLHandler) UpdateStatusPenjualan(c *gin.Context) {
	idAdmin := c.MustGet("id").(string)

	if idAdmin == "" {
		utils.FailureOrErrorResponse(c, http.StatusUnauthorized, "credential not found", errors.New("unathorized access, please login first"))
		return
	}

	admin, err := pdlh.AdminRepository.GetByID(idAdmin)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "credentials not found", err)
		return
	}
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

	input.EditorStatus = admin.Username

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