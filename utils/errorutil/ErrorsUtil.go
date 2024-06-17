package errorutil

import (
	"errors"
	"short-code/model/do"
)

var (
	NotFindError      = NewBizErr("未能找到数据", 404)
	SystemError       = NewSysErr(errors.New("系统内部异常,请联系管理员"))
	LimiterError      = NewBizErr("限流保护,请稍后再试", 40011)
	CodeDepletedError = NewBizErr("获取短码超时，资源可能已经耗尽，请联系系统管理员", 40010)
	RedisAndCKError   = NewBizErr("Redis和ClickHouse都故障了，无法存储数据", 50001)
)

func NewSysErr(err error) *do.ShortCodeError {
	return &do.ShortCodeError{Err: err, Code: 500}
}

func NewBizErr(msg string, code int) *do.ShortCodeError {
	return &do.ShortCodeError{Err: errors.New(msg), Code: code}
}
