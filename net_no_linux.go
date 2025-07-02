//go:build !linux
// +build !linux

package gonetutil

import (
	"errors"
	"net"
)

func CheckTCPNetPref(conn *net.TCPConn) (TCPInfo, error) {
	return TCPInfo{}, errors.New("TCP info not supported on non-Linux platforms")
}
