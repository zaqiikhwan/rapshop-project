package repo

import (
	"rapsshop-project/entities"

	"gorm.io/gorm"
)

type repoMetodePembayaran struct {
	db *gorm.DB
}

func NewRepoMetodePembayaran(db *gorm.DB) entities.MetodePembayaranRepository {
	return &repoMetodePembayaran{db: db}
}

func (rmp *repoMetodePembayaran) Create(newMetode entities.MetodePembayaran) error {
	if err := rmp.db.Create(&newMetode).Error; err != nil {
		return err
	}
	return nil
}

func (rmp *repoMetodePembayaran) GetAll() ([]entities.MetodePembayaran, error) {
	var allMethod []entities.MetodePembayaran

	if err := rmp.db.Find(&allMethod).Error; err != nil {
		return allMethod, err
	}

	return allMethod, nil
}

func (rmp *repoMetodePembayaran) GetByIndex(index int) (entities.MetodePembayaran, error) {
	var detailMethod entities.MetodePembayaran
	if err := rmp.db.Where("index_pembayaran = ?", index).Take(&detailMethod).Error; err != nil {
		return detailMethod, err
	}
	return detailMethod, nil
}


func (rmp *repoMetodePembayaran) GetByJenis(jenis string) (entities.MetodePembayaran, error) {
	var detailMethod entities.MetodePembayaran
	if err := rmp.db.Where("jenis_pembayaran = ?", jenis).Take(&detailMethod).Error; err != nil {
		return detailMethod, err
	}
	return detailMethod, nil
}

func (rmp *repoMetodePembayaran) UpdateKredensial(jenis string, patchKredensial entities.MetodePembayaran) error {
	if err := rmp.db.Where("jenis_pembayaran = ?", jenis).Updates(patchKredensial).Error; err != nil {
		return err
	}
	return nil
}