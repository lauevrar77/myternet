package dns

import (
	"fmt"
	"github.com/miekg/dns"
	"gorm.io/gorm"
	"myternet/domain"
)

type Handler struct {
	resolver    *DNSResolver
	queryStream chan DNSQuery
}

func NewDNSHandler(db *gorm.DB) Handler {
	queryStream := make(chan DNSQuery, 100)
	changeStream := make(chan []domain.BlockedDNSDomain, 0)
	updater := &dnsUpdater{
		db:             db,
		changeStream:   changeStream,
		blockedDomains: make([]domain.BlockedDNSDomain, 0),
	}
	go updater.Update()
	resolver := &DNSResolver{
		db:                db,
		changeStream:      changeStream,
		queryStream:       queryStream,
		aRecordsBlockList: make(map[string]string),
		updater:           updater,
	}
	go resolver.Run()
	return Handler{
		resolver:    resolver,
		queryStream: queryStream,
	}
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	go fmt.Printf("Received DNS query: %v\n", r)
	resultChan := make(chan DNSResponse, 0)
	h.queryStream <- DNSQuery{
		Query:      r,
		ResultChan: resultChan,
	}
	result := <-resultChan
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	w.WriteMsg(&result.Response)
}
