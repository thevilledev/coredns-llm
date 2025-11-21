package llm

import (
	"context"
	"fmt"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

// LLMClient defines the interface for querying an LLM.
type LLMClient interface {
	Chat(ctx context.Context, prompt string) (string, error)
}

// Handler implements the CoreDNS plugin.Handler interface.
type Handler struct {
	Next      plugin.Handler
	ChunkSize int
	Client    LLMClient
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

	answer, err := h.Client.Chat(ctx, prompt)
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
