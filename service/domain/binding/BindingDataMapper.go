package binding

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"short-code/global"
	"short-code/utils/errorutil"
	"sync"
	"time"
)

var (
	unSaveBindingData  []BindingData
	lazySaveMut        = sync.Mutex{}
	SaveCKErrorCodeKey = "SaveCKErrorCode"
	exTime             = time.Hour * 720
	localCacheExTime   = time.Hour * 168
)

type BindingData struct {
	ShotCode   string `gorm:"column:shot_code"`
	Message    string `gorm:"column:message"`
	CreateTime int64  `gorm:"column:create_time"`
}

func (BindingData) TableName() string {
	if global.CONF.ShotCode.DataTable.BindingData != "" {
		return global.CONF.ShotCode.DataTable.BindingData
	}
	return "short_code_binding_data"
}

type BindingDataMapper struct {
}

func (e *BindingDataMapper) cacheAndSave(shotCode *string, data *string) error {
	redisErr := e.cacheInRedis(shotCode, data)
	ckErr := e.lazySaveInCK(*shotCode, *data)
	if redisErr != nil && ckErr != nil {
		return errors.New(errorutil.RedisAndCKError.GetMsg())
	}
	e.cacheInLocal(shotCode, data)
	return nil
}

// 500TPS基本满足需求，chan遇到宕机可能会丢失更多的数据
func (e *BindingDataMapper) lazySaveInCK(shotCode string, data string) error {
	defer lazySaveMut.Unlock()
	lazySaveMut.Lock()
	bindingData := BindingData{ShotCode: shotCode, Message: data, CreateTime: time.Now().UnixMilli()}
	unSaveBindingData = append(unSaveBindingData, bindingData)
	if len(unSaveBindingData) >= global.CONF.ShotCode.BatchFlushSize {
		return e.flush()
	}
	return nil
}

func (e *BindingDataMapper) flushInLock() {
	defer lazySaveMut.Unlock()
	lazySaveMut.Lock()
	//因故障导致存储的数据量不会很大，分布式环境中多存储了几份问题不大
	var stagingData []BindingData
	for _, v := range *e.listStagingFromRedis() {
		unSaveBindingData = append(unSaveBindingData, v)
		stagingData = append(stagingData, v)
	}
	err := e.flush()
	if err == nil {
		e.delStagingInRedis(&stagingData)
	}
}

func (e *BindingDataMapper) flush() error {
	if len(unSaveBindingData) == 0 {
		return nil
	}
	//global.LOG.Info("开始批量刷写数据，数组数量{2};", len(unSaveBindingData))
	result := global.DB.Create(&unSaveBindingData)
	if result.Error != nil {
		e.stagingInRedis(unSaveBindingData)
		return result.Error
	}
	unSaveBindingData = []BindingData{}
	return nil
}

func (e *BindingDataMapper) stagingInRedis(datas []BindingData) {
	//global.LOG.Info("数据暂存Redis size = ", len(datas))
	if len(datas) == 0 {
		return
	}
	var ctx = context.Background()
	nm := map[string]string{}
	for _, v := range datas {
		nm[v.ShotCode] = v.Message
	}
	global.RedisClient.HSet(ctx, SaveCKErrorCodeKey, nm)
}

func (e *BindingDataMapper) listStagingFromRedis() *[]BindingData {
	var ctx = context.Background()
	m, _ := global.RedisClient.HGetAll(ctx, SaveCKErrorCodeKey).Result()
	datas := []BindingData{}
	for k, v := range m {
		bindingData := BindingData{ShotCode: k, Message: v, CreateTime: time.Now().UnixMilli()}
		datas = append(datas, bindingData)
	}
	//global.LOG.Info("Redis 暂存的数据 size = ", len(datas))
	return &datas
}

func (e *BindingDataMapper) delStagingInRedis(datas *[]BindingData) {
	var ctx = context.Background()
	for _, v := range *datas {
		global.RedisClient.HDel(ctx, SaveCKErrorCodeKey, v.ShotCode)
	}
}

func (e *BindingDataMapper) getByCode(shotCode *string) (string, error) {
	var queryList []BindingData
	result := global.DB.Limit(1).Where(" shot_code = ?", shotCode).Order("create_time desc ").Find(&queryList)
	if result.Error != nil {
		return "", result.Error
	}
	if len(queryList) == 0 {
		return "", nil
	}
	return queryList[0].Message, nil
}

func (e *BindingDataMapper) getLast() (*BindingData, error) {
	var queryList []BindingData
	result := global.DB.Limit(1).Where(" 1 = 1 ").Order("create_time desc ").Find(&queryList)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(queryList) == 0 {
		return nil, nil
	}
	return &queryList[0], nil
}

func (e *BindingDataMapper) listLtCreateTime(createTime int64, limit int) (*[]BindingData, error) {
	var queryList []BindingData
	result := global.DB.Limit(limit).Where(" create_time < ?", createTime).Order("create_time desc ").Find(&queryList)
	if result.Error != nil {
		return nil, result.Error
	}
	return &queryList, nil
}

func (e *BindingDataMapper) cacheInLocal(shotCode *string, data *string) {
_:
	global.LocalCache.SetWithExpire(*shotCode, data, localCacheExTime)
}

func (e *BindingDataMapper) cacheInRedis(shotCode *string, data *string) error {
	var ctx = context.Background()
	_, err := global.RedisClient.Set(ctx, *shotCode, data, exTime).Result()
	if err != nil {
		return err
	}
	global.RedisClient.Set(ctx, *data, shotCode, exTime)
	global.RedisClient.Set(ctx, *shotCode+"#MkS", "1", exTime)
	return nil
}

func (e *BindingDataMapper) existsMarkBind(shotCode *string) (bool, error) {
	var ctx = context.Background()
	mark, err := global.RedisClient.Get(ctx, *shotCode+"#MkS").Result()
	if err == redis.Nil {
		return true, nil
	} else if err != nil {
		return false, err
	}
	if mark == "1" {
		return true, nil
	} else {
		return false, nil
	}

}

func (e *BindingDataMapper) markUnBind(shotCode *string) {
	var ctx = context.Background()
	global.RedisClient.Set(ctx, *shotCode+"#MkS", "0", exTime)
}
