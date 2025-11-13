// ui/MainForm.go - 简化版主窗口表单定义
// 本文件定义了应用程序的主窗口界面，包括UI组件、事件处理和数据管理功能
// 主窗口采用现代化设计，提供Excel操作、JSON处理、HTTP请求、数据库管理和多线程测试等功能
package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"windows-gui-app/interfaces"
	"windows-gui-app/managers"

	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

// MainForm 主窗口结构体
// 封装了主窗口的所有UI组件和管理器实例，负责界面展示和用户交互
type MainForm struct {
	*vcl.TForm // 嵌入VCL表单基类，继承表单的基本功能

	// UI组件 - 界面布局相关
	PanelMain    *vcl.TPanel // 主面板 - 作为所有UI组件的容器
	PanelButtons *vcl.TPanel // 按钮面板 - 包含所有功能按钮和参数配置
	PanelStatus  *vcl.TPanel // 状态面板 - 包含状态栏
	PanelTable   *vcl.TPanel // 表格面板 - 包含数据表格和输入编辑框

	// 菜单组件 - 功能操作相关
	MainMenu      *vcl.TMainMenu // 主菜单栏
	FileMenu      *vcl.TMenuItem // 文件菜单
	EditMenu      *vcl.TMenuItem // 编辑菜单
	ToolsMenu     *vcl.TMenuItem // 工具菜单
	HelpMenu      *vcl.TMenuItem // 帮助菜单
	ClearEditItem *vcl.TMenuItem // 清空编辑框菜单项

	// 按钮组件 - 功能操作相关
	BtnGetSelections *vcl.TButton // 获取选中状态按钮 - 用于获取单选框和多选框的选中状态
	BtnImportExcel   *vcl.TButton // 导入Excel按钮 - 用于导入外部Excel文件
	BtnImportImage   *vcl.TButton // 导入图片按钮 - 用于导入图片文件
	BtnExcel         *vcl.TButton // Excel操作按钮 - 用于测试Excel相关功能
	BtnJSON          *vcl.TButton // JSON操作按钮 - 用于测试JSON处理功能
	BtnHTTP          *vcl.TButton // HTTP操作按钮 - 用于测试网络请求功能
	BtnDatabase      *vcl.TButton // 数据库操作按钮 - 用于测试数据库功能
	BtnConcurrent    *vcl.TButton // 多线程操作按钮 - 用于测试并发任务执行
	BtnResizeColumns *vcl.TButton // 调整列宽按钮 - 用于调整表格列宽
	BtnSaveConfig    *vcl.TButton // 保存配置按钮 - 用于保存当前配置到文件
	BtnWebServer     *vcl.TButton // Web服务器按钮 - 用于启动/停止Web服务器
	BtnTCPServer     *vcl.TButton // TCP服务器按钮 - 用于启动/停止TCP服务
	BtnClose         *vcl.TButton // 关闭按钮 - 用于关闭应用程序

	// 输入组件 - 参数配置相关
	EditThreadCount *vcl.TSpinEdit // 线程数量编辑框 - 用于设置并发测试的线程数
	EditTaskCount   *vcl.TSpinEdit // 任务数量编辑框 - 用于设置并发测试的任务数
	LabelThread     *vcl.TLabel    // 线程数量标签 - 线程数量编辑框的说明标签
	LabelTask       *vcl.TLabel    // 任务数量标签 - 任务数量编辑框的说明标签

	// TCP客户端连接相关
	EditTCPServerIP   *vcl.TEdit   // TCP服务器IP编辑框 - 用于输入TCP服务器IP地址
	EditTCPServerPort *vcl.TEdit   // TCP服务器端口编辑框 - 用于输入TCP服务器端口
	BtnConnectTCP     *vcl.TButton // 连接TCP服务器按钮 - 用于连接到指定的TCP服务器
	LabelTCPServer    *vcl.TLabel  // TCP服务器标签 - TCP服务器连接配置的说明标签

	// TCP数据发送相关
	EditTCPData    *vcl.TEdit   // TCP数据编辑框 - 用于输入要发送到TCP服务器的数据
	BtnSendTCPData *vcl.TButton // 发送TCP数据按钮 - 用于发送数据到TCP服务器
	LabelTCPData   *vcl.TLabel  // TCP数据标签 - TCP数据发送的说明标签
	TCPConn        net.Conn     // TCP连接 - 用于存储与TCP服务器的连接

	// 单选框组件 - 选项配置相关
	RadioOption1 *vcl.TRadioButton // 单选框1 - 用于选项1
	RadioOption2 *vcl.TRadioButton // 单选框2 - 用于选项2

	// 多选框组件 - 选项配置相关
	CheckBox1 *vcl.TCheckBox // 多选框1 - 用于选项1
	CheckBox2 *vcl.TCheckBox // 多选框2 - 用于选项2

	// 状态和显示组件 - 数据展示相关
	StatusBar *vcl.TStatusBar  // 状态栏 - 显示应用程序当前状态
	TableData *vcl.TStringGrid // 数据表格 - 用于展示表格数据
	ImageBox  *vcl.TImage      // 图片框 - 用于显示用户导入的图片
	EditInput *vcl.TMemo       // 输入编辑框 - 用于显示日志和操作结果
	LabelInfo *vcl.TLabel      // 信息标签 - 显示应用程序信息

	// 管理器 - 功能实现相关
	ExcelManager      *managers.ExcelManager      // Excel管理器 - 处理Excel文件导入导出
	JSONManager       *managers.JSONManager       // JSON管理器 - 处理JSON数据解析生成
	HTTPManager       *managers.HTTPManager       // HTTP管理器 - 处理网络请求和文件下载
	DBManager         *managers.DatabaseManager   // 数据库管理器 - 处理SQLite数据库操作
	ConcurrentManager *managers.ConcurrentManager // 并发管理器 - 处理多线程任务执行
	WebServerManager  *managers.WebServerManager  // Web服务器管理器 - 处理HTTP和WebSocket服务
	TCPServerManager  *managers.TCPServerManager  // TCP服务器管理器 - 处理TCP服务

	// 时间相关
	lastUpdate time.Time // 最后更新时间 - 用于跟踪状态更新时间
}

// NewMainForm 创建新的主窗口实例
// 返回一个初始化完成的MainForm指针，包含所有必要的管理器实例
// 该函数负责创建窗口基础结构、初始化管理器并设置基本属性
func NewMainForm() *MainForm {
	form := &MainForm{
		lastUpdate: time.Now(),
	}

	// 创建窗口
	form.TForm = vcl.Application.CreateForm()
	form.TForm.SetCaption("GO VCL 多功能演示程序")
	form.TForm.SetWidth(1024)
	form.TForm.SetHeight(568)
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

	concurrentManager := managers.NewConcurrentManager(nil)
	concurrentManager.SetUIInstance(form)

	webServerManager := managers.NewWebServerManager(nil)
	webServerManager.SetUIInstance(form)

	tcpServerManager := managers.NewTCPServerManager(nil)
	tcpServerManager.SetUIInstance(form)

	form.ExcelManager = excelManager
	form.JSONManager = jsonManager
	form.HTTPManager = httpManager
	form.DBManager = dbManager
	form.ConcurrentManager = concurrentManager
	form.WebServerManager = webServerManager
	form.TCPServerManager = tcpServerManager

	return form
}

// Show 显示窗口
// 该方法负责创建用户界面、应用样式并显示窗口
// 在调用此方法前，窗口的所有组件和管理器应该已经初始化完成
func (f *MainForm) Show() {
	f.createUI()
	f.applyUIStyles()
	f.loadConfig() // 加载配置文件

	// 在应用UI样式后，再次设置emoji兼容字体，确保按钮上的emoji图标能够正确显示
	// 这是因为皮肤加载可能会覆盖按钮的字体设置
	skinManager := managers.GetSkinManager()
	if skinManager != nil && skinManager.IsLoaded() {
		skinManager.SetEmojiCompatibleFont()
		// 再次应用按钮样式，确保emoji字体设置生效
		f.applyButtonStyles()
	}

	f.TForm.Show()
}

