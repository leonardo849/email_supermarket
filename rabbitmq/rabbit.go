package rabbitmq

import (
	"github.com/leonardo849/email_supermarket/logger"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"strings"

	"github.com/leonardo849/utils_for_backend/pkg/email_dto"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var isRabbitMQon = false
const key = "email"

type clientI interface {
	sendEmail(input email_dto.SendEmailDTO) error
	consumerEmail()
	loadEnvVars() error
}

type client struct {
	conn *amqp.Connection
	ch *amqp.Channel
	from string
	password string
	smtpPort string
	smtpHost string
	auth smtp.Auth
}



var rabbitClient *client

func ConnectToRabbitMQ() error {
	isRabbitOn := strings.ToLower(os.Getenv("RABBIT_ON"))
	if isRabbitOn == "true" {
		isRabbitMQon = true
	}
	logger.ZapLogger.Info("Is rabbit going to be on?", zap.Bool("rabbit value", isRabbitMQon))
	if !isRabbitMQon {
		startConsumer()
		return fmt.Errorf("rabbit mq is going to be off")
	}

	uriRabbit := os.Getenv("RABBIT_URI")
	if uriRabbit == "" {
		err := fmt.Errorf("rabbit_uri is empty")
		logger.ZapLogger.Error("error in get rabbit uri", zap.Error(err))
		startConsumer()
		return err
	}

	var conn *amqp.Connection
	const maxTries = 11
	secondDelay := os.Getenv("SECOND_DELAY")
	secondInt := 1
	var err error
	if secondDelay != "" {
		secondInt, err = strconv.Atoi(secondDelay)
		if err != nil {
			secondInt = 1
		}
	}
	for i := 0; i < maxTries; i++ {
		conn, err = amqp.Dial(uriRabbit)
		if err != nil {
			if i == maxTries-1 {
				startConsumer()
				return err
			}
			logger.ZapLogger.Error("error in Connect To rabbit mq", zap.Error(err))
			time.Sleep(time.Duration(secondInt) * time.Second)
		} else {
			break
		}
	}

	logger.ZapLogger.Info("conn is estabilished")
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		logger.ZapLogger.Info("closing conn")
		logger.ZapLogger.Error("error in get channel", zap.Error(err))
		startConsumer()
		return err
	}
	logger.ZapLogger.Info("channel is ready")

	rabbitClient = &client{
		conn: conn,
		ch:   ch,
		smtpHost: "smtp.gmail.com",
		smtpPort: "587",
	}
	if err := rabbitClient.loadEnvVars(); err != nil {
		os.Exit(1)
	}
	startConsumer()
	return nil
}

func startConsumer() {
	c := getRabbitMQClient()
	go c.consumerEmail()
	select {}
}


func getRabbitMQClient() clientI {
	if isRabbitMQon && rabbitClient != nil {
		return  rabbitClient
	} else {
		return  &fakeClient{}
	}
}