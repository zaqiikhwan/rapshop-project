package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"rapsshop-project/database/mysql"
	"rapsshop-project/entities"
	"rapsshop-project/model"
	"rapsshop-project/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	qris string = "qris"
	gopay string = "gopay"
	shopeepay = "shopeepay"
	bca = "bca"
	bri = "bri"
	bni = "bni"
)

type pembelianHandler struct {
	ServicePembelianDL model.PembelianDLUsecase
	AdminRepository model.AdminRepository
	PaymentUsecase entities.MetodePembayaranUsecase
}
func NewPembelianHandler(r *gin.RouterGroup, usecaseBeliDL model.PembelianDLUsecase, adminRepo model.AdminRepository, usecasePayment entities.MetodePembayaranUsecase, jwtMiddleware gin.HandlerFunc){
	pembelianHandler := &pembelianHandler{ServicePembelianDL: usecaseBeliDL, AdminRepository: adminRepo, PaymentUsecase: usecasePayment }
	r.POST("/pembelian", pembelianHandler.HandlerPembelian)
	r.POST("/new/pembelian", pembelianHandler.NewHandlerPembelian)
	r.POST("/pembelian/status", pembelianHandler.HandlerStatus)
	r.GET("/pembelians",jwtMiddleware ,pembelianHandler.GetAllDataPembelian)
	r.GET("/pembelian/total",jwtMiddleware ,pembelianHandler.GetTotalPembelian)
	r.GET("/pembelian/:id", pembelianHandler.GetDetailPembelian) // detail data dari database
	r.GET("/pembelian/status/:id", pembelianHandler.GetStatus) // detail status dari midtrans
	r.PATCH("/pembelian/:id", jwtMiddleware,pembelianHandler.UpdateStatusPengiriman)
	r.PATCH("/pembelian/button/:id", pembelianHandler.NewUpdateButton)
	r.PATCH("/pembelian/confirm/:id", jwtMiddleware, pembelianHandler.NewUpdateConfirmPayment)
	r.Static("/public", "./public/payment")
	r.POST("/upload", pembelianHandler.UploadFile)
}

// catetan!!
// di sini masih perlu recognition lebih banyak
// belum clean...

func (ph *pembelianHandler) HandlerPembelian(c *gin.Context) {
	var paymentType string
	var jenisBank string
	var input entities.PembelianDL

	if err := c.BindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "bad request for binding input", err)
		return
	}
	input.ID = uuid.NewString()

	if input.MetodeTransfer == 4 || input.MetodeTransfer == 5 || input.MetodeTransfer == 6 {
		paymentType = "bank_transfer"
	}

	if input.MetodeTransfer == 1 {
		paymentType = qris
	} else if input.MetodeTransfer == 2 {
		paymentType = gopay
	} else if input.MetodeTransfer == 3 {
		paymentType = shopeepay
	} else if input.MetodeTransfer == 4 {
		jenisBank = bca
	} else if input.MetodeTransfer == 5 {
		jenisBank = bri
	} else if input.MetodeTransfer == 6 {
		jenisBank = bni
	}

	var harga entities.StockDL
	if err := mysql.InitDatabase().Order("id desc").First(&harga).Error; err != nil {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "price not found", err)
		return
	}

	midtransData := model.NewMidtransData(paymentType, jenisBank, input, harga)
	result, totalTransaksi := midtransData.IniDataPembelian()
	data, err := json.Marshal(&result)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed when change data to json", err)
		return
	}

	payload := strings.NewReader(string(data))

	req, err := http.NewRequest("POST", os.Getenv("MIDTRANS"), payload)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed when make new request", err)
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(os.Getenv("AUTHORIZATION_VALUE")))))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when make new transcation request (issue from server)", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when read body json", err)
		return
	}

	type responsTransaction struct {
		ID string `json:"id"`
		StatusCode string `json:"status_code"`
		StatusMessage string `json:"status_message"`
	}

	fmt.Println(len(body))

	if len(body) == 115 {
		var responseBody responsTransaction
		err := json.Unmarshal(body, &responseBody)
		if err != nil {
			utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "something wrong happen when unmarshal the json", err)
			return
		}

		if responseBody.StatusCode == "401" {
			utils.FailureOrErrorResponse(c, http.StatusUnauthorized, "unknown authorization value for selected environment", errors.New(responseBody.StatusMessage))
			return
		}
	
	}

	var responseBody any
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "something wrong happen when unmarshal the json", err)
		return
	}

	input.JumlahTransaksi = totalTransaksi
	input.HargaBeli = harga.HargaBeliDL
	if err := ph.ServicePembelianDL.CreateDataPembelianMidtrans(input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed create data penmbelian to database", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "transaction successfully created", responseBody)
}

func (ph *pembelianHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "get form err: " + err.Error(), err)
		return
	}

	splitFileName := strings.Split(file.Filename, ".")

	if splitFileName[1] != "png" && splitFileName[1] != "jpg" && splitFileName[1] != "jpeg" && splitFileName[1] != "heic" && splitFileName[1] != "heif" {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "format picture not allowed", err)
		return
	}

	var linkImage string

	if file != nil {
		splitFileName := strings.Split(file.Filename, ".")

		rand.Seed(time.Now().Unix())
		str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

		shuff := []rune(str)

		rand.Shuffle(len(shuff), func(i, j int) {
			shuff[i], shuff[j] = shuff[j], shuff[i]
		})

		
		file.Filename = (string(shuff) + "." + splitFileName[1])

		if err := c.SaveUploadedFile(file, "./public/payment/" + file.Filename); err != nil {
			utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed upload file", err)
			return
		}

		middleLink := "/api/v1/public/"

		linkImage = os.Getenv("HOST_URL") + middleLink + file.Filename
	}

	utils.SuccessResponse(c, http.StatusOK, "success upload file", linkImage)
}

