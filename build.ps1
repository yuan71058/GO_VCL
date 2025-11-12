# 编译和运行Windows GUI应用的简化脚本

Write-Host "开始编译Windows GUI应用..." -ForegroundColor Green

# 创建输出目录
if (-not (Test-Path "dist")) {
    New-Item -ItemType Directory -Path "dist" | Out-Null
}

# 编译应用
$buildResult = & go build -tags tempdll -ldflags "-w -s -H=windowsgui" -o "dist\windows-gui-app.exe"

if ($LASTEXITCODE -eq 0) {
    Write-Host "编译成功!" -ForegroundColor Green
    
    # 检查UPX是否可用
    $upxAvailable = Get-Command upx -ErrorAction SilentlyContinue
    if ($upxAvailable) {
        Write-Host "使用UPX压缩可执行文件..." -ForegroundColor Yellow
        
        # 获取原始文件大小
        $originalSize = (Get-Item "dist\windows-gui-app.exe").Length
        Write-Host "原始文件大小: $([math]::Round($originalSize / 1MB, 2)) MB" -ForegroundColor Cyan
        
        # 使用UPX压缩
        & upx --best --lzma "dist\windows-gui-app.exe"
        
        if ($LASTEXITCODE -eq 0) {
            # 获取压缩后文件大小
            $compressedSize = (Get-Item "dist\windows-gui-app.exe").Length
            $compressionRatio = [math]::Round((1 - $compressedSize / $originalSize) * 100, 2)
            Write-Host "压缩后文件大小: $([math]::Round($compressedSize / 1MB, 2)) MB" -ForegroundColor Cyan
            Write-Host "压缩率: $compressionRatio%" -ForegroundColor Green
        } else {
            Write-Host "UPX压缩失败，继续使用未压缩的文件" -ForegroundColor Yellow
        }
    } else {
        Write-Host "未找到UPX，跳过压缩步骤" -ForegroundColor Yellow
        Write-Host "提示: 可以从 https://upx.github.io/ 下载UPX以减小可执行文件大小" -ForegroundColor Gray
    }
    
    # 启动应用
    Write-Host "启动GUI应用程序..." -ForegroundColor Yellow
    Start-Process -FilePath "dist\windows-gui-app.exe" -WindowStyle Normal
    
    Write-Host "程序已启动!" -ForegroundColor Cyan
} else {
    Write-Host "编译失败!" -ForegroundColor Red
    exit 1
}