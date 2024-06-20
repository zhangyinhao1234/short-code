package task

type StartUpLoadCache struct {
}

func (e *StartUpLoadCache) LoadCacheInLocal() {
	go bindingService.LoadBindCacheInLocal()
	unUseCodeService.Load()
}
