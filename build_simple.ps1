# Build and compress Windows GUI application with UPX

Write-Host "Building Windows GUI application with UPX compression..." -ForegroundColor Cyan

# Clean previous build
if (Test-Path "dist") {
    Remove-Item "dist" -Recurse -Force
}
if (Test-Path "windows-gui-app.exe") {
    Remove-Item "windows-gui-app.exe" -Force
}
if (Test-Path "defaultRes_windows_386.syso") {
    Remove-Item "defaultRes_windows_386.syso" -Force
}

# Create output directory
if (-not (Test-Path "dist")) {
    New-Item -ItemType Directory -Path "dist" | Out-Null
}

# Check if windres is available and create resource file
$windresAvailable = Get-Command windres -ErrorAction SilentlyContinue
if ($windresAvailable) {
    Write-Host "Creating resource file..." -ForegroundColor Yellow
    & windres.exe -i app.rc -o defaultRes_windows_386.syso -F pe-i386
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Warning: Resource file creation failed, continuing without custom resources" -ForegroundColor Yellow
    }
} else {
    Write-Host "Warning: windres not found, skipping custom resource creation" -ForegroundColor Yellow
    Write-Host "You can install MinGW-w64 to enable custom icons and version info" -ForegroundColor Gray
}

# Build the application
Write-Host "Building application..." -ForegroundColor Yellow
& go build -tags tempdll -ldflags "-w -s -H=windowsgui" -o "dist/windows-gui-app.exe"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    Read-Host "Press any key to exit"
    exit 1
}

# Check if UPX is available
$upxAvailable = Get-Command upx -ErrorAction SilentlyContinue
if (-not $upxAvailable) {
    Write-Host "UPX not found, skipping compression" -ForegroundColor Yellow
    Write-Host "You can download UPX from https://upx.github.io/" -ForegroundColor Gray
    Read-Host "Press any key to exit"
    exit 0
}

# Compress with UPX
Write-Host "Compressing with UPX..." -ForegroundColor Yellow
& upx --best --lzma "dist/windows-gui-app.exe"
if ($LASTEXITCODE -ne 0) {
    Write-Host "UPX compression failed" -ForegroundColor Red
    Read-Host "Press any key to exit"
    exit 1
}

# Display file sizes
$originalSize = 15394304  # Known original size
$compressedSize = (Get-Item "dist/windows-gui-app.exe").Length
$compressionRatio = [math]::Round((1 - $compressedSize / $originalSize) * 100, 2)

Write-Host "Build and compression completed successfully!" -ForegroundColor Green
Write-Host "Original size: $([math]::Round($originalSize / 1MB, 2)) MB" -ForegroundColor Cyan
Write-Host "Compressed size: $([math]::Round($compressedSize / 1MB, 2)) MB" -ForegroundColor Cyan
Write-Host "Compression ratio: $compressionRatio%" -ForegroundColor Green

Read-Host "Press any key to exit"