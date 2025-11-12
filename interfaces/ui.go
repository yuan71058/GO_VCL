// Package interfaces 定义应用程序中使用的接口
// 接口定义了不同模块之间的契约，实现了模块间的解耦和依赖倒置
package interfaces

// UIInterface 定义UI界面接口
// 该接口规定了UI组件必须实现的基本功能，包括状态更新、数据获取和日志记录
// 通过接口抽象，使得UI组件可以与业务逻辑模块解耦，提高代码的可测试性和可维护性
type UIInterface interface {
	// UpdateStatus 更新状态栏显示
	// 参数:
	//   - status: 要在状态栏显示的文本信息
	// 功能: 将提供的文本信息显示在应用程序的状态栏中，用于向用户反馈当前操作状态
	UpdateStatus(status string)
	
	// SetTableData 设置表格数据
	// 参数:
	//   - data: 二维字符串数组，表示要显示在表格中的数据
	// 功能: 将提供的数据填充到UI表格组件中，替换当前显示的所有数据
	// 数据格式: 每个字符串切片代表表格的一行，切片中的每个元素代表该行的单元格数据
	SetTableData(data [][]string)
	
	// GetTableData 获取表格数据
	// 返回值:
	//   - [][]string: 当前表格中显示的所有数据，以二维字符串数组形式返回
	// 功能: 从UI表格组件中提取当前显示的所有数据，用于数据处理或保存操作
	// 数据格式: 与SetTableData方法相同，每个字符串切片代表表格的一行
	GetTableData() [][]string
	
	// AddLog 添加日志信息
	// 参数:
	//   - logText: 要添加到日志区域的文本信息
	// 功能: 将提供的文本信息添加到应用程序的日志显示区域，通常带有时间戳
	// 用途: 记录操作历史、错误信息或调试信息，便于用户了解程序运行状态
	AddLog(logText string)
}