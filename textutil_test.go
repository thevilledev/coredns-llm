package llm

import (
	"strings"
	"testing"

	"github.com/miekg/dns"
)

func TestDecodeQName(t *testing.T) {
	cases := map[string]string{
		"hello.world.":          "hello world",
		"what\\032is\\032up.":   "what is up",
		"slash\\046dot.":        "slash.dot",
		"octal\\049\\050\\051.": "octal123",
	}
	for in, want := range cases {
		got := decodeQName(in)
		if got != want {
			t.Fatalf("%q => %q; want %q", in, got, want)
		}
	}
}

func TestAddTXTChunks(t *testing.T) {
	m := new(dns.Msg)
	var text strings.Builder
	for range 500 {
		text.WriteString("a")
	}
	addTXTChunks(m, "example.", text.String(), 200)
	if len(m.Answer) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(m.Answer))
	}
	if rr, ok := m.Answer[0].(*dns.TXT); !ok || rr.Hdr.Name != "example." {
		t.Fatalf("unexpected first rr: %#v", m.Answer[0])
	}
}
