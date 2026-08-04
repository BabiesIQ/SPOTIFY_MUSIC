/*
 * SPOTIFY_MUSIC - Telegram Music Bot
 *  Copyright (c) 2025-2026 BabiesIQ
 *
 *  Licensed under GNU GPL v3
 *  See https://github.com/BabiesIQ/SPOTIFY_MUSIC
 */

package handlers

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/BabiesIQ/SPOTIFY_MUSIC/bot/store"

	td "github.com/BabiesIQ/gotdbot"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

type AppStats struct {
	Uptime     string
	Goroutines int
	GoVersion  string

	AppMemUsed string
	AppHeap    string
	GCCount    uint32
	GCPause    string

	MemLimit  string
	DiskUsed  string
	DiskTotal string

	SystemCPU string
	AppCPU    string

	SystemMemUsed  string
	SystemMemTotal string
	CPUCores       int
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func readContainerMemLimit() uint64 {
	if data, err := os.ReadFile("/sys/fs/cgroup/memory.max"); err == nil {
		val := strings.TrimSpace(string(data))
		if val != "max" {
			if v, err := strconv.ParseUint(val, 10, 64); err == nil {
				return v
			}
		}
	}

	if data, err := os.ReadFile("/sys/fs/cgroup/memory/memory.limit_in_bytes"); err == nil {
		if v, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64); err == nil && v < (1<<60) {
			return v
		}
	}
	return 0
}

func storageDiskUsage(path string) (used, total string) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return "N/A", "N/A"
	}

	totalBytes := stat.Blocks * uint64(stat.Bsize)
	freeBytes := stat.Bfree * uint64(stat.Bsize)
	usedBytes := totalBytes - freeBytes

	return formatBytes(usedBytes), formatBytes(totalBytes)
}

func systemMemoryStats() (used, total string) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return "N/A", "N/A"
	}
	return formatBytes(v.Used), formatBytes(v.Total)
}

func measureMemoryStats() (used, heap string, gc uint32, pause string) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	return formatBytes(ms.Alloc),
		formatBytes(ms.HeapAlloc),
		ms.NumGC,
		(time.Duration(ms.PauseTotalNs) * time.Nanosecond).String()
}

func systemCPUPercent() string {
	p, err := cpu.Percent(500*time.Millisecond, false)
	if err != nil || len(p) == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%.2f%%", p[0])
}

func measureCPUPercent() string {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return "N/A"
	}

	v, err := p.CPUPercent()
	if err != nil {
		return "N/A"
	}

	return fmt.Sprintf("%.2f%%", v)
}

func collectAppStats() *AppStats {
	memUsed, heap, gcCount, gcPause := measureMemoryStats()
	sysMemUsed, sysMemTotal := systemMemoryStats()

	root := "/"
	if runtime.GOOS == "windows" {
		root = "C:\\"
	}

	dUsed, dTotal := storageDiskUsage(root)

	stats := &AppStats{
		Uptime:     time.Since(startTime).Round(time.Second).String(),
		Goroutines: runtime.NumGoroutine(),
		GoVersion:  runtime.Version(),

		AppMemUsed: memUsed,
		AppHeap:    heap,
		GCCount:    gcCount,
		GCPause:    gcPause,

		DiskUsed:  dUsed,
		DiskTotal: dTotal,

		SystemCPU:      systemCPUPercent(),
		AppCPU:         measureCPUPercent(),
		SystemMemUsed:  sysMemUsed,
		SystemMemTotal: sysMemTotal,
		CPUCores:       runtime.NumCPU(),
	}

	if limit := readContainerMemLimit(); limit > 0 {
		stats.MemLimit = formatBytes(limit)
	}

	return stats
}
func statsHandler(c *td.Client, m *td.Message) error {
	if !isDeveloper(c, m) {
		return td.EndGroups
	}

	msg := m
	sysMsg, err := msg.ReplyText(c, "Collecting system statistics...", nil)
	if err != nil {
		return err
	}

	stats := collectAppStats()

	chats, _ := db.Instance.GetAllChats()
	users, _ := db.Instance.GetAllUsers()

	memLine := fmt.Sprintf("• Ram usage: %s\n", stats.AppMemUsed)
	if stats.MemLimit != "" {
		memLine = fmt.Sprintf("• Ram usage: %s | %s\n", stats.AppMemUsed, stats.MemLimit)
	}

	text := fmt.Sprintf(
		"<b>%s — Runtime Status</b>\n"+
			"────────────────────────────────────\n\n"+
			"<b>System</b>\n"+
			"• CPU usage: %s (%d cores)\n"+
			"• Ram usage: %s | %s\n"+
			"• Storage: %s | %s\n\n"+
			"<b>Application</b>\n"+
			"• Uptime: %s\n"+
			"• Goroutines: %d\n"+
			"• Go AppVersion: %s\n"+
			"• CPU usage: %s\n"+
			"%s"+
			"• Heap: %s\n"+
			"• GC Runs: %d (pause %s)\n\n"+
			"<b>Database</b>\n"+
			"• Chats: %d\n"+
			"• Users: %d\n\n"+
			"────────────────────────────────────",

		c.Me.FirstName,

		stats.SystemCPU,
		stats.CPUCores,
		stats.SystemMemUsed,
		stats.SystemMemTotal,
		stats.DiskUsed,
		stats.DiskTotal,

		stats.Uptime,
		stats.Goroutines,
		stats.GoVersion,
		stats.AppCPU,

		memLine,

		stats.AppHeap,
		stats.GCCount,
		stats.GCPause,

		len(chats),
		len(users),
	)

	_, _ = sysMsg.EditText(c, text, &td.EditTextMessageOpts{ParseMode: "HTML"})
	return nil
}
