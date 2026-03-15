.PHONY: build clean install uninstall test run help

# 变量
BINARY_NAME=jobrunner
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=cmd/jobrunner/main.go
CONFIG_SOURCE=config/config.yaml
CONFIG_DEST=/etc/jobrunner/config.yaml

# 默认目标
all: build

# 构建
build:
	@echo "构建 $(BINARY_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "构建完成：$(BINARY_PATH)"

# 清理
clean:
	@echo "清理构建文件..."
	@rm -rf bin
	@echo "清理完成"

# 安装服务
install: build
	@echo "安装服务..."
	@sudo mkdir -p /usr/local/bin
	@sudo cp $(BINARY_PATH) /usr/local/bin/$(BINARY_NAME)
	@sudo mkdir -p /etc/jobrunner
	@sudo cp $(CONFIG_SOURCE) $(CONFIG_DEST)
	@sudo $(BINARY_NAME) --command install --config $(CONFIG_DEST)
	@echo "服务安装完成"
	@echo "使用 'sudo systemctl start $(BINARY_NAME)' 启动服务"
	@echo "使用 'sudo systemctl status $(BINARY_NAME)' 查看状态"

# 更新服务
update: build
	@echo "更新服务..."
	@sudo cp $(BINARY_PATH) /usr/local/bin/$(BINARY_NAME)
	@sudo systemctl restart $(BINARY_NAME)
	@echo "服务更新完成"
	@echo "使用 'sudo systemctl status $(BINARY_NAME)' 查看状态"

# 卸载服务
uninstall:
	@echo "卸载服务..."
	@sudo systemctl stop $(BINARY_NAME) 2>/dev/null || true
	@sudo $(BINARY_NAME) --command uninstall 2>/dev/null || true
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@sudo rm -rf /etc/jobrunner
	@echo "服务卸载完成"

# 运行（前台）
run: build
	@echo "运行 $(BINARY_NAME)..."
	@./$(BINARY_PATH) --config $(CONFIG_SOURCE)

# 测试
test:
	@echo "运行测试..."
	@go test -v ./pkg/...

# 格式化
fmt:
	@echo "格式化代码..."
	@go fmt ./...

# 下载依赖
deps:
	@echo "下载依赖..."
	@go mod tidy

# 帮助
help:
	@echo "jobrunner 构建脚本"
	@echo ""
	@echo "可用命令:"
	@echo "  make build      - 构建程序"
	@echo "  make clean      - 清理构建文件"
	@echo "  make install    - 安装为 systemd 服务"
	@echo "  make update     - 更新服务（重新构建并重启）"
	@echo "  make uninstall  - 卸载服务"
	@echo "  make run        - 前台运行程序"
	@echo "  make test       - 运行测试"
	@echo "  make fmt        - 格式化代码"
	@echo "  make deps       - 下载依赖"
	@echo "  make help       - 显示帮助"
