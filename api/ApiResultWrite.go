package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"short-code/global"
	"short-code/model/do"
	. "short-code/model/dto"
	"short-code/utils/errorutil"
)

func writeResult(context *gin.Context, shotErr *do.ShortCodeError, data interface{}) {
	if shotErr != nil {
		global.LOG.Error(shotErr.Err)
		if shotErr.IsBizErr() {
			if shotErr.IsNotFindErr() {
				context.JSON(http.StatusNotFound, ErrorResult(shotErr.Code, shotErr.GetMsg()))
			} else {
				context.JSON(http.StatusBadRequest, ErrorResult(shotErr.Code, shotErr.GetMsg()))
			}
		} else {
			context.JSON(http.StatusInternalServerError, ErrorResult(shotErr.Code, errorutil.SystemError.GetMsg()))
		}
		return
	}
	context.JSON(http.StatusOK, SuccessResult(data))
}
