// managers/excel_manager.go - Excel操作管理器
// 本文件提供了Excel文件的读取、写入、转换和验证等功能
// 使用excelize库作为底层Excel操作引擎，支持.xlsx和.xlsm格式
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
// 封装了Excel文件的常用操作，包括导入、导出、信息获取和格式转换等功能
// 通过UIInterface接口与用户界面交互，提供操作反馈
type ExcelManager struct {
	uiInstance interfaces.UIInterface // UI实例接口，用于显示操作状态和结果
}

// NewExcelManager 创建Excel管理器实例
// 参数:
//   - ui: UI实例，用于显示操作结果（可以为nil，之后通过SetUIInstance设置）
//
// 返回:
//   - *ExcelManager: Excel管理器实例
func NewExcelManager(ui interfaces.UIInterface) *ExcelManager {
	return &ExcelManager{
		uiInstance: ui,
	}
}

// SetUIInstance 设置UI实例
// 允许在创建ExcelManager后设置或更改UI实例
// 参数:
//   - ui: UI实例接口实现
func (em *ExcelManager) SetUIInstance(ui interfaces.UIInterface) {
	em.uiInstance = ui
}

// ImportExcel 导入Excel文件到UI表格
// 读取指定Excel文件的第一个工作表数据，并填充到UI表格中
// 自动过滤空行，确保导入的数据有效性
// 参数:
//   - filePath: Excel文件路径
//
// 返回:
//   - error: 操作错误，nil表示成功
func (em *ExcelManager) ImportExcel(filePath string) error {
	// 检查文件路径是否为空
	if filePath == "" {
		return fmt.Errorf("Excel文件路径不能为空")
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("Excel文件不存在: %s", filePath)
	}

	// 更新UI状态，通知用户开始读取文件
	em.uiInstance.UpdateStatus("正在读取Excel文件...")

	// 使用excelize库打开Excel文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Printf("关闭Excel文件时出错: %v", closeErr)
		}
	}()

	// 获取所有工作表名称
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		return fmt.Errorf("Excel文件中没有工作表")
	}

	// 获取第一个工作表的数据
	sheetName := sheetNames[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("读取工作表数据失败: %v", err)
	}

	if len(rows) == 0 {
		return fmt.Errorf("工作表中没有数据")
	}

	// 更新UI状态，通知用户正在更新表格
	em.uiInstance.UpdateStatus("正在更新UI表格...")

	// 准备表格数据，过滤空行
	tableData := make([][]string, 0)
	for _, row := range rows {
		// 检查行中是否包含非空数据
		hasData := false
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				hasData = true
				break
			}
		}
		// 只保留包含数据的行
		if hasData {
			tableData = append(tableData, row)
		}
	}

	// 更新UI表格数据
	em.uiInstance.SetTableData(tableData)

	// 更新状态栏，显示导入结果
	em.uiInstance.UpdateStatus(fmt.Sprintf("成功导入Excel文件，共%d行数据", len(tableData)))

	// 记录操作日志
	log.Printf("成功导入Excel文件: %s，工作表: %s，数据行数: %d", filePath, sheetName, len(tableData))
	return nil
}

// ExportExcel 从UI表格导出数据到Excel文件
// 将UI表格中的数据或提供的数据导出到Excel文件
// 自动创建目录结构，确保文件路径有效
// 参数:
//   - filePath: 输出Excel文件路径
//   - data: 要导出的数据（如果为nil则从UI表格获取）
//
// 返回:
//   - error: 操作错误，nil表示成功
func (em *ExcelManager) ExportExcel(filePath string, data [][]string) error {
	// 检查文件路径是否为空
	if filePath == "" {
		return fmt.Errorf("输出文件路径不能为空")
	}

	// 确保文件扩展名为.xlsx格式
	if !strings.HasSuffix(strings.ToLower(filePath), ".xlsx") {
		filePath += ".xlsx"
	}

	// 创建目录结构（如果不存在）
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 更新UI状态，通知用户开始导出
	em.uiInstance.UpdateStatus("正在导出Excel文件...")

	// 创建新的Excel文件
	f := excelize.NewFile()
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Printf("关闭Excel文件时出错: %v", closeErr)
		}
	}()

	// 获取数据源
	var tableData [][]string
	if data != nil {
		tableData = data
	} else {
		// 从UI表格获取数据
		tableData = em.uiInstance.GetTableData()
	}

	if len(tableData) == 0 {
		return fmt.Errorf("没有数据可以导出")
	}

	// 写入数据到Excel，添加详细的索引验证和调试信息
	for rowIndex := 0; rowIndex < len(tableData); rowIndex++ {
		// 检查输入数据有效性
		if tableData[rowIndex] == nil {
			log.Printf("警告: 第%d行数据为nil，跳过此行", rowIndex)
			continue
		}
		
		// 检查是否超过Excel行数限制（Excel 2007及以后版本限制为1,048,576行）
		if rowIndex >= 1048576 {
			return fmt.Errorf("行数超过Excel限制: %d", rowIndex+1)
		}
		
		rowData := tableData[rowIndex]
		log.Printf("处理第%d行数据，列数: %d", rowIndex+1, len(rowData))
		
		for colIndex := 0; colIndex < len(rowData); colIndex++ {
			// 检查是否超过Excel列数限制（Excel 2007及以后版本限制为16,384列）
			if colIndex >= 16384 {
				return fmt.Errorf("列数超过Excel限制: %d", colIndex+1)
			}
			
			cellValue := rowData[colIndex]
			
			// 确保单元格值不是nil
			if cellValue == "" {
				cellValue = "" // 确保是空字符串而不是nil
			}
			
			// 转换列索引为Excel列名，Excel中行和列索引从1开始
			excelCol := colIndex + 1
			excelRow := rowIndex + 1
			
			// 额外的安全检查，确保索引在有效范围内
			if excelCol < 1 || excelCol > 16384 || excelRow < 1 || excelRow > 1048576 {
				return fmt.Errorf("Excel索引超出范围: 列=%d, 行=%d", excelCol, excelRow)
			}
			
			// 使用defer捕获可能发生的panic，增强程序健壮性
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Excel设置单元格异常 (列=%d, 行=%d): %v", excelCol, excelRow, r)
				}
			}()
			
			// 将行列索引转换为Excel单元格地址（如A1, B2等）
			cellName, err := excelize.CoordinatesToCellName(excelCol, excelRow)
			if err != nil {
				return fmt.Errorf("转换单元格地址失败 (列=%d, 行=%d): %v", excelCol, excelRow, err)
			}
			
			// 安全地设置单元格值
			if len(cellValue) > 0 {
				f.SetCellValue("Sheet1", cellName, cellValue)
				log.Printf("设置单元格 %s: '%s'", cellName, cellValue)
			} else {
				f.SetCellValue("Sheet1", cellName, "")
				log.Printf("设置单元格 %s: (空值)", cellName)
			}
		}
	}

	// 保存文件到指定路径
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("保存Excel文件失败: %v", err)
	}

	// 更新状态栏，显示导出结果
	em.uiInstance.UpdateStatus(fmt.Sprintf("成功导出Excel文件: %s，共%d行数据", filePath, len(tableData)))

	// 记录操作日志
	log.Printf("成功导出Excel文件: %s，数据行数: %d", filePath, len(tableData))
	return nil
}

