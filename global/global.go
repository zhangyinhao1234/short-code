package global

import (
	"github.com/bluele/gcache"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
	"short-code/conf"
)

var (
	DB             *gorm.DB
	CONF           conf.AppConf
	Logger         = logrus.New() // 初始化日志对象
	LOG            *logrus.Entry
	RedisClient    *redis.ClusterClient
	LocalCache     gcache.Cache
	DBQueryLimiter *rate.Limiter
)
