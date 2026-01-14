package provider

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
)

type SpamAssassin struct {
	addr    string
	timeout time.Duration
	maxsize int
}

func (p *SpamAssassin) Name() string {
	return "spamassassin"
}

func (p *SpamAssassin) ValidateConfig(config map[string]string) error {
	// host and port optional; set defaults
	host := config["host"]
	if host == "" {
		host = "127.0.0.1"
	}
	port := config["port"]
	if port == "" {
		port = "783"
	}
	p.addr = net.JoinHostPort(host, port)

	// timeout optional (seconds)
	if config["timeout"] == "" {
		p.timeout = 5 * time.Second
	} else {
		t, err := strconv.ParseFloat(config["timeout"], 64)
		if err != nil || t <= 0 {
			return errors.New("spamassassin timeout must be a positive number (seconds)")
		}
		p.timeout = time.Duration(t * float64(time.Second))
	}

	// maxsize required (positive integer)
	n, err := strconv.ParseInt(config["maxsize"], 10, 64)
	if err != nil || n < 1 {
		return errors.New("spamassassin maxsize must be a positive integer")
	}
	p.maxsize = int(n)

	return nil
}

func (p *SpamAssassin) Init(config map[string]string) error {
	if err := p.ValidateConfig(config); err != nil {
		return err
	}
	// nothing else to init for TCP spamd client
	return nil
}

func (p *SpamAssassin) Analyze(msg imap.Message) (int, error) {
	// Prefer sending the original raw message if available.
	var rawBytes []byte
	if len(msg.Raw) > p.maxsize {
		logx.Debugf("spamassassin: truncating raw message for message #%d (%s)", msg.UID, msg.Subject)
		rawBytes = msg.Raw[:p.maxsize]
	} else {
		rawBytes = msg.Raw
	}

	// connect to spamd
	conn, err := net.DialTimeout("tcp", p.addr, p.timeout)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = conn.Close()
	}()

	// build REPORT request per SPAMC protocol
	reqHeader := fmt.Sprintf("REPORT SPAMC/1.5\r\nContent-length: %d\r\n\r\n", len(rawBytes))
	if _, err := conn.Write([]byte(reqHeader)); err != nil {
		return 0, err
	}
	if _, err := conn.Write(rawBytes); err != nil {
		return 0, err
	}

	// read response (headers + body). Use reader and read all available.
	if err := conn.SetReadDeadline(time.Now().Add(p.timeout)); err != nil {
		return 0, err
	}
	reader := bufio.NewReader(conn)

	// read status line
	if _, err = reader.ReadString('\n'); err != nil {
		return 0, err
	}

	// read headers until blank line, capture Content-length if present
	contentLen := -1
	score := ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				if v, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
					contentLen = v
				}
			}
		} else if strings.HasPrefix(strings.ToLower(line), "spam:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				score = strings.TrimSpace(parts[1])
			}
		}
	}

	// read body
	var body string
	if contentLen >= 0 {
		buf := make([]byte, contentLen)
		if _, err = io.ReadFull(reader, buf); err != nil {
			return 0, err
		}
		body = string(buf)
	} else {
		// read till EOF / timeout; distinguish EOF/timeout from other errors
		var sb strings.Builder
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					// finished reading
					break
				}
				var ne net.Error
				if errors.As(err, &ne) && ne.Timeout() {
					// timeout is acceptable as end of response
					break
				}
				// other errors are real failures
				return 0, err
			}
			sb.WriteString(line)
		}
		body = sb.String()
	}

	logx.Debugf("spamassassin score for message #%d (%s): Spam: %s", msg.UID, msg.Subject, score)
	logx.Debugf("spamassassin response for message #%d (%s):\n%s", msg.UID, msg.Subject, body)

	// try to extract score from score header:
	// 1) True ; 7.8 / 5.0
	// 2) False ; 0.0 / 5.0
	reScore := regexp.MustCompile(`(True|False) ; ([+-]?[0-9]*\.?[0-9]+) / ([+-]?[0-9]*\.?[0-9]+)`)
	var scoreF float64
	found := false

	if m := reScore.FindStringSubmatch(score); len(m) == 4 {
		if v, err := strconv.ParseFloat(m[2], 64); err == nil {
			scoreF = v
			found = true
		}
	}

	if !found {
		return 0, errors.New("could not parse spamassassin score from response")
	}

	// map spamassassin score (float, typically -2..10) to 0..100
	const minScore = -2.0
	const maxScore = 10.0
	norm := (scoreF - minScore) / (maxScore - minScore) * 100.0
	if norm < 0 {
		norm = 0
	}
	if norm > 100 {
		norm = 100
	}
	intScore := int(norm)

	logx.Debugf("spamassassin raw score for message #%d (%s): %v -> %d/100", msg.UID, msg.Subject, scoreF, intScore)
	return intScore, nil
}
