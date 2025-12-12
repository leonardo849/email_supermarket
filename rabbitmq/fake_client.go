package rabbitmq

import (
	"github.com/leonardo849/utils_for_backend/pkg/email_dto"
	"github.com/leonardo849/email_supermarket/logger"
	"strings"
	"time"
)

type fakeClient struct {
}

func (c *fakeClient) loadEnvVars()  error {
	logger.ZapLogger.Info("[fake] loading env vars")
	return  nil
}

func (c *fakeClient) sendEmail(input email_dto.SendEmailDTO) error {
	logger.ZapLogger.Info("[fake] " + "sending email to " + strings.Join(input.To, "") + " subject " + input.Subject)
	return  nil
}

func (c *fakeClient) consumerEmail() {
	msg := make(chan email_dto.SendEmailDTO)
	go func() {
		for email := range msg {
			err := c.sendEmail(email)
			if err != nil {

			}
		}
	}()
	go func() {
		for {
			msg <- email_dto.SendEmailDTO{
				To:      []string{"user@example.com"},
				Subject: "test",
				Text:    "test message",
			}
			time.Sleep(time.Second * 10)
			
		} 
	}()
	for range msg {}
}