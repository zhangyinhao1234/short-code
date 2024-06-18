package initialize

import "short-code/task"

var (
	flushBindDataTask = task.FlushBindDataTask{}
	startUpLoadCache  = task.StartUpLoadCache{}
)

func Tasks() {
	flushBindDataTask.Run()
}

func LoadCacheInLocal() {
	startUpLoadCache.LoadBindCacheInLocal()
}
