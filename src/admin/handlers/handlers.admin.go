package handlers

import (
	"errors"
	"net/http"
	"os"
	"rapsshop-project/model"
	"rapsshop-project/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct{
	model.AdminUsecase
}

func NewAdminHandler(r *gin.RouterGroup, au model.AdminUsecase, jwtMiddleware gin.HandlerFunc) {
	adminHandler := &AdminHandler{AdminUsecase: au}
	api := r.Group("/admin") 
	{
		api.POST("/register", adminHandler.RegisterAdmin)
		api.POST("/login", adminHandler.LoginAdmin)
	}
	r.GET("/profile", jwtMiddleware, adminHandler.GetProfile)
}

func (ah *AdminHandler) RegisterAdmin(c *gin.Context) {
	var input model.NewAdmin
	
	if err := c.BindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "input not binding with json", err)
		return
	}
	if input.Token != os.Getenv("TOKEN") {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "token not match", errors.New("token is not valid with the key"))
		return 
	}

	if err := ah.Register(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed register new admin", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "success register new admin", input.Nama)
}

func (ah *AdminHandler) LoginAdmin(c *gin.Context) {
	var input model.AdminLogin

	if err := c.BindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "input not binding with json", err)
		return
	}

	token, err := ah.Login(&input)
	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "username is not exist", err)
		return
	}

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "username or password was wrong", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "login successful", map[string]string{"token": token})
}

func (ah *AdminHandler) GetProfile(c *gin.Context) {
	id := c.MustGet("id").(string)

	if id == "" {
		utils.FailureOrErrorResponse(c, http.StatusUnauthorized, "please, login first", errors.New("credential is not available"))
		return
	}

	admin, err := ah.AdminUsecase.Profile(string(id))

	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "admin is not exist", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "successfully get data admin", admin)
}