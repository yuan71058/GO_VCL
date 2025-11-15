# Web服务器功能使用指南

## 概述

本程序已添加了Web服务器功能，支持启动HTTP和WebSocket服务。用户可以通过图形界面轻松控制Web服务器的启动和停止。Web服务器提供RESTful API接口，支持实时双向通信，适用于构建Web应用和实时数据推送场景。

## 功能特点

1. **HTTP服务**：提供RESTful API接口
2. **WebSocket服务**：支持实时双向通信
3. **客户端管理**：跟踪和管理连接的客户端
4. **API路由**：内置多个API端点，包括状态查询、客户端列表等
5. **实时数据推送**：支持向客户端推送实时数据
6. **消息队列管理**：提供高效的消息队列机制
7. **跨域支持**：支持CORS，便于前端应用集成

## 使用方法

1. 启动程序后，在主界面中找到"🌐 Web服务器"按钮（位于第四排按钮中）
2. 点击按钮启动Web服务器
3. 服务器启动后，按钮文本将变为"⏹️ 停止服务器"
4. 表格中的"Web服务器"行状态将更新为"运行中"
5. 日志区域将显示服务器启动信息和服务器详情
6. 再次点击按钮可停止服务器

## API端点

服务器启动后，将提供以下API端点：

### HTTP API
- `GET /` - 服务器状态信息
- `GET /api/status` - 服务器状态
- `GET /api/clients` - 连接的客户端列表
- `POST /api/clients/{id}/message` - 向指定客户端发送消息
- `POST /api/broadcast` - 向所有客户端广播消息
- `GET /api/stats` - 服务器统计信息

### WebSocket API
- `WebSocket /ws` - WebSocket连接端点
- 支持以下消息类型：
  - `ping` - 心跳检测
  - `message` - 文本消息
  - `data` - JSON数据
  - `broadcast` - 广播消息

## 默认配置

- HTTP端口：8080
- WebSocket端口：8081
- 服务器地址：localhost

## 注意事项

1. 确保端口8080和8081未被其他程序占用
2. 防火墙可能会阻止外部访问，如需远程访问请配置防火墙规则
3. 服务器日志将显示在程序日志区域中

## 技术实现

Web服务器功能使用以下技术实现：

- HTTP服务：基于标准库net/http
- WebSocket服务：使用gorilla/websocket库
- 客户端管理：内置客户端连接池
- API路由：使用标准库mux路由器
- 消息队列：使用channel实现高效消息传递
- 并发处理：使用goroutine处理并发请求
- 跨域处理：使用CORS中间件支持跨域请求

## 使用示例

### JavaScript客户端示例
```javascript
// 连接WebSocket
const ws = new WebSocket('ws://localhost:8081/ws');

// 接收消息
ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('收到消息:', data);
};

// 发送消息
ws.send(JSON.stringify({
    type: 'message',
    content: 'Hello Server'
}));

// 心跳检测
setInterval(() => {
    ws.send(JSON.stringify({type: 'ping'}));
}, 30000);
```

### HTTP API调用示例
```bash
# 获取服务器状态
curl http://localhost:8080/api/status

# 获取客户端列表
curl http://localhost:8080/api/clients

# 向指定客户端发送消息
curl -X POST http://localhost:8080/api/clients/client1/message \
     -H "Content-Type: application/json" \
     -d '{"message": "Hello Client"}'

# 广播消息
curl -X POST http://localhost:8080/api/broadcast \
     -H "Content-Type: application/json" \
     -d '{"message": "Broadcast Message"}'
```

## 故障排除

如果遇到问题，请检查：

1. **端口占用**：
   - 检查端口8080和8081是否被其他程序占用
   - 使用`netstat -ano | findstr :8080`命令检查端口状态
   - 如有占用，可修改代码中的默认端口

2. **防火墙设置**：
   - 确保防火墙允许端口8080和8081的入站连接
   - 在Windows防火墙中添加例外规则

3. **程序日志**：
   - 查看程序日志中的错误信息
   - 检查服务器启动是否成功

4. **权限问题**：
   - 确保程序有足够的权限绑定端口
   - 以管理员身份运行程序

5. **WebSocket连接问题**：
   - 检查浏览器控制台是否有错误信息
   - 确认WebSocket URL格式正确
   - 检查是否使用了正确的协议(ws://或wss://)

## 更新日志

### 2025-11-14
- 完善Web服务器功能，支持HTTP和WebSocket服务
- 实现RESTful API接口，支持状态查询和客户端管理
- 添加WebSocket实时双向通信功能
- 实现实时数据推送机制
- 添加消息队列管理，提高通信效率
- 添加CORS支持，便于前端应用集成
- 完善API文档，添加详细的使用示例

### 2025-11-12
- 初始版本发布
- 实现基本Web服务器功能
- 支持HTTP和WebSocket服务
- 提供基础API端点

## 扩展功能

可以考虑添加的功能：

1. **HTTPS支持**：添加SSL/TLS加密支持
2. **用户认证**：实现基于JWT的身份验证
3. **API限流**：添加请求频率限制
4. **数据持久化**：支持数据持久化存储
5. **负载均衡**：支持多实例负载均衡
6. **监控面板**：添加Web监控面板
7. **API文档**：集成Swagger API文档
8. **日志管理**：实现结构化日志记录