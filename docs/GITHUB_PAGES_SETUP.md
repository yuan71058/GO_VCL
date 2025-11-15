# GO_VCL GitHub Pages 设置指南

本指南将帮助您为GO_VCL项目设置GitHub Pages主页。

## 前提条件

1. 您已有一个GitHub账户
2. 您已将GO_VCL项目推送到GitHub仓库

## 设置步骤

### 方法一：使用GitHub Pages设置（推荐）

1. 登录GitHub，进入您的GO_VCL仓库
2. 点击仓库顶部的"Settings"选项卡
3. 在左侧菜单中找到"Pages"选项
4. 在"Source"部分，选择"Deploy from a branch"
5. 在"Branch"下拉菜单中，选择"main"分支
6. 在文件夹下拉菜单中，选择"/docs"文件夹
7. 点击"Save"保存设置

几分钟后，您的GitHub Pages网站将在以下地址可用：
`https://[您的用户名].github.io/GO_VCL/`

### 方法二：使用GitHub Actions自动部署

1. 确保您已在仓库中创建了`.github/workflows/deploy-pages.yml`文件（已包含在项目中）
2. 按照方法一的步骤设置GitHub Pages
3. 当您推送代码到main分支时，GitHub Actions将自动部署您的网站

## 自定义域名（可选）

如果您想使用自定义域名：

1. 在`docs/github-pages/`目录下创建一个`CNAME`文件
2. 在文件中添加您的域名，例如：`www.yourdomain.com`
3. 在您的域名提供商处配置DNS记录
4. 在GitHub Pages设置中启用自定义域名

## 网站内容

您的GitHub Pages网站包含以下内容：

- 项目概述和功能介绍
- 应用程序截图展示
- 快速开始指南
- 技术架构说明
- 下载和安装说明

## 更新网站

要更新网站内容：

1. 修改`docs/github-pages/index.html`文件
2. 提交并推送更改到GitHub
3. 等待GitHub Actions自动部署（如果使用方法二）
4. 或等待GitHub Pages自动更新（如果使用方法一）

## 故障排除

如果网站无法正常显示：

1. 检查GitHub Pages设置是否正确
2. 确保所有文件路径正确
3. 查看GitHub Actions日志（如果使用方法二）
4. 确保仓库是公开的（私有仓库需要升级到GitHub Pro才能使用GitHub Pages）

## 联系支持

如果您在设置过程中遇到问题，请：

1. 查看GitHub Pages官方文档
2. 在GitHub仓库中创建Issue
3. 联系项目维护者

---

祝您使用愉快！