// createUI 创建用户界面
func (f *MainForm) createUI() {
	// 创建主菜单
	f.createMainMenu()

	// 创建主面板
	f.PanelMain = vcl.NewPanel(f.TForm)
	f.PanelMain.SetParent(f.TForm)
	f.PanelMain.SetAlign(types.AlClient)
	f.PanelMain.SetBevelOuter(types.BvNone)

	// 创建按钮面板
	f.PanelButtons = vcl.NewPanel(f.TForm)
	f.PanelButtons.SetParent(f.PanelMain)
	f.PanelButtons.SetAlign(types.AlTop)
	f.PanelButtons.SetHeight(190) // 增加高度以适应4排按钮布局
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
	// f.StatusBar.SetTextBuf("就绪")
	// f.StatusBar.Show()

	// f.StatusBar.Panels().Add()
	// f.StatusBar.Panels().Items(0).SetText("就绪")

	// 创建表格面板
	f.PanelTable = vcl.NewPanel(f.TForm)
	f.PanelTable.SetParent(f.PanelMain)
	f.PanelTable.SetAlign(types.AlClient)
	f.PanelTable.SetBevelOuter(types.BvNone)

	// 创建左侧面板，包含数据表格
	leftPanel := vcl.NewPanel(f.TForm)
	leftPanel.SetParent(f.PanelTable)
	leftPanel.SetAlign(types.AlLeft)
	leftPanel.SetWidth(600)
	leftPanel.SetBevelOuter(types.BvNone)

	// 创建数据表格（修复索引越界问题）
	f.TableData = vcl.NewStringGrid(f.TForm)
	f.TableData.SetParent(leftPanel)
	f.TableData.SetAlign(types.AlClient)

	// 首先设置固定行列
	f.TableData.SetFixedRows(1)
	f.TableData.SetFixedCols(0) // 设置为0，允许调整所有列

	// 设置表格选项，确保支持列宽调整
	gridOptions := types.TGridOptions(0) // 从零开始设置选项
	gridOptions = gridOptions | types.TGridOptions(types.GoRowSelect) |
		types.TGridOptions(types.GoColSizing) |
		types.TGridOptions(types.GoThumbTracking) |
		types.TGridOptions(types.GoColMoving) |
		types.TGridOptions(types.GoTabs) |
		types.TGridOptions(types.GoRowMoving) |
		types.TGridOptions(types.GoDrawFocusSelected)
	f.TableData.SetOptions(gridOptions)

	// 设置行列数
	f.TableData.SetRowCount(7) // 增加一行以容纳TCP服务器
	f.TableData.SetColCount(3)

	// 设置默认列宽
	f.TableData.SetDefaultColWidth(100)

	// 再次设置固定行列，确保设置生效
	f.TableData.SetFixedRows(1)
	f.TableData.SetFixedCols(0)

	log.Printf("表格初始化选项已设置: %v", gridOptions)

	// 设置初始列宽
	f.TableData.SetColWidths(0, 100) // 第一列宽度
	f.TableData.SetColWidths(1, 80)  // 第二列宽度
	f.TableData.SetColWidths(2, 420) // 第三列宽度

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
		{"Web服务器", "待测试", "启动/停止HTTP和WebSocket服务"},
		{"TCP服务器", "待测试", "启动/停止TCP服务"},
		{"多线程操作", "待测试", "并发任务执行测试"},
	}

	// 安全地设置所有单元格
	for r, row := range tableData {
		if r < 0 || r >= 7 { // 更新为7行
			log.Printf("跳过无效行索引: %d", r)
			continue
		}
		for c, cell := range row {
			if c < 0 || c >= 3 {
				log.Printf("跳过无效列索引: %d", c)
				continue
			}

			// 使用列, 行的参数顺序调用SetCells
			colParam := int32(c) // 列在前
			rowParam := int32(r) // 行在后

			log.Printf("设置表格单元格: (列=%d, 行=%d) = '%s'", colParam, rowParam, cell)
			f.TableData.SetCells(colParam, rowParam, cell)
		}
	}

	// 创建右侧面板，包含图片框和输入编辑框
	rightPanel := vcl.NewPanel(f.TForm)
	rightPanel.SetParent(f.PanelTable)
	rightPanel.SetAlign(types.AlClient)
	rightPanel.SetBevelOuter(types.BvNone)

	// 创建图片框容器面板，用于显示边框
	imagePanel := vcl.NewPanel(f.TForm)
	imagePanel.SetParent(rightPanel)
	imagePanel.SetAlign(types.AlTop)
	imagePanel.SetHeight(200)
	imagePanel.SetBevelOuter(types.BvLowered) // 设置边框样式

	// 创建图片框
	f.ImageBox = vcl.NewImage(f.TForm)
	f.ImageBox.SetParent(imagePanel)
	f.ImageBox.SetAlign(types.AlClient) // 填充整个面板
	f.ImageBox.SetCenter(true)
	f.ImageBox.SetStretch(true)
	f.ImageBox.SetProportional(true)

	// 创建输入编辑框
	f.EditInput = vcl.NewMemo(f.TForm)
	f.EditInput.SetParent(rightPanel)
	f.EditInput.SetAlign(types.AlClient)
	f.EditInput.SetScrollBars(types.SsBoth)
	f.EditInput.SetReadOnly(true)
	f.EditInput.SetText("欢迎使用 GO VCL 多功能演示程序！\n\n本程序提供以下功能：\n\n1. Excel操作 - 导入/导出Excel文件\n2. JSON操作 - 解析/生成JSON数据\n3. HTTP操作 - 网络请求和文件下载\n4. 数据库操作 - SQLite数据库管理\n5. Web服务器 - 启动/停止HTTP和WebSocket服务\n6. TCP服务器 - 启动/停止TCP服务\n7. 图片导入 - 导入并显示图片\n\n请点击上方按钮开始使用。")

	// 设置事件处理
	f.setupEvents()
}

// createMainMenu 创建主菜单
func (f *MainForm) createMainMenu() {
	// 创建主菜单栏
	f.MainMenu = vcl.NewMainMenu(f.TForm)

	// 创建文件菜单
	f.FileMenu = vcl.NewMenuItem(f.TForm)
	f.FileMenu.SetCaption("文件(&F)")
	f.MainMenu.Items().Add(f.FileMenu)

	// 创建编辑菜单
	f.EditMenu = vcl.NewMenuItem(f.TForm)
	f.EditMenu.SetCaption("编辑(&E)")
	f.MainMenu.Items().Add(f.EditMenu)

	// 添加清空编辑框菜单项到编辑菜单
	f.ClearEditItem = vcl.NewMenuItem(f.TForm)
	f.ClearEditItem.SetCaption("清空编辑框(&C)")
	f.ClearEditItem.SetOnClick(f.onClearEditClick)
	f.EditMenu.Add(f.ClearEditItem)

	// 创建工具菜单
	f.ToolsMenu = vcl.NewMenuItem(f.TForm)
	f.ToolsMenu.SetCaption("工具(&T)")
	f.MainMenu.Items().Add(f.ToolsMenu)

	// 创建帮助菜单
	f.HelpMenu = vcl.NewMenuItem(f.TForm)
	f.HelpMenu.SetCaption("帮助(&H)")
	f.MainMenu.Items().Add(f.HelpMenu)
}