// GetExcelInfo 获取Excel文件信息
// 提取Excel文件的基本信息，包括工作表数量、数据行数、文件大小等
// 参数:
//   - filePath: Excel文件路径
//
// 返回:
//   - map[string]interface{}: Excel文件信息
//   - error: 操作错误，nil表示成功
func (em *ExcelManager) GetExcelInfo(filePath string) (map[string]interface{}, error) {
	// 检查文件路径
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
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Printf("关闭Excel文件时出错: %v", closeErr)
		}
	}()

	// 创建信息映射
	info := make(map[string]interface{})

	// 获取工作表信息
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

	// 获取文件系统信息
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		info["文件名"] = filepath.Base(filePath)
		info["文件大小"] = fmt.Sprintf("%.2f KB", float64(fileInfo.Size())/1024)
		info["修改时间"] = fileInfo.ModTime().Format("2006-01-02 15:04:05")
	}

	return info, nil
}

// ConvertToJSON 将Excel数据转换为JSON格式
// 读取Excel文件的第一个工作表数据，并将其转换为JSON格式
// 每行数据转换为JSON对象，列索引作为键名
// 参数:
//   - filePath: Excel文件路径
//
// 返回:
//   - string: JSON格式的数据
//   - error: 操作错误，nil表示成功
func (em *ExcelManager) ConvertToJSON(filePath string) (string, error) {
	// 检查文件路径
	if filePath == "" {
		return "", fmt.Errorf("Excel文件路径不能为空")
	}

	// 读取Excel数据
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return "", fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Printf("关闭Excel文件时出错: %v", closeErr)
		}
	}()

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

	// 转换为JSON格式
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
				// 将每行数据转换为键值对映射
				rowMap := make(map[string]string)
				for colIndex, cellValue := range row {
					// 使用列索引作为键（列1, 列2, ...）
					rowMap[fmt.Sprintf("列%d", colIndex+1)] = cellValue
				}
				result = append(result, rowMap)
			}
		}
	}

	// 序列化为格式化的JSON字符串
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("转换为JSON失败: %v", err)
	}

	return string(jsonData), nil
}

// ValidateExcelFile 验证Excel文件格式
// 检查文件是否存在、扩展名是否正确、是否可以正常打开
// 参数:
//   - filePath: Excel文件路径
//
// 返回:
//   - bool: 是否为有效的Excel文件
//   - string: 验证结果描述
func (em *ExcelManager) ValidateExcelFile(filePath string) (bool, string) {
	// 检查文件路径
	if filePath == "" {
		return false, "文件路径为空"
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, "文件不存在"
	}

	// 检查文件扩展名是否为支持的格式
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xlsx" && ext != ".xlsm" {
		return false, "不支持的文件格式，只支持.xlsx和.xlsm格式"
	}

	// 尝试打开文件，验证文件格式是否正确
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return false, fmt.Sprintf("文件格式错误或损坏: %v", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Printf("关闭Excel文件时出错: %v", closeErr)
		}
	}()

	// 检查是否有工作表
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		return false, "Excel文件中没有工作表"
	}

	return true, "文件格式验证通过"
}

// GetDefaultFilePath 获取默认的Excel文件路径
// 返回项目目录下的一个默认Excel文件路径，用于演示和测试
// 返回:
//   - string: 默认文件路径
func (em *ExcelManager) GetDefaultFilePath() string {
	// 返回项目目录下的data子目录中的demo.xlsx文件
	return "data" + string(os.PathSeparator) + "demo.xlsx"
}