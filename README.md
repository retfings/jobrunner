# jobrunner

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

一个用 Go 语言实现的 Linux 定时任务服务工具，按固定间隔检查目录状态并自动创建序号目录。

## 功能特性

- 注册为 systemd 服务
- 使用固定间隔轮询检查目录状态
- 当目录不为空时自动创建下一个序号目录
- 按时间规则自动生成目录（如 `/var/jobs/2026/03/14/001`）
- 使用 YAML 配置文件
- 每天 0 点重置序号

## 目录结构

```
jobrunner/
├── cmd/
│   └── jobrunner/
│       └── main.go          # 程序入口
├── pkg/
│   ├── monitor/
│   │   └── monitor.go       # 目录监控器
│   ├── generator/
│   │   └── generator.go     # 目录生成器
│   ├── config/
│   │   └── config.go        # 配置加载与解析
│   └── service/
│       └── service.go       # systemd 服务管理
├── config/
│   └── config.yaml          # 默认配置文件
├── go.mod
├── go.sum
├── Makefile                 # 构建脚本
└── install-service.sh       # 服务安装脚本
```

## 配置说明

配置文件 (`config/config.yaml`)：

```yaml
# 基础路径 - 目录将在此路径下创建
base_path: "/var/jobs"

# 检查间隔 - 每隔多久检查当前目录是否为空
check_interval: 30s

# 初始延迟 - 启动后延迟多久创建首个目录
initial_delay: 5s
```

## 目录命名规则

目录格式：`{base_path}/{YYYY}/{MM}/{DD}/{NNN}`

- YYYY: 年（4 位）
- MM: 月（2 位）
- DD: 日（2 位）
- NNN: 序号（3 位，从 001 开始递增）

示例：
- 2026 年 3 月 14 日第 1 个目录：`/var/jobs/2026/03/14/001`
- 2026 年 3 月 14 日第 2 个目录：`/var/jobs/2026/03/14/002`

## 构建和安装

### 构建项目

```bash
make build
```

### 安装为 systemd 服务

```bash
sudo make install
```

或使用安装脚本：

```bash
sudo ./install-service.sh
```

### 手动运行

```bash
./bin/jobrunner --config config/config.yaml
```

## 服务管理

```bash
# 启动服务
sudo systemctl start jobrunner

# 停止服务
sudo systemctl stop jobrunner

# 重启服务
sudo systemctl restart jobrunner

# 查看状态
sudo systemctl status jobrunner

# 查看日志
sudo journalctl -u jobrunner -f

# 卸载服务
sudo make uninstall
```

## 命令行参数

```
-config string
      配置文件路径
-command string
      命令：run/install/start/stop/restart/uninstall (default "run")
-version
      显示版本
```

## 工作流程

1. 服务启动 → 自动创建当天首个目录（如 `/var/jobs/2026/03/14/001`）
2. 每隔 N 秒检查当前目录是否为空
3. 如果目录不为空 → 立即创建下一个序号目录
4. 每天 0 点重置序号，从 001 重新开始

## 验证方式

1. 构建项目：`make build`
2. 安装服务：`sudo make install`
3. 启动服务：`sudo systemctl start jobrunner`
4. 查看状态：`sudo systemctl status jobrunner`
5. 查看日志：`journalctl -u jobrunner -f`
6. 测试：往当前目录放入文件，观察是否自动创建下一个目录

## 开发

```bash
# 运行测试
go test -v ./pkg/...

# 格式化代码
go fmt ./...

# 下载依赖
go mod tidy
```

## License

MIT
