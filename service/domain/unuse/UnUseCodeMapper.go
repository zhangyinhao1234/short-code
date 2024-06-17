package unuse

import (
	"short-code/global"
)

type UnUseCode struct {
	ShotCode     int64 `gorm:"column:shot_code"`
	SerialNumber int64 `gorm:"column:serial_number"`
}

func (UnUseCode) TableName() string {
	if global.CONF.ShotCode.DataTable.UnUseCode != "" {
		return global.CONF.ShotCode.DataTable.UnUseCode
	}
	return "short_code_code"
}

type UnUseCodeMapper struct {
	ShotCode     string `gorm:"column:shot_code"`
	SerialNumber int64  `gorm:"column:serial_number"`
}

func (e *UnUseCodeMapper) listShotCodeFromDB(SerialNumber int64) (*[]UnUseCode, error) {
	var queryList []UnUseCode
	result := global.DB.Limit(int(global.CONF.ShotCode.CacheSize)).Where(" serial_number >= ?", SerialNumber).Order("serial_number asc ").Find(&queryList)
	if result.Error != nil {
		return nil, result.Error
	}
	return &queryList, nil
}
