// managers/webserver_manager.go - Web服务器管理器
// 该文件实现了Web服务器管理器，提供HTTP和WebSocket服务器功能
// 支持启动/停止服务器，处理HTTP请求和WebSocket连接
package managers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"windows-gui-app/interfaces"
)

// WebSocketClient WebSocket客户端信息
type WebSocketClient struct {
	ID       string
	Conn     *websocket.Conn
	SendChan chan []byte
}

// WebServerManager Web服务器管理器
// 封装了HTTP和WebSocket服务器功能，提供统一的Web服务接口
// 与UI组件交互，将服务器状态和操作结果展示在用户界面上
type WebServerManager struct {
	uiInstance      interfaces.UIInterface // UI实例接口，用于更新界面状态和数据
	httpServer      *http.Server          // HTTP服务器实例
	wsUpgrader      websocket.Upgrader    // WebSocket升级器
	clients         map[string]*WebSocketClient // WebSocket客户端连接
	clientsMutex    sync.RWMutex          // 客户端映射的读写锁
	isRunning       bool                  // 服务器运行状态
	isRunningMutex  sync.RWMutex          // 运行状态的读写锁
	httpPort        int                   // HTTP服务器端口
	wsPort          int                   // WebSocket服务器端口
	serverStartTime time.Time             // 服务器启动时间
}

// NewWebServerManager 创建Web服务器管理器实例
// 初始化WebSocket升级器和客户端映射，设置默认端口
// 参数:
//   - ui: UI实例，用于显示操作结果（可以为nil，之后通过SetUIInstance设置）
//
// 返回:
//   - *WebServerManager: Web服务器管理器实例
func NewWebServerManager(ui interfaces.UIInterface) *WebServerManager {
	return &WebServerManager{
		uiInstance: ui,
		wsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源的WebSocket连接
			},
		},
		clients:    make(map[string]*WebSocketClient),
		httpPort:   8080, // 默认HTTP端口
		wsPort:     8081, // 默认WebSocket端口
		isRunning:  false,
	}
}

// SetUIInstance 设置UI实例
// 允许在创建管理器后设置UI实例，用于解耦合UI和管理器的创建顺序
// 参数:
//   - ui: UI实例接口
func (wsm *WebServerManager) SetUIInstance(ui interfaces.UIInterface) {
	wsm.uiInstance = ui
}

// SetPorts 设置HTTP和WebSocket端口
// 参数:
//   - httpPort: HTTP服务器端口
//   - wsPort: WebSocket服务器端口
func (wsm *WebServerManager) SetPorts(httpPort, wsPort int) {
	wsm.httpPort = httpPort
	wsm.wsPort = wsPort
}

// GetServerInfo 获取服务器信息
// 返回:
//   - map[string]interface{}: 包含服务器状态、端口、客户端数量等信息
func (wsm *WebServerManager) GetServerInfo() map[string]interface{} {
	wsm.isRunningMutex.RLock()
	isRunning := wsm.isRunning
	wsm.isRunningMutex.RUnlock()

	wsm.clientsMutex.RLock()
	clientCount := len(wsm.clients)
	wsm.clientsMutex.RUnlock()

	info := map[string]interface{}{
		"运行状态":       isRunning,
		"HTTP端口":      wsm.httpPort,
		"WebSocket端口": wsm.wsPort,
		"连接客户端数":    clientCount,
	}

	if isRunning {
		info["启动时间"] = wsm.serverStartTime.Format("2006-01-02 15:04:05")
		info["运行时长"] = time.Since(wsm.serverStartTime).String()
	}

	return info
}

