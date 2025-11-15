# BAT批处理文件格式说明

## 格式要求

本项目中的所有BAT批处理文件必须使用**PC格式**（Windows格式），而不是UNIX格式。

### PC格式与UNIX格式的区别

- **PC格式**：使用回车符(CR)和换行符(LF)作为行结束符，即`\r\n`
- **UNIX格式**：仅使用换行符(LF)作为行结束符，即`\n`

### 为什么必须使用PC格式

1. **兼容性**：Windows命令解释器(cmd.exe)期望使用PC格式的批处理文件
2. **可靠性**：UNIX格式的批处理文件在Windows上可能导致执行错误
3. **一致性**：保持项目中所有批处理文件格式统一

### 如何确保使用PC格式

#### 方法1：使用提供的转换工具

项目提供了`convert_bat_to_pc.bat`工具，可以将UNIX格式的bat文件转换为PC格式：

```cmd
convert_bat_to_pc.bat script.bat
```

#### 方法2：使用文本编辑器

在创建或编辑bat文件时，确保文本编辑器设置为Windows/PC格式：
- Notepad++：编辑->文档格式转换->转换为Windows格式
- VS Code：点击右下角的CRLF/LF，确保选择CRLF
- Windows记事本：默认保存为PC格式

#### 方法3：使用PowerShell命令

```powershell
# 将UNIX格式转换为PC格式
(Get-Content 'script.bat') | Set-Content -NoNewline 'script.bat'
```

### 验证文件格式

可以使用以下方法验证bat文件是否为PC格式：

#### 使用PowerShell

```powershell
# 检查文件格式
Get-Content 'script.bat' -Raw | Select-String "`r`n"
```

#### 使用十六进制编辑器

- PC格式的行结束符显示为：`0D 0A`
- UNIX格式的行结束符显示为：`0A`

### 项目中的批处理文件

当前项目中的批处理文件：
- `build_and_run.bat` - 构建并运行应用程序
- `build_and_compress.bat` - 构建并压缩应用程序
- `convert_bat_to_pc.bat` - 将UNIX格式bat文件转换为PC格式

所有这些文件都应该使用PC格式。

### 注意事项

1. 从Git仓库克隆的文件可能默认为UNIX格式，需要转换
2. 在跨平台开发环境中，特别注意保持bat文件的PC格式
3. 提交代码前，确保所有bat文件都是PC格式

### 相关资源

- [Windows批处理文件文档](https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/windows-commands)
- [行结束符说明](https://en.wikipedia.org/wiki/Newline)