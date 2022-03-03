package util

import "strings"

//ParseAddress Get the IP address or domain name, for example: 127.0.0.1:8000 -> 127.0.0.1
func ParseAddress(domain string) string {
	return domain[:strings.LastIndex(domain, ":")]
}
