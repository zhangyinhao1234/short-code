package server

import (
	"short-code/api"
	"short-code/initialize"
)

type QueryAndBindServer struct {
	DefaultServer
}

func (e *QueryAndBindServer) StartUp() {
	e.Init()
	initialize.Tasks()
	initialize.LoadCacheInLocal()
	api.BindRequest(e.router)
	e.Run()
}
