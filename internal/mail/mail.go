package mail

import (
	"crypto/rand"
	"errors"
	"fmt"
	netmail "net/mail"
	"net/smtp"
	"strings"
	"time"

	"smtp-service/internal/config"
)

// ErrInvalidInput means the recipient or subject is empty or holds a line break.
var ErrInvalidInput = errors.New("invalid recipient or subject")

func Send(cfg *config.Config, to string, subject string, body string) error {
	if to == "" || strings.ContainsAny(to+subject, "\r\n") {
		return ErrInvalidInput
	}

	addr := fmt.Sprintf("%s:%d", cfg.SMTP.Host, cfg.SMTP.Port)
	auth := smtp.PlainAuth("", cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Host)
	msg := buildMessage(cfg.SMTP.From, to, subject, body, time.Now())

	return smtp.SendMail(addr, auth, cfg.SMTP.From, []string{to}, []byte(msg))
}

func buildMessage(from string, to string, subject string, body string, now time.Time) string {
	return "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Date: " + now.Format(time.RFC1123Z) + "\r\n" +
		"Message-ID: " + messageID(from) + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + body
}

// messageID returns a unique id in the form <random@sender-domain>.
func messageID(from string) string {
	domain := "localhost"
	parsed, err := netmail.ParseAddress(from)
	if err == nil {
		at := strings.LastIndex(parsed.Address, "@")
		if at >= 0 && at < len(parsed.Address)-1 {
			domain = parsed.Address[at+1:]
		}
	}

	return "<" + rand.Text() + "@" + domain + ">"
}
