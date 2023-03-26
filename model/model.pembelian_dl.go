package model

import (
	"rapsshop-project/entities"

	"github.com/midtrans/midtrans-go"
)

type RekapTotalPembelian struct {
	Tanggal string `json:"tanggal"`
	JumlahDL int `json:"jumlah_dl"`
}

type PembelianDLRepository interface {
	Create(input entities.PembelianDL) error
	GetAll(_startInt int , _endInt int) ([]entities.PembelianDL, int, error)
	UpdateStatus(input entities.PembelianDL, id string) error
	GetByID(id string) (entities.PembelianDL, error)
	GetTotalPembelian(date string) ([]RekapTotalPembelian, error)
}

type PembelianDLUsecase interface {
	CreateDataPembelian(input entities.PembelianDL) error
	GetAllPembelian(_startInt int, _endInt int) ([]entities.PembelianDL, int, error)
	UpdateStatusPembayaran(id string) error
	GetDetailByID(id string)(entities.PembelianDL, error)
	GetTotal(date string) ([]RekapTotalPembelian, error)
	UpdateStatusPengiriman(id string, input entities.PembelianDL) error 
}

type MidtransData struct {
	typePayment string 
	jenisBank string
	newPembelian entities.PembelianDL
	harga entities.StockDL
}

func NewMidtransData(typePayment string, jenisBank string, newPembelian entities.PembelianDL, harga entities.StockDL) *MidtransData {
	return &MidtransData{
		typePayment: typePayment,
		jenisBank: jenisBank,
		newPembelian: newPembelian,
		harga: harga,
	}
}

func (m *MidtransData) IniDataPembelian() (map[string]any, int64) {
	transactionDetailsContent := map[string]any{}
	transactionDetailsContent["order_id"] = m.newPembelian.ID
	payload := map[string]any{}
	if m.newPembelian.JumlahDL > 0 && m.newPembelian.JumlahDL < 100 {
		var Items = []midtrans.ItemDetails{
			{
				ID:    m.newPembelian.ID,
				Price: int64(m.harga.HargaBeliDL) ,
				Qty:   int32(m.newPembelian.JumlahDL),
				Name:  "Item DL",
			},
		}
		payload["item_details"] = Items
		transactionDetailsContent["gross_amount"] = (Items[0].Price * int64(Items[0].Qty)) 
	} else if  m.newPembelian.JumlahDL % 100 == 0 && m.newPembelian.JumlahDL > 0 {
		var Items = []midtrans.ItemDetails{
			{
				ID:    m.newPembelian.ID,
				Price: int64(m.harga.HargaBeliBGL) ,
				Qty:   int32(m.newPembelian.JumlahDL) / 100,
				Name:  "Item BGL",
			},
		}
		payload["item_details"] = Items
		transactionDetailsContent["gross_amount"] = (Items[0].Price * int64(Items[0].Qty)) 
	} else if m.newPembelian.JumlahDL > 100 {
		var Items = []midtrans.ItemDetails{
			{
				ID:    m.newPembelian.ID,
				Price: int64(m.harga.HargaBeliBGL) ,
				Qty:   int32(m.newPembelian.JumlahDL / 100),
				Name:  "Item BGL",
			},
			{
				ID:    m.newPembelian.ID,
				Price: int64(m.harga.HargaBeliDL) ,
				Qty:   int32(m.newPembelian.JumlahDL % 100),
				Name:  "Item DL",
			},
		}
		payload["item_details"] = Items
		transactionDetailsContent["gross_amount"] = (Items[0].Price * int64(Items[0].Qty)) + (Items[1].Price * int64(Items[1].Qty))
	}

	customerDetails := map[string]any{}
	customerDetails["first_name"] = m.newPembelian.Nama
	customerDetails["phone"] = m.newPembelian.WA

	gopayContent := map[string]any{}
	shopeepayContent := map[string]any{}
	transferBankContent := map[string]any{}
	qrisContent := map[string]string{}
	if m.typePayment == "gopay" {
		gopayContent["enable_callback"] = true
		gopayContent["callback_url"] = "https://dlcheap.com"
	} else if m.typePayment == "shopeepay" {
		shopeepayContent["enable_callback"] = true
		shopeepayContent["callback_url"] = "https://dlcheap.com"
	} else if m.typePayment == "bank_transfer" {
		transferBankContent["bank"] = m.jenisBank
	} else if m.typePayment == "qris" {
		qrisContent["acquirer"] = "gopay"
	}


	payload["payment_type"] = m.typePayment
	payload["transaction_details"] = transactionDetailsContent
	
	payload["customer_details"] = customerDetails
	if len(gopayContent) != 0 {
		payload["gopay"] = gopayContent
	} else if len(shopeepayContent) != 0 {
		payload["shopeepay"] = shopeepayContent
	} else if len(transferBankContent) != 0 {
		payload["bank_transfer"] = transferBankContent
	} else {
		payload["qris"] = qrisContent
	}
	return payload, transactionDetailsContent["gross_amount"].(int64)
}
