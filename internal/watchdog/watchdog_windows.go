//go:build windows

package watchdog

import (
	"context"
	"os"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func monitorParent(ctx context.Context, cancel context.CancelFunc) {
	ppid := os.Getppid()
	if ppid <= 0 {
		return
	}

	origTime := getCreationTime(ppid)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(ppid))
			if err != nil {
				cancel()
				return
			}
			if origTime != 0 {
				ct := getCreationTimeFromHandle(handle)
				windows.CloseHandle(handle)
				if ct != origTime {
					cancel()
					return
				}
			} else {
				windows.CloseHandle(handle)
			}
		}
	}
}

func getCreationTime(pid int) uint64 {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return 0
	}
	defer windows.CloseHandle(handle)
	return getCreationTimeFromHandle(handle)
}

func getCreationTimeFromHandle(handle windows.Handle) uint64 {
	var creation, exit, kernel, user windows.Filetime
	err := windows.GetProcessTimes(handle, &creation, &exit, &kernel, &user)
	if err != nil {
		return 0
	}
	return uint64(creation.HighDateTime)<<32 | uint64(creation.LowDateTime)
}

// getParentPID gets the parent PID using NtQueryInformationProcess.
func getParentPID(pid int) int {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return 0
	}
	defer windows.CloseHandle(handle)

	ntdll := windows.NewLazySystemDLL("ntdll.dll")
	ntQuery := ntdll.NewProc("NtQueryInformationProcess")

	type processBasicInfo struct {
		ExitStatus                   uintptr
		PebBaseAddress               uintptr
		AffinityMask                 uintptr
		BasePriority                 int32
		UniqueProcessId              uintptr
		InheritedFromUniqueProcessId uintptr
	}

	var pbi processBasicInfo
	var retLen uint32
	r, _, _ := ntQuery.Call(
		uintptr(handle),
		0, // ProcessBasicInformation
		uintptr(unsafe.Pointer(&pbi)),
		unsafe.Sizeof(pbi),
		uintptr(unsafe.Pointer(&retLen)),
	)
	if r != 0 {
		return 0
	}
	return int(pbi.InheritedFromUniqueProcessId)
}
