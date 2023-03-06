package main

import (
	"log"
	"net/http"
	"os"
	"rapsshop-project/database/mysql"
	"rapsshop-project/middleware"
	
	// admin 
	adminHandler "rapsshop-project/src/admin/handlers"
	adminRepo "rapsshop-project/src/admin/repo"
	adminUsecase "rapsshop-project/src/admin/service"

	// pembelian
	pembelianDLHandler "rapsshop-project/src/pembelian_dl/handlers"

	// stock_dl
	stockDLRepo "rapsshop-project/src/stock_dl/repo"
	stockDLUsecase "rapsshop-project/src/stock_dl/service"
	stockDLHandler "rapsshop-project/src/stock_dl/handlers"

	// harga_dl
	hargaDLRepo "rapsshop-project/src/harga_dl/repo"
	hargaDLUsecase "rapsshop-project/src/harga_dl/service"
	hargaDLHandler "rapsshop-project/src/harga_dl/handlers"

	// testimoni
	testiHandler "rapsshop-project/src/testimoni/handlers"
	testiRepo "rapsshop-project/src/testimoni/repo"
	testiUsecase "rapsshop-project/src/testimoni/service"

	// sosmed
	sosmedHandler "rapsshop-project/src/sosmed/handlers"
	sosmedRepo "rapsshop-project/src/sosmed/repo"
	sosmedUsecase "rapsshop-project/src/sosmed/service"

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

	if db == nil {
		log.Fatal("failed to connect database\n")
	}
	

	// uncomment for change to release mode
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()
	
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


	pembelianDLHandler.NewPembelianHandler(api)

	r.Run()
}