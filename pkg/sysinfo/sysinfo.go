// Package sysinfo reports host-level CPU, memory, and disk utilization for
// the System Status admin page.
package sysinfo

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// Snapshot is a point-in-time read of host resource usage.
type Snapshot struct {
	CPUPercent  float64 `json:"cpu_percent"`
	MemTotal    uint64  `json:"mem_total"`
	MemUsed     uint64  `json:"mem_used"`
	MemPercent  float64 `json:"mem_percent"`
	DiskTotal   uint64  `json:"disk_total"`
	DiskUsed    uint64  `json:"disk_used"`
	DiskPercent float64 `json:"disk_percent"`
	UptimeSecs  uint64  `json:"uptime_seconds"`
}

// Get samples the host's current CPU/RAM/disk usage. diskPath is the
// filesystem mount to report disk usage for (defaults to "/").
//
// The CPU sample blocks for ~300ms to average usage over that window —
// gopsutil's instantaneous (zero-interval) sample is unreliable on the
// first call. That's an acceptable cost for an admin status page polled
// every few seconds, not a hot path.
func Get(diskPath string) Snapshot {
	var snap Snapshot

	if pct, err := cpu.Percent(300*time.Millisecond, false); err == nil && len(pct) > 0 {
		snap.CPUPercent = pct[0]
	}

	if vm, err := mem.VirtualMemory(); err == nil {
		snap.MemTotal = vm.Total
		snap.MemUsed = vm.Used
		snap.MemPercent = vm.UsedPercent
	}

	if diskPath == "" {
		diskPath = "/"
	}
	if du, err := disk.Usage(diskPath); err == nil {
		snap.DiskTotal = du.Total
		snap.DiskUsed = du.Used
		snap.DiskPercent = du.UsedPercent
	}

	if hi, err := host.Info(); err == nil {
		snap.UptimeSecs = hi.Uptime
	}

	return snap
}
