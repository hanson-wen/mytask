package util

import (
	"errors"
	"github.com/shirou/gopsutil/cpu"
	"net"
)

// GetIPV4 get ipv4
func GetIPV4() (ip string, err error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addresses {
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		}
	}
	return "", errors.New("ip not found")
}

// GetCpuUsageRate cpu used percent
func GetCpuUsageRate() int8 {
	percent, _ := cpu.Percent(0, false)
	return int8(int(percent[0]))
}
