#!/bin/bash
# jobrunner 服务更新脚本

set -e

echo "=== jobrunner 服务更新 ==="

echo "1. 构建新版本..."
make build

echo "2. 停止服务..."
sudo systemctl stop jobrunner

echo "3. 更新二进制文件..."
sudo cp bin/jobrunner /usr/local/bin/jobrunner

echo "4. 启动服务..."
sudo systemctl start jobrunner

echo "5. 检查状态..."
sudo systemctl status jobrunner --no-pager

echo ""
echo "=== 更新完成 ==="
