package unuse

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"short-code/global"
	"short-code/model/do"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	base52Chars                 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	PollTimeOut                 = 3000
	CurrentSerialNumberCacheKey = "ShotCode_CurrentSerialNumber"
)

var (
	unUseCodeMapper           = UnUseCodeMapper{}
	currentSerialNumberMapper = CurrentSerialNumberMapper{}

	exTime = time.Hour * 24
	codeCh = make(chan string, 150000)

	loadingCodeMark int64
	size            int64
)

type UnUseCodeService struct {
}

func (e *UnUseCodeService) Poll() (*string, *do.ShortCodeError) {
	if atomic.LoadInt64(&size) < global.CONF.ShotCode.SafetyStock && atomic.CompareAndSwapInt64(&loadingCodeMark, 0, 1) {
		go func() {
			e.Load()
			atomic.StoreInt64(&loadingCodeMark, 0)
		}()
	}
	str := <-codeCh
	atomic.AddInt64(&size, -1)
	return &str, nil
}

func (e *UnUseCodeService) Load() error {
	if atomic.LoadInt64(&size) > global.CONF.ShotCode.SafetyStock {
		return nil
	}
	//global.LOG.Debug("准备加载数据，当前剩余短码 = ", e.size)
	number, err := e.getSerialNumber()
	if err != nil {
		return err
	}
	codes, err := unUseCodeMapper.listShortCodeFromDB(number)
	if err != nil {
		return err
	}
	for _, code := range *codes {
		c := e.parseNumber2Str(code.ShotCode)
		codeCh <- *c
	}
	atomic.AddInt64(&size, int64(len(*codes)))
	return nil
}

func (e *UnUseCodeService) getSerialNumber() (int64, error) {
	var ctx = context.Background()
	currentNum, err := e.getCurrentSerialNumber()
	if err != nil {
		return 0, err
	}
	nextNumber := currentNum + global.CONF.ShotCode.CacheSize
	if nextNumber > global.CONF.ShotCode.TotalSize {
		nextNumber = 0
	}
	t1 := time.Now().UnixMilli()
	for {
		if time.Now().UnixMilli()-t1 > PollTimeOut {
			return 0, errors.New("获取短码序号超时,请联系管理员！！")
		}
		success, err := global.RedisClient.SetNX(ctx, e.getSerialNumberCacheKey(nextNumber), nextNumber, exTime).Result()
		if err != nil {
			return 0, err
		}
		if success {
			if nextNumber > global.CONF.ShotCode.TotalSize {
				return 0, errors.New("短码耗尽,请联系管理员！")
			}
			global.RedisClient.Set(ctx, CurrentSerialNumberCacheKey, nextNumber, exTime)
			go currentSerialNumberMapper.saveCurrentSerialNumberInDB(nextNumber)
			return nextNumber, nil
		}
		nextNumber = nextNumber + global.CONF.ShotCode.CacheSize
	}
}

func (e *UnUseCodeService) getCurrentSerialNumber() (int64, error) {
	var ctx = context.Background()
	number, err := global.RedisClient.Get(ctx, CurrentSerialNumberCacheKey).Result()
	if err == redis.Nil {
		number, err = currentSerialNumberMapper.getCurrentSerialNumberFromDB()
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}
	currentNum, _ := strconv.ParseInt(number, 10, 64)
	return currentNum, nil
}

func (e *UnUseCodeService) getSerialNumberCacheKey(SerialNumber int64) string {
	return CurrentSerialNumberCacheKey + "#" + strconv.FormatInt(SerialNumber, 10)
}

func (e *UnUseCodeService) parseNumber2Str(decimal int64) *string {
	if decimal == 0 {
		str := "AAAAAA"
		return &str
	}
	result := ""
	for decimal > 0 {
		remainder := decimal % 52
		result = string(base52Chars[remainder]) + result
		decimal /= 52
	}
	for len(result) < 6 {
		result = "A" + result
	}
	return &result
}
