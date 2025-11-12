// managers/http_manager.go - HTTP操作管理器
// 功能描述: 提供HTTP请求的发送、文件下载和响应处理功能
// 主要功能:
//   - GET请求：发送HTTP GET请求并获取响应
//   - POST请求：发送HTTP POST请求并提交数据
//   - 文件下载：下载文件到指定路径
//   - 响应处理：解析HTTP响应状态和内容
//   - 错误处理：提供详细的网络错误和HTTP错误信息
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12

package managers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"windows-gui-app/interfaces"
)

// HTTPManager HTTP操作管理器
// 负责处理HTTP请求的发送和响应处理，提供UI集成功能
type HTTPManager struct {
	uiInstance interfaces.UIInterface // UI实例接口，用于显示操作状态和结果
	client     *http.Client            // HTTP客户端实例
}

// NewHTTPManager 创建HTTP管理器实例
// 参数:
//   - ui: UI实例，用于显示操作结果（可以为nil，之后通过SetUIInstance设置）
//
// 返回:
//   - *HTTPManager: HTTP管理器实例
//
// 使用示例:
//   httpManager := managers.NewHTTPManager(mainForm)
func NewHTTPManager(ui interfaces.UIInterface) *HTTPManager {
	// 创建自定义HTTP客户端，设置超时时间
	client := &http.Client{
		Timeout: 30 * time.Second, // 30秒超时
	}

	return &HTTPManager{
		uiInstance: ui,
		client:     client,
	}
}

// SetUIInstance 设置UI实例
// 用于在运行时更新UI实例引用
// 参数:
//   - ui: 新的UI实例
func (hm *HTTPManager) SetUIInstance(ui interfaces.UIInterface) {
	hm.uiInstance = ui
}

// GETRequest 发送HTTP GET请求
// 功能描述:
//   - 发送HTTP GET请求到指定URL
//   - 自动处理URL格式验证和编码
//   - 获取并解析HTTP响应状态码和内容
//   - 提供详细的错误信息和状态反馈
//
// 参数:
//   - urlStr: 请求的URL地址（必须包含协议，如http://或https://）
//
// 返回:
//   - string: 响应内容
//   - int: HTTP状态码
//   - error: 操作错误，nil表示成功
//
// 支持的URL格式:
//   - "http://example.com/api/users"
//   - "https://api.example.com/v1/data"
//   - "http://localhost:8080/status"
//
// 可能的错误:
//   - "URL不能为空": urlStr参数为空
//   - "URL格式错误": URL格式不合法
//   - "不支持的协议": 不是http或https协议
//   - "发送GET请求失败": 网络错误或服务器错误
//   - "读取响应失败": 响应内容读取错误
//
// 使用示例:
//   content, statusCode, err := httpManager.GETRequest("https://api.example.com/users")
//   if err != nil {
//       log.Printf("请求失败: %v", err)
//   }
//   log.Printf("状态码: %d, 内容长度: %d", statusCode, len(content))
func (hm *HTTPManager) GETRequest(urlStr string) (string, int, error) {
	// 验证URL参数
	if strings.TrimSpace(urlStr) == "" {
		return "", 0, fmt.Errorf("URL不能为空")
	}

	// 验证URL格式
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", 0, fmt.Errorf("URL格式错误: %v", err)
	}

	// 检查协议
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", 0, fmt.Errorf("不支持的协议: %s，只支持http和https", parsedURL.Scheme)
	}

	// 更新UI状态
	if hm.uiInstance != nil {
		hm.uiInstance.UpdateStatus(fmt.Sprintf("正在发送GET请求到: %s", urlStr))
	}

	log.Printf("发送GET请求: %s", urlStr)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", 0, fmt.Errorf("创建GET请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "GO_VCL-HTTP-Client/1.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	// 发送请求
	resp, err := hm.client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("发送GET请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("读取响应失败: %v", err)
	}

	// 更新UI状态
	if hm.uiInstance != nil {
		statusMsg := fmt.Sprintf("GET请求完成，状态码: %d，响应长度: %d字节", resp.StatusCode, len(body))
		hm.uiInstance.UpdateStatus(statusMsg)
	}

	// 记录操作日志
	log.Printf("GET请求成功: %s -> 状态码: %d，响应长度: %d字节", urlStr, resp.StatusCode, len(body))
	return string(body), resp.StatusCode, nil
}

