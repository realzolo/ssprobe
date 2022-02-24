package model

type AuthModel struct {
	Code      int8   `json:"code"`
	Host      string `json:"host"`
	IpVersion string `json:"ip_version"`
}
