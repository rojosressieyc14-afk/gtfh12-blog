#!/bin/bash
set -euo pipefail

echo "=== PulseBlog 环境检查 ==="

check_cmd() {
    if command -v "$1" >/dev/null 2>&1; then
        echo "  [OK] $1: $($1 --version 2>&1 | head -1)"
    else
        echo "  [FAIL] $1: 未安装"
        return 1
    fi
}

FAIL=0
check_cmd docker || FAIL=1
check_cmd docker-compose || FAIL=1
check_cmd node || FAIL=1
check_cmd npm || FAIL=1
check_cmd go || FAIL=1

if [ "$FAIL" -eq 1 ]; then
    echo "请安装缺失的依赖后重试。"
    exit 1
fi

echo ""
echo "所有依赖已就绪。"
echo "运行 ./scripts/start-prod.sh 启动生产环境。"
