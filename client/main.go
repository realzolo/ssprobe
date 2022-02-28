package main

import (
	"encoding/json"
	"server-monitor/client/monitor"
	"server-monitor/common/model"
	"time"
)

var authRes *monitor.AuthResult

func main() {
	// Authenticate the client.
	authRes = monitor.RequestAuth()
	// Realtime data.
	monitor.GetRealtimeData()
	// Stay connected and push data.
	pushData()
}

func pushData() {
	for {
		_ip, _ipVersion, _location := monitor.GetIP()
		_os, _process, _uptime := monitor.GetHost()
		_memTotal, _memUsed, _memUsedPct := monitor.GetMemory()
		_swapMemTotal, _swapMemUsed, _swapMemUsedPct := monitor.GetSwapMemory()
		_hddTotal, _hddUsed, _hddUsedPct := monitor.GetHDDSize()
		_cpuCount, _cpuUsedPct := monitor.GetCPU()
		_load1, _load5, _load15 := monitor.GetLoad()
		_ping10000, _ := monitor.PingTime.Load("10000")
		_ping10010, _ := monitor.PingTime.Load("10010")
		_ping10086, _ := monitor.PingTime.Load("10086")
		_lostRate10000, _ := monitor.LostRate.Load("10000")
		_lostRate10010, _ := monitor.LostRate.Load("10010")
		_lostRate10086, _ := monitor.LostRate.Load("10086")
		osModel := model.OSModel{
			Name:           authRes.Name,
			Host:           _ip,
			IPVersion:      _ipVersion,
			State:          true,
			OS:             _os,
			Location:       _location,
			Uptime:         _uptime,
			Load1:          _load1,
			Load5:          _load5,
			Load15:         _load15,
			MemTotal:       _memTotal,
			MemUsed:        _memUsed,
			MemUsedPct:     _memUsedPct,
			SwapMemTotal:   _swapMemTotal,
			SwapMemUsed:    _swapMemUsed,
			SwapMemUsedPct: _swapMemUsedPct,
			HddTotal:       _hddTotal,
			HddUsed:        _hddUsed,
			HddUsedPct:     _hddUsedPct,
			CpuCount:       _cpuCount,
			CpuUsedPct:     _cpuUsedPct,
			NetDownSpeed:   monitor.NetInfo["byteRecv"],
			NetUpSpeed:     monitor.NetInfo["byteSent"],
			ByteRecvTotal:  monitor.NetInfo["byteTotalRecv"],
			ByteSentTotal:  monitor.NetInfo["byteTotalSent"],
			Ping10000:      _ping10000.(int),
			Ping10010:      _ping10010.(int),
			Ping10086:      _ping10086.(int),
			LostRate10000:  _lostRate10000.(int),
			LostRate10010:  _lostRate10010.(int),
			LostRate10086:  _lostRate10086.(int),
			Tcp:            0, // TODO
			Udp:            0, // TODO
			Process:        _process,
		}
		bytes, _ := json.Marshal(osModel)
		(*authRes.Conn).Write(bytes)
		time.Sleep(time.Second)
	}
}
