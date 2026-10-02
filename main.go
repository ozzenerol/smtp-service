package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/textproto"
	"smtp-service/internal/config"
	"smtp-service/internal/mail"
)

type sendFunc func(cfg *config.Config, to string, subject string, body string) error

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	http.HandleFunc("POST /send", sendHandler(cfg, mail.Send))

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// sendHandler never logs the address, subject or body. It logs the recipient count and the outcome.
func sendHandler(cfg *config.Config, send sendFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			To      string `json:"to"`
			Subject string `json:"subject"`
			Body    string `json:"body"`
		}

		var err error
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		err = send(cfg, req.To, req.Subject, req.Body)
		if err != nil {
			log.Printf("Email send failed (1 recipient), status %d: %s", http.StatusBadGateway, safeError(err))
			http.Error(w, "send failed", http.StatusBadGateway)
			return
		}

		log.Printf("Email sent (1 recipient), status %d", http.StatusNoContent)

		w.WriteHeader(http.StatusNoContent)
	}
}

// safeError describes err without any address.
// SMTP servers often echo the recipient address in their replies, so only the reply code is kept.
// Other errors are kept only when they cannot hold an address.
func safeError(err error) string {
	var smtpErr *textproto.Error
	if errors.As(err, &smtpErr) {
		return fmt.Sprintf("smtp reply %d", smtpErr.Code)
	}

	var netErr net.Error
	if errors.As(err, &netErr) || errors.Is(err, mail.ErrInvalidInput) {
		return err.Error()
	}

	return "details withheld"
}
