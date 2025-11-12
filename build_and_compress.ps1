# 编译和UPX压缩Windows GUI应用的脚本

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Windows GUI应用 编译和UPX压缩脚本" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 设置变量
$ProjectName = "windows-gui-app"
$OutputDir = "dist"
$OutputExe = "$ProjectName.exe"

# 清理之前的构建
Write-Host "[1/5] 清理之前的构建文件..." -ForegroundColor Yellow
if (Test-Path $OutputDir) {
    Remove-Item $OutputDir -Recurse -Force
}
if (Test-Path $OutputExe) {
    Remove-Item $OutputExe -Force
}

# 创建输出目录
Write-Host "[2/5] 创建输出目录..." -ForegroundColor Yellow
if (-not (Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir | Out-Null
}

# 检查Go环境
Write-Host "[3/5] 检查Go环境..." -ForegroundColor Yellow
try {
    $goVersion = & go version
    Write-Host "  $goVersion" -ForegroundColor Green
} catch {
    Write-Host "错误: Go环境未正确配置" -ForegroundColor Red
    Read-Host "按任意键退出"
    exit 1
}

# 下载依赖（如果需要）
Write-Host "[4/5] 下载依赖包..." -ForegroundColor Yellow
& go mod download
if ($LASTEXITCODE -ne 0) {
    Write-Host "警告: 依赖下载可能有问题，但继续编译..." -ForegroundColor Yellow
}

# 编译为Windows GUI应用（无控制台窗口）
Write-Host "[5/5] 编译为Windows GUI应用（无控制台窗口）..." -ForegroundColor Yellow
& go build -tags tempdll -ldflags "-w -s -H=windowsgui" -o "$OutputDir\$OutputExe"
if ($LASTEXITCODE -ne 0) {
    Write-Host "编译失败!" -ForegroundColor Red
    Read-Host "按任意键退出"
    exit 1
}

# 检查编译结果
Write-Host "检查编译结果..." -ForegroundColor Yellow
if (-not (Test-Path "$OutputDir\$OutputExe")) {
    Write-Host "错误: 可执行文件未生成" -ForegroundColor Red
    Read-Host "按任意键退出"
    exit 1
}

# 显示文件信息
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "编译成功!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "可执行文件: $OutputDir\$OutputExe" -ForegroundColor White
$fileInfo = Get-Item "$OutputDir\$OutputExe"
Write-Host "文件大小: $($fileInfo.Length) 字节" -ForegroundColor White

# 使用UPX压缩可执行文件
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "使用UPX压缩可执行文件..." -ForegroundColor Yellow
$upxAvailable = Get-Command upx -ErrorAction SilentlyContinue
if ($upxAvailable) {
    # 获取原始文件大小
    $originalSize = (Get-Item "$OutputDir\$OutputExe").Length
    Write-Host "原始文件大小: $([math]::Round($originalSize / 1MB, 2)) MB" -ForegroundColor Cyan
    
    # 使用UPX压缩
    & upx --best --lzma "$OutputDir\$OutputExe"
    
    if ($LASTEXITCODE -eq 0) {
        # 获取压缩后文件大小
        $compressedSize = (Get-Item "$OutputDir\$OutputExe").Length
        $compressionRatio = [math]::Round((1 - $compressedSize / $originalSize) * 100, 2)
        Write-Host "压缩后文件大小: $([math]::Round($compressedSize / 1MB, 2)) MB" -ForegroundColor Cyan
        Write-Host "压缩率: $compressionRatio%" -ForegroundColor Green
        
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host "压缩完成!" -ForegroundColor Green
        Write-Host "========================================" -ForegroundColor Cyan
    } else {
        Write-Host "UPX压缩失败，继续使用未压缩的文件" -ForegroundColor Yellow
    }
} else {
    Write-Host "未找到UPX，跳过压缩步骤" -ForegroundColor Yellow
    Write-Host "提示: 可以从 https://upx.github.io/ 下载UPX以减小可执行文件大小" -ForegroundColor Gray
}

Write-Host ""
Write-Host "程序已准备就绪!" -ForegroundColor Green
Write-Host "可执行文件位于: $OutputDir\$OutputExe" -ForegroundColor White
Write-Host ""
Read-Host "Press any key to exit"