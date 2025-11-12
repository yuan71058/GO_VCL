// utils/helpers.go - 工具函数库
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

// StringToMD5 将字符串转换为MD5哈希
// 参数:
//   - input: 输入字符串
// 返回:
//   - string: MD5哈希值
func StringToMD5(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// FileExists 检查文件是否存在
// 参数:
//   - filePath: 文件路径
// 返回:
//   - bool: 文件是否存在
func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateDirectory 创建目录（递归创建）
// 参数:
//   - dirPath: 目录路径
// 返回:
//   - error: 创建错误，nil表示成功
func CreateDirectory(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// GetFileExtension 获取文件扩展名
// 参数:
//   - filePath: 文件路径
// 返回:
//   - string: 文件扩展名（小写）
func GetFileExtension(filePath string) string {
	ext := filepath.Ext(filePath)
	return strings.ToLower(ext)
}

// IsValidEmail 验证邮箱格式
// 参数:
//   - email: 邮箱字符串
// 返回:
//   - bool: 是否为有效邮箱
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// IsValidURL 验证URL格式
// 参数:
//   - url: URL字符串
// 返回:
//   - bool: 是否为有效URL
func IsValidURL(url string) bool {
	pattern := `^https?://[^\s/$.?#].[^\s]*$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

// ParseInt 字符串转换为整数
// 参数:
//   - str: 要转换的字符串
//   - defaultValue: 默认值
// 返回:
//   - int: 转换后的整数
func ParseInt(str string, defaultValue int) int {
	if value, err := strconv.Atoi(str); err == nil {
		return value
	}
	return defaultValue
}

// ParseFloat 字符串转换为浮点数
// 参数:
//   - str: 要转换的字符串
//   - defaultValue: 默认值
// 返回:
//   - float64: 转换后的浮点数
func ParseFloat(str string, defaultValue float64) float64 {
	if value, err := strconv.ParseFloat(str, 64); err == nil {
		return value
	}
	return defaultValue
}

// ToTitleCase 转换为标题格式
// 参数:
//   - text: 输入文本
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
// 参数:
//   - text: 输入文本
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
// 参数:
//   - format: 时间格式，默认为"2006-01-02 15:04:05"
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
// 参数:
//   - text: 要截断的字符串
//   - maxLen: 最大长度
//   - suffix: 后缀，默认为"..."
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
// 参数:
//   - text: 主字符串
//   - substring: 子字符串
//   - caseSensitive: 是否区分大小写
// 返回:
//   - bool: 是否包含
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
// 参数:
//   - text: 输入文本
// 返回:
//   - []string: 行数组
func SplitLines(text string) []string {
	return strings.Split(text, "\n")
}

// JoinLines 将行数组合并为字符串
// 参数:
//   - lines: 行数组
//   - separator: 分隔符，默认为换行符
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
// 参数:
//   - text: 输入文本
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
// 参数:
//   - filePath: 文件路径
// 返回:
//   - int64: 文件大小，-1表示文件不存在
func GetFileSize(filePath string) int64 {
	if info, err := os.Stat(filePath); err == nil {
		return info.Size()
	}
	return -1
}

// FormatFileSize 格式化文件大小
// 参数:
//   - size: 文件大小（字节）
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
// 参数:
//   - text: 输入字符串
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
// 参数:
//   - text: 要检查的字符串
// 返回:
//   - bool: 是否为空
func IsEmpty(text string) bool {
	return strings.TrimSpace(text) == ""
}

// SafeSlice 安全切片操作
// 参数:
//   - slice: 源切片
//   - start: 起始索引
//   - end: 结束索引（不包括）
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
// 参数:
//   - slice: 源切片
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
// 参数:
//   - slice: 源切片
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
// 参数:
//   - text: 输入文本
//   - maxLen: 最大长度
//   - defaultValue: 默认值
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
