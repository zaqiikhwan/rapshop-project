package main

import (
	"log"
	"net/http"
	"rapsshop-project/admin/handlers"
	"rapsshop-project/admin/repo"
	"rapsshop-project/admin/usecase"
	"rapsshop-project/database/mysql"
	"rapsshop-project/middleware"
	testiRepo "rapsshop-project/testimoni/repo"
	testiUsecase "rapsshop-project/testimoni/usecase"
	testiHandler "rapsshop-project/testimoni/handlers"

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
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// health check route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api/v1")

	adminRepo := repo.NewAdminRepository(db)
	adminUsecase := usecase.NewAdminUsecase(adminRepo)
	handlers.NewAdminHandler(api, adminUsecase, jwtMiddleware)

	testiRepo := testiRepo.NewTestimoniRepository(db)
	testiUsecase := testiUsecase.NewTestimoniUsecase(testiRepo)
	testiHandler.NewTestimoniHandler(api, testiUsecase, jwtMiddleware)

	r.Run()
}