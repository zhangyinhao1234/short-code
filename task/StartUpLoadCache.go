package task

type StartUpLoadCache struct {
}

func (e *StartUpLoadCache) LoadBindCacheInLocal() {
	bindDataService.LoadBindCacheInLocal()
}
