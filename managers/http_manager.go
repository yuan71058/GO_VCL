// managers/http_manager.go - HTTP操作管理器
// 该文件实现了HTTP操作管理器，提供HTTP请求、文件下载、API测试等功能
// 支持GET、POST等HTTP方法，并提供了丰富的辅助函数用于处理HTTP请求和响应
package managers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"windows-gui-app/interfaces"
)

// HTTPManager HTTP操作管理器
// 封装了HTTP客户端，提供统一的HTTP操作接口，支持GET、POST请求和文件下载
// 与UI组件交互，将操作结果展示在用户界面上
type HTTPManager struct {
	uiInstance interfaces.UIInterface // UI实例接口，用于更新界面状态和数据
	client     *resty.Client          // HTTP客户端实例，用于执行HTTP请求
}

// NewHTTPManager 创建HTTP管理器实例
// 初始化HTTP客户端，设置默认超时时间，并关联UI实例
// 参数:
//   - ui: UI实例，用于显示操作结果（可以为nil，之后通过SetUIInstance设置）
//
// 返回:
//   - *HTTPManager: HTTP管理器实例
func NewHTTPManager(ui interfaces.UIInterface) *HTTPManager {
	client := resty.New()
	// 设置默认超时时间为30秒，避免长时间等待
	client.SetTimeout(30 * time.Second)

	return &HTTPManager{
		uiInstance: ui,
		client:     client,
	}
}

// SetUIInstance 设置UI实例
// 允许在创建管理器后设置UI实例，用于解耦合UI和管理器的创建顺序
// 参数:
//   - ui: UI实例接口
func (hm *HTTPManager) SetUIInstance(ui interfaces.UIInterface) {
	hm.uiInstance = ui
}

// GETRequest 发起HTTP GET请求
// 使用封装的HTTP客户端发送GET请求，返回原始HTTP响应
// 参数:
//   - url: 请求的URL地址
//
// 返回:
//   - *http.Response: HTTP响应对象
//   - error: 错误信息，nil表示请求成功
func (hm *HTTPManager) GETRequest(url string) (*http.Response, error) {
	resp, err := hm.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	return resp.RawResponse, nil
}

// POSTRequest 执行HTTP POST请求
// 支持多种数据类型，自动处理JSON序列化和表单数据格式
// 将请求结果展示在UI表格中，并更新状态栏
// 参数:
//   - urlStr: 请求URL
//   - data: POST数据，支持字符串、map或[]byte
//   - contentType: 内容类型，如"application/json"
//
// 返回:
//   - error: 操作错误，nil表示成功
func (hm *HTTPManager) POSTRequest(urlStr string, data interface{}, contentType string) error {
	if strings.TrimSpace(urlStr) == "" {
		return fmt.Errorf("URL不能为空")
	}

	// 验证URL格式，确保URL符合HTTP规范
	if _, err := url.ParseRequestURI(urlStr); err != nil {
		return fmt.Errorf("URL格式错误: %v", err)
	}

	// 更新UI状态，提示用户正在发送请求
	hm.uiInstance.UpdateStatus("正在发送POST请求...")

	// 设置POST数据 - 准备请求体
	var reqBody interface{}
	var useFormData bool

	switch v := data.(type) {
	case string:
		reqBody = v
	case map[string]string:
		reqBody = v
		useFormData = true // map[string]string使用表单数据格式
	case map[string]interface{}:
		jsonStr, err := convertToJSONString(v)
		if err != nil {
			return fmt.Errorf("JSON序列化失败: %v", err)
		}
		reqBody = jsonStr
	case []byte:
		reqBody = v
	default:
		return fmt.Errorf("不支持的数据类型: %T", v)
	}

	// 发送请求
	var resp *resty.Response
	var err error

	if useFormData {
		// 使用表单数据格式
		resp, err = hm.client.R().
			SetFormData(data.(map[string]string)).
			Post(urlStr)
	} else {
		// 使用普通请求体
		resp, err = hm.client.R().
			SetHeader("Content-Type", func() string {
				if contentType != "" {
					return contentType
				}
				// 根据数据类型设置默认的内容类型
				if str, ok := data.(string); ok && isJSON(str) {
					return "application/json"
				}
				return "application/x-www-form-urlencoded"
			}()).
			SetBody(reqBody).
			Post(urlStr)
	}
	if err != nil {
		return fmt.Errorf("POST请求失败: %v", err)
	}

	// 检查响应状态
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return fmt.Errorf("服务器响应错误，状态码: %d", resp.StatusCode())
	}

	// 获取响应内容
	body := resp.String()

	// 准备显示数据
	var displayData [][]string
	displayData = append(displayData, []string{"字段", "值"})
	displayData = append(displayData, []string{"状态码", fmt.Sprintf("%d", resp.StatusCode())})
	displayData = append(displayData, []string{"URL", urlStr})
	displayData = append(displayData, []string{"请求时间", time.Now().Format("2006-01-02 15:04:05")})
	displayData = append(displayData, []string{"请求数据大小", fmt.Sprintf("%d字节", getDataSize(data))})
	displayData = append(displayData, []string{"响应大小", fmt.Sprintf("%d字节", len(body))})

	// 添加内容类型信息
	if contentType != "" {
		displayData = append(displayData, []string{"请求类型", contentType})
	}

	// 响应内容预览
	if len(body) > 0 {
		if isJSONResponse(body) {
			displayData = append(displayData, []string{"响应类型", "JSON"})
		} else if isHTMLResponse(body) {
			displayData = append(displayData, []string{"响应类型", "HTML"})
		} else {
			displayData = append(displayData, []string{"响应类型", "纯文本"})
		}
		displayData = append(displayData, []string{"响应内容", truncateString(body, 200)})
	}

	// 更新UI表格显示
	hm.uiInstance.SetTableData(displayData)

	// 更新状态
	hm.uiInstance.UpdateStatus(fmt.Sprintf("POST请求成功，状态码: %d", resp.StatusCode()))

	log.Printf("成功执行POST请求: %s，状态码: %d", urlStr, resp.StatusCode())
	return nil
}

