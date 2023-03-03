package main

import (
	"log"
	"net/http"
	"os"
	"rapsshop-project/database/mysql"
	"rapsshop-project/middleware"
	
	adminHandler "rapsshop-project/src/admin/handlers"
	adminRepo "rapsshop-project/src/admin/repo"
	adminUsecase "rapsshop-project/src/admin/service"

	testiHandler "rapsshop-project/src/testimoni/handlers"
	testiRepo "rapsshop-project/src/testimoni/repo"
	testiUsecase "rapsshop-project/src/testimoni/service"

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

	testiRepo := testiRepo.NewTestimoniRepository(db)
	testiUsecase := testiUsecase.NewTestimoniUsecase(testiRepo)
	testiHandler.NewTestimoniHandler(api, testiUsecase, jwtMiddleware)

	sosmedRepo := sosmedRepo.NewSosmedRepository(db)
	sosmedUsecase := sosmedUsecase.NewSosmedUsecase(sosmedRepo)
	sosmedHandler.NewAdminHandler(api, sosmedUsecase, jwtMiddleware)

	r.Run()
}