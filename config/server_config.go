package config

import (
	"os"

	"github.com/lfmexi/tcpgateway/handlers"
	"github.com/lfmexi/tcpgateway/server"
)

func handler() server.Handler {
	return handlers.NewConnectionHandler(portsDeviceTypeMap, sessionService(), publisherService(), packetEventSubscriberFactory())
}

// ConfigureServer creates a new tcp server with a configured connection handler
func ConfigureServer() server.Server {
	if envhost := os.Getenv("SERVER_HOST"); envhost != "" {
		configuration.Server.Host = envhost
	}

	connectionHandler := handler()

	return server.NewTCPServer(configuration.Server.Host, portsDeviceTypeMap, connectionHandler)
}