// DownloadFile 下载文件
// 从指定URL下载文件并保存到本地路径，自动创建必要的目录结构
// 将下载结果展示在UI表格中，并更新状态栏
// 参数:
//   - urlStr: 下载URL
//   - filePath: 本地保存路径
//
// 返回:
//   - error: 操作错误，nil表示成功
func (hm *HTTPManager) DownloadFile(urlStr, filePath string) error {
	if strings.TrimSpace(urlStr) == "" {
		return fmt.Errorf("下载URL不能为空")
	}

	// 验证URL格式，确保URL符合HTTP规范
	if _, err := url.ParseRequestURI(urlStr); err != nil {
		return fmt.Errorf("URL格式错误: %v", err)
	}

	// 更新UI状态，提示用户正在下载文件
	hm.uiInstance.UpdateStatus("正在下载文件...")

	// 发送GET请求下载文件
	resp, err := hm.client.R().Get(urlStr)
	if err != nil {
		return fmt.Errorf("下载文件失败: %v", err)
	}

	// 检查响应状态
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("下载请求失败，状态码: %d", resp.StatusCode())
	}

	// 获取文件内容
	fileData := resp.Body()

	// 保存文件
	if err := hm.saveFile(filePath, fileData); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	// 更新UI
	displayData := [][]string{
		[]string{"字段", "值"},
		[]string{"下载URL", urlStr},
		[]string{"保存路径", filePath},
		[]string{"文件大小", fmt.Sprintf("%.2fKB", float64(len(fileData))/1024)},
		[]string{"下载时间", time.Now().Format("2006-01-02 15:04:05")},
	}
	hm.uiInstance.SetTableData(displayData)

	// 更新状态
	hm.uiInstance.UpdateStatus(fmt.Sprintf("文件下载成功: %.2fKB", float64(len(fileData))/1024))

	log.Printf("成功下载文件: %s，大小: %d字节", urlStr, len(fileData))
	return nil
}

