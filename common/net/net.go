package net

import (
	"fmt"
	"net"
)

// 获取可用端口
func GetAvailablePort() int {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		fmt.Println(err)
		return 0
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

// 判断端口是否可以（未被占用）
func IsPortAvailable(port int) bool {
	address := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("port %s is taken: %s", address, err)
		return false
	}

	defer listener.Close()
	return true
}
