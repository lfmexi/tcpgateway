package main

import "github.com/lfmexi/tcpgateway/config"

func main() {
	server := config.ConfigureServer()
	server.Listen()
}