// createButtons 创建按钮组件
func (f *MainForm) createButtons() {
	// 创建一个现代化的按钮面板布局
	// 使用表格布局来优化按钮排列

	// 创建参数配置区域面板
	configPanel := vcl.NewPanel(f.TForm)
	configPanel.SetParent(f.PanelButtons)
	configPanel.SetAlign(types.AlLeft)
	configPanel.SetWidth(330) // 增加宽度以容纳多选框
	configPanel.SetBevelOuter(types.BvNone)
	configPanel.SetHeight(150) // 增加高度以容纳TCP客户端连接配置和数据发送控件

	// 创建参数标签和编辑框
	f.LabelThread = vcl.NewLabel(f.TForm)
	f.LabelThread.SetParent(configPanel)
	f.LabelThread.SetCaption("线程数量:")
	f.LabelThread.SetLeft(10)
	f.LabelThread.SetTop(10)
	f.LabelThread.SetWidth(70)

	f.EditThreadCount = vcl.NewSpinEdit(f.TForm)
	f.EditThreadCount.SetParent(configPanel)
	f.EditThreadCount.SetLeft(85)
	f.EditThreadCount.SetTop(8)
	f.EditThreadCount.SetWidth(80)
	f.EditThreadCount.SetMinValue(1)
	f.EditThreadCount.SetMaxValue(100)
	f.EditThreadCount.SetValue(8)

	f.LabelTask = vcl.NewLabel(f.TForm)
	f.LabelTask.SetParent(configPanel)
	f.LabelTask.SetCaption("任务数量:")
	f.LabelTask.SetLeft(10)
	f.LabelTask.SetTop(35)
	f.LabelTask.SetWidth(70)

	f.EditTaskCount = vcl.NewSpinEdit(f.TForm)
	f.EditTaskCount.SetParent(configPanel)
	f.EditTaskCount.SetLeft(85)
	f.EditTaskCount.SetTop(33)
	f.EditTaskCount.SetWidth(80)
	f.EditTaskCount.SetMinValue(1)
	f.EditTaskCount.SetMaxValue(1000)
	f.EditTaskCount.SetValue(50)

	// 创建单选框
	f.RadioOption1 = vcl.NewRadioButton(f.TForm)
	f.RadioOption1.SetParent(configPanel)
	f.RadioOption1.SetCaption("选项 1")
	f.RadioOption1.SetLeft(10)
	f.RadioOption1.SetTop(60)
	f.RadioOption1.SetWidth(70)
	f.RadioOption1.SetChecked(true) // 默认选中第一个选项

	f.RadioOption2 = vcl.NewRadioButton(f.TForm)
	f.RadioOption2.SetParent(configPanel)
	f.RadioOption2.SetCaption("选项 2")
	f.RadioOption2.SetLeft(85)
	f.RadioOption2.SetTop(60)
	f.RadioOption2.SetWidth(70)

	// 创建多选框
	f.CheckBox1 = vcl.NewCheckBox(f.TForm)
	f.CheckBox1.SetParent(configPanel)
	f.CheckBox1.SetCaption("多选 1")
	f.CheckBox1.SetLeft(170)
	f.CheckBox1.SetTop(60)
	f.CheckBox1.SetWidth(70)

	f.CheckBox2 = vcl.NewCheckBox(f.TForm)
	f.CheckBox2.SetParent(configPanel)
	f.CheckBox2.SetCaption("多选 2")
	f.CheckBox2.SetLeft(245)
	f.CheckBox2.SetTop(60)
	f.CheckBox2.SetWidth(70)

	// 创建TCP客户端连接配置
	f.LabelTCPServer = vcl.NewLabel(f.TForm)
	f.LabelTCPServer.SetParent(configPanel)
	f.LabelTCPServer.SetCaption("TCP服务器:")
	f.LabelTCPServer.SetLeft(10)
	f.LabelTCPServer.SetTop(85)
	f.LabelTCPServer.SetWidth(70)

	f.EditTCPServerIP = vcl.NewEdit(f.TForm)
	f.EditTCPServerIP.SetParent(configPanel)
	f.EditTCPServerIP.SetLeft(85)
	f.EditTCPServerIP.SetTop(83)
	f.EditTCPServerIP.SetWidth(80)
	f.EditTCPServerIP.SetText("127.0.0.1")

	f.LabelTCPServer = vcl.NewLabel(f.TForm)
	f.LabelTCPServer.SetParent(configPanel)
	f.LabelTCPServer.SetCaption(":")
	f.LabelTCPServer.SetLeft(170)
	f.LabelTCPServer.SetTop(85)
	f.LabelTCPServer.SetWidth(10)

	f.EditTCPServerPort = vcl.NewEdit(f.TForm)
	f.EditTCPServerPort.SetParent(configPanel)
	f.EditTCPServerPort.SetLeft(180)
	f.EditTCPServerPort.SetTop(83)
	f.EditTCPServerPort.SetWidth(50)
	f.EditTCPServerPort.SetText("8082")

	f.BtnConnectTCP = vcl.NewButton(f.TForm)
	f.BtnConnectTCP.SetParent(configPanel)
	f.BtnConnectTCP.SetCaption("连接TCP")
	f.BtnConnectTCP.SetLeft(240)
	f.BtnConnectTCP.SetTop(83)
	f.BtnConnectTCP.SetWidth(75)
	f.BtnConnectTCP.SetHeight(22)
	f.BtnConnectTCP.SetOnClick(f.onConnectTCPClick)

	// 创建TCP数据发送配置
	f.LabelTCPData = vcl.NewLabel(f.TForm)
	f.LabelTCPData.SetParent(configPanel)
	f.LabelTCPData.SetCaption("发送数据:")
	f.LabelTCPData.SetLeft(10)
	f.LabelTCPData.SetTop(110)
	f.LabelTCPData.SetWidth(70)

	f.EditTCPData = vcl.NewEdit(f.TForm)
	f.EditTCPData.SetParent(configPanel)
	f.EditTCPData.SetLeft(85)
	f.EditTCPData.SetTop(108)
	f.EditTCPData.SetWidth(145)
	f.EditTCPData.SetText("Hello TCP Server!")

	f.BtnSendTCPData = vcl.NewButton(f.TForm)
	f.BtnSendTCPData.SetParent(configPanel)
	f.BtnSendTCPData.SetCaption("发送数据")
	f.BtnSendTCPData.SetLeft(240)
	f.BtnSendTCPData.SetTop(108)
	f.BtnSendTCPData.SetWidth(75)
	f.BtnSendTCPData.SetHeight(22)
	f.BtnSendTCPData.SetOnClick(f.onSendTCPDataClick)

	// 创建功能按钮面板
	buttonPanel := vcl.NewPanel(f.TForm)
	buttonPanel.SetParent(f.PanelButtons)
	buttonPanel.SetAlign(types.AlClient)
	buttonPanel.SetBevelOuter(types.BvNone)

	// 设置按钮面板高度以适应4排按钮的布局
	buttonPanel.SetHeight(180) // 增加高度以适应4排按钮的布局

	// 创建按钮 - 使用网格布局优化视觉效果

	// 第一行按钮 (3个)
	f.BtnSaveConfig = vcl.NewButton(f.TForm)
	f.BtnSaveConfig.SetParent(buttonPanel)
	f.BtnSaveConfig.SetCaption("💾 保存配置")
	f.BtnSaveConfig.SetWidth(100)
	f.BtnSaveConfig.SetHeight(30)
	f.BtnSaveConfig.SetLeft(10)
	f.BtnSaveConfig.SetTop(5)
	f.BtnSaveConfig.SetOnClick(f.onSaveConfigClick)

	f.BtnGetSelections = vcl.NewButton(f.TForm)
	f.BtnGetSelections.SetParent(buttonPanel)
	f.BtnGetSelections.SetCaption("🔍 获取选中状态")
	f.BtnGetSelections.SetWidth(100)
	f.BtnGetSelections.SetHeight(30)
	f.BtnGetSelections.SetLeft(115)
	f.BtnGetSelections.SetTop(5)
	f.BtnGetSelections.SetOnClick(f.onGetSelectionsClick)

	f.BtnImportExcel = vcl.NewButton(f.TForm)
	f.BtnImportExcel.SetParent(buttonPanel)
	f.BtnImportExcel.SetCaption("📊 导入Excel")
	f.BtnImportExcel.SetWidth(100)
	f.BtnImportExcel.SetHeight(30)
	f.BtnImportExcel.SetLeft(220)
	f.BtnImportExcel.SetTop(5)
	f.BtnImportExcel.SetOnClick(f.onImportExcelClick)

	// 第二行按钮 (3个)
	f.BtnImportImage = vcl.NewButton(f.TForm)
	f.BtnImportImage.SetParent(buttonPanel)
	f.BtnImportImage.SetCaption("🖼️ 导入图片")
	f.BtnImportImage.SetWidth(100)
	f.BtnImportImage.SetHeight(30)
	f.BtnImportImage.SetLeft(10)
	f.BtnImportImage.SetTop(45)
	f.BtnImportImage.SetOnClick(f.onImportImageClick)

	f.BtnExcel = vcl.NewButton(f.TForm)
	f.BtnExcel.SetParent(buttonPanel)
	f.BtnExcel.SetCaption("📈 Excel操作")
	f.BtnExcel.SetWidth(100)
	f.BtnExcel.SetHeight(30)
	f.BtnExcel.SetLeft(115)
	f.BtnExcel.SetTop(45)
	f.BtnExcel.SetOnClick(f.onExcelClick)

	f.BtnJSON = vcl.NewButton(f.TForm)
	f.BtnJSON.SetParent(buttonPanel)
	f.BtnJSON.SetCaption("📋 JSON操作")
	f.BtnJSON.SetWidth(100)
	f.BtnJSON.SetHeight(30)
	f.BtnJSON.SetLeft(220)
	f.BtnJSON.SetTop(45)
	f.BtnJSON.SetOnClick(f.onJSONClick)

	// 第三行按钮 (3个)
	f.BtnHTTP = vcl.NewButton(f.TForm)
	f.BtnHTTP.SetParent(buttonPanel)
	f.BtnHTTP.SetCaption("🌐 HTTP操作")
	f.BtnHTTP.SetWidth(100)
	f.BtnHTTP.SetHeight(30)
	f.BtnHTTP.SetLeft(10)
	f.BtnHTTP.SetTop(85)
	f.BtnHTTP.SetOnClick(f.onHTTPClick)

	f.BtnDatabase = vcl.NewButton(f.TForm)
	f.BtnDatabase.SetParent(buttonPanel)
	f.BtnDatabase.SetCaption("🗄️ 数据库操作")
	f.BtnDatabase.SetWidth(100)
	f.BtnDatabase.SetHeight(30)
	f.BtnDatabase.SetLeft(115)
	f.BtnDatabase.SetTop(85)
	f.BtnDatabase.SetOnClick(f.onDatabaseClick)

	f.BtnConcurrent = vcl.NewButton(f.TForm)
	f.BtnConcurrent.SetParent(buttonPanel)
	f.BtnConcurrent.SetCaption("⚡ 多线程测试")
	f.BtnConcurrent.SetWidth(100)
	f.BtnConcurrent.SetHeight(30)
	f.BtnConcurrent.SetLeft(220)
	f.BtnConcurrent.SetTop(85)
	f.BtnConcurrent.SetOnClick(f.onConcurrentClick)

	// 第四行按钮 (4个)
	f.BtnResizeColumns = vcl.NewButton(f.TForm)
	f.BtnResizeColumns.SetParent(buttonPanel)
	f.BtnResizeColumns.SetCaption("📏 调整列宽")
	f.BtnResizeColumns.SetWidth(90)
	f.BtnResizeColumns.SetHeight(30)
	f.BtnResizeColumns.SetLeft(10)
	f.BtnResizeColumns.SetTop(125)
	f.BtnResizeColumns.SetOnClick(f.onResizeColumnsClick)

	f.BtnWebServer = vcl.NewButton(f.TForm)
	f.BtnWebServer.SetParent(buttonPanel)
	f.BtnWebServer.SetCaption("🌐 Web服务器")
	f.BtnWebServer.SetWidth(90)
	f.BtnWebServer.SetHeight(30)
	f.BtnWebServer.SetLeft(105)
	f.BtnWebServer.SetTop(125)
	f.BtnWebServer.SetOnClick(f.onWebServerClick)

	f.BtnTCPServer = vcl.NewButton(f.TForm)
	f.BtnTCPServer.SetParent(buttonPanel)
	f.BtnTCPServer.SetCaption("🔌 TCP服务")
	f.BtnTCPServer.SetWidth(90)
	f.BtnTCPServer.SetHeight(30)
	f.BtnTCPServer.SetLeft(200)
	f.BtnTCPServer.SetTop(125)
	f.BtnTCPServer.SetOnClick(f.onTCPServerClick)

	f.BtnClose = vcl.NewButton(f.TForm)
	f.BtnClose.SetParent(buttonPanel)
	f.BtnClose.SetCaption("❌ 关闭")
	f.BtnClose.SetWidth(90)
	f.BtnClose.SetHeight(30)
	f.BtnClose.SetLeft(295)
	f.BtnClose.SetTop(125)
	f.BtnClose.SetOnClick(f.onCloseClick)

	// 创建信息标签，放置在关闭按钮后面
	f.LabelInfo = vcl.NewLabel(f.TForm)
	f.LabelInfo.SetParent(buttonPanel)
	f.LabelInfo.SetCaption("GO VCL 多功能演示程序 v1.0")
	f.LabelInfo.Font().SetSize(14)
	f.LabelInfo.Font().SetColor(0x0000FF) // 设置字体颜色为红色
	fontStyle := f.LabelInfo.Font().Style()
	fontStyle = fontStyle | types.TFontStyles(types.FsBold)
	f.LabelInfo.Font().SetStyle(fontStyle)
	f.LabelInfo.SetLeft(400)          // 修改位置到关闭按钮后面
	f.LabelInfo.SetTop(125)           // 调整位置到第四行
	f.LabelInfo.SetColor(0xFFFFFF)    // 设置背景色为白色
	f.LabelInfo.SetTransparent(false) // 确保背景色不透明
}