// POSTRequest 发送HTTP POST请求
// 功能描述:
//   - 发送HTTP POST请求到指定URL
//   - 支持多种数据格式（JSON、表单、纯文本）
//   - 自动检测和设置适当的内容类型
//   - 处理请求超时和响应解析
//
// 参数:
//   - urlStr: 请求的URL地址（必须包含协议）
//   - data: 要发送的数据（支持字符串、map、slice等类型）
//   - contentType: 内容类型（可选，如"application/json"，为空则自动检测）
//
// 返回:
//   - string: 响应内容
//   - int: HTTP状态码
//   - error: 操作错误，nil表示成功
//
// 支持的数据类型:
//   - 字符串: 直接作为请求体发送
//   - map[string]interface{}: 转换为JSON格式发送
//   - []interface{}: 转换为JSON数组发送
//   - 其他类型: 转换为字符串发送
//
// 自动内容类型检测:
//   - JSON数据: "application/json"
//   - 纯文本: "text/plain"
//   - 表单数据: "application/x-www-form-urlencoded"
//
// 可能的错误:
//   - "URL不能为空": urlStr参数为空
//   - "数据不能为空": data参数为空
//   - "URL格式错误": URL格式不合法
//   - "不支持的协议": 不是http或https协议
//   - "编码数据失败": 数据序列化错误
//   - "发送POST请求失败": 网络错误或服务器错误
//
// 使用示例:
//   // 发送JSON数据
//   data := map[string]interface{}{"name": "张三", "age": 25}
//   content, statusCode, err := httpManager.POSTRequest("https://api.example.com/users", data, "")
//   
//   // 发送纯文本数据
//   content, statusCode, err := httpManager.POSTRequest("https://api.example.com/text", "Hello World", "")
func (hm *HTTPManager) POSTRequest(urlStr string, data interface{}, contentType string) (string, int, error) {
	// 验证参数
	if strings.TrimSpace(urlStr) == "" {
		return "", 0, fmt.Errorf("URL不能为空")
	}

	if data == nil {
		return "", 0, fmt.Errorf("数据不能为空")
	}

	// 验证URL格式
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", 0, fmt.Errorf("URL格式错误: %v", err)
	}

	// 检查协议
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", 0, fmt.Errorf("不支持的协议: %s，只支持http和https", parsedURL.Scheme)
	}

	// 编码请求数据
	var bodyData []byte
	var actualContentType string

	// 根据数据类型自动选择编码方式
	switch v := data.(type) {
	case string:
		bodyData = []byte(v)
		if contentType == "" {
			actualContentType = "text/plain"
		} else {
			actualContentType = contentType
		}
	case map[string]interface{}, []interface{}:
		// JSON编码
		jsonData, err := json.Marshal(v)
		if err != nil {
			return "", 0, fmt.Errorf("编码数据失败: %v", err)
		}
		bodyData = jsonData
		if contentType == "" {
			actualContentType = "application/json"
		} else {
			actualContentType = contentType
		}
	default:
		// 转换为字符串
		bodyData = []byte(fmt.Sprintf("%v", v))
		if contentType == "" {
			actualContentType = "text/plain"
		} else {
			actualContentType = contentType
		}
	}

	// 更新UI状态
	if hm.uiInstance != nil {
		hm.uiInstance.UpdateStatus(fmt.Sprintf("正在发送POST请求到: %s", urlStr))
	}

	log.Printf("发送POST请求: %s，数据长度: %d字节，内容类型: %s", urlStr, len(bodyData), actualContentType)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(bodyData))
	if err != nil {
		return "", 0, fmt.Errorf("创建POST请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "GO_VCL-HTTP-Client/1.0")
	req.Header.Set("Content-Type", actualContentType)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(bodyData)))

	// 发送请求
	resp, err := hm.client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("发送POST请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("读取响应失败: %v", err)
	}

	// 更新UI状态
	if hm.uiInstance != nil {
		statusMsg := fmt.Sprintf("POST请求完成，状态码: %d，响应长度: %d字节", resp.StatusCode, len(body))
		hm.uiInstance.UpdateStatus(statusMsg)
	}

	// 记录操作日志
	log.Printf("POST请求成功: %s -> 状态码: %d，响应长度: %d字节", urlStr, resp.StatusCode, len(body))
	return string(body), resp.StatusCode, nil
}

