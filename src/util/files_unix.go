//go:build !windows

package util

import (
	"os"
	"syscall"
)

func GetFileDevice(path string) (uint64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	sys := fileInfo.Sys()
	if sys == nil {
		return 0, nil // Not supported on this OS
	}

	sysStat, ok := sys.(*syscall.Stat_t)
	if !ok {
		return 0, nil // Not a syscall.Stat_t
	}

	return sysStat.Dev, nil
}
