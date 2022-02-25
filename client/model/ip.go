package model

type IPModel struct {
	IP   string `json:"ip"`
	Info struct {
		Country string `json:"country"`
	}
}
