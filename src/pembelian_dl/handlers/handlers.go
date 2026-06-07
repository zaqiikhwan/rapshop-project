package handlers

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"rapsshop-project/entities"
	"rapsshop-project/model"
	"rapsshop-project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type pembelianHandler struct {
	ServicePembelianDL model.PembelianDLUsecase
	AdminRepository    model.AdminRepository
	PaymentUsecase     model.MetodePembayaranUsecase
	StockDLUsecase     model.StockDLUsecase
}

func NewPembelianHandler(r *gin.RouterGroup, usecaseBeliDL model.PembelianDLUsecase, adminRepo model.AdminRepository, usecasePayment model.MetodePembayaranUsecase, stockDLUsecase model.StockDLUsecase, jwtMiddleware gin.HandlerFunc) {
	pembelianHandler := &pembelianHandler{ServicePembelianDL: usecaseBeliDL, AdminRepository: adminRepo, PaymentUsecase: usecasePayment, StockDLUsecase: stockDLUsecase}
	r.POST("/pembelian", pembelianHandler.HandlerPembelian)
	r.POST("/new/pembelian", pembelianHandler.NewHandlerPembelian)
	r.POST("/pembelian/status", pembelianHandler.HandlerStatus)
	r.GET("/checkout/preview", pembelianHandler.GetCheckoutPreview)
	r.GET("/pembelians", jwtMiddleware, pembelianHandler.GetAllDataPembelian)
	r.GET("/pembelian/total", jwtMiddleware, pembelianHandler.GetTotalPembelian)
	r.GET("/pembelian/:id/tracking", pembelianHandler.GetTrackingPembelian)
	r.GET("/pembelian/:id", pembelianHandler.GetDetailPembelian) // detail data dari database
	r.GET("/pembelian/status/:id", pembelianHandler.GetStatus)   // detail status dari midtrans
	r.PATCH("/pembelian/:id", jwtMiddleware, pembelianHandler.UpdateStatusPengiriman)
	r.PATCH("/pembelian/button/:id", pembelianHandler.NewUpdateButton)
	r.PATCH("/pembelian/confirm/:id", jwtMiddleware, pembelianHandler.NewUpdateConfirmPayment)
	r.Static("/public", "./public/payment")
	r.PATCH("/upload/:id", pembelianHandler.UploadFile)
}

// catetan!!
// di sini masih perlu recognition lebih banyak
// belum clean...

func (ph *pembelianHandler) HandlerPembelian(c *gin.Context) {
	var input entities.PembelianDL

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "bad request for binding input", err)
		return
	}
	input.ID = uuid.NewString()

	responseBody, err := ph.ServicePembelianDL.ChargeAndCreate(input)
	if err != nil {
		ph.respondPurchaseCreationError(c, "failed to create transaction", err, http.StatusBadGateway)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "transaction successfully created", responseBody)
}

func (ph *pembelianHandler) GetCheckoutPreview(c *gin.Context) {
	jumlahDL, err := strconv.Atoi(c.Query("jumlah_dl"))
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed when convert jumlah_dl to int", err)
		return
	}

	preview, err := ph.ServicePembelianDL.GetCheckoutPreview(jumlahDL)
	if err != nil {
		ph.respondPurchaseCreationError(c, "failed create checkout preview", err, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success create checkout preview", preview)
}

func (ph *pembelianHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	id := c.Param("id")

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "get form err: "+err.Error(), err)
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed open uploaded file", err)
		return
	}
	defer openedFile.Close()

	sniff := make([]byte, 512)
	n, err := openedFile.Read(sniff)
	if err != nil && !errors.Is(err, io.EOF) {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed read uploaded file", err)
		return
	}

	ext, err := model.ValidatePaymentProofUpload(file.Filename, file.Size, sniff[:n])
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "format picture not allowed", err)
		return
	}

	file.Filename = uuid.NewString() + ext
	if err := c.SaveUploadedFile(file, "./public/payment/"+file.Filename); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed upload file", err)
		return
	}

	linkImage := os.Getenv("HOST_URL") + "/api/v1/public/" + file.Filename

	var input entities.PembelianDL

	input.BuktiPembayaran = linkImage

	if err := ph.ServicePembelianDL.UpdateTambahBukti(id, input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed add bukti pembayaran", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success upload file", linkImage)
}

func (ph *pembelianHandler) NewHandlerPembelian(c *gin.Context) {
	var input entities.PembelianDL

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "bad request for binding input", err)
		return
	}

	input.ID = uuid.NewString()

	if err := ph.ServicePembelianDL.CreateDataPembelianManual(input); err != nil {
		ph.respondPurchaseCreationError(c, "failed create new data pembelian", err, http.StatusInternalServerError)
		return
	}

	paymentMethod, err := ph.PaymentUsecase.GetDetailPembayaranByIndex(input.MetodeTransfer)

	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "payment method not found", err)
		return
	}
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed fetch payment method", err)
		return
	}

	paymentInstruction := model.NewManualPaymentInstruction(paymentMethod)
	response := model.NewPembelianManualCreateResponse(input.ID, paymentInstruction)

	utils.SuccessResponse(c, http.StatusCreated, "transaction successfully created", response)
}

