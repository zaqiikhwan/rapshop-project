package handlers

import (
	"net/http"
	"rapsshop-project/model"
	"rapsshop-project/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type envGrowtopiaHandler struct {
	EnvGrowtopiaUsecase model.GrowtopiaEnvUsecase
}

func NewEnvGrowtopiaHandler(r *gin.RouterGroup, egh model.GrowtopiaEnvUsecase, jwtMiddleware gin.HandlerFunc) {
	envGrowtopiaHandler := &envGrowtopiaHandler{EnvGrowtopiaUsecase: egh}
	r.POST("/env", jwtMiddleware, envGrowtopiaHandler.CreateNewEnv)
	r.GET("/env", jwtMiddleware, envGrowtopiaHandler.GetLatestEnv)
	r.PATCH("/env", jwtMiddleware, envGrowtopiaHandler.UpdateLatestEnv)
}

func (egh *envGrowtopiaHandler) CreateNewEnv(c *gin.Context) {
	var input model.InputGrowtopiaEnv

	if err := c.BindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "must bind with json", err)
		return
	}

	if err := egh.EnvGrowtopiaUsecase.CreateNewEnv(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when add new env growtopia", err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "success add new env growtopia", input)
}

func (egh *envGrowtopiaHandler) GetLatestEnv(c *gin.Context) {
	envGrowtopia, err := egh.EnvGrowtopiaUsecase.GetLatestEnv(); 
	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "env growtopia not found", err)
		return
	}
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed to fetch env growtopia data", err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "success add new env growtopia", envGrowtopia)
}

func (egh *envGrowtopiaHandler) UpdateLatestEnv(c *gin.Context) {
	var input model.InputGrowtopiaEnv

	if err := c.BindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "must bind with json", err)
		return
	}

	updatedEnv, err := egh.EnvGrowtopiaUsecase.UpdateLatestEnv(&input); 
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when add new env growtopia", err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "success add new env growtopia", updatedEnv)
}
