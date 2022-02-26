package model

type OSModel struct {
	Name           string  `json:"name"`
	Host           string  `json:"host"`
	IPVersion      string  `json:"ip_version"`
	State          bool    `json:"state"`
	OS             string  `json:"os"`
	Location       string  `json:"location"`
	Uptime         uint64  `json:"uptime"` // Unit: second
	Load1          float64 `json:"load_1"`
	Load5          float64 `json:"load_5"`
	Load15         float64 `json:"load_15"`
	MemTotal       uint64  `json:"mem_total"` // Unit: B(Byte)
	MemUsed        uint64  `json:"mem_used"`  // Unit: B(Byte)
	MemUsedPct     uint64  `json:"mem_used_pct"`
	SwapMemTotal   uint64  `json:"swap_mem_total"` // Unit: B(Byte)
	SwapMemUsed    uint64  `json:"swap_mem_used"`  // Unit: B(Byte)
	SwapMemUsedPct uint64  `json:"swap_mem_used_pct"`
	HddTotal       uint64  `json:"hdd_total"` // Unit: B(Byte)
	HddUsed        uint64  `json:"hdd_used"`  // Unit: B(Byte)
	HddUsedPct     uint64  `json:"hdd_used_pct"`
	CpuCount       int     `json:"cpu_count"`
	CpuUsedPct     uint64  `json:"cpu_used_pct"`
	NetDownSpeed   uint64  `json:"net_down_speed"`  // Unit: B/s
	NetUpSpeed     uint64  `json:"net_up_speed"`    // Unit: B/s
	ByteRecvTotal  uint64  `json:"byte_recv_total"` // Unit: B(Byte)
	ByteSentTotal  uint64  `json:"byte_sent_total"` // Unit: B(Byte)
	Ping10000      int     `json:"ping_10000"`      // Unit: ms
	Ping10010      int     `json:"ping_10010"`      // Unit: ms
	Ping10086      int     `json:"ping_10086"`      // Unit: ms
	LostRate10000  int     `json:"lost_rate_10000"`
	LostRate10010  int     `json:"lost_rate_10010"`
	LostRate10086  int     `json:"lost_rate_10086"`
	Tcp            uint    `json:"tcp"`
	Udp            uint    `json:"udp"`
	Process        uint64  `json:"process"`
}
