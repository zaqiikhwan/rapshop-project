package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

	// payment method
	pmHandler "rapsshop-project/src/metode_pembayaran/handlers"
	pmRepo "rapsshop-project/src/metode_pembayaran/repo"
	pmUsecase "rapsshop-project/src/metode_pembayaran/service"

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
	allowedOrigins := map[string]bool{}
	for _, o := range strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",") {
		if o = strings.TrimSpace(o); o != "" {
			allowedOrigins[o] = true
		}
	}
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
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
	stockDLHandler.NewStockDLHandler(api, stockDLUsecase, jwtMiddleware)

	testiRepo := testiRepo.NewTestimoniRepository(db)
	testiUsecase := testiUsecase.NewTestimoniUsecase(testiRepo)
	testiHandler.NewTestimoniHandler(api, testiUsecase, jwtMiddleware)

	sosmedRepo := sosmedRepo.NewSosmedRepository(db)
	sosmedUsecase := sosmedUsecase.NewSosmedUsecase(sosmedRepo)
	sosmedHandler.NewSosmedHandler(api, sosmedUsecase, jwtMiddleware)

	hargaDLRepo := hargaDLRepo.NewHargaDLRepository(db)
	hargaDLUsecase := hargaDLUsecase.NewHargaDLUsecase(hargaDLRepo)
	hargaDLHandler.NewHargaDLHandler(api, hargaDLUsecase, jwtMiddleware)

	envGrowtopiaRepo := envGrowtopiaRepo.NewEnvGrowtopiaRepo(db)
	envGrowtopiaUsecase := envGrowtopiaUsecase.NewEnvGrowtopiaUsecase(envGrowtopiaRepo)
	envGrowtopiaHandler.NewEnvGrowtopiaHandler(api, envGrowtopiaUsecase, jwtMiddleware)

	jualDLRepo := jualDLRepo.NewPenjualanDLRepository(db)
	jualDLUsecase := jualDLUsecase.NewPenjualanDLUsecase(db, jualDLRepo, stockDLRepo)
	jualDLHandler.NewPenjualanDLHandler(api, jualDLUsecase, adminRepo, stockDLUsecase, jwtMiddleware)

	paymentMethodRepo := pmRepo.NewRepoMetodePembayaran(db)
	paymentMethodUsecase := pmUsecase.NewMetodePembayaranUsecase(paymentMethodRepo)
	pmHandler.NewMetodePembayaranHandler(api, paymentMethodUsecase, jwtMiddleware)

	pembelianDLRepo := pembelianDLRepo.NewRepoPembelianDL(db)
	pembelianDLUsecase := pembelianDLUsecase.NewServicePembelianDL(db, pembelianDLRepo, &midtransDriver, stockDLRepo)
	pembelianDLHandler.NewPembelianHandler(api, pembelianDLUsecase, adminRepo, paymentMethodUsecase, stockDLUsecase, jwtMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen failed: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %s\n", err)
	}
}
