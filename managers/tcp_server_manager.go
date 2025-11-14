package managers

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	
	"windows-gui-app/interfaces"
)

// TCPClient 表示TCP客户端连接
type TCPClient struct {
	ID     string
	Conn   net.Conn
	Addr   string
	Active bool
}

// TCPServerManager TCP服务器管理器
type TCPServerManager struct {
	uiInstance      interfaces.UIInterface
	listener        net.Listener
	clients         map[string]*TCPClient
	clientsMutex    sync.RWMutex
	isRunning       bool
	serverMutex     sync.RWMutex
	port            int
	host            string
	connectionsMade int
}

// NewTCPServerManager 创建新的TCP服务器管理器
func NewTCPServerManager(ui interfaces.UIInterface) *TCPServerManager {
	return &TCPServerManager{
		uiInstance:   ui,
		clients:      make(map[string]*TCPClient),
		isRunning:    false,
		port:         8082, // 默认端口
		host:         "localhost",
	}
}

// SetUIInstance 设置UI实例
func (m *TCPServerManager) SetUIInstance(ui interfaces.UIInterface) {
	m.uiInstance = ui
}

// StartServer 启动TCP服务器
func (m *TCPServerManager) StartServer() error {
	m.serverMutex.Lock()
	defer m.serverMutex.Unlock()

	if m.isRunning {
		return fmt.Errorf("TCP服务器已在运行")
	}

	// 创建监听器
	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("启动TCP服务器失败: %v", err)
	}

	m.listener = listener
	m.isRunning = true

	// 启动接受连接的goroutine
	go m.acceptConnections()

	if m.uiInstance != nil {
		m.uiInstance.AddLog(fmt.Sprintf("TCP服务器已启动，监听地址: %s", addr))
	}

	return nil
}

// StopServer 停止TCP服务器
func (m *TCPServerManager) StopServer() error {
	m.serverMutex.Lock()
	defer m.serverMutex.Unlock()

	if !m.isRunning {
		return fmt.Errorf("TCP服务器未运行")
	}

	// 关闭监听器
	if m.listener != nil {
		m.listener.Close()
	}

	// 关闭所有客户端连接
	m.clientsMutex.Lock()
	for _, client := range m.clients {
		if client.Conn != nil {
			client.Conn.Close()
		}
	}
	m.clients = make(map[string]*TCPClient)
	m.clientsMutex.Unlock()

	m.isRunning = false

	if m.uiInstance != nil {
		m.uiInstance.AddLog("TCP服务器已停止")
	}

	return nil
}

// IsRunning 检查服务器是否运行中
func (m *TCPServerManager) IsRunning() bool {
	m.serverMutex.RLock()
	defer m.serverMutex.RUnlock()
	return m.isRunning
}

// acceptConnections 接受客户端连接
func (m *TCPServerManager) acceptConnections() {
	defer func() {
		if r := recover(); r != nil {
			if m.uiInstance != nil {
				m.uiInstance.AddLog(fmt.Sprintf("TCP服务器接受连接时发生错误: %v", r))
			}
		}
	}()

	for {
		conn, err := m.listener.Accept()
		if err != nil {
			// 服务器关闭时会产生错误，这是正常的
			if m.isRunning {
				if m.uiInstance != nil {
					m.uiInstance.AddLog(fmt.Sprintf("接受TCP连接失败: %v", err))
				}
			}
			return
		}

		// 创建客户端
		clientID := fmt.Sprintf("tcp_%d", m.connectionsMade+1)
		client := &TCPClient{
			ID:     clientID,
			Conn:   conn,
			Addr:   conn.RemoteAddr().String(),
			Active: true,
		}

		// 添加到客户端列表
		m.clientsMutex.Lock()
		m.clients[clientID] = client
		m.clientsMutex.Unlock()

		m.connectionsMade++

		if m.uiInstance != nil {
			m.uiInstance.AddLog(fmt.Sprintf("新的TCP客户端连接: %s (%s)", clientID, client.Addr))
		}

		// 启动处理客户端的goroutine
		go m.handleClient(client)
	}
}

