@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo   Windows GUI应用 无窗口编译和运行脚本
echo ========================================
echo.

:: 设置变量
set PROJECT_NAME=windows-gui-app
set OUTPUT_DIR=dist
set OUTPUT_EXE=%PROJECT_NAME%.exe

:: 清理之前的构建
echo [1/7] 清理之前的构建文件...
if exist %OUTPUT_DIR%\%OUTPUT_EXE% del /q %OUTPUT_DIR%\%OUTPUT_EXE%
if exist %OUTPUT_EXE% del /q %OUTPUT_EXE%
if exist defaultRes_windows_386.syso del /q defaultRes_windows_386.syso

:: 创建输出目录
echo [2/7] 创建输出目录...
if not exist %OUTPUT_DIR% mkdir %OUTPUT_DIR%

:: 检查windres工具并创建资源文件
echo [3/7] 检查windres工具并创建资源文件...
where windres >nul 2>&1
if %errorlevel% equ 0 (
    echo 创建资源文件...
    windres.exe -i app.rc -o defaultRes_windows_386.syso -F pe-i386
    if %errorlevel% neq 0 (
        echo 警告: 资源文件创建失败，继续编译但不包含自定义图标
    )
) else (
    echo 警告: windres未找到，跳过自定义资源创建
    echo 提示: 安装MinGW-w64以启用自定义图标和版本信息
)

:: 检查Go环境
echo [4/7] 检查Go环境...
go version
if %errorlevel% neq 0 (
    echo 错误: Go环境未正确配置
    pause
    exit /b 1
)

:: 下载依赖（如果需要）
echo [5/7] 下载依赖包...
go mod download
if %errorlevel% neq 0 (
    echo 警告: 依赖下载可能有问题，但继续编译...
)

:: 编译为Windows GUI应用（无控制台窗口）
echo [6/7] 编译为Windows GUI应用（无控制台窗口）...
go build -tags tempdll -ldflags "-w -s -H=windowsgui" -o %OUTPUT_DIR%\%OUTPUT_EXE%
if %errorlevel% neq 0 (
    echo 编译失败!
    pause
    exit /b 1
)

:: 使用UPX压缩可执行文件
echo [6.5/7] 使用UPX压缩可执行文件...
upx --best --lzma %OUTPUT_DIR%\%OUTPUT_EXE% >nul 2>&1
if %errorlevel% equ 0 (
    echo UPX压缩成功!
) else (
    echo UPX压缩失败或未安装UPX，跳过压缩步骤
    echo 提示: 可以从 https://upx.github.io/ 下载UPX以减小可执行文件大小
)

:: 检查编译结果
echo [7/7] 检查编译结果...
if not exist %OUTPUT_DIR%\%OUTPUT_EXE% (
    echo 错误: 可执行文件未生成
    pause
    exit /b 1
)

:: 显示文件信息
echo.
echo ========================================
echo 编译成功! 
echo ========================================
echo 可执行文件: %OUTPUT_DIR%\%OUTPUT_EXE%
echo 文件大小: 
for %%I in (%OUTPUT_DIR%\%OUTPUT_EXE%) do echo   %%~zI 字节
echo.

:: 运行程序
echo 正在启动应用程序...
echo ========================================
start "" "%OUTPUT_DIR%\%OUTPUT_EXE%"

echo.
echo 程序已启动! 可以关闭此窗口。
echo 应用程序将出现在任务栏和桌面上。
echo.
:: pause