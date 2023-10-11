package dns

import (
	"fmt"
	"github.com/miekg/dns"
	"gorm.io/gorm"
	"myternet/domain"
	"net"
	"time"
)

type dnsChange struct {
	domain string
	add    bool
}

type dnsUpdater struct {
	db             *gorm.DB
	blockedDomains []domain.BlockedDNSDomain
	changeStream   chan []domain.BlockedDNSDomain
}

func (u *dnsUpdater) Update() {
	tick := time.Tick(time.Second)
	for range tick {
		u.update()
	}
}

func (u *dnsUpdater) update() {
	var blockedDomains []domain.BlockedDNSDomain
	u.db.Find(&blockedDomains)
	if len(blockedDomains) != len(u.blockedDomains) {
		u.changeStream <- blockedDomains
	} else {
		for _, blockedDomain := range blockedDomains {
			if alreadyBlocked := u.find(blockedDomain.Domain); alreadyBlocked == nil {
				u.changeStream <- blockedDomains
				break
			} else if alreadyBlocked.IsActive != blockedDomain.IsActive {
				u.changeStream <- blockedDomains
				break
			}
		}
	}
	u.blockedDomains = blockedDomains
}

func (u *dnsUpdater) find(domain string) *domain.BlockedDNSDomain {
	for _, blockedDomain := range u.blockedDomains {
		if blockedDomain.Domain == domain {
			return &blockedDomain
		}
	}
	return nil
}

type DNSResponse struct {
	Response dns.Msg
	Error    error
}

type DNSQuery struct {
	Query      *dns.Msg
	ResultChan chan DNSResponse
}

type DNSResolver struct {
	db                *gorm.DB
	changeStream      chan []domain.BlockedDNSDomain
	queryStream       chan DNSQuery
	aRecordsBlockList map[string]string
	updater           *dnsUpdater
}

func (r *DNSResolver) Run() {
	for {
		select {
		case change := <-r.changeStream:
			r.updateCache(change)
		case query := <-r.queryStream:
			result, err := r.resolve(query.Query)
			query.ResultChan <- DNSResponse{
				Response: result,
				Error:    err,
			}
			close(query.ResultChan)
		}
	}
}

func (r *DNSResolver) updateCache(blockedDomains []domain.BlockedDNSDomain) {
	r.aRecordsBlockList = make(map[string]string)
	for _, blockedDomain := range blockedDomains {
		if blockedDomain.IsActive {
			r.aRecordsBlockList[blockedDomain.Domain] = "0.0.0.0"
		}
	}
	fmt.Printf("Updated DNS cache with %v\n", r.aRecordsBlockList)
}

func (r *DNSResolver) resolve(m *dns.Msg) (dns.Msg, error) {
	msg := dns.Msg{}
	msg.SetReply(m)
	switch m.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, ok := r.aRecordsBlockList[domain]
		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})
		} else {
			answer, _ := r.resolveRemote(m)
			for _, a := range answer {
				msg.Answer = append(msg.Answer, a)
			}
		}
	default:
		answer, _ := r.resolveRemote(m)
		for _, a := range answer {
			msg.Answer = append(msg.Answer, a)
		}
	}
	return msg, nil
}

func (r *DNSResolver) resolveRemote(m *dns.Msg) ([]dns.RR, error) {
	c := new(dns.Client)
	in, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return in.Answer, nil
}
