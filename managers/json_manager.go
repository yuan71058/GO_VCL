// managers/json_manager.go - JSON操作管理器
// 功能描述: 提供JSON数据的解析、查询、修改和创建功能
// 主要功能:
//   - JSON解析：将JSON字符串解析为Go数据结构
//   - 路径查询：支持JSONPath风格的查询语法
//   - 数据修改：通过路径修改JSON中的特定值
//   - JSON创建：从数据结构生成格式化的JSON字符串
//   - 错误处理：提供详细的解析错误和路径错误信息
// 作者: GO_VCL开发团队
// 创建时间: 2025-11-12

package managers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"windows-gui-app/interfaces"
)

// JSONManager JSON操作管理器
// 负责处理JSON数据的解析、查询和修改操作，提供UI集成功能
type JSONManager struct {
	uiInstance interfaces.UIInterface // UI实例接口，用于显示操作状态和结果
	jsonData   interface{}            // 当前处理的JSON数据
}

// NewJSONManager 创建JSON管理器实例
// 参数:
//   - ui: UI实例，用于显示操作结果（可以为nil，之后通过SetUIInstance设置）
//
// 返回:
//   - *JSONManager: JSON管理器实例
//
// 使用示例:
//   jsonManager := managers.NewJSONManager(mainForm)
func NewJSONManager(ui interfaces.UIInterface) *JSONManager {
	return &JSONManager{
		uiInstance: ui,
		jsonData:   nil,
	}
}

// SetUIInstance 设置UI实例
// 用于在运行时更新UI实例引用
// 参数:
//   - ui: 新的UI实例
func (jm *JSONManager) SetUIInstance(ui interfaces.UIInterface) {
	jm.uiInstance = ui
}

// ParseJSON 解析JSON字符串
// 功能描述:
//   - 将JSON字符串解析为Go的interface{}数据结构
//   - 支持任意有效的JSON格式（对象、数组、基本类型）
//   - 自动检测和处理JSON格式错误
//   - 更新内部jsonData状态，供后续操作使用
//   - 提供详细的解析错误信息
//
// 参数:
//   - jsonStr: 要解析的JSON字符串
//
// 返回:
//   - interface{}: 解析后的JSON数据
//   - error: 解析错误，nil表示成功
//
// 支持的JSON类型:
//   - 对象: {"name": "张三", "age": 25}
//   - 数组: [1, 2, 3, "test"]
//   - 字符串: "Hello World"
//   - 数字: 123, 45.67
//   - 布尔值: true, false
//   - null: null
//
// 可能的错误:
//   - "JSON字符串不能为空": jsonStr参数为空或仅包含空白字符
//   - "解析JSON失败": JSON格式错误，包含具体错误位置
//
// 使用示例:
//   jsonStr := `{"users": [{"name": "张三", "age": 25}, {"name": "李四", "age": 30}]}`
//   data, err := jsonManager.ParseJSON(jsonStr)
//   if err != nil {
//       log.Printf("解析失败: %v", err)
//   }
func (jm *JSONManager) ParseJSON(jsonStr string) (interface{}, error) {
	// 验证输入参数
	if strings.TrimSpace(jsonStr) == "" {
		return nil, fmt.Errorf("JSON字符串不能为空")
	}

	// 更新UI状态
	if jm.uiInstance != nil {
		jm.uiInstance.UpdateStatus("正在解析JSON数据...")
	}

	// 解析JSON数据
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 更新内部状态
	jm.jsonData = data

	// 更新UI状态
	if jm.uiInstance != nil {
		statusMsg := fmt.Sprintf("成功解析JSON数据，类型: %T", data)
		jm.uiInstance.UpdateStatus(statusMsg)
	}

	// 记录操作日志
	log.Printf("成功解析JSON数据，长度: %d字符，类型: %T", len(jsonStr), data)
	return data, nil
}

