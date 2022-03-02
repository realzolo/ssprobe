package monitor

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"ssprobe-client/model"
	"strings"
)

var ipApis = []string{
	"https://api.ip.sb/geoip",
	"https://ipapi.co/json",
	"http://ip-api.com/json",
}

const v6Api = "http://v6.ipv6-test.com/api/myip.php"

var cache = struct {
	hasCached       bool
	cachedIP        string
	cachedIPVersion string
	cachedLocation  string
}{}

// GetIP Get the external IP address of the local host.
func GetIP() (string, string, string) {
	if cache.hasCached {
		return cache.cachedIP, cache.cachedIPVersion, cache.cachedLocation
	}
	for _, api := range ipApis {
		response, err := http.Get(api)
		if err != nil {
			continue
		}
		var ipModel model.IPModel
		bytes, _ := ioutil.ReadAll(response.Body)
		_ = json.Unmarshal(bytes, &ipModel)

		ipVersion := "IPv4"
		if isIPv6() {
			ipVersion = "IPv6"
		}

		// Adapt to the third API above.
		if len([]rune(ipModel.IP)) == 0 {
			if len([]rune(ipModel.IP_)) == 0 {
				continue
			}
			ipModel.IP = ipModel.IP_
			ipModel.CountryCode = ipModel.CountryCode_
		}

		cache.hasCached = true
		cache.cachedIP = ipModel.IP
		cache.cachedIPVersion = ipVersion
		cache.cachedLocation = strings.ToUpper(ipModel.CountryCode)
		return cache.cachedIP, cache.cachedIPVersion, cache.cachedLocation
	}
	return "", "", ""
}

func isIPv6() bool {
	_, err := http.Get(v6Api)
	return err == nil
}
