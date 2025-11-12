// managers/excel_manager.go - Excel操作管理器
// 功能描述: 提供Excel文件的导入导出功能，支持.xlsx格式
// 主要功能:
//   - Excel文件导入：读取Excel文件内容并显示在UI表格中
//   - Excel文件导出：将UI表格数据导出为Excel文件
//   - 数据验证：检查文件路径、数据格式和Excel限制
//   - 错误处理：提供详细的错误信息和状态反馈
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12

package managers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
	"windows-gui-app/interfaces"
)

// ExcelManager Excel操作管理器
// 负责处理Excel文件的读写操作，提供UI集成功能
type ExcelManager struct {
	uiInstance interfaces.UIInterface // UI实例接口，用于显示操作状态和结果
}

// NewExcelManager 创建Excel管理器实例
// 参数:
//   - ui: UI实例，用于显示操作结果（可以为nil，之后通过SetUIInstance设置）
//
// 返回:
//   - *ExcelManager: Excel管理器实例
//
// 使用示例:
//   excelManager := managers.NewExcelManager(mainForm)
func NewExcelManager(ui interfaces.UIInterface) *ExcelManager {
	return &ExcelManager{
		uiInstance: ui,
	}
}

// SetUIInstance 设置UI实例
// 用于在运行时更新UI实例引用
// 参数:
//   - ui: 新的UI实例
func (em *ExcelManager) SetUIInstance(ui interfaces.UIInterface) {
	em.uiInstance = ui
}

// ImportExcel 导入Excel文件到UI表格
// 功能描述:
//   - 读取指定路径的Excel文件（.xlsx格式）
//   - 解析第一个工作表的数据内容
//   - 将数据转换为表格格式并显示在UI中
//   - 提供详细的导入状态和错误信息
//
// 参数:
//   - filePath: Excel文件路径（必须存在且为.xlsx格式）
//
// 返回:
//   - error: 操作错误，nil表示成功
//
// 可能的错误:
//   - "Excel文件路径不能为空": filePath参数为空
//   - "Excel文件不存在": 指定路径的文件不存在
//   - "打开Excel文件失败": 文件格式错误或损坏
//   - "Excel文件中没有工作表": 文件为空或没有工作表
//   - "工作表中没有数据": 第一个工作表为空
//
// 使用示例:
//   err := excelManager.ImportExcel("C:\\data\\employees.xlsx")
//   if err != nil {
//       log.Printf("导入失败: %v", err)
//   }
func (em *ExcelManager) ImportExcel(filePath string) error {
	// 验证输入参数
	if filePath == "" {
		return fmt.Errorf("Excel文件路径不能为空")
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("Excel文件不存在: %s", filePath)
	}

	// 更新UI状态，通知用户开始导入操作
	if em.uiInstance != nil {
		em.uiInstance.UpdateStatus("正在读取Excel文件...")
	}

	// 打开Excel文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	// 获取工作表列表
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		return fmt.Errorf("Excel文件中没有工作表")
	}

	// 默认读取第一个工作表
	sheetName := sheetNames[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("读取工作表数据失败: %v", err)
	}

	if len(rows) == 0 {
		return fmt.Errorf("工作表中没有数据")
	}

	// 更新UI状态
	if em.uiInstance != nil {
		em.uiInstance.UpdateStatus("正在更新UI表格...")
	}

	// 数据预处理：过滤空行
	tableData := make([][]string, 0)
	for rowIndex, row := range rows {
		// 检查是否为有效数据行（非空行）
		hasData := false
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				hasData = true
				break
			}
		}
		// 只保留包含有效数据的行
		if hasData {
			tableData = append(tableData, row)
		} else {
			log.Printf("跳过空行: 第%d行", rowIndex+1)
		}
	}

	// 更新UI表格数据
	if em.uiInstance != nil {
		em.uiInstance.SetTableData(tableData)
		// 更新状态信息
		statusMsg := fmt.Sprintf("成功导入Excel文件，共%d行数据", len(tableData))
		em.uiInstance.UpdateStatus(statusMsg)
	}

	// 记录操作日志
	log.Printf("成功导入Excel文件: %s，工作表: %s，数据行数: %d", filePath, sheetName, len(tableData))
	return nil
}

