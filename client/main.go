package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	psnet "github.com/shirou/gopsutil/v3/net"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime"
	"server-monitor-client/model"
	"strconv"
	"strings"
	"sync"
	"time"
)

var netInfo = map[string]uint64{
	"byteSent":      0,
	"byteRecv":      0,
	"byteTotalRecv": 0,
	"byteTotalSent": 0,
	"clock":         0,
	"diff":          0,
}

var (
	name   string = "onezol.com"
	server string = "127.0.0.1"
	port   int    = 3384
	token  string = "123456"
)
var (
	pingTime sync.Map
	lostRate sync.Map
)
var conn *net.Conn

func init() {
	pingTime.Store("10000", 0)
	pingTime.Store("10010", 0)
	pingTime.Store("10086", 0)
	lostRate.Store("10000", 0)
	lostRate.Store("10010", 0)
	lostRate.Store("10086", 0)
}

func main() {
	// Collect config.
	collectConfig()
	// Connect to server.
	requestAuth()
	// Realtime data.
	GetRealtimeData()
	// Stay connected and push data.
	keepConnAndPushData()
}

// collectConfig Collect client configuration information.
func collectConfig() {
	fmt.Print("Client Name: ")
	fmt.Scanln(&name)
	fmt.Print("Server ip: ")
	fmt.Scanln(&server)
	fmt.Print("Server port: ")
	fmt.Scanln(&port)
	fmt.Print("Token: ")
	fmt.Scanln(&token)
	if len(name) == 0 || len(server) == 0 || len(token) == 0 {
		fmt.Println("Data cannot be null, Please re-run this program!")
	}
}

// connectToServer Use the socket connect to the server.
func requestAuth() {
	_conn, err := net.Dial("tcp", server+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("Failed to connect to server.", err)
	}
	// Authentication.
	bytes, _ := json.Marshal(token)
	_, err = _conn.Write(bytes)
	if err != nil {
		log.Fatal("Authentication failed! ", err)
	}
	var buf = make([]byte, 1024)
	n, err := _conn.Read(buf)
	if err != nil {
		log.Fatal("Authentication failed! ", err)
	}
	var resModel = struct {
		Code int `json:"code"`
	}{}
	json.Unmarshal(buf[:n], &resModel)
	// The token is incorrect.
	if resModel.Code == -1 {
		log.Fatal("Client authentication failed, token is incorrect!")
	}
	conn = &_conn
	log.Println("Server connection successful!")
}

func keepConnAndPushData() {
	_ip, _ipVersion, _location := GetIP()
	for {
		_os, _process, _uptime := GetHost()
		_memTotal, _memUsed, _memUsedPct := GetMemory()
		_swapMemTotal, _swapMemUsed, _swapMemUsedPct := GetSwapMemory()
		_hddTotal, _hddUsed, _hddUsedPct := GetHDDSize()
		_cpuCount, _cpuUsedPct := GetCPU()
		_load1, _load5, _load15 := GetLoad()
		_ping10000, _ := pingTime.Load("10000")
		_ping10010, _ := pingTime.Load("10010")
		_ping10086, _ := pingTime.Load("10086")
		_lostRate10000, _ := lostRate.Load("10000")
		_lostRate10010, _ := lostRate.Load("10010")
		_lostRate10086, _ := lostRate.Load("10086")
		osModel := model.OSModel{
			Name:           name,
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
			NetDownSpeed:   netInfo["byteRecv"],
			NetUpSpeed:     netInfo["byteSent"],
			ByteRecvTotal:  netInfo["byteTotalRecv"],
			ByteSentTotal:  netInfo["byteTotalSent"],
			Ping10000:      _ping10000.(int),
			Ping10010:      _ping10010.(int),
			Ping10086:      _ping10086.(int),
			LostRate10000:  _lostRate10000.(int),
			LostRate10010:  _lostRate10010.(int),
			LostRate10086:  _lostRate10086.(int),
			Tcp:            0,
			Udp:            0,
			Process:        _process,
		}
		bytes, _ := json.Marshal(osModel)
		fmt.Println(len(bytes))
		(*conn).Write(bytes)
		time.Sleep(time.Second)
	}
}

// GetMemory Get the usage of memory
func GetMemory() (uint64, uint64, uint64) {
	m, _ := mem.VirtualMemory()
	return m.Total, m.Used, uint64(m.UsedPercent)
}

// GetSwapMemory Get the usage of swap memory
func GetSwapMemory() (uint64, uint64, uint64) {
	m, _ := mem.SwapMemory()
	return m.Total, m.Used, uint64(m.UsedPercent)
}

