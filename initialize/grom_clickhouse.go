package initialize

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"short-code/global"
	"time"
)

func GRomClickHouse() {
	dsn := global.CONF.ClickHouse.DSN
	logMode := logger.Default.LogMode(logger.Info)
	if "error" == global.CONF.Logger.Level {
		logMode = logger.Default.LogMode(logger.Error)
	}
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{
		Logger: logMode,
	})

	if "debug" == global.CONF.Logger.Level {
		db.Debug()
	}
	if err != nil {
		global.LOG.Fatalf("创建数据库连接失败:%v", err)
	} else {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.SetMaxIdleConns(32)
			sqlDB.SetMaxOpenConns(32)
			sqlDB.SetConnMaxLifetime(time.Hour * 1)
		} else {
			global.LOG.Fatalf("配置连接池失败:%v", err)
		}
		global.DB = db
	}
}
