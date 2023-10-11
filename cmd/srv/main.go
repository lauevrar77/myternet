package main

import (
	"fmt"
	"log"
	"myternet/domain"
	"myternet/infrastructure/services/db"
	mdns "myternet/infrastructure/services/dns"
	"myternet/web/api"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
	"gorm.io/gorm"
)

func runWeb(db *gorm.DB) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/dns/blocklist/add", api.AddBlockedDNSDomain(db))
	r.Run()
}

func main() {
	db, err := db.Connect(domain.BlockedDNSDomain{})
	if err != nil {
		log.Fatalf("Failed to connect to database %s\n", err.Error())
	}
	go runWeb(db)
	fmt.Println(db)
	srv := &dns.Server{Addr: ":" + strconv.Itoa(8053), Net: "udp"}
	dnsHandler := mdns.NewDNSHandler(db)
	srv.Handler = &dnsHandler
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}
