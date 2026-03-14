#!/bin/bash

# jobrunner 服务安装脚本
# 用法：sudo ./install-service.sh [config_path]

set -e

BINARY_NAME="jobrunner"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_PATH="${SCRIPT_DIR}/bin/${BINARY_NAME}"
CONFIG_SOURCE="${SCRIPT_DIR}/config/config.yaml"
CONFIG_DEST="/etc/jobrunner/config.yaml"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否以 root 运行
if [[ $EUID -ne 0 ]]; then
    log_error "此脚本必须以 sudo 运行"
    exit 1
fi

# 检查二进制文件
if [[ ! -f "${BINARY_PATH}" ]]; then
    log_warn "二进制文件不存在，先构建..."
    make build
fi

# 创建目录
log_info "创建配置目录..."
mkdir -p /etc/jobrunner
mkdir -p /usr/local/bin

# 复制文件
log_info "复制二进制文件..."
cp "${BINARY_PATH}" "/usr/local/bin/${BINARY_NAME}"
chmod +x "/usr/local/bin/${BINARY_NAME}"

# 复制配置文件
CONFIG_PATH="${1:-${CONFIG_DEST}}"
if [[ ! -f "${CONFIG_PATH}" ]]; then
    log_info "复制配置文件..."
    cp "${CONFIG_SOURCE}" "${CONFIG_DEST}"
else
    log_info "使用现有配置文件：${CONFIG_PATH}"
fi

# 安装服务
log_info "注册 systemd 服务..."
/usr/local/bin/${BINARY_NAME} --command install --config "${CONFIG_PATH}"

# 完成
echo ""
log_info "服务安装完成!"
echo ""
echo "可用命令:"
echo "  sudo systemctl start ${BINARY_NAME}    - 启动服务"
echo "  sudo systemctl stop ${BINARY_NAME}     - 停止服务"
echo "  sudo systemctl restart ${BINARY_NAME}  - 重启服务"
echo "  sudo systemctl status ${BINARY_NAME}   - 查看状态"
echo "  sudo journalctl -u ${BINARY_NAME} -f   - 查看日志"
echo ""
