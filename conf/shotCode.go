package conf

type ShotCode struct {
	CacheSize              int64     `yaml:"cacheSize"`
	SafetyStock            int64     `yaml:"safetyStock"`
	BindDataLocalCacheSize int       `yaml:"bindDataLocalCacheSize"`
	TotalSize              int64     `yaml:"totalSize"`
	BatchFlushSize         int       `yaml:"batchFlushSize"`
	DbQueryLimit           int64     `yaml:"dbQueryLimit"`
	DataTable              DataTable `yaml:"dataTable"`
}

type DataTable struct {
	BindingData         string `yaml:"BindingData"`
	UnUseCode           string `yaml:"UnUseCode"`
	CurrentSerialNumber string `yaml:"CurrentSerialNumber"`
}
