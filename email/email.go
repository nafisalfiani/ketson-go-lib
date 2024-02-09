package email

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/nafisalfiani/ketson-go-lib/codes"
	"github.com/nafisalfiani/ketson-go-lib/errors"
	"github.com/nafisalfiani/ketson-go-lib/log"
	"gopkg.in/gomail.v2"
)

const (
	emailRawHeaderFrom = "From"
	emailRawHeaderTo   = "To"
	emailRawHeaderCc   = "Cc"
	emailRawHeaderBcc  = "Bcc"
	emailRawSubject    = "Subject"
)

type Interface interface {
	SendEmail(ctx context.Context, params SendEmailParams) error
}

type Config struct {
	Smtp SmtpConfig
}

type SmtpConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	TLSConfig struct {
		InsecureSkipVerify bool
	}
}

type email struct {
	dialer *gomail.Dialer
	config Config
	log    log.Interface
}

func Init(cfg Config, log log.Interface) Interface {
	dialer := gomail.NewDialer(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password)
	if cfg.Smtp.TLSConfig.InsecureSkipVerify {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &email{
		dialer: dialer,
		config: cfg,
		log:    log,
	}
}

func (e *email) SendEmail(ctx context.Context, param SendEmailParams) error {
	if param.BodyType == "" {
		param.BodyType = BodyContentTypePlain
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader(emailRawHeaderFrom, fmt.Sprintf("%s <%s>", param.SenderName, param.SenderEmail))
	mailer.SetHeader(emailRawHeaderTo, param.Recipients.ToEmails...)
	mailer.SetHeader(emailRawHeaderCc, param.Recipients.CCEmails...)
	mailer.SetHeader(emailRawHeaderBcc, param.Recipients.BCCEmails...)
	mailer.SetHeader(emailRawSubject, param.Subject)
	mailer.SetBody(param.BodyType, param.Body)
	for hk, hv := range param.Headers {
		mailer.SetHeader(hk, hv)
	}
	for i := range param.Attachments {
		mailer.Attach(param.Attachments[i])
	}

	if err := e.dialer.DialAndSend(mailer); err != nil {
		return errors.NewWithCode(codes.CodeSendEmailFailed, "failed to send email, with err: %v", err)
	}

	return nil
}
