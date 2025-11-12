// utils/helpers.go - 工具函数库
// 该文件提供了常用的工具函数集合，包括字符串处理、文件操作、数据转换和验证等功能
// 这些函数旨在简化常见编程任务，提高代码复用性和可维护性
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// StringToMD5 将字符串转换为MD5哈希值
// 使用MD5算法计算字符串的哈希值，返回32位十六进制字符串
// 注意：MD5不适用于密码存储等安全场景，仅用于数据完整性校验
// 参数:
//   - input: 要计算哈希的输入字符串
//
// 返回:
//   - string: 32位十六进制MD5哈希值
func StringToMD5(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// FileExists 检查文件是否存在
// 通过os.Stat函数检查文件系统中的文件是否存在，处理可能的错误情况
// 该函数可以检查文件和目录的存在性
// 参数:
//   - filePath: 要检查的文件或目录路径
//
// 返回:
//   - bool: 文件或目录存在返回true，否则返回false
func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateDirectory 创建目录（递归创建）
// 使用os.MkdirAll递归创建目录及其所有父目录，设置权限为0755
// 如果目录已存在，不会报错，直接返回成功
// 参数:
//   - dirPath: 要创建的目录路径
//
// 返回:
//   - error: 创建错误，nil表示成功
func CreateDirectory(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// GetFileExtension 获取文件扩展名
// 使用filepath.Ext提取文件扩展名，并转换为小写形式
// 扩展名包括点号（如".txt"），如果文件没有扩展名则返回空字符串
// 参数:
//   - filePath: 文件路径
//
// 返回:
//   - string: 文件扩展名（小写，包含点号）
func GetFileExtension(filePath string) string {
	ext := filepath.Ext(filePath)
	return strings.ToLower(ext)
}

// IsValidEmail 验证邮箱格式
// 使用正则表达式验证邮箱地址是否符合标准格式
// 支持常见的邮箱格式，包括字母、数字、点号、下划线、百分号、加号和减号
// 参数:
//   - email: 要验证的邮箱字符串
//
// 返回:
//   - bool: 邮箱格式有效返回true，否则返回false
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// IsValidURL 验证URL格式
// 使用正则表达式验证URL是否符合标准格式，支持HTTP和HTTPS协议
// 检查URL是否以http://或https://开头，并包含有效的域名部分
// 参数:
//   - url: 要验证的URL字符串
//
// 返回:
//   - bool: URL格式有效返回true，否则返回false
func IsValidURL(url string) bool {
	pattern := `^https?://[^\s/$.?#].[^\s]*$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

// ParseInt 字符串转换为整数
// 安全地将字符串转换为整数，转换失败时返回默认值
// 处理常见的转换错误，避免程序因格式错误而崩溃
// 参数:
//   - str: 要转换的字符串
//   - defaultValue: 转换失败时返回的默认值
//
// 返回:
//   - int: 转换后的整数或默认值
func ParseInt(str string, defaultValue int) int {
	if value, err := strconv.Atoi(str); err == nil {
		return value
	}
	return defaultValue
}

// ParseFloat 字符串转换为浮点数
// 安全地将字符串转换为64位浮点数，转换失败时返回默认值
// 处理常见的转换错误，避免程序因格式错误而崩溃
// 参数:
//   - str: 要转换的字符串
//   - defaultValue: 转换失败时返回的默认值
//
// 返回:
//   - float64: 转换后的浮点数或默认值
func ParseFloat(str string, defaultValue float64) float64 {
	if value, err := strconv.ParseFloat(str, 64); err == nil {
		return value
	}
	return defaultValue
}

// ToTitleCase 转换为标题格式
// 将输入文本转换为标题格式，每个单词的首字母大写，其余字母保持原样
// 适用于将普通文本转换为标题或标题栏显示格式
// 参数:
//   - text: 要转换的输入文本
//
// 返回:
//   - string: 标题格式文本
func ToTitleCase(text string) string {
	words := strings.Fields(text)
	for i, word := range words {
		runes := []rune(word)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

// ToCamelCase 转换为驼峰格式
// 将输入文本转换为驼峰命名格式，第一个单词首字母小写，后续单词首字母大写
// 非字母数字字符被替换为空格，然后转换为驼峰格式
// 适用于将普通文本转换为编程变量名或标识符
// 参数:
//   - text: 要转换的输入文本
//
// 返回:
//   - string: 驼峰格式文本
func ToCamelCase(text string) string {
	// 替换非字母数字字符为空格
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	words := re.Split(text, -1)
	
	var result strings.Builder
	for i, word := range words {
		if word == "" {
			continue
		}
		
		// 转换为小写
		lowerWord := strings.ToLower(word)
		
		// 第一个单词的首字母保持小写，后续单词首字母大写
		if i > 0 && len(lowerWord) > 0 {
			upperBytes := []byte(lowerWord)
			upperBytes[0] = byte(unicode.ToUpper(rune(lowerWord[0])))
			result.Write(upperBytes)
		} else {
			result.WriteString(lowerWord)
		}
	}
	
	return result.String()
}

// GetCurrentTimeStr 获取当前时间字符串
// 获取当前时间并按照指定格式格式化为字符串，默认格式为"2006-01-02 15:04:05"
// 使用Go语言特有的时间格式布局，参考时间必须是2006年1月2日15:04:05
// 参数:
//   - format: 可选的时间格式字符串，默认为"2006-01-02 15:04:05"
//
// 返回:
//   - string: 格式化的时间字符串
func GetCurrentTimeStr(format ...string) string {
	t := time.Now()
	layout := "2006-01-02 15:04:05"
	if len(format) > 0 && format[0] != "" {
		layout = format[0]
	}
	return t.Format(layout)
}

// TruncateString 截断字符串
// 将字符串截断到指定长度，如果超过最大长度则添加后缀
// 适用于显示长文本的摘要或预览，保持界面整洁
// 参数:
//   - text: 要截断的字符串
//   - maxLen: 最大长度
//   - suffix: 可选的后缀，默认为"..."
//
// 返回:
//   - string: 截断后的字符串
func TruncateString(text string, maxLen int, suffix ...string) string {
	if len(text) <= maxLen {
		return text
	}
	
	suf := "..."
	if len(suffix) > 0 {
		suf = suffix[0]
	}
	
	truncatedLen := maxLen - len(suf)
	if truncatedLen < 0 {
		return suf
	}
	
	return text[:truncatedLen] + suf
}

// Contains 检查字符串是否包含子字符串
// 检查主字符串是否包含指定的子字符串，可选择是否区分大小写
// 提供比标准库更灵活的大小写控制选项
// 参数:
//   - text: 主字符串
//   - substring: 要查找的子字符串
//   - caseSensitive: 可选参数，是否区分大小写，默认为true
//
// 返回:
//   - bool: 包含子字符串返回true，否则返回false
func Contains(text, substring string, caseSensitive ...bool) bool {
	isCaseSensitive := true
	if len(caseSensitive) > 0 {
		isCaseSensitive = caseSensitive[0]
	}
	
	if isCaseSensitive {
		return strings.Contains(text, substring)
	}
	
	return strings.Contains(strings.ToLower(text), strings.ToLower(substring))
}

// SplitLines 按行分割字符串
// 将多行文本按换行符分割为行数组，适用于处理多行文本内容
// 参数:
//   - text: 输入文本
//
// 返回:
//   - []string: 行数组
func SplitLines(text string) []string {
	return strings.Split(text, "\n")
}

// JoinLines 将行数组合并为字符串
// 将行数组使用指定分隔符合并为字符串，默认分隔符为换行符
// 参数:
//   - lines: 行数组
//   - separator: 可选的分隔符，默认为换行符
//
// 返回:
//   - string: 合并后的字符串
func JoinLines(lines []string, separator ...string) string {
	sep := "\n"
	if len(separator) > 0 {
		sep = separator[0]
	}
	return strings.Join(lines, sep)
}

// CleanWhitespace 清理空白字符
// 清理文本中的多余空白字符，包括首尾空白和连续空白
// 适用于规范化用户输入或处理从外部获取的文本数据
// 参数:
//   - text: 输入文本
//
// 返回:
//   - string: 清理后的文本
func CleanWhitespace(text string) string {
	// 移除首尾空白
	text = strings.TrimSpace(text)
	
	// 将多个连续空白替换为单个空格
	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")
	
	return text
}

// GetFileSize 获取文件大小（字节）
// 获取指定文件的大小，以字节为单位，如果文件不存在则返回-1
// 适用于检查文件大小或计算存储空间使用情况
// 参数:
//   - filePath: 文件路径
//
// 返回:
//   - int64: 文件大小（字节），-1表示文件不存在
func GetFileSize(filePath string) int64 {
	if info, err := os.Stat(filePath); err == nil {
		return info.Size()
	}
	return -1
}

// FormatFileSize 格式化文件大小
// 将字节大小的文件格式化为人类可读的形式，自动选择合适的单位（B、KB、MB、GB）
// 适用于在用户界面中显示文件大小信息
// 参数:
//   - size: 文件大小（字节）
//
// 返回:
//   - string: 格式化后的文件大小
func FormatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
	}
}

// ReverseString 反转字符串
// 将输入字符串按字符顺序反转，适用于处理需要倒序显示的文本
// 支持Unicode字符，正确处理多字节字符
// 参数:
//   - text: 输入字符串
//
// 返回:
//   - string: 反转后的字符串
func ReverseString(text string) string {
	runes := []rune(text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsEmpty 检查字符串是否为空（包含空白字符）
// 检查字符串是否为空或只包含空白字符，适用于验证用户输入
// 参数:
//   - text: 要检查的字符串
//
// 返回:
//   - bool: 字符串为空或只包含空白字符返回true，否则返回false
func IsEmpty(text string) bool {
	return strings.TrimSpace(text) == ""
}

// SafeSlice 安全切片操作
// 提供安全的切片操作，自动处理索引越界情况，避免panic
// 适用于处理不确定长度的切片，确保程序稳定性
// 参数:
//   - slice: 源切片
//   - start: 起始索引
//   - end: 结束索引（不包括）
//
// 返回:
//   - []interface{}: 安全切片
func SafeSlice(slice []interface{}, start, end int) []interface{} {
	if start < 0 {
		start = 0
	}
	if end > len(slice) {
		end = len(slice)
	}
	if start >= end {
		return []interface{}{}
	}
	
	result := make([]interface{}, end-start)
	copy(result, slice[start:end])
	return result
}

// ConvertSlice 转换切片类型
// 将interface{}类型的切片转换为字符串切片，使用fmt.Sprintf进行格式化
// 适用于处理来自动态数据源的切片，统一转换为字符串类型
// 参数:
//   - slice: 源切片
//
// 返回:
//   - []string: 字符串切片
func ConvertSlice(slice []interface{}) []string {
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = fmt.Sprintf("%v", v)
	}
	return result
}

// RemoveDuplicates 移除切片中的重复元素
// 移除字符串切片中的重复元素，保持原始顺序，使用map记录已出现的元素
// 适用于需要去重的数据处理场景，如标签列表、关键词列表等
// 参数:
//   - slice: 源切片
//
// 返回:
//   - []string: 去重后的切片
func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// SafeString 安全字符串操作
// 提供安全的字符串处理，包括空值检查、长度限制和默认值处理
// 适用于处理用户输入或外部数据源的数据，确保数据安全性和一致性
// 参数:
//   - text: 输入文本
//   - maxLen: 最大长度，0表示不限制
//   - defaultValue: 可选的默认值，当输入为空时使用
//
// 返回:
//   - string: 安全字符串
func SafeString(text string, maxLen int, defaultValue ...string) string {
	if IsEmpty(text) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	
	if maxLen > 0 && len(text) > maxLen {
		return TruncateString(text, maxLen)
	}
	
	return text
}