package model

import (
	"errors"
	"net/url"
	"os"
	"rapsshop-project/entities"
	"strings"
	"time"

	"github.com/midtrans/midtrans-go"
	"gorm.io/gorm"
)

var (
	ErrInvalidJumlahDL   = errors.New("jumlah_dl must be greater than 0")
	ErrInsufficientStock = errors.New("insufficient stock")
)

const StatusPembayaranChallenge = "challange" //nolint:misspell // legacy persisted status value; do not change without migration.

type RekapTotalPembelian struct {
	Tanggal  string `json:"tanggal"`
	JumlahDL int    `json:"jumlah_dl"`
}

type PembelianDLRepository interface {
	Create(input entities.PembelianDL) error
	GetAll(_startInt int, _endInt int, queue string) ([]entities.PembelianDL, int, error)
	UpdateByID(input entities.PembelianDL, id string) error
	GetByID(id string) (entities.PembelianDL, error)
	GetTotalPembelian(date string) ([]RekapTotalPembelian, error)
	MarkPaid(id string) (bool, error)
	MarkShipped(id string, editor string) (bool, error)
	WithTx(tx *gorm.DB) PembelianDLRepository
}

type PembelianDLUsecase interface {
	// CreateDataPembelian(world string, nama string, grow_id string, jenis_item bool, jumlah_dl int, wa string, metode_transfer int, gambar string, id string) error
	GetTotal(date string) ([]RekapTotalPembelian, error)
	GetDetailByID(id string) (entities.PembelianDL, error)
	GetTrackingByID(id string) (PembelianTrackingResponse, error)
	GetCheckoutPreview(jumlahDL int) (PembelianCheckoutPreviewResponse, error)
	GetAllPembelian(_startInt int, _endInt int, queue string) ([]entities.PembelianDL, int, error)
	CreateDataPembelian(input entities.PembelianDL) error
	CreateDataPembelianManual(input entities.PembelianDL) error
	ChargeAndCreate(input entities.PembelianDL) (map[string]any, error)
	GetLiveStatus(id string) (map[string]any, error)
	UpdateStatusPembayaran(id string) error
	UpdateTambahBukti(id string, input entities.PembelianDL) error
	UpdateStatusPengiriman(id string, input entities.PembelianDL) error
	UpdateStatusButtonBayar(id string, input entities.PembelianDL) error
	UpdateStatusPembayaranAdmin(id string, input entities.PembelianDL) error
}

type PembelianSupportResponse struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
	Link    string `json:"link,omitempty"`
}

type PembelianTrackingResponse struct {
	ID                    string                   `json:"id"`
	World                 string                   `json:"world"`
	Nama                  string                   `json:"nama"`
	GrowID                string                   `json:"grow_id"`
	JenisItem             string                   `json:"jenis_item"`
	JumlahDL              int                      `json:"jumlah_dl"`
	MetodeTransfer        int                      `json:"metode_transfer"`
	JumlahTransaksi       int64                    `json:"jumlah_transaksi"`
	StatusPembayaran      string                   `json:"status_pembayaran"`
	StatusPembayaranLabel string                   `json:"status_pembayaran_label"`
	StatusPengiriman      bool                     `json:"status_pengiriman"`
	StatusPengirimanLabel string                   `json:"status_pengiriman_label"`
	BuktiPembayaran       string                   `json:"bukti_pembayaran"`
	Support               PembelianSupportResponse `json:"support"`
	CreatedAt             time.Time                `json:"created_at"`
	UpdatedAt             time.Time                `json:"updated_at"`
}

type PembelianCheckoutBreakdown struct {
	BGLQuantity int   `json:"bgl_quantity"`
	DLQuantity  int   `json:"dl_quantity"`
	BGLSubtotal int64 `json:"bgl_subtotal"`
	DLSubtotal  int64 `json:"dl_subtotal"`
}

type PembelianCheckoutPreviewResponse struct {
	JumlahDL      int                        `json:"jumlah_dl"`
	StockDL       int                        `json:"stock_dl"`
	StockEnough   bool                       `json:"stock_enough"`
	CanCheckout   bool                       `json:"can_checkout"`
	HargaBeliDL   int                        `json:"harga_beli_dl"`
	HargaBeliBGL  int                        `json:"harga_beli_bgl"`
	Breakdown     PembelianCheckoutBreakdown `json:"breakdown"`
	TotalPayment  int64                      `json:"total_payment"`
	FailureReason string                     `json:"failure_reason,omitempty"`
}

type PembelianManualCreateResponse struct {
	IDTransaksi     string                           `json:"id_transaksi"`
	Payment         ManualPaymentInstructionResponse `json:"payment"`
	UploadProofPath string                           `json:"upload_proof_path"`
	TrackingPath    string                           `json:"tracking_path"`
}

func NewPembelianManualCreateResponse(id string, payment ManualPaymentInstructionResponse) PembelianManualCreateResponse {
	return PembelianManualCreateResponse{
		IDTransaksi:     id,
		Payment:         payment,
		UploadProofPath: "/api/v1/upload/" + id,
		TrackingPath:    "/api/v1/pembelian/" + id + "/tracking",
	}
}

