package initialize

import (
	"golang.org/x/time/rate"
	"short-code/global"
)

func Limiter() {
	global.DBQueryLimiter = rate.NewLimiter(rate.Limit(global.CONF.ShotCode.DbQueryLimit), 8)
}
