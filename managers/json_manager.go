// managers/json_manager.go - JSON操作管理器
// 该文件实现了JSON操作管理器，提供JSON解析、路径查询、值设置、文件操作等功能
// 支持JSON数据的创建、读取、修改、保存和验证，并提供了丰富的辅助函数用于处理JSON数据
package managers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"windows-gui-app/interfaces"

	"github.com/tidwall/gjson"
)

// JSONManager JSON操作管理器
// 封装了JSON操作相关的功能，提供JSON解析、路径查询、值设置、文件操作等接口
// 与UI组件交互，将操作结果展示在用户界面上
type JSONManager struct {
	uiInstance interfaces.UIInterface // UI实例接口，用于更新界面状态和数据
}

// NewJSONManager 创建JSON管理器实例
// 初始化JSON管理器，并关联UI实例
// 参数:
//   - ui: UI实例，用于显示操作结果（可以为nil，之后通过SetUIInstance设置）
//
// 返回:
//   - *JSONManager: JSON管理器实例
func NewJSONManager(ui interfaces.UIInterface) *JSONManager {
	return &JSONManager{
		uiInstance: ui,
	}
}

// SetUIInstance 设置UI实例
// 允许在创建管理器后设置UI实例，用于解耦合UI和管理器的创建顺序
// 参数:
//   - ui: UI实例接口
func (jm *JSONManager) SetUIInstance(ui interfaces.UIInterface) {
	jm.uiInstance = ui
}

// ParseJSON 解析JSON字符串
// 验证JSON格式并解析为键值对，将结果展示在UI表格中
// 使用gjson库进行高效解析，支持复杂的JSON结构
// 参数:
//   - jsonData: JSON字符串
//
// 返回:
//   - error: 操作错误，nil表示成功
func (jm *JSONManager) ParseJSON(jsonData string) error {
	if strings.TrimSpace(jsonData) == "" {
		return fmt.Errorf("JSON数据不能为空")
	}

	// 更新UI状态，提示用户正在解析JSON
	jm.uiInstance.UpdateStatus("正在解析JSON数据...")

	// 验证JSON格式
	var rawData interface{}
	if err := json.Unmarshal([]byte(jsonData), &rawData); err != nil {
		return fmt.Errorf("JSON格式错误: %v", err)
	}

	// 使用gjson解析并美化输出
	result := gjson.Parse(jsonData)

	// 获取键值对列表
	var data [][]string
	data = append(data, []string{"键", "值", "类型"})

	// 遍历根级属性
	result.ForEach(func(key, value gjson.Result) bool {
		data = append(data, []string{
			key.String(),
			value.String(),
			value.Type.String(),
		})
		return true
	})

	// 更新UI表格显示
	jm.uiInstance.SetTableData(data)

	// 更新状态
	jm.uiInstance.UpdateStatus(fmt.Sprintf("成功解析JSON数据，共有%d个属性", len(data)-1))

	log.Printf("成功解析JSON数据，键值对数量: %d", len(data)-1)
	return nil
}

// GetValueByPath 根据路径获取JSON值
// 使用gjson库的路径查询功能，支持复杂的JSON路径表达式
// 路径格式示例："user.name"、"items.0.name"、"users.#.name"等
// 参数:
//   - jsonData: JSON字符串
//   - path: 路径，如 "user.name" 或 "items.0.name"
//
// 返回:
//   - string: 获取到的值
//   - error: 操作错误，nil表示成功
func (jm *JSONManager) GetValueByPath(jsonData, path string) (string, error) {
	if strings.TrimSpace(jsonData) == "" {
		return "", fmt.Errorf("JSON数据不能为空")
	}

	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("路径不能为空")
	}

	// 解析JSON
	result := gjson.Get(jsonData, path)
	if !result.Exists() {
		return "", fmt.Errorf("路径 '%s' 不存在", path)
	}

	// 更新UI状态
	jm.uiInstance.UpdateStatus(fmt.Sprintf("获取路径 '%s' 的值: %s", path, result.String()))

	log.Printf("成功获取JSON路径 '%s' 的值: %s", path, result.String())
	return result.String(), nil
}

