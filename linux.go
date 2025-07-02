//go:build linux
// +build linux

package gonetutil

import (
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"syscall"
)

func CheckTCPNetPref(conn *net.TCPConn) (TCPInfo, error) {
	rawConn, err := conn.SyscallConn()
	if err != nil {
		return TCPInfo{}, fmt.Errorf("failed to get raw conn: %v", err)
	}

	var (
		info TCPInfo
		e    error
	)
	err = rawConn.Control(func(fd uintptr) {
		// 1. 获取发送缓冲区大小（SO_SNDBUF）
		sndBufSize, _ := unix.GetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUF)

		// 2. 获取 TCP_INFO
		var t *unix.TCPInfo
		t, e = unix.GetsockoptTCPInfo(int(fd), syscall.SOL_TCP, unix.TCP_INFO)
		if e != nil {
			return
		}

		info.State = t.State
		info.Ca_state = t.Ca_state
		info.Retransmits = t.Retransmits
		info.Probes = t.Probes
		info.Backoff = t.Backoff
		info.Options = t.Options
		info.Rto = t.Rto
		info.Ato = t.Ato
		info.Snd_mss = t.Snd_mss
		info.Rcv_mss = t.Rcv_mss
		info.Unacked = t.Unacked
		info.Sacked = t.Sacked
		info.Lost = t.Lost
		info.Retrans = t.Retrans
		info.Fackets = t.Fackets
		info.Last_data_sent = t.Last_data_sent
		info.Last_ack_sent = t.Last_ack_sent
		info.Last_data_recv = t.Last_data_recv
		info.Last_ack_recv = t.Last_ack_recv
		info.Pmtu = t.Pmtu
		info.Rcv_ssthresh = t.Rcv_ssthresh
		info.Rtt = t.Rtt
		info.Rttvar = t.Rttvar
		info.Snd_ssthresh = t.Snd_ssthresh
		info.Snd_cwnd = t.Snd_cwnd
		info.Advmss = t.Advmss
		info.Reordering = t.Reordering
		info.Rcv_rtt = t.Rcv_rtt
		info.Rcv_space = t.Rcv_space
		info.Total_retrans = t.Total_retrans
		info.SndBufSize = sndBufSize
		return
	})
	if err != nil {
		return TCPInfo{}, err
	}
	return info, err
}
