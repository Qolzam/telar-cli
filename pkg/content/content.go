package content

import (
	"log"
	"net/url"
	"strings"
)

func GetDomainFromURI(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}
	parts := strings.Split(u.Hostname(), ".")
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
	return domain
}
