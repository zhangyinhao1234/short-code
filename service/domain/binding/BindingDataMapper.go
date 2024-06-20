package binding

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
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

	ckFault            = false
	lazySaveChPool     []chan BindingData
	lazySaveChPoolSize = 2
	lazySaveWG         sync.WaitGroup
)

type BindingData struct {
	Code       string `gorm:"column:code"`
	Message    string `gorm:"column:message"`
	CreateTime int64  `gorm:"column:create_time"`
}

func (BindingData) TableName() string {
	return "sc_binding_data"
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

func (e *BindingDataMapper) lazySaveInCK(shotCode string, data string) error {
	if ckFault {
		return errors.New("ClickHouse上次被标记为故障")
	}
	bindingData := BindingData{Code: shotCode, Message: data, CreateTime: time.Now().UnixMilli()}
	ascii := int([]rune(shotCode[0:5])[0])
	lazySaveCh := lazySaveChPool[ascii%lazySaveChPoolSize]
	lazySaveWG.Add(1)
	lazySaveCh <- bindingData
	return nil
}

func (e *BindingDataMapper) startLazySaveChConsumer() {
	for i := 0; i < lazySaveChPoolSize; i++ {
		lazySaveCh := make(chan BindingData, 2000)
		lazySaveChPool = append(lazySaveChPool, lazySaveCh)
		e.lazySaveChConsumer(lazySaveCh)
	}
}

func (e *BindingDataMapper) closeLazySaveChConsumer() {
	for _, ch := range lazySaveChPool {
		close(ch)
	}
}

func (e *BindingDataMapper) lazySaveChConsumer(lazySaveCh chan BindingData) {
	go func() {
		var datas []BindingData
		for v := range lazySaveCh {
			if "#FlushCKNow#" == v.Code {
				//global.LOG.Info("接受数据持久化到CK指令")
				//因故障导致存储的数据量不会很大，分布式环境中多存储了几份问题不大
				var stagingData []BindingData
				for _, v := range *e.listStagingFromRedis() {
					datas = append(datas, v)
					stagingData = append(stagingData, v)
				}
				result := global.DB.Create(&datas)
				if result.Error == nil {
					e.delStagingInRedis(&stagingData)
					datas = datas[0:0]
				}
			} else {
				datas = append(datas, v)
			}

			lazySaveWG.Done()
			if len(datas) >= global.CONF.ShortCode.BatchFlushSize {
				global.LOG.Info("数据持久化到CK size = ", len(datas))
				var replica []BindingData
				for _, v := range datas {
					replica = append(replica, v)
				}
				datas = datas[0:0]
				go func() {
					result := global.DB.Create(&replica)
					if result.Error != nil {
						e.stagingInRedis(replica)
						ckFault = true
					} else {
						ckFault = false
					}
				}()
			}
		}
	}()
}

func (e *BindingDataMapper) flushCK() {
	lazySaveWG.Add(lazySaveChPoolSize)
	for _, ch := range lazySaveChPool {
		bindingData := BindingData{Code: "#FlushCKNow#"}
		ch <- bindingData
	}
}

func (e *BindingDataMapper) stagingInRedis(datas []BindingData) {
	//global.LOG.Info("数据暂存Redis size = ", len(datas))
	if len(datas) == 0 {
		return
	}
	var ctx = context.Background()
	nm := map[string]string{}
	for _, v := range datas {
		nm[v.Code] = v.Message
	}
	global.RedisClient.HSet(ctx, SaveCKErrorCodeKey, nm)
}

func (e *BindingDataMapper) listStagingFromRedis() *[]BindingData {
	var ctx = context.Background()
	m, _ := global.RedisClient.HGetAll(ctx, SaveCKErrorCodeKey).Result()
	var datas []BindingData
	for k, v := range m {
		bindingData := BindingData{Code: k, Message: v, CreateTime: time.Now().UnixMilli()}
		datas = append(datas, bindingData)
	}
	//global.LOG.Info("Redis 暂存的数据 size = ", len(datas))
	return &datas
}

func (e *BindingDataMapper) delStagingInRedis(datas *[]BindingData) {
	var ctx = context.Background()
	for _, v := range *datas {
		global.RedisClient.HDel(ctx, SaveCKErrorCodeKey, v.Code)
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
	global.LocalCache.SetWithExpire(*shotCode, data, localCacheExTime)
}

func (e *BindingDataMapper) cacheInRedis(shotCode *string, data *string) error {
	var ctx = context.Background()
	_, err := global.RedisClient.Set(ctx, *shotCode, data, exTime).Result()
	if err != nil {
		return err
	}
	global.RedisClient.Set(ctx, e.md5(data), shotCode, exTime)
	return nil
}

func (e *BindingDataMapper) md5(data *string) string {
	re := md5.Sum([]byte(*data))
	md5str := fmt.Sprintf("%x", re)
	return md5str
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
