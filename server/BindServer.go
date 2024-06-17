package server

import (
	"short-code/api"
)

type BindServer struct {
	DefaultServer
}

func (e *BindServer) StartUp() {
	e.Init()
	api.BindWriteCodeRequest(e.router)
	e.Run()
}
