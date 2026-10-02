package mail

import (
	"bufio"
	"net"
	netmail "net/mail"
	"strconv"
	"strings"
	"testing"

	"smtp-service/internal/config"
)

// fakeSMTP accepts one message on 127.0.0.1 and sends its DATA on the returned channel.
func fakeSMTP(t *testing.T) (string, uint16, <-chan string) {
	t.Helper()

	var err error
	var ln net.Listener
	ln, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	data := make(chan string, 1)
	go func() {
		var err error
		var conn net.Conn
		conn, err = ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		r := bufio.NewReader(conn)
		reply := func(s string) { conn.Write([]byte(s + "\r\n")) }
		reply("220 fake ESMTP")
		for {
			var line string
			line, err = r.ReadString('\n')
			if err != nil {
				return
			}
			cmd := strings.ToUpper(strings.TrimSpace(line))
			switch {
			case strings.HasPrefix(cmd, "EHLO"):
				reply("250-fake")
				reply("250 AUTH PLAIN")
			case strings.HasPrefix(cmd, "AUTH"):
				reply("235 ok")
			case cmd == "DATA":
				reply("354 go ahead")
				var sb strings.Builder
				for {
					var l string
					l, err = r.ReadString('\n')
					if err != nil {
						return
					}
					if l == ".\r\n" {
						break
					}
					sb.WriteString(l)
				}
				data <- sb.String()
				reply("250 queued")
			case cmd == "QUIT":
				reply("221 bye")
				return
			default:
				reply("250 ok")
			}
		}
	}()

	var host, portStr string
	host, portStr, err = net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("split address: %v", err)
	}
	var port int
	port, err = strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parse port: %v", err)
	}

	return host, uint16(port), data
}

func TestSendWritesDateAndMessageID(t *testing.T) {
	host, port, data := fakeSMTP(t)

	var cfg config.Config
	cfg.SMTP.Host = host
	cfg.SMTP.Port = port
	cfg.SMTP.Username = "user"
	cfg.SMTP.Password = "pass"
	cfg.SMTP.From = "noreply@example.com"

	var err error
	err = Send(&cfg, "rcpt@example.org", "Hi", "<p>Hello</p>")
	if err != nil {
		t.Fatalf("Send: %v", err)
	}

	var msg *netmail.Message
	msg, err = netmail.ReadMessage(strings.NewReader(<-data))
	if err != nil {
		t.Fatalf("parse message: %v", err)
	}

	_, err = msg.Header.Date()
	if err != nil {
		t.Errorf("Date header %q: %v", msg.Header.Get("Date"), err)
	}

	id := msg.Header.Get("Message-ID")
	if !strings.HasPrefix(id, "<") || !strings.HasSuffix(id, "@example.com>") {
		t.Errorf("Message-ID = %q, want <random@example.com>", id)
	}
}
