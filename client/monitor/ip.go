package monitor

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"server-monitor/client/model"
	"strings"
)

var ipApis = []string{
	"https://api.ip.sb/geoip",
	"https://ipapi.co/json",
	"http://ip-api.com/json",
}

const v6Api = "http://v6.ipv6-test.com/api/myip.php"

// GetIP Get the external IP address of the local host.
func GetIP() (string, string, string) {
	for _, api := range ipApis {
		response, err := http.Get(api)
		if err != nil {
			continue
		}
		bytes, _ := ioutil.ReadAll(response.Body)
		var ipModel model.IPModel
		json.Unmarshal(bytes, &ipModel)

		ipVersion := "IPv4"
		isV6 := isIPv6()
		if isV6 {
			ipVersion = "IPv6"
		}

		// Adapt to the third IP API above.
		if len(ipModel.IP) == 0 {
			ipModel.IP = ipModel.IP_
			ipModel.CountryCode = ipModel.CountryCode_
		}

		return ipModel.IP, ipVersion, strings.ToUpper(ipModel.CountryCode)
	}
	return "", "", ""
}

func isIPv6() bool {
	_, err := http.Get(v6Api)
	return err == nil
}
