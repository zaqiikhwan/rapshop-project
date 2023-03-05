package mysql

import (
	"fmt"
	"log"
	"os"
	"rapsshop-project/entities"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("init db failed, %s\n", err)
	}

	var admin entities.Admin
	var testimoni entities.Testimoni
	var sosmed entities.Sosmed
	var stock_dl entities.StockDL

	err = db.AutoMigrate(admin, testimoni, sosmed, stock_dl)
	if err != nil {
		log.Fatalf("failed to migrate, %s\n", err)
	}

	return db
}