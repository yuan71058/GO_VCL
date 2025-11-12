# GO_VCL Windows GUI应用程序 - API文档

## 📚 概述

本文档详细描述了GO_VCL Windows GUI应用程序中各个功能模块的API接口，包括Excel操作、JSON处理、HTTP请求和数据库操作等核心功能。

## 🏗️ 项目架构

### 核心模块结构
```
windows-gui-app/
├── managers/          # 功能管理器
│   ├── excel_manager.go    # Excel文件操作
│   ├── json_manager.go     # JSON数据处理
│   ├── http_manager.go     # HTTP网络请求
│   └── database_manager.go # 数据库操作
├── handlers/          # 事件处理器
│   └── button_handlers.go  # 按钮事件处理
├── ui/               # 用户界面
│   └── MainForm.go    # 主窗口界面
├── interfaces/       # 接口定义
│   └── ui.go          # UI接口规范
└── utils/            # 工具函数
    └── helpers.go     # 辅助函数
```

## 📊 Excel管理器 (ExcelManager)

### 概述
Excel管理器负责处理Excel文件的导入导出操作，支持.xlsx格式，提供完整的数据读写功能。

### 构造函数
```go
func NewExcelManager(ui interfaces.UIInterface) *ExcelManager
```
**参数:**
- `ui`: UI实例接口，用于显示操作状态和结果

**返回:**
- `*ExcelManager`: Excel管理器实例

**示例:**
```go
excelManager := managers.NewExcelManager(uiInstance)
```

### 主要方法

#### 1. 导入Excel文件
```go
func (em *ExcelManager) ImportExcel(filePath string) error
```
**功能:** 从指定路径导入Excel文件并显示在UI表格中

**参数:**
- `filePath`: Excel文件的完整路径

**返回:**
- `error`: 操作错误，nil表示成功

**错误处理:**
- 文件路径为空：返回"Excel文件路径不能为空"
- 文件不存在：返回"Excel文件不存在: [文件路径]"
- 文件格式错误：返回"打开Excel文件失败: [错误信息]"
- 工作表无数据：返回"工作表中没有数据"

**使用示例:**
```go
err := excelManager.ImportExcel("C:\\data\\sample.xlsx")
if err != nil {
    log.Printf("导入失败: %v", err)
}
```

#### 2. 导出Excel文件
```go
func (em *ExcelManager) ExportExcel(filePath string, data [][]string) error
```
**功能:** 将数据导出到Excel文件

**参数:**
- `filePath`: 输出Excel文件路径（自动添加.xlsx扩展名）
- `data`: 要导出的二维字符串数组，nil表示从UI表格获取

**返回:**
- `error`: 操作错误，nil表示成功

**特殊功能:**
- 自动添加.xlsx扩展名
- 自动创建输出目录
- 支持1048576行×16384列（Excel限制）
- 空值自动转换为空字符串

**使用示例:**
```go
data := [][]string{
    {"姓名", "年龄", "城市"},
    {"张三", "25", "北京"},
    {"李四", "30", "上海"},
}
err := excelManager.ExportExcel("C:\\output\\report.xlsx", data)
```

### 数据格式规范

**导入数据格式:**
```
工作表1:
┌──────┬─────┬──────┐
│ 姓名 │ 年龄 │ 城市 │
├──────┼─────┼──────┤
│ 张三 │ 25  │ 北京 │
│ 李四 │ 30  │ 上海 │
└──────┴─────┴──────┘
```

**导出数据格式:**
```go
[][]string{
    {"列1", "列2", "列3"},     // 表头
    {"数据1", "数据2", "数据3"}, // 数据行
    {"数据4", "数据5", "数据6"}, // 数据行
}
```

## 📝 JSON管理器 (JSONManager)

### 概述
JSON管理器提供JSON数据的解析、生成、修改和文件操作功能，支持复杂嵌套结构。

### 构造函数
```go
func NewJSONManager(ui interfaces.UIInterface) *JSONManager
```

### 主要方法

#### 1. 解析JSON字符串
```go
func (jm *JSONManager) ParseJSON(jsonData string) error
```
**功能:** 解析JSON字符串并在表格中显示键值对

**参数:**
- `jsonData`: JSON格式字符串

**返回:**
- `error`: 解析错误，nil表示成功

**支持格式:**
```json
{
    "name": "张三",
    "age": 25,
    "address": {
        "city": "北京",
        "street": "长安街"
    }
}
```

#### 2. 根据路径获取值
```go
func (jm *JSONManager) GetValueByPath(jsonData, path string) (string, error)
```
**功能:** 使用点号路径获取JSON中的值

