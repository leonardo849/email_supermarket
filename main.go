package main

import (
	"log"

	"github.com/leonardo849/email_supermarket/config"
	"github.com/leonardo849/email_supermarket/logger"
	"github.com/leonardo849/email_supermarket/rabbitmq"
)

func main() {
	if err := config.SetupEnvVar(); err != nil {
		log.Panic(err.Error())
	}
	if err := logger.StartLogger(); err != nil {
		logger.ZapLogger.Fatal(err.Error())
	}
	if err := rabbitmq.ConnectToRabbitMQ(); err != nil {
		logger.ZapLogger.Error(err.Error())
	}

}