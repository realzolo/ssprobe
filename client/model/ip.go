package model

type IPModel struct {
	IP           string `json:"ip"`
	CountryCode  string `json:"country_code"`
	IP_          string `json:"query,omitempty"`
	CountryCode_ string `json:"countryCode,omitempty"`
}
