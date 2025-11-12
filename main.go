// main.go - 主程序入口，Windows GUI应用
package main

import (
	"log"
	"runtime"
	"syscall"

	ui "windows-gui-app/ui"

	"github.com/ying32/govcl/vcl"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	getConsoleWindow = kernel32.NewProc("GetConsoleWindow")
	showWindow       = kernel32.NewProc("ShowWindow")
)

// 隐藏Windows控制台窗口
func hideConsoleWindow() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("隐藏控制台窗口失败: %v", r)
		}
	}()

	if runtime.GOOS == "windows" {
		consoleWindow, _, _ := getConsoleWindow.Call()
		if consoleWindow != 0 {
			// SW_HIDE = 0，如果失败则忽略错误继续运行
			showWindow.Call(consoleWindow, 0)
		}
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("程序异常: %v", r)
		}
	}()

	// Windows下隐藏控制台窗口
	hideConsoleWindow()

	// 初始化应用程序
	vcl.Application.Initialize()

	// 设置应用程序属性
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.SetTitle("Windows GUI应用")

	// 创建主窗口
	mainForm := ui.NewMainForm()

	// 显示窗口
	mainForm.Show()

	// 运行应用程序
	vcl.Application.Run()
}
