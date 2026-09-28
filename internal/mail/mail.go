package mail

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"

	"smtp-service/internal/config"
)

func Send(cfg *config.Config, to string, subject string, body string) error {
	if to == "" || strings.ContainsAny(to+subject, "\r\n") {
		return errors.New("invalid recipient or subject")
	}

	addr := fmt.Sprintf("%s:%d", cfg.SMTP.Host, cfg.SMTP.Port)
	auth := smtp.PlainAuth("", cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Host)

	msg := "From: " + cfg.SMTP.From + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + body

	return smtp.SendMail(addr, auth, cfg.SMTP.From, []string{to}, []byte(msg))
}
