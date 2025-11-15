# GitHub Pages 配置

本目录包含GO_VCL项目的GitHub Pages网站文件。

## 启用GitHub Pages

1. 在您的GitHub仓库中，进入"Settings"选项卡
2. 在左侧菜单中找到"Pages"选项
3. 在"Source"部分，选择"Deploy from a branch"
4. 在"Branch"下拉菜单中，选择"main"分支
5. 在文件夹下拉菜单中，选择"/docs"文件夹
6. 点击"Save"保存设置

## 文件结构

```
docs/
└── github-pages/
    └── index.html    # 主页文件
```

## 自定义域名（可选）

如果您想使用自定义域名，可以在`docs/github-pages/`目录下创建一个`CNAME`文件，内容为您的域名。

## 注意事项

- GitHub Pages会自动渲染`/docs`文件夹中的文件
- 确保所有相对路径都相对于`/docs/github-pages/`目录
- 图片和其他资源文件应使用正确的相对路径引用