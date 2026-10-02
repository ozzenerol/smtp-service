package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"smtp-service/internal/config"
	"smtp-service/internal/mail"
)

const (
	testTo      = "customer@example.org"
	testSubject = "Your quote 4711"
	testBody    = "<p>Secret quote details</p>"
)

func TestSendHandlerLogsNoAddress(t *testing.T) {
	cases := []struct {
		name    string
		sendErr error
		status  int
	}{
		{"success", nil, http.StatusNoContent},
		{"smtp reply echoes address", &textproto.Error{Code: 550, Msg: "5.1.1 <" + testTo + ">: Recipient address rejected"}, http.StatusBadGateway},
		{"wrapped smtp reply", fmt.Errorf("rcpt: %w", &textproto.Error{Code: 452, Msg: "too many recipients for " + testTo}), http.StatusBadGateway},
		{"network error", &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}, http.StatusBadGateway},
		{"invalid input", mail.ErrInvalidInput, http.StatusBadGateway},
		{"unknown error with address", errors.New("cannot send to " + testTo + " about " + testSubject), http.StatusBadGateway},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var logs bytes.Buffer
			prev := log.Writer()
			log.SetOutput(&logs)
			defer log.SetOutput(prev)

			var gotTo string
			send := func(cfg *config.Config, to string, subject string, body string) error {
				gotTo = to
				return tc.sendErr
			}

			payload := `{"to":"` + testTo + `","subject":"` + testSubject + `","body":"` + testBody + `"}`
			req := httptest.NewRequest(http.MethodPost, "/send", strings.NewReader(payload))
			rec := httptest.NewRecorder()
			sendHandler(&config.Config{}, send)(rec, req)

			if rec.Code != tc.status {
				t.Errorf("status = %d, want %d", rec.Code, tc.status)
			}
			if gotTo != testTo {
				t.Errorf("send got to = %q, want %q", gotTo, testTo)
			}

			out := logs.String()
			if out == "" {
				t.Fatal("no log line written")
			}
			if !strings.Contains(out, "1 recipient") {
				t.Errorf("log has no recipient count: %q", out)
			}
			if !strings.Contains(out, fmt.Sprintf("status %d", tc.status)) {
				t.Errorf("log has no status: %q", out)
			}
			for _, secret := range []string{"@", testTo, testSubject, testBody} {
				if strings.Contains(out, secret) {
					t.Errorf("log contains %q: %q", secret, out)
				}
			}
		})
	}
}
