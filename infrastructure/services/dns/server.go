package dns

import (
	"net"

	"github.com/miekg/dns"
)

var domainsToAddresses map[string]string = map[string]string{
	"google.com.":       "1.2.3.4",
	"jameshfisher.com.": "104.198.14.52",
}

type Handler struct {
	internalResolver map[string]string
}

func NewDNSHandler(internalResolver map[string]string) Handler {
	return Handler{
		internalResolver: internalResolver,
	}
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, ok := h.internalResolver[domain]
		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})
		} else {
			answer, _ := Resolve(r)
			for _, a := range answer {
				msg.Answer = append(msg.Answer, a)
			}
		}
	default:
		answer, _ := Resolve(r)
		for _, a := range answer {
			msg.Answer = append(msg.Answer, a)
		}
	}
	w.WriteMsg(&msg)
}
