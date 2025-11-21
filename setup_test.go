package llm

import (
	"testing"

	"github.com/coredns/caddy"
)

func TestParse_Success(t *testing.T) {
	corefile := `llm {
		endpoint https://openrouter.ai/api/v1/chat/completions
		model openai/gpt-4o-mini
		chunk_size 128
		timeout 5
	}`
	c := caddy.NewTestController("dns", corefile)
	cfg, err := parse(c)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if cfg.EndpointURL == "" || cfg.Model == "" || cfg.ChunkSize != 128 || cfg.TimeoutSeconds != 5 {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
}

func TestParse_Errors(t *testing.T) {
	// missing endpoint
	c1 := caddy.NewTestController("dns", `llm { model m }`)
	if _, err := parse(c1); err == nil {
		t.Fatalf("expected error for missing endpoint")
	}

	// missing model
	c2 := caddy.NewTestController("dns", `llm { endpoint https://x }`)
	if _, err := parse(c2); err == nil {
		t.Fatalf("expected error for missing model")
	}

	// unknown property
	c3 := caddy.NewTestController("dns", `llm { foo bar }`)
	if _, err := parse(c3); err == nil {
		t.Fatalf("expected error for unknown property")
	}

	// invalid numbers
	c4 := caddy.NewTestController("dns", `llm { endpoint https://x model m chunk_size -1 }`)
	if _, err := parse(c4); err == nil {
		t.Fatalf("expected error for bad chunk_size")
	}
	c5 := caddy.NewTestController("dns", `llm { endpoint https://x model m timeout 0 }`)
	if _, err := parse(c5); err == nil {
		t.Fatalf("expected error for bad timeout")
	}
}
