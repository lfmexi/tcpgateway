package main

import "bitbucket.org/challengerdevs/gpsdriver/config"

func main() {
	server := config.ConfigureServer()
	server.Listen()
}
