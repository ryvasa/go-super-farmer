package mail_handler

import (
	"encoding/json"

	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	mail_usecase "github.com/ryvasa/go-super-farmer/service_mail/usecase"
)

type MailHandler interface {
	ConsumerHandler() error
}

type MailHandlerImpl struct {
	rabbitMQ messages.RabbitMQ
	usecase  mail_usecase.MailUsecase
}

type EmailMessage struct {
	To  string `json:"to"`
	OTP string `json:"otp"`
}

func NewMailHandler(rabbitMQ messages.RabbitMQ, usecase mail_usecase.MailUsecase) MailHandler {
	return &MailHandlerImpl{
		rabbitMQ: rabbitMQ,
		usecase:  usecase,
	}
}

func (h *MailHandlerImpl) ConsumerHandler() error {
	mail, err := h.rabbitMQ.ConsumeMessages("mail-queue")
	if err != nil {
		logrus.Log.Fatalf("failed to consume messages: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range mail {
			var emailMsg EmailMessage
			if err := json.Unmarshal(d.Body, &emailMsg); err != nil {
				logrus.Log.Errorf("failed to unmarshal message: %v", err)
				continue
			}

			if err := h.usecase.SendOTPEmail(emailMsg.To, emailMsg.OTP); err != nil {
				logrus.Log.Errorf("failed to send email: %v", err)
			} else {
				logrus.Log.Infof("Successfully sent OTP email to: %s", emailMsg.To)
			}
		}
	}()

	logrus.Log.Info("Consumer Handler Started")
	<-forever

	return nil
}
