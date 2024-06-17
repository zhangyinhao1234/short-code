package dto

type BaseResult struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

func SuccessResult(data interface{}) *BaseResult {
	return &BaseResult{data, "调用成功", 200}
}

func ErrorResult(code int, message string) *BaseResult {
	return &BaseResult{nil, message, code}
}

func DefaultErrorResult(message string) *BaseResult {
	return &BaseResult{nil, message, 500}
}
