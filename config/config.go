package config

import (
	"fmt"
	"os"

	"bitbucket.org/challengerdevs/gpsdriver/handlers"
	"bitbucket.org/challengerdevs/gpsdriver/readers"
	"bitbucket.org/challengerdevs/gpsdriver/server"
	"bitbucket.org/challengerdevs/gpsdriver/writers"
)

var host = "0.0.0.0"
var port = "8889"

func readerServiceFactory() readers.ReaderServiceFactory {
	return readers.NewContinuousReaderServiceFactory()
}

func writerServiceFactory() writers.WriterServiceFactory {
	return writers.NewAsyncWriterServiceFactory()
}

func handler() server.Handler {
	return handlers.NewConnectionHandler(sessionService(), publisherService(), packetEventSubscriberFactory(), readerServiceFactory(), writerServiceFactory())
}

// ConfigureServer creates a new tcp server with a configured connection handler
func ConfigureServer() server.Server {
	if envhost := os.Getenv("SERVER_HOST"); envhost != "" {
		host = envhost
	}

	if envport := os.Getenv("SERVER_PORT"); envport != "" {
		port = envport
	}

	connectionHandler := handler()

	return server.NewTCPServer(fmt.Sprintf("%s:%s", host, port), connectionHandler)
}
