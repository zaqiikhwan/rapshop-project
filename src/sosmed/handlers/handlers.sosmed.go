package handlers

import (
	"net/http"
	"rapsshop-project/model"
	"rapsshop-project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type sosmedHandler struct {
	SosmedUsecase model.SosmedUsecase
}

func NewAdminHandler(r *gin.RouterGroup, su model.SosmedUsecase, jwtMiddleware gin.HandlerFunc) {
	sosmedHandler := &sosmedHandler{SosmedUsecase: su}
	r.POST("/platform", jwtMiddleware, sosmedHandler.CreateNewSosmed)
	r.GET("/platforms", sosmedHandler.GetAllSosmed)
	r.GET("/platform/:id", jwtMiddleware, sosmedHandler.GetDetailSosmedByID)
	r.PATCH("/platform/:id", jwtMiddleware, sosmedHandler.UpdateSosmedByID)
	r.DELETE("/platform/:id", jwtMiddleware, sosmedHandler.DeleteSosmedByID)
}

func (sh *sosmedHandler) CreateNewSosmed(c *gin.Context) {
	var input model.InputSosmed

	if err := c.BindJSON(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "must bind with json", err)
		return
	}

	if err := sh.SosmedUsecase.CreateSosmed(&input); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed add new social media platform", err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "success add new social media platform", input.Username)
}

func (sh *sosmedHandler) GetAllSosmed(c *gin.Context) {
	allSosmed, err := sh.SosmedUsecase.GetAllSosmed()

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed when query all data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success fetch all platform data", allSosmed)
}

func (sh *sosmedHandler) GetDetailSosmedByID(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "error convert string to uint", err)
		return
	}

	sosmed, err := sh.SosmedUsecase.GetSosmedByID(uint(idUint))

	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "sosmed detail not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "sosmed found", sosmed)
}

func(sh *sosmedHandler) UpdateSosmedByID(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "error convert string to uint", err)
		return
	}

	var updateSosmed model.InputSosmed

	if err := c.BindJSON(&updateSosmed); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "must bind with json", err)
		return
	}

	result, err := sh.SosmedUsecase.UpdateSosmedByID(uint(idUint), &updateSosmed)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed to update sosmed data", err)
		return
	}

	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "sosmed not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success update data medsos", result)
}

func(sh *sosmedHandler) DeleteSosmedByID(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "error convert string to uint", err)
		return
	}

	if err := sh.SosmedUsecase.DeleteSosmedByID(uint(idUint)); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed to delete data sosmed", err)
		return	
	}

	utils.SuccessResponse(c, http.StatusOK, "success delete data sosmed", nil)
}