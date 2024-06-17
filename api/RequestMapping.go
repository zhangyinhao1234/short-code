package api

import "github.com/gin-gonic/gin"

type RequestMapping interface {
	BindRequest(router *gin.Engine)
}

var (
	queryApi RequestMapping = &QueryApi{}
	bindApi  RequestMapping = &BindApi{}
	//add other controller
)

func BindRequest(router *gin.Engine) {
	BindQueryRequest(router)
	BindWriteCodeRequest(router)
}

func BindQueryRequest(router *gin.Engine) {
	queryApi.BindRequest(router)
}

func BindWriteCodeRequest(router *gin.Engine) {
	bindApi.BindRequest(router)
}