// setupEvents 设置事件处理
func (f *MainForm) setupEvents() {
	// 设置表格事件处理
	if f.TableData != nil {
		// TStringGrid没有直接的列宽调整事件，但我们可以使用鼠标事件来检测
		f.TableData.SetOnMouseDown(f.onTableMouseDown)
		f.TableData.SetOnMouseUp(f.onTableMouseUp)
	}

	// 设置单选框事件处理
	if f.RadioOption1 != nil {
		f.RadioOption1.SetOnClick(f.onRadioOption1Click)
	}

	if f.RadioOption2 != nil {
		f.RadioOption2.SetOnClick(f.onRadioOption2Click)
	}

	// 设置多选框事件处理
	if f.CheckBox1 != nil {
		f.CheckBox1.SetOnClick(f.onCheckBox1Click)
	}

	if f.CheckBox2 != nil {
		f.CheckBox2.SetOnClick(f.onCheckBox2Click)
	}
}

// onTableMouseDown 表格鼠标按下事件
func (f *MainForm) onTableMouseDown(sender vcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
	// 记录鼠标按下时的位置，用于后续判断是否进行了列宽调整
	log.Printf("表格鼠标按下事件: 按钮=%v, 位置=(%d,%d)", button, x, y)

	// 检查是否在列边界附近（用于列宽调整）
	if f.TableData != nil {
		// 获取表格选项，确认列宽调整已启用
		options := f.TableData.Options()
		isColSizingEnabled := (options & types.TGridOptions(types.GoColSizing)) != 0
		log.Printf("列宽调整选项已启用: %v", isColSizingEnabled)

		// 如果列宽调整未启用，重新设置表格选项
		if !isColSizingEnabled {
			log.Printf("重新启用列宽调整选项...")
			newOptions := types.TGridOptions(0)
			newOptions = newOptions | types.TGridOptions(types.GoRowSelect) |
				types.TGridOptions(types.GoColSizing) |
				types.TGridOptions(types.GoThumbTracking) |
				types.TGridOptions(types.GoColMoving) |
				types.TGridOptions(types.GoTabs) |
				types.TGridOptions(types.GoRowMoving) |
				types.TGridOptions(types.GoDrawFocusSelected)
			f.TableData.SetOptions(newOptions)
			log.Printf("表格选项已重新设置: %v", newOptions)
		}
	}
}

// onTableMouseUp 表格鼠标释放事件
// 该方法处理表格组件的鼠标释放事件，主要用于完成表格操作后的状态更新
// 特别关注列宽调整后的状态检查和日志记录
func (f *MainForm) onTableMouseUp(sender vcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
	// 记录鼠标释放时的位置
	log.Printf("表格鼠标释放事件: 按钮=%v, 位置=(%d,%d)", button, x, y)

	// 检查表格选项是否仍然包含列宽调整
	if f.TableData != nil {
		options := f.TableData.Options()
		isColSizingEnabled := (options & types.TGridOptions(types.GoColSizing)) != 0
		log.Printf("列宽调整选项状态: %v", isColSizingEnabled)

		// 输出当前列宽信息
		colCount := f.TableData.ColCount()
		var colWidths []int32
		for i := int32(0); i < colCount; i++ {
			width := f.TableData.ColWidths(i)
			colWidths = append(colWidths, width)
		}
		log.Printf("当前列宽: %v", colWidths)
	}
}

// 事件处理函数

// onShow 窗口显示事件
// 该方法在窗口显示时被调用，用于初始化窗口状态和显示欢迎信息
// 这是窗口生命周期中的重要事件，通常用于完成界面初始化后的准备工作
func (f *MainForm) onShow(sender vcl.IObject) {
	f.UpdateStatus("程序已启动就绪")
	log.Println("主窗口已显示")
}

// onClose 窗口关闭事件
// 该方法在用户尝试关闭窗口时被调用，用于处理程序退出前的清理工作
// 通过设置action参数控制窗口关闭行为，CaFree表示关闭后释放窗口资源
func (f *MainForm) onClose(sender vcl.IObject, action *types.TCloseAction) {
	log.Println("程序正在退出...")
	f.UpdateStatus("程序退出")
	*action = types.CaFree // 释放窗口
}

// onCloseClick 关闭按钮点击事件
// 该方法处理用户点击关闭按钮的操作，直接调用窗口的Close方法
// 这是用户主动退出应用程序的主要方式之一
func (f *MainForm) onCloseClick(sender vcl.IObject) {
	f.TForm.Close()
}

// onImportExcelClick 导入Excel按钮事件
func (f *MainForm) onImportExcelClick(sender vcl.IObject) {
	f.AddLog("=== Excel文件导入开始 ===")

	// 创建文件对话框
	dialog := vcl.NewOpenDialog(f.TForm)
	dialog.SetFilter("Excel文件 (*.xlsx)|*.xlsx|所有文件 (*.*)|*.*")
	dialog.SetTitle("选择要导入的Excel文件")

	// 显示文件选择对话框
	if dialog.Execute() {
		filePath := dialog.FileName()
		f.AddLog(fmt.Sprintf("选择文件: %s", filePath))

		// 调用Excel管理器的导入功能
		if err := f.ExcelManager.ImportExcel(filePath); err != nil {
			f.AddLog(fmt.Sprintf("导入Excel文件失败: %v", err))
			f.UpdateStatus("Excel文件导入失败")
		} else {
			f.AddLog("Excel文件导入成功！")
			f.AddLog("数据已显示在下方表格中")
			f.UpdateStatus("Excel文件导入完成")
		}
	} else {
		f.AddLog("用户取消了文件选择")
		f.UpdateStatus("Excel文件导入已取消")
	}

	f.AddLog("=== Excel文件导入结束 ===")
}

// onImportImageClick 导入图片按钮事件
func (f *MainForm) onImportImageClick(sender vcl.IObject) {
	f.AddLog("=== 图片导入开始 ===")

	// 创建文件对话框
	dialog := vcl.NewOpenDialog(f.TForm)
	dialog.SetFilter("图片文件 (*.jpg;*.jpeg;*.png;*.bmp;*.gif)|*.jpg;*.jpeg;*.png;*.bmp;*.gif|所有文件 (*.*)|*.*")
	dialog.SetTitle("选择要导入的图片文件")

	// 显示文件选择对话框
	if dialog.Execute() {
		filePath := dialog.FileName()
		f.AddLog(fmt.Sprintf("选择图片文件: %s", filePath))

		// 加载图片到图片框
		picture := vcl.NewPicture()
		picture.LoadFromFile(filePath)
		f.ImageBox.Picture().Assign(picture)
		f.AddLog("图片加载成功！")
		f.AddLog(fmt.Sprintf("图片尺寸: %d x %d", picture.Width(), picture.Height()))
		f.UpdateStatus("图片导入完成")
	} else {
		f.AddLog("用户取消了文件选择")
		f.UpdateStatus("图片导入已取消")
	}

	f.AddLog("=== 图片导入结束 ===")
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
			func() int32 { defer func() {}(); return f.TableData.RowCount() }(),
			func() int32 { defer func() {}(); return f.TableData.ColCount() }())
	}

	log.Println("数据库操作测试完成")
}

