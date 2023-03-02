package handlers

import (
	"net/http"
	"rapsshop-project/model"
	"rapsshop-project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type testimoniHandler struct {
	TestimoniUsecase model.TestimoniUsecase
}

func NewTestimoniHandler(r *gin.RouterGroup, tu model.TestimoniUsecase, jwtMiddleware gin.HandlerFunc) {
	testimoniHandler := &testimoniHandler{TestimoniUsecase: tu}

	r.POST("/testimoni", jwtMiddleware, testimoniHandler.CreateTestimoni)
	r.GET("/testimonis", testimoniHandler.GetAllTestimoni)
	r.GET("/testimoni/:id", testimoniHandler.GetTestimoniByID)
	r.PATCH("/testimoni/:id", testimoniHandler.UpdateTestimoniByID)
	r.DELETE("/testimoni/:id", testimoniHandler.DeleteTestimoniByID)
}


func (th *testimoniHandler) CreateTestimoni(c *gin.Context) {
	testimoni := c.PostForm("testimoni")
	jumlahDLStr := c.PostForm("jumlah_dl")
	image, err := c.FormFile("gambar")
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	JumlahDL, err := strconv.Atoi(jumlahDLStr)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	if err := th.TestimoniUsecase.CreateTestimoni(image, testimoni, JumlahDL); err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "successfully create new testimoni", testimoni)
}

func (th *testimoniHandler) GetAllTestimoni(c *gin.Context) {
	allTestimoni, err := th.TestimoniUsecase.GetAllTestimoni()
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "successfully get all testimoni", allTestimoni)
}

func (th *testimoniHandler) GetTestimoniByID(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "error convert string to uint", err)
		return
	}

	detailTesti, err := th.TestimoniUsecase.GetTestimoniByID(uint(idUint))

	if err == gorm.ErrRecordNotFound {
		utils.FailureOrErrorResponse(c, http.StatusNotFound, "testimoni not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "testimoni found", detailTesti)
}

func (th *testimoniHandler) UpdateTestimoniByID(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "error convert string to uint", err)
		return
	}

	testimoni := c.PostForm("testimoni")
	jumlahDLStr := c.PostForm("jumlah_dl")
	image, _ := c.FormFile("gambar")

	JumlahDL, _ := strconv.Atoi(jumlahDLStr)

	updateTestimoni, err := th.TestimoniUsecase.UpdateTestimoniByID(uint(idUint),image, testimoni, JumlahDL)
	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "successfully create new testimoni", updateTestimoni)
}

func (th *testimoniHandler) DeleteTestimoniByID(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		utils.FailureOrErrorResponse(c, http.StatusBadRequest, "error convert string to uint", err)
		return
	}

	if err := th.TestimoniUsecase.DeleteTestimoniByID(uint(idUint)); err != nil{
		utils.FailureOrErrorResponse(c, http.StatusInternalServerError, "failed to delete testimoni by id", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "successfully deteled testimoni", nil)
}