package unuse

import (
	"short-code/initialize"
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
	var service = UnUseCodeService{}
	service.Poll()
}

func TestLoad1(t *testing.T) {
	var service = UnUseCodeService{}
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
	t.Log("加载的短码数量 = ", service.size)
	m1, count := getCodesMap(service)
	t.Log("短码遍历次数 = ", count)
	t.Log("MAP中短码数量 = ", len(m1))
}

func getCodesMap(service UnUseCodeService) (map[string]int, int) {
	var m1 = map[string]int{}
	count := 0
	for service.currentNode != nil {
		m1[service.currentNode.shotCode] = 0
		service.currentNode = service.currentNode.next
		count++
	}
	return m1, count
}

func TestPoll(t *testing.T) {
	var service = UnUseCodeService{}
	runSize := 30000
	var wg = sync.WaitGroup{}
	wg.Add(runSize)
	var mcode = map[string]int{}
	var mt = sync.Mutex{}
	for i := 0; i < runSize; i++ {
		go func() {
			code, _ := service.Poll()
			mt.Lock()
			mcode[*code] = 0
			mt.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	t.Log("获取的Code数量 = ", len(mcode))
	t.Log("剩下的code数量 = ", service.size)
	m1, count := getCodesMap(service)
	t.Log("map遍历次数 = ", count)
	t.Log("map剩余数量 = ", len(m1))
}
