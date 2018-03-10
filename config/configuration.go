package config

import "github.com/BurntSushi/toml"

type configProperties struct {
	Server         serverProperties
	KafkaConsumers map[string]consumerProperties `toml:"kafka_consumers"`
	KafkaProducers map[string]producerProperties `toml:"kafka_producers"`
}

type serverProperties struct {
	Host string
}

type consumerProperties struct {
	Broker          string
	GroupID         string `toml:"group.id"`
	AutoOffsetReset string `toml:"auto.offset.reset"`
}

type producerProperties struct {
	Broker string
}

var configuration configProperties

func init() {
	if _, err := toml.DecodeFile("config.toml", &configuration); err != nil {
		panic(err)
	}
}
