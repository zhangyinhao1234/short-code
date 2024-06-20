package unuse

import (
	"short-code/global"
	"strconv"
	"time"
)

type CurrentSerialNumber struct {
	SerialNumber int64 `gorm:"column:serial_number"`
	CreateTime   int64 `gorm:"column:create_time"`
}

func (CurrentSerialNumber) TableName() string {
	return "sc_current_serial_number"
}

type CurrentSerialNumberMapper struct {
}

func (e *CurrentSerialNumberMapper) getCurrentSerialNumberFromDB() (string, error) {
	var queryList []CurrentSerialNumber
	result := global.DB.Limit(1).Order("create_time desc ").Find(&queryList)
	if result.Error != nil {
		return "", result.Error
	}
	return strconv.FormatInt(queryList[0].SerialNumber, 10), nil
}

func (e *CurrentSerialNumberMapper) saveCurrentSerialNumberInDB(serialNumber int64) {
	currentSerialNumber := CurrentSerialNumber{SerialNumber: serialNumber, CreateTime: time.Now().UnixMilli()}
	global.DB.Create(currentSerialNumber)
}
