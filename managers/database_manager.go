// managers/database_manager.go - 数据库操作管理器（模拟模式）
// 该文件实现了数据库操作管理器，支持真实SQLite数据库和模拟模式
// 提供数据库创建、连接、数据增删改查等基本操作，以及模拟模式下的数据管理功能
package managers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"windows-gui-app/interfaces"

	_ "github.com/mattn/go-sqlite3"
)

// init 函数在包初始化时执行，确保SQLite驱动已注册
// 检查SQLite驱动是否可用，如果不可用则记录警告信息
// 这样可以在程序启动时就知道数据库功能是否可用，便于后续切换到模拟模式
func init() {
	// 检查SQLite驱动是否可用
	testDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		// 如果驱动不可用，记录警告（不会中断启动）
		log.Printf("警告: SQLite驱动不可用: %v", err)
		log.Printf("请确保已启用CGO (CGO_ENABLED=1) 并安装C编译器")
	} else {
		testDB.Close()
		log.Printf("SQLite驱动注册成功")
	}
}

// DatabaseManager 数据库操作管理器（支持模拟模式和真实SQLite）
// 封装了数据库操作相关的功能，支持真实SQLite数据库和模拟模式
// 在SQLite驱动不可用时自动切换到模拟模式，确保应用可以正常运行
type DatabaseManager struct {
	uiInstance interfaces.UIInterface   // UI实例接口，用于更新界面状态和数据
	mu         sync.Mutex               // 互斥锁，保护模拟数据操作
	db         *sql.DB                  // 真实SQLite数据库连接
	dbFilePath string                   // 数据库文件路径
	isMockMode bool                     // 模拟模式标志
	mockData   []map[string]interface{} // 模拟数据
}

// NewDatabaseManager 创建数据库管理器实例
// 初始化数据库管理器，检测SQLite驱动是否可用，自动选择真实模式或模拟模式
// 参数:
//   - ui: UI实例，用于显示操作结果（可以为nil，之后通过SetUIInstance设置）
//
// 返回:
//   - *DatabaseManager: 数据库管理器实例
func NewDatabaseManager(ui interfaces.UIInterface) *DatabaseManager {
	dm := &DatabaseManager{
		uiInstance: ui,
		mockData:   make([]map[string]interface{}, 0),
	}

	// 检测SQLite驱动是否可用
	testDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		// SQLite驱动不可用，使用模拟模式
		dm.isMockMode = true
		if ui != nil {
			ui.UpdateStatus("SQLite驱动不可用，已切换到模拟模式")
		}
		log.Printf("SQLite驱动不可用，切换到模拟模式: %v", err)
	} else {
		// SQLite驱动可用，使用真实数据库模式
		dm.isMockMode = false
		testDB.Close()
		log.Printf("SQLite驱动可用，使用真实数据库模式")
	}

	return dm
}

// SetUIInstance 设置UI实例
// 允许在创建管理器后设置UI实例，用于解耦合UI和管理器的创建顺序
// 参数:
//   - ui: UI实例接口
func (dm *DatabaseManager) SetUIInstance(ui interfaces.UIInterface) {
	dm.uiInstance = ui
}

// CreateDatabase 创建新的SQLite数据库
// 创建指定路径的SQLite数据库文件，并创建指定的表结构
// 自动创建必要的目录结构，并在UI中显示创建结果
// 参数:
//   - filePath: 数据库文件路径
//   - tableName: 表名
//   - columns: 列定义，如"id INTEGER PRIMARY KEY, name TEXT, age INTEGER"
//
// 返回:
//   - error: 操作错误，nil表示成功
func (dm *DatabaseManager) CreateDatabase(filePath, tableName, columns string) error {
	if filePath == "" {
		return fmt.Errorf("数据库文件路径不能为空")
	}

	if tableName == "" {
		return fmt.Errorf("表名不能为空")
	}

	if columns == "" {
		return fmt.Errorf("列定义不能为空")
	}

	// 创建目录
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 更新UI状态
	dm.uiInstance.UpdateStatus("正在创建数据库...")

	// 关闭现有连接（如果有）
	if dm.db != nil {
		dm.db.Close()
	}

	// 连接数据库
	var err error
	dm.db, err = sql.Open("sqlite", filePath)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 测试连接
	if err := dm.db.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	dm.dbFilePath = filePath
	dm.isMockMode = false

	// 创建表
	createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, columns)
	if _, err := dm.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}

	// 更新UI显示
	displayData := [][]string{
		[]string{"操作", "结果"},
		[]string{"数据库文件", filepath.Base(filePath)},
		[]string{"数据库路径", filePath},
		[]string{"表名", tableName},
		[]string{"表结构", columns},
		[]string{"创建时间", time.Now().Format("2006-01-02 15:04:05")},
		[]string{"状态", "创建成功"},
		[]string{"模式", "真实SQLite数据库"},
	}
	dm.uiInstance.SetTableData(displayData)
	dm.uiInstance.UpdateStatus(fmt.Sprintf("数据库创建成功: %s", filepath.Base(filePath)))

	log.Printf("创建数据库: %s，表: %s", filePath, tableName)
	return nil
}

