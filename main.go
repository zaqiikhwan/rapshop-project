package main

import (
	"log"
	"net/http"
	"os"
	"rapsshop-project/database/mysql"
	"rapsshop-project/lib"
	"rapsshop-project/middleware"

	// admin
	adminRepo "rapsshop-project/src/admin/repo"
	adminUsecase "rapsshop-project/src/admin/service"
	adminHandler "rapsshop-project/src/admin/handlers"

	// pembelian
	pembelianDLRepo "rapsshop-project/src/pembelian_dl/repo"
	pembelianDLUsecase "rapsshop-project/src/pembelian_dl/service"
	pembelianDLHandler "rapsshop-project/src/pembelian_dl/handlers"

	// penjualan dl
	jualDLRepo "rapsshop-project/src/penjualan_dl/repo"
	jualDLHandler "rapsshop-project/src/penjualan_dl/handlers"
	jualDLUsecase "rapsshop-project/src/penjualan_dl/service"

	// stock_dl
	stockDLRepo "rapsshop-project/src/stock_dl/repo"
	stockDLUsecase "rapsshop-project/src/stock_dl/service"
	stockDLHandler "rapsshop-project/src/stock_dl/handlers"

	// harga_dl
	hargaDLRepo "rapsshop-project/src/harga_dl/repo"
	hargaDLUsecase "rapsshop-project/src/harga_dl/service"
	hargaDLHandler "rapsshop-project/src/harga_dl/handlers"

	// testimoni
	testiRepo "rapsshop-project/src/testimoni/repo"
	testiUsecase "rapsshop-project/src/testimoni/service"
	testiHandler "rapsshop-project/src/testimoni/handlers"

	// sosmed
	sosmedRepo "rapsshop-project/src/sosmed/repo"
	sosmedUsecase "rapsshop-project/src/sosmed/service"
	sosmedHandler "rapsshop-project/src/sosmed/handlers"

	// env growtopia
	envGrowtopiaRepo "rapsshop-project/src/env_growtopia/repo"
	envGrowtopiaUsecase "rapsshop-project/src/env_growtopia/service"
	envGrowtopiaHandler "rapsshop-project/src/env_growtopia/handlers"

	// payment method
	pmRepo "rapsshop-project/src/metode_pembayaran/repo"
	pmUsecase "rapsshop-project/src/metode_pembayaran/service"
	pmHandler "rapsshop-project/src/metode_pembayaran/handlers"


	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env file\n")
	}
	db := mysql.InitDatabase()
	jwtMiddleware := middleware.NewAuthMiddleware()
	midtransDriver := lib.NewMidtransDriver()

	if db == nil {
		log.Fatal("failed to connect database\n")
	}
	

	// uncomment for change to release mode
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatus(204)
		} else {
			c.Next()
		}
	})
	
	// health check route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api/v1")

	adminRepo := adminRepo.NewAdminRepository(db)
	adminUsecase := adminUsecase.NewAdminUsecase(adminRepo)
	adminHandler.NewAdminHandler(api, adminUsecase, jwtMiddleware)

	stockDLRepo := stockDLRepo.NewStockDLRepository(db)
	stockDLUsecase := stockDLUsecase.NewStockDLUsecase(stockDLRepo)
	stockDLHandler.NewAdminHandler(api, stockDLUsecase, jwtMiddleware)

	testiRepo := testiRepo.NewTestimoniRepository(db)
	testiUsecase := testiUsecase.NewTestimoniUsecase(testiRepo)
	testiHandler.NewTestimoniHandler(api, testiUsecase, jwtMiddleware)

	sosmedRepo := sosmedRepo.NewSosmedRepository(db)
	sosmedUsecase := sosmedUsecase.NewSosmedUsecase(sosmedRepo)
	sosmedHandler.NewAdminHandler(api, sosmedUsecase, jwtMiddleware)

	hargaDLRepo := hargaDLRepo.NewHargaDLRepository(db)
	hargaDLUsecase := hargaDLUsecase.NewHargaDLUsecase(hargaDLRepo)
	hargaDLHandler.NewHargaDLHandler(api, hargaDLUsecase, jwtMiddleware)

	envGrowtopiaRepo := envGrowtopiaRepo.NewEnvGrowtopiaRepo(db)
	envGrowtopiaUsecase := envGrowtopiaUsecase.NewEnvGrowtopiaUsecase(envGrowtopiaRepo)
	envGrowtopiaHandler.NewEnvGrowtopiaHandler(api, envGrowtopiaUsecase, jwtMiddleware)

	jualDLRepo := jualDLRepo.NewPenjualanDLRepository(db)
	jualDLUsecase := jualDLUsecase.NewTestimoniUsecase(jualDLRepo, stockDLUsecase)
	jualDLHandler.NewPenjualanDLHandler(api, jualDLUsecase, adminRepo, jwtMiddleware)

	paymentMethodRepo := pmRepo.NewRepoMetodePembayaran(db)
	paymentMethodUsecase := pmUsecase.NewMetodePembayaranUsecase(paymentMethodRepo)
	pmHandler.NewMetodePembayaranHandler(api, paymentMethodUsecase, jwtMiddleware)

	pembelianDLRepo := pembelianDLRepo.NewRepoPembelianDL(db)
	pembelianDLUsecase := pembelianDLUsecase.NewServicePembelianDL(pembelianDLRepo, &midtransDriver, stockDLUsecase)
	pembelianDLHandler.NewPembelianHandler(api, pembelianDLUsecase, adminRepo, paymentMethodUsecase, jwtMiddleware)

	r.Run()
}