func NewPembelianCheckoutPreview(jumlahDL int, stock entities.StockDL) (PembelianCheckoutPreviewResponse, error) {
	if jumlahDL <= 0 {
		return PembelianCheckoutPreviewResponse{}, ErrInvalidJumlahDL
	}

	bglQuantity := jumlahDL / 100
	dlQuantity := jumlahDL % 100
	bglSubtotal := int64(bglQuantity * stock.HargaBeliBGL)
	dlSubtotal := int64(dlQuantity * stock.HargaBeliDL)
	stockEnough := jumlahDL <= stock.StockDL

	response := PembelianCheckoutPreviewResponse{
		JumlahDL:     jumlahDL,
		StockDL:      stock.StockDL,
		StockEnough:  stockEnough,
		CanCheckout:  stockEnough,
		HargaBeliDL:  stock.HargaBeliDL,
		HargaBeliBGL: stock.HargaBeliBGL,
		Breakdown: PembelianCheckoutBreakdown{
			BGLQuantity: bglQuantity,
			DLQuantity:  dlQuantity,
			BGLSubtotal: bglSubtotal,
			DLSubtotal:  dlSubtotal,
		},
		TotalPayment: bglSubtotal + dlSubtotal,
	}
	if !stockEnough {
		response.FailureReason = ErrInsufficientStock.Error()
	}
	return response, nil
}

func NewPembelianTrackingResponse(data entities.PembelianDL) PembelianTrackingResponse {
	statusPengiriman := false
	if data.StatusPengiriman != nil {
		statusPengiriman = *data.StatusPengiriman
	}

	jenisItem := "DL"
	if data.JenisItem {
		jenisItem = "BGL"
	}

	return PembelianTrackingResponse{
		ID:                    data.ID,
		World:                 data.World,
		Nama:                  data.Nama,
		GrowID:                data.GrowID,
		JenisItem:             jenisItem,
		JumlahDL:              data.JumlahDL,
		MetodeTransfer:        data.MetodeTransfer,
		JumlahTransaksi:       data.JumlahTransaksi,
		StatusPembayaran:      data.StatusPembayaran,
		StatusPembayaranLabel: PaymentStatusLabel(data.StatusPembayaran, data.BuktiPembayaran),
		StatusPengiriman:      statusPengiriman,
		StatusPengirimanLabel: DeliveryStatusLabel(statusPengiriman),
		BuktiPembayaran:       data.BuktiPembayaran,
		Support:               NewPembelianSupportResponse(data.ID),
		CreatedAt:             data.CreatedAt,
		UpdatedAt:             data.UpdatedAt,
	}
}

func NewPembelianSupportResponse(orderID string) PembelianSupportResponse {
	message := "Halo admin, saya ingin bertanya tentang order " + orderID
	response := PembelianSupportResponse{
		Channel: "whatsapp",
		Message: message,
	}

	number := whatsappNumberDigits(os.Getenv("SUPPORT_WHATSAPP_NUMBER"))
	if number != "" {
		response.Link = "https://wa.me/" + number + "?text=" + url.QueryEscape(message)
	}

	return response
}

func whatsappNumberDigits(number string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, number)
}

func PaymentStatusLabel(status string, buktiPembayaran string) string {
	if status == "belum_dibayar" && buktiPembayaran == "" {
		return "Waiting for proof upload"
	}
	if status == "belum_dibayar" && buktiPembayaran != "" {
		return "Waiting for admin confirmation"
	}

	switch status {
	case "pending":
		return "Payment pending"
	case "success", "dibayar":
		return "Payment confirmed"
	case "deny":
		return "Payment denied"
	case "failure":
		return "Payment failed"
	case StatusPembayaranChallenge:
		return "Under review"
	default:
		return "Waiting for payment"
	}
}

func DeliveryStatusLabel(statusPengiriman bool) string {
	if statusPengiriman {
		return "Delivered"
	}
	return "Waiting for delivery"
}

type MidtransData struct {
	typePayment  string
	jenisBank    string
	newPembelian entities.PembelianDL
	harga        entities.StockDL
}

func NewMidtransData(typePayment string, jenisBank string, newPembelian entities.PembelianDL, harga entities.StockDL) *MidtransData {
	return &MidtransData{
		typePayment:  typePayment,
		jenisBank:    jenisBank,
		newPembelian: newPembelian,
		harga:        harga,
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
				Price: int64(m.harga.HargaBeliDL),
				Qty:   int32(m.newPembelian.JumlahDL),
				Name:  "Item DL",
			},
		}
		payload["item_details"] = Items
		transactionDetailsContent["gross_amount"] = (Items[0].Price * int64(Items[0].Qty))
	} else if m.newPembelian.JumlahDL%100 == 0 && m.newPembelian.JumlahDL > 0 {
		var Items = []midtrans.ItemDetails{
			{
				ID:    m.newPembelian.ID,
				Price: int64(m.harga.HargaBeliBGL),
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
				Price: int64(m.harga.HargaBeliBGL),
				Qty:   int32(m.newPembelian.JumlahDL / 100),
				Name:  "Item BGL",
			},
			{
				ID:    m.newPembelian.ID,
				Price: int64(m.harga.HargaBeliDL),
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
