package model

type OSModel struct {
	Name        string  `json:"name"`
	Host        string  `json:"host"`
	State       bool    `json:"state"`
	Platform    string  `json:"platform"`
	Location    string  `json:"location"`
	Uptime      uint    `json:"uptime"`
	Online4     bool    `json:"online4"`
	Online6     bool    `json:"online6"`
	Load1       float64 `json:"load_1"`
	Load5       float64 `json:"load_5"`
	Load15      float64 `json:"load_15"`
	MemoryTotal uint    `json:"memory_total"`
	MemoryUsed  uint    `json:"memory_used"`
	SwapTotal   uint    `json:"swap_total"`
	SwapUsed    uint    `json:"swap_used"`
	HddTotal    uint    `json:"hdd_total"`
	HddUsed     uint    `json:"hdd_used"`
	Cpu         float32 `json:"cpu"`
	NetworkRx   uint    `json:"network_rx"`
	NetworkTx   uint    `json:"network_tx"`
	NetworkIn   uint    `json:"network_in"`
	NetworkOut  uint    `json:"network_out"`
	Ping10010   float32 `json:"ping_10010"`
	Ping189     float32 `json:"ping_189"`
	Ping10086   float32 `json:"ping_10086"`
	Time10010   uint    `json:"time_10010"`
	Time189     uint    `json:"time_189"`
	Time10086   uint    `json:"time_10086"`
	Tcp         uint    `json:"tcp"`
	Udp         uint    `json:"udp"`
	Process     uint    `json:"process"`
	Thread      uint    `json:"thread"`
}