// onRadioOption1Click 单选框1选中事件
func (f *MainForm) onRadioOption1Click(sender vcl.IObject) {
	if f.RadioOption1 != nil && f.RadioOption1.Checked() {
		f.AddLog("单选框1被选中: 选项 1")
		log.Println("单选框1被选中: 选项 1")

		// 确保单选框2不被选中
		if f.RadioOption2 != nil {
			f.RadioOption2.SetChecked(false)
		}
	}
}

// onRadioOption2Click 单选框2选中事件
func (f *MainForm) onRadioOption2Click(sender vcl.IObject) {
	if f.RadioOption2 != nil && f.RadioOption2.Checked() {
		f.AddLog("单选框2被选中: 选项 2")
		log.Println("单选框2被选中: 选项 2")

		// 确保单选框1不被选中
		if f.RadioOption1 != nil {
			f.RadioOption1.SetChecked(false)
		}
	}
}

// onCheckBox1Click 多选框1点击事件
func (f *MainForm) onCheckBox1Click(sender vcl.IObject) {
	if f.CheckBox1 != nil {
		if f.CheckBox1.Checked() {
			f.AddLog("多选框1被选中: 多选 1")
			log.Println("多选框1被选中: 多选 1")
		} else {
			f.AddLog("多选框1取消选中: 多选 1")
			log.Println("多选框1取消选中: 多选 1")
		}
	}
}

// onCheckBox2Click 多选框2点击事件
func (f *MainForm) onCheckBox2Click(sender vcl.IObject) {
	if f.CheckBox2 != nil {
		if f.CheckBox2.Checked() {
			f.AddLog("多选框2被选中: 多选 2")
			log.Println("多选框2被选中: 多选 2")
		} else {
			f.AddLog("多选框2取消选中: 多选 2")
			log.Println("多选框2取消选中: 多选 2")
		}
	}
}

// onClearEditClick 清空编辑框菜单项点击事件
func (f *MainForm) onClearEditClick(sender vcl.IObject) {
	if f.EditInput != nil {
		f.EditInput.SetText("")
		f.AddLog("编辑框内容已清空")
		log.Println("编辑框内容已清空")
	}
}

// onGetSelectionsClick 获取选中状态按钮点击事件
func (f *MainForm) onGetSelectionsClick(sender vcl.IObject) {
	f.AddLog("=== 获取控件选中状态 ===")

	// 获取单选框状态
	if f.RadioOption1 != nil {
		isChecked := f.RadioOption1.Checked()
		status := "未选中"
		if isChecked {
			status = "已选中"
		}
		f.AddLog(fmt.Sprintf("单选框1: %s", status))
		log.Printf("单选框1状态: %s", status)
	} else {
		f.AddLog("单选框1: 控件未初始化")
	}

	if f.RadioOption2 != nil {
		isChecked := f.RadioOption2.Checked()
		status := "未选中"
		if isChecked {
			status = "已选中"
		}
		f.AddLog(fmt.Sprintf("单选框2: %s", status))
		log.Printf("单选框2状态: %s", status)
	} else {
		f.AddLog("单选框2: 控件未初始化")
	}

	// 获取多选框状态
	if f.CheckBox1 != nil {
		isChecked := f.CheckBox1.Checked()
		status := "未选中"
		if isChecked {
			status = "已选中"
		}
		f.AddLog(fmt.Sprintf("多选框1: %s", status))
		log.Printf("多选框1状态: %s", status)
	} else {
		f.AddLog("多选框1: 控件未初始化")
	}

	if f.CheckBox2 != nil {
		isChecked := f.CheckBox2.Checked()
		status := "未选中"
		if isChecked {
			status = "已选中"
		}
		f.AddLog(fmt.Sprintf("多选框2: %s", status))
		log.Printf("多选框2状态: %s", status)
	} else {
		f.AddLog("多选框2: 控件未初始化")
	}

	f.AddLog("=== 选中状态获取完成 ===")
	f.UpdateStatus("已获取所有控件选中状态")
}

// onResizeColumnsClick 调整列宽按钮事件
func (f *MainForm) onResizeColumnsClick(sender vcl.IObject) {
	f.AddLog("=== 表格列宽调整测试开始 ===")

	if f.TableData == nil {
		f.AddLog("表格组件未初始化")
		return
	}

	// 获取当前列数
	colCount := f.TableData.ColCount()
	f.AddLog(fmt.Sprintf("当前表格列数: %d", colCount))

	// 检查表格选项，确保列宽调整已启用
	options := f.TableData.Options()
	isColSizingEnabled := (options & types.TGridOptions(types.GoColSizing)) != 0
	isThumbTrackingEnabled := (options & types.TGridOptions(types.GoThumbTracking)) != 0

	f.AddLog(fmt.Sprintf("列宽调整选项状态: GoColSizing=%v, GoThumbTracking=%v", isColSizingEnabled, isThumbTrackingEnabled))

	// 如果列宽调整未启用，重新设置表格选项
	if !isColSizingEnabled || !isThumbTrackingEnabled {
		f.AddLog("重新启用列宽调整选项...")
		newOptions := types.TGridOptions(0)
		newOptions = newOptions | types.TGridOptions(types.GoRowSelect) |
			types.TGridOptions(types.GoColSizing) |
			types.TGridOptions(types.GoThumbTracking) |
			types.TGridOptions(types.GoColMoving) |
			types.TGridOptions(types.GoTabs) |
			types.TGridOptions(types.GoRowMoving) |
			types.TGridOptions(types.GoDrawFocusSelected)
		f.TableData.SetOptions(newOptions)
		f.AddLog(fmt.Sprintf("表格选项已重新设置: %v", newOptions))
	}

	// 确保固定列设置为0，允许调整所有列
	f.TableData.SetFixedCols(0)
	f.AddLog("固定列数已设置为0，允许调整所有列")

	// 设置每列的宽度为不同值，以便测试手动调整
	for i := int32(0); i < colCount; i++ {
		newWidth := int32(80 + i*30) // 每列宽度递增
		f.TableData.SetColWidths(i, newWidth)
		f.AddLog(fmt.Sprintf("设置第%d列宽度为: %d", i+1, newWidth))
	}

	// 强制刷新表格显示
	f.TableData.Invalidate()
	f.AddLog("表格显示已刷新")

	// 输出最终列宽信息
	var finalColWidths []int32
	for i := int32(0); i < colCount; i++ {
		width := f.TableData.ColWidths(i)
		finalColWidths = append(finalColWidths, width)
	}
	f.AddLog(fmt.Sprintf("最终列宽: %v", finalColWidths))

	f.AddLog("表格列宽已调整，请尝试手动拖动列边界调整宽度")
	f.AddLog("提示：将鼠标移动到列边界处，当光标变为双向箭头时拖动调整宽度")
	f.AddLog("=== 表格列宽调整测试结束 ===")
}

// onWebServerClick Web服务器按钮事件
func (f *MainForm) onWebServerClick(sender vcl.IObject) {
	if f.WebServerManager == nil {
		f.AddLog("Web服务器管理器未初始化")
		return
	}

	f.AddLog("=== Web服务器操作开始 ===")

	// 检查服务器当前状态
	isRunning := f.WebServerManager.IsRunning()
	if isRunning {
		// 服务器正在运行，执行停止操作
		f.AddLog("正在停止Web服务器...")
		if err := f.WebServerManager.StopServer(); err != nil {
			f.AddLog(fmt.Sprintf("停止Web服务器失败: %v", err))
			f.UpdateStatus("停止Web服务器失败")
		} else {
			f.AddLog("Web服务器已成功停止")
			f.BtnWebServer.SetCaption("🌐 Web服务器")
			f.UpdateStatus("Web服务器已停止")

			// 更新表格状态
			if f.TableData != nil && int(f.TableData.RowCount()) > 5 && int(f.TableData.ColCount()) > 1 {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("更新表格状态时发生异常: %v", r)
					}
				}()
				f.TableData.SetCells(1, 5, "已停止") // Web服务器在第6行(索引5)
				log.Printf("表格状态已更新: (列=1, 行=5) = '已停止'")
			}
		}
	} else {
		// 服务器未运行，执行启动操作
		f.AddLog("正在启动Web服务器...")
		if err := f.WebServerManager.StartServer(); err != nil {
			f.AddLog(fmt.Sprintf("启动Web服务器失败: %v", err))
			f.UpdateStatus("启动Web服务器失败")
		} else {
			f.AddLog("Web服务器已成功启动")
			f.BtnWebServer.SetCaption("⏹️ 停止服务器")
			f.UpdateStatus("Web服务器运行中")

			// 获取并显示服务器信息
			serverInfo := f.WebServerManager.GetServerInfo()
			f.AddLog("服务器信息:")
			for key, value := range serverInfo {
				f.AddLog(fmt.Sprintf("  %s: %v", key, value))
			}

			// 更新表格状态
			if f.TableData != nil && int(f.TableData.RowCount()) > 5 && int(f.TableData.ColCount()) > 1 {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("更新表格状态时发生异常: %v", r)
					}
				}()
				f.TableData.SetCells(1, 5, "运行中") // Web服务器在第6行(索引5)
				log.Printf("表格状态已更新: (列=1, 行=5) = '运行中'")
			}
		}
	}

	f.AddLog("=== Web服务器操作结束 ===")
	log.Println("Web服务器操作完成")
}