// TestAPI 测试公共API接口
// 使用JSONPlaceholder API进行测试，验证HTTP客户端的基本功能
// 将测试结果展示在UI表格中，并更新状态栏
// 返回:
//   - error: 操作错误，nil表示成功
func (hm *HTTPManager) TestAPI() error {
	// 更新UI状态，提示用户正在测试API
	hm.uiInstance.UpdateStatus("正在测试公共API...")

	// 测试JSONPlaceholder API
	testURL := "https://jsonplaceholder.typicode.com/posts/1"

	// 发送GET请求
	resp, err := hm.client.R().Get(testURL)
	if err != nil {
		return fmt.Errorf("API测试失败: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("API响应错误，状态码: %d", resp.StatusCode())
	}

	// 解析响应
	body := resp.String()

	// 准备显示数据
	var data [][]string
	data = append(data, []string{"API测试项目", "结果"})
	data = append(data, []string{"测试URL", testURL})
	data = append(data, []string{"状态码", fmt.Sprintf("%d", resp.StatusCode())})
	data = append(data, []string{"响应时间", time.Now().Format("2006-01-02 15:04:05")})
	data = append(data, []string{"响应数据", truncateString(body, 300)})

	// 更新UI
	hm.uiInstance.SetTableData(data)
	hm.uiInstance.UpdateStatus("API测试成功")

	log.Printf("API测试成功: %s", testURL)
	return nil
}

// GetRequestInfo 获取请求详细信息
// 将HTTP请求的各个组成部分整理为键值对，便于展示和调试
// 参数:
//   - method: HTTP方法
//   - urlStr: 请求URL
//   - headers: 请求头
//   - data: 请求数据
//
// 返回:
//   - map[string]string: 请求信息键值对
func (hm *HTTPManager) GetRequestInfo(method, urlStr string, headers map[string]string, data interface{}) map[string]string {
	info := make(map[string]string)

	info["方法"] = method
	info["URL"] = urlStr
	info["时间"] = time.Now().Format("2006-01-02 15:04:05")

	if headers != nil {
		info["头部数量"] = fmt.Sprintf("%d个", len(headers))
	}

	if data != nil {
		info["数据大小"] = fmt.Sprintf("%d字节", getDataSize(data))
	}

	return info
}

// SetTimeout 设置请求超时时间
// 动态调整HTTP客户端的超时时间，适应不同网络环境
// 参数:
//   - timeout: 超时时间
func (hm *HTTPManager) SetTimeout(timeout time.Duration) {
	hm.client.SetTimeout(timeout)
	hm.uiInstance.UpdateStatus(fmt.Sprintf("设置请求超时时间: %v", timeout))
}

// GetClient 获取HTTP客户端实例
// 返回底层的HTTP客户端实例，允许进行更高级的配置
// 返回:
//   - *resty.Client: HTTP客户端实例
func (hm *HTTPManager) GetClient() *resty.Client {
	return hm.client
}

// 辅助函数

// isJSONResponse 检查响应是否为JSON格式
// 通过检查响应内容的前缀字符判断是否为JSON格式
// 参数:
//   - body: 响应内容字符串
//
// 返回:
//   - bool: 是否为JSON格式
func isJSONResponse(body string) bool {
	lowerBody := strings.ToLower(body)
	return strings.HasPrefix(strings.TrimSpace(lowerBody), "{") ||
		strings.HasPrefix(strings.TrimSpace(lowerBody), "[")
}

// isHTMLResponse 检查响应是否为HTML格式
// 通过检查响应内容是否包含HTML标签判断是否为HTML格式
// 参数:
//   - body: 响应内容字符串
//
// 返回:
//   - bool: 是否为HTML格式
func isHTMLResponse(body string) bool {
	lowerBody := strings.ToLower(body)
	return strings.Contains(lowerBody, "<html") ||
		strings.Contains(lowerBody, "<!doctype")
}

// isJSON 检查字符串是否为JSON格式
// 通过检查字符串的第一个字符判断是否为JSON格式
// 参数:
//   - str: 待检查的字符串
//
// 返回:
//   - bool: 是否为JSON格式
func isJSON(str string) bool {
	if strings.TrimSpace(str) == "" {
		return false
	}
	firstChar := strings.TrimSpace(str)[0]
	return firstChar == '{' || firstChar == '['
}

// truncateString 截断字符串
// 当字符串长度超过指定限制时，截断并添加省略号
// 参数:
//   - str: 待截断的字符串
//   - maxLen: 最大长度
//
// 返回:
//   - string: 截断后的字符串
func truncateString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	return str[:maxLen] + "..."
}

// convertToJSONString 将map转换为JSON字符串
// 将map[string]interface{}类型的数据序列化为JSON字符串
// 参数:
//   - data: 待转换的数据
//
// 返回:
//   - string: JSON字符串
//   - error: 错误信息
func convertToJSONString(data interface{}) (string, error) {
	// 转换为JSON字符串
	switch v := data.(type) {
	case map[string]interface{}:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("JSON序列化失败: %v", err)
		}
		return string(jsonBytes), nil
	default:
		return "", fmt.Errorf("不支持的数据类型: %T", v)
	}
}

// getDataSize 获取数据大小
// 计算不同类型数据的字节大小，用于显示和统计
// 参数:
//   - data: 待计算的数据
//
// 返回:
//   - int: 数据大小（字节）
func getDataSize(data interface{}) int {
	switch v := data.(type) {
	case string:
		return len([]byte(v))
	case []byte:
		return len(v)
	case map[string]string:
		return len(v)
	case map[string]interface{}:
		return len(v)
	default:
		return 0
	}
}

// saveFile 保存文件到本地
// 将字节数据保存到指定文件路径，自动创建必要的目录结构
// 参数:
//   - filePath: 文件路径
//   - data: 文件数据
//
// 返回:
//   - error: 错误信息
func (hm *HTTPManager) saveFile(filePath string, data []byte) error {
	// 创建目录
	dir := getDirectory(filePath)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %v", err)
		}
	}

	// 写入文件
	return os.WriteFile(filePath, data, 0644)
}

// getDirectory 获取文件目录
// 从完整文件路径中提取目录部分
// 参数:
//   - filePath: 文件路径
//
// 返回:
//   - string: 目录路径
func getDirectory(filePath string) string {
	return filepath.Dir(filePath)
}