// StartServer 启动Web服务器
// 启动HTTP和WebSocket服务器，设置路由和处理函数
// 返回:
//   - error: 操作错误，nil表示成功
func (wsm *WebServerManager) StartServer() error {
	wsm.isRunningMutex.Lock()
	defer wsm.isRunningMutex.Unlock()

	if wsm.isRunning {
		return fmt.Errorf("服务器已在运行中")
	}

	// 创建HTTP路由器
	mux := http.NewServeMux()

	// 设置HTTP路由
	mux.HandleFunc("/", wsm.handleHome)
	mux.HandleFunc("/api/status", wsm.handleStatusAPI)
	mux.HandleFunc("/api/time", wsm.handleTimeAPI)
	mux.HandleFunc("/ws", wsm.handleWebSocket)

	// 创建HTTP服务器
	wsm.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", wsm.httpPort),
		Handler: mux,
	}

	// 启动HTTP服务器
	go func() {
		if err := wsm.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP服务器错误: %v", err)
			if wsm.uiInstance != nil {
				wsm.uiInstance.AddLog(fmt.Sprintf("HTTP服务器错误: %v", err))
			}
		}
	}()

	// 标记服务器为运行状态
	wsm.isRunning = true
	wsm.serverStartTime = time.Now()

	// 更新UI状态
	if wsm.uiInstance != nil {
		wsm.uiInstance.AddLog(fmt.Sprintf("Web服务器已启动"))
		wsm.uiInstance.AddLog(fmt.Sprintf("HTTP服务: http://localhost:%d", wsm.httpPort))
		wsm.uiInstance.AddLog(fmt.Sprintf("WebSocket服务: ws://localhost:%d/ws", wsm.wsPort))
		wsm.uiInstance.UpdateStatus("Web服务器运行中")
	}

	log.Printf("Web服务器已启动，HTTP端口: %d, WebSocket端口: %d", wsm.httpPort, wsm.wsPort)
	return nil
}

// StopServer 停止Web服务器
// 关闭HTTP服务器和所有WebSocket连接
// 返回:
//   - error: 操作错误，nil表示成功
func (wsm *WebServerManager) StopServer() error {
	wsm.isRunningMutex.Lock()
	defer wsm.isRunningMutex.Unlock()

	if !wsm.isRunning {
		return fmt.Errorf("服务器未运行")
	}

	// 关闭HTTP服务器
	if wsm.httpServer != nil {
		if err := wsm.httpServer.Close(); err != nil {
			log.Printf("关闭HTTP服务器时出错: %v", err)
		}
	}

	// 关闭所有WebSocket连接
	wsm.clientsMutex.Lock()
	for id, client := range wsm.clients {
		close(client.SendChan)
		client.Conn.Close()
		delete(wsm.clients, id)
	}
	wsm.clientsMutex.Unlock()

	// 标记服务器为停止状态
	wsm.isRunning = false

	// 更新UI状态
	if wsm.uiInstance != nil {
		wsm.uiInstance.AddLog("Web服务器已停止")
		wsm.uiInstance.UpdateStatus("Web服务器已停止")
	}

	log.Println("Web服务器已停止")
	return nil
}

// IsRunning 检查服务器是否正在运行
// 返回:
//   - bool: 服务器运行状态
func (wsm *WebServerManager) IsRunning() bool {
	wsm.isRunningMutex.RLock()
	defer wsm.isRunningMutex.RUnlock()
	return wsm.isRunning
}

