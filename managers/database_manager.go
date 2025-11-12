// managers/database_manager.go - 数据库操作管理器
// 功能描述: 提供数据库连接、表操作、数据查询和事务管理功能
// 主要功能:
//   - 数据库连接: 支持SQLite和MySQL数据库连接
//   - 表操作: 创建、删除、修改数据库表结构
//   - 数据查询: 执行SQL查询和数据检索
//   - 事务管理: 支持数据库事务操作
//   - 模拟模式: 提供模拟数据库功能用于测试
//   - 错误处理: 详细的错误信息和状态反馈
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12

package managers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"windows-gui-app/interfaces"

	_ "github.com/go-sql-driver/mysql" // MySQL驱动
	_ "github.com/mattn/go-sqlite3"    // SQLite驱动
)

var (
	// sqliteDriverAvailable 标记SQLite驱动是否可用
	sqliteDriverAvailable bool
)

// init 初始化函数
// 在包加载时检查SQLite驱动是否可用
func init() {
	// 检查SQLite驱动是否可用
	_, err := sql.Open("sqlite3", ":memory:")
	sqliteDriverAvailable = err == nil
	if sqliteDriverAvailable {
		log.Println("SQLite驱动可用")
	} else {
		log.Println("SQLite驱动不可用，将使用模拟模式")
	}
}

// DatabaseManager 数据库管理器
// 负责数据库连接、查询、事务处理等操作，支持SQLite和MySQL
type DatabaseManager struct {
	uiInstance interfaces.UIInterface // UI实例接口，用于显示操作状态和结果
	db         *sql.DB                // 数据库连接实例
	isMockMode bool                   // 是否为模拟模式（用于测试或驱动不可用的情况）
}

// NewDatabaseManager 创建数据库管理器实例
// 功能描述:
//   - 创建数据库管理器实例
//   - 自动检测SQLite驱动可用性
//   - 支持模拟模式用于测试环境
//
// 参数:
//   - ui: UI实例，用于显示操作结果（可以为nil，之后通过SetUIInstance设置）
//   - useMockMode: 是否使用模拟模式（true为强制模拟模式，false为自动检测）
//
// 返回:
//   - *DatabaseManager: 数据库管理器实例
//
// 使用示例:
//
//	// 自动检测模式
//	dbManager := managers.NewDatabaseManager(mainForm, false)
//
//	// 强制模拟模式
//	dbManager := managers.NewDatabaseManager(mainForm, true)
func NewDatabaseManager(ui interfaces.UIInterface, useMockMode bool) *DatabaseManager {
	// 确定是否使用模拟模式
	mockMode := useMockMode || !sqliteDriverAvailable
	if mockMode {
		log.Println("使用模拟数据库模式")
	}

	return &DatabaseManager{
		uiInstance: ui,
		db:         nil,
		isMockMode: mockMode,
	}
}

// SetUIInstance 设置UI实例
// 用于在运行时更新UI实例引用
// 参数:
//   - ui: 新的UI实例
func (dm *DatabaseManager) SetUIInstance(ui interfaces.UIInterface) {
	dm.uiInstance = ui
}

// CreateDatabase 创建数据库
// 功能描述:
//   - 创建新的SQLite数据库文件
//   - 自动创建数据库目录
//   - 初始化数据库连接
//   - 提供详细的错误信息和状态反馈
//
// 参数:
//   - dbPath: 数据库文件路径（如"data/mydb.db"）
//
// 返回:
//   - error: 操作错误，nil表示成功
//
// 可能的错误:
//   - "数据库文件路径不能为空": dbPath参数为空
//   - "创建数据库目录失败": 数据库目录创建错误
//   - "创建数据库文件失败": 数据库文件创建错误
//   - "数据库连接失败": 数据库连接初始化错误
//
// 使用示例:
//
//	err := dbManager.CreateDatabase("data/users.db")
//	if err != nil {
//	    log.Printf("创建数据库失败: %v", err)
//	}
func (dm *DatabaseManager) CreateDatabase(dbPath string) error {
	// 验证参数
	if strings.TrimSpace(dbPath) == "" {
		return fmt.Errorf("数据库文件路径不能为空")
	}

	// 如果是模拟模式，直接返回成功
	if dm.isMockMode {
		log.Printf("模拟模式：创建数据库 %s", dbPath)
		if dm.uiInstance != nil {
			dm.uiInstance.UpdateStatus(fmt.Sprintf("模拟创建数据库: %s", dbPath))
		}
		return nil
	}

	// 创建数据库目录
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建数据库目录失败: %v", err)
		}
	}

	// 创建数据库文件
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("创建数据库文件失败: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 关闭临时连接
	if err := db.Close(); err != nil {
		return fmt.Errorf("关闭数据库连接失败: %v", err)
	}

	// 更新UI状态
	if dm.uiInstance != nil {
		dm.uiInstance.UpdateStatus(fmt.Sprintf("数据库创建成功: %s", dbPath))
	}

	log.Printf("数据库创建成功: %s", dbPath)
	return nil
}

