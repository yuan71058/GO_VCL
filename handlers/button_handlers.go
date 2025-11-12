// handlers/button_handlers.go - 按钮事件处理器
// 功能描述: 处理GUI界面中各种按钮的点击事件，协调各管理器完成相应功能
// 主要功能:
//   - Excel文件操作: 导入导出Excel文件
//   - JSON数据处理: 解析、保存、加载JSON数据
//   - HTTP网络请求: 发送GET/POST请求，下载文件
//   - 数据库操作: 连接数据库，执行查询
//   - 文件操作: 文件选择、路径处理
//   - 错误处理: 统一的错误处理和用户反馈
//   - 状态更新: 实时更新UI状态信息
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"windows-gui-app/interfaces"
	"windows-gui-app/managers"
)

// ButtonHandlers 按钮事件处理器
// 负责处理GUI界面中各种按钮的点击事件，协调各管理器完成相应功能
type ButtonHandlers struct {
	ui            interfaces.UIInterface    // UI实例接口，用于显示操作结果和状态
	excelManager  *managers.ExcelManager    // Excel文件管理器，处理Excel文件操作
	jsonManager   *managers.JSONManager     // JSON数据管理器，处理JSON数据操作
	httpManager   *managers.HTTPManager     // HTTP网络管理器，处理网络请求
	databaseManager *managers.DatabaseManager // 数据库管理器，处理数据库操作
}

// NewButtonHandlers 创建按钮事件处理器实例
// 功能描述:
//   - 创建按钮事件处理器实例
//   - 初始化各管理器引用
//   - 建立UI与管理器的连接
//
// 参数:
//   - ui: UI实例，用于显示操作结果和状态
//   - excelManager: Excel文件管理器实例
//   - jsonManager: JSON数据管理器实例
//   - httpManager: HTTP网络管理器实例
//   - databaseManager: 数据库管理器实例
//
// 返回:
//   - *ButtonHandlers: 按钮事件处理器实例
//
// 使用示例:
//   handlers := handlers.NewButtonHandlers(
//       mainForm,
//       excelManager,
//       jsonManager,
//       httpManager,
//       databaseManager,
//   )
func NewButtonHandlers(ui interfaces.UIInterface, excelManager *managers.ExcelManager, jsonManager *managers.JSONManager, httpManager *managers.HTTPManager, databaseManager *managers.DatabaseManager) *ButtonHandlers {
	return &ButtonHandlers{
		ui:              ui,
		excelManager:    excelManager,
		jsonManager:     jsonManager,
		httpManager:     httpManager,
		databaseManager: databaseManager,
	}
}

// OnExcelImport 处理Excel导入按钮点击事件
// 功能描述:
//   - 打开文件选择对话框，选择Excel文件
//   - 调用Excel管理器导入文件内容
//   - 显示导入结果或错误信息
//   - 更新UI状态信息
//
// 参数:
//   - 无
//
// 返回:
//   - 无
//
// 使用示例:
//   // 绑定到按钮点击事件
//   btnExcelImport.OnClick(func(sender vcl.Object) {
//       handlers.OnExcelImport()
//   })
func (bh *ButtonHandlers) OnExcelImport() {
	log.Println("开始执行Excel导入操作")
	bh.ui.UpdateStatus("正在选择Excel文件...")

	// 显示文件选择对话框
	filePath := bh.ui.ShowOpenDialog("Excel文件|*.xlsx;*.xls")
	if filePath == "" {
		bh.ui.UpdateStatus("用户取消了文件选择")
		log.Println("用户取消了Excel文件选择")
		return
	}

	bh.ui.UpdateStatus("正在导入Excel文件...")
	log.Printf("选择的Excel文件: %s", filePath)

	// 执行导入操作
	err := bh.excelManager.ImportExcel(filePath)
	if err != nil {
		errorMsg := fmt.Sprintf("Excel导入失败: %v", err)
		bh.ui.UpdateStatus(errorMsg)
		bh.ui.ShowErrorMessage(errorMsg)
		log.Printf("Excel导入失败: %v", err)
		return
	}

	successMsg := fmt.Sprintf("Excel文件导入成功: %s", filepath.Base(filePath))
	bh.ui.UpdateStatus(successMsg)
	bh.ui.ShowSuccessMessage(successMsg)
	log.Printf("Excel文件导入成功: %s", filePath)
}

