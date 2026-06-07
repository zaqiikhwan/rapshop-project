package model

import "rapsshop-project/entities"

type InputMetodePembayaran struct {
	IndexPembayaran      int    `json:"index_pembayaran"`
	JenisPembayaran      string `json:"jenis_pembayaran"`
	KredensialPembayaran string `json:"kredensial_pembayaran"`
	Pemilik              string `json:"pemilik"`
}

type CheckoutPaymentOptionResponse struct {
	IndexPembayaran      int    `json:"index_pembayaran"`
	JenisPembayaran      string `json:"jenis_pembayaran"`
	CheckoutType         string `json:"checkout_type"`
	Provider             string `json:"provider"`
	Bank                 string `json:"bank,omitempty"`
	Pemilik              string `json:"pemilik,omitempty"`
	RequiresPaymentProof bool   `json:"requires_payment_proof"`
	Enabled              bool   `json:"enabled"`
}

type CheckoutPaymentOptionsResponse struct {
	Gateway []CheckoutPaymentOptionResponse `json:"gateway"`
	Manual  []CheckoutPaymentOptionResponse `json:"manual"`
}

type ManualPaymentInstructionResponse struct {
	IndexPembayaran      int    `json:"index_pembayaran"`
	JenisPembayaran      string `json:"jenis_pembayaran"`
	CheckoutType         string `json:"checkout_type"`
	Provider             string `json:"provider"`
	Pemilik              string `json:"pemilik,omitempty"`
	RequiresPaymentProof bool   `json:"requires_payment_proof"`
}

type MetodePembayaranRepository interface {
	Create(newMetode entities.MetodePembayaran) error
	GetAll() ([]entities.MetodePembayaran, error)
	GetByID(id uint) (entities.MetodePembayaran, error)
	GetByIndex(index int) (entities.MetodePembayaran, error)
	UpdateKredensialByID(id uint, patchKredensial entities.MetodePembayaran) error
	DeleteByID(id uint) error
}

type MetodePembayaranUsecase interface {
	CreateNewPembayaran(input *InputMetodePembayaran) error
	GetAllPembayaran() ([]entities.MetodePembayaran, error)
	GetCheckoutOptions() (CheckoutPaymentOptionsResponse, error)
	GetDetailPembayaranByIndex(index int) (entities.MetodePembayaran, error)
	GetDetailPembayaranByID(id uint) (entities.MetodePembayaran, error)
	PatchDetailPembayaranByID(id uint, input *InputMetodePembayaran) error
	DeletePembayaranByID(id uint) error
}

func NewCheckoutPaymentOptions(manualMethods []entities.MetodePembayaran, gatewayEnabled bool) CheckoutPaymentOptionsResponse {
	options := CheckoutPaymentOptionsResponse{
		Gateway: []CheckoutPaymentOptionResponse{
			{IndexPembayaran: 1, JenisPembayaran: "QRIS", CheckoutType: "gateway", Provider: "midtrans", Enabled: gatewayEnabled},
			{IndexPembayaran: 2, JenisPembayaran: "GoPay", CheckoutType: "gateway", Provider: "midtrans", Enabled: gatewayEnabled},
			{IndexPembayaran: 3, JenisPembayaran: "ShopeePay", CheckoutType: "gateway", Provider: "midtrans", Enabled: gatewayEnabled},
			{IndexPembayaran: 4, JenisPembayaran: "BCA Virtual Account", CheckoutType: "gateway", Provider: "midtrans", Bank: "bca", Enabled: gatewayEnabled},
			{IndexPembayaran: 5, JenisPembayaran: "BRI Virtual Account", CheckoutType: "gateway", Provider: "midtrans", Bank: "bri", Enabled: gatewayEnabled},
			{IndexPembayaran: 6, JenisPembayaran: "BNI Virtual Account", CheckoutType: "gateway", Provider: "midtrans", Bank: "bni", Enabled: gatewayEnabled},
		},
		Manual: []CheckoutPaymentOptionResponse{},
	}

	for _, method := range manualMethods {
		if method.IndexPembayaran == nil {
			continue
		}
		options.Manual = append(options.Manual, CheckoutPaymentOptionResponse{
			IndexPembayaran:      *method.IndexPembayaran,
			JenisPembayaran:      method.JenisPembayaran,
			CheckoutType:         "manual",
			Provider:             "manual_transfer",
			Pemilik:              method.Pemilik,
			RequiresPaymentProof: true,
			Enabled:              true,
		})
	}

	return options
}

func NewManualPaymentInstruction(method entities.MetodePembayaran) ManualPaymentInstructionResponse {
	index := 0
	if method.IndexPembayaran != nil {
		index = *method.IndexPembayaran
	}

	return ManualPaymentInstructionResponse{
		IndexPembayaran:      index,
		JenisPembayaran:      method.JenisPembayaran,
		CheckoutType:         "manual",
		Provider:             "manual_transfer",
		Pemilik:              method.Pemilik,
		RequiresPaymentProof: true,
	}
}
