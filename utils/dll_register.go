package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

// DLLRegister DLL注册相关结构体
type DLLRegister struct {
	dllPath    string
	registered bool
}

// NewDLLRegister 创建DLL注册对象
func NewDLLRegister() *DLLRegister {
	return &DLLRegister{
		registered: false,
	}
}

// RegisterDLL 注册大漠插件DLL
func (dr *DLLRegister) RegisterDLL(dllPath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(dllPath); os.IsNotExist(err) {
		return fmt.Errorf("DLL文件不存在: %s", dllPath)
	}
	
	// 获取绝对路径
	absPath, err := filepath.Abs(dllPath)
	if err != nil {
		return fmt.Errorf("获取DLL绝对路径失败: %v", err)
	}
	
	// 检查是否已注册
	if dr.IsRegistered(absPath) {
		dr.dllPath = absPath
		dr.registered = true
		return nil
	}
	
	// 使用regsvr32注册DLL
	cmd := exec.Command("regsvr32", "/s", absPath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("注册DLL失败: %v", err)
	}
	
	dr.dllPath = absPath
	dr.registered = true
	return nil
}

// UnregisterDLL 注销大漠插件DLL
func (dr *DLLRegister) UnregisterDLL(dllPath string) error {
	// 获取绝对路径
	_, err := filepath.Abs(dllPath)
	if err != nil {
		return fmt.Errorf("获取DLL绝对路径失败: %v", err)
	}
	
	// 使用regsvr32 /u注销DLL
	cmd := exec.Command("regsvr32", "/s", "/u", dllPath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("注销DLL失败: %v", err)
	}
	
	// 如果是当前注册的DLL，更新状态
	if dr.dllPath == dllPath {
		dr.registered = false
		dr.dllPath = ""
	}
	
	return nil
}

