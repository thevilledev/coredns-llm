package llm

import (
	"context"
	"testing"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type mockLLM struct {
	reply string
	err   error
}

func (m *mockLLM) Chat(ctx context.Context, prompt string) (string, error) { return m.reply, m.err }

type nopWriter struct{ dns.ResponseWriter }

func (w nopWriter) WriteMsg(*dns.Msg) error { return nil }

func TestServeDNS_TXT(t *testing.T) {
	h := &Handler{ChunkSize: 200, Client: &mockLLM{reply: "ok"}}
	req := new(dns.Msg)
	req.SetQuestion("hello.world.", dns.TypeTXT)
	code, err := h.ServeDNS(context.Background(), nopWriter{}, req)
	if err != nil || code != dns.RcodeSuccess {
		t.Fatalf("ServeDNS failed: code=%d err=%v", code, err)
	}
}

func TestServeDNS_Passthrough(t *testing.T) {
	// NextOrFailure requires Next to be non-nil; set Next to a no-op handler returning NXDOMAIN
	h := &Handler{ChunkSize: 200, Client: &mockLLM{reply: "ok"}, Next: plugin.HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
		return dns.RcodeNameError, nil
	})}
	req := new(dns.Msg)
	req.SetQuestion("hello.world.", dns.TypeA)
	code, _ := h.ServeDNS(context.Background(), nopWriter{}, req)
	if code != dns.RcodeNameError {
		t.Fatalf("expected passthrough NXDOMAIN, got %d", code)
	}
}
