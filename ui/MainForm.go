// ui/MainForm.go - 简化版主窗口表单定义
package ui

import (
	"fmt"
	"log"
	"time"

	"windows-gui-app/interfaces"
	"windows-gui-app/managers"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

// MainForm 主窗口结构体
type MainForm struct {
	*vcl.TForm

	// UI组件
	PanelMain    *vcl.TPanel // 主面板
	PanelButtons *vcl.TPanel // 按钮面板
	PanelStatus  *vcl.TPanel // 状态面板
	PanelTable   *vcl.TPanel // 表格面板

	// 按钮组件
	BtnExcel    *vcl.TButton // Excel操作按钮
	BtnJSON     *vcl.TButton // JSON操作按钮
	BtnHTTP     *vcl.TButton // HTTP操作按钮
	BtnDatabase *vcl.TButton // 数据库操作按钮
	BtnClose    *vcl.TButton // 关闭按钮

	// 状态和显示组件
	StatusBar *vcl.TStatusBar  // 状态栏
	TableData *vcl.TStringGrid // 数据表格
	EditInput *vcl.TMemo       // 输入编辑框
	LabelInfo *vcl.TLabel      // 信息标签

	// 管理器
	ExcelManager *managers.ExcelManager
	JSONManager  *managers.JSONManager
	HTTPManager  *managers.HTTPManager
	DBManager    *managers.DatabaseManager

	// 时间相关
	lastUpdate time.Time // 最后更新时间
}

// NewMainForm 创建新的主窗口实例
func NewMainForm() *MainForm {
	form := &MainForm{
		lastUpdate: time.Now(),
	}

	// 创建窗口
	form.TForm = vcl.Application.CreateForm()
	form.TForm.SetCaption("GO VCL 多功能演示程序")
	form.TForm.SetWidth(1024)
	form.TForm.SetHeight(768)
	form.TForm.SetPosition(types.PoScreenCenter)
	form.TForm.SetOnClose(form.onClose)
	form.TForm.SetOnShow(form.onShow)

	// 初始化管理器（注意：这里需要小心循环引用）
	// 先创建管理器并设置UI实例
	excelManager := managers.NewExcelManager(nil)
	excelManager.SetUIInstance(form)

	jsonManager := managers.NewJSONManager(nil)
	jsonManager.SetUIInstance(form)

	httpManager := managers.NewHTTPManager(nil)
	httpManager.SetUIInstance(form)

	dbManager := managers.NewDatabaseManager(nil)
	dbManager.SetUIInstance(form)

	form.ExcelManager = excelManager
	form.JSONManager = jsonManager
	form.HTTPManager = httpManager
	form.DBManager = dbManager

	return form
}

// Show 显示窗口
func (f *MainForm) Show() {
	f.createUI()
	f.TForm.Show()
}

// createUI 创建用户界面
func (f *MainForm) createUI() {
	// 创建主面板
	f.PanelMain = vcl.NewPanel(f.TForm)
	f.PanelMain.SetParent(f.TForm)
	f.PanelMain.SetAlign(types.AlClient)
	f.PanelMain.SetBevelOuter(types.BvNone)

	// 创建按钮面板
	f.PanelButtons = vcl.NewPanel(f.TForm)
	f.PanelButtons.SetParent(f.PanelMain)
	f.PanelButtons.SetAlign(types.AlTop)
	f.PanelButtons.SetHeight(80)
	f.PanelButtons.SetBevelOuter(types.BvNone)

	// 创建按钮
	f.createButtons()

	// 创建状态面板
	f.PanelStatus = vcl.NewPanel(f.TForm)
	f.PanelStatus.SetParent(f.PanelMain)
	f.PanelStatus.SetAlign(types.AlBottom)
	f.PanelStatus.SetHeight(30)
	f.PanelStatus.SetBevelOuter(types.BvNone)

	// 创建状态栏
	f.StatusBar = vcl.NewStatusBar(f.TForm)
	f.StatusBar.SetParent(f.PanelStatus)
	f.StatusBar.SetAlign(types.AlClient)
	f.StatusBar.Panels().Add()
	f.StatusBar.Panels().Items(0).SetText("就绪")

	// 创建表格面板
	f.PanelTable = vcl.NewPanel(f.TForm)
	f.PanelTable.SetParent(f.PanelMain)
	f.PanelTable.SetAlign(types.AlClient)
	f.PanelTable.SetBevelOuter(types.BvNone)

	// 创建数据表格（修复索引越界问题）
	f.TableData = vcl.NewStringGrid(f.TForm)
	f.TableData.SetParent(f.PanelTable)
	f.TableData.SetAlign(types.AlLeft)
	f.TableData.SetWidth(400)
	gridOptions := f.TableData.Options()
	gridOptions = gridOptions | types.TGridOptions(types.GoRowSelect)
	f.TableData.SetOptions(gridOptions)
	f.TableData.SetRowCount(5)
	f.TableData.SetColCount(3)
	
	// 统一使用列, 行的参数顺序添加表格数据
	// 确保每次调用都在有效范围内
	defer func() {
		if r := recover(); r != nil {
			log.Printf("表格初始化时发生panic: %v", r)
		}
	}()
	
	// 使用安全的索引设置表格数据
	tableData := [][]string{
		{"功能", "状态", "描述"},
		{"Excel操作", "待测试", "导入/导出Excel文件"},
		{"JSON操作", "待测试", "解析/生成JSON数据"},
		{"HTTP操作", "待测试", "网络请求和下载"},
		{"数据库操作", "待测试", "SQLite数据库管理"},
	}
	
	// 安全地设置所有单元格
	for r, row := range tableData {
		if r < 0 || r >= 5 {
			log.Printf("跳过无效行索引: %d", r)
			continue
		}
		for c, cell := range row {
			if c < 0 || c >= 3 {
				log.Printf("跳过无效列索引: %d", c)
				continue
			}
			
			// 使用列, 行的参数顺序调用SetCells
			colParam := int32(c)  // 列在前
			rowParam := int32(r)  // 行在后
			
			log.Printf("设置表格单元格: (列=%d, 行=%d) = '%s'", colParam, rowParam, cell)
			f.TableData.SetCells(colParam, rowParam, cell)
		}
	}

	// 创建输入编辑框
	f.EditInput = vcl.NewMemo(f.TForm)
	f.EditInput.SetParent(f.PanelTable)
	f.EditInput.SetAlign(types.AlClient)
	f.EditInput.SetScrollBars(types.SsBoth)
	f.EditInput.SetReadOnly(true)
	f.EditInput.SetText("欢迎使用 GO VCL 多功能演示程序！\n\n本程序提供以下功能：\n\n1. Excel操作 - 导入/导出Excel文件\n2. JSON操作 - 解析/生成JSON数据\n3. HTTP操作 - 网络请求和文件下载\n4. 数据库操作 - SQLite数据库管理\n\n请点击上方按钮开始使用。")

	// 创建信息标签
	f.LabelInfo = vcl.NewLabel(f.TForm)
	f.LabelInfo.SetParent(f.PanelButtons)
	f.LabelInfo.SetCaption("GO VCL 多功能演示程序 v1.0")
	f.LabelInfo.Font().SetSize(14)
	fontStyle := f.LabelInfo.Font().Style()
	fontStyle = fontStyle | types.TFontStyles(types.FsBold)
	f.LabelInfo.Font().SetStyle(fontStyle)
	f.LabelInfo.SetLeft(20)
	f.LabelInfo.SetTop(10)

	// 设置事件处理
	f.setupEvents()
}

// createButtons 创建按钮组件
func (f *MainForm) createButtons() {
	// Excel操作按钮
	f.BtnExcel = vcl.NewButton(f.TForm)
	f.BtnExcel.SetParent(f.PanelButtons)
	f.BtnExcel.SetCaption("Excel操作")
	f.BtnExcel.SetWidth(120)
	f.BtnExcel.SetHeight(40)
	f.BtnExcel.SetLeft(200)
	f.BtnExcel.SetTop(25)
	f.BtnExcel.SetOnClick(f.onExcelClick)

	// JSON操作按钮
	f.BtnJSON = vcl.NewButton(f.TForm)
	f.BtnJSON.SetParent(f.PanelButtons)
	f.BtnJSON.SetCaption("JSON操作")
	f.BtnJSON.SetWidth(120)
	f.BtnJSON.SetHeight(40)
	f.BtnJSON.SetLeft(330)
	f.BtnJSON.SetTop(25)
	f.BtnJSON.SetOnClick(f.onJSONClick)

	// HTTP操作按钮
	f.BtnHTTP = vcl.NewButton(f.TForm)
	f.BtnHTTP.SetParent(f.PanelButtons)
	f.BtnHTTP.SetCaption("HTTP操作")
	f.BtnHTTP.SetWidth(120)
	f.BtnHTTP.SetHeight(40)
	f.BtnHTTP.SetLeft(460)
	f.BtnHTTP.SetTop(25)
	f.BtnHTTP.SetOnClick(f.onHTTPClick)

	// 数据库操作按钮
	f.BtnDatabase = vcl.NewButton(f.TForm)
	f.BtnDatabase.SetParent(f.PanelButtons)
	f.BtnDatabase.SetCaption("数据库操作")
	f.BtnDatabase.SetWidth(120)
	f.BtnDatabase.SetHeight(40)
	f.BtnDatabase.SetLeft(590)
	f.BtnDatabase.SetTop(25)
	f.BtnDatabase.SetOnClick(f.onDatabaseClick)

	// 关闭按钮
	f.BtnClose = vcl.NewButton(f.TForm)
	f.BtnClose.SetParent(f.PanelButtons)
	f.BtnClose.SetCaption("关闭")
	f.BtnClose.SetWidth(80)
	f.BtnClose.SetHeight(40)
	f.BtnClose.SetLeft(750)
	f.BtnClose.SetTop(25)
	f.BtnClose.SetOnClick(f.onCloseClick)
}

// setupEvents 设置事件处理
func (f *MainForm) setupEvents() {
	// 这里可以设置各种事件处理
}

// 事件处理函数

// onShow 窗口显示事件
func (f *MainForm) onShow(sender vcl.IObject) {
	f.UpdateStatus("程序已启动就绪")
	log.Println("主窗口已显示")
}

// onClose 窗口关闭事件
func (f *MainForm) onClose(sender vcl.IObject, action *types.TCloseAction) {
	log.Println("程序正在退出...")
	f.UpdateStatus("程序退出")
	*action = types.CaFree // 释放窗口
}

// onCloseClick 关闭按钮点击事件
func (f *MainForm) onCloseClick(sender vcl.IObject) {
	f.TForm.Close()
}

// onExcelClick Excel操作按钮事件
func (f *MainForm) onExcelClick(sender vcl.IObject) {
	if f.ExcelManager == nil {
		f.AddLog("Excel管理器未初始化")
		return
	}

	// 测试Excel功能
	f.AddLog("=== Excel操作测试开始 ===")

	// 1. 测试获取默认Excel文件路径
	defaultPath := f.ExcelManager.GetDefaultFilePath()
	f.AddLog(fmt.Sprintf("默认Excel文件路径: %s", defaultPath))

	// 2. 测试创建示例Excel文件
	f.AddLog("创建示例Excel文件...")
	excelData := [][]string{
		{"姓名", "年龄", "城市", "职业"},
		{"张三", "28", "北京", "软件工程师"},
		{"李四", "32", "上海", "产品经理"},
		{"王五", "25", "广州", "设计师"},
		{"赵六", "30", "深圳", "数据分析师"},
	}

	// 创建Excel文件
	if err := f.ExcelManager.ExportExcel("test_data.xlsx", excelData); err != nil {
		f.AddLog(fmt.Sprintf("创建Excel文件失败: %v", err))
	} else {
		f.AddLog("示例Excel文件创建成功: test_data.xlsx")
	}

	// 3. 测试Excel信息获取
	if info, err := f.ExcelManager.GetExcelInfo("test_data.xlsx"); err != nil {
		f.AddLog(fmt.Sprintf("获取Excel信息失败: %v", err))
	} else {
		f.AddLog(fmt.Sprintf("Excel文件信息: %+v", info))
	}

	// 4. 测试JSON转换
	if jsonData, err := f.ExcelManager.ConvertToJSON("test_data.xlsx"); err != nil {
		f.AddLog(fmt.Sprintf("Excel转JSON失败: %v", err))
	} else {
		f.AddLog(fmt.Sprintf("Excel转JSON成功: %s", jsonData))
	}

	// 更新表格显示
	f.TableData.SetCells(1, 1, "完成")
	f.UpdateStatus("Excel操作测试完成")
	f.AddLog("=== Excel操作测试结束 ===")

	log.Println("Excel操作测试完成")
}

// onJSONClick JSON操作按钮事件
func (f *MainForm) onJSONClick(sender vcl.IObject) {
	f.AddLog("=== JSON操作测试开始 ===")

	// 1. 测试JSON解析
	f.AddLog("1. 测试JSON数据解析...")
	defaultJSON := f.JSONManager.GetDefaultJSON()
	f.AddLog("默认JSON数据已准备")

	// 解析默认JSON
	if err := f.JSONManager.ParseJSON(defaultJSON); err != nil {
		f.AddLog(fmt.Sprintf("JSON解析失败: %v", err))
		return
	}
	f.AddLog("JSON解析成功")

	// 2. 测试JSON信息获取
	f.AddLog("2. 测试JSON信息获取...")
	if info, err := f.JSONManager.GetJSONInfo(defaultJSON); err != nil {
		f.AddLog(fmt.Sprintf("获取JSON信息失败: %v", err))
	} else {
		// 显示JSON信息到表格
		var infoData [][]string
		infoData = append(infoData, []string{"属性", "值"})
		for key, value := range info {
			infoData = append(infoData, []string{key, fmt.Sprintf("%v", value)})
		}
		f.SetTableData(infoData)
		f.AddLog("JSON信息已显示在表格中")
	}

	// 3. 测试JSON路径查找
	f.AddLog("3. 测试JSON路径查找...")
	if pathValue, err := f.JSONManager.GetValueByPath(defaultJSON, "姓名"); err != nil {
		f.AddLog(fmt.Sprintf("路径查找失败: %v", err))
	} else {
		f.AddLog(fmt.Sprintf("找到 '姓名' 字段的值: %s", pathValue))
	}

	// 4. 测试JSON值修改
	f.AddLog("4. 测试JSON值修改...")
	if updatedJSON, err := f.JSONManager.SetValueByPath(defaultJSON, "年龄", "26"); err != nil {
		f.AddLog(fmt.Sprintf("修改JSON值失败: %v", err))
	} else {
		f.AddLog("JSON值修改成功")
		// 显示修改后的JSON
		f.SetCurrentData(updatedJSON)
	}

	// 5. 测试JSON文件保存和加载
	f.AddLog("5. 测试JSON文件操作...")
	if err := f.JSONManager.SaveJSONToFile(defaultJSON, "test_data.json"); err != nil {
		f.AddLog(fmt.Sprintf("保存JSON文件失败: %v", err))
	} else {
		f.AddLog("JSON文件保存成功: test_data.json")

		// 加载JSON文件
		if loadedJSON, err := f.JSONManager.LoadJSONFromFile("test_data.json"); err != nil {
			f.AddLog(fmt.Sprintf("加载JSON文件失败: %v", err))
		} else {
			f.AddLog("JSON文件加载成功")
			f.AddLog(fmt.Sprintf("加载的JSON大小: %d 字节", len(loadedJSON)))
		}
	}

	// 6. 测试JSON格式验证
	f.AddLog("6. 测试JSON格式验证...")
	if isValid, message := f.JSONManager.ValidateJSON(defaultJSON); isValid {
		f.AddLog("JSON格式验证通过")
	} else {
		f.AddLog(fmt.Sprintf("JSON格式验证失败: %s", message))
	}

	f.UpdateStatus("JSON操作测试完成")
	f.AddLog("=== JSON操作测试结束 ===")

	// 更新表格状态
	f.TableData.SetCells(1, 2, "完成")

	log.Println("JSON操作测试完成")
}

// onHTTPClick HTTP操作按钮事件
func (f *MainForm) onHTTPClick(sender vcl.IObject) {
	f.AddLog("=== HTTP操作测试开始 ===")

	// 1. 测试公共API接口
	f.AddLog("1. 测试公共API接口...")
	if err := f.HTTPManager.TestAPI(); err != nil {
		f.AddLog(fmt.Sprintf("API测试失败: %v", err))
		return
	}
	f.AddLog("API测试成功")

	// 2. 测试GET请求
	f.AddLog("2. 测试HTTP GET请求...")
	getURL := "https://jsonplaceholder.typicode.com/todos/1"
	if resp, err := f.HTTPManager.GETRequest(getURL); err != nil {
		f.AddLog(fmt.Sprintf("GET请求失败: %v", err))
	} else {
		f.AddLog(fmt.Sprintf("GET请求成功，状态码: %d", resp.StatusCode))
	}

	// 3. 测试POST请求 - JSON数据
	f.AddLog("3. 测试POST请求（JSON数据）...")
	postData := map[string]interface{}{
		"title":  "测试POST请求",
		"body":   "这是一个测试数据",
		"userId": 1,
	}
	if err := f.HTTPManager.POSTRequest("https://jsonplaceholder.typicode.com/posts", postData, "application/json"); err != nil {
		f.AddLog(fmt.Sprintf("POST请求失败: %v", err))
	} else {
		f.AddLog("POST请求（JSON）成功")
	}

	// 4. 测试POST请求 - 表单数据
	f.AddLog("4. 测试POST请求（表单数据）...")
	formData := map[string]string{
		"name":  "张三",
		"email": "zhangsan@example.com",
		"phone": "13800138000",
	}
	if err := f.HTTPManager.POSTRequest("https://jsonplaceholder.typicode.com/posts", formData, "application/x-www-form-urlencoded"); err != nil {
		f.AddLog(fmt.Sprintf("POST表单请求失败: %v", err))
	} else {
		f.AddLog("POST请求（表单）成功")
	}

	// 5. 测试文件下载
	f.AddLog("5. 测试文件下载...")
	downloadURL := "https://jsonplaceholder.typicode.com/posts/1"
	if err := f.HTTPManager.DownloadFile(downloadURL, "downloads/api_response.json"); err != nil {
		f.AddLog(fmt.Sprintf("文件下载失败: %v", err))
	} else {
		f.AddLog("文件下载成功")
	}

	// 6. 显示HTTP管理器信息
	f.AddLog("6. 显示HTTP客户端配置...")
	_ = f.HTTPManager.GetClient() // 获取HTTP客户端实例
	// 注意：resty客户端的超时设置在内部管理，这里显示配置信息

	// 获取测试请求的详细信息
	requestInfo := f.HTTPManager.GetRequestInfo("GET", getURL, nil, nil)
	f.AddLog("请求信息:")
	for key, value := range requestInfo {
		f.AddLog(fmt.Sprintf("  %s: %s", key, value))
	}

	f.UpdateStatus("HTTP操作测试完成")
	f.AddLog("=== HTTP操作测试结束 ===")

	// 更新表格状态
	f.TableData.SetCells(1, 3, "完成")

	log.Println("HTTP操作测试完成")
}

// onDatabaseClick 数据库操作按钮事件
func (f *MainForm) onDatabaseClick(sender vcl.IObject) {
	f.AddLog("=== 数据库操作测试开始 ===")

	// 1. 创建示例数据库表
	f.AddLog("1. 创建示例数据库表...")
	if err := f.DBManager.CreateDefaultTable(); err != nil {
		f.AddLog(fmt.Sprintf("创建数据库表失败: %v", err))
		return
	}
	f.AddLog("示例数据库表创建成功")

	// 2. 插入测试数据
	f.AddLog("2. 插入测试数据...")
	testData := []map[string]interface{}{
		{
			"name":  "张三",
			"age":   28,
			"email": "zhangsan@example.com",
		},
		{
			"name":  "李四",
			"age":   32,
			"email": "lisi@example.com",
		},
		{
			"name":  "王五",
			"age":   25,
			"email": "wangwu@example.com",
		},
	}

	for i, data := range testData {
		f.AddLog(fmt.Sprintf("插入第%d条数据: %s", i+1, data["name"]))
		if err := f.DBManager.InsertData("users", data); err != nil {
			f.AddLog(fmt.Sprintf("插入数据失败: %v", err))
		} else {
			f.AddLog("插入数据成功")
		}
	}

	// 3. 查询数据
	f.AddLog("3. 查询数据...")
	if rows, err := f.DBManager.QueryData("users", "*", "", 100); err != nil {
		f.AddLog(fmt.Sprintf("查询数据失败: %v", err))
	} else {
		f.AddLog(fmt.Sprintf("查询成功，获取到 %d 条记录", len(rows)))
		// 显示第一条记录作为示例
		if len(rows) > 0 {
			f.AddLog("第一条记录数据:")
			for key, value := range rows[0] {
				f.AddLog(fmt.Sprintf("  %s: %v", key, value))
			}
		}
	}

	// 4. 更新数据
	f.AddLog("4. 更新数据...")
	updateData := map[string]interface{}{
		"age": 26,
	}
	if err := f.DBManager.UpdateData("users", "name = ?", updateData, []interface{}{"张三"}); err != nil {
		f.AddLog(fmt.Sprintf("更新数据失败: %v", err))
	} else {
		f.AddLog("更新数据成功")
	}

	// 5. 验证更新结果
	f.AddLog("5. 验证更新结果...")
	if rows, err := f.DBManager.QueryData("users", "*", "name = ?", 1); err != nil {
		f.AddLog(fmt.Sprintf("验证更新失败: %v", err))
	} else if len(rows) > 0 {
		f.AddLog("张三的更新后数据:")
		for key, value := range rows[0] {
			f.AddLog(fmt.Sprintf("  %s: %v", key, value))
		}
	}

	// 6. 删除数据
	f.AddLog("6. 删除数据...")
	if err := f.DBManager.DeleteData("users", "name = ?", []interface{}{"王五"}); err != nil {
		f.AddLog(fmt.Sprintf("删除数据失败: %v", err))
	} else {
		f.AddLog("删除成功（模拟模式）")
	}

	// 7. 获取数据库信息
	f.AddLog("7. 获取数据库信息...")
	dbInfo := f.DBManager.GetDatabaseInfo()
	f.AddLog("数据库信息:")
	for key, value := range dbInfo {
		f.AddLog(fmt.Sprintf("  %s: %v", key, value))
	}

	// 8. 统计查询
	f.AddLog("8. 统计查询...")
	if rows, err := f.DBManager.QueryData("users", "*", "", 100); err != nil {
		f.AddLog(fmt.Sprintf("统计查询失败: %v", err))
	} else {
		f.AddLog(fmt.Sprintf("数据库中剩余记录数: %d", len(rows)))
	}

	f.UpdateStatus("数据库操作测试完成")
	f.AddLog("=== 数据库操作测试结束 ===")

	// 安全地更新表格状态（使用正确的参数顺序：列在前，行在后）
	if f.TableData != nil && int(f.TableData.RowCount()) > 4 && int(f.TableData.ColCount()) > 1 {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("更新表格状态时发生异常: %v", r)
			}
		}()
		f.TableData.SetCells(1, 4, "完成")
		log.Printf("表格状态已更新: (列=1, 行=4) = '完成'")
	} else {
		log.Printf("表格尺寸不足以更新状态，表格尺寸: %dx%d", 
			func() int32 { defer func(){}(); return f.TableData.RowCount() }(),
			func() int32 { defer func(){}(); return f.TableData.ColCount() }())
	}

	log.Println("数据库操作测试完成")
}