// IsRegistered 检查DLL是否已注册
func (dr *DLLRegister) IsRegistered(dllPath string) bool {
	// 检查注册表
	keyPath := `SOFTWARE\Classes\TypeLib\`
	
	// 打开注册表键
	key, err := openRegistryKey(syscall.HKEY_CLASSES_ROOT, keyPath)
	if err != nil {
		return false
	}
	defer closeRegistryKey(key)
	
	// 枚举子键，查找大漠插件相关的注册信息
	var subKeyNames []string
	err = enumRegistryKeys(key, &subKeyNames)
	if err != nil {
		return false
	}
	
	// 检查是否包含大漠插件的注册信息
	for _, subKeyName := range subKeyNames {
		subKey, err := openRegistryKey(key, subKeyName)
		if err != nil {
			continue
		}
		
		// 获取默认值
		var value string
		err = queryRegistryValue(subKey, "", &value)
		closeRegistryKey(subKey)
		
		if err == nil && strings.Contains(value, "dm") {
			return true
		}
	}
	
	return false
}

// GetDLLPath 获取当前注册的DLL路径
func (dr *DLLRegister) GetDLLPath() string {
	return dr.dllPath
}

// IsCurrentDLLRegistered 检查当前DLL是否已注册
func (dr *DLLRegister) IsCurrentDLLRegistered() bool {
	return dr.registered
}

// CheckDLLExists 检查DLL文件是否存在
func CheckDLLExists(dllPath string) bool {
	_, err := os.Stat(dllPath)
	return !os.IsNotExist(err)
}

// CheckDLLArchitecture 检查DLL架构是否匹配
func CheckDLLArchitecture(dllPath string) (bool, error) {
	file, err := os.Open(dllPath)
	if err != nil {
		return false, fmt.Errorf("打开DLL文件失败: %v", err)
	}
	defer file.Close()
	
	// 读取DOS头
	var dosHeader [64]byte
	_, err = file.Read(dosHeader[:])
	if err != nil {
		return false, fmt.Errorf("读取DOS头失败: %v", err)
	}
	
	// 检查MZ签名
	if string(dosHeader[0:2]) != "MZ" {
		return false, fmt.Errorf("无效的DLL文件")
	}
	
	// 获取PE头偏移
	peOffset := int(dosHeader[60]) | int(dosHeader[61])<<8 | int(dosHeader[62])<<16 | int(dosHeader[63])<<24
	
	// 读取PE头
	_, err = file.Seek(int64(peOffset), 0)
	if err != nil {
		return false, fmt.Errorf("定位PE头失败: %v", err)
	}
	
	var peHeader [24]byte
	_, err = file.Read(peHeader[:])
	if err != nil {
		return false, fmt.Errorf("读取PE头失败: %v", err)
	}
	
	// 检查PE签名
	if string(peHeader[0:2]) != "PE" {
		return false, fmt.Errorf("无效的PE文件")
	}
	
	// 获取机器类型
	machine := uint16(peHeader[4]) | uint16(peHeader[5])<<8
	
	// 检查是否为x86或x64
	isX86 := machine == 0x014c // IMAGE_FILE_MACHINE_I386
	isX64 := machine == 0x8664 // IMAGE_FILE_MACHINE_AMD64
	
	// 获取当前进程架构 - 使用更简单的方法
	// 在Windows上，我们可以通过环境变量来判断
	currentIsX64 := os.Getenv("PROCESSOR_ARCHITECTURE") == "AMD64" ||
		(os.Getenv("PROCESSOR_ARCHITEW6432") != "" && os.Getenv("PROCESSOR_ARCHITEW6432") == "AMD64")
	
	// 检查架构是否匹配
	if (currentIsX64 && isX86) || (!currentIsX64 && isX64) {
		return false, fmt.Errorf("DLL架构与当前进程不匹配")
	}
	
	return true, nil
}

// CheckDLLDependencies 检查DLL依赖项
func CheckDLLDependencies(dllPath string) ([]string, error) {
	// 这里应该实现DLL依赖项检查
	// 由于Go标准库没有直接支持，这里提供一个简化实现
	
	// 常见的大漠插件依赖项
	commonDeps := []string{
		"msvcr120.dll",
		"msvcp120.dll",
		"kernel32.dll",
		"user32.dll",
		"gdi32.dll",
		"ole32.dll",
		"oleaut32.dll",
	}
	
	var missingDeps []string
	
	for _, dep := range commonDeps {
		// 检查系统目录中是否存在依赖项
		systemDir := getSystemDirectory()
		depPath := filepath.Join(systemDir, dep)
		
		if !CheckDLLExists(depPath) {
			missingDeps = append(missingDeps, dep)
		}
	}
	
	return missingDeps, nil
}

// getSystemDirectory 获取系统目录
func getSystemDirectory() string {
	var buffer [256]uint16
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getSystemDirectory := kernel32.NewProc("GetSystemDirectoryW")
	ret, _, _ := getSystemDirectory.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
	)
	n := uint32(ret)
	if n == 0 {
		return "C:\\Windows\\System32"
	}
	return syscall.UTF16ToString(buffer[:n])
}

// 注册表相关函数
var (
	modadvapi32 = syscall.NewLazyDLL("advapi32.dll")
	procRegOpenKeyExW = modadvapi32.NewProc("RegOpenKeyExW")
	procRegCloseKey = modadvapi32.NewProc("RegCloseKey")
	procRegEnumKeyExW = modadvapi32.NewProc("RegEnumKeyExW")
	procRegQueryValueExW = modadvapi32.NewProc("RegQueryValueExW")
)

const (
	KEY_READ = 0x20019
	KEY_ENUMERATE_SUB_KEYS = 0x8
)

// openRegistryKey 打开注册表键
func openRegistryKey(key syscall.Handle, subKey string) (syscall.Handle, error) {
	var k syscall.Handle
	subKeyPtr, err := syscall.UTF16PtrFromString(subKey)
	if err != nil {
		return 0, err
	}
	
	ret, _, err := procRegOpenKeyExW.Call(
		uintptr(key),
		uintptr(unsafe.Pointer(subKeyPtr)),
		0,
		KEY_READ|KEY_ENUMERATE_SUB_KEYS,
		uintptr(unsafe.Pointer(&k)),
	)
	
	if ret != 0 {
		return 0, err
	}
	
	return k, nil
}

// closeRegistryKey 关闭注册表键
func closeRegistryKey(key syscall.Handle) error {
	ret, _, err := procRegCloseKey.Call(uintptr(key))
	if ret != 0 {
		return err
	}
	return nil
}

// enumRegistryKeys 枚举注册表键
func enumRegistryKeys(key syscall.Handle, names *[]string) error {
	var index uint32 = 0
	var name [256]uint16
	var nameLen uint32 = 256
	
	for {
		ret, _, err := procRegEnumKeyExW.Call(
			uintptr(key),
			uintptr(index),
			uintptr(unsafe.Pointer(&name[0])),
			uintptr(unsafe.Pointer(&nameLen)),
			uintptr(0),
			uintptr(0),
			uintptr(0),
			uintptr(0),
		)
		
		if ret != 0 {
			if ret == 259 { // ERROR_NO_MORE_ITEMS
				break
			}
			return err
		}
		
		nameStr := syscall.UTF16ToString(name[:nameLen])
		*names = append(*names, nameStr)
		
		index++
		nameLen = 256
	}
	
	return nil
}

// queryRegistryValue 查询注册表值
func queryRegistryValue(key syscall.Handle, valueName string, value *string) error {
	valueNamePtr, err := syscall.UTF16PtrFromString(valueName)
	if err != nil {
		return err
	}
	
	var valueType uint32
	var data [256]byte
	var dataLen uint32 = 256
	
	ret, _, err := procRegQueryValueExW.Call(
		uintptr(key),
		uintptr(unsafe.Pointer(valueNamePtr)),
		uintptr(0),
		uintptr(unsafe.Pointer(&valueType)),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(unsafe.Pointer(&dataLen)),
	)
	
	if ret != 0 {
		return err
	}
	
	if valueType == 1 { // REG_SZ
		*value = syscall.UTF16ToString((*[256]uint16)(unsafe.Pointer(&data[0]))[:dataLen/2])
	}
	
	return nil
}