// OnExcelExport 处理Excel导出按钮点击事件
// 功能描述:
//   - 打开文件保存对话框，选择导出位置
//   - 获取当前表格数据
//   - 调用Excel管理器导出数据到文件
//   - 显示导出结果或错误信息
//
// 参数:
//   - 无
//
// 返回:
//   - 无
//
// 使用示例:
//   // 绑定到按钮点击事件
//   btnExcelExport.OnClick(func(sender vcl.Object) {
//       handlers.OnExcelExport()
//   })
func (bh *ButtonHandlers) OnExcelExport() {
	log.Println("开始执行Excel导出操作")
	bh.ui.UpdateStatus("正在选择导出位置...")

	// 显示文件保存对话框
	filePath := bh.ui.ShowSaveDialog("Excel文件|*.xlsx")
	if filePath == "" {
		bh.ui.UpdateStatus("用户取消了导出位置选择")
		log.Println("用户取消了Excel导出位置选择")
		return
	}

	// 确保文件扩展名正确
	if !strings.HasSuffix(filePath, ".xlsx") && !strings.HasSuffix(filePath, ".xls") {
		filePath += ".xlsx"
	}

	bh.ui.UpdateStatus("正在导出Excel文件...")
	log.Printf("选择的导出位置: %s", filePath)

	// 获取当前表格数据
	tableData := bh.ui.GetTableData()
	if len(tableData) == 0 {
		warningMsg := "没有可导出的数据"
		bh.ui.UpdateStatus(warningMsg)
		bh.ui.ShowWarningMessage(warningMsg)
		log.Println("Excel导出警告: 没有可导出的数据")
		return
	}

	// 执行导出操作
	err := bh.excelManager.ExportExcel(filePath, tableData)
	if err != nil {
		errorMsg := fmt.Sprintf("Excel导出失败: %v", err)
		bh.ui.UpdateStatus(errorMsg)
		bh.ui.ShowErrorMessage(errorMsg)
		log.Printf("Excel导出失败: %v", err)
		return
	}

	successMsg := fmt.Sprintf("Excel文件导出成功: %s", filepath.Base(filePath))
	bh.ui.UpdateStatus(successMsg)
	bh.ui.ShowSuccessMessage(successMsg)
	log.Printf("Excel文件导出成功: %s", filePath)
}

// OnJSONParse 处理JSON解析按钮点击事件
// 功能描述:
//   - 打开文件选择对话框，选择JSON文件
//   - 读取JSON文件内容
//   - 调用JSON管理器解析数据
//   - 显示解析结果或错误信息
//
// 参数:
//   - 无
//
// 返回:
//   - 无
//
// 使用示例:
//   // 绑定到按钮点击事件
//   btnJSONParse.OnClick(func(sender vcl.Object) {
//       handlers.OnJSONParse()
//   })
func (bh *ButtonHandlers) OnJSONParse() {
	log.Println("开始执行JSON解析操作")
	bh.ui.UpdateStatus("正在选择JSON文件...")

	// 显示文件选择对话框
	filePath := bh.ui.ShowOpenDialog("JSON文件|*.json")
	if filePath == "" {
		bh.ui.UpdateStatus("用户取消了文件选择")
		log.Println("用户取消了JSON文件选择")
		return
	}

	bh.ui.UpdateStatus("正在读取JSON文件...")
	log.Printf("选择的JSON文件: %s", filePath)

	// 读取文件内容
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		errorMsg := fmt.Sprintf("读取JSON文件失败: %v", err)
		bh.ui.UpdateStatus(errorMsg)
		bh.ui.ShowErrorMessage(errorMsg)
		log.Printf("读取JSON文件失败: %v", err)
		return
	}

	bh.ui.UpdateStatus("正在解析JSON数据...")

	// 解析JSON数据
	err = bh.jsonManager.ParseJSON(string(jsonData))
	if err != nil {
		errorMsg := fmt.Sprintf("JSON解析失败: %v", err)
		bh.ui.UpdateStatus(errorMsg)
		bh.ui.ShowErrorMessage(errorMsg)
		log.Printf("JSON解析失败: %v", err)
		return
	}

	// 获取解析结果
	jsonDataMap := bh.jsonManager.GetData()
	if jsonDataMap == nil {
		errorMsg := "JSON数据为空或格式错误"
		bh.ui.UpdateStatus(errorMsg)
		bh.ui.ShowErrorMessage(errorMsg)
		log.Println("JSON解析结果为空")
		return
	}

	// 显示解析结果
	displayData := bh.convertJSONToTableData(jsonDataMap)
	bh.ui.SetTableData(displayData)

	successMsg := fmt.Sprintf("JSON文件解析成功: %s", filepath.Base(filePath))
	bh.ui.UpdateStatus(successMsg)
	bh.ui.ShowSuccessMessage(successMsg)
	log.Printf("JSON文件解析成功: %s", filePath)
}

