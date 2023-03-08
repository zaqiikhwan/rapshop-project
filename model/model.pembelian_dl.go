package model

import (
	"fmt"
	"rapsshop-project/entities"

	"github.com/google/uuid"
)

type PembelianDLRepository interface {
	Create(input entities.PembelianDL) error
}

type MidtransData struct {
	typePayment string 
	jenisBank string
	newPembelian entities.PembelianDL
	harga entities.HargaDL
}

func NewMidtransData(typePayment string, jenisBank string, newPembelian entities.PembelianDL, harga entities.HargaDL) *MidtransData {
	return &MidtransData{
		typePayment: typePayment,
		jenisBank: jenisBank,
		newPembelian: newPembelian,
		harga: harga,
	}
}

func (m *MidtransData) IniDataPembelian() map[string]any {
	transactionDetailsContent := map[string]any{}
	transactionDetailsContent["order_id"] = fmt.Sprintf("order %v", uuid.New())
	
	listData := []map[string]any{}
	itemDetailsContent := map[string]any{}

	itemDetailsContent["id"] = fmt.Sprintf("order ID: %v", m.newPembelian.ID)
	// itemDetailsContent["price"] = m.harga.HargaBeliDL
	// itemDetailsContent["price"] = m.harga.HargaBeliDL
	// itemDetailsContent["quantity"] = m.newPembelian.JumlahDL
	// transactionDetailsContent["gross_amount"] = itemDetailsContent["price"].(int) * itemDetailsContent["quantity"].(int)
	// transactionDetailsContent["gross_amount"] = ((m.newPembelian.JumlahDL % 100) * m.harga.HargaBeliDL) + ((m.newPembelian.JumlahDL - (m.newPembelian.JumlahDL % 100)) * (m.harga.HargaBeliDL - 100))
	if m.newPembelian.JumlahDL % 100 == 0 && m.newPembelian.JumlahDL > 0 {
		itemDetailsContent["price"] = m.harga.HargaBeliBGL
		itemDetailsContent["quantity"] = m.newPembelian.JumlahDL / 100
		transactionDetailsContent["gross_amount"] = (m.newPembelian.JumlahDL / 100) * m.harga.HargaBeliBGL
	} else if m.newPembelian.JumlahDL > 100 {
		itemDetailsContent["quantity"] = m.newPembelian.JumlahDL
		itemDetailsContent["price"] = m.harga.HargaBeliDL
		transactionDetailsContent["gross_amount"] = ((m.newPembelian.JumlahDL % 100) * m.harga.HargaBeliDL) + ((m.newPembelian.JumlahDL - (m.newPembelian.JumlahDL % 100)) * (m.harga.HargaBeliDL - 100))
	} else if m.newPembelian.JumlahDL > 0 && m.newPembelian.JumlahDL < 100 {
		itemDetailsContent["quantity"] = m.newPembelian.JumlahDL
		itemDetailsContent["price"] = m.harga.HargaBeliDL
		transactionDetailsContent["gross_amount"] = m.newPembelian.JumlahDL * m.harga.HargaBeliDL
	}

	fmt.Println(itemDetailsContent["quantity"].(int))
	fmt.Println(itemDetailsContent["price"].(int))
	fmt.Println(transactionDetailsContent["gross_amount"].(int))
 	itemDetailsContent["name"] = m.newPembelian.Nama

	listData = append(listData, itemDetailsContent)

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

	payload := map[string]any{}
	payload["payment_type"] = m.typePayment
	payload["transaction_details"] = transactionDetailsContent
	payload["item_details"] = listData
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
	return payload
}
