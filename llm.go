package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

// Config holds Corefile configuration for the plugin.
type Config struct {
	EndpointURL    string
	Model          string
	ChunkSize      int
	TimeoutSeconds int
}

func defaultConfig() Config {
	return Config{
		ChunkSize:      200,
		TimeoutSeconds: 15,
	}
}

// Handler implements the CoreDNS plugin.Handler interface.
type Handler struct {
	Next        plugin.Handler
	EndpointURL string
	Model       string
	ChunkSize   int
	APIKey      string
	HTTPClient  *http.Client
}

func (h *Handler) Name() string { return "llm" }

func (h *Handler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	if len(r.Question) == 0 {
		return dns.RcodeFormatError, nil
	}
	q := r.Question[0]
	if q.Qtype != dns.TypeTXT {
		return plugin.NextOrFailure(h.Name(), h.Next, ctx, w, r)
	}

	prompt := decodeQName(q.Name)
	if prompt == "" {
		return dns.RcodeNameError, nil
	}

	answer, err := h.queryLLM(ctx, prompt)
	if err != nil {
		answer = fmt.Sprintf("Error: %v", err)
	}

	reply := new(dns.Msg)
	reply.SetReply(r)
	addTXTChunks(reply, q.Name, answer, h.ChunkSize)

	if err := w.WriteMsg(reply); err != nil {
		return dns.RcodeServerFailure, err
	}
	return dns.RcodeSuccess, nil
}

// decodeQName converts a DNS QNAME into a user prompt by unescaping octal escapes
// and replacing dots between labels with spaces.
func decodeQName(qname string) string {
	qname = strings.TrimSuffix(qname, ".")
	// RFC1035 allows \DDD (octal). dig uses \032 for spaces, etc.
	var b strings.Builder
	for i := 0; i < len(qname); i++ {
		ch := qname[i]
		if ch == '\\' {
			// Try to parse up to 3 octal digits.
			if i+3 < len(qname) {
				o1, o2, o3 := qname[i+1], qname[i+2], qname[i+3]
				if isOctal(o1) && isOctal(o2) && isOctal(o3) {
					v := (int(o1-'0') << 6) | (int(o2-'0') << 3) | int(o3-'0')
					b.WriteByte(byte(v))
					i += 3
					continue
				}
			}
			// Not octal; treat next char as literal if present.
			if i+1 < len(qname) {
				i++
				b.WriteByte(qname[i])
				continue
			}
		}
		if ch == '.' {
			b.WriteByte(' ')
			continue
		}
		b.WriteByte(ch)
	}
	return strings.TrimSpace(b.String())
}

func isOctal(b byte) bool { return b >= '0' && b <= '7' }

func addTXTChunks(m *dns.Msg, qname, text string, limit int) {
	if limit <= 0 {
		limit = 200
	}
	chunks := splitUTF8ByByteLimit(text, limit)
	for i, c := range chunks {
		label := c
		if len(chunks) > 1 {
			label = fmt.Sprintf("[%d/%d] %s", i+1, len(chunks), c)
		}
		rr := &dns.TXT{Hdr: dns.RR_Header{Name: qname, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0}, Txt: []string{label}}
		m.Answer = append(m.Answer, rr)
	}
}

func splitUTF8ByByteLimit(s string, limit int) []string {
	var res []string
	var b strings.Builder
	for _, r := range s {
		l := utf8.RuneLen(r)
		if l > limit { // extremely rare, but guard
			continue
		}
		if b.Len()+l > limit {
			res = append(res, b.String())
			b.Reset()
		}
		b.WriteRune(r)
	}
	if b.Len() > 0 {
		res = append(res, b.String())
	}
	if len(res) == 0 {
		return []string{""}
	}
	return res
}

// OpenAI-compatible request/response
type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (h *Handler) queryLLM(ctx context.Context, prompt string) (string, error) {
	payload := chatRequest{
		Model:  h.Model,
		Stream: false,
		Messages: []chatMessage{{
			Role:    "user",
			Content: prompt,
		}},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.EndpointURL, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+h.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("upstream status %d", resp.StatusCode)
	}
	var cr chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return "", err
	}
	if len(cr.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	return cr.Choices[0].Message.Content, nil
}

// parsePositiveInt converts a string to a positive int.
func parsePositiveInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("must be a positive integer")
	}
	return n, nil
}