// OnJSONSave 处理JSON保存按钮点击事件
// 功能描述:
//   - 获取当前表格数据
//   - 转换为JSON格式
//   - 打开文件保存对话框，选择保存位置
//   - 保存JSON数据到文件
//
// 参数:
//   - 无
//
// 返回:
//   - 无
//
// 使用示例:
//   // 绑定到按钮点击事件
//   btnJSONSave.OnClick(func(sender vcl.Object) {
//       handlers.OnJSONSave()
//   })
func (bh *ButtonHandlers) OnJSONSave() {
	log.Println("开始执行JSON保存操作")
	bh.ui.UpdateStatus("正在准备JSON数据...")

	// 获取当前表格数据
	tableData := bh.ui.GetTableData()
	if len(tableData) == 0 {
		warningMsg := "没有可保存的数据"
		bh.ui.UpdateStatus(warningMsg)
		bh.ui.ShowWarningMessage(warningMsg)
		log.Println("JSON保存警告: 没有可保存的数据")
		return
	}

	// 将表格数据转换为JSON格式
	jsonData, err := bh.convertTableDataToJSON(tableData)
	if err != nil {
		errorMsg := fmt.Sprintf("数据转换失败: %v", err)
		bh.ui.UpdateStatus(errorMsg)
		bh.ui.ShowErrorMessage(errorMsg)
		log.Printf("表格数据转换为JSON失败: %v", err)
		return
	}

	bh.ui.UpdateStatus("正在选择保存位置...")

	// 显示文件保存对话框
	filePath := bh.ui.ShowSaveDialog("JSON文件|*.json")
	if filePath == "" {
		bh.ui.UpdateStatus("用户取消了保存位置选择")
		log.Println("用户取消了JSON保存位置选择")
		return
	}

	// 确保文件扩展名正确
	if !strings.HasSuffix(filePath, ".json") {
		filePath += ".json"
	}

	bh.ui.UpdateStatus("正在保存JSON文件...")
	log.Printf("选择的保存位置: %s", filePath)

	// 保存JSON数据到文件
	err = os.WriteFile(filePath, []byte(jsonData), 0644)
	if err != nil {
		errorMsg := fmt.Sprintf("JSON文件保存失败: %v", err)
		bh.ui.UpdateStatus(errorMsg)
		bh.ui.ShowErrorMessage(errorMsg)
		log.Printf("JSON文件保存失败: %v", err)
		return
	}

	// 更新JSON管理器中的数据
	bh.jsonManager.ParseJSON(jsonData)

	successMsg := fmt.Sprintf("JSON文件保存成功: %s", filepath.Base(filePath))
	bh.ui.UpdateStatus(successMsg)
	bh.ui.ShowSuccessMessage(successMsg)
	log.Printf("JSON文件保存成功: %s", filePath)
}

// OnJSONLoad 处理JSON加载按钮点击事件
// 功能描述:
//   - 打开文件选择对话框，选择JSON文件
//   - 读取并解析JSON文件
//   - 将数据转换为表格格式并显示
//
// 参数:
//   - 无
//
// 返回:
//   - 无
//
// 使用示例:
//   // 绑定到按钮点击事件
//   btnJSONLoad.OnClick(func(sender vcl.Object) {
//       handlers.OnJSONLoad()
//   })
func (bh *ButtonHandlers) OnJSONLoad() {
	log.Println("开始执行JSON加载操作")
	bh.ui.UpdateStatus("正在选择JSON文件...")

	// 显示文件选择对话框
	filePath := bh.ui.ShowOpenDialog("JSON文件|*.json")
	if filePath == "" {
		bh.ui.UpdateStatus("用户取消了文件选择")
		log.Println("用户取消了JSON文件选择")
		return
	}

	bh.ui.UpdateStatus("正在加载JSON文件...")
	log.Printf("选择的JSON文件: %s", filePath)

	// 读取文件内容
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		errorMsg := fmt.Sprintf("读取JSON文件失败: %v", err)
		bh.ui.UpdateStatus(errorMsg)
		bh.ui.ShowErrorMessage(errorMsg)
		log.Printf("读取JSON文件失败: %v", err)
		return
	}

	bh.ui.UpdateStatus("正在解析JSON数据...")

	// 解析JSON数据
	err = bh.jsonManager.ParseJSON(string(jsonData))
	if err != nil {
		errorMsg := fmt.Sprintf("JSON解析失败: %v", err)
		bh.ui.UpdateStatus(errorMsg)
		bh.ui.ShowErrorMessage(errorMsg)
		log.Printf("JSON解析失败: %v", err)
		return
	}

	// 获取解析结果并显示
	jsonDataMap := bh.jsonManager.GetData()
	if jsonDataMap != nil {
		displayData := bh.convertJSONToTableData(jsonDataMap)
		bh.ui.SetTableData(displayData)
	}

	successMsg := fmt.Sprintf("JSON文件加载成功: %s", filepath.Base(filePath))
	bh.ui.UpdateStatus(successMsg)
	bh.ui.ShowSuccessMessage(successMsg)
	log.Printf("JSON文件加载成功: %s", filePath)
}

