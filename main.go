// main.go - 主程序入口，Windows GUI应用
// 本文件是Windows GUI应用程序的入口点，负责应用程序的初始化、
// 主窗口创建和应用程序事件循环的启动
package main

import (
	"log"
	"runtime"
	"syscall"

	ui "windows-gui-app/ui"

	"github.com/ying32/govcl/vcl"
)

// Windows API相关变量
// 用于调用Windows系统函数来控制控制台窗口的显示和隐藏
var (
	// 加载kernel32.dll动态链接库，该库包含Windows系统核心API
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	// 获取控制台窗口句柄的函数指针
	getConsoleWindow = kernel32.NewProc("GetConsoleWindow")
	// 显示或隐藏窗口的函数指针
	showWindow       = kernel32.NewProc("ShowWindow")
)

// hideConsoleWindow 隐藏Windows控制台窗口
// 在Windows GUI应用程序中，通常不需要显示控制台窗口，
// 此函数通过调用Windows API将控制台窗口隐藏，提供更纯净的GUI体验
func hideConsoleWindow() {
	// 使用defer和recover捕获可能的异常，防止程序因隐藏窗口失败而崩溃
	defer func() {
		if r := recover(); r != nil {
			log.Printf("隐藏控制台窗口失败: %v", r)
		}
	}()

	// 仅在Windows操作系统下执行隐藏控制台窗口操作
	if runtime.GOOS == "windows" {
		// 获取当前控制台窗口的句柄
		consoleWindow, _, _ := getConsoleWindow.Call()
		// 如果成功获取到窗口句柄（非0值），则隐藏该窗口
		if consoleWindow != 0 {
			// SW_HIDE = 0，表示隐藏窗口
			// 忽略可能的错误继续运行，因为隐藏窗口失败不应影响程序正常功能
			showWindow.Call(consoleWindow, 0)
		}
	}
}

// main 应用程序主函数
// 负责整个应用程序的生命周期管理，包括初始化、窗口创建和事件循环
func main() {
	// 使用defer和recover捕获全局异常，防止程序意外崩溃
	defer func() {
		if r := recover(); r != nil {
			log.Printf("程序异常: %v", r)
		}
	}()

	// 在Windows环境下隐藏控制台窗口，提供纯GUI体验
	hideConsoleWindow()

	// 初始化VCL应用程序框架，这是使用govcl库的必要步骤
	vcl.Application.Initialize()

	// 设置应用程序在任务栏上显示主窗口图标
	// 这使得应用程序在任务栏中有更好的可见性和用户体验
	vcl.Application.SetMainFormOnTaskBar(true)
	
	// 设置应用程序标题，该标题将显示在窗口标题栏和任务栏中
	vcl.Application.SetTitle("Windows GUI应用")

	// 创建主窗口实例，主窗口是应用程序的主要用户界面
	// NewMainForm函数在ui包中定义，负责创建和初始化主窗口
	mainForm := ui.NewMainForm()

	// 显示主窗口，使其对用户可见
	// Show方法不仅显示窗口，还会触发OnShow事件
	mainForm.Show()

	// 启动应用程序消息循环，处理窗口事件和用户交互
	// 这是一个阻塞调用，直到应用程序退出才会返回
	vcl.Application.Run()
}