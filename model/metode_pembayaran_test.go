package model

import (
	"encoding/json"
	"rapsshop-project/entities"
	"strings"
	"testing"
)

func TestNewCheckoutPaymentOptions(t *testing.T) {
	index := 7
	options := NewCheckoutPaymentOptions([]entities.MetodePembayaran{
		{
			IndexPembayaran:      &index,
			JenisPembayaran:      "BCA Manual",
			KredensialPembayaran: "secret-account-number",
			Pemilik:              "RapsShop",
		},
	}, true)

	if len(options.Gateway) != 6 {
		t.Fatalf("gateway options got %d, want 6", len(options.Gateway))
	}
	if !options.Gateway[0].Enabled {
		t.Fatal("gateway option enabled got false, want true")
	}
	if options.Gateway[0].RequiresPaymentProof {
		t.Fatal("gateway requires proof got true, want false")
	}
	if len(options.Manual) != 1 {
		t.Fatalf("manual options got %d, want 1", len(options.Manual))
	}
	manual := options.Manual[0]
	if manual.IndexPembayaran != index {
		t.Fatalf("manual index got %d, want %d", manual.IndexPembayaran, index)
	}
	if manual.JenisPembayaran != "BCA Manual" {
		t.Fatalf("manual jenis got %q, want BCA Manual", manual.JenisPembayaran)
	}
	if manual.Provider != "manual_transfer" {
		t.Fatalf("manual provider got %q, want manual_transfer", manual.Provider)
	}
	if !manual.RequiresPaymentProof {
		t.Fatal("manual requires proof got false, want true")
	}
	encoded, err := json.Marshal(options)
	if err != nil {
		t.Fatalf("failed marshal options: %v", err)
	}
	if strings.Contains(string(encoded), "secret-account-number") {
		t.Fatal("checkout options exposed manual payment credentials")
	}
}

func TestNewCheckoutPaymentOptionsSkipsManualMethodsWithoutIndex(t *testing.T) {
	options := NewCheckoutPaymentOptions([]entities.MetodePembayaran{
		{JenisPembayaran: "Broken Manual"},
	}, false)

	if len(options.Manual) != 0 {
		t.Fatalf("manual options got %d, want 0", len(options.Manual))
	}
	if options.Gateway[0].Enabled {
		t.Fatal("gateway enabled got true, want false")
	}
}

func TestNewPembelianManualCreateResponseDoesNotExposePaymentCredentials(t *testing.T) {
	index := 7
	payment := NewManualPaymentInstruction(entities.MetodePembayaran{
		IndexPembayaran:      &index,
		JenisPembayaran:      "BCA Manual",
		KredensialPembayaran: "secret-account-number",
		Pemilik:              "RapsShop",
	})
	response := NewPembelianManualCreateResponse("order-123", payment)

	if response.IDTransaksi != "order-123" {
		t.Fatalf("id_transaksi got %q, want order-123", response.IDTransaksi)
	}
	if response.Payment.IndexPembayaran != index {
		t.Fatalf("payment index got %d, want %d", response.Payment.IndexPembayaran, index)
	}
	if response.Payment.JenisPembayaran != "BCA Manual" {
		t.Fatalf("payment jenis got %q, want BCA Manual", response.Payment.JenisPembayaran)
	}
	if response.UploadProofPath != "/api/v1/upload/order-123" {
		t.Fatalf("upload proof path got %q", response.UploadProofPath)
	}
	if response.TrackingPath != "/api/v1/pembelian/order-123/tracking" {
		t.Fatalf("tracking path got %q", response.TrackingPath)
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed marshal manual create response: %v", err)
	}
	body := string(encoded)
	if strings.Contains(body, "kredensial_pembayaran") {
		t.Fatal("manual create response exposed credential field")
	}
	if strings.Contains(body, "secret-account-number") {
		t.Fatal("manual create response exposed credential value")
	}
}
