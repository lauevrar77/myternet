package main

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
	"myternet/domain"
	"myternet/infrastructure/services/db"
	mdns "myternet/infrastructure/services/dns"
	"strconv"
)

func main() {
	db, err := db.Connect(domain.BlockedDNSDomain{})
	if err != nil {
		log.Fatalf("Failed to connect to database %s\n", err.Error())
	}
	fmt.Println(db)
	srv := &dns.Server{Addr: ":" + strconv.Itoa(53), Net: "udp"}
	srv.Handler = &mdns.Handler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}