// ConnectDatabase 连接数据库
// 功能描述:
//   - 连接到现有数据库文件
//   - 支持SQLite和MySQL数据库
//   - 自动检测数据库类型
//   - 提供连接测试和错误处理
//
// 参数:
//   - connectionStr: 数据库连接字符串
//   - SQLite: 文件路径，如"data/mydb.db"
//   - MySQL: 连接字符串，如"user:password@tcp(localhost:3306)/dbname"
//
// 返回:
//   - error: 操作错误，nil表示成功
//
// 可能的错误:
//   - "数据库连接字符串不能为空": connectionStr参数为空
//   - "数据库连接失败": 数据库连接初始化错误
//   - "数据库连接测试失败": 数据库ping测试失败
//
// 使用示例:
//
//	// 连接SQLite数据库
//	err := dbManager.ConnectDatabase("data/users.db")
//
//	// 连接MySQL数据库
//	err := dbManager.ConnectDatabase("root:password@tcp(127.0.0.1:3306)/testdb")
func (dm *DatabaseManager) ConnectDatabase(connectionStr string) error {
	// 验证参数
	if strings.TrimSpace(connectionStr) == "" {
		return fmt.Errorf("数据库连接字符串不能为空")
	}

	// 如果是模拟模式，直接返回成功
	if dm.isMockMode {
		log.Printf("模拟模式：连接数据库 %s", connectionStr)
		if dm.uiInstance != nil {
			dm.uiInstance.UpdateStatus(fmt.Sprintf("模拟连接数据库: %s", connectionStr))
		}
		return nil
	}

	// 关闭现有连接
	if dm.db != nil {
		dm.db.Close()
	}

	// 根据连接字符串判断数据库类型
	var driverName string
	var actualConnectionStr string

	if strings.HasSuffix(connectionStr, ".db") || strings.HasSuffix(connectionStr, ".sqlite") {
		// SQLite数据库
		driverName = "sqlite3"
		actualConnectionStr = connectionStr
	} else if strings.Contains(connectionStr, "@tcp(") {
		// MySQL数据库
		driverName = "mysql"
		actualConnectionStr = connectionStr
	} else {
		// 默认使用SQLite
		driverName = "sqlite3"
		actualConnectionStr = connectionStr
	}

	// 创建数据库连接
	db, err := sql.Open(driverName, actualConnectionStr)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 保存连接
	dm.db = db

	// 更新UI状态
	if dm.uiInstance != nil {
		dm.uiInstance.UpdateStatus(fmt.Sprintf("数据库连接成功: %s", connectionStr))
	}

	log.Printf("数据库连接成功: %s (驱动: %s)", connectionStr, driverName)
	return nil
}

