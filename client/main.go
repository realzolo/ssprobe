package main

import (
	"context"
	"encoding/json"
	"net"
	"ssprobe-client/monitor"
	"ssprobe-common/model"
	"time"
)

func main() {
	for {
		ctx, cancelFunc := context.WithCancel(context.Background())
		// Authenticate the client.
		authRes := monitor.RequestAuth()
		if !authRes.Ok {
			time.Sleep(time.Second * 60)
			continue
		}
		// Realtime data.
		monitor.GetRealtimeData(ctx)
		// Stay connected and push data.
		pushDataToServer(authRes.Name, authRes.Conn)
		// End other goroutines.
		cancelFunc()
		time.Sleep(time.Second * 60)
	}
}

func pushDataToServer(name string, conn *net.Conn) {
	var maxRetries = 10
	for {
		_ip, _ipVersion, _location := monitor.GetIP()
		_platform, _process, _uptime := monitor.GetHost()
		_memTotal, _memUsed, _memUsedPct := monitor.GetMemory()
		_swapMemTotal, _swapMemUsed, _swapMemUsedPct := monitor.GetSwapMemory()
		_hddTotal, _hddUsed, _hddUsedPct := monitor.GetHDDSize()
		_cpuCount, _cpuUsedPct := monitor.GetCPU()
		_load1, _load5, _load15 := monitor.GetLoad()

		_netDownSpeed := monitor.NetInfo["byteRecv"]
		_nettUpSpeed := monitor.NetInfo["byteSent"]
		_byteRecvTotal := monitor.NetInfo["byteTotalRecv"]
		_byteSentTotal := monitor.NetInfo["byteTotalSent"]

		_ping10000, _ := monitor.PingTime.Load("10000")
		_ping10010, _ := monitor.PingTime.Load("10010")
		_ping10086, _ := monitor.PingTime.Load("10086")
		_lostRate10000, _ := monitor.LostRate.Load("10000")
		_lostRate10010, _ := monitor.LostRate.Load("10010")
		_lostRate10086, _ := monitor.LostRate.Load("10086")

		osModel := model.OSModel{
			Name:           name,
			Host:           _ip,
			IPVersion:      _ipVersion,
			State:          true,
			Platform:       _platform,
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
			NetDownSpeed:   _netDownSpeed,
			NetUpSpeed:     _nettUpSpeed,
			ByteRecvTotal:  _byteRecvTotal,
			ByteSentTotal:  _byteSentTotal,
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

		// Try up to 10 times more.
		if _, err := (*conn).Write(bytes); err != nil {
			if maxRetries--; maxRetries == 0 {
				return
			}
		}

		time.Sleep(time.Second * 2)
	}
}