// OnHTTPRequest 处理HTTP请求按钮点击事件
// 功能描述:
//   - 获取用户输入的URL
//   - 解析URL并验证格式
//   - 发送HTTP GET或POST请求
//   - 显示请求结果或错误信息
//
// 参数:
//   - 无
//
// 返回:
//   - 无
//
// 使用示例:
//   // 绑定到按钮点击事件
//   btnHTTPRequest.OnClick(func(sender vcl.Object) {
//       handlers.OnHTTPRequest()
//   })
func (bh *ButtonHandlers) OnHTTPRequest() {
	log.Println("开始执行HTTP请求操作")
	bh.ui.UpdateStatus("正在准备HTTP请求...")

	// 获取用户输入的URL
	requestURL := bh.ui.GetInputText()
	if requestURL == "" {
		warningMsg := "请输入请求URL"
		bh.ui.UpdateStatus(warningMsg)
		bh.ui.ShowWarningMessage(warningMsg)
		log.Println("HTTP请求警告: 未输入URL")
		return
	}

	// 验证URL格式
	parsedURL, err := url.Parse(requestURL)
	if err != nil || parsedURL.Scheme == "" {
		errorMsg := "请输入有效的URL（如：https://api.example.com）"
		bh.ui.UpdateStatus(errorMsg)
		bh.ui.ShowErrorMessage(errorMsg)
		log.Printf("HTTP请求URL格式错误: %s", requestURL)
		return
	}

	bh.ui.UpdateStatus("正在发送HTTP请求...")
	log.Printf("发送HTTP请求到: %s", requestURL)

	// 发送HTTP请求（根据URL判断使用GET还是POST）
	var response string
	var requestErr error

	if strings.Contains(strings.ToUpper(requestURL), "POST") || strings.Contains(requestURL, "api/") {
		// 对于包含API的URL，使用POST请求
		response, requestErr = bh.httpManager.POSTRequest(requestURL, "application/json", "{}")
	} else {
		// 默认使用GET请求
		response, requestErr = bh.httpManager.GETRequest(requestURL)
	}

	if requestErr != nil {
		errorMsg := fmt.Sprintf("HTTP请求失败: %v", requestErr)
		bh.ui.UpdateStatus(errorMsg)
		bh.ui.ShowErrorMessage(errorMsg)
		log.Printf("HTTP请求失败: %v", requestErr)
		return
	}

	// 显示响应结果
	responseData := [][]string{
		{"响应内容"},
		{response},
	}
	bh.ui.SetTableData(responseData)

	successMsg := fmt.Sprintf("HTTP请求成功: %s", requestURL)
	bh.ui.UpdateStatus(successMsg)
	bh.ui.ShowSuccessMessage(successMsg)
	log.Printf("HTTP请求成功: %s，响应长度: %d字符", requestURL, len(response))
}

