package test

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func TestTCPClient() {
	// 连接到TCP服务器
	conn, err := net.Dial("tcp", "localhost:8082")
	if err != nil {
		fmt.Printf("连接TCP服务器失败: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("已连接到TCP服务器")

	// 读取欢迎消息
	welcomeMsg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("读取欢迎消息失败: %v\n", err)
	} else {
		// 确保UTF-8编码处理
		welcomeMsg = strings.TrimSpace(welcomeMsg)
		fmt.Printf("服务器欢迎消息: %s\n", welcomeMsg)
	}

	// 发送测试消息 - 确保UTF-8编码
	testMsg := "Hello TCP Server!"
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
		// 确保UTF-8编码处理
		response = strings.TrimSpace(response)
		fmt.Printf("服务器响应: %s\n", response)
	}

	// 等待一段时间，然后断开连接
	time.Sleep(2 * time.Second)
	fmt.Println("客户端即将断开连接")
}