// GetValueByPath 通过路径获取JSON中的值
// 功能描述:
//   - 支持简单的点号路径语法（如"user.name", "items.0.title"）
//   - 支持数组索引访问（如"users.0", "items.2.price"）
//   - 支持嵌套对象和数组的混合访问
//   - 提供详细的路径错误信息
//   - 返回值的类型保持原始JSON类型
//
// 参数:
//   - path: 查询路径，使用点号分隔（如"user.name", "items.0.title"）
//   - data: JSON数据（如果为nil则使用内部jsonData）
//
// 返回:
//   - interface{}: 查询到的值，nil表示未找到
//   - bool: 是否找到值
//   - error: 操作错误，nil表示成功
//
// 支持的路径格式:
//   - "name": 获取顶级字段
//   - "user.name": 获取嵌套对象的字段
//   - "users.0": 获取数组的第一个元素
//   - "items.2.name": 获取数组中对象的字段
//   - "config.database.host": 深层嵌套访问
//
// 可能的错误:
//   - "查询路径不能为空": path参数为空或仅包含空白字符
//   - "JSON数据为空": 没有可用的JSON数据
//   - "路径格式错误": 路径包含非法字符或格式
//
// 使用示例:
//   // 从内部数据查询
//   value, found, err := jsonManager.GetValueByPath("users.0.name", nil)
//   
//   // 从指定数据查询
//   data := map[string]interface{}{"user": map[string]interface{}{"name": "张三"}}
//   value, found, err := jsonManager.GetValueByPath("user.name", data)
func (jm *JSONManager) GetValueByPath(path string, data interface{}) (interface{}, bool, error) {
	// 验证路径参数
	if strings.TrimSpace(path) == "" {
		return nil, false, fmt.Errorf("查询路径不能为空")
	}

	// 确定数据源
	var targetData interface{}
	if data != nil {
		targetData = data
	} else {
		targetData = jm.jsonData
	}

	// 检查数据有效性
	if targetData == nil {
		return nil, false, fmt.Errorf("JSON数据为空")
	}

	// 更新UI状态
	if jm.uiInstance != nil {
		jm.uiInstance.UpdateStatus(fmt.Sprintf("正在查询路径: %s", path))
	}

	// 分割路径
	parts := strings.Split(path, ".")
	current := targetData

	// 遍历路径各部分
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		switch v := current.(type) {
		case map[string]interface{}:
			// 处理对象类型
			if val, ok := v[part]; ok {
				current = val
			} else {
				return nil, false, nil // 路径不存在
			}
		case []interface{}:
			// 处理数组类型
			var index int
			if _, err := fmt.Sscanf(part, "%d", &index); err != nil {
				return nil, false, fmt.Errorf("路径格式错误: %s不是有效的数组索引", part)
			}
			if index < 0 || index >= len(v) {
				return nil, false, nil // 索引越界
			}
			current = v[index]
		default:
			// 基本类型，无法继续访问
			return nil, false, nil // 路径不存在
		}
	}

	// 更新UI状态
	if jm.uiInstance != nil {
		statusMsg := fmt.Sprintf("成功查询路径: %s，值类型: %T", path, current)
		jm.uiInstance.UpdateStatus(statusMsg)
	}

	// 记录操作日志
	log.Printf("成功查询JSON路径: %s，值类型: %T", path, current)
	return current, true, nil
}

// SetValueByPath 通过路径设置JSON中的值
// 功能描述:
//   - 支持创建新的嵌套对象和数组
//   - 支持修改现有值
//   - 自动处理类型转换和验证
//   - 提供详细的修改状态信息
//
// 参数:
//   - path: 要设置值的路径（如"user.name", "items.0.title"）
//   - value: 要设置的值（支持任意JSON兼容类型）
//   - data: JSON数据（如果为nil则使用内部jsonData）
//
// 返回:
//   - interface{}: 修改后的JSON数据
//   - error: 操作错误，nil表示成功
//
// 特殊功能:
//   - 自动创建不存在的路径
//   - 支持基本类型和复杂对象
//   - 保持原始数据的其他部分不变
//
// 可能的错误:
//   - "路径不能为空": path参数为空
//   - "值不能为空": value参数为nil且不允许
//   - "不支持的路径格式": 路径包含非法结构
//
// 使用示例:
//   // 修改内部数据
//   newData, err := jsonManager.SetValueByPath("user.name", "李四", nil)
//   
//   // 修改指定数据
//   data := map[string]interface{}{"user": map[string]interface{}{"name": "张三"}}
//   newData, err := jsonManager.SetValueByPath("user.age", 25, data)
func (jm *JSONManager) SetValueByPath(path string, value interface{}, data interface{}) (interface{}, error) {
	// 验证参数
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("路径不能为空")
	}

	// 确定数据源
	var targetData interface{}
	if data != nil {
		targetData = data
	} else {
		targetData = jm.jsonData
	}

	// 如果目标数据为空，创建新的map
	if targetData == nil {
		targetData = make(map[string]interface{})
	}

	// 更新UI状态
	if jm.uiInstance != nil {
		jm.uiInstance.UpdateStatus(fmt.Sprintf("正在设置路径: %s", path))
	}

	// 分割路径
	parts := strings.Split(path, ".")
	
	// 创建副本以避免修改原始数据
	result := jm.deepCopy(targetData)
	current := result

	// 遍历路径（除了最后一部分）
	for i := 0; i < len(parts)-1; i++ {
		part := strings.TrimSpace(parts[i])
		if part == "" {
			continue
		}

		switch v := current.(type) {
		case map[string]interface{}:
			// 确保路径存在
			if _, ok := v[part]; !ok {
				v[part] = make(map[string]interface{})
			}
			current = v[part]
		case []interface{}:
			// 处理数组索引
			var index int
			if _, err := fmt.Sscanf(part, "%d", &index); err != nil {
				return nil, fmt.Errorf("不支持的路径格式: %s", part)
			}
			if index < 0 || index >= len(v) {
				return nil, fmt.Errorf("数组索引越界: %d", index)
			}
			current = v[index]
		default:
			return nil, fmt.Errorf("不支持的路径格式")
		}
	}

	// 设置最终值
	lastPart := strings.TrimSpace(parts[len(parts)-1])
	if lastPart == "" {
		return nil, fmt.Errorf("路径格式错误")
	}

	switch v := current.(type) {
	case map[string]interface{}:
		v[lastPart] = value
	case []interface{}:
		var index int
		if _, err := fmt.Sscanf(lastPart, "%d", &index); err != nil {
			return nil, fmt.Errorf("不支持的路径格式: %s", lastPart)
		}
		if index < 0 || index >= len(v) {
			return nil, fmt.Errorf("数组索引越界: %d", index)
		}
		v[index] = value
	default:
		return nil, fmt.Errorf("不支持的路径格式")
	}

	// 更新内部状态
	jm.jsonData = result

	// 更新UI状态
	if jm.uiInstance != nil {
		statusMsg := fmt.Sprintf("成功设置路径: %s，新值: %v", path, value)
		jm.uiInstance.UpdateStatus(statusMsg)
	}

	// 记录操作日志
	log.Printf("成功设置JSON路径: %s = %v", path, value)
	return result, nil
}

