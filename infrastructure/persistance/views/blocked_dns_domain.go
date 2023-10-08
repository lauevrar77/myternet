package views

import (
	"fmt"
	"myternet/domain"

	"gorm.io/gorm"
)

func BlockedDNSDomainToInternalResolver(db *gorm.DB) map[string]string {
	var domains []domain.BlockedDNSDomain
	db.Find(&domains)
	internalResolver := make(map[string]string)
	for _, dom := range domains {
		if dom.IsActive {
			internalResolver[dom.Domain] = "0.0.0.0"
		}
	}
	fmt.Println(internalResolver)
	return internalResolver
}
