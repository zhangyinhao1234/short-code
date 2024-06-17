package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "short-code/model/dto"
	"short-code/service/choreography"
)

var (
	getDataService = choreography.GetDataService{}
)

type QueryApi struct {
}

func (e *QueryApi) BindRequest(router *gin.Engine) {
	router.GET("/getByCode", e.getByShotCode)
}

func (e *QueryApi) getByShotCode(context *gin.Context) {
	code, find := context.GetQuery("code")
	if !find {
		context.JSON(http.StatusBadRequest, ErrorResult(http.StatusBadRequest, "缺少参数"))
		return
	}
	data, shotErr := getDataService.GetByShotCode(&code)
	writeResult(context, shotErr, data)
}