// SetValueByPath 根据路径设置JSON值
// 解析JSON为map结构，递归设置指定路径的值，并返回更新后的JSON字符串
// 支持多级路径设置，自动创建不存在的中间对象
// 参数:
//   - jsonData: 原始JSON字符串
//   - path: 路径，如 "user.age"
//   - value: 要设置的值
//
// 返回:
//   - string: 更新后的JSON字符串
//   - error: 操作错误，nil表示成功
func (jm *JSONManager) SetValueByPath(jsonData, path, value string) (string, error) {
	if strings.TrimSpace(jsonData) == "" {
		return "", fmt.Errorf("JSON数据不能为空")
	}

	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("路径不能为空")
	}

	// 解析JSON到map结构
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return "", fmt.Errorf("JSON格式错误: %v", err)
	}

	// 解析路径
	keys := strings.Split(path, ".")
	if len(keys) == 0 {
		return "", fmt.Errorf("无效的路径格式")
	}

	// 设置值
	if err := jm.setValueInMap(data, keys, value); err != nil {
		return "", fmt.Errorf("设置值失败: %v", err)
	}

	// 重新序列化为JSON
	updatedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化JSON失败: %v", err)
	}

	// 更新UI
	jm.uiInstance.UpdateStatus(fmt.Sprintf("成功更新路径 '%s' 的值", path))

	log.Printf("成功设置JSON路径 '%s' 的值为: %s", path, value)
	return string(updatedJSON), nil
}

// setValueInMap 在嵌套map中设置值
// 递归遍历嵌套的map结构，设置指定路径的值
// 如果中间路径不存在，自动创建新的map对象
// 参数:
//   - data: map结构数据
//   - keys: 路径键数组
//   - value: 要设置的值
//
// 返回:
//   - error: 错误信息
func (jm *JSONManager) setValueInMap(data map[string]interface{}, keys []string, value string) error {
	if len(keys) == 1 {
		// 最后一个键，直接设置值
		data[keys[0]] = value
		return nil
	}

	// 获取或创建嵌套map
	currentKey := keys[0]
	if _, exists := data[currentKey]; !exists {
		data[currentKey] = make(map[string]interface{})
	}

	// 类型断言
	nestedMap, ok := data[currentKey].(map[string]interface{})
	if !ok {
		return fmt.Errorf("路径 '%s' 无法访问，类型不匹配", strings.Join(keys, "."))
	}

	// 递归设置值
	return jm.setValueInMap(nestedMap, keys[1:], value)
}

// CreateJSON 创建新的JSON数据
// 将任意类型的数据序列化为格式化的JSON字符串
// 支持结构体、map、slice等Go数据类型
// 参数:
//   - data: 要转换为JSON的数据
//
// 返回:
//   - string: JSON字符串
//   - error: 操作错误，nil表示成功
func (jm *JSONManager) CreateJSON(data interface{}) (string, error) {
	if data == nil {
		return "", fmt.Errorf("数据不能为空")
	}

	// 序列化为JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("创建JSON失败: %v", err)
	}

	// 更新UI状态
	jm.uiInstance.UpdateStatus("成功创建新的JSON数据")

	log.Printf("成功创建JSON数据，大小: %d字节", len(jsonData))
	return string(jsonData), nil
}

// SaveJSONToFile 保存JSON到文件
// 将JSON字符串保存到指定文件路径，自动创建必要的目录结构
// 确保文件扩展名为.json，并设置适当的文件权限
// 参数:
//   - jsonData: JSON字符串
//   - filePath: 文件路径
//
// 返回:
//   - error: 操作错误，nil表示成功
func (jm *JSONManager) SaveJSONToFile(jsonData, filePath string) error {
	if strings.TrimSpace(jsonData) == "" {
		return fmt.Errorf("JSON数据不能为空")
	}

	if filePath == "" {
		return fmt.Errorf("文件路径不能为空")
	}

	// 确保文件扩展名
	if !strings.HasSuffix(strings.ToLower(filePath), ".json") {
		filePath += ".json"
	}

	// 创建目录
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 更新UI状态，提示用户正在保存文件
	jm.uiInstance.UpdateStatus("正在保存JSON文件...")

	// 写入文件
	if err := os.WriteFile(filePath, []byte(jsonData), 0644); err != nil {
		return fmt.Errorf("保存JSON文件失败: %v", err)
	}

	// 更新状态
	jm.uiInstance.UpdateStatus(fmt.Sprintf("成功保存JSON文件: %s", filePath))

	log.Printf("成功保存JSON文件: %s，大小: %d字节", filePath, len(jsonData))
	return nil
}

