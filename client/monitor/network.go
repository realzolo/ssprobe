package monitor

import (
	"bytes"
	"container/list"
	"context"
	"encoding/binary"
	"fmt"
	psnet "github.com/shirou/gopsutil/v3/net"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

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
	lock     sync.Mutex
)

func init() {
	PingTime.Store("10000", 0)
	PingTime.Store("10010", 0)
	PingTime.Store("10086", 0)
	LostRate.Store("10000", 0)
	LostRate.Store("10010", 0)
	LostRate.Store("10086", 0)
}

func GetRealtimeData(ctx context.Context) {
	go PingThread(ctx, CT, "10000")
	go PingThread(ctx, CU, "10010")
	go PingThread(ctx, CM, "10086")
	go NetSpeed(ctx)
}

// NetSpeed Monitor network speed.
// TODO: The upload speed is not accurate, maybe this is a BUG.
func NetSpeed(ctx context.Context) {
	var totalRecv uint64 = 0
	var totalSent uint64 = 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			counters, _ := psnet.IOCounters(true)
			for _, counter := range counters {
				totalRecv += counter.BytesRecv
				totalSent += counter.BytesSent
			}
			now := uint64(time.Now().Unix())

			lock.Lock()
			NetInfo["diff"] = now - NetInfo["clock"]
			NetInfo["clock"] = now
			NetInfo["byteRecv"] = (totalRecv - NetInfo["byteTotalRecv"]) / NetInfo["diff"]
			NetInfo["byteSent"] = (totalSent - NetInfo["byteTotalSent"]) / NetInfo["diff"]
			NetInfo["byteTotalRecv"] = totalRecv
			NetInfo["byteTotalSent"] = totalSent
			totalRecv = 0
			totalSent = 0
			lock.Unlock()

			time.Sleep(time.Second * 2)
		}
	}
}

// PingThread Start a Ping thread to monitor the response time and packet loss rate of the destination host.
func PingThread(ctx context.Context, host string, mark string) {
	var (
		icmp          = ICMP{8, 0, 0, 0, 0}
		remoteAddr, _ = net.ResolveIPAddr("ip", host)
	)
	conn, err := net.DialIP("ip4:icmp", nil, remoteAddr)
	if err != nil {
		log.Panicln(err)
		return
	}
	defer conn.Close()

	var (
		buffer      bytes.Buffer
		originBytes = make([]byte, 1024)
	)
	_ = binary.Write(&buffer, binary.BigEndian, icmp)
	_ = binary.Write(&buffer, binary.BigEndian, originBytes[0:4])
	b := buffer.Bytes()
	binary.BigEndian.PutUint16(b[2:], checkSum(b))

	var (
		queue      = list.New()
		lostPacket = 0
		recv       = make([]byte, 1024)
	)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			lock.Lock()
			// The packet is sent to the destination address. If the packet fails to be sent, packet loss occurs.
			if _, err = conn.Write(buffer.Bytes()); err != nil {
				logger.LogWithFormat("[Ping][%s]: %s", mark, err.Error())
				lostPacket++
				enqueue(queue, &lostPacket, 0, mark)
				lock.Unlock()
				time.Sleep(time.Second * 2)
				continue
			}

			timeStart := time.Now().UnixMilli()
			_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
			// If no response is received, the Ping fails.
			if _, err = conn.Read(recv); err != nil {
				logger.LogWithFormat("[Ping][%s]: %s", mark, err.Error())
				lostPacket++
				enqueue(queue, &lostPacket, 0, mark)
				lock.Unlock()
				time.Sleep(time.Second * 2)
				continue
			}
			timeEnd := time.Now().UnixMilli()
			timeCost := int(timeEnd - timeStart)
			enqueue(queue, &lostPacket, timeCost, mark)
			lock.Unlock()
			time.Sleep(time.Second * 2)
		}
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

	lostPacketRateStr := fmt.Sprintf("%.f", (float64(*lostPacket)/float64(queue.Len()))*100)
	lostPacketRate, _ := strconv.Atoi(lostPacketRateStr)
	// TODO: sometimes lostPacketRate is inaccurate.
	if lostPacketRate > 100 {
		logger.OnlyLog("--------Exception--------")
		logger.OnlyLog("lostPacket: " + (string(rune(*lostPacket))))
		logger.OnlyLog("queue.Len: " + string(rune(queue.Len())))
		logger.OnlyLog("----------------")
	}
	PingTime.Store(mark, value)
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
