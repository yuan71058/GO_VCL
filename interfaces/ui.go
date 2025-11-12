package interfaces

// UIInterface 定义UI界面接口
type UIInterface interface {
	// 更新状态栏显示
	UpdateStatus(status string)
	
	// 设置表格数据
	SetTableData(data [][]string)
	
	// 获取表格数据
	GetTableData() [][]string
	
	// 添加日志信息
	AddLog(logText string)
}