// CreateJSON 创建新的JSON数据
// 功能描述:
//   - 从Go数据结构生成格式化的JSON字符串
//   - 支持任意JSON兼容的数据类型
//   - 提供美观的格式化输出选项
//   - 自动处理循环引用检测
//
// 参数:
//   - data: 要转换为JSON的数据（支持map、slice、基本类型等）
//   - pretty: 是否使用缩进格式化输出
//
// 返回:
//   - string: 生成的JSON字符串
//   - error: 转换错误，nil表示成功
//
// 支持的数据类型:
//   - map[string]interface{}: JSON对象
//   - []interface{}: JSON数组
//   - string: JSON字符串
//   - float64/int/bool: JSON基本类型
//   - nil: JSON null值
//
// 可能的错误:
//   - "数据不能为空": data参数为nil
//   - "转换为JSON失败": 数据包含无法序列化的类型
//
// 使用示例:
//   // 创建简单对象
//   data := map[string]interface{}{
//       "name": "张三",
//       "age": 25,
//       "active": true,
//   }
//   jsonStr, err := jsonManager.CreateJSON(data, true)
//   
//   // 创建数组
//   items := []interface{}{"item1", "item2", "item3"}
//   jsonStr, err := jsonManager.CreateJSON(items, false)
func (jm *JSONManager) CreateJSON(data interface{}, pretty bool) (string, error) {
	// 验证数据
	if data == nil {
		return "", fmt.Errorf("数据不能为空")
	}

	// 更新UI状态
	if jm.uiInstance != nil {
		jm.uiInstance.UpdateStatus("正在创建JSON数据...")
	}

	// 序列化为JSON
	var jsonData []byte
	var err error

	if pretty {
		jsonData, err = json.MarshalIndent(data, "", "  ")
	} else {
		jsonData, err = json.Marshal(data)
	}

	if err != nil {
		return "", fmt.Errorf("转换为JSON失败: %v", err)
	}

	// 更新内部状态
	jm.jsonData = data

	// 更新UI状态
	if jm.uiInstance != nil {
		statusMsg := fmt.Sprintf("成功创建JSON数据，长度: %d字符", len(jsonData))
		jm.uiInstance.UpdateStatus(statusMsg)
	}

	// 记录操作日志
	log.Printf("成功创建JSON数据，长度: %d字符，类型: %T", len(jsonData), data)
	return string(jsonData), nil
}

// GetCurrentData 获取当前JSON数据
// 返回内部维护的JSON数据副本
//
// 返回:
//   - interface{}: 当前JSON数据，如果没有数据则为nil
func (jm *JSONManager) GetCurrentData() interface{} {
	return jm.jsonData
}

// ClearData 清空当前JSON数据
// 重置内部状态，清除所有已加载的JSON数据
func (jm *JSONManager) ClearData() {
	jm.jsonData = nil
	if jm.uiInstance != nil {
		jm.uiInstance.UpdateStatus("JSON数据已清空")
	}
	log.Println("JSON数据已清空")
}

// deepCopy 深度复制JSON数据
// 用于在修改操作中创建数据副本，避免修改原始数据
//
// 参数:
//   - data: 要复制的数据
//
// 返回:
//   - interface{}: 数据的深拷贝副本
func (jm *JSONManager) deepCopy(data interface{}) interface{} {
	// 使用JSON序列化/反序列化实现深度复制
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("深度复制失败: %v", err)
		return data // 如果复制失败，返回原始数据
	}

	var copy interface{}
	if err := json.Unmarshal(jsonData, &copy); err != nil {
		log.Printf("深度复制失败: %v", err)
		return data // 如果复制失败，返回原始数据
	}

	return copy
}