// LoadJSONFromFile 从文件加载JSON
// 从指定文件路径读取JSON数据，验证格式并格式化输出
// 自动处理文件不存在和格式错误的情况
// 参数:
//   - filePath: 文件路径
//
// 返回:
//   - string: JSON字符串
//   - error: 操作错误，nil表示成功
func (jm *JSONManager) LoadJSONFromFile(filePath string) (string, error) {
	if filePath == "" {
		return "", fmt.Errorf("文件路径不能为空")
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("文件不存在: %s", filePath)
	}

	// 更新UI状态，提示用户正在加载文件
	jm.uiInstance.UpdateStatus("正在加载JSON文件...")

	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取JSON文件失败: %v", err)
	}

	// 验证JSON格式
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return "", fmt.Errorf("JSON文件格式错误: %v", err)
	}

	// 美化输出
	formattedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("格式化JSON失败: %v", err)
	}

	// 更新状态
	fileInfo, _ := os.Stat(filePath)
	jm.uiInstance.UpdateStatus(fmt.Sprintf("成功加载JSON文件: %s，大小: %.2fKB",
		filepath.Base(filePath), float64(fileInfo.Size())/1024))

	log.Printf("成功加载JSON文件: %s", filePath)
	return string(formattedJSON), nil
}

// ValidateJSON 验证JSON格式
// 检查JSON字符串的格式是否正确，并返回详细的验证结果
// 支持对象和数组两种基本JSON结构
// 参数:
//   - jsonData: JSON字符串
//
// 返回:
//   - bool: 是否有效
//   - string: 验证结果描述
func (jm *JSONManager) ValidateJSON(jsonData string) (bool, string) {
	if strings.TrimSpace(jsonData) == "" {
		return false, "JSON数据为空"
	}

	// 尝试解析
	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return false, fmt.Sprintf("JSON格式错误: %v", err)
	}

	// 检查数据结构
	result := gjson.Parse(jsonData)
	if !result.IsObject() && !result.IsArray() {
		return false, "JSON必须是对象或数组格式"
	}

	return true, "JSON格式验证通过"
}

// GetJSONInfo 获取JSON数据信息
// 分析JSON数据的结构、类型、大小等统计信息
// 返回包含详细信息的map，便于展示和分析
// 参数:
//   - jsonData: JSON字符串
//
// 返回:
//   - map[string]interface{}: JSON信息
//   - error: 操作错误，nil表示成功
func (jm *JSONManager) GetJSONInfo(jsonData string) (map[string]interface{}, error) {
	if strings.TrimSpace(jsonData) == "" {
		return nil, fmt.Errorf("JSON数据不能为空")
	}

	// 解析JSON
	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, fmt.Errorf("JSON格式错误: %v", err)
	}

	info := make(map[string]interface{})

	// 使用gjson获取详细信息
	result := gjson.Parse(jsonData)

	// 基本信息
	info["数据类型"] = result.Type.String()
	info["数据大小"] = len(jsonData)
	info["字符数"] = len([]rune(jsonData))

	// 统计信息
	if result.IsObject() {
		// 获取对象的所有键
		keys := getObjectKeys(jsonData)
		info["对象属性数"] = len(keys)
	} else if result.IsArray() {
		info["数组元素数"] = result.Array()
	}

	// 获取示例数据
	if result.IsObject() {
		keys := getObjectKeys(jsonData)
		if len(keys) > 0 {
			info["前3个属性"] = keys[:min(3, len(keys))]
		}
	} else if result.IsArray() && result.Array() != nil {
		info["数组类型"] = "多种类型"
		if len(result.Array()) > 0 {
			info["第一个元素类型"] = result.Array()[0].Type.String()
		}
	}

	return info, nil
}

// min 辅助函数
// 返回两个整数中的较小值
// 参数:
//   - a: 第一个整数
//   - b: 第二个整数
//
// 返回:
//   - int: 较小的整数值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetDefaultJSON 获取默认的示例JSON数据
// 返回一个包含常见数据类型的示例JSON，用于演示和测试
// 包括字符串、数字、数组、嵌套对象等多种数据结构
// 返回:
//   - string: 默认JSON数据
func (jm *JSONManager) GetDefaultJSON() string {
	// 默认示例JSON数据
	defaultData := map[string]interface{}{
		"姓名": "张三",
		"年龄": 25,
		"邮箱": "zhangsan@example.com",
		"技能": []string{
			"Go语言",
			"数据库",
			"Web开发",
		},
		"地址": map[string]interface{}{
			"城市": "北京",
			"区县": "朝阳区",
			"街道": "建国路",
		},
		"工作": map[string]interface{}{
			"公司": "科技公司",
			"职位": "软件工程师",
			"经验": 3,
		},
	}

	jsonData, _ := json.MarshalIndent(defaultData, "", "  ")
	return string(jsonData)
}

// getObjectKeys 获取JSON对象的键列表
// 解析JSON字符串并返回所有顶级键的列表
// 参数:
//   - jsonData: JSON字符串
//
// 返回:
//   - []string: 键列表
func getObjectKeys(jsonData string) []string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil
	}

	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	return keys
}