// onTCPServerClick TCP服务器按钮事件
func (f *MainForm) onTCPServerClick(sender vcl.IObject) {
	if f.TCPServerManager == nil {
		f.AddLog("TCP服务器管理器未初始化")
		return
	}

	f.AddLog("=== TCP服务器操作开始 ===")

	// 检查服务器当前状态
	isRunning := f.TCPServerManager.IsRunning()
	if isRunning {
		// 服务器正在运行，执行停止操作
		f.AddLog("正在停止TCP服务器...")
		if err := f.TCPServerManager.StopServer(); err != nil {
			f.AddLog(fmt.Sprintf("停止TCP服务器失败: %v", err))
			f.UpdateStatus("停止TCP服务器失败")
		} else {
			f.AddLog("TCP服务器已成功停止")
			f.BtnTCPServer.SetCaption("🔌 TCP服务")
			f.UpdateStatus("TCP服务器已停止")

			// 更新表格状态
			if f.TableData != nil && int(f.TableData.RowCount()) > 6 && int(f.TableData.ColCount()) > 1 {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("更新表格状态时发生异常: %v", r)
					}
				}()
				f.TableData.SetCells(1, 6, "已停止") // TCP服务器在第7行(索引6)
				log.Printf("表格状态已更新: (列=1, 行=6) = '已停止'")
			}
		}
	} else {
		// 服务器未运行，执行启动操作
		f.AddLog("正在启动TCP服务器...")
		if err := f.TCPServerManager.StartServer(); err != nil {
			f.AddLog(fmt.Sprintf("启动TCP服务器失败: %v", err))
			f.UpdateStatus("启动TCP服务器失败")
		} else {
			f.AddLog("TCP服务器已成功启动")
			f.BtnTCPServer.SetCaption("⏹️ 停止服务")
			f.UpdateStatus("TCP服务器运行中")

			// 获取并显示服务器信息
			serverInfo := f.TCPServerManager.GetServerInfo()
			f.AddLog("服务器信息:")
			for key, value := range serverInfo {
				f.AddLog(fmt.Sprintf("  %s: %v", key, value))
			}

			// 更新表格状态
			if f.TableData != nil && int(f.TableData.RowCount()) > 6 && int(f.TableData.ColCount()) > 1 {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("更新表格状态时发生异常: %v", r)
					}
				}()
				f.TableData.SetCells(1, 6, "运行中") // TCP服务器在第7行(索引6)
				log.Printf("表格状态已更新: (列=1, 行=6) = '运行中'")
			}
		}
	}

	f.AddLog("=== TCP服务器操作结束 ===")
	log.Println("TCP服务器操作完成")
}

// onConnectTCPClick 连接TCP服务器按钮事件
func (f *MainForm) onConnectTCPClick(sender vcl.IObject) {
	// 如果已有连接，先关闭
	if f.TCPConn != nil {
		f.TCPConn.Close()
		f.TCPConn = nil
		f.AddLog("已关闭之前的TCP连接")
	}

	// 获取用户输入的IP和端口
	ip := f.EditTCPServerIP.Text()
	port := f.EditTCPServerPort.Text()

	if ip == "" || port == "" {
		f.AddLog("错误: 请输入TCP服务器IP和端口")
		return
	}

	// 验证端口号是否为数字
	if _, err := strconv.Atoi(port); err != nil {
		f.AddLog("错误: 端口号必须是数字")
		return
	}

	f.AddLog("=== TCP客户端连接开始 ===")
	f.AddLog(fmt.Sprintf("正在连接到TCP服务器: %s:%s", ip, port))

	// 使用goroutine异步执行连接操作，避免阻塞UI线程
	go func() {
		defer func() {
			if r := recover(); r != nil {
				f.AddLog(fmt.Sprintf("TCP连接异常: %v", r))
			}
		}()

		// 创建TCP连接
		address := fmt.Sprintf("%s:%s", ip, port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			f.AddLog(fmt.Sprintf("连接TCP服务器失败: %v", err))
			return
		}

		// 存储连接到结构体字段
		f.TCPConn = conn
		f.AddLog(fmt.Sprintf("成功连接到TCP服务器: %s", address))

		// 读取服务器的欢迎消息
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			f.AddLog(fmt.Sprintf("读取服务器消息失败: %v", err))
			return
		}

		serverMessage := string(buffer[:n])
		f.AddLog(fmt.Sprintf("服务器消息: %s", serverMessage))

		// 发送测试消息
		testMessage := "Hello from TCP Client"
		_, err = conn.Write([]byte(testMessage))
		if err != nil {
			f.AddLog(fmt.Sprintf("发送消息失败: %v", err))
			return
		}

		f.AddLog(fmt.Sprintf("已发送消息: %s", testMessage))

		// 读取服务器响应
		n, err = conn.Read(buffer)
		if err != nil {
			f.AddLog(fmt.Sprintf("读取服务器响应失败: %v", err))
			return
		}

		serverResponse := string(buffer[:n])
		f.AddLog(fmt.Sprintf("服务器响应: %s", serverResponse))

		f.AddLog("=== TCP客户端连接结束 ===")
	}()
}

// onSendTCPDataClick 发送TCP数据按钮事件
func (f *MainForm) onSendTCPDataClick(sender vcl.IObject) {
	// 检查是否有TCP连接
	if f.TCPConn == nil {
		f.AddLog("错误: 未连接到TCP服务器，请先点击'连接TCP'按钮")
		return
	}

	// 获取用户输入的数据
	data := f.EditTCPData.Text()
	if data == "" {
		f.AddLog("错误: 请输入要发送的数据")
		return
	}

	f.AddLog("=== 发送TCP数据开始 ===")
	f.AddLog(fmt.Sprintf("发送数据: %s", data))

	// 使用goroutine异步执行发送操作，避免阻塞UI线程
	go func() {
		defer func() {
			if r := recover(); r != nil {
				f.AddLog(fmt.Sprintf("发送TCP数据异常: %v", r))
			}
		}()

		// 发送数据
		_, err := f.TCPConn.Write([]byte(data))
		if err != nil {
			f.AddLog(fmt.Sprintf("发送数据失败: %v", err))
			// 发送失败时关闭连接
			f.TCPConn.Close()
			f.TCPConn = nil
			return
		}

		f.AddLog("数据发送成功")

		// 读取服务器响应
		buffer := make([]byte, 1024)
		n, err := f.TCPConn.Read(buffer)
		if err != nil {
			f.AddLog(fmt.Sprintf("读取服务器响应失败: %v", err))
			// 读取失败时关闭连接
			f.TCPConn.Close()
			f.TCPConn = nil
			return
		}

		serverResponse := string(buffer[:n])
		f.AddLog(fmt.Sprintf("服务器响应: %s", serverResponse))

		f.AddLog("=== 发送TCP数据结束 ===")
	}()
}

// onConcurrentClick 多线程操作按钮事件
func (f *MainForm) onConcurrentClick(sender vcl.IObject) {
	if f.ConcurrentManager == nil {
		f.AddLog("并发管理器未初始化")
		return
	}

	// 获取用户输入的参数
	threadCount := int(f.EditThreadCount.Value())
	taskCount := int(f.EditTaskCount.Value())

	f.AddLog("=== 多线程操作测试开始 ===")
	f.AddLog(fmt.Sprintf("配置参数 - 线程数量: %d, 任务数量: %d", threadCount, taskCount))

	// 1. 测试HTTP并发任务
	f.AddLog("1. 测试HTTP并发任务...")

	// 使用goroutine异步执行，避免阻塞UI线程
	go func() {
		defer func() {
			if r := recover(); r != nil {
				f.AddLog(fmt.Sprintf("HTTP并发任务异常: %v", r))
			}
		}()

		if err := f.ConcurrentManager.RunConcurrentTasks("http", taskCount); err != nil {
			f.AddLog(fmt.Sprintf("HTTP并发任务失败: %v", err))
		} else {
			f.AddLog("HTTP并发任务执行完成")
		}
	}()

	// 2. 测试计算并发任务
	f.AddLog("2. 测试计算并发任务...")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				f.AddLog(fmt.Sprintf("计算并发任务异常: %v", r))
			}
		}()

		if err := f.ConcurrentManager.RunConcurrentTasks("compute", taskCount); err != nil {
			f.AddLog(fmt.Sprintf("计算并发任务失败: %v", err))
		} else {
			f.AddLog("计算并发任务执行完成")
		}
	}()

	// 3. 测试I/O并发任务
	f.AddLog("3. 测试I/O并发任务...")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				f.AddLog(fmt.Sprintf("I/O并发任务异常: %v", r))
			}
		}()

		if err := f.ConcurrentManager.RunConcurrentTasks("io", taskCount); err != nil {
			f.AddLog(fmt.Sprintf("I/O并发任务失败: %v", err))
		} else {
			f.AddLog("I/O并发任务执行完成")
		}
	}()

	// 4. 测试混合并发任务
	f.AddLog("4. 测试混合并发任务...")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				f.AddLog(fmt.Sprintf("混合并发任务异常: %v", r))
			}
		}()

		if err := f.ConcurrentManager.RunConcurrentTasks("mixed", taskCount); err != nil {
			f.AddLog(fmt.Sprintf("混合并发任务失败: %v", err))
		} else {
			f.AddLog("混合并发任务执行完成")
		}

		// 最后更新状态和表格显示
		f.UpdateStatus("多线程操作测试完成")
		f.AddLog("=== 多线程操作测试结束 ===")
		f.AddLog("多线程测试完毕")
	}()

	// 安全地更新表格状态（使用正确的参数顺序：列在前，行在后）
	if f.TableData != nil && int(f.TableData.RowCount()) > 5 && int(f.TableData.ColCount()) > 1 {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("更新表格状态时发生异常: %v", r)
			}
		}()
		f.TableData.SetCells(1, 5, "完成")
		log.Printf("表格状态已更新: (列=1, 行=5) = '完成'")
	} else {
		log.Printf("表格尺寸不足以更新状态，表格尺寸: %dx%d",
			func() int32 { defer func() {}(); return f.TableData.RowCount() }(),
			func() int32 { defer func() {}(); return f.TableData.ColCount() }())
	}

	log.Println("多线程操作测试完成")
}

