# 构建脚本说明

本项目提供了多种构建脚本，满足不同场景下的构建需求。

## 构建脚本列表

### 1. build_and_run.bat
- **功能**: 编译并立即运行程序
- **使用场景**: 开发测试阶段，快速验证代码更改
- **命令**: 
  ```batch
  go build -tags tempdll -ldflags "-w -s -H=windowsgui" -o dist/windows-gui-app.exe
  dist/windows-gui-app.exe
  ```

### 2. build_and_compress.bat
- **功能**: 编译并使用UPX压缩程序
- **使用场景**: 发布准备，减小可执行文件体积
- **命令**:
  ```batch
  go build -tags tempdll -ldflags "-w -s -H=windowsgui" -o dist/windows-gui-app.exe
  upx --best --lzma dist/windows-gui-app.exe
  ```

### 3. build.ps1
- **功能**: PowerShell版本的构建脚本
- **使用场景**: PowerShell环境下构建
- **命令**:
  ```powershell
  go build -tags tempdll -ldflags "-w -s -H=windowsgui" -o dist/windows-gui-app.exe
  ```

### 4. build_and_run.ps1
- **功能**: PowerShell版本的构建运行脚本
- **使用场景**: PowerShell环境下开发测试
- **命令**:
  ```powershell
  go build -tags tempdll -ldflags "-w -s -H=windowsgui" -o dist/windows-gui-app.exe
  & .\dist\windows-gui-app.exe
  ```

## 构建参数说明

### 编译参数
- `-tags tempdll`: 嵌入liblcl.dll到可执行文件中
- `-ldflags "-w -s -H=windowsgui"`: 
  - `-w`: 去除调试信息
  - `-s`: 去除符号表
  - `-H=windowsgui`: 设置为Windows GUI应用程序

### UPX压缩参数
- `--best`: 使用最佳压缩算法
- `--lzma`: 使用LZMA压缩算法

## 构建环境要求

- **Go版本**: 1.22或更高版本
- **操作系统**: Windows 10/11
- **UPX工具**: 使用压缩脚本时需要安装UPX

## 构建输出

所有构建脚本都会在`dist`目录下生成`windows-gui-app.exe`可执行文件。

- **未压缩大小**: 约15.5MB
- **压缩后大小**: 约5.1MB
- **压缩率**: 67.9%

## 注意事项

1. 首次构建前请确保已执行`go mod tidy`安装依赖
2. 构建过程中会自动创建`dist`目录
3. 压缩构建需要额外安装UPX工具
4. 构建前请确保没有其他程序占用输出文件

## 故障排除

### 构建失败
- 检查Go版本是否符合要求
- 确认网络连接正常，能够下载依赖
- 检查磁盘空间是否充足

### UPX压缩失败
- 确认已安装UPX并添加到PATH环境变量
- 尝试使用不带压缩的构建脚本
- 检查UPX版本是否兼容当前系统

### 运行失败
- 确认目标系统是否为Windows
- 检查是否有足够的系统权限
- 尝试以管理员身份运行