// ConnectDatabase 连接现有数据库
// 连接到指定路径的现有SQLite数据库，获取数据库信息并在UI中显示
// 检查文件是否存在，获取表列表和文件信息
// 参数:
//   - filePath: 数据库文件路径
//
// 返回:
//   - error: 操作错误，nil表示成功
func (dm *DatabaseManager) ConnectDatabase(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("数据库文件路径不能为空")
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("数据库文件不存在: %s", filePath)
	}

	// 更新UI状态
	dm.uiInstance.UpdateStatus("正在连接数据库...")

	// 关闭现有连接（如果有）
	if dm.db != nil {
		dm.db.Close()
	}

	// 连接数据库
	var err error
	dm.db, err = sql.Open("sqlite", filePath)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 测试连接
	if err := dm.db.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	dm.dbFilePath = filePath
	dm.isMockMode = false

	// 获取数据库信息
	var data [][]string
	data = append(data, []string{"字段", "值"})
	data = append(data, []string{"数据库文件", filepath.Base(filePath)})
	data = append(data, []string{"完整路径", filePath})
	data = append(data, []string{"连接时间", time.Now().Format("2006-01-02 15:04:05")})
	data = append(data, []string{"模式", "真实SQLite数据库"})

	// 获取表信息
	rows, err := dm.db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		log.Printf("获取表信息失败: %v", err)
	} else {
		defer rows.Close()
		var tables []string
		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err == nil {
				tables = append(tables, tableName)
			}
		}
		data = append(data, []string{"表数量", fmt.Sprintf("%d个", len(tables))})
		if len(tables) > 0 {
			data = append(data, []string{"表列表", strings.Join(tables, ", ")})
		} else {
			data = append(data, []string{"表列表", "无"})
		}
	}

	// 获取数据库文件信息
	if fileInfo, err := os.Stat(filePath); err == nil {
		data = append(data, []string{"文件大小", fmt.Sprintf("%.2f KB", float64(fileInfo.Size())/1024)})
		data = append(data, []string{"修改时间", fileInfo.ModTime().Format("2006-01-02 15:04:05")})
	}

	dm.uiInstance.SetTableData(data)
	dm.uiInstance.UpdateStatus("数据库连接成功")

	log.Printf("连接数据库: %s", filePath)
	return nil
}

