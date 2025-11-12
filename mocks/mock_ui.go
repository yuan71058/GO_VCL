// mocks/mock_ui.go - UI模拟接口
// 功能描述: 为测试提供模拟的UI接口
// 主要功能:
//   - 模拟UI状态更新
//   - 模拟进度条显示
//   - 模拟消息显示
//   - 模拟错误提示
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12

package mocks

import (
	"fmt"
	"log"
	"time"
)

// MockUI 模拟UI接口
// 功能描述: 提供测试用的UI接口模拟
// 字段说明:
//   - statusUpdates: 状态更新记录
//   - progressUpdates: 进度更新记录
//   - messages: 消息记录
//   - errors: 错误记录
//   - lastStatus: 最后状态
//   - lastProgress: 最后进度
//   - isProcessing: 处理状态
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
type MockUI struct {
	statusUpdates   []string
	progressUpdates []int
	messages        []string
	errors          []string
	lastStatus      string
	lastProgress    int
	isProcessing    bool
}

// NewMockUI 创建新的模拟UI实例
// 功能描述: 创建并初始化模拟UI接口
// 返回值:
//   - *MockUI: 模拟UI实例
// 使用示例:
//   mockUI := mocks.NewMockUI()
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func NewMockUI() *MockUI {
	return &MockUI{
		statusUpdates:   make([]string, 0),
		progressUpdates: make([]int, 0),
		messages:        make([]string, 0),
		errors:          make([]string, 0),
		lastStatus:      "Ready",
		lastProgress:    0,
		isProcessing:    false,
	}
}

// SetStatus 设置状态信息
// 功能描述: 模拟设置状态栏信息
// 参数:
//   - status: 状态文本
// 返回值:
//   - error: 错误信息（模拟接口始终返回nil）
// 使用示例:
//   err := mockUI.SetStatus("Processing...")
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) SetStatus(status string) error {
	m.statusUpdates = append(m.statusUpdates, status)
	m.lastStatus = status
	log.Printf("[MockUI] Status: %s", status)
	return nil
}

// SetProgress 设置进度
// 功能描述: 模拟设置进度条进度
// 参数:
//   - progress: 进度值（0-100）
// 返回值:
//   - error: 错误信息（模拟接口始终返回nil）
// 使用示例:
//   err := mockUI.SetProgress(50)
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) SetProgress(progress int) error {
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	
	m.progressUpdates = append(m.progressUpdates, progress)
	m.lastProgress = progress
	log.Printf("[MockUI] Progress: %d%%", progress)
	return nil
}

// ShowMessage 显示消息
// 功能描述: 模拟显示消息对话框
// 参数:
//   - message: 消息文本
//   - title: 消息标题
// 返回值:
//   - error: 错误信息（模拟接口始终返回nil）
// 使用示例:
//   err := mockUI.ShowMessage("操作完成", "成功")
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) ShowMessage(message, title string) error {
	fullMessage := fmt.Sprintf("[%s] %s", title, message)
	m.messages = append(m.messages, fullMessage)
	log.Printf("[MockUI] Message: %s", fullMessage)
	return nil
}

// ShowError 显示错误
// 功能描述: 模拟显示错误对话框
// 参数:
//   - err: 错误对象
//   - title: 错误标题
// 返回值:
//   - error: 错误信息（模拟接口始终返回nil）
// 使用示例:
//   err := mockUI.ShowError(errors.New("操作失败"), "错误")
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) ShowError(err error, title string) error {
	if err != nil {
		errorMsg := fmt.Sprintf("[%s] %v", title, err)
		m.errors = append(m.errors, errorMsg)
		log.Printf("[MockUI] Error: %s", errorMsg)
	}
	return nil
}

// SetProcessing 设置处理状态
// 功能描述: 模拟设置处理状态
// 参数:
//   - processing: 是否正在处理
// 返回值:
//   - error: 错误信息（模拟接口始终返回nil）
// 使用示例:
//   err := mockUI.SetProcessing(true)
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) SetProcessing(processing bool) error {
	m.isProcessing = processing
	status := "Idle"
	if processing {
		status = "Processing"
	}
	log.Printf("[MockUI] Processing: %s", status)
	return nil
}

// GetStatusHistory 获取状态历史
// 功能描述: 获取所有状态更新记录
// 返回值:
//   - []string: 状态更新历史列表
// 使用示例:
//   history := mockUI.GetStatusHistory()
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) GetStatusHistory() []string {
	return append([]string{}, m.statusUpdates...)
}

// GetProgressHistory 获取进度历史
// 功能描述: 获取所有进度更新记录
// 返回值:
//   - []int: 进度更新历史列表
// 使用示例:
//   history := mockUI.GetProgressHistory()
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) GetProgressHistory() []int {
	return append([]int{}, m.progressUpdates...)
}

// GetMessages 获取消息记录
// 功能描述: 获取所有消息记录
// 返回值:
//   - []string: 消息记录列表
// 使用示例:
//   messages := mockUI.GetMessages()
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) GetMessages() []string {
	return append([]string{}, m.messages...)
}

// GetErrors 获取错误记录
// 功能描述: 获取所有错误记录
// 返回值:
//   - []string: 错误记录列表
// 使用示例:
//   errors := mockUI.GetErrors()
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) GetErrors() []string {
	return append([]string{}, m.errors...)
}

