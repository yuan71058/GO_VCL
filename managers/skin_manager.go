// skin_manager.go - 皮肤管理器，负责加载和管理应用程序皮肤
// 本文件封装了SkinH_EL.dll的调用方法，提供皮肤加载和管理功能
package managers

import (
	"log"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

// SkinManager 皮肤管理器结构体
// 封装了皮肤加载和管理的相关方法和状态
type SkinManager struct {
	// SkinH_DLL 动态链接库句柄
	// 用于加载和调用SkinH_EL.dll中的函数
	SkinH_DLL *syscall.LazyDLL
	
	// 是否已加载皮肤标志
	// 跟踪皮肤加载状态，避免重复加载
	isLoaded bool
}

// NewSkinManager 创建新的皮肤管理器实例
// 返回初始化的皮肤管理器对象
func NewSkinManager() *SkinManager {
	return &SkinManager{
		isLoaded: false,
	}
}

// LoadSkinLibrary 加载皮肤库
// 加载SkinH_EL.dll动态链接库，初始化皮肤功能
// 返回是否加载成功
func (sm *SkinManager) LoadSkinLibrary() bool {
	// 使用defer和recover捕获可能的异常，防止程序因加载DLL失败而崩溃
	defer func() {
		if r := recover(); r != nil {
			log.Printf("加载皮肤库失败: %v", r)
		}
	}()

	// 加载SkinH_EL.dll动态链接库
	sm.SkinH_DLL = syscall.NewLazyDLL("SkinH_EL.dll")
	
	// 检查DLL是否加载成功
	if sm.SkinH_DLL == nil {
		log.Println("无法加载SkinH_EL.dll")
		return false
	}
	
	log.Println("SkinH_EL.dll加载成功")
	return true
}

// SkinH_Attach 加载皮肤文件
// 从指定路径加载皮肤文件，应用到当前应用程序
// 参数:
//   - skinPath: 皮肤文件的路径，如果为空则使用默认皮肤
// 返回是否加载成功
func (sm *SkinManager) SkinH_Attach(skinPath string) bool {
	// 检查DLL是否已加载
	if sm.SkinH_DLL == nil {
		log.Println("皮肤库未加载，请先调用LoadSkinLibrary")
		return false
	}

	// 获取SkinH_Attach函数指针
	proc := sm.SkinH_DLL.NewProc("SkinH_Attach")
	if proc == nil {
		log.Println("无法获取SkinH_Attach函数")
		return false
	}

	// 将皮肤路径转换为UTF16指针，Windows API需要
	skinPathPtr, err := syscall.UTF16PtrFromString(skinPath)
	if err != nil {
		log.Printf("皮肤路径转换失败: %v", err)
		return false
	}

	// 调用SkinH_Attach函数加载皮肤
	// 参数1: 皮肤文件路径指针
	// 参数2: 保留参数，通常为0
	ret, _, err := proc.Call(uintptr(unsafe.Pointer(skinPathPtr)), 0)
	
	// 检查返回值，非0表示成功
	if ret == 0 {
		log.Printf("加载皮肤失败: %v", err)
		return false
	}

	log.Printf("皮肤加载成功: %s", skinPath)
	sm.isLoaded = true
	return true
}

// SkinH_SetAero 启用Aero效果
// 启用Windows AERO效果，使界面更加美观
// 返回是否设置成功
func (sm *SkinManager) SkinH_SetAero() bool {
	// 检查DLL是否已加载
	if sm.SkinH_DLL == nil {
		log.Println("皮肤库未加载，请先调用LoadSkinLibrary")
		return false
	}

	// 获取SkinH_SetAero函数指针
	proc := sm.SkinH_DLL.NewProc("SkinH_SetAero")
	if proc == nil {
		log.Println("无法获取SkinH_SetAero函数")
		return false
	}

	// 调用SkinH_SetAero函数启用AERO效果
	// 参数: 保留参数，通常为0
	ret, _, err := proc.Call(0)
	
	// 检查返回值，非0表示成功
	if ret == 0 {
		log.Printf("启用AERO效果失败: %v", err)
		return false
	}

	log.Println("AERO效果启用成功")
	return true
}

// SkinH_AdjustAero 调整AERO效果参数
// 根据官方文档调整AERO效果的各种参数
// 参数:
//   - nAlpha: 透明度 (0-255, 默认值0)
//   - nShwDark: 亮度 (0-255, 默认值0)
//   - nShwSharp: 锐度 (0-12, 默认值0)
//   - nShwSize: 阴影大小 (0-18, 默认值0)
//   - nX: 水平偏移 (默认值0)
//   - nY: 垂直偏移 (默认值0)
//   - nRed: 红色分量 (0-255, 默认值0)
//   - nGreen: 绿色分量 (0-255, 默认值0)
//   - nBlue: 蓝色分量 (0-255, 默认值0)
// 返回是否设置成功
func (sm *SkinManager) SkinH_AdjustAero(nAlpha, nShwDark, nShwSharp, nShwSize, nX, nY, nRed, nGreen, nBlue int) bool {
	// 检查DLL是否已加载
	if sm.SkinH_DLL == nil {
		log.Println("皮肤库未加载，请先调用LoadSkinLibrary")
		return false
	}

	// 获取SkinH_AdjustAero函数指针
	proc := sm.SkinH_DLL.NewProc("SkinH_AdjustAero")
	if proc == nil {
		log.Println("无法获取SkinH_AdjustAero函数")
		return false
	}

	// 调用SkinH_AdjustAero函数调整AERO效果参数
	// 参数顺序: 透明度, 亮度, 锐度, 阴影大小, 水平偏移, 垂直偏移, 红色分量, 绿色分量, 蓝色分量
	ret, _, err := proc.Call(
		uintptr(nAlpha),
		uintptr(nShwDark),
		uintptr(nShwSharp),
		uintptr(nShwSize),
		uintptr(nX),
		uintptr(nY),
		uintptr(nRed),
		uintptr(nGreen),
		uintptr(nBlue),
	)
	
	// 检查返回值，非0表示成功
	if ret == 0 {
		log.Printf("调整AERO效果参数失败: %v", err)
		return false
	}

	log.Printf("AERO效果参数调整成功")
	return true
}

// LoadDefaultSkin 加载默认皮肤
// 加载应用程序默认的皮肤文件，通常在程序启动时调用
// 返回是否加载成功
func (sm *SkinManager) LoadDefaultSkin() bool {
	// 首先加载皮肤库
	if !sm.LoadSkinLibrary() {
		return false
	}

	// 获取当前执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("获取程序路径失败: %v", err)
		return false
	}

	// 构建默认皮肤文件路径
	// 假设皮肤文件与可执行文件在同一目录下，文件名为skinh.she
	skinPath := filepath.Join(filepath.Dir(exePath), "skinh.she")

	// 加载皮肤文件
	if !sm.SkinH_Attach(skinPath) {
		log.Println("加载默认皮肤失败，尝试使用内置皮肤")
		// 如果外部皮肤文件加载失败，尝试使用空字符串加载内置皮肤
		return sm.SkinH_Attach("")
	}

	// 启用AERO效果
	sm.SkinH_SetAero()

	// 调整AERO效果参数，使用推荐的默认值
	// 透明度: 200, 亮度: 0, 锐度: 3, 阴影大小: 5
	// 水平偏移: 0, 垂直偏移: 0, 颜色分量: 0
	sm.SkinH_AdjustAero(200, 0, 3, 5, 0, 0, 0, 0, 0)

	log.Println("默认皮肤加载完成")
	return true
}

// IsLoaded 检查皮肤是否已加载
// 返回皮肤加载状态
func (sm *SkinManager) IsLoaded() bool {
	return sm.isLoaded
}

// 全局皮肤管理器实例
var globalSkinManager = NewSkinManager()

// InitSkin 初始化皮肤系统
// 在程序启动时调用，加载默认皮肤
// 返回是否初始化成功
func InitSkin() bool {
	return globalSkinManager.LoadDefaultSkin()
}

// GetSkinManager 获取全局皮肤管理器实例
// 返回全局皮肤管理器对象
func GetSkinManager() *SkinManager {
	return globalSkinManager
}