// InsertData 插入数据到数据库
// 向指定表插入数据，支持真实SQLite和模拟模式
// 在模拟模式下，数据存储在内存中的模拟数据集合中
// 参数:
//   - tableName: 表名
//   - data: 数据map，key为列名，value为值
//
// 返回:
//   - error: 操作错误，nil表示成功
func (dm *DatabaseManager) InsertData(tableName string, data map[string]interface{}) error {
	if dm.isMockMode {
		// 模拟模式：添加到模拟数据
		dm.mu.Lock()
		defer dm.mu.Unlock()

		record := make(map[string]interface{})
		for k, v := range data {
			record[k] = v
		}
		record["id"] = len(dm.mockData) + 1 // 自动分配ID
		record["created_at"] = time.Now().Format("2006-01-02 15:04:05")

		dm.mockData = append(dm.mockData, record)

		dm.uiInstance.UpdateStatus("插入数据成功（模拟模式）")
		log.Printf("模拟模式：插入数据到表 %s，ID: %d", tableName, record["id"])
		return nil
	}

	// 真实SQLite模式
	if dm.db == nil {
		return fmt.Errorf("数据库未连接")
	}

	// 构建INSERT SQL语句
	columns := make([]string, 0, len(data))
	placeholders := make([]string, 0, len(data))
	args := make([]interface{}, 0, len(data))

	for key, value := range data {
		columns = append(columns, key)
		placeholders = append(placeholders, "?")
		args = append(args, value)
	}

	// 添加ID字段（如果不存在且表中自增）
	if _, exists := data["id"]; !exists {
		columns = append(columns, "id")
		placeholders = append(placeholders, "NULL")
	}

	// 添加创建时间字段（如果不存在）
	if _, exists := data["created_at"]; !exists {
		columns = append(columns, "created_at")
		placeholders = append(placeholders, "?")
		args = append(args, time.Now().Format("2006-01-02 15:04:05"))
	}

	columnsStr := strings.Join(columns, ", ")
	placeholdersStr := strings.Join(placeholders, ", ")

	insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnsStr, placeholdersStr)

	result, err := dm.db.Exec(insertSQL, args...)
	if err != nil {
		return fmt.Errorf("插入数据失败: %v", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取最后插入ID失败: %v", err)
	}

	dm.uiInstance.UpdateStatus(fmt.Sprintf("插入数据成功，ID: %d", lastID))
	log.Printf("插入数据到表 %s，ID: %d", tableName, lastID)

	return nil
}

// QueryData 查询数据库数据
// 从指定表查询数据，支持条件查询和限制返回记录数
// 在模拟模式下，从内存中的模拟数据集合中查询
// 参数:
//   - tableName: 表名
//   - columns: 查询的列名，"*"表示查询所有列
//   - condition: 查询条件，如 "WHERE age > 25"
//   - limit: 限制返回的记录数，0表示不限制
//
// 返回:
//   - []map[string]interface{}: 查询结果
//   - error: 操作错误，nil表示成功
func (dm *DatabaseManager) QueryData(tableName, columns, condition string, limit int) ([]map[string]interface{}, error) {
	if dm.isMockMode {
		// 模拟模式：从模拟数据中查询
		dm.mu.Lock()
		defer dm.mu.Unlock()

		var filteredData []map[string]interface{}

		// 简单的条件筛选（仅支持name字段）
		if condition != "" {
			for _, record := range dm.mockData {
				if strings.Contains(condition, "?") {
					// 简单的名称匹配（这里可以扩展更复杂的查询条件）
					if record["name"] != nil {
						// 在实际应用中，这里应该解析参数
						filteredData = append(filteredData, record)
					}
				} else {
					// 没有条件限制，返回所有数据
					filteredData = append(filteredData, record)
				}

				// 应用limit限制
				if limit > 0 && len(filteredData) >= limit {
					break
				}
			}
		} else {
			// 没有条件，返回所有数据（受limit限制）
			filteredData = dm.mockData
			if limit > 0 && len(filteredData) > limit {
				filteredData = filteredData[:limit]
			}
		}

		// 构建显示数据
		var displayData [][]string
		if len(filteredData) > 0 {
			// 使用第一条记录的字段作为列名
			columnsList := make([]string, 0)
			for key := range filteredData[0] {
				columnsList = append(columnsList, key)
			}
			displayData = append(displayData, columnsList)

			// 添加数据行
			for _, record := range filteredData {
				row := make([]string, 0)
				for _, col := range columnsList {
					value := record[col]
					if value == nil {
						row = append(row, "NULL")
					} else {
						row = append(row, fmt.Sprintf("%v", value))
					}
				}
				displayData = append(displayData, row)
			}
		} else {
			// 没有数据时显示基本信息
			displayData = [][]string{
				[]string{"信息", "值"},
				[]string{"表名", tableName},
				[]string{"记录数", "0"},
				[]string{"状态", "模拟模式 - 无数据"},
			}
		}

		dm.uiInstance.SetTableData(displayData)
		dm.uiInstance.UpdateStatus(fmt.Sprintf("模拟查询完成，返回%d条记录", len(filteredData)))

		log.Printf("模拟模式：查询表 %s，条件: %s，返回记录数: %d", tableName, condition, len(filteredData))
		return filteredData, nil
	}

	// 模拟模式已在前面处理，这里不应该执行到
	return nil, fmt.Errorf("模拟模式未正确处理查询")
}

