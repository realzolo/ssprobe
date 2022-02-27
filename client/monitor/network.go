package monitor

import (
	"container/list"
	"fmt"
	psnet "github.com/shirou/gopsutil/v3/net"
	"net"
	"strconv"
	"sync"
	"time"
)

var NetInfo = map[string]uint64{
	"byteSent":      0,
	"byteRecv":      0,
	"byteTotalRecv": 0,
	"byteTotalSent": 0,
	"clock":         0,
	"diff":          0,
}
var (
	PingTime sync.Map
	LostRate sync.Map
)

func init() {
	PingTime.Store("10000", 0)
	PingTime.Store("10010", 0)
	PingTime.Store("10086", 0)
	LostRate.Store("10000", 0)
	LostRate.Store("10010", 0)
	LostRate.Store("10086", 0)
}

// NetSpeed Monitor network speed.
// TODO: There is a bug, the upload speed is not accurate.
func NetSpeed() {
	for {
		counters, _ := psnet.IOCounters(true)
		var totalRecv uint64 = 0
		var totalSent uint64 = 0
		for _, counter := range counters {
			totalRecv += counter.BytesRecv
			totalSent += counter.BytesSent
		}
		now := uint64(time.Now().Unix())
		NetInfo["diff"] = now - NetInfo["clock"]
		NetInfo["clock"] = now
		NetInfo["byteRecv"] = (totalRecv - NetInfo["byteTotalRecv"]) / NetInfo["diff"]
		NetInfo["byteSent"] = (totalSent - NetInfo["byteTotalSent"]) / NetInfo["diff"]
		NetInfo["byteTotalRecv"] = totalRecv
		NetInfo["byteTotalSent"] = totalSent
		time.Sleep(time.Second)
	}
}

// PingThread Start a Ping thread to monitor the response time and packet loss rate of the destination host.
// TODO BUG: Packet loss rate calculated incorrectly!
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

		PingTime.Store(mark, queue.Front().Value)
		if queue.Len() > 10 {
			numStr := fmt.Sprintf("%.f", (float64(lostPacket)/float64(queue.Len()))*100)
			num, err := strconv.Atoi(numStr)
			if err != nil {
				LostRate.Store(mark, 0)
			}
			LostRate.Store(mark, num)
		}
		if err == nil {
			conn.Close()
		}
		time.Sleep(time.Second)
	}
}

func GetRealtimeData() {
	const (
		CT = "www.189.cn"
		CU = "www.10010.com"
		CM = "www.10086.cn"
	)
	go PingThread(CT, 80, "10000")
	go PingThread(CU, 80, "10010")
	go PingThread(CM, 80, "10086")
	go NetSpeed()
}
