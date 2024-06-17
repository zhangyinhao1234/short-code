package main

import (
	"short-code/server"
)

func main() {
	var webServer server.Server = &server.BindServer{}
	webServer.StartUp()
}
