package main

import (
	"log"
	"net/http"
	"os"
	"rapsshop-project/database/mysql"
	"rapsshop-project/lib"
	"rapsshop-project/middleware"

	// admin
	adminHandler "rapsshop-project/src/admin/handlers"
	adminRepo "rapsshop-project/src/admin/repo"
	adminUsecase "rapsshop-project/src/admin/service"

	// pembelian
	pembelianDLHandler "rapsshop-project/src/pembelian_dl/handlers"
	pembelianDLRepo "rapsshop-project/src/pembelian_dl/repo"
	pembelianDLUsecase "rapsshop-project/src/pembelian_dl/service"

	// penjualan dl
	jualDLHandler "rapsshop-project/src/penjualan_dl/handlers"
	jualDLRepo "rapsshop-project/src/penjualan_dl/repo"
	jualDLUsecase "rapsshop-project/src/penjualan_dl/service"

	// stock_dl
	stockDLHandler "rapsshop-project/src/stock_dl/handlers"
	stockDLRepo "rapsshop-project/src/stock_dl/repo"
	stockDLUsecase "rapsshop-project/src/stock_dl/service"

	// harga_dl
	hargaDLHandler "rapsshop-project/src/harga_dl/handlers"
	hargaDLRepo "rapsshop-project/src/harga_dl/repo"
	hargaDLUsecase "rapsshop-project/src/harga_dl/service"

	// testimoni
	testiHandler "rapsshop-project/src/testimoni/handlers"
	testiRepo "rapsshop-project/src/testimoni/repo"
	testiUsecase "rapsshop-project/src/testimoni/service"

	// sosmed
	sosmedHandler "rapsshop-project/src/sosmed/handlers"
	sosmedRepo "rapsshop-project/src/sosmed/repo"
	sosmedUsecase "rapsshop-project/src/sosmed/service"

	// env growtopia
	envGrowtopiaHandler "rapsshop-project/src/env_growtopia/handlers"
	envGrowtopiaRepo "rapsshop-project/src/env_growtopia/repo"
	envGrowtopiaUsecase "rapsshop-project/src/env_growtopia/service"

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
	jualDLUsecase := jualDLUsecase.NewTestimoniUsecase(jualDLRepo)
	jualDLHandler.NewPenjualanDLHandler(api, jualDLUsecase, jwtMiddleware)

	pembelianDLRepo := pembelianDLRepo.NewRepoPembelianDL(db)
	pembelianDLUsecase := pembelianDLUsecase.NewServicePembelianDL(pembelianDLRepo, &midtransDriver)
	pembelianDLHandler.NewPembelianHandler(api, pembelianDLUsecase)

	r.Run()
}