// func (ph *pembelianHandler) NewHandlerPembelian(c *gin.Context) {
// 	var input entities.PembelianDL

// 	world := c.PostForm("world")
// 	nama := c.PostForm("nama")
// 	grow_id := c.PostForm("grow_id")
// 	jenis_item := c.PostForm("jenis_item")
// 	jumlah_dl := c.PostForm("jumlah_dl")
// 	wa := c.PostForm("wa")
// 	metode_transfer := c.PostForm("metode_transfer")
// 	gambar, err := c.FormFile("gambar")

// 	// png, jpg, jpeg,heif,heic
// 	if err != nil {
// 		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "get form err: " + err.Error(), err)
// 		return
// 	}

// 	splitFileName := strings.Split(gambar.Filename, ".")

// 	if splitFileName[1] != "png" && splitFileName[1] != "jpg" && splitFileName[1] != "jpeg" && splitFileName[1] != "heic" && splitFileName[1] != "heif" {
// 		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "format picture not allowed", err)
// 		return
// 	}

// 	var linkImage string

// 	if gambar != nil {
// 		splitFileName := strings.Split(gambar.Filename, ".")

// 		rand.Seed(time.Now().Unix())
// 		str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

// 		shuff := []rune(str)

// 		rand.Shuffle(len(shuff), func(i, j int) {
// 			shuff[i], shuff[j] = shuff[j], shuff[i]
// 		})
// 		gambar.Filename = (string(shuff) + "." + splitFileName[1])

// 		if err := c.SaveUploadedFile(gambar, "./public/payment/" + gambar.Filename); err != nil {
// 			utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed upload file", err)
// 			return
// 		}

// 		middleLink := "/api/v1/public/"

// 		linkImage = os.Getenv("HOST_URL") + middleLink + gambar.Filename
// 	}

// 	jenisItemBoolean, err := strconv.ParseBool(jenis_item)

// 	if err != nil {
// 		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed convert jenis_item to boolean", err)
// 		return
// 	}

// 	jumlahDLInt, err := strconv.Atoi(jumlah_dl)
// 	if err != nil {
// 		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed convert jumlah_dl to int", err)
// 		return
// 	}

// 	metodePembayaranInt, err := strconv.Atoi(metode_transfer)
// 	if err != nil {
// 		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed convert jumlah_dl to int", err)
// 		return
// 	}
// 	input.ID = uuid.NewString()

// 	if err := ph.ServicePembelianDL.CreateDataPembelian(world, nama, grow_id, jenisItemBoolean, jumlahDLInt, wa, metodePembayaranInt, linkImage, input.ID); err != nil {
// 		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed create new data pembelian", err)
// 		return
// 	}

// 	paymentMethod, err := ph.PaymentUsecase.GetDetailPembayaranByIndex(metodePembayaranInt)

// 	if err == gorm.ErrRecordNotFound {
// 		utils.FailureOrErrorResponse(c, http.StatusNotFound, "payment method not found", err)
// 		return
// 	}
// 	utils.SuccessResponse(c, http.StatusCreated, "transaction successfully created", map[string]any{"id_transaksi":input.ID, "payment":paymentMethod})
// }

func (ph *pembelianHandler) NewHandlerPembelian(c *gin.Context) {
	var input entities.PembelianDL

	if err := c.BindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "bad request for binding input", err)
		return
	}

	input.ID = uuid.NewString()

	if err := ph.ServicePembelianDL.CreateDataPembelian(input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed create new data pembelian", err)
		return
	}

	paymentMethod, err := ph.PaymentUsecase.GetDetailPembayaranByIndex(input.MetodeTransfer)

	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "payment method not found", err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "transaction successfully created", map[string]any{"id_transaksi":input.ID, "payment":paymentMethod})
}

func (ph *pembelianHandler) NewUpdateButton(c *gin.Context) {
	id := c.Param("id")

	var input entities.PembelianDL

	if err := c.BindJSON(&input); err != nil {
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
	id := c.Param("id")

	var input entities.PembelianDL

	if err := c.BindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "bad request for binding input", err)
		return
	}

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
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "order_id not found", nil)
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

	allData, lenData, err := ph.ServicePembelianDL.GetAllPembelian(_startInt, _endInt)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when fetch all data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success fetch all data", gin.H{"data": allData,"total": lenData})
}

func (ph *pembelianHandler) GetStatus(c *gin.Context) {
	id := c.Param("id")
	linkRequest := fmt.Sprintf("https://api.midtrans.com/v2/%s/status", id)
	req, err := http.NewRequest("GET", linkRequest, nil)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed when make new request", err)
		return
	}

	// add header key "Accept" and value "application/json"
	req.Header.Add("Accept", "application/json")
	// add header key "Content-Type, "application/json""
	req.Header.Add("Content-Type", "application/json")


	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(os.Getenv("AUTHORIZATION_VALUE")))))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when make new transcation request (issue from server)", err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when read body json", err)
		return
	}

	var responseBody any
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "something wrong happen when unmarshal the json", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "transaction found", responseBody)
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

	if err := c.BindJSON(&input); err != nil {
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