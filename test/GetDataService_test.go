package test

import (
	"short-code/global"
	"short-code/service/choreography"
	"short-code/service/domain/binding"
	"sync"
	"testing"
	"time"
)

func TestGetV1(t *testing.T) {
	var getDataService = choreography.GetDataService{}
	code := "AABTUg"
	data, err := getDataService.GetByShotCode(&code)
	if err != nil {
		t.Log(err.GetMsg())
	}
	t.Log("get val = ", data)
	data, err = getDataService.GetByShotCode(&code)
	if err != nil {
		t.Log(err.GetMsg())
	}
	t.Log("get val = ", data)

	time.Sleep(time.Second * 2)
}

func TestGetV2(t *testing.T) {
	var getDataService = choreography.GetDataService{}
	runSize := 6000
	var queryList []binding.BindingData
	global.DB.Limit(runSize).Where(" 1=1 ").Order("create_time desc ").Find(&queryList)

	var wg = sync.WaitGroup{}
	wg.Add(runSize)
	for i := 0; i < runSize; i++ {
		go func(codes []binding.BindingData) {
			_, err := getDataService.GetByShotCode(&codes[i].ShotCode)
			if err != nil {
				t.Log(err.GetMsg())
			}
			wg.Done()
		}(queryList)
	}
	wg.Wait()
	time.Sleep(time.Second * 2)

	wg = sync.WaitGroup{}
	wg.Add(runSize)
	for i := 0; i < runSize; i++ {
		go func(codes []binding.BindingData) {
			_, err := getDataService.GetByShotCode(&codes[i].ShotCode)
			if err != nil {
				t.Log(err.GetMsg())
			}
			wg.Done()
		}(queryList)
	}
	wg.Wait()

	time.Sleep(time.Second * 5)
}