// handleHome 处理主页请求
func (wsm *WebServerManager) handleHome(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>GO VCL Web服务器</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .container { max-width: 800px; margin: 0 auto; }
        .status { background-color: #f0f0f0; padding: 10px; border-radius: 5px; margin-bottom: 20px; }
        .websocket { border: 1px solid #ccc; padding: 10px; margin-top: 20px; }
        #messages { height: 200px; overflow-y: scroll; border: 1px solid #ccc; padding: 10px; margin-top: 10px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>GO VCL Web服务器</h1>
        <div class="status">
            <h2>服务器状态</h2>
            <p>HTTP端口: ` + fmt.Sprintf("%d", wsm.httpPort) + `</p>
            <p>WebSocket端口: ` + fmt.Sprintf("%d", wsm.wsPort) + `</p>
            <p>启动时间: ` + wsm.serverStartTime.Format("2006-01-02 15:04:05") + `</p>
            <p>运行时长: ` + time.Since(wsm.serverStartTime).String() + `</p>
        </div>
        
        <div class="websocket">
            <h2>WebSocket测试</h2>
            <input type="text" id="messageInput" placeholder="输入消息">
            <button onclick="sendMessage()">发送消息</button>
            <div id="messages"></div>
        </div>
    </div>
    
    <script>
        const ws = new WebSocket('ws://localhost:` + fmt.Sprintf("%d", wsm.wsPort) + `/ws');
        const messages = document.getElementById('messages');
        
        ws.onmessage = function(event) {
            const message = document.createElement('div');
            message.textContent = '收到: ' + event.data;
            messages.appendChild(message);
            messages.scrollTop = messages.scrollHeight;
        };
        
        function sendMessage() {
            const input = document.getElementById('messageInput');
            if (input.value) {
                ws.send(input.value);
                const message = document.createElement('div');
                message.textContent = '发送: ' + input.value;
                messages.appendChild(message);
                messages.scrollTop = messages.scrollHeight;
                input.value = '';
            }
        }
    </script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// handleStatusAPI 处理状态API请求
func (wsm *WebServerManager) handleStatusAPI(w http.ResponseWriter, r *http.Request) {
	wsm.clientsMutex.RLock()
	clientCount := len(wsm.clients)
	wsm.clientsMutex.RUnlock()

	status := map[string]interface{}{
		"status":     "running",
		"http_port":  wsm.httpPort,
		"ws_port":    wsm.wsPort,
		"clients":    clientCount,
		"start_time": wsm.serverStartTime.Format(time.RFC3339),
		"uptime":     time.Since(wsm.serverStartTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// handleTimeAPI 处理时间API请求
func (wsm *WebServerManager) handleTimeAPI(w http.ResponseWriter, r *http.Request) {
	timeData := map[string]interface{}{
		"current_time": time.Now().Format(time.RFC3339),
		"unix_time":    time.Now().Unix(),
		"formatted":    time.Now().Format("2006-01-02 15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(timeData)
}

// handleWebSocket 处理WebSocket连接
func (wsm *WebServerManager) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := wsm.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	// 生成客户端ID
	clientID := fmt.Sprintf("client_%d", time.Now().UnixNano())

	// 创建客户端
	client := &WebSocketClient{
		ID:       clientID,
		Conn:     conn,
		SendChan: make(chan []byte, 256),
	}

	// 添加客户端到映射
	wsm.clientsMutex.Lock()
	wsm.clients[clientID] = client
	wsm.clientsMutex.Unlock()

	// 更新UI状态
	if wsm.uiInstance != nil {
		wsm.uiInstance.AddLog(fmt.Sprintf("新的WebSocket连接: %s", clientID))
	}

	// 启动读取和写入goroutine
	go wsm.readPump(client)
	go wsm.writePump(client)
}

// readPump 处理WebSocket读取
func (wsm *WebServerManager) readPump(client *WebSocketClient) {
	defer func() {
		client.Conn.Close()
		wsm.clientsMutex.Lock()
		delete(wsm.clients, client.ID)
		wsm.clientsMutex.Unlock()
		
		if wsm.uiInstance != nil {
			wsm.uiInstance.AddLog(fmt.Sprintf("WebSocket连接断开: %s", client.ID))
		}
	}()

	client.Conn.SetReadLimit(512)
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			break
		}

		// 广播消息给所有客户端
		wsm.broadcastMessage(message)
		
		if wsm.uiInstance != nil {
			wsm.uiInstance.AddLog(fmt.Sprintf("收到WebSocket消息 (%s): %s", client.ID, string(message)))
		}
	}
}

// writePump 处理WebSocket写入
func (wsm *WebServerManager) writePump(client *WebSocketClient) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.SendChan:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 添加队列中的其他消息
			n := len(client.SendChan)
			for i := 0; i < n; i++ {
				w.Write(<-client.SendChan)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// broadcastMessage 广播消息给所有WebSocket客户端
func (wsm *WebServerManager) broadcastMessage(message []byte) {
	wsm.clientsMutex.RLock()
	defer wsm.clientsMutex.RUnlock()

	for _, client := range wsm.clients {
		select {
		case client.SendChan <- message:
		default:
			close(client.SendChan)
			client.Conn.Close()
			delete(wsm.clients, client.ID)
		}
	}
}

// GetDefaultConfig 获取默认配置
// 返回:
//   - map[string]interface{}: 默认配置信息
func (wsm *WebServerManager) GetDefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"http_port": 8080,
		"ws_port":   8081,
		"enable_cors": true,
		"max_connections": 100,
	}
}

// SaveConfig 保存配置到文件
// 参数:
//   - filePath: 配置文件路径
// 返回:
//   - error: 操作错误，nil表示成功
func (wsm *WebServerManager) SaveConfig(filePath string) error {
	config := map[string]interface{}{
		"http_port": wsm.httpPort,
		"ws_port":   wsm.wsPort,
	}

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	return os.WriteFile(filePath, configJSON, 0644)
}

// LoadConfig 从文件加载配置
// 参数:
//   - filePath: 配置文件路径
// 返回:
//   - error: 操作错误，nil表示成功
func (wsm *WebServerManager) LoadConfig(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 应用配置
	if httpPort, ok := config["http_port"].(float64); ok {
		wsm.httpPort = int(httpPort)
	}

	if wsPort, ok := config["ws_port"].(float64); ok {
		wsm.wsPort = int(wsPort)
	}

	return nil
}