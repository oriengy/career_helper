#!/bin/bash

echo "=== 端口检查脚本 ==="
echo ""

# 检查 3000 端口（前端）
echo "检查端口 3000 (前端)..."
if command -v lsof &> /dev/null; then
    lsof -i :3000
elif command -v netstat &> /dev/null; then
    netstat -tulpn | grep :3000
elif command -v ss &> /dev/null; then
    ss -tulpn | grep :3000
fi

echo ""

# 检查 8082 端口（后端）
echo "检查端口 8082 (后端)..."
if command -v lsof &> /dev/null; then
    lsof -i :8082
elif command -v netstat &> /dev/null; then
    netstat -tulpn | grep :8082
elif command -v ss &> /dev/null; then
    ss -tulpn | grep :8082
fi

echo ""
echo "=== Docker 容器状态 ==="
docker ps -a

echo ""
echo "=== 解决方案 ==="
echo "如果端口被占用，可以："
echo "1. 停止占用端口的进程: kill -9 <PID>"
echo "2. 停止所有 Docker 容器: docker compose down"
echo "3. 修改 docker-compose.yml 中的端口映射"
echo ""
