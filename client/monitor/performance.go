package monitor

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"strconv"
	"time"
)

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