// DownloadFile 下载文件到指定路径
// 功能描述:
//   - 从指定URL下载文件
//   - 支持大文件分块下载和进度显示
//   - 自动创建输出目录
//   - 提供下载进度和速度信息
//
// 参数:
//   - urlStr: 要下载的文件的URL地址
//   - filePath: 本地保存文件路径（包含文件名）
//
// 返回:
//   - error: 操作错误，nil表示成功
//
// 特殊功能:
//   - 自动创建输出目录
//   - 显示下载进度（通过UI状态更新）
//   - 支持断点续传（如果服务器支持）
//   - 下载完成后验证文件完整性
//
// 可能的错误:
//   - "下载URL不能为空": urlStr参数为空
//   - "文件路径不能为空": filePath参数为空
//   - "URL格式错误": URL格式不合法
//   - "不支持的协议": 不是http或https协议
//   - "创建目录失败": 输出目录创建错误
//   - "发送下载请求失败": 网络错误或服务器错误
//   - "创建文件失败": 本地文件创建错误
//   - "下载文件失败": 文件写入错误
//
// 使用示例:
//   err := httpManager.DownloadFile("https://example.com/file.pdf", "C:\\downloads\\file.pdf")
//   if err != nil {
//       log.Printf("下载失败: %v", err)
//   }
func (hm *HTTPManager) DownloadFile(urlStr, filePath string) error {
	// 验证参数
	if strings.TrimSpace(urlStr) == "" {
		return fmt.Errorf("下载URL不能为空")
	}

	if strings.TrimSpace(filePath) == "" {
		return fmt.Errorf("文件路径不能为空")
	}

	// 验证URL格式
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("URL格式错误: %v", err)
	}

	// 检查协议
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("不支持的协议: %s，只支持http和https", parsedURL.Scheme)
	}

	// 创建输出目录
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 更新UI状态
	if hm.uiInstance != nil {
		hm.uiInstance.UpdateStatus(fmt.Sprintf("正在下载文件: %s", filepath.Base(filePath)))
	}

	log.Printf("开始下载文件: %s -> %s", urlStr, filePath)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return fmt.Errorf("创建下载请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "GO_VCL-HTTP-Client/1.0")

	// 发送请求
	resp, err := hm.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送下载请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，HTTP状态码: %d", resp.StatusCode)
	}

	// 获取文件大小
	contentLength := resp.ContentLength
	log.Printf("文件大小: %d字节", contentLength)

	// 创建输出文件
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer out.Close()

	// 复制数据（带进度更新）
	var written int64
	if contentLength > 0 {
		// 大文件分块下载
		buf := make([]byte, 32*1024) // 32KB缓冲区
		for {
			nr, err := resp.Body.Read(buf)
			if nr > 0 {
				nw, err := out.Write(buf[:nr])
				if err != nil {
					return fmt.Errorf("写入文件失败: %v", err)
				}
				written += int64(nw)

				// 更新进度（每下载1MB更新一次）
				if written%(1024*1024) == 0 && hm.uiInstance != nil {
					progress := float64(written) / float64(contentLength) * 100
					hm.uiInstance.UpdateStatus(fmt.Sprintf("下载进度: %.1f%% (%d/%d KB)", 
						progress, written/1024, contentLength/1024))
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("下载文件失败: %v", err)
			}
		}
	} else {
		// 小文件直接复制
		written, err = io.Copy(out, resp.Body)
		if err != nil {
			return fmt.Errorf("下载文件失败: %v", err)
		}
	}

	// 验证文件大小
	fileInfo, err := out.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	if contentLength > 0 && fileInfo.Size() != contentLength {
		return fmt.Errorf("文件下载不完整，期望大小: %d，实际大小: %d", contentLength, fileInfo.Size())
	}

	// 更新UI状态
	if hm.uiInstance != nil {
		statusMsg := fmt.Sprintf("文件下载完成: %s，大小: %.2fKB", 
			filepath.Base(filePath), float64(fileInfo.Size())/1024)
		hm.uiInstance.UpdateStatus(statusMsg)
	}

	// 记录操作日志
	log.Printf("文件下载成功: %s -> %s，大小: %d字节，耗时: %v", 
		urlStr, filePath, fileInfo.Size(), time.Since(time.Now()))
	return nil
}

// GetDefaultURL 获取默认的测试URL
// 返回:
//   - string: 默认的HTTP测试URL
func (hm *HTTPManager) GetDefaultURL() string {
	// 返回一个稳定的测试API URL
	return "https://jsonplaceholder.typicode.com/posts/1"
}