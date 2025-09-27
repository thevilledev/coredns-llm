package llm

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/miekg/dns"
)

// decodeQName converts a DNS QNAME into a user prompt by unescaping octal escapes
// and replacing dots between labels with spaces.
func decodeQName(qname string) string {
	qname = strings.TrimSuffix(qname, ".")
	var b strings.Builder
	for i := 0; i < len(qname); i++ {
		ch := qname[i]
		if ch == 0x1a { // some inputs (e.g., Go string "\032") become 0x1a; treat as space
			b.WriteByte(' ')
			continue
		}
		if ch == '\\' {
			if i+3 < len(qname) {
				d1, d2, d3 := qname[i+1], qname[i+2], qname[i+3]
				if isDigit(d1) && isDigit(d2) && isDigit(d3) {
					v := int(d1-'0')*100 + int(d2-'0')*10 + int(d3-'0')
					if v >= 0 && v <= 255 {
						b.WriteByte(byte(v))
						i += 3
						continue
					}
				}
			}
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

func isDigit(b byte) bool { return b >= '0' && b <= '9' }

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
		if l > limit {
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
