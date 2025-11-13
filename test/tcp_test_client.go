package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// 从命令行参数获取服务器地址
	serverAddr := "localhost:8082"
	if len(os.Args) > 1 {
		serverAddr = os.Args[1]
	}

	// 连接到TCP服务器
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("连接TCP服务器失败: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("已连接到TCP服务器 (%s)\n", serverAddr)

	// 读取欢迎消息
	welcomeMsg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("读取欢迎消息失败: %v\n", err)
	} else {
		fmt.Printf("服务器欢迎消息: %s", welcomeMsg)
	}

	// 发送测试消息
	testMsg := "Hello TCP Server from test client!"
	_, err = conn.Write([]byte(testMsg))
	if err != nil {
		fmt.Printf("发送消息失败: %v\n", err)
		return
	}
	fmt.Printf("已发送消息: %s\n", testMsg)

	// 读取服务器响应
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("读取服务器响应失败: %v\n", err)
	} else {
		fmt.Printf("服务器响应: %s", response)
	}

	// 等待一段时间，然后断开连接
	time.Sleep(2 * time.Second)
	fmt.Println("客户端即将断开连接")
}