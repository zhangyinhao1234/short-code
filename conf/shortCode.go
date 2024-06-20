package conf

type ShortCode struct {
	CacheSize                         int64 `yaml:"cacheSize"`
	SafetyStock                       int64 `yaml:"safetyStock"`
	BindDataLocalCacheSize            int   `yaml:"bindDataLocalCacheSize"`
	TotalSize                         int64 `yaml:"totalSize"`
	BatchFlushSize                    int   `yaml:"batchFlushSize"`
	DbQueryLimit                      int64 `yaml:"dbQueryLimit"`
	StartUpLoadBindDataLocalCacheSize int64 `yaml:"startUpLoadBindDataLocalCacheSize"`
}