// UpdateStatus 更新状态栏显示
func (f *MainForm) UpdateStatus(status string) {
	if f.StatusBar != nil && f.StatusBar.Panels().Count() > 0 {
		f.StatusBar.Panels().Items(0).SetText(fmt.Sprintf("[%s] %s",
			time.Now().Format("15:04:05"), status))
	}
	f.lastUpdate = time.Now()

	log.Printf("状态更新: %s", status)
}

// GetCurrentData 获取当前显示的数据
func (f *MainForm) GetCurrentData() string {
	if f.EditInput != nil {
		return f.EditInput.Text()
	}
	return ""
}

// SetCurrentData 设置当前显示的数据
func (f *MainForm) SetCurrentData(data string) {
	if f.EditInput != nil {
		f.EditInput.SetText(data)
	}
}

// GetTableData 获取表格数据 - 彻底修复索引越界问题
func (f *MainForm) GetTableData() [][]string {
	if f.TableData == nil {
		log.Printf("警告: 表格组件为nil，返回空数据")
		return [][]string{}
	}

	// 使用defer捕获所有可能的panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("获取表格数据时发生panic，已恢复: %v", r)
		}
	}()

	// 获取表格实际尺寸（添加详细调试信息）
	rowCount := int(f.TableData.RowCount())
	colCount := int(f.TableData.ColCount())

	log.Printf("开始获取表格数据: 表格尺寸=%dx%d", rowCount, colCount)

	// 严格的安全检查
	if rowCount <= 0 {
		log.Printf("表格行数为0，返回空数据")
		return [][]string{}
	}
	
	if colCount <= 0 {
		log.Printf("表格列数为0，返回空数据")
		return [][]string{}
	}

	// 防止超过Excel限制
	if rowCount > 100 {
		rowCount = 100
		log.Printf("表格行数超过限制，已截断为100行")
	}
	
	if colCount > 50 {
		colCount = 50
		log.Printf("表格列数超过限制，已截断为50列")
	}

	var data [][]string
	
	// 使用安全的索引遍历，添加详细的调试信息
	for r := 0; r < rowCount; r++ {
		var row []string
		
		log.Printf("处理第%d行数据", r)
		
		// 确保行索引有效
		if r >= 0 && r < int(f.TableData.RowCount()) {
			for c := 0; c < colCount; c++ {
				// 确保列索引有效
				if c >= 0 && c < int(f.TableData.ColCount()) {
					cellValue := ""
					
					// 安全地调用Cells方法
					defer func() {
						if r := recover(); r != nil {
							log.Printf("单元格访问异常: 行=%d, 列=%d, 错误=%v", r, c, r)
							cellValue = ""
						}
					}()
					
					// 使用正确的参数顺序：列, 行
					// StringGrid的Cells方法是列在前，行在后
					colParam := int32(c)
					rowParam := int32(r)
					
					log.Printf("获取单元格 (列=%d, 行=%d)", colParam, rowParam)
					
					cellValue = f.TableData.Cells(colParam, rowParam)
					
					// 确保返回值不是nil或空指针
					if cellValue == "" {
						cellValue = ""
					}
					
					log.Printf("单元格值: '%s'", cellValue)
					
					row = append(row, cellValue)
				} else {
					log.Printf("跳过无效列索引: %d", c)
					row = append(row, "")
				}
			}
			data = append(data, row)
		} else {
			log.Printf("跳过无效行索引: %d", r)
		}
	}

	// 验证数据完整性
	if len(data) == 0 {
		log.Printf("获取表格数据为空，创建默认数据")
		data = [][]string{{"Excel操作", "就绪", "点击测试"}}
	}

	// 输出完整的数据调试信息
	log.Printf("成功获取表格数据: %d行 x %d列", len(data), len(data[0]))
	for i, row := range data {
		log.Printf("第%d行: %v", i, row)
	}
	
	return data
}

