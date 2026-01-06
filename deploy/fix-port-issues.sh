#!/bin/bash

# 端口问题自动修复脚本

echo "======================================"
echo "   端口问题诊断与修复脚本"
echo "======================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. 停止所有相关容器
echo -e "${YELLOW}步骤 1: 停止所有 Docker 容器...${NC}"
docker compose down 2>/dev/null || true
docker stop $(docker ps -aq) 2>/dev/null || true
echo -e "${GREEN}✓ 容器已停止${NC}"
echo ""

# 2. 检查端口占用
echo -e "${YELLOW}步骤 2: 检查端口占用...${NC}"
PORT_3000_PID=$(lsof -ti:3000 2>/dev/null)
PORT_8082_PID=$(lsof -ti:8082 2>/dev/null)

if [ -n "$PORT_3000_PID" ]; then
    echo -e "${RED}✗ 端口 3000 被进程 $PORT_3000_PID 占用${NC}"
    read -p "是否终止该进程? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        kill -9 $PORT_3000_PID
        echo -e "${GREEN}✓ 进程已终止${NC}"
    fi
else
    echo -e "${GREEN}✓ 端口 3000 可用${NC}"
fi

if [ -n "$PORT_8082_PID" ]; then
    echo -e "${RED}✗ 端口 8082 被进程 $PORT_8082_PID 占用${NC}"
    read -p "是否终止该进程? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        kill -9 $PORT_8082_PID
        echo -e "${GREEN}✓ 进程已终止${NC}"
    fi
else
    echo -e "${GREEN}✓ 端口 8082 可用${NC}"
fi
echo ""

# 3. 清理 Docker 资源
echo -e "${YELLOW}步骤 3: 清理 Docker 资源...${NC}"
read -p "是否清理未使用的 Docker 资源? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker system prune -f
    echo -e "${GREEN}✓ Docker 资源已清理${NC}"
fi
echo ""

# 4. 重新启动服务
echo -e "${YELLOW}步骤 4: 启动服务...${NC}"
cd "$(dirname "$0")/.."

# 检查是否仍有端口被占用
PORT_3000_CHECK=$(lsof -ti:3000 2>/dev/null)
PORT_8082_CHECK=$(lsof -ti:8082 2>/dev/null)

if [ -n "$PORT_3000_CHECK" ] || [ -n "$PORT_8082_CHECK" ]; then
    echo -e "${YELLOW}! 检测到端口仍被占用，使用备用端口配置${NC}"
    docker compose -f deploy/docker-compose.alternative-ports.yml up -d
    echo ""
    echo -e "${GREEN}✓ 服务已启动（备用端口）${NC}"
    echo "前端地址: http://localhost:3001"
    echo "后端地址: http://localhost:8083"
else
    docker compose up -d
    echo ""
    echo -e "${GREEN}✓ 服务已启动（默认端口）${NC}"
    echo "前端地址: http://localhost:3000"
    echo "后端地址: http://localhost:8082"
fi
echo ""

# 5. 检查服务状态
echo -e "${YELLOW}步骤 5: 检查服务状态...${NC}"
sleep 3
docker compose ps
echo ""

# 6. 显示日志
echo -e "${YELLOW}查看实时日志（按 Ctrl+C 退出）:${NC}"
echo "  docker compose logs -f"
echo ""
echo -e "${GREEN}======================================"
echo "   修复完成！"
echo "======================================${NC}"
