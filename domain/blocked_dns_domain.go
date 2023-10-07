package domain

import "gorm.io/gorm"

type BlockedDNSDomain struct {
	gorm.Model
	Domain   string
	IsActive bool
}