// GetLastStatus 获取最后状态
// 功能描述: 获取最后设置的状态
// 返回值:
//   - string: 最后状态文本
// 使用示例:
//   status := mockUI.GetLastStatus()
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) GetLastStatus() string {
	return m.lastStatus
}

// GetLastProgress 获取最后进度
// 功能描述: 获取最后设置的进度
// 返回值:
//   - int: 最后进度值
// 使用示例:
//   progress := mockUI.GetLastProgress()
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) GetLastProgress() int {
	return m.lastProgress
}

// IsCurrentlyProcessing 获取当前处理状态
// 功能描述: 获取当前是否正在处理
// 返回值:
//   - bool: 处理状态
// 使用示例:
//   processing := mockUI.IsCurrentlyProcessing()
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) IsCurrentlyProcessing() bool {
	return m.isProcessing
}

// Reset 重置模拟UI状态
// 功能描述: 清除所有历史记录并重置状态
// 使用示例:
//   mockUI.Reset()
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) Reset() {
	m.statusUpdates = make([]string, 0)
	m.progressUpdates = make([]int, 0)
	m.messages = make([]string, 0)
	m.errors = make([]string, 0)
	m.lastStatus = "Ready"
	m.lastProgress = 0
	m.isProcessing = false
	log.Println("[MockUI] Reset completed")
}

// SimulateOperation 模拟操作过程
// 功能描述: 模拟一个完整的操作过程，包括状态更新和进度变化
// 参数:
//   - operationName: 操作名称
//   - duration: 模拟持续时间（秒）
//   - success: 是否模拟成功
// 使用示例:
//   err := mockUI.SimulateOperation("File Import", 2, true)
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) SimulateOperation(operationName string, duration int, success bool) error {
	log.Printf("[MockUI] Starting simulated operation: %s", operationName)
	
	// 开始操作
	m.SetProcessing(true)
	m.SetStatus(fmt.Sprintf("Starting %s...", operationName))
	m.SetProgress(0)
	
	// 模拟操作过程
	steps := duration * 2 // 每500ms更新一次
	for i := 1; i <= steps; i++ {
		progress := (i * 100) / steps
		m.SetProgress(progress)
		m.SetStatus(fmt.Sprintf("%s in progress... %d%%", operationName, progress))
		time.Sleep(500 * time.Millisecond)
	}
	
	// 完成操作
	m.SetProcessing(false)
	m.SetProgress(100)
	
	if success {
		m.SetStatus(fmt.Sprintf("%s completed successfully", operationName))
		m.ShowMessage(fmt.Sprintf("%s operation completed", operationName), "Success")
	} else {
		m.SetStatus(fmt.Sprintf("%s failed", operationName))
		m.ShowError(fmt.Errorf("%s operation failed", operationName), "Error")
	}
	
	log.Printf("[MockUI] Completed simulated operation: %s (success: %v)", operationName, success)
	return nil
}

// GetStatistics 获取统计信息
// 功能描述: 获取UI操作的统计信息
// 返回值:
//   - map[string]interface{}: 统计信息映射
// 使用示例:
//   stats := mockUI.GetStatistics()
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) GetStatistics() map[string]interface{} {
	return map[string]interface{}{
		"total_status_updates": len(m.statusUpdates),
		"total_progress_updates": len(m.progressUpdates),
		"total_messages":         len(m.messages),
		"total_errors":           len(m.errors),
		"last_status":            m.lastStatus,
		"last_progress":          m.lastProgress,
		"is_processing":          m.isProcessing,
		"current_time":           time.Now().Format("2006-01-02 15:04:05"),
	}
}

// PrintStatistics 打印统计信息
// 功能描述: 打印详细的统计信息到日志
// 使用示例:
//   mockUI.PrintStatistics()
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) PrintStatistics() {
	stats := m.GetStatistics()
	log.Println("=== MockUI Statistics ===")
	for key, value := range stats {
		log.Printf("%s: %v", key, value)
	}
	log.Println("=======================")
}

// AddLog 添加日志信息
// 功能描述: 模拟添加日志信息到界面
// 参数:
//   - logText: 日志文本内容
// 使用示例:
//   mockUI.AddLog("操作开始")
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) AddLog(logText string) {
	log.Printf("[MockUI] Log: %s", logText)
}

// UpdateStatus 更新状态栏显示
// 功能描述: 模拟更新状态栏显示
// 参数:
//   - status: 状态文本
// 使用示例:
//   mockUI.UpdateStatus("就绪")
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) UpdateStatus(status string) {
	m.SetStatus(status)
}

// SetTableData 设置表格数据
// 功能描述: 模拟设置表格数据
// 参数:
//   - data: 二维字符串数组表示的表格数据
// 使用示例:
//   mockUI.SetTableData([][]string{{"A1", "B1"}, {"A2", "B2"}})
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) SetTableData(data [][]string) {
	log.Printf("[MockUI] Table data set: %d rows", len(data))
}

// GetTableData 获取表格数据
// 功能描述: 模拟获取表格数据
// 返回值:
//   - [][]string: 二维字符串数组表示的表格数据
// 使用示例:
//   data := mockUI.GetTableData()
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12
func (m *MockUI) GetTableData() [][]string {
	return [][]string{{"Mock", "Data"}, {"Test", "Row"}}
}