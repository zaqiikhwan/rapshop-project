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
