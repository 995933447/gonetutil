package gonetutil

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestCheckTcpNetPref(testing *testing.T) {
	conn, err := net.DialTimeout("tcp", "google.com:80", time.Second*5)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	info, err := CheckTCPNetPref(conn.(*net.TCPConn))
	if err != nil {
		fmt.Println("GetTCPInfo error:", err)
		return
	}

	fmt.Println("==== TCP 状态诊断 ====")
	fmt.Printf("RTT: %.3f ms\n", float64(info.Rtt)/1000)
	fmt.Printf("RTT波动: %.3f ms\n", float64(info.Rttvar)/1000)
	fmt.Printf("重传次数: %d\n", info.Total_retrans)
	fmt.Printf("未确认段数（Unacked）: %d\n", info.Unacked)
	fmt.Printf("剩余: %d\n", info.Rcv_space)
	fmt.Printf("拥塞窗口(cwnd): %d\n", info.Snd_cwnd)
	fmt.Printf("慢启动阈值(ssthresh): %d\n", info.Snd_ssthresh)
	fmt.Printf("最大窗口(MaxSeg): %d\n", info.Snd_mss)
	fmt.Printf("发送缓冲区大小(SO_SNDBUF): %d\n", info.SndBufSize)
	fmt.Println("======================")
}
