package test

import (
	"short-code/initialize"
	"short-code/service/domain/unuse"
	"sync"
	"testing"
)

func init() {
	initialize.ConfByEnv("local")
	initialize.Logrus()

	initialize.RedisClient()
	initialize.LocalCache()
	initialize.GRomClickHouse()
	initialize.Limiter()
}

func TestReadTime(t *testing.T) {
	var service = unuse.UnUseCodeService{}
	service.Poll()
}

func TestLoad1(t *testing.T) {
	var service = unuse.UnUseCodeService{}
	var wg = sync.WaitGroup{}
	runSize := 10
	wg.Add(runSize)
	for i := 0; i < runSize; i++ {
		go func() {
			err := service.Load()
			wg.Done()
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
	t.Log("加载的短码数量 = ")
}

func TestPoll(t *testing.T) {
	var service = unuse.UnUseCodeService{}
	runSize := 30000
	var wg = sync.WaitGroup{}
	wg.Add(runSize)

	for i := 0; i < runSize; i++ {
		go func() {
			service.Poll()
			wg.Done()
		}()
	}
	wg.Wait()
	t.Log("剩下的code数量 = ")

}