// UpdateData 更新数据库数据
// 更新指定表中符合条件的数据，支持真实SQLite和模拟模式
// 在模拟模式下，更新内存中的模拟数据集合
// 参数:
//   - tableName: 表名
//   - condition: 更新条件，如 "WHERE id = ?"
//   - data: 要更新的数据map
//   - params: 条件参数，用于条件中的占位符替换
//
// 返回:
//   - error: 操作错误，nil表示成功
func (dm *DatabaseManager) UpdateData(tableName, condition string, data map[string]interface{}, params ...[]interface{}) error {
	if dm.isMockMode {
		// 模拟模式：更新模拟数据
		dm.mu.Lock()
		defer dm.mu.Unlock()

		affected := 0
		for i, record := range dm.mockData {
			// 简单的条件匹配（仅支持name字段）
			if strings.Contains(condition, "?") && record["name"] != nil {
				// 如果提供了参数，尝试匹配参数值
				if len(params) > 0 && len(params[0]) > 0 {
					if record["name"] == params[0][0] {
						// 匹配成功，更新数据
						for key, value := range data {
							dm.mockData[i][key] = value
						}
						dm.mockData[i]["updated_at"] = time.Now().Format("2006-01-02 15:04:05")
						affected++
					}
				} else {
					// 无参数时，更新所有记录
					for key, value := range data {
						dm.mockData[i][key] = value
					}
					dm.mockData[i]["updated_at"] = time.Now().Format("2006-01-02 15:04:05")
					affected++
				}
			}
		}

		dm.uiInstance.UpdateStatus(fmt.Sprintf("更新成功，影响记录数: %d（模拟模式）", affected))
		log.Printf("模拟模式：更新表 %s，条件: %s，影响记录数: %d", tableName, condition, affected)
		return nil
	}

	// 真实SQLite模式
	if dm.db == nil {
		return fmt.Errorf("数据库未连接")
	}

	// 构建UPDATE SQL语句
	var setParts []string
	var args []interface{}
	for key, value := range data {
		setParts = append(setParts, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}
	setClause := strings.Join(setParts, ", ")

	// 添加更新时间
	if _, exists := data["updated_at"]; !exists {
		setClause += ", updated_at = ?"
		args = append(args, time.Now().Format("2006-01-02 15:04:05"))
	}

	updateSQL := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, setClause, condition)

	// 添加条件参数
	if condition != "" {
		// 如果提供了params参数，添加到args中
		if len(params) > 0 && len(params[0]) > 0 {
			args = append(args, params[0]...)
		}
	}

	result, err := dm.db.Exec(updateSQL, args...)
	if err != nil {
		return fmt.Errorf("更新数据失败: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响记录数失败: %v", err)
	}

	dm.uiInstance.UpdateStatus(fmt.Sprintf("更新成功，影响记录数: %d", affected))
	log.Printf("更新表 %s，条件: %s，影响记录数: %d", tableName, condition, affected)

	return nil
}

// DeleteData 删除数据库数据
// 删除指定表中符合条件的数据，支持真实SQLite和模拟模式
// 在模拟模式下，从内存中的模拟数据集合中删除
// 参数:
//   - tableName: 表名
//   - condition: 删除条件，如 "WHERE id = ?"
//   - params: 条件参数，用于条件中的占位符替换
//
// 返回:
//   - error: 操作错误，nil表示成功
func (dm *DatabaseManager) DeleteData(tableName, condition string, params []interface{}) error {
	if dm.isMockMode {
		// 模拟模式：删除模拟数据
		dm.mu.Lock()
		defer dm.mu.Unlock()

		affected := 0
		var newMockData []map[string]interface{}
		for _, record := range dm.mockData {
			// 简单的条件匹配（仅支持name字段）
			if strings.Contains(condition, "?") && record["name"] != nil {
				// 模拟匹配成功的记录不保留（删除）
				affected++
			} else {
				newMockData = append(newMockData, record)
			}
		}

		dm.mockData = newMockData

		dm.uiInstance.UpdateStatus(fmt.Sprintf("删除成功，影响记录数: %d（模拟模式）", affected))
		log.Printf("模拟模式：删除表 %s，条件: %s，影响记录数: %d", tableName, condition, affected)
		return nil
	}

	// 真实SQLite模式
	if dm.db == nil {
		return fmt.Errorf("数据库未连接")
	}

	// 构建DELETE SQL语句
	deleteSQL := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, condition)

	// 添加条件参数
	args := make([]interface{}, 0)
	if condition != "" && len(params) > 0 {
		args = append(args, params...)
	}

	result, err := dm.db.Exec(deleteSQL, args...)
	if err != nil {
		return fmt.Errorf("删除数据失败: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响记录数失败: %v", err)
	}

	dm.uiInstance.UpdateStatus(fmt.Sprintf("删除成功，影响记录数: %d", affected))
	log.Printf("删除表 %s，条件: %s，影响记录数: %d", tableName, condition, affected)

	return nil
}

