# 局域网唤醒工具 (Wake on LAN Tool)

一个简单易用的局域网唤醒工具，支持通过Web界面管理主机、发送WOL唤醒包、检测服务器状态以及测试端口连通性。

## 功能特性

- ✨ 主机管理（添加、编辑、删除）
- ⚡ WOL唤醒发送Magic Packet
- 🔍 服务器状态检测（ping/SSH/RDP并发检测）
- 🔌 端口连通性测试（3秒超时）
- 📤 数据导入导出（JSON格式）
- 🔐 密码保护功能
- 🌙 深色/浅色主题切换, 兼容移动端

## 系统要求

- Windows 操作系统
- Go 1.24.11 或更高版本
- 现代浏览器

## 安装运行

```bash
# 编译
build.bat

# 运行（默认端口999，密码adminer）
dwol.exe

# 自定义端口和密码
dwol.exe -p 8080 -pwd mypassword
```

访问 `http://localhost:999` 打开Web界面。

## 使用说明

### 主机管理
点击"➕ 添加主机"，填写MAC地址（必填）、主机地址、端口、设备名称等信息。

### 唤醒主机
点击主机列表中的"⚡ 唤醒"按钮，确认后发送WOL Magic Packet。

### 检查状态
点击"🔍 检查状态"按钮，并发检测所有主机状态（ping → SSH 22 → RDP 3389）。

### 测试端口
点击"🔌 测试端口"按钮，输入端口号测试连通性。

### 数据导入导出
点击"📤 导出数据"导出JSON文件，点击"📥 导入数据"导入JSON文件。

### 主题切换
点击右上角 🌙/☀️ 图标切换主题。

## 注意事项

1. 目标主机需支持并启用WOL功能
2. 唤醒主机和被唤醒主机必须在同一局域网
3. 确保防火墙允许相关端口通信
4. MAC地址支持多种格式：`00:11:22:33:44:55`、`00-11-22-33-44-55`、`001122334455`

## 技术栈

- 后端：Go 1.24.11
- 前端：Vue 3
- 样式：原生CSS

## 其他说明
- 页面效果图见`appimg`目录
- <img src="https://gcore.jsdelivr.net/gh/dhjz/dwol@master/appimg/app1.jpg" style="width: 340px;"/>
- <img src="https://gcore.jsdelivr.net/gh/dhjz/dwol@master/appimg/app2.jpg" style="width: 340px;"/>
