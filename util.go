package gonetutil

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
)

func IsPortAvailable(port int) (bool, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false, err
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		// 只检测 IPv4 和 IPv6 的可用地址
		if ip == nil {
			continue
		}

		l, err := net.Listen("tcp", net.JoinHostPort(ip.String(), fmt.Sprintf("%d", port)))
		if err != nil {
			if isAddrInUseErr(err) {
				return false, nil // 已经被占用
			}
			if !isAddrNotAvailableErr(err) {
				return false, err // 其他错误
			}
		} else {
			err = l.Close()
			if err != nil {
				return false, err
			}
		}
	}

	// 还要检测 0.0.0.0 绑定情况
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		if isAddrInUseErr(err) {
			return false, nil
		}
		if !isAddrNotAvailableErr(err) {
			return false, err
		}
	} else {
		err = l.Close()
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func isAddrInUseErr(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var sysErr *os.SyscallError
		if errors.As(opErr.Err, &sysErr) {
			return errors.Is(sysErr.Err, syscall.EADDRINUSE)
		}
	}
	return false
}

func isAddrNotAvailableErr(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var sysErr *os.SyscallError
		if errors.As(opErr.Err, &sysErr) {
			return errors.Is(sysErr.Err, syscall.EADDRNOTAVAIL)
		}
	}
	return false
}
func GetMacAddrs() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var macAddrs []string
	for _, netInterface := range interfaces {
		if macAddr := netInterface.HardwareAddr.String(); macAddr != "" {
			macAddrs = append(macAddrs, macAddr)
		}
	}
	return macAddrs, nil
}

func GetLocalIpsWithoutLoopback() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	var ips []string
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && IsLocalIpV4(ipNet.IP) {
			ips = append(ips, ipNet.IP.String())
		}
	}
	return ips, nil
}

const (
	InnerIp            = "$inner_ip"
	NetInterfacePrefix = "$iface"
)

func EvalVarToParseIp(ipVar string) (string, error) {
	if ipVar == InnerIp {
		innerIps, err := GetLocalIpsWithoutLoopback()
		if err != nil {
			return "", err
		}
		if len(innerIps) == 0 {
			return "", errors.New("not found inner ip")
		}
		return innerIps[0], nil
	}

	if strings.HasPrefix(ipVar, NetInterfacePrefix) {
		iface := ipVar[len(NetInterfacePrefix):]
		ip, err := GetIpV4ByIFace(iface)
		if err != nil {
			return "", err
		}
		return ip, nil
	}

	return ipVar, nil
}

func GetIpV4ByIFace(name string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		if iface.Name != name {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			if ip4 := ipNet.IP.To4(); ip4 != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", errors.New(fmt.Sprintf("not found face %s", name))
}

func IsLocalIpV4(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		}
	}
	return false
}