// SetTableData 设置表格数据（修复索引越界问题）
func (f *MainForm) SetTableData(data [][]string) {
	if f.TableData == nil {
		log.Printf("警告: 表格组件为nil，无法设置数据")
		return
	}
	
	if len(data) == 0 {
		log.Printf("数据为空，无法设置表格")
		return
	}

	// 额外的边界检查
	if len(data) == 0 {
		log.Printf("数据为空数组，无法设置表格")
		return
	}
	
	if len(data[0]) == 0 {
		log.Printf("第一行数据为空，无法设置表格")
		return
	}

	// 设置表格尺寸（添加最大值限制）
	rowCount := int32(len(data))
	colCount := int32(len(data[0]))
	
	// 添加Excel限制检查
	maxRows := int32(100)
	maxCols := int32(50)
	
	if rowCount > maxRows {
		rowCount = maxRows
		log.Printf("行数超过限制，截断为%d行", maxRows)
	}
	
	if colCount > maxCols {
		colCount = maxCols
		log.Printf("列数超过限制，截断为%d列", maxCols)
	}

	// 添加安全检查
	if rowCount <= 0 || colCount <= 0 {
		log.Printf("表格尺寸无效: %dx%d", rowCount, colCount)
		return
	}

	log.Printf("设置表格尺寸: %dx%d", rowCount, colCount)
	
	// 使用安全的方式设置表格尺寸
	defer func() {
		if r := recover(); r != nil {
			log.Printf("设置表格尺寸时发生panic: %v", r)
		}
	}()
	
	f.TableData.SetRowCount(rowCount)
	f.TableData.SetColCount(colCount)

	// 填充数据，添加多层异常处理
	for r, row := range data {
		// 确保行索引有效
		if r < 0 || r >= int(rowCount) {
			log.Printf("跳过无效行索引: %d", r)
			continue
		}
		
		for c, cell := range row {
			// 确保列索引有效
			if c < 0 || c >= int(colCount) {
				log.Printf("跳过无效列索引: %d", c)
				continue
			}
			
			// 确保单元格值不是nil
			cellValue := ""
			if cell != "" {
				cellValue = cell
			}
			
			// 使用安全的参数顺序：列在前，行在后
			colParam := int32(c)  // 列索引
			rowParam := int32(r)  // 行索引
			
			log.Printf("设置表格单元格 (列=%d, 行=%d) = '%s'", colParam, rowParam, cellValue)
			
			// 使用defer捕获可能的panic
			defer func() {
				if r := recover(); r != nil {
					log.Printf("设置表格数据时发生异常: 行=%d, 列=%d, 值='%s', 错误=%v", r, c, cellValue, r)
				}
			}()
			
			f.TableData.SetCells(colParam, rowParam, cellValue)
		}
	}

	log.Printf("表格数据已更新，%d行 x %d列", len(data), len(data[0]))
}

// AddLog 添加日志信息
func (f *MainForm) AddLog(logText string) {
	currentText := f.GetCurrentData()
	newText := currentText + "\n\n" + fmt.Sprintf("[%s] %s",
		time.Now().Format("15:04:05"), logText)
	f.SetCurrentData(newText)

	log.Printf("UI日志: %s", logText)
}

// 确保MainForm实现UIInterface接口
var _ interfaces.UIInterface = (*MainForm)(nil)
