package dto

type BindDto struct {
	Data string `json:"data" validate:"required" label:"绑定数据" `
}