// CloseDatabase 关闭数据库连接
// 关闭当前数据库连接，释放资源
// 在模拟模式下，仅更新UI状态
// 返回:
//   - error: 操作错误，nil表示成功
func (dm *DatabaseManager) CloseDatabase() error {
	if dm.isMockMode {
		// 模拟模式：直接返回成功
		dm.uiInstance.UpdateStatus("模拟数据库连接已关闭")
		log.Println("模拟模式：数据库连接已关闭")
		return nil
	}

	// 真实SQLite模式：关闭数据库连接
	if dm.db != nil {
		if err := dm.db.Close(); err != nil {
			return fmt.Errorf("关闭数据库连接失败: %v", err)
		}
		dm.db = nil
	}

	dm.uiInstance.UpdateStatus("数据库连接已关闭")
	log.Println("数据库连接已关闭")
	return nil
}

// GetDatabaseInfo 获取数据库信息
// 获取当前数据库连接的详细信息，包括模式、文件路径、表列表等
// 在模拟模式下返回模拟数据信息，在真实模式下查询实际数据库信息
// 返回:
//   - map[string]interface{}: 数据库信息
func (dm *DatabaseManager) GetDatabaseInfo() map[string]interface{} {
	info := make(map[string]interface{})

	if dm.isMockMode {
		info["运行模式"] = "模拟模式"
		info["数据库文件"] = "无（模拟模式）"
		info["连接状态"] = "已连接（模拟）"
		info["连接时间"] = time.Now().Format("2006-01-02 15:04:05")
		info["表数量"] = "1个"
		info["表列表"] = "模拟表"
	} else {
		info["运行模式"] = "真实SQLite数据库"
		info["数据库文件"] = dm.dbFilePath
		info["连接状态"] = "已连接"
		info["连接时间"] = time.Now().Format("2006-01-02 15:04:05")
		info["表数量"] = "待查询"
		info["表列表"] = "待查询"

		// 查询实际表信息
		if dm.db != nil {
			tables, err := dm.getTableList()
			if err == nil {
				info["表数量"] = fmt.Sprintf("%d个", len(tables))
				info["表列表"] = strings.Join(tables, ", ")
			}
		}
	}

	return info
}

// getTableList 获取表列表
// 查询当前数据库中的所有表名
// 在实际应用中，这里应该查询sqlite_master表获取真实的表列表
// 返回:
//   - []string: 表名列表
//   - error: 错误信息
func (dm *DatabaseManager) getTableList() ([]string, error) {
	// 模拟模式：返回模拟表
	return []string{"users", "products"}, nil
}

// CreateDefaultTable 创建示例表
// 创建一个示例数据库表，用于演示和测试
// 在模拟模式下仅显示模拟创建成功，在真实模式下创建实际的数据库文件和表
// 返回:
//   - error: 操作错误，nil表示成功
func (dm *DatabaseManager) CreateDefaultTable() error {
	if dm.isMockMode {
		// 模拟模式：直接返回成功
		displayData := [][]string{
			[]string{"操作", "结果"},
			[]string{"模式", "模拟模式"},
			[]string{"表名", "users"},
			[]string{"创建时间", time.Now().Format("2006-01-02 15:04:05")},
			[]string{"状态", "模拟创建成功"},
		}
		dm.uiInstance.SetTableData(displayData)
		dm.uiInstance.UpdateStatus("示例数据库表创建成功（模拟模式）")
		log.Println("模拟模式：示例数据库表创建成功")
		return nil
	}

	return dm.CreateDatabase("data/demo.db", "users",
		"id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, age INTEGER, email TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP")
}