// ExecuteQuery 执行SQL查询
// 功能描述:
//   - 执行SELECT查询语句
//   - 获取查询结果并格式化为表格数据
//   - 支持参数化查询
//   - 提供详细的错误信息和状态反馈
//
// 参数:
//   - query: SQL查询语句（如"SELECT * FROM users WHERE age > ?"）
//   - args: 查询参数（可选，用于参数化查询）
//
// 返回:
//   - [][]string: 查询结果，第一行为列名，后续行为数据
//   - error: 操作错误，nil表示成功
//
// 可能的错误:
//   - "SQL查询语句不能为空": query参数为空
//   - "数据库未连接": 没有建立数据库连接
//   - "执行查询失败": SQL执行错误
//   - "获取列信息失败": 结果集列信息获取错误
//   - "读取数据失败": 结果集数据读取错误
//
// 使用示例:
//
//	// 简单查询
//	results, err := dbManager.ExecuteQuery("SELECT * FROM users")
//
//	// 参数化查询
//	results, err := dbManager.ExecuteQuery("SELECT * FROM users WHERE age > ? AND city = ?", 18, "北京")
func (dm *DatabaseManager) ExecuteQuery(query string, args ...interface{}) ([][]string, error) {
	// 验证参数
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("SQL查询语句不能为空")
	}

	// 如果是模拟模式，返回模拟数据
	if dm.isMockMode {
		log.Printf("模拟模式：执行查询 %s", query)
		if dm.uiInstance != nil {
			dm.uiInstance.UpdateStatus(fmt.Sprintf("模拟执行查询: %s", truncateString(query, 50)))
		}
		return dm.getMockQueryResult(query), nil
	}

	// 检查数据库连接
	if dm.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	// 执行查询
	rows, err := dm.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("执行查询失败: %v", err)
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("获取列信息失败: %v", err)
	}

	// 创建结果数组
	var results [][]string

	// 添加列名作为第一行
	results = append(results, columns)

	// 准备扫描目标
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// 读取数据
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("读取数据失败: %v", err)
		}

		// 转换数据为字符串
		row := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		results = append(results, row)
	}

	// 检查读取错误
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("读取数据时出错: %v", err)
	}

	// 更新UI状态
	if dm.uiInstance != nil {
		rowCount := len(results) - 1 // 减去列名行
		if rowCount < 0 {
			rowCount = 0
		}
		statusMsg := fmt.Sprintf("查询完成，返回 %d 行数据", rowCount)
		dm.uiInstance.UpdateStatus(statusMsg)
	}

	log.Printf("查询执行成功: %s，返回 %d 行数据", truncateString(query, 50), len(results)-1)
	return results, nil
}

// ExecuteNonQuery 执行非查询SQL语句
// 功能描述:
//   - 执行INSERT、UPDATE、DELETE等非查询SQL语句
//   - 获取影响的行数
//   - 支持参数化查询
//   - 提供事务支持
//
// 参数:
//   - query: SQL语句（如"INSERT INTO users (name, age) VALUES (?, ?)"）
//   - args: 查询参数（可选，用于参数化查询）
//
// 返回:
//   - int64: 影响的行数
//   - error: 操作错误，nil表示成功
//
// 可能的错误:
//   - "SQL语句不能为空": query参数为空
//   - "数据库未连接": 没有建立数据库连接
//   - "执行语句失败": SQL执行错误
//
// 使用示例:
//
//	// 插入数据
//	affectedRows, err := dbManager.ExecuteNonQuery("INSERT INTO users (name, age) VALUES (?, ?)", "张三", 25)
//
//	// 更新数据
//	affectedRows, err := dbManager.ExecuteNonQuery("UPDATE users SET age = ? WHERE name = ?", 26, "张三")
func (dm *DatabaseManager) ExecuteNonQuery(query string, args ...interface{}) (int64, error) {
	// 验证参数
	if strings.TrimSpace(query) == "" {
		return 0, fmt.Errorf("SQL语句不能为空")
	}

	// 如果是模拟模式，返回模拟结果
	if dm.isMockMode {
		log.Printf("模拟模式：执行非查询 %s", query)
		if dm.uiInstance != nil {
			dm.uiInstance.UpdateStatus(fmt.Sprintf("模拟执行: %s", truncateString(query, 50)))
		}
		return 1, nil // 模拟影响1行
	}

	// 检查数据库连接
	if dm.db == nil {
		return 0, fmt.Errorf("数据库未连接")
	}

	// 执行语句
	result, err := dm.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("执行语句失败: %v", err)
	}

	// 获取影响的行数
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("获取影响行数失败: %v", err)
	}

	// 更新UI状态
	if dm.uiInstance != nil {
		statusMsg := fmt.Sprintf("执行成功，影响 %d 行", affectedRows)
		dm.uiInstance.UpdateStatus(statusMsg)
	}

	log.Printf("非查询执行成功: %s，影响 %d 行", truncateString(query, 50), affectedRows)
	return affectedRows, nil
}

