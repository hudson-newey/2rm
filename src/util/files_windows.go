//go:build windows

package util

func GetFileDevice(path string) (uint64, error) {
	return 0, nil
}
