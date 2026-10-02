package mail

import (
	netmail "net/mail"
	"strings"
	"testing"
	"time"
)

func TestBuildMessageSetsDateAndMessageID(t *testing.T) {
	now := time.Date(2026, 10, 2, 9, 30, 0, 0, time.UTC)
	raw := buildMessage("Sender <noreply@example.com>", "rcpt@example.org", "Hi", "<p>Hello</p>", now)

	var err error
	var msg *netmail.Message
	msg, err = netmail.ReadMessage(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("parse message: %v", err)
	}

	var date time.Time
	date, err = msg.Header.Date()
	if err != nil {
		t.Fatalf("parse Date header %q: %v", msg.Header.Get("Date"), err)
	}
	if !date.Equal(now) {
		t.Errorf("Date = %v, want %v", date, now)
	}

	id := msg.Header.Get("Message-ID")
	if !strings.HasPrefix(id, "<") || !strings.HasSuffix(id, "@example.com>") {
		t.Errorf("Message-ID = %q, want <random@example.com>", id)
	}
	if id == "<@example.com>" {
		t.Errorf("Message-ID has no random part: %q", id)
	}
}

func TestMessageIDIsUnique(t *testing.T) {
	seen := make(map[string]bool)
	for range 1000 {
		id := messageID("noreply@example.com")
		if seen[id] {
			t.Fatalf("duplicate Message-ID %q", id)
		}
		seen[id] = true
	}
}

func TestMessageIDFallsBackWhenFromIsInvalid(t *testing.T) {
	id := messageID("not an address")
	if !strings.HasSuffix(id, "@localhost>") {
		t.Errorf("Message-ID = %q, want suffix @localhost>", id)
	}
}