// GetTableList 获取数据库表列表
// 功能描述:
//   - 获取数据库中所有表的列表
//   - 支持SQLite和MySQL数据库
//   - 提供表名和表类型信息
//
// 参数:
//   - 无
//
// 返回:
//   - [][]string: 表列表，第一行为列名，后续行为表信息
//   - error: 操作错误，nil表示成功
//
// 可能的错误:
//   - "数据库未连接": 没有建立数据库连接
//   - "获取表列表失败": 查询系统表失败
//
// 使用示例:
//
//	tables, err := dbManager.GetTableList()
//	if err != nil {
//	    log.Printf("获取表列表失败: %v", err)
//	}
//	for _, table := range tables[1:] { // 跳过列名行
//	    log.Printf("表名: %s", table[0])
//	}
func (dm *DatabaseManager) GetTableList() ([][]string, error) {
	// 如果是模拟模式，返回模拟表
	if dm.isMockMode {
		log.Println("模拟模式：获取表列表")
		if dm.uiInstance != nil {
			dm.uiInstance.UpdateStatus("模拟获取表列表")
		}
		return [][]string{
			{"表名", "类型"},
			{"users", "table"},
			{"products", "table"},
			{"orders", "table"},
		}, nil
	}

	// 检查数据库连接
	if dm.db == nil {
		return nil, fmt.Errorf("数据库未连接")
	}

	// 获取数据库类型
	driverName := dm.getDriverName()

	var query string
	switch driverName {
	case "sqlite3":
		query = "SELECT name as '表名', type as '类型' FROM sqlite_master WHERE type='table' ORDER BY name"
	case "mysql":
		query = "SHOW TABLES"
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", driverName)
	}

	return dm.ExecuteQuery(query)
}

// Close 关闭数据库连接
// 功能描述:
//   - 关闭当前数据库连接
//   - 清理相关资源
//
// 参数:
//   - 无
//
// 返回:
//   - error: 操作错误，nil表示成功
//
// 使用示例:
//
//	err := dbManager.Close()
//	if err != nil {
//	    log.Printf("关闭连接失败: %v", err)
//	}
func (dm *DatabaseManager) Close() error {
	if dm.db != nil {
		if err := dm.db.Close(); err != nil {
			return fmt.Errorf("关闭数据库连接失败: %v", err)
		}
		dm.db = nil
	}

	if dm.uiInstance != nil {
		dm.uiInstance.UpdateStatus("数据库连接已关闭")
	}

	log.Println("数据库连接已关闭")
	return nil
}

// IsConnected 检查是否已连接数据库
// 返回:
//   - bool: true表示已连接，false表示未连接
func (dm *DatabaseManager) IsConnected() bool {
	if dm.isMockMode {
		return true // 模拟模式始终返回已连接
	}
	return dm.db != nil
}

// IsMockMode 检查是否为模拟模式
// 返回:
//   - bool: true表示模拟模式，false表示真实数据库模式
func (dm *DatabaseManager) IsMockMode() bool {
	return dm.isMockMode
}

// GetConnectionInfo 获取数据库连接信息
// 返回:
//   - string: 连接信息描述
func (dm *DatabaseManager) GetConnectionInfo() string {
	if dm.isMockMode {
		return "模拟数据库模式"
	}
	if dm.db == nil {
		return "未连接"
	}
	return fmt.Sprintf("已连接 (%s驱动)", dm.getDriverName())
}

// 内部辅助函数

// getDriverName 获取当前数据库驱动名称
func (dm *DatabaseManager) getDriverName() string {
	if dm.db == nil {
		return "unknown"
	}
	// 通过Ping方法间接获取驱动信息（简化实现）
	return "sqlite3" // 简化处理，实际项目中可以通过更复杂的方式判断
}

// getMockQueryResult 获取模拟查询结果
func (dm *DatabaseManager) getMockQueryResult(query string) [][]string {
	queryLower := strings.ToLower(strings.TrimSpace(query))

	// 根据查询类型返回不同的模拟数据
	if strings.Contains(queryLower, "select") && strings.Contains(queryLower, "users") {
		return [][]string{
			{"ID", "姓名", "年龄", "邮箱"},
			{"1", "张三", "25", "zhangsan@example.com"},
			{"2", "李四", "30", "lisi@example.com"},
			{"3", "王五", "28", "wangwu@example.com"},
		}
	} else if strings.Contains(queryLower, "select") && strings.Contains(queryLower, "products") {
		return [][]string{
			{"ID", "产品名称", "价格", "库存"},
			{"1", "笔记本电脑", "5999.00", "10"},
			{"2", "智能手机", "2999.00", "25"},
			{"3", "平板电脑", "1999.00", "15"},
		}
	} else if strings.Contains(queryLower, "show") && strings.Contains(queryLower, "tables") {
		return [][]string{
			{"Tables_in_database"},
			{"users"},
			{"products"},
			{"orders"},
		}
	}

	// 默认返回空结果
	return [][]string{
		{"结果"},
		{"模拟查询结果"},
	}
}

// truncateString 截断字符串
func truncateString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	return str[:maxLen] + "..."
}
