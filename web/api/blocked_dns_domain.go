package api

import (
	"myternet/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddBlockedDomain struct {
	Domain string `json:"domain"`
}

func AddBlockedDNSDomain(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		domainCommand := AddBlockedDomain{}
		if err := c.BindJSON(&domainCommand); err != nil {
			c.Error(err)
			return
		}
		dom := domain.BlockedDNSDomain{
			Domain:   domainCommand.Domain,
			IsActive: true,
		}
		db.Create(&dom)
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}