func (ph *pembelianHandler) respondPurchaseCreationError(c *gin.Context, message string, err error, fallbackStatus int) {
	switch {
	case errors.Is(err, model.ErrInvalidJumlahDL):
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, message, err)
	case errors.Is(err, model.ErrInsufficientStock):
		utils.FailureOrErrorResponse(c, http.StatusConflict, message, err)
	case errors.Is(err, gorm.ErrRecordNotFound):
		utils.FailureOrErrorResponse(c, http.StatusNotFound, message, err)
	default:
		utils.FailureOrErrorResponse(c, fallbackStatus, message, err)
	}
}

func (ph *pembelianHandler) NewUpdateButton(c *gin.Context) {
	id := c.Param("id")

	var input entities.PembelianDL

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "bad request for binding input", err)
		return
	}

	if err := ph.ServicePembelianDL.UpdateStatusButtonBayar(id, input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed update button bayar", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success patch button bayar", input.ButtonBayar)
}

func (ph *pembelianHandler) NewUpdateConfirmPayment(c *gin.Context) {
	idAdmin := c.MustGet("id").(string)

	if idAdmin == "" {
		utils.FailureOrErrorResponse(c, http.StatusUnauthorized, "credential not found", errors.New("unathorized access, please login first"))
		return
	}

	admin, err := ph.AdminRepository.GetByID(idAdmin)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "credentials not found", err)
		return
	}

	id := c.Param("id")

	var input entities.PembelianDL

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "bad request for binding input", err)
		return
	}

	input.EditorStatus = admin.Username

	if err := ph.ServicePembelianDL.UpdateStatusPembayaranAdmin(id, input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed update konfirmasi bayar", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success patch konfirmasi bayar", input.StatusPembayaran)
}

func (ph *pembelianHandler) HandlerStatus(c *gin.Context) {
	var notifPayload map[string]interface{}
	err := json.NewDecoder(c.Request.Body).Decode(&notifPayload)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed when decode json payload", err)
		return
	}
	orderId, exist := notifPayload["order_id"].(string)
	if !exist {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "order_id not found", errors.New("order_id missing in payload"))
		return
	}

	// verify Midtrans signature: sha512(order_id + status_code + gross_amount + ServerKey)
	statusCode, _ := notifPayload["status_code"].(string)
	grossAmount, _ := notifPayload["gross_amount"].(string)
	signatureKey, _ := notifPayload["signature_key"].(string)
	sum := sha512.Sum512([]byte(orderId + statusCode + grossAmount + os.Getenv("AUTHORIZATION_VALUE")))
	if hex.EncodeToString(sum[:]) != signatureKey {
		utils.FailureOrErrorResponse(c, http.StatusUnauthorized, "invalid midtrans signature", errors.New("signature mismatch"))
		return
	}

	if err := ph.ServicePembelianDL.UpdateStatusPembayaran(orderId); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed update status pembelian", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "successfully update status pembelian", nil)
}

func (ph *pembelianHandler) GetAllDataPembelian(c *gin.Context) {
	_start := c.Query("_start")
	_end := c.Query("_end")
	queue := c.Query("queue")

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

	allData, lenData, err := ph.ServicePembelianDL.GetAllPembelian(_startInt, _endInt, queue)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when fetch all data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success fetch all data", gin.H{"data": allData, "total": lenData})
}

func (ph *pembelianHandler) GetStatus(c *gin.Context) {
	id := c.Param("id")
	responseBody, err := ph.ServicePembelianDL.GetLiveStatus(id)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed fetch transaction status", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "transaction found", responseBody)
}

func (ph *pembelianHandler) GetTrackingPembelian(c *gin.Context) {
	id := c.Param("id")
	tracking, err := ph.ServicePembelianDL.GetTrackingByID(id)
	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "data pembelian not found", err)
		return
	}
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed fetch tracking data pembelian", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success fetch tracking data pembelian", tracking)
}

func (ph *pembelianHandler) GetDetailPembelian(c *gin.Context) {
	id := c.Param("id")
	allData, err := ph.ServicePembelianDL.GetDetailByID(id)
	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "data pembelian not found", err)
		return
	}

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed fetch data pembelian", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success fetch all data", allData)
}

func (ph *pembelianHandler) GetTotalPembelian(c *gin.Context) {
	_date := c.Query("_date")

	rekapBeli, err := ph.ServicePembelianDL.GetTotal(_date)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed fetch total pembelian data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success fetch total pembelian data", rekapBeli)
}

func (ph *pembelianHandler) UpdateStatusPengiriman(c *gin.Context) {
	idAdmin := c.MustGet("id").(string)

	if idAdmin == "" {
		utils.FailureOrErrorResponse(c, http.StatusUnauthorized, "credential not found", errors.New("unathorized access, please login first"))
		return
	}

	admin, err := ph.AdminRepository.GetByID(idAdmin)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "credentials not found", err)
		return
	}
	id := c.Param("id")

	var input entities.PembelianDL

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "bad request for binding input", err)
		return
	}

	input.EditorStatus = admin.Username

	if err := ph.ServicePembelianDL.UpdateStatusPengiriman(id, input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed update data status pengiriman", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success update data pengiriman", input.StatusPengiriman)
}
