package monitor

import (
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"runtime"
	"strings"
)

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
		return stat.Total, stat.Used, uint64(stat.UsedPercent)
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
