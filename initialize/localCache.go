package initialize

import (
	"github.com/bluele/gcache"
	"short-code/global"
)

func LocalCache() {
	global.LocalCache = gcache.New(global.CONF.ShortCode.BindDataLocalCacheSize).LRU().Build()
}
