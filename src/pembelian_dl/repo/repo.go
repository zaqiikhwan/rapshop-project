package repo

import (
	"rapsshop-project/entities"
	"rapsshop-project/model"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type repoPembelianDL struct {
	db *gorm.DB
}

func NewRepoPembelianDL(db *gorm.DB) model.PembelianDLRepository {
	return &repoPembelianDL{db: db}
}

func (rp *repoPembelianDL) WithTx(tx *gorm.DB) model.PembelianDLRepository {
	return &repoPembelianDL{db: tx}
}

// MarkPaid flips status_pembayaran to "success" only if it isn't already, returning
// true when this call performed the transition (enables exactly-once stock decrement).
func (rp *repoPembelianDL) MarkPaid(id string) (bool, error) {
	res := rp.db.Model(&entities.PembelianDL{}).
		Where("id = ? AND status_pembayaran <> ?", id, "success").
		Update("status_pembayaran", "success")
	return res.RowsAffected == 1, res.Error
}

// MarkShipped flips status_pengiriman to true only on the not-yet-shipped -> shipped
// transition, returning true when this call performed it.
func (rp *repoPembelianDL) MarkShipped(id string, editor string) (bool, error) {
	shipped := true
	res := rp.db.Model(&entities.PembelianDL{}).
		Where("id = ? AND (status_pengiriman = ? OR status_pengiriman IS NULL)", id, false).
		Updates(entities.PembelianDL{EditorStatus: editor, StatusPengiriman: &shipped})
	return res.RowsAffected == 1, res.Error
}

func (rp *repoPembelianDL) Create(input entities.PembelianDL) error {
	if err := rp.db.Create(&input).Error; err != nil {
		return err
	}
	return nil
}

func (rp *repoPembelianDL) GetAll(_startInt int, _endInt int, queue string) ([]entities.PembelianDL, int, error) {
	var allPembelian []entities.PembelianDL
	var total int64
	if _startInt < 1 {
		_startInt = 1
	}
	if _endInt < _startInt {
		_endInt = _startInt
	}

	query := rp.db.Model(&entities.PembelianDL{})
	switch queue {
	case "pending_payment":
		query = query.Where("(status_pembayaran = ? OR status_pembayaran = ?) AND (bukti_pembayaran = ? OR bukti_pembayaran IS NULL)", "belum_dibayar", "pending", "")
	case "proof_uploaded":
		query = query.Where("status_pembayaran = ? AND bukti_pembayaran <> ?", "belum_dibayar", "")
	case "waiting_delivery":
		query = query.Where("(status_pembayaran = ? OR status_pembayaran = ?) AND (status_pengiriman = ? OR status_pengiriman IS NULL)", "success", "dibayar", false)
	case "delivered":
		query = query.Where("status_pengiriman = ?", true)
	case "failed":
		query = query.Where("status_pembayaran = ? OR status_pembayaran = ?", "deny", "failure")
	case "review":
		query = query.Where("status_pembayaran = ?", model.StatusPembayaranChallenge)
	default:
		query = query.Where("status_pembayaran = ? OR status_pembayaran = ?", "success", "dibayar")
	}

	if err := query.Count(&total).Error; err != nil {
		return allPembelian, 0, err
	}
	if err := query.Order("created_at desc").Offset(_startInt - 1).Limit(_endInt - _startInt + 1).Find(&allPembelian).Error; err != nil {
		return allPembelian, 0, err
	}
	return allPembelian, int(total), nil
}

func (rp *repoPembelianDL) GetByID(id string) (entities.PembelianDL, error) {
	var model entities.PembelianDL
	if err := rp.db.Where("id = ?", id).Take(&model).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (rp *repoPembelianDL) UpdateByID(input entities.PembelianDL, id string) error {
	var statusPayment entities.PembelianDL
	if err := rp.db.Where("id = ?", id).Model(&statusPayment).Updates(input).Error; err != nil {
		return err
	}
	return nil
}

func (rp *repoPembelianDL) GetTotalPembelian(date string) ([]model.RekapTotalPembelian, error) {

	var totalPembelianTunggal model.RekapTotalPembelian
	var totalPembelian []model.RekapTotalPembelian

	arrayDate := strings.Split(date, "-")
	arrayYearInt, _ := strconv.Atoi(arrayDate[0])
	arrayMonthInt, _ := strconv.Atoi(arrayDate[1])

	var arrayTanggalStr []string
	if arrayYearInt%4 == 0 && arrayMonthInt == 2 {
		jumlahHari := 29
		arrayTanggalStr = GeneratorTanggal(arrayDate, jumlahHari)
	} else if arrayMonthInt == 2 {
		jumlahHari := 28
		arrayTanggalStr = GeneratorTanggal(arrayDate, jumlahHari)
	} else if (arrayMonthInt%2 != 0 && (arrayMonthInt != 9 && arrayMonthInt != 11)) || arrayMonthInt == 8 || arrayMonthInt == 10 || arrayMonthInt == 12 {
		jumlahHari := 31
		arrayTanggalStr = GeneratorTanggal(arrayDate, jumlahHari)
	} else {
		jumlahHari := 30
		arrayTanggalStr = GeneratorTanggal(arrayDate, jumlahHari)
	}

	for i := 0; i < len(arrayTanggalStr); i++ {

		queryDate := "%" + arrayTanggalStr[i] + "%"

		if err := rp.db.Raw("select sum(jumlah_dl) from pembelian_dls where created_at LIKE ? and (status_pembayaran = ? or status_pembayaran = ?)", queryDate, "success", "dibayar").Scan(&totalPembelianTunggal.JumlahDL).Error; err != nil {
			totalPembelianTunggal.JumlahDL = 0
		}

		totalPembelianTunggal.Tanggal = arrayTanggalStr[i][8:10]
		totalPembelian = append(totalPembelian, totalPembelianTunggal)
	}

	return totalPembelian, nil
}

func GeneratorTanggal(arrayDate []string, jumlahHari int) []string {
	var arrayTanggalStrCheck []string

	for i := 1; i <= jumlahHari; i++ {
		arrDate := [1]int{i} // {1, 2, 3, 4, 5, 6, 7, 8, 9, 10} // ini buat tanggal
		var arrDateStr [1]string
		var tglStr string
		arrDateStr = [1]string{strconv.Itoa(arrDate[0])}
		if arrDate[0] < 10 {
			combine := [3]string{"0", arrDateStr[0]}
			tglStr = strings.Join(combine[0:2], "")
		} else {
			tglStr = arrDateStr[0]
		}

		combine := [3]string{arrayDate[0], arrayDate[1], tglStr}
		arrayTanggalStrCheck = append(arrayTanggalStrCheck, strings.Join(combine[0:3], "-"))
	}
	return arrayTanggalStrCheck
}
