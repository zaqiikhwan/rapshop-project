package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"rapsshop-project/entities"
	"rapsshop-project/model"
	"rapsshop-project/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// Qris Endpoint (production)
	// productionEnv string = "https://api.midtrans.com/v2/charge"

	// Qris Endpoint (sandbox)
	sandboxEnv string = "https://api.sandbox.midtrans.com/v2/charge"

	qris string = "qris"
	gopay string = "gopay"
	shopeepay = "shopeepay"
	bca = "bca"
	bri = "bri"
	bni = "bni"
)


// 1 -> qris
// 2 -> gopay
// 3 -> shopeepay
// 4 -> bca
// 5 -> bri
// 6 -> bni

type pembelianHandler struct {}

func NewPembelianHandler(r *gin.RouterGroup){
	pembelianHandler := &pembelianHandler{}
	r.POST("/pembelian", pembelianHandler.HandlerPembelian)
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

	midtransData := model.NewMidtransData(paymentType, jenisBank, input)
	result := midtransData.IniDataPembelian()
	data, err := json.Marshal(&result)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "failed when change data to json", err)
		return
	}

	payload := strings.NewReader(string(data))
	// fmt.Println("json = ", string(data))

	req, err := http.NewRequest("POST", sandboxEnv, payload)
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	type responsTransaction struct {
		ID string `json:"id"`
		StatusCode string `json:"status_code"`
		StatusMessage string `json:"status_message"`
	}

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

	utils.SuccessResponse(c, http.StatusCreated, "transaction successfully created", responseBody)
}

func (ph *pembelianHandler) HandlerNotifikasi(c *gin.Context) {

}