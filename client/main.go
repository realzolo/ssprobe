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
	"net"
	"net/http"
	"runtime"
	"server-monitor-client/model"
	"strconv"
	"strings"
	"sync"
	"time"
)

var netSpeed = map[string]uint64{
	"byteSent":      0,
	"byteRecv":      0,
	"byteTotalRecv": 0,
	"byteTotalSent": 0,
	"clock":         0,
	"diff":          0,
}

var pingTime sync.Map
var lostRate sync.Map

func init() {
	pingTime.Store("10000", 0)
	pingTime.Store("10010", 0)
	pingTime.Store("10086", 0)
	lostRate.Store("10000", 0.0)
	lostRate.Store("10010", 0.0)
	lostRate.Store("10086", 0.0)
}

func main() {
	go GetRealtimeData()
	for {
		if value, ok := pingTime.Load("10000"); ok {
			fmt.Printf("ping time: %v\t", value)
		}
		if value, ok := lostRate.Load("10000"); ok {
			fmt.Printf("lost rate: %v\n", value)
		}
		if value, ok := pingTime.Load("10010"); ok {
			fmt.Printf("ping time: %v\t", value)
		}
		if value, ok := lostRate.Load("10010"); ok {
			fmt.Printf("lost rate: %v\n", value)
		}
		time.Sleep(time.Second)
	}
}

// GetMemory Get the usage of memory
func GetMemory() (uint64, uint64, float64) {
	m, _ := mem.VirtualMemory()
	return m.Total, m.Used, m.UsedPercent
}

// GetSwapMemory Get the usage of swap memory
func GetSwapMemory() (uint64, uint64, float64) {
	m, _ := mem.SwapMemory()
	return m.Total, m.Used, m.UsedPercent
}

// GetHDDSize Get disk capacity information.
func GetHDDSize() (uint64, uint64, float64) {
	platform := strings.ToLower(runtime.GOOS)
	if strings.Contains(platform, "linux") { // Linux
		stat, _ := disk.Usage("/")
		return stat.Total, stat.Used, stat.UsedPercent
	} else if strings.Contains(platform, "win") { // Windows
		_disks, _ := disk.Partitions(false)
		var (
			total       uint64  = 0
			used        uint64  = 0
			usedPercent float64 = 0.0
		)
		for _, _disk := range _disks {
			stat, _ := disk.Usage(_disk.Device)
			total += stat.Total
			used += stat.Used
		}
		usedPercent = float64(used) / float64(total)
		return total, used, usedPercent
	}
	return 0, 0, 0
}

// GetCPU Get CPU information.
// Return the number and percentage of CPU usage.
func GetCPU() (int, float64) {
	counts, _ := cpu.Counts(false)
	percent, _ := cpu.Percent(time.Second, false)
	return counts, percent[0]
}

// GetLoad Get CPU load information.
func GetLoad() (float64, float64, float64) {
	avg, _ := load.Avg()
	return avg.Load1, avg.Load5, avg.Load15
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
		netSpeed["diff"] = nowClock - netSpeed["clock"]
		netSpeed["clock"] = nowClock
		netSpeed["byteRecv"] = (totalRecv - netSpeed["byteTotalRecv"]) / netSpeed["diff"]
		netSpeed["byteSent"] = (totalSent - netSpeed["byteTotalSent"]) / netSpeed["diff"]
		netSpeed["byteTotalRecv"] = totalRecv
		netSpeed["byteTotalSent"] = totalSent
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
			if backElement.Value.(uint) <= 10 {
				lostPacket -= 1
			}
			queue.Remove(backElement)
		}
		queue.PushFront(uint(end - start))

		pingTime.Store(mark, (queue.Front().Value).(uint))
		if queue.Len() > 10 {
			lostRate.Store(mark, (float32(lostPacket)/float32(queue.Len()))*100)
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
}
