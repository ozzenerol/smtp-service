package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"smtp-service/internal/config"
	"smtp-service/internal/mail"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	http.HandleFunc("POST /send", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			To      string `json:"to"`
			Subject string `json:"subject"`
			Body    string `json:"body"`
		}

		if json.NewDecoder(r.Body).Decode(&req) != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if sendErr := mail.Send(cfg, req.To, req.Subject, req.Body); sendErr != nil {
			log.Printf("send failed: %v", sendErr)
			http.Error(w, "send failed", http.StatusBadGateway)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