// GetHDDSize Get disk capacity information.
func GetHDDSize() (uint64, uint64, uint64) {
	platform := strings.ToLower(runtime.GOOS)
	if strings.Contains(platform, "linux") { // Linux
		stat, _ := disk.Usage("/")
		return stat.Total, stat.Used, uint64(stat.UsedPercent * 100)
	} else if strings.Contains(platform, "win") { // Windows
		_disks, _ := disk.Partitions(false)
		var (
			total       uint64 = 0
			used        uint64 = 0
			usedPercent uint64 = 0
		)
		for _, _disk := range _disks {
			stat, _ := disk.Usage(_disk.Device)
			total += stat.Total
			used += stat.Used
		}
		usedPercent = uint64((float64(used) / float64(total)) * 100)
		return total, used, usedPercent
	}
	return 0, 0, 0
}

// GetCPU Get CPU information.
// Return the number and percentage of CPU usage.
func GetCPU() (int, uint64) {
	counts, _ := cpu.Counts(false)
	percent, _ := cpu.Percent(time.Second, false)
	return counts, uint64(percent[0])
}

// GetLoad Get CPU load information.
func GetLoad() (float64, float64, float64) {
	avg, _ := load.Avg()
	load1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", avg.Load1), 64)
	load5, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", avg.Load1), 64)
	load15, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", avg.Load1), 64)
	return load1, load5, load15
}

// GetHost Get current host information.
// Return the type, number of processes, and running time of the system.
func GetHost() (string, uint64, uint64) {
	_host, _ := host.Info()
	return _host.OS, _host.Procs, _host.Uptime
}

// GetIP Get the external IP address of the local host.
func GetIP() (string, string, string) {
	response, err := http.Get("https://api.vvhan.com/api/getIpInfo")
	if err != nil {
		return "0.0.0.0", "", ""
	}
	bytes, _ := ioutil.ReadAll(response.Body)
	var ipModel model.IPModel
	json.Unmarshal(bytes, &ipModel)
	var ipVersion = "IPv4"
	if strings.Count(ipModel.IP, ":") >= 1 {
		ipVersion = "IPv6"
	}
	return ipModel.IP, ipVersion, ipModel.Info.Country
}

// NetSpeed Monitor network speed.
// TODO: There is an exception, the upload speed is not accurate.
func NetSpeed() {
	for {
		counters, _ := psnet.IOCounters(true)
		var totalRecv uint64 = 0
		var totalSent uint64 = 0
		for _, counter := range counters {
			totalRecv += counter.BytesRecv
			totalSent += counter.BytesSent
		}
		nowClock := uint64(time.Now().Unix())
		netInfo["diff"] = nowClock - netInfo["clock"]
		netInfo["clock"] = nowClock
		netInfo["byteRecv"] = (totalRecv - netInfo["byteTotalRecv"]) / netInfo["diff"]
		netInfo["byteSent"] = (totalSent - netInfo["byteTotalSent"]) / netInfo["diff"]
		netInfo["byteTotalRecv"] = totalRecv
		netInfo["byteTotalSent"] = totalSent
		time.Sleep(time.Second)
	}
}

// PingThread Start a Ping thread to monitor the response time and packet loss rate of the destination host.
func PingThread(host string, port int, mark string) {
	queue := list.New()
	lostPacket := 0
	for {
		// Send the request and start the timer.
		start := time.Now().UnixMilli()
		conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
		// Dial-up request failed, packet lost!
		if err != nil {
			lostPacket += 1
		}
		// Request completed, stop the timer.
		end := time.Now().UnixMilli()

		if queue.Len() > 100 {
			backElement := queue.Back()
			if backElement.Value.(int) <= 10 { // An empty packet.
				lostPacket -= 1
			}
			queue.Remove(backElement)
		}
		queue.PushFront(int(end - start))

		pingTime.Store(mark, queue.Front().Value)
		if queue.Len() > 10 {
			lostRate.Store(mark, int((float64(lostPacket)/float64(queue.Len()))*100))
		}
		if err == nil {
			conn.Close()
		}
		time.Sleep(time.Second)
	}
}

func GetRealtimeData() {
	const (
		CT = "www.189"
		CU = "www.10010.com"
		CM = "www.10086.cn"
	)
	go PingThread(CT, 80, "10000")
	go PingThread(CU, 80, "10010")
	go PingThread(CM, 80, "10086")
	go NetSpeed()
}
