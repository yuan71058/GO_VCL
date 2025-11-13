@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo   Windows GUI Application Build and UPX Compression Script
echo ========================================
echo.

:: Set variables
set PROJECT_NAME=windows-gui-app
set OUTPUT_DIR=dist
set OUTPUT_EXE=%PROJECT_NAME%.exe

:: Clean previous build
echo [1/6] Cleaning previous build files...
if exist %OUTPUT_DIR%\%OUTPUT_EXE% del /q %OUTPUT_DIR%\%OUTPUT_EXE%
if exist %OUTPUT_EXE% del /q %OUTPUT_EXE%
if exist defaultRes_windows_386.syso del /q defaultRes_windows_386.syso

:: Create output directory
echo [2/6] Creating output directory...
if not exist %OUTPUT_DIR% mkdir %OUTPUT_DIR%

:: Check windres tool and create resource file
echo [3/6] Checking windres tool and creating resource file...
where windres >nul 2>&1
if %errorlevel% equ 0 (
    echo Creating resource file...
    windres.exe -i app.rc -o defaultRes_windows_386.syso -F pe-i386
    if %errorlevel% neq 0 (
        echo Warning: Resource file creation failed, continuing compilation without custom icon
    )
) else (
    echo Warning: windres not found, skipping custom resource creation
    echo Tip: Install MinGW-w64 to enable custom icons and version information
)

:: Check Go environment
echo [4/6] Checking Go environment...
go version
if %errorlevel% neq 0 (
    echo Error: Go environment not properly configured
    pause
    exit /b 1
)

:: Download dependencies (if needed)
echo [5/6] Downloading dependency packages...
go mod download
if %errorlevel% neq 0 (
    echo Warning: Dependency download may have issues, but continuing compilation...
)

:: Compile to Windows GUI application (no console window)
echo [6/6] Compiling to Windows GUI application (no console window)...
go build -tags tempdll -ldflags "-w -s -H=windowsgui" -o %OUTPUT_DIR%\%OUTPUT_EXE%
if %errorlevel% neq 0 (
    echo Compilation failed!
    pause
    exit /b 1
)

:: Check compilation result
echo Checking compilation result...
if not exist %OUTPUT_DIR%\%OUTPUT_EXE% (
    echo Error: Executable file not generated
    pause
    exit /b 1
)

:: Display file information
echo.
echo ========================================
echo Compilation successful!
echo ========================================
echo Executable file: %OUTPUT_DIR%\%OUTPUT_EXE%
echo File size:
for %%I in (%OUTPUT_DIR%\%OUTPUT_EXE%) do echo   %%~zI bytes
echo.

:: Use UPX to compress executable file
echo ========================================
echo Using UPX to compress executable file...
upx --best --lzma %OUTPUT_DIR%\%OUTPUT_EXE% >nul 2>&1
if %errorlevel% equ 0 (
    echo UPX compression successful!
    
    :: Display compressed file size
    echo Compressed file size:
    for %%I in (%OUTPUT_DIR%\%OUTPUT_EXE%) do echo   %%~zI bytes
    
    echo.
    echo ========================================
    echo Compression complete!
    echo ========================================
) else (
    echo UPX compression failed or UPX not installed, skipping compression step
    echo Tip: You can download UPX from https://upx.github.io/ to reduce executable file size
)

echo.
echo Program is ready!
echo Executable file located at: %OUTPUT_DIR%\%OUTPUT_EXE%
echo.
pause