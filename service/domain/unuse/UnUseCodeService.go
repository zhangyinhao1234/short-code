package unuse

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"short-code/global"
	"short-code/model/do"
	"short-code/utils/errorutil"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	base52Chars                 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	PollTimeOut                 = 3000
	CurrentSerialNumberCacheKey = "ShotCode_CurrentSerialNumber"
)

var (
	loadMut = sync.Mutex{}
	pollMut = sync.Mutex{}

	unUseCodeMapper           = UnUseCodeMapper{}
	currentSerialNumberMapper = CurrentSerialNumberMapper{}

	exTime = time.Hour * 24
)

type UnUseCodeService struct {
	currentNode     *ShotCodeNode
	lastNode        *ShotCodeNode
	size            int64
	loadingCodeMark int64
}

type ShotCodeNode struct {
	shotCode string
	next     *ShotCodeNode
}

func (e *UnUseCodeService) Poll() (*string, *do.ShortCodeError) {
	defer func() {
		pollMut.Unlock()
	}()
	pollMut.Lock()
	t1 := time.Now().UnixMilli()
	for e.currentNode == nil {
		if time.Now().UnixMilli()-t1 > PollTimeOut {
			return nil, errorutil.CodeDepletedError
		}
		e.loadSwapLoadingCodeMark()
	}
	for e.currentNode.next == nil {
		if time.Now().UnixMilli()-t1 > PollTimeOut {
			return nil, errorutil.CodeDepletedError
		}
		e.loadSwapLoadingCodeMark()
	}
	shotCode := e.currentNode.shotCode
	e.currentNode = e.currentNode.next
	afterSize := atomic.AddInt64(&e.size, -1)

	if afterSize < global.CONF.ShotCode.SafetyStock && atomic.CompareAndSwapInt64(&e.loadingCodeMark, 0, 1) {
		//global.LOG.Info("达到安全库存，加载短码数据")
		go func() {
			e.Load()
			atomic.StoreInt64(&e.loadingCodeMark, 0)
		}()
	}
	return &shotCode, nil
}

func (e *UnUseCodeService) loadSwapLoadingCodeMark() {
	if atomic.CompareAndSwapInt64(&e.loadingCodeMark, 0, 1) {
		go func() {
			e.Load()
			atomic.StoreInt64(&e.loadingCodeMark, 0)
		}()
	}
}

func (e *UnUseCodeService) Load() error {
	defer func() {
		loadMut.Unlock()
	}()
	loadMut.Lock()
	if e.size > global.CONF.ShotCode.SafetyStock {
		return nil
	}
	//global.LOG.Debug("准备加载数据，当前剩余短码 = ", e.size)
	number, err := e.getSerialNumber()
	if err != nil {
		return err
	}
	codes, err := unUseCodeMapper.listShotCodeFromDB(number)
	if err != nil {
		return err
	}
	var first_ *ShotCodeNode = nil
	var move_ *ShotCodeNode = nil
	for _, code := range *codes {
		node := ShotCodeNode{shotCode: *e.parseNumber2Str(code.ShotCode), next: nil}
		if first_ == nil {
			first_ = &node
		} else {
			move_.next = &node
		}
		move_ = &node
	}
	if e.currentNode == nil {
		e.currentNode = first_
	}
	if e.lastNode == nil {
		e.lastNode = move_
	} else {
		e.lastNode.next = first_
		e.lastNode = move_
	}
	atomic.AddInt64(&e.size, int64(len(*codes)))
	//global.LOG.Debug("加载数据完成，当前剩余短码 = ", e.size)
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