// UpdateStatus 更新状态栏显示
func (f *MainForm) UpdateStatus(status string) {
	// if f.StatusBar != nil && f.StatusBar.Panels().Count() > 0 {
	// 	f.StatusBar.Panels().Items(0).SetText(fmt.Sprintf("[%s] %s",
	// 		time.Now().Format("15:04:05"), status))
	// }
	// f.lastUpdate = time.Now()

	// log.Printf("状态更新: %s", status)
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
			colParam := int32(c) // 列索引
			rowParam := int32(r) // 行索引

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
	newText := currentText + "\n" + logText
	f.SetCurrentData(newText)

	// 自动滚动到底部，显示最新内容
	if f.EditInput != nil {
		// 计算文本总长度
		textLength := len(newText)
		if textLength > 0 {
			// 设置选择起始位置为文本末尾，实现自动滚动到底部
			f.EditInput.SetSelStart(int32(textLength))
			f.EditInput.SetSelLength(0)
		}
	}

	log.Printf("UI日志: %s", logText)
}

// applyUIStyles 应用现代化UI样式
func (f *MainForm) applyUIStyles() {
	// 设置主窗口样式
	f.applyWindowStyles()

	// 应用面板样式
	f.applyPanelStyles()

	// 应用按钮样式
	f.applyButtonStyles()

	// 应用输入组件样式
	f.applyInputStyles()

	// 应用表格样式
	f.applyTableStyles()

	// 应用状态栏样式
	f.applyStatusBarStyles()
}

// applyWindowStyles 应用窗口样式
func (f *MainForm) applyWindowStyles() {
	// 设置窗口标题和图标
	f.TForm.SetCaption("GO VCL 多功能演示程序 v1.0")

	// 设置窗口初始位置
	f.TForm.SetPosition(types.PoScreenCenter)

	// 设置窗口最大化最小化按钮可用
	f.TForm.SetBorderIcons(types.BiSystemMenu | types.BiMinimize | types.BiMaximize)
}

// applyPanelStyles 应用面板样式
func (f *MainForm) applyPanelStyles() {
	// 主面板 - 基本样式
	if f.PanelMain != nil {
		f.PanelMain.SetParentBackground(false)
	}

	// 按钮面板 - 添加边框效果
	if f.PanelButtons != nil {
		f.PanelButtons.SetParentBackground(false)
		f.PanelButtons.SetBevelOuter(types.BvLowered)
	}

	// 状态面板 - 基本样式
	if f.PanelStatus != nil {
		f.PanelStatus.SetParentBackground(false)
	}

	// 表格面板 - 基本样式
	if f.PanelTable != nil {
		f.PanelTable.SetParentBackground(false)
		f.PanelTable.SetBevelOuter(types.BvRaised)
	}
}

// applyButtonStyles 应用按钮样式
func (f *MainForm) applyButtonStyles() {
	buttons := []*vcl.TButton{
		f.BtnSaveConfig, f.BtnGetSelections, f.BtnImportExcel,
		f.BtnImportImage, f.BtnExcel, f.BtnJSON,
		f.BtnHTTP, f.BtnDatabase, f.BtnConcurrent,
		f.BtnResizeColumns, f.BtnWebServer, f.BtnTCPServer,
		f.BtnClose, f.BtnConnectTCP, f.BtnSendTCPData,
	}

	for _, button := range buttons {
		if button != nil {
			// 设置按钮字体样式
			font := button.Font()
			font.SetSize(10)
			font.SetStyle(types.TFontStyles(types.FsBold))
			// 设置支持emoji的字体，确保emoji图标能够正确显示
			font.SetName("Segoe UI Emoji")

			// 设置按钮悬停效果
			button.SetParentFont(false)
		}
	}
}

// applyInputStyles 应用输入组件样式
func (f *MainForm) applyInputStyles() {
	// 美化信息标签
	if f.LabelInfo != nil {
		font := f.LabelInfo.Font()
		font.SetSize(16)
		font.SetStyle(types.TFontStyles(types.FsBold | types.FsItalic))
		f.LabelInfo.SetFont(font)
		f.LabelInfo.SetTransparent(false)
	}

	// 美化线程数量和任务数量标签
	labels := []*vcl.TLabel{f.LabelThread, f.LabelTask}
	for _, label := range labels {
		if label != nil {
			font := label.Font()
			font.SetSize(10)
			font.SetStyle(types.TFontStyles(types.FsBold))
			label.SetFont(font)
		}
	}

	// 美化SpinEdit组件
	spins := []*vcl.TSpinEdit{f.EditThreadCount, f.EditTaskCount}
	for _, spin := range spins {
		if spin != nil {
			font := spin.Font()
			font.SetSize(10)
			spin.SetFont(font)
			spin.SetParentFont(false)
		}
	}

	// 美化输入编辑框
	if f.EditInput != nil {
		font := f.EditInput.Font()
		font.SetSize(10)
		font.SetStyle(types.TFontStyles(types.FsNormal))
		f.EditInput.SetFont(font)
		f.EditInput.SetReadOnly(true)
		// 设置滚动条样式
		f.EditInput.SetScrollBars(types.SsBoth)
	}
}

// applyTableStyles 应用表格样式
func (f *MainForm) applyTableStyles() {
	if f.TableData != nil {
		// 设置表格字体
		font := f.TableData.Font()
		font.SetSize(10)
		font.SetStyle(types.TFontStyles(types.FsNormal))
		f.TableData.SetFont(font)

		// 设置标题行样式
		f.TableData.SetFixedRows(1)
		f.TableData.SetFixedCols(0) // 修改为0，允许调整所有列

		// 设置网格线，保持列宽调整选项
		// 使用正确的Options设置方式，确保包含列宽调整选项
		gridOptions := types.TGridOptions(0)
		gridOptions = gridOptions | types.TGridOptions(types.GoFixedVertLine) |
			types.TGridOptions(types.GoFixedHorzLine) |
			types.TGridOptions(types.GoVertLine) |
			types.TGridOptions(types.GoHorzLine) |
			types.TGridOptions(types.GoRangeSelect) |
			types.TGridOptions(types.GoDrawFocusSelected) |
			types.TGridOptions(types.GoRowSizing) |
			types.TGridOptions(types.GoColSizing) | // 列大小调整 - 这是关键选项
			types.TGridOptions(types.GoRowMoving) |
			types.TGridOptions(types.GoColMoving) |
			types.TGridOptions(types.GoEditing) |
			types.TGridOptions(types.GoTabs) |
			types.TGridOptions(types.GoRowSelect) |
			types.TGridOptions(types.GoAlwaysShowEditor) |
			types.TGridOptions(types.GoThumbTracking) // 拖动跟踪 - 这是关键选项

		f.TableData.SetOptions(gridOptions)

		// 确保表格允许列宽调整
		f.TableData.SetDefaultColWidth(100)

		// 再次设置固定行列，确保设置生效
		f.TableData.SetFixedRows(1)
		f.TableData.SetFixedCols(0)

		log.Printf("表格选项已设置，包含列宽调整选项")

		// 强制刷新表格显示
		f.TableData.Invalidate()
	}
}

// applyStatusBarStyles 应用状态栏样式
func (f *MainForm) applyStatusBarStyles() {
	if f.StatusBar != nil {
		// 设置状态栏面板样式
		if f.StatusBar.Panels().Count() > 0 {
			panel := f.StatusBar.Panels().Items(0)
			panel.SetText("就绪 - 版本 1.0 ")
		}
	}
}

