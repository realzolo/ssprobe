package monitor

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	psnet "github.com/shirou/gopsutil/v3/net"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	CT = "www.189.cn"
	CU = "www.10010.com"
	CM = "221.130.33.52" // www.10086.cn "Ping" is prohibited.
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

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func GetRealtimeData() {
	go PingThread(CT, "10000")
	go PingThread(CU, "10010")
	go PingThread(CM, "10086")
	go NetSpeed()
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
		if NetInfo["diff"] == 0 {
			NetInfo["diff"] = 1
		}
		NetInfo["clock"] = now
		NetInfo["byteRecv"] = (totalRecv - NetInfo["byteTotalRecv"]) / NetInfo["diff"]
		NetInfo["byteSent"] = (totalSent - NetInfo["byteTotalSent"]) / NetInfo["diff"]
		NetInfo["byteTotalRecv"] = totalRecv
		NetInfo["byteTotalSent"] = totalSent
		time.Sleep(time.Second)
	}
}

// PingThread Start a Ping thread to monitor the response time and packet loss rate of the destination host.
func PingThread(host string, mark string) {
	var (
		icmp          ICMP
		remoteAddr, _ = net.ResolveIPAddr("ip", host)
	)
	conn, err := net.DialIP("ip4:icmp", nil, remoteAddr)
	if err != nil {
		log.Panicln(err)
		return
	}
	defer conn.Close()

	icmp = ICMP{8, 0, 0, 0, 0}
	var (
		buffer      bytes.Buffer
		originBytes = make([]byte, 1024)
	)
	binary.Write(&buffer, binary.BigEndian, icmp)
	binary.Write(&buffer, binary.BigEndian, originBytes[0:4])
	b := buffer.Bytes()
	binary.BigEndian.PutUint16(b[2:], checkSum(b))

	var (
		queue      = list.New()
		lostPacket = 0
		recv       = make([]byte, 1024)
	)
	for {
		// The packet is sent to the destination address. If the packet fails to be sent, packet loss occurs.
		if _, err = conn.Write(buffer.Bytes()); err != nil {
			log.Printf("%s: %v\n", mark, err)
			lostPacket++
			enqueue(queue, &lostPacket, 0, mark)
			time.Sleep(time.Second)
			continue
		}

		timeStart := time.Now().UnixMilli()
		conn.SetReadDeadline(time.Now().Add(time.Second * 3))
		_, err = conn.Read(recv)
		// If no response is received, the Ping fails.
		if err != nil {
			log.Printf("%s: %v\n", mark, err)
			lostPacket++
			enqueue(queue, &lostPacket, 0, mark)
			time.Sleep(time.Second)
			continue
		}
		timeEnd := time.Now().UnixMilli()
		timeCost := int(timeEnd - timeStart)
		enqueue(queue, &lostPacket, timeCost, mark)
		time.Sleep(time.Second)
	}
}

func enqueue(queue *list.List, lostPacket *int, value int, mark string) {
	if queue.Len() > 100 { // The maximum length of the queue is 100.
		backElement := queue.Back()
		if backElement.Value.(int) == 0 { // An empty packet.
			*lostPacket--
		}
		queue.Remove(backElement)
	}
	queue.PushFront(value)
	PingTime.Store(mark, queue.Front().Value)
	lostPacketRateStr := fmt.Sprintf("%.f", (float64(*lostPacket)/float64(queue.Len()))*100)
	lostPacketRate, _ := strconv.Atoi(lostPacketRateStr)
	LostRate.Store(mark, lostPacketRate)
}

func checkSum(data []byte) (rt uint16) {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index]) << 8
	}
	rt = uint16(sum) + uint16(sum>>16)

	return ^rt
}
