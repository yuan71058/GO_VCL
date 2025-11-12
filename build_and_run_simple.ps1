# 设置变量
$AppName = "windows-gui-app"
$DistDir = "dist"
$Executable = "$AppName.exe"
$ExecutablePath = "$DistDir\$Executable"

Write-Host "Starting build and run Windows GUI application..." -ForegroundColor Green

# 1. Clean old build files
Write-Host "Step 1: Cleaning old build files..."
if (Test-Path $DistDir) {
    Remove-Item -Recurse -Force $DistDir
    Write-Host "Deleted old dist directory"
}

# 2. Create build directory
Write-Host "Step 2: Creating build directory..."
New-Item -ItemType Directory -Force -Path $DistDir | Out-Null
Write-Host "Created dist directory"

# 3. Check Go environment
Write-Host "Step 3: Checking Go environment..."
try {
    $goVersion = go version
    Write-Host "Go environment check passed: $goVersion"
} catch {
    Write-Host "Error: Go environment not found, please install Go first" -ForegroundColor Red
    exit 1
}

# 4. Download dependencies
Write-Host "Step 4: Downloading dependencies..."
try {
    go mod tidy
    Write-Host "Dependencies downloaded"
} catch {
    Write-Host "Error: Failed to download dependencies" -ForegroundColor Red
    exit 1
}

# 5. Compile Windows GUI application
Write-Host "Step 5: Compiling Windows GUI application..."
try {
    go build -tags tempdll -ldflags "-w -s -H=windowsgui" -o $ExecutablePath
    Write-Host "Compilation successful"
} catch {
    Write-Host "Error: Compilation failed" -ForegroundColor Red
    exit 1
}

# 6. Check compilation result
if (-not (Test-Path $ExecutablePath)) {
    Write-Host "Error: Cannot find the compiled executable file" -ForegroundColor Red
    exit 1
}

# 7. Display file information
$fileInfo = Get-Item $ExecutablePath
Write-Host "File size: $($fileInfo.Length) bytes"
Write-Host "Creation time: $($fileInfo.CreationTime)"

# 8. Compress executable with UPX
Write-Host "Step 8: Compressing executable with UPX..."
try {
    # Check if UPX is available
    $upxVersion = upx --version 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "UPX version: $upxVersion"
        
        # Get original file size
        $originalSize = $fileInfo.Length
        
        # Perform compression
        upx --best --lzma $ExecutablePath
        
        # Get compressed file size
        $compressedFileInfo = Get-Item $ExecutablePath
        $compressedSize = $compressedFileInfo.Length
        
        # Calculate compression ratio
        $compressionRatio = [math]::Round(($compressedSize / $originalSize) * 100, 2)
        
        Write-Host "Compression complete: $originalSize -> $compressedSize bytes (compression ratio: $compressionRatio%)"
    } else {
        Write-Host "Warning: UPX not found, skipping compression step" -ForegroundColor Yellow
    }
} catch {
    Write-Host "Warning: UPX compression failed, continuing execution" -ForegroundColor Yellow
}

# 9. Launch application
Write-Host "Step 9: Launching application..."
try {
    Start-Process $ExecutablePath -WindowStyle Normal
    Write-Host "Application launched"
} catch {
    Write-Host "Error: Failed to launch application" -ForegroundColor Red
    exit 1
}

Write-Host "Build and run complete!" -ForegroundColor Green