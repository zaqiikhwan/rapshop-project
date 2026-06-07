package model

import (
	"rapsshop-project/entities"
	"testing"
)

// Locks the agreed pricing rule:
//
//	1–99 DL  -> priced per DL (HargaBeliDL)
//	multiple of 100 -> priced per BGL (HargaBeliBGL), qty = DL/100
//	mixed   -> BGL portion + DL remainder
func TestIniDataPembelianGrossAmount(t *testing.T) {
	harga := entities.StockDL{HargaBeliDL: 100, HargaBeliBGL: 9000}

	cases := []struct {
		name   string
		jumlah int
		want   int64
	}{
		{"sub-100 DL", 50, 5000},     // 50 * 100
		{"exact BGL", 100, 9000},     // 1 * 9000
		{"mixed BGL+DL", 250, 23000}, // 2*9000 + 50*100
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			md := NewMidtransData("qris", "", entities.PembelianDL{ID: "x", JumlahDL: tc.jumlah}, harga)
			_, total := md.IniDataPembelian()
			if total != tc.want {
				t.Fatalf("jumlah=%d: got %d, want %d", tc.jumlah, total, tc.want)
			}
		})
	}
}

func TestNewPembelianCheckoutPreview(t *testing.T) {
	stock := entities.StockDL{StockDL: 500, HargaBeliDL: 100, HargaBeliBGL: 9000}

	cases := []struct {
		name        string
		jumlahDL    int
		wantBGL     int
		wantDL      int
		wantTotal   int64
		wantAllowed bool
		wantReason  string
	}{
		{"sub-100 DL", 50, 0, 50, 5000, true, ""},
		{"exact BGL", 100, 1, 0, 9000, true, ""},
		{"mixed BGL+DL", 250, 2, 50, 23000, true, ""},
		{"insufficient stock", 600, 6, 0, 54000, false, ErrInsufficientStock.Error()},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewPembelianCheckoutPreview(tc.jumlahDL, stock)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Breakdown.BGLQuantity != tc.wantBGL {
				t.Fatalf("bgl quantity got %d, want %d", got.Breakdown.BGLQuantity, tc.wantBGL)
			}
			if got.Breakdown.DLQuantity != tc.wantDL {
				t.Fatalf("dl quantity got %d, want %d", got.Breakdown.DLQuantity, tc.wantDL)
			}
			if got.TotalPayment != tc.wantTotal {
				t.Fatalf("total got %d, want %d", got.TotalPayment, tc.wantTotal)
			}
			if got.CanCheckout != tc.wantAllowed {
				t.Fatalf("can_checkout got %t, want %t", got.CanCheckout, tc.wantAllowed)
			}
			if got.FailureReason != tc.wantReason {
				t.Fatalf("failure_reason got %q, want %q", got.FailureReason, tc.wantReason)
			}
		})
	}
}

func TestNewPembelianCheckoutPreviewRejectsInvalidQuantity(t *testing.T) {
	_, err := NewPembelianCheckoutPreview(0, entities.StockDL{StockDL: 500})
	if err != ErrInvalidJumlahDL {
		t.Fatalf("got %v, want %v", err, ErrInvalidJumlahDL)
	}
}

func TestPaymentStatusLabel(t *testing.T) {
	cases := []struct {
		name   string
		status string
		proof  string
		want   string
	}{
		{"manual waiting upload", "belum_dibayar", "", "Waiting for proof upload"},
		{"manual waiting confirmation", "belum_dibayar", "https://example.com/proof.jpg", "Waiting for admin confirmation"},
		{"gateway pending", "pending", "", "Payment pending"},
		{"gateway success", "success", "", "Payment confirmed"},
		{"manual approved", "dibayar", "https://example.com/proof.jpg", "Payment confirmed"},
		{"gateway denied", "deny", "", "Payment denied"},
		{"gateway failed", "failure", "", "Payment failed"},
		{"gateway review", StatusPembayaranChallenge, "", "Under review"},
		{"unknown", "unexpected", "", "Waiting for payment"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := PaymentStatusLabel(tc.status, tc.proof)
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestNewPembelianTrackingResponse(t *testing.T) {
	t.Setenv("SUPPORT_WHATSAPP_NUMBER", "+62 812-3456-7890")

	shipped := true
	response := NewPembelianTrackingResponse(entities.PembelianDL{
		ID:               "order-1",
		JenisItem:        true,
		JumlahDL:         100,
		StatusPembayaran: "dibayar",
		StatusPengiriman: &shipped,
		BuktiPembayaran:  "https://example.com/proof.jpg",
		JumlahTransaksi:  9000,
		MetodeTransfer:   7,
		World:            "WORLD",
		Nama:             "Buyer",
		GrowID:           "GrowID",
	})

	if response.JenisItem != "BGL" {
		t.Fatalf("jenis_item got %q, want BGL", response.JenisItem)
	}
	if !response.StatusPengiriman {
		t.Fatal("status_pengiriman got false, want true")
	}
	if response.StatusPengirimanLabel != "Delivered" {
		t.Fatalf("status_pengiriman_label got %q, want Delivered", response.StatusPengirimanLabel)
	}
	if response.StatusPembayaranLabel != "Payment confirmed" {
		t.Fatalf("status_pembayaran_label got %q, want Payment confirmed", response.StatusPembayaranLabel)
	}
	if response.Support.Channel != "whatsapp" {
		t.Fatalf("support.channel got %q, want whatsapp", response.Support.Channel)
	}
	if response.Support.Message != "Halo admin, saya ingin bertanya tentang order order-1" {
		t.Fatalf("support.message got %q", response.Support.Message)
	}
	if response.Support.Link != "https://wa.me/6281234567890?text=Halo+admin%2C+saya+ingin+bertanya+tentang+order+order-1" {
		t.Fatalf("support.link got %q", response.Support.Link)
	}
}

func TestNewPembelianSupportResponseOmitsLinkWithoutConfiguredNumber(t *testing.T) {
	t.Setenv("SUPPORT_WHATSAPP_NUMBER", "")

	response := NewPembelianSupportResponse("order-1")

	if response.Link != "" {
		t.Fatalf("support.link got %q, want empty", response.Link)
	}
}
