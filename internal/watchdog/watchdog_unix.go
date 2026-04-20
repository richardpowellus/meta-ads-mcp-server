//go:build !windows

package watchdog

import (
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

func getStartTime(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return ""
	}
	// Start time is field 22 (0-indexed), after the comm field in parens
	s := string(data)
	idx := strings.LastIndex(s, ")")
	if idx < 0 || idx+2 >= len(s) {
		return ""
	}
	fields := strings.Fields(s[idx+2:])
	if len(fields) < 20 {
		return ""
	}
	return fields[19]
}

func monitorParent(ctx context.Context, cancel context.CancelFunc) {
	ppid := os.Getppid()
	if ppid <= 1 {
		return
	}
	origTime := getStartTime(ppid)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Check if parent still exists
			if err := syscall.Kill(ppid, 0); err != nil {
				cancel()
				return
			}
			// Check for PID recycling
			if origTime != "" && getStartTime(ppid) != origTime {
				cancel()
				return
			}
		}
	}
}
