package server

import (
	"short-code/api"
)

type QueryServer struct {
	DefaultServer
}

func (e *QueryServer) StartUp() {
	e.Init()
	api.BindQueryRequest(e.router)
	e.Run()
}
