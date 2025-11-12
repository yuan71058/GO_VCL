# GO_VCL Windows GUI应用程序 v1.0.0 发布说明

## 概述
GO_VCL是一个基于Go语言开发的Windows GUI应用程序，使用Govcl UI库，提供完整的图形界面和多功能集成。

## 主要功能
- **GUI界面**: 完整的Windows窗口应用程序，包含标签、按钮、表格、编辑框、状态栏等组件
- **Excel操作**: 支持xlsx文件导入导出，包含新增的导入功能
- **JSON处理**: 支持JSON文件解析和生成
- **HTTP请求**: 支持GET/POST请求，集成测试API
- **数据库操作**: 8个完整的数据库操作流程（模拟模式）
- **多线程支持**: 多线程任务管理，支持并发执行和结果收集
- **图片导入**: 支持用户选择图片文件并在带边框的图片框中显示
- **菜单功能**: 包含文件、编辑、工具和帮助菜单，提供清空编辑框等快捷操作
- **UI组件**: 单选框、多选框、获取选中状态按钮等交互组件

## 技术特点
- **模块化设计**: 采用管理器模式，各功能模块独立实现，易于扩展和维护
- **错误处理**: 完善的异常处理机制，确保程序稳定运行
- **代码注释**: 所有代码添加了详细的中文注释，提高可读性和可维护性
- **独立部署**: 通过liblcl.dll嵌入技术，实现程序独立运行无需外部依赖
- **压缩优化**: 使用UPX压缩技术，减小可执行文件大小约67.9%

## 系统要求
- Windows 10/11操作系统
- 无需额外依赖（已嵌入liblcl.dll）

## 安装与使用
1. 下载`windows-gui-app.exe`文件
2. 双击运行程序
3. 使用界面上的按钮和菜单进行各项操作

## 文件结构
```
GO_VCL/
├── dist/windows-gui-app.exe    # 可执行文件（已压缩）
├── docs/                       # 文档目录
│   ├── USER_MANUAL.md         # 用户手册
│   ├── ICON_AND_MANIFEST_GUIDE.md  # 图标和清单指南
│   ├── IMPORT_EXCEL_FEATURE.md     # Excel导入功能说明
│   └── UPX_COMPRESSION_FEATURE.md # UPX压缩功能说明
├── BUILD_SCRIPTS_README.md    # 构建脚本说明
└── 说明文档.md                 # 项目说明文档
```

## 更新日志
### v1.0.0 (2025-11-12)
- 初始发布版本
- 实现所有核心功能模块
- 完成UI界面设计和交互功能
- 添加完整的错误处理机制
- 实现liblcl.dll嵌入和UPX压缩功能
- 完善代码注释和文档

## 技术栈
- **主框架**: github.com/ying32/govcl
- **Excel操作**: github.com/xuri/excelize/v2  
- **JSON处理**: github.com/tidwall/gjson
- **HTTP请求**: github.com/kirinlabs/HttpRequest
- **数据库**: SQLite3（模拟模式）

## 许可证
本项目采用MIT许可证。

## 联系方式
如有问题或建议，请通过GitHub Issues联系。