// OnFileDownload 处理文件下载按钮点击事件
// 功能描述:
//   - 获取用户输入的下载URL
//   - 打开文件保存对话框，选择保存位置
//   - 下载文件并保存到指定位置
//   - 显示下载结果或错误信息
//
// 参数:
//   - 无
//
// 返回:
//   - 无
//
// 使用示例:
//   // 绑定到按钮点击事件
//   btnFileDownload.OnClick(func(sender vcl.Object) {
//       handlers.OnFileDownload()
//   })
func (bh *ButtonHandlers) OnFileDownload() {
	log.Println("开始执行文件下载操作")
	bh.ui.UpdateStatus("正在准备文件下载...")

	// 获取用户输入的下载URL
	downloadURL := bh.ui.GetInputText()
	if downloadURL == "" {
		warningMsg := "请输入下载URL"
		bh.ui.UpdateStatus(warningMsg)
		bh.ui.ShowWarningMessage(warningMsg)
		log.Println("文件下载警告: 未输入URL")
		return
	}

	// 验证URL格式
	parsedURL, err := url.Parse(downloadURL)
	if err != nil || parsedURL.Scheme == "" {
		errorMsg := "请输入有效的下载URL"
		bh.ui.UpdateStatus(errorMsg)
		bh.ui.ShowErrorMessage(errorMsg)
		log.Printf("文件下载URL格式错误: %s", downloadURL)
		return
	}

	bh.ui.UpdateStatus("正在选择保存位置...")

	// 从URL中提取文件名
	fileName := filepath.Base(parsedURL.Path)
	if fileName == "" || fileName == "/" {
		fileName = "downloaded_file"
	}

	// 显示文件保存对话框，默认文件名从URL提取
	filePath := bh.ui.ShowSaveDialogWithName("所有文件|*.*", fileName)
	if filePath == "" {
		bh.ui.UpdateStatus("用户取消了保存位置选择")
		log.Println("用户取消了文件下载保存位置选择")
		return
	}

	bh.ui.UpdateStatus("正在下载文件...")
	log.Printf("开始下载文件: %s -> %s", downloadURL, filePath)

	// 下载文件
	err = bh.httpManager.DownloadFile(downloadURL, filePath)
	if err != nil {
		errorMsg := fmt.Sprintf("文件下载失败: %v", err)
		bh.ui.UpdateStatus(errorMsg)
		bh.ui.ShowErrorMessage(errorMsg)
		log.Printf("文件下载失败: %v", err)
		return
	}

	// 显示下载结果信息
	fileInfo := [][]string{
		{"属性", "值"},\t	{"文件名", filepath.Base(filePath)},
		{"保存路径", filePath},
		{"下载时间", time.Now().Format("2006-01-02 15:04:05")},
		{"状态", "下载成功"},
	}
	bh.ui.SetTableData(fileInfo)

	successMsg := fmt.Sprintf("文件下载成功: %s", filepath.Base(filePath))
	bh.ui.UpdateStatus(successMsg)
	bh.ui.ShowSuccessMessage(successMsg)
	log.Printf("文件下载成功: %s", filePath)
}

// 辅助函数

// convertJSONToTableData 将JSON数据转换为表格数据格式
// 参数:
//   - jsonData: JSON数据映射
//
// 返回:
//   - [][]string: 表格数据格式，第一行为列名，后续行为数据
func (bh *ButtonHandlers) convertJSONToTableData(jsonData map[string]interface{}) [][]string {
	if jsonData == nil || len(jsonData) == 0 {
		return [][]string{
			{"信息"},
			{"无数据"},
		}
	}

	var tableData [][]string

	// 如果是数组类型，处理为表格格式
	if dataArray, ok := jsonData["data"].([]interface{}); ok && len(dataArray) > 0 {
		// 获取列名（从第一个元素）
		if firstItem, ok := dataArray[0].(map[string]interface{}); ok {
			var headers []string
			for key := range firstItem {
				headers = append(headers, key)
			}
			tableData = append(tableData, headers)

			// 添加数据行
			for _, item := range dataArray {
				if itemMap, ok := item.(map[string]interface{}); ok {
					var row []string
					for _, header := range headers {
						value := ""
						if val, exists := itemMap[header]; exists && val != nil {
							value = fmt.Sprintf("%v", val)
						}
						row = append(row, value)
					}
					tableData = append(tableData, row)
				}
			}
		}
	} else {
		// 处理为键值对格式
		tableData = append(tableData, []string{"键", "值"})
		for key, value := range jsonData {
			valueStr := ""
			if value != nil {
				valueStr = fmt.Sprintf("%v", value)
			}
			tableData = append(tableData, []string{key, valueStr})
		}
	}

	return tableData
}

// convertTableDataToJSON 将表格数据转换为JSON格式
// 参数:
//   - tableData: 表格数据格式，第一行为列名，后续行为数据
//
// 返回:
//   - string: JSON格式字符串
//   - error: 转换错误
func (bh *ButtonHandlers) convertTableDataToJSON(tableData [][]string) (string, error) {
	if len(tableData) == 0 {
		return "{}", nil
	}

	// 如果只有一行，处理为键值对格式
	if len(tableData) == 1 {
		result := make(map[string]string)
		for i, cell := range tableData[0] {
			result[fmt.Sprintf("column_%d", i)] = cell
		}
		jsonBytes, err := json.Marshal(result)
		return string(jsonBytes), err
	}

	// 处理为数组格式
	var result []map[string]string
	headers := tableData[0]

	for i := 1; i < len(tableData); i++ {
		row := make(map[string]string)
		for j, cell := range tableData[i] {
			if j < len(headers) {
				row[headers[j]] = cell
			}
		}
		result = append(result, row)
	}

	// 包装为data字段
	finalResult := map[string]interface{}{
		"data": result,
		"count": len(result),
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	jsonBytes, err := json.MarshalIndent(finalResult, "", "  ")
	return string(jsonBytes), err
}