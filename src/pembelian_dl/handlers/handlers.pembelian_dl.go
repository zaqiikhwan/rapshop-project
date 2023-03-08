package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"rapsshop-project/database/mysql"
	"rapsshop-project/entities"
	"rapsshop-project/model"
	"rapsshop-project/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// Qris Endpoint (production)
	productionEnv string = "https://api.midtrans.com/v2/charge"

	// Qris Endpoint (sandbox)
	// sandboxEnv string = "https://api.sandbox.midtrans.com/v2/charge"

	qris string = "qris"
	gopay string = "gopay"
	shopeepay = "shopeepay"
	bca = "bca"
	bri = "bri"
	bni = "bni"
)

type pembelianHandler struct {
	ServicePembelianDL model.PembelianDLUsecase
}
func NewPembelianHandler(r *gin.RouterGroup, usecaseBeliDL model.PembelianDLUsecase){
	pembelianHandler := &pembelianHandler{ServicePembelianDL: usecaseBeliDL}
	r.POST("/pembelian", pembelianHandler.HandlerPembelian)
	r.GET("/pembelians", pembelianHandler.GetAllDataPembelian)
	r.GET("/pembelian/:id", pembelianHandler.GetDetailPembelian)
	r.POST("/pembelian/status", pembelianHandler.HandlerStatus)
	r.GET("/status", pembelianHandler.GetStatus)
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

	var harga entities.HargaDL
	if err := mysql.InitDatabase().Order("id desc").First(&harga).Error; err != nil {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "price not found", err)
		return
	}

	// https://api.sandbox.midtrans.com/v2/96d23a08-bdb5-4282-ae3a-6cfaa1eed867/status

	midtransData := model.NewMidtransData(paymentType, jenisBank, input, harga)
	result := midtransData.IniDataPembelian()
	data, err := json.Marshal(&result)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed when change data to json", err)
		return
	}

	payload := strings.NewReader(string(data))
	// fmt.Println("json = ", string(data))

	req, err := http.NewRequest("POST", productionEnv, payload)
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

	if err := ph.ServicePembelianDL.CreateDataPembelian(input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed create data penmbelian to database", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "transaction successfully created", responseBody)
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
	if err := ph.ServicePembelianDL.UpdateStatusPembelian(orderId); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed update status pembelian", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "successfully update status pembelian", nil)
}

func (ph *pembelianHandler) GetAllDataPembelian(c *gin.Context) {
	allData, err := ph.ServicePembelianDL.GetAllPembelian()
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when fetch all data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success fetch all data", allData)
}

func (ph *pembelianHandler) GetStatus(c *gin.Context) {
	req, err := http.NewRequest("GET", "https://api.sandbox.midtrans.com/v2/96d23a08-bdb5-4282-ae3a-6cfaa1eed867/status", nil)
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

	utils.SuccessResponse(c, http.StatusCreated, "transaction successfully created", responseBody)
}

func (ph *pembelianHandler) GetDetailPembelian(c *gin.Context) {
	id := c.Param("id")
	allData, err := ph.ServicePembelianDL.GetDetailByID(id)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when fetch all data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success fetch all data", allData)
}