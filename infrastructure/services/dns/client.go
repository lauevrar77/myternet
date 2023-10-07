package dns

import (
	"fmt"
	"github.com/miekg/dns"
)

func Resolve(m *dns.Msg) ([]dns.RR, error) {
	c := new(dns.Client)
	in, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return in.Answer, nil
}
