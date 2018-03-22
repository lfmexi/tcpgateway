package main

import "bitbucket.org/challengerdevs/tcpgateway/config"

func main() {
	server := config.ConfigureServer()
	server.Listen()
}