**参数:**
- `jsonData`: JSON字符串
- `path`: 路径，如 "user.name" 或 "items.0.name"

**返回:**
- `string`: 获取到的值
- `error`: 操作错误

**示例:**
```go
value, err := jsonManager.GetValueByPath(jsonStr, "user.profile.age")
```

#### 3. 根据路径设置值
```go
func (jm *JSONManager) SetValueByPath(jsonData, path, value string) (string, error)
```
**功能:** 修改JSON中指定路径的值

**返回:**
- `string`: 更新后的完整JSON字符串

#### 4. 创建JSON数据
```go
func (jm *JSONManager) CreateJSON(data interface{}) (string, error)
```
**功能:** 将Go数据结构转换为JSON字符串

**支持数据类型:**
- map[string]interface{}
- []interface{}
- 基本类型（string, int, float64, bool）

#### 5. 文件操作
```go
func (jm *JSONManager) SaveJSONToFile(jsonData, filePath string) error
func (jm *JSONManager) LoadJSONFromFile(filePath string) (string, error)
```

## 🌐 HTTP管理器 (HTTPManager)

### 概述
HTTP管理器提供完整的HTTP请求功能，支持GET、POST等方法，自动处理JSON响应。

### 构造函数
```go
func NewHTTPManager(ui interfaces.UIInterface) *HTTPManager
```

### 主要方法

#### 1. GET请求
```go
func (hm *HTTPManager) GETRequest(url string) (*http.Response, error)
```

#### 2. POST请求
```go
func (hm *HTTPManager) POSTRequest(urlStr string, data interface{}, contentType string) error
```
**支持数据类型:**
- `string`: 纯文本或JSON字符串
- `map[string]string`: 表单数据
- `map[string]interface{}`: JSON数据（自动序列化）
- `[]byte`: 二进制数据

**自动内容类型检测:**
- JSON数据: `application/json`
- 表单数据: `application/x-www-form-urlencoded`
- 自定义: 通过contentType参数指定

#### 3. 文件下载
```go
func (hm *HTTPManager) DownloadFile(urlStr, filePath string) error
```

### 响应处理
**自动检测响应类型:**
- JSON响应: 解析并显示结构化数据
- HTML响应: 显示页面标题和主要内容
- 纯文本: 直接显示内容

**状态信息显示:**
- HTTP状态码
- 响应时间
- 数据大小
- 内容类型

## 🗄️ 数据库管理器 (DatabaseManager)

### 概述
数据库管理器支持SQLite数据库操作，提供双模式运行：真实数据库模式和模拟模式。

### 运行模式

#### 1. 真实数据库模式
- 使用SQLite驱动进行实际数据库操作
- 支持完整的SQL语法
- 数据持久化存储

#### 2. 模拟模式
- 当SQLite驱动不可用时自动切换
- 使用内存数据结构模拟数据库操作
- 提供相同的功能接口

### 构造函数
```go
func NewDatabaseManager(ui interfaces.UIInterface) *DatabaseManager
```

### 主要方法

#### 1. 创建数据库
```go
func (dm *DatabaseManager) CreateDatabase(filePath, tableName, columns string) error
```
**参数说明:**
- `filePath`: 数据库文件路径
- `tableName`: 表名
- `columns`: 列定义，如 "id INTEGER PRIMARY KEY, name TEXT, age INTEGER"

#### 2. 连接数据库
```go
func (dm *DatabaseManager) ConnectDatabase(filePath string) error
```

#### 3. 数据操作
```go
func (dm *DatabaseManager) InsertData(tableName string, data map[string]interface{}) error
func (dm *DatabaseManager) UpdateData(tableName, condition string, data map[string]interface{}) error
func (dm *DatabaseManager) DeleteData(tableName, condition string) error
func (dm *DatabaseManager) QueryData(tableName, condition string) ([]map[string]interface{}, error)
```

### 数据库操作示例

