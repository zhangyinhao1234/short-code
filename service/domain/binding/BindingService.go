package binding

import (
	"context"
	"github.com/redis/go-redis/v9"
	"short-code/global"
	"short-code/model/do"
	"short-code/utils/errorutil"
)

var (
	bindingDataMapper = BindingDataMapper{}
)

type BindingService struct {
}

func init() {
	bindingDataMapper.startLazySaveChConsumer()
}

func (e *BindingService) Binding(shotCode *string, data *string) (*string, *do.ShortCodeError) {
	var ctx = context.Background()
	oldCode, err := global.RedisClient.Get(ctx, bindingDataMapper.md5(data)).Result()
	if err == redis.Nil {
		if err = bindingDataMapper.cacheAndSave(shotCode, data); err != nil {
			return nil, errorutil.RedisAndCKError
		}
		return shotCode, nil
	}
	return &oldCode, nil
}

func (e *BindingService) Flush() {
	bindingDataMapper.flushCK()
}

func (e *BindingService) DestroyAndFlush() {
	global.LOG.Info("服务被销毁,准备持久化数据")
	lazySaveWG.Wait()
	bindingDataMapper.flushCK()
	lazySaveWG.Wait()
	bindingDataMapper.closeLazySaveChConsumer()
}

func (e *BindingService) GetByShotCode(shotCode *string) (*string, *do.ShortCodeError) {
	value, _ := global.LocalCache.Get(*shotCode)
	if value != nil {
		str := value.(*string)
		return str, nil
	}
	var ctx = context.Background()
	var data string
	var err error

	data, err = global.RedisClient.Get(ctx, *shotCode).Result()
	if err == redis.Nil {
		if bind, err := bindingDataMapper.existsMarkBind(shotCode); err != nil {
			return nil, errorutil.NewSysErr(err)
		} else if !bind {
			return nil, errorutil.NotFindError
		}
		if !global.DBQueryLimiter.Allow() {
			return nil, errorutil.LimiterError
		}
		if data, err = bindingDataMapper.getByCode(shotCode); err != nil {
			return nil, errorutil.NewSysErr(err)
		} else if "" == data {
			bindingDataMapper.markUnBind(shotCode)
			return nil, errorutil.NotFindError
		}
		go func() {
			bindingDataMapper.cacheInLocal(shotCode, &data)
		_:
			bindingDataMapper.cacheInRedis(shotCode, &data)
		}()
	} else if err != nil {
		if !global.DBQueryLimiter.Allow() {
			return nil, errorutil.LimiterError
		}
		global.LOG.Error("Redis故障,通过短码获取数据失败,从数据库加载数据,进行限流", err)
		if data, err = bindingDataMapper.getByCode(shotCode); err != nil {
			return nil, errorutil.NewSysErr(err)
		} else if "" == data {
			bindingDataMapper.markUnBind(shotCode)
			return nil, errorutil.NotFindError
		}
	}
	if data == "" {
		return nil, errorutil.NotFindError
	}
	bindingDataMapper.cacheInLocal(shotCode, &data)
	return &data, nil
}

func (e *BindingService) LoadBindCacheInLocal() {
	if global.CONF.ShortCode.StartUpLoadBindDataLocalCacheSize <= 0 {
		return
	}
	bdata, err := bindingDataMapper.getLast()
	if err != nil {
		global.LOG.Error("获取最后绑定数据异常", err)
		return
	}
	if bdata == nil {
		return
	}
	lastCreateTime := bdata.CreateTime
	limit := 200000
	inSize := int64(0)
	for inSize < global.CONF.ShortCode.StartUpLoadBindDataLocalCacheSize {
		datas, err := bindingDataMapper.listLtCreateTime(lastCreateTime, limit)
		if err != nil {
			global.LOG.Error("加载绑定数据到本地内存读取数据异常", err)
			return
		}
		if len(*datas) == 0 {
			return
		}
		for _, d := range *datas {
			bindingDataMapper.cacheInLocal(&d.Code, &d.Message)
			lastCreateTime = d.CreateTime
		}
		//global.LOG.Info("加载数量,加载下一个时间", inSize, lastCreateTime)
		inSize = inSize + int64(200000)
	}
}
