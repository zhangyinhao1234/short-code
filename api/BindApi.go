package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "short-code/model/dto"
	"short-code/service/choreography"
	"short-code/utils/validator"
)

var (
	bindDataService = choreography.BindDataService{}
)

type BindApi struct {
}

func (e *BindApi) BindRequest(router *gin.Engine) {
	router.POST("/bind", e.bind)
}

func (e *BindApi) bind(context *gin.Context) {
	var bindData BindDto
	_ = context.ShouldBindJSON(&bindData)
	if err := validator.Validate(&bindData); err != nil {
		context.JSON(http.StatusBadRequest, DefaultErrorResult(err.Error()))
		return
	}
	shotCode, shotErr := bindDataService.Bind(&bindData.Data)
	writeResult(context, shotErr, shotCode)
}