// ExportExcel 从UI表格导出数据到Excel文件
// 功能描述:
//   - 获取UI表格中的数据（或指定的数据）
//   - 创建新的Excel文件并写入数据
//   - 支持自动创建目录和文件扩展名处理
//   - 提供Excel格式限制检查和错误处理
//
// 参数:
//   - filePath: 输出Excel文件路径（自动添加.xlsx扩展名）
//   - data: 要导出的数据（如果为nil则从UI表格获取）
//
// 返回:
//   - error: 操作错误，nil表示成功
//
// 特殊功能:
//   - 自动添加.xlsx扩展名（如果没有）
//   - 自动创建输出目录（如果不存在）
//   - 支持1048576行×16384列（Excel限制）
//   - 空值自动转换为空字符串
//
// 可能的错误:
//   - "输出文件路径不能为空": filePath参数为空
//   - "创建目录失败": 输出目录创建失败
//   - "没有数据可以导出": 数据源为空
//   - "行数超过Excel限制": 数据行数超过1,048,576
//   - "列数超过Excel限制": 数据列数超过16,384
//   - "转换单元格地址失败": Excel内部地址转换错误
//
// 使用示例:
//   // 从UI表格导出
//   err := excelManager.ExportExcel("C:\\output\\report.xlsx", nil)
//   
//   // 导出指定数据
//   data := [][]string{
//       {"姓名", "年龄", "部门"},
//       {"张三", "25", "技术部"},
//       {"李四", "30", "销售部"},
//   }
//   err := excelManager.ExportExcel("C:\\output\\employees.xlsx", data)
func (em *ExcelManager) ExportExcel(filePath string, data [][]string) error {
	// 验证文件路径
	if filePath == "" {
		return fmt.Errorf("输出文件路径不能为空")
	}

	// 确保文件扩展名为.xlsx
	if !strings.HasSuffix(strings.ToLower(filePath), ".xlsx") {
		filePath += ".xlsx"
		log.Printf("自动添加.xlsx扩展名，新路径: %s", filePath)
	}

	// 创建输出目录（如果不存在）
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 更新UI状态
	if em.uiInstance != nil {
		em.uiInstance.UpdateStatus("正在导出Excel文件...")
	}

	// 创建新的Excel文件
	f := excelize.NewFile()
	defer f.Close()

	// 获取数据源
	var tableData [][]string
	if data != nil {
		// 使用指定的数据
		tableData = data
	} else if em.uiInstance != nil {
		// 从UI表格获取数据
		tableData = em.uiInstance.GetTableData()
	}

	// 验证数据有效性
	if len(tableData) == 0 {
		return fmt.Errorf("没有数据可以导出")
	}

	// 写入数据到Excel，添加详细的索引验证和调试信息
	for rowIndex := 0; rowIndex < len(tableData); rowIndex++ {
		// 数据行验证
		if tableData[rowIndex] == nil {
			log.Printf("警告: 第%d行数据为nil，跳过此行", rowIndex+1)
			continue
		}
		
		// Excel行数限制检查（Excel 2007+ 最大支持1,048,576行）
		if rowIndex >= 1048576 {
			return fmt.Errorf("行数超过Excel限制: %d", rowIndex+1)
		}
		
		rowData := tableData[rowIndex]
		log.Printf("处理第%d行数据，列数: %d", rowIndex+1, len(rowData))
		
		// 处理每一列数据
		for colIndex := 0; colIndex < len(rowData); colIndex++ {
			// Excel列数限制检查（Excel 2007+ 最大支持16,384列）
			if colIndex >= 16384 {
				return fmt.Errorf("列数超过Excel限制: %d", colIndex+1)
			}
			
			cellValue := rowData[colIndex]
			
			// 确保单元格值有效（避免nil值）
			if cellValue == "" {
				cellValue = "" // 确保是空字符串而不是nil
			}
			
			// 转换索引为Excel单元格地址（Excel使用1-based索引）
			excelCol := colIndex + 1
			excelRow := rowIndex + 1
			
			// 地址范围安全检查
			if excelCol < 1 || excelCol > 16384 || excelRow < 1 || excelRow > 1048576 {
				return fmt.Errorf("Excel索引超出范围: 列=%d, 行=%d", excelCol, excelRow)
			}
			
			// 使用defer捕获和处理可能的panic
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Excel设置单元格异常 (列=%d, 行=%d): %v", excelCol, excelRow, r)
				}
			}()
			
			// 生成Excel单元格地址（如"A1", "B2"等）
			cellName, err := excelize.CoordinatesToCellName(excelCol, excelRow)
			if err != nil {
				return fmt.Errorf("转换单元格地址失败 (列=%d, 行=%d): %v", excelCol, excelRow, err)
			}
			
			// 设置单元格值，包含错误处理
			if err := f.SetCellValue("Sheet1", cellName, cellValue); err != nil {
				log.Printf("设置单元格值失败 (单元格=%s, 值=%s): %v", cellName, cellValue, err)
				// 继续处理其他单元格，不中断整个导出过程
			}
		}
	}

	// 保存Excel文件
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("保存Excel文件失败: %v", err)
	}

	// 更新UI状态和完成信息
	if em.uiInstance != nil {
		statusMsg := fmt.Sprintf("成功导出Excel文件: %s", filePath)
		em.uiInstance.UpdateStatus(statusMsg)
	}

	// 记录操作成功日志
	log.Printf("成功导出Excel文件: %s，数据行数: %d", filePath, len(tableData))
	return nil
}

