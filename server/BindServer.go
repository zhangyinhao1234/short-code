package server

import (
	"short-code/api"
	"short-code/initialize"
)

type BindServer struct {
	DefaultServer
}

func (e *BindServer) StartUp() {
	e.Init()
	initialize.Tasks()
	api.BindWriteCodeRequest(e.router)
	e.Run()
}
