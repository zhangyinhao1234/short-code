package server

import (
	"short-code/api"
)

type QueryAndBindServer struct {
	DefaultServer
}

func (e *QueryAndBindServer) StartUp() {
	e.Init()
	api.BindRequest(e.router)
	e.Run()
}