// GetExcelInfo 获取Excel文件信息
// 参数:
//   - filePath: Excel文件路径
//
// 返回:
//   - map[string]interface{}: Excel文件信息
//   - error: 操作错误，nil表示成功
func (em *ExcelManager) GetExcelInfo(filePath string) (map[string]interface{}, error) {
	if filePath == "" {
		return nil, fmt.Errorf("Excel文件路径不能为空")
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Excel文件不存在: %s", filePath)
	}

	// 打开Excel文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	// 获取文件信息
	info := make(map[string]interface{})

	// 工作表信息
	sheetNames := f.GetSheetList()
	info["工作表数量"] = len(sheetNames)
	info["工作表列表"] = sheetNames

	// 获取第一个工作表的详细信息
	if len(sheetNames) > 0 {
		sheetName := sheetNames[0]
		rows, err := f.GetRows(sheetName)
		if err == nil {
			info["数据行数"] = len(rows)
			if len(rows) > 0 {
				info["列数"] = len(rows[0])
			}
		}
	}

	// 文件信息
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		info["文件名"] = filepath.Base(filePath)
		info["文件大小"] = fmt.Sprintf("%.2f KB", float64(fileInfo.Size())/1024)
		info["修改时间"] = fileInfo.ModTime().Format("2006-01-02 15:04:05")
	}

	return info, nil
}

// ConvertToJSON 将Excel数据转换为JSON格式
// 参数:
//   - filePath: Excel文件路径
//
// 返回:
//   - string: JSON格式的数据
//   - error: 操作错误，nil表示成功
func (em *ExcelManager) ConvertToJSON(filePath string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("Excel文件路径不能为空")
	}

	// 读取Excel数据
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return "", fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	// 获取第一个工作表的数据
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		return "", fmt.Errorf("Excel文件中没有工作表")
	}

	rows, err := f.GetRows(sheetNames[0])
	if err != nil {
		return "", fmt.Errorf("读取工作表数据失败: %v", err)
	}

	if len(rows) == 0 {
		return "", fmt.Errorf("工作表中没有数据")
	}

	// 转换为JSON
	var result []map[string]string
	for _, row := range rows {
		if len(row) > 0 {
			// 过滤完全空的行
			hasData := false
			for _, cell := range row {
				if strings.TrimSpace(cell) != "" {
					hasData = true
					break
				}
			}
			if hasData {
				rowMap := make(map[string]string)
				for colIndex, cellValue := range row {
					// 使用列索引作为键
					rowMap[fmt.Sprintf("列%d", colIndex+1)] = cellValue
				}
				result = append(result, rowMap)
			}
		}
	}

	// 序列化为JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("转换为JSON失败: %v", err)
	}

	return string(jsonData), nil
}

// ValidateExcelFile 验证Excel文件格式
// 参数:
//   - filePath: Excel文件路径
//
// 返回:
//   - bool: 是否为有效的Excel文件
//   - string: 验证结果描述
func (em *ExcelManager) ValidateExcelFile(filePath string) (bool, string) {
	if filePath == "" {
		return false, "文件路径为空"
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, "文件不存在"
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xlsx" && ext != ".xlsm" {
		return false, "不支持的文件格式，只支持.xlsx和.xlsm格式"
	}

	// 尝试打开文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return false, fmt.Sprintf("文件格式错误或损坏: %v", err)
	}
	defer f.Close()

	// 检查是否有工作表
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		return false, "Excel文件中没有工作表"
	}

	return true, "文件格式验证通过"
}

// GetDefaultFilePath 获取默认的Excel文件路径
// 返回:
//   - string: 默认文件路径
func (em *ExcelManager) GetDefaultFilePath() string {
	// 返回项目目录下的一个默认Excel文件路径
	return "data" + string(os.PathSeparator) + "demo.xlsx"
}