// handleClient 处理客户端连接
func (m *TCPServerManager) handleClient(client *TCPClient) {
	defer func() {
		// 连接关闭时清理客户端
		m.clientsMutex.Lock()
		client.Active = false
		delete(m.clients, client.ID)
		m.clientsMutex.Unlock()

		if client.Conn != nil {
			client.Conn.Close()
		}

		if m.uiInstance != nil {
			m.uiInstance.AddLog(fmt.Sprintf("TCP客户端断开连接: %s", client.ID))
		}
	}()

	// 发送欢迎消息 - 确保使用UTF-8编码
	welcomeMsg := fmt.Sprintf("欢迎连接到TCP服务器! 您的客户端ID: %s\n", client.ID)
	_, err := client.Conn.Write([]byte(welcomeMsg))
	if err != nil && m.uiInstance != nil {
		m.uiInstance.AddLog(fmt.Sprintf("发送欢迎消息失败 (%s): %v", client.ID, err))
	}

	// 使用缓冲读取器提高读取效率
	reader := bufio.NewReader(client.Conn)
	
	for {
		// 读取客户端数据，直到遇到换行符
		data, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF && m.uiInstance != nil {
				m.uiInstance.AddLog(fmt.Sprintf("读取TCP客户端数据失败 (%s): %v", client.ID, err))
			}
			break
		}

		// 处理接收到的数据 - 确保UTF-8编码处理
		// 去除可能的换行符和回车符
		data = strings.TrimSpace(data)
		
		if m.uiInstance != nil {
			m.uiInstance.AddLog(fmt.Sprintf("收到TCP客户端消息 (%s): %s", client.ID, data))
		}

		// 发送响应 - 确保使用UTF-8编码
		response := fmt.Sprintf("服务器收到消息: %s\n", data)
		_, err = client.Conn.Write([]byte(response))
		if err != nil && m.uiInstance != nil {
			m.uiInstance.AddLog(fmt.Sprintf("发送响应失败 (%s): %v", client.ID, err))
			break
		}
	}
}

// GetServerInfo 获取服务器信息
func (m *TCPServerManager) GetServerInfo() map[string]interface{} {
	m.serverMutex.RLock()
	m.clientsMutex.RLock()
	defer m.serverMutex.RUnlock()
	defer m.clientsMutex.RUnlock()

	info := map[string]interface{}{
		"running":          m.isRunning,
		"host":             m.host,
		"port":             m.port,
		"connections_made": m.connectionsMade,
		"active_clients":   len(m.clients),
	}

	if m.isRunning {
		info["listen_address"] = fmt.Sprintf("%s:%d", m.host, m.port)
	}

	return info
}

// GetClients 获取客户端列表
func (m *TCPServerManager) GetClients() map[string]interface{} {
	m.clientsMutex.RLock()
	defer m.clientsMutex.RUnlock()

	clients := make(map[string]interface{})
	for id, client := range m.clients {
		clients[id] = map[string]interface{}{
			"id":     client.ID,
			"addr":   client.Addr,
			"active": client.Active,
		}
	}

	return clients
}

// BroadcastMessage 广播消息给所有客户端
func (m *TCPServerManager) BroadcastMessage(message string) error {
	if !m.IsRunning() {
		return fmt.Errorf("TCP服务器未运行")
	}

	m.clientsMutex.RLock()
	defer m.clientsMutex.RUnlock()

	for _, client := range m.clients {
		if client.Active && client.Conn != nil {
			// 确保消息以换行符结尾，并使用UTF-8编码
			if !strings.HasSuffix(message, "\n") {
				message = message + "\n"
			}
			_, err := client.Conn.Write([]byte(message))
			if err != nil {
				if m.uiInstance != nil {
					m.uiInstance.AddLog(fmt.Sprintf("向TCP客户端发送消息失败 (%s): %v", client.ID, err))
				}
			}
		}
	}

	return nil
}

// SendMessageToClient 向指定客户端发送消息
func (m *TCPServerManager) SendMessageToClient(clientID, message string) error {
	if !m.IsRunning() {
		return fmt.Errorf("TCP服务器未运行")
	}

	m.clientsMutex.RLock()
	client, exists := m.clients[clientID]
	m.clientsMutex.RUnlock()

	if !exists {
		return fmt.Errorf("客户端不存在: %s", clientID)
	}

	if !client.Active || client.Conn == nil {
		return fmt.Errorf("客户端连接已断开: %s", clientID)
	}

	// 确保消息以换行符结尾，并使用UTF-8编码
	if !strings.HasSuffix(message, "\n") {
		message = message + "\n"
	}
	_, err := client.Conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("向客户端发送消息失败 (%s): %v", clientID, err)
	}

	return nil
}

// SetPort 设置服务器端口
func (m *TCPServerManager) SetPort(port int) {
	if !m.IsRunning() {
		m.port = port
	}
}

// SetHost 设置服务器主机
func (m *TCPServerManager) SetHost(host string) {
	if !m.IsRunning() {
		m.host = host
	}
}

// GetPort 获取服务器端口
func (m *TCPServerManager) GetPort() int {
	return m.port
}

// GetHost 获取服务器主机
func (m *TCPServerManager) GetHost() string {
	return m.host
}