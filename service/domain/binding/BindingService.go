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

func (e *BindingService) Binding(shotCode *string, data *string) (*string, *do.ShortCodeError) {
	var ctx = context.Background()
	oldCode, err := global.RedisClient.Get(ctx, *data).Result()
	if err == redis.Nil {
		if err = bindingDataMapper.cacheAndSave(shotCode, data); err != nil {
			return nil, errorutil.RedisAndCKError
		}
		return shotCode, nil
	}
	return &oldCode, nil
}

func (e *BindingService) Flush() {
	bindingDataMapper.flushInLock()
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