// autoResizeColumns 自动调整表格列宽以适应内容
func (f *MainForm) autoResizeColumns() {
	if f.TableData == nil {
		log.Printf("表格组件未初始化，无法自动调整列宽")
		return
	}

	// 检查表格选项，确保列宽调整已启用
	options := f.TableData.Options()
	isColSizingEnabled := (options & types.TGridOptions(types.GoColSizing)) != 0

	if !isColSizingEnabled {
		log.Printf("列宽调整选项未启用，重新设置表格选项...")
		// 重新设置表格选项，确保列宽调整可用
		gridOptions := types.TGridOptions(0)
		gridOptions = gridOptions | types.TGridOptions(types.GoFixedVertLine) |
			types.TGridOptions(types.GoFixedHorzLine) |
			types.TGridOptions(types.GoVertLine) |
			types.TGridOptions(types.GoHorzLine) |
			types.TGridOptions(types.GoRangeSelect) |
			types.TGridOptions(types.GoDrawFocusSelected) |
			types.TGridOptions(types.GoRowSizing) |
			types.TGridOptions(types.GoColSizing) | // 列大小调整 - 这是关键选项
			types.TGridOptions(types.GoRowMoving) |
			types.TGridOptions(types.GoColMoving) |
			types.TGridOptions(types.GoEditing) |
			types.TGridOptions(types.GoTabs) |
			types.TGridOptions(types.GoRowSelect) |
			types.TGridOptions(types.GoAlwaysShowEditor) |
			types.TGridOptions(types.GoThumbTracking) // 拖动跟踪 - 这是关键选项
		f.TableData.SetOptions(gridOptions)
		log.Printf("表格选项已重新设置，包含列宽调整选项")
	}

	// 获取表格的行数和列数
	rowCount := f.TableData.RowCount()
	colCount := f.TableData.ColCount()
	log.Printf("开始自动调整列宽，表格大小: %d行 x %d列", rowCount, colCount)

	// 确保固定列设置为0，允许调整所有列
	f.TableData.SetFixedCols(0)

	// 遍历每一列，计算最佳宽度
	for col := int32(0); col < colCount; col++ {
		maxWidth := 50 // 设置最小宽度为50像素

		// 遍历每一行，计算该列中每个单元格的宽度
		for row := int32(0); row < rowCount; row++ {
			// 获取单元格文本
			cellText := f.TableData.Cells(col, row)

			// 计算文本宽度（简单估算，每个字符约8像素宽度）
			textWidth := len(cellText)*8 + 20 // 加20像素作为边距

			// 更新最大宽度
			if textWidth > maxWidth {
				maxWidth = textWidth
			}
		}

		// 设置列宽，但不超过300像素
		if maxWidth > 300 {
			maxWidth = 300
		}

		// 设置列宽
		f.TableData.SetColWidths(col, int32(maxWidth))
		log.Printf("第%d列宽度已设置为: %d", col+1, maxWidth)
	}

	// 强制刷新表格显示
	f.TableData.Invalidate()
	log.Printf("表格列宽自动调整完成，显示已刷新")
}

// 确保MainForm实现UIInterface接口
// 这行代码确保MainForm结构体完全实现了UIInterface接口中定义的所有方法
// 如果MainForm没有实现UIInterface接口的所有方法，编译器将在此处报错
// 这是一种编译时检查机制，确保接口实现的完整性
var _ interfaces.UIInterface = (*MainForm)(nil)

// onSaveConfigClick 保存配置按钮点击事件
// 该方法处理用户点击保存配置按钮的操作，将当前单选框、多选框、线程数量和任务数量的值保存到配置文件中
func (f *MainForm) onSaveConfigClick(sender vcl.IObject) {
	f.AddLog("=== 保存配置开始 ===")

	// 获取当前配置值
	threadCount := f.EditThreadCount.Value()
	taskCount := f.EditTaskCount.Value()
	radioOption1Checked := f.RadioOption1.Checked()
	radioOption2Checked := f.RadioOption2.Checked()
	checkBox1Checked := f.CheckBox1.Checked()
	checkBox2Checked := f.CheckBox2.Checked()

	// 记录当前配置值
	f.AddLog(fmt.Sprintf("线程数量: %d", threadCount))
	f.AddLog(fmt.Sprintf("任务数量: %d", taskCount))
	f.AddLog(fmt.Sprintf("单选框1状态: %t", radioOption1Checked))
	f.AddLog(fmt.Sprintf("单选框2状态: %t", radioOption2Checked))
	f.AddLog(fmt.Sprintf("多选框1状态: %t", checkBox1Checked))
	f.AddLog(fmt.Sprintf("多选框2状态: %t", checkBox2Checked))

	// 创建配置数据结构
	configData := map[string]interface{}{
		"ThreadCount":    threadCount,
		"TaskCount":      taskCount,
		"RadioOption1":   radioOption1Checked,
		"RadioOption2":   radioOption2Checked,
		"CheckBox1":      checkBox1Checked,
		"CheckBox2":      checkBox2Checked,
		"LastUpdateTime": time.Now().Format("2006-01-02 15:04:05"),
	}

	// 使用JSON管理器保存配置到文件
	if f.JSONManager != nil {
		// 将配置数据转换为JSON字符串
		configJSON, err := f.JSONManager.CreateJSON(configData)
		if err != nil {
			f.AddLog(fmt.Sprintf("配置数据转换为JSON失败: %v", err))
			f.UpdateStatus("配置保存失败")
			f.AddLog("=== 保存配置结束 ===")
			return
		}

		// 保存JSON到配置文件
		configFileName := "app_config.json"
		if err := f.JSONManager.SaveJSONToFile(configJSON, configFileName); err != nil {
			f.AddLog(fmt.Sprintf("保存配置文件失败: %v", err))
			f.UpdateStatus("配置保存失败")
		} else {
			f.AddLog(fmt.Sprintf("配置已成功保存到文件: %s", configFileName))
			f.UpdateStatus("配置保存成功")
		}
	} else {
		f.AddLog("JSON管理器未初始化，无法保存配置")
		f.UpdateStatus("配置保存失败")
	}

	f.AddLog("=== 保存配置结束 ===")
}

// loadConfig 加载配置文件
// 该方法在程序启动时调用，从配置文件中读取保存的配置并应用到界面控件
func (f *MainForm) loadConfig() {
	f.AddLog("=== 加载配置开始 ===")

	configFileName := "app_config.json"

	// 检查配置文件是否存在
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		f.AddLog(fmt.Sprintf("配置文件不存在，使用默认配置: %s", configFileName))
		f.UpdateStatus("使用默认配置")
		f.AddLog("=== 加载配置结束 ===")
		return
	}

	// 从文件加载JSON
	if f.JSONManager == nil {
		f.AddLog("JSON管理器未初始化，无法加载配置")
		f.UpdateStatus("配置加载失败")
		f.AddLog("=== 加载配置结束 ===")
		return
	}

	jsonStr, err := f.JSONManager.LoadJSONFromFile(configFileName)
	if err != nil {
		f.AddLog(fmt.Sprintf("加载配置文件失败: %v", err))
		f.UpdateStatus("配置加载失败")
		f.AddLog("=== 加载配置结束 ===")
		return
	}

	// 解析JSON
	var configData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &configData); err != nil {
		f.AddLog(fmt.Sprintf("解析配置JSON失败: %v", err))
		f.UpdateStatus("配置解析失败")
		f.AddLog("=== 加载配置结束 ===")
		return
	}

	// 从配置中获取值并设置到控件
	if threadCount, ok := configData["ThreadCount"].(float64); ok {
		f.EditThreadCount.SetValue(int32(threadCount))
		f.AddLog(fmt.Sprintf("线程数量已设置为: %d", int32(threadCount)))
	}

	if taskCount, ok := configData["TaskCount"].(float64); ok {
		f.EditTaskCount.SetValue(int32(taskCount))
		f.AddLog(fmt.Sprintf("任务数量已设置为: %d", int32(taskCount)))
	}

	if radioOption1, ok := configData["RadioOption1"].(bool); ok {
		f.RadioOption1.SetChecked(radioOption1)
		f.AddLog(fmt.Sprintf("单选框1状态已设置为: %t", radioOption1))
	}

	if radioOption2, ok := configData["RadioOption2"].(bool); ok {
		f.RadioOption2.SetChecked(radioOption2)
		f.AddLog(fmt.Sprintf("单选框2状态已设置为: %t", radioOption2))
	}

	if checkBox1, ok := configData["CheckBox1"].(bool); ok {
		f.CheckBox1.SetChecked(checkBox1)
		f.AddLog(fmt.Sprintf("多选框1状态已设置为: %t", checkBox1))
	}

	if checkBox2, ok := configData["CheckBox2"].(bool); ok {
		f.CheckBox2.SetChecked(checkBox2)
		f.AddLog(fmt.Sprintf("多选框2状态已设置为: %t", checkBox2))
	}

	// 显示最后更新时间
	if lastUpdateTime, ok := configData["LastUpdateTime"].(string); ok {
		f.AddLog(fmt.Sprintf("配置最后更新时间: %s", lastUpdateTime))
	}

	// 更新状态栏
	f.UpdateStatus("配置加载成功")
	f.AddLog("=== 加载配置结束 ===")
}
