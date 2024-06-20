package unuse

import (
	"short-code/global"
)

type UnUseCode struct {
	Code         int64 `gorm:"column:code"`
	SerialNumber int64 `gorm:"column:serial_number"`
}

func (UnUseCode) TableName() string {
	return "sc_code"
}

type UnUseCodeMapper struct {
	ShotCode     string `gorm:"column:shot_code"`
	SerialNumber int64  `gorm:"column:serial_number"`
}

func (e *UnUseCodeMapper) listShortCodeFromDB(SerialNumber int64) (*[]UnUseCode, error) {
	var queryList []UnUseCode
	result := global.DB.Limit(int(global.CONF.ShortCode.CacheSize)).Where(" serial_number >= ?", SerialNumber).Order("serial_number asc ").Select("code").Find(&queryList)
	if result.Error != nil {
		return nil, result.Error
	}
	return &queryList, nil
}
