package config

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

var amqpConn *amqp.Connection

var amqpHost = "localhost:5672"

var amqpuser = "guest"
var amqppassword = "guest"

func getAmqpConnectionString() string {
	if envhost := os.Getenv("AMQP_HOST"); envhost != "" {
		amqpHost = envhost
	}

	if envuser := os.Getenv("AMQP_USER"); envuser != "" {
		amqpuser = envuser
	}

	if envpass := os.Getenv("AMQP_PASS"); envpass != "" {
		amqppassword = envpass
	}

	return "amqp://" + amqpuser + ":" + amqppassword + "@" + amqpHost + "/"
}

func amqpConnection() *amqp.Connection {
	if amqpConn == nil {
		connString := getAmqpConnectionString()
		log.Printf("Starting amqp connection on %s", connString)
		conn, err := amqp.Dial(connString)
		if err != nil {
			panic(err)
		}

		amqpConn = conn
	}

	return amqpConn
}
