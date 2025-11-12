# 图标和Manifest修改指南

## 概述

根据Govcl官方文档，我们已成功修改了程序的图标和manifest信息，使应用程序具有自定义图标和版本信息。

## 实现方法

### 1. 创建资源文件

创建了`app.rc`资源文件，包含以下内容：
- 图标资源：使用`rgb.ico`作为应用程序图标
- Manifest文件：指定了`app.manifest`文件
- 版本信息：包含文件版本、产品版本、公司名称等信息

### 2. 创建Manifest文件

创建了`app.manifest`文件，包含以下内容：
- DPI感知设置：启用高DPI支持
- Windows通用控件依赖：确保现代UI控件可用
- 执行级别：设置为`asInvoker`（不需要管理员权限）
- Windows兼容性：支持Windows Vista到Windows 10

### 3. 修改编译脚本

修改了`build_simple.bat`脚本，添加了以下功能：
- 检查windres工具是否可用
- 使用windres编译资源文件为`.syso`格式
- 在编译过程中包含资源文件

## 文件说明

### app.rc
资源脚本文件，定义了应用程序的图标、manifest和版本信息。

### app.manifest
XML格式的manifest文件，定义了应用程序的执行要求和兼容性设置。

### defaultRes_windows_386.syso
编译后的资源文件，由windres从app.rc生成，在Go编译过程中自动链接到可执行文件中。

## 使用方法

1. **修改图标**：
   - 替换`rgb.ico`文件为您自己的图标文件
   - 确保图标文件为`.ico`格式，建议包含多种尺寸（16x16, 32x32, 48x48, 256x256）

2. **修改版本信息**：
   - 编辑`app.rc`文件中的版本信息部分
   - 修改文件版本、产品版本、公司名称等字段

3. **修改执行级别**：
   - 编辑`app.manifest`文件中的`requestedExecutionLevel`元素
   - 设置为`requireAdministrator`以请求管理员权限

4. **编译程序**：
   - 运行`build_simple.bat`脚本
   - 脚本会自动处理资源文件的编译和链接

## 注意事项

1. **windres工具**：
   - 需要安装MinGW-w64或TDM-GCC以获取windres工具
   - 如果没有windres，程序仍可编译，但不会包含自定义图标和版本信息

2. **资源文件限制**：
   - 一个Go工程中只能有一个.syso文件
   - 如果使用自定义资源文件，不能再导入winappres包

3. **图标格式**：
   - 图标必须为.ico格式
   - 建议使用专业的图标编辑工具创建多尺寸图标

## 验证方法

1. **查看图标**：
   - 编译后的可执行文件应显示自定义图标
   - 运行程序时，任务栏和窗口标题栏应显示自定义图标

2. **查看版本信息**：
   - 右键点击可执行文件，选择"属性"
   - 在"详细信息"选项卡中查看版本信息

3. **查看兼容性**：
   - 程序应在Windows Vista到Windows 10上正常运行
   - 在高DPI显示器上，界面应清晰显示

## 故障排除

1. **图标未显示**：
   - 检查图标文件是否为.ico格式
   - 确认windres工具已正确安装
   - 查看编译日志中是否有资源文件创建错误

2. **版本信息未显示**：
   - 检查app.rc文件中的版本信息格式是否正确
   - 确认资源文件已成功编译

3. **程序无法运行**：
   - 检查manifest文件格式是否正确
   - 尝试移除manifest文件，使用默认设置

## 参考资料

- [Govcl官方文档](https://gitee.com/ying32/govcl/wikis/pages?sort_id=410058&doc_id=102420)
- [Microsoft Manifest文档](https://docs.microsoft.com/en-us/windows/win32/sbscs/application-manifests)
- [Windows资源文件文档](https://docs.microsoft.com/en-us/windows/win32/menurc/resource-files)