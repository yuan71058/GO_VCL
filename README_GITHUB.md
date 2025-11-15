# GO_VCL Windows GUI应用程序

[![Release Version](https://img.shields.io/github/release/yuan71058/GO_VCL.svg)](https://github.com/yuan71058/GO_VCL/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Windows-lightgrey.svg)](https://www.microsoft.com/windows)
[![Version](https://img.shields.io/badge/Version-v1.3.2-red.svg)](https://github.com/yuan71058/GO_VCL/releases)

GO_VCL是一个功能丰富的Windows桌面GUI应用程序，基于Go语言和Govcl UI框架开发，集成了Excel处理、JSON操作、HTTP请求和数据库管理等多种功能，提供直观易用的图形界面。

## 📸 功能展示

### 主界面
![应用程序界面截图](image.png)

### 应用图标
![应用图标](rgb.ico)

## 🌟 主要功能

### 📊 Excel操作
- **文件导入导出**: 支持xlsx格式的Excel文件导入和导出
- **数据处理**: 提供完整的Excel数据读取、写入和修改功能
- **批量操作**: 支持批量数据处理和格式转换

### 🔧 JSON操作
- **数据解析**: 完整的JSON数据解析和验证功能
- **文件操作**: JSON文件的读取、写入和格式化
- **路径查询**: 支持JSONPath查询和数据提取

### 🌐 HTTP网络操作
- **多方法支持**: GET、POST、PUT、DELETE等HTTP方法
- **文件下载**: 支持文件下载和进度显示
- **请求管理**: 自定义请求头、参数和超时设置

### 💾 数据库操作
- **SQLite管理**: 完整的SQLite数据库操作功能
- **CRUD操作**: 创建、读取、更新、删除数据记录
- **模拟模式**: 在无SQLite驱动环境下的模拟操作

### 🖼️ 图像处理
- **图片导入**: 支持多种图片格式的导入和显示
- **图像显示**: 在应用程序界面中显示导入的图片
- **边框效果**: 为图片框添加美观的边框效果

### 🎛️ UI组件
- **单选框**: 提供选项配置功能
- **多选框**: 支持多选项配置
- **菜单栏**: 包含文件、编辑、工具和帮助菜单
- **状态栏**: 实时显示操作状态和结果

### 🌐 TCP服务器
- **TCP通信**: 完整的TCP服务器功能，支持多客户端连接
- **消息处理**: 支持客户端消息接收和广播
- **UTF-8编码**: 解决了服务器重启后出现乱码的问题
- **连接管理**: 实时显示客户端连接状态和信息

### 🌐 Web服务器
- **HTTP服务**: 提供RESTful API接口，支持状态查询和客户端管理
- **WebSocket服务**: 支持实时双向通信，适用于实时数据推送
- **客户端管理**: 跟踪和管理连接的客户端
- **API路由**: 内置多个API端点，包括状态查询、客户端列表等

### ⚡ 高级特性
- **多线程支持**: 并发任务管理和结果收集
- **独立部署**: 通过DLL嵌入技术实现无需外部依赖
- **压缩优化**: 使用UPX技术减小可执行文件大小约67.9%

## 📁 项目结构

```
GO_VCL/
├── main.go                      # 程序入口文件
├── MainForm.go                  # 主窗口界面实现
├── ui.go                        # UI接口定义
├── app.manifest                 # 应用程序清单
├── app.rc                       # 资源文件
├── rgb.ico                      # 应用程序图标
├── defaultRes_windows_386.syso  # 默认资源文件
├── dist/                        # 构建输出目录
│   ├── windows-gui-app.exe      # 可执行文件
│   ├── SkinH_EL.dll             # 皮肤引擎库
│   ├── liblcl.dll               # VCL运行时库
│   ├── Aero.she                 # 皮肤文件
│   ├── MACOS-白色.she           # 皮肤文件
│   ├── QQ2008.she               # 皮肤文件
│   └── skinh.she                # 皮肤文件
├── managers/                    # 功能管理器模块
│   ├── excel_manager.go         # Excel操作管理器
│   ├── json_manager.go          # JSON操作管理器
│   ├── http_manager.go          # HTTP操作管理器
│   ├── database_manager.go      # 数据库操作管理器
│   ├── concurrent_manager.go    # 并发任务管理器
│   ├── tcp_server_manager.go    # TCP服务器管理器
│   ├── webserver_manager.go     # Web服务器管理器
│   └── skin_manager.go          # 皮肤管理器
├── test/                        # 测试文件目录
│   └── test_tcp_client.go       # TCP客户端测试
├── docs/                        # 文档目录
│   ├── USER_MANUAL.md           # 用户手册
│   ├── ICON_AND_MANIFEST_GUIDE.md # 图标和清单指南
│   ├── IMPORT_EXCEL_FEATURE.md  # Excel导入功能说明
│   ├── TCP_SERVER_FEATURE.md    # TCP服务器功能说明
│   ├── WEBSERVER_FEATURE.md     # Web服务器功能说明
│   ├── IMAGE_GUIDE.md           # 图片添加指南
│   ├── BAT_FORMAT_GUIDE.md      # 批处理脚本格式指南
│   └── UPX_COMPRESSION_FEATURE.md # UPX压缩功能说明
├── 说明文档.md                   # 项目详细说明文档
├── README.md                    # 项目说明文档
├── BUILD_SCRIPTS_README.md      # 构建脚本说明
├── build_and_run.bat            # 编译并运行脚本
└── build_and_compress.bat       # 编译并压缩脚本
```

## 🚀 快速开始

### 系统要求
- **操作系统**: Windows 10/11 (x64)
- **运行环境**: 无需额外依赖（已嵌入所有必需组件）
- **内存**: 最少512MB可用内存
- **磁盘**: 至少20MB可用空间

### 下载安装

1. **访问发布页面**:
   - 打开浏览器访问: https://github.com/yuan71058/GO_VCL/releases
   - 下载最新版本的`windows-gui-app.exe`文件

2. **运行程序**:
   - 将下载的文件放置到任意目录
   - 双击`windows-gui-app.exe`启动程序
   - 无需安装，即下即用

### 从源码构建

如果您想从源码构建程序，请按照以下步骤操作：

1. **环境准备**:
   ```bash
   # 安装Go 1.20或更高版本
   # 确保已配置好GOPATH和PATH环境变量
   ```

2. **获取源码**:
   ```bash
   git clone https://github.com/yuan71058/GO_VCL.git
   cd GO_VCL
   ```

3. **安装依赖**:
   ```bash
   go mod tidy
   ```

4. **构建程序**:
   ```bash
   # 使用提供的构建脚本
   .\build_and_run.bat
   
   # 或手动构建
   go build -tags tempdll -ldflags "-w -s -H=windowsgui" -o dist/windows-gui-app.exe
   ```

## 📖 使用指南

### 主界面介绍

程序启动后，您将看到包含以下组件的主界面：

- **顶部菜单栏**: 文件、编辑、工具、帮助菜单
- **左侧按钮区**: 各功能模块的操作按钮
- **中央表格区**: 数据显示和编辑区域
- **右侧图片区**: 图片导入和显示区域
- **底部配置区**: 单选框、多选框等配置选项
- **底部状态栏**: 操作状态和结果信息

### 功能操作说明

#### Excel操作
1. 点击"📊 Excel"按钮进行Excel文件操作
2. 使用"📥 导入Excel"按钮选择并导入xlsx文件
3. 数据将自动显示在中央表格区域
4. 支持数据的编辑和导出功能

#### JSON操作
1. 点击"🔧 JSON"按钮进行JSON数据处理
2. 在弹出的对话框中输入或粘贴JSON数据
3. 程序将自动解析并验证JSON格式
4. 支持JSON数据的格式化和保存

#### HTTP请求
1. 点击"🌐 HTTP"按钮进行网络请求
2. 输入目标URL和选择请求方法
3. 可添加自定义请求头和参数
4. 查看服务器响应和状态码

#### 数据库操作
1. 点击"💾 数据库"按钮进行数据库操作
2. 程序将自动创建或连接SQLite数据库
3. 执行各种CRUD操作（创建、读取、更新、删除）
4. 查询结果将显示在表格区域

#### 图片导入
1. 点击"🖼️ 导入图片"按钮选择图片文件
2. 支持常见图片格式（JPG、PNG、BMP等）
3. 图片将显示在右侧图片框中
4. 图片框带有美观的边框效果

#### TCP服务器
1. 点击"🔌 TCP服务"按钮启动TCP服务器
2. 服务器默认监听端口8082
3. 使用"📡 连接TCP"按钮连接到服务器
4. 使用"📤 发送数据"按钮向服务器发送消息
5. 服务器支持多客户端同时连接
6. 消息采用UTF-8编码，避免了乱码问题

#### Web服务器
1. 点击"🌐 Web服务器"按钮启动Web服务器
2. HTTP服务默认监听端口8080
3. WebSocket服务默认监听端口8081
4. 服务器提供RESTful API接口
5. 支持实时双向WebSocket通信
6. 可通过API查询服务器状态和客户端信息

#### 获取选中状态
1. 点击"📋 获取选中状态"按钮
2. 程序将收集所有单选框和多选框的状态
3. 结果将显示在底部状态栏中

#### 清空编辑框
1. 通过菜单栏"编辑" → "清空编辑框"
2. 或使用快捷键（如果已定义）
3. 所有编辑框内容将被清空

## 🛠️ 技术架构

### 架构设计
- **分层架构**: UI层、业务逻辑层、数据访问层清晰分离
- **模块化设计**: 各功能模块独立开发和维护，降低耦合度
- **事件驱动**: 基于事件处理机制的用户交互模式，直接在MainForm结构体中实现所有事件处理逻辑
- **管理器模式**: 使用专门的管理器类处理各功能模块

### 核心技术
- **Go语言**: 高性能、类型安全的编程语言
- **Govcl框架**: 跨平台的VCL GUI框架
- **SQLite**: 轻量级嵌入式数据库
- **Excelize**: 高性能Excel文件处理库
- **GJSON**: 高效JSON解析库
- **HttpRequest**: 简洁易用的HTTP客户端库

### 代码质量
- **详细注释**: 所有函数和结构体都有详细的中文注释
- **错误处理**: 完善的错误处理和异常捕获机制
- **类型安全**: 充分利用Go语言的类型系统
- **内存管理**: 自动的内存管理和垃圾回收
- **编码规范**: 遵循Go语言官方编码规范

## 🔧 高级配置

### 构建选项

项目提供了多种构建脚本，满足不同需求：

1. **build_and_run.bat**: 编译并立即运行程序
2. **build_and_compress.bat**: 编译并使用UPX压缩程序
3. **build.ps1**: PowerShell版本的构建脚本
4. **build_and_run.ps1**: PowerShell版本的构建运行脚本

### 自定义配置

您可以通过修改以下文件来自定义程序行为：

- **app.manifest**: 应用程序清单，包含UAC和兼容性设置
- **app.rc**: 资源文件，可自定义图标和版本信息
- **rgb.ico**: 应用程序图标，可替换为自定义图标

### 压缩优化

程序使用UPX压缩技术减小可执行文件大小：

- **原始大小**: 约15.5MB
- **压缩后大小**: 约5.1MB
- **压缩率**: 67.9%

压缩后的程序功能完全相同，但文件大小显著减小，便于分发和下载。

## 🐛 故障排除

### 常见问题

1. **程序无法启动**
   - 检查Windows版本是否为Windows 10/11
   - 确认下载的文件完整，没有损坏
   - 尝试以管理员身份运行

2. **功能操作失败**
   - 检查文件路径是否包含特殊字符
   - 确认目标文件没有被其他程序占用
   - 查看状态栏的错误信息提示

3. **Excel操作问题**
   - 确认Excel文件格式为xlsx
   - 检查文件是否被其他程序打开
   - 尝试使用较小的Excel文件测试

4. **HTTP请求失败**
   - 检查网络连接是否正常
   - 确认URL地址是否正确
   - 某些网站可能不允许程序化访问

### 日志和调试

程序内置了详细的日志输出功能：

1. **状态栏日志**: 操作结果和状态信息实时显示
2. **错误提示**: 友好的错误信息和解决建议
3. **操作反馈**: 每个操作都有明确的成功/失败反馈

### 性能优化

- **大文件处理**: 对于大型Excel文件，处理可能需要较长时间
- **内存使用**: 程序会自动管理内存，长时间使用后建议重启
- **并发限制**: 多线程操作有数量限制，避免同时执行过多任务

## 📋 版本历史

### v1.3.2 (2025-11-14)
- **TCP服务器功能增强**：修复TCP服务器停止后再次启动出现乱码的问题
- **UTF-8编码支持**：为TCP服务器添加UTF-8编码处理，确保中文字符正确传输
- **Web服务器功能**：新增Web服务器功能，支持HTTP和WebSocket服务
- **API接口完善**：提供RESTful API接口，支持状态查询和客户端管理
- **文档更新**：添加TCP服务器和Web服务器功能文档，完善用户指南

### v1.3.1 (2025-11-12)
- 项目结构优化：清理无用文件和脚本，删除重复文件
- 皮肤文件管理：将皮肤文件迁移到dist目录，统一管理
- 批处理脚本优化：保留最终版本，删除冗余脚本
- 文档更新：更新项目说明文档，记录最新清理工作

### v1.3.0 (2025-11-12)
- **代码清理与优化**
  - 删除未使用的handlers/button_handlers.go文件（675行代码）
  - 将所有事件处理逻辑整合到MainForm结构体中
  - 简化项目结构，提高代码内聚性
  - 更新项目文档，保持文档与代码同步
- **架构优化**
  - 采用更直接的事件处理方式
  - 减少不必要的抽象层
  - 提高代码可维护性

### v1.2.0 (2025-11-12)
- **界面优化与美化**
  - 优化界面布局和控件排列
  - 添加图标和美化界面元素
  - 改进用户体验
- **功能增强**
  - 完善单选框和多选框功能
  - 增强菜单功能
  - 添加图片导入功能
- **技术改进**
  - 将liblcl.dll嵌入到可执行文件中
  - 添加UPX压缩支持
  - 优化构建脚本

### v1.1.0 (2025-11-11)
- **多线程支持**
  - 添加并发任务管理器
  - 支持HTTP、计算、I/O和混合并发任务
  - 实现任务进度监控
- **数据库功能增强**
  - 完善数据库操作功能
  - 添加事务支持
  - 修复索引越界问题
- **错误处理改进**
  - 添加全局错误处理机制
  - 改进异常信息显示

### v1.0.0 (2025-11-10)
- **初始版本发布**
  - 实现基本的Excel操作功能
  - 实现基本的JSON操作功能
  - 实现基本的HTTP操作功能
  - 实现基本的数据库操作功能
  - 创建基本的GUI界面

## 🤝 贡献指南

我们欢迎社区贡献！如果您想为项目做出贡献，请遵循以下步骤：

1. Fork本项目到您的GitHub账户
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建Pull Request

### 代码贡献规范

- 遵循现有的代码风格和命名约定
- 为新功能添加适当的注释和文档
- 确保所有测试通过
- 提交信息应清晰描述更改内容

## 📄 许可证

本项目基于MIT许可证开源，详情请参阅[LICENSE](LICENSE)文件。

## 👨‍💻 作者

**yuan71058** - *项目所有者* - [GitHub](https://github.com/yuan71058)

## 📞 联系支持

如有问题或建议，请通过以下方式联系：

- **GitHub Issues**: [提交问题](https://github.com/yuan71058/GO_VCL/issues)
- **项目主页**: [GitHub仓库](https://github.com/yuan71058/GO_VCL)
- **发布页面**: [GitHub Releases](https://github.com/yuan71058/GO_VCL/releases)

## 🙏 致谢

感谢以下开源项目和贡献者：

- [Govcl](https://github.com/ying32/govcl) - Go语言的VCL绑定
- [Excelize](https://github.com/qax-os/excelize) - Go语言Excel处理库
- [GJSON](https://github.com/tidwall/gjson) - Go语言JSON解析库
- [HttpRequest](https://github.com/kirinlabs/HttpRequest) - Go语言HTTP客户端
- [Go-SQLite3](https://github.com/mattn/go-sqlite3) - Go语言SQLite驱动

---

<div align="center">
  <p>© 2025 GO_VCL Windows GUI应用程序. 保留所有权利.</p>
  <p>使用 ❤️ 和 Go 语言构建</p>
  <p><strong>最后更新</strong>: 2025-11-14 (v1.3.2)</p>
</div>