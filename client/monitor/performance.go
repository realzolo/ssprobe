package monitor

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"strconv"
	"strings"
)

// GetCPU get CPU information.
// Return the number and percentage of CPU usage.
func GetCPU() (int, uint64) {
	counts, _ := cpu.Counts(false)
	percents, _ := cpu.Percent(0, true)
	var pct = 0.0
	for _, p := range percents {
		pct += p
	}
	var percent = uint64(pct)
	if percent > 100 {
		percent = 100
	}
	return counts, percent
}

// GetLoad get CPU load information.
func GetLoad() (float64, float64, float64) {
	avg, _ := load.Avg()
	load1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", avg.Load1), 64)
	load5, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", avg.Load1), 64)
	load15, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", avg.Load1), 64)
	return load1, load5, load15
}

// GetHost get current host information.
// Return the type, number of processes, and running time of the system.
func GetHost() (string, uint64, uint64) {
	_host, _ := host.Info()
	platform := strings.ToLower(_host.Platform)
	if strings.Contains(platform, "windows") {
		platform = "windows"
	}
	return platform, _host.Procs, _host.Uptime
}