**创建表:**
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    age INTEGER,
    email TEXT UNIQUE
)
```

**插入数据:**
```go
data := map[string]interface{}{
    "name": "张三",
    "age": 25,
    "email": "zhangsan@example.com",
}
err := dbManager.InsertData("users", data)
```

**查询数据:**
```go
results, err := dbManager.QueryData("users", "age > 20")
```

## 🔘 按钮事件处理器 (ButtonHandlers)

### 概述
按钮事件处理器统一管理所有UI按钮的点击事件，协调各个功能模块。

### 构造函数
```go
func NewButtonHandlers(ui *ui.MainForm) *ButtonHandlers
```

### 事件处理方法

#### Excel相关事件
```go
func (bh *ButtonHandlers) OnExcelImport(filePath, sheetName string) error
func (bh *ButtonHandlers) OnExcelExport(filePath, sheetName string, data [][]string) error
```

#### JSON相关事件
```go
func (bh *ButtonHandlers) OnJSONParse(jsonText string) error
func (bh *ButtonHandlers) OnJSONSave(filePath string, data map[string]interface{}) error
func (bh *ButtonHandlers) OnJSONLoad(filePath string) error
```

#### HTTP相关事件
```go
func (bh *ButtonHandlers) OnHTTPRequest(method, url string, headers map[string]string, body string) error
func (bh *ButtonHandlers) OnFileDownload(url, filePath string) error
```

#### 数据库相关事件
```go
func (bh *ButtonHandlers) OnDatabaseCreate(dbPath, tableName, columns string) error
func (bh *ButtonHandlers) OnDatabaseConnect(dbPath string) error
func (bh *ButtonHandlers) OnDatabaseInsert(tableName string, data map[string]interface{}) error
func (bh *ButtonHandlers) OnDatabaseQuery(tableName, condition string) error
```

## 🎯 使用示例

### 完整工作流程示例

```go
package main

import (
    "windows-gui-app/managers"
    "windows-gui-app/ui"
)

func main() {
    // 创建UI实例
    mainForm := ui.NewMainForm()
    
    // 创建功能管理器
    excelManager := managers.NewExcelManager(mainForm)
    jsonManager := managers.NewJSONManager(mainForm)
    httpManager := managers.NewHTTPManager(mainForm)
    dbManager := managers.NewDatabaseManager(mainForm)
    
    // Excel操作示例
    err := excelManager.ImportExcel("data.xlsx")
    if err != nil {
        // 处理错误
    }
    
    // JSON操作示例
    jsonData := `{"name": "测试", "value": 123}`
    err = jsonManager.ParseJSON(jsonData)
    
    // HTTP请求示例
    err = httpManager.POSTRequest("https://api.example.com/data", 
        map[string]interface{}{
            "key": "value",
            "number": 42,
        }, "")
    
    // 数据库操作示例
    err = dbManager.CreateDatabase("app.db", "records", 
        "id INTEGER PRIMARY KEY, name TEXT, created_at DATETIME")
}
```

## ⚠️ 错误处理

### 常见错误类型

1. **文件操作错误**
   - 文件不存在
   - 路径权限不足
   - 文件格式不支持

2. **网络请求错误**
   - 连接超时
   - HTTP状态码错误
   - JSON解析失败

3. **数据库错误**
   - SQL语法错误
   - 连接失败
   - 数据约束违反

### 错误处理最佳实践
```go
result, err := manager.SomeOperation()
if err != nil {
    // 记录详细错误信息
    log.Printf("操作失败: %v", err)
    
    // 向用户显示友好错误信息
    ui.ShowError("操作失败，请检查输入数据")
    
    // 根据错误类型采取不同处理
    switch err.(type) {
    case *FileNotFoundError:
        // 文件相关处理
    case *NetworkError:
        // 网络相关处理
    }
}
```

## 🔧 高级配置

### HTTP客户端配置
```go
httpManager := managers.NewHTTPManager(ui)
httpManager.client.
    SetTimeout(60 * time.Second).
    SetRetryCount(3).
    SetRedirectPolicy(resty.FlexibleRedirectPolicy(5))
```

### 数据库连接配置
```go
dbManager := managers.NewDatabaseManager(ui)
// 设置连接池参数
dbManager.db.SetMaxOpenConns(25)
dbManager.db.SetMaxIdleConns(5)
dbManager.db.SetConnMaxLifetime(5 * time.Minute)
```

## 📋 性能优化建议

1. **批量操作**: 尽量使用批量插入而非单条插入
2. **连接复用**: HTTP客户端和数据库连接保持复用
3. **内存管理**: 大文件处理时使用流式读写
4. **并发控制**: 合理控制并发请求数量
5. **缓存策略**: 对频繁访问的数据实施缓存

## 🚀 扩展开发

### 添加新功能模块
1. 创建新的管理器结构体
2. 实现核心功能方法
3. 添加UI接口支持
4. 在按钮处理器中集成
5. 更新事件处理逻辑

### 自定义UI组件
1. 扩展UI接口定义
2. 实现组件创建和事件绑定
3. 添加状态更新方法
4. 集成到主窗口

---

**文档版本**: 1.0  
**最后更新**: 2025年11月12日  
**维护团队**: GO_VCL开发团队