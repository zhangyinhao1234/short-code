package server

import (
	"short-code/api"
	"short-code/initialize"
)

type QueryServer struct {
	DefaultServer
}

func (e *QueryServer) StartUp() {
	e.Init()
	initialize.LoadCacheInLocal()
	api.BindQueryRequest(e.router)
	e.Run()
}
