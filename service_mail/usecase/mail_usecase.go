package mail_usecase

import (
	"fmt"
	"strconv"

	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"gopkg.in/gomail.v2"
)

type MailUsecase interface {
	SendOTPEmail(to string, otp string) error
}

type MailUsecaseImpl struct {
	env *env.Env
}

func NewMailUsecase(env *env.Env) MailUsecase {
	return &MailUsecaseImpl{env: env}
}

func (s *MailUsecaseImpl) SendOTPEmail(to string, otp string) error {
	logrus.Log.Infof("Sending OTP email to: %s", to)
	m := gomail.NewMessage()
	m.SetHeader("From", s.env.Email.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "OTP Verification")

	// HTML template untuk email
	htmlBody := fmt.Sprintf(`
        <html>
            <body>
                <h2>OTP Verification</h2>
                <p>Your OTP code is: <strong>%s</strong></p>
                <p>This code will expire in 5 minutes.</p>
                <p>If you didn't request this code, please ignore this email.</p>
            </body>
        </html>
    `, otp)

	m.SetBody("text/html", htmlBody)
	port, err := strconv.Atoi(s.env.SMTP.Port)
	if err != nil {
		return fmt.Errorf("failed to parse port: %v", err)
	}
	d := gomail.NewDialer(
		s.env.SMTP.Host,
		port,
		s.env.Email.From,
		s.env.Email.Password,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
