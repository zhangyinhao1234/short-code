package test

import (
	"github.com/google/uuid"
	"short-code/global"
	"short-code/initialize"
	"short-code/service/choreography"
	"sync"
	"testing"
	"time"
)

func init() {
	initialize.ConfByEnv("local")
	initialize.Logrus()

	initialize.RedisClient()
	initialize.LocalCache()
	initialize.GRomClickHouse()
	initialize.Limiter()
}

func TestBind(t *testing.T) {
	global.LOG.Info("start ........")
	var bindDataService = choreography.BindDataService{}
	runSize := 20
	for i := 0; i < runSize; i++ {
		id := uuid.New().String() + uuid.New().String() + uuid.New().String() + uuid.New().String()
		_, err := bindDataService.Bind(&id)
		if err != nil {
			t.Log(err)
		}

	}
	bindDataService.Flush()
	global.LOG.Info("end ........")
	//time.Sleep(time.Second * 5)
}

func TestBind_V2(t *testing.T) {
	var bindDataService = choreography.BindDataService{}
	runSize := 300
	var wg = sync.WaitGroup{}
	wg.Add(runSize)
	for i := 0; i < runSize; i++ {
		//time.Sleep(time.Millisecond * 3)
		go func() {
			id := uuid.New().String() + uuid.New().String() + uuid.New().String() + uuid.New().String()
			_, err := bindDataService.Bind(&id)
			if err != nil {
				t.Log(err.GetMsg())
			}
			wg.Done()
		}()
	}
	wg.Wait()

	bindDataService.Flush()

	time.Sleep(time.Second * 5)

}
