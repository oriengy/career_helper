# 部署指南

## 问题诊断

### 1. 检查端口占用

在服务器上运行：

```bash
# 方法1：使用 lsof
sudo lsof -i :3000
sudo lsof -i :8082

# 方法2：使用 netstat
sudo netstat -tulpn | grep :3000
sudo netstat -tulpn | grep :8082

# 方法3：使用 ss
sudo ss -tulpn | grep :3000
sudo ss -tulpn | grep :8082
```

### 2. 停止占用端口的服务

```bash
# 停止旧的 Docker 容器
docker compose down

# 查看所有容器（包括已停止的）
docker ps -a

# 强制删除所有容器
docker rm -f $(docker ps -aq)

# 如果端口被其他进程占用，找到 PID 后停止
sudo kill -9 <PID>
```

### 3. 清理 Docker 资源

```bash
# 清理未使用的容器、网络、镜像
docker system prune -a

# 清理未使用的卷
docker volume prune
```

## 部署方案

### 方案 1：默认端口（推荐本地开发）

使用根目录的 `docker-compose.yml`：

```bash
cd /www/wwwroot/career_helper-main
docker compose up -d
```

- 前端：http://localhost:3000
- 后端：http://localhost:8082

### 方案 2：备用端口（端口冲突时）

使用 `deploy/docker-compose.alternative-ports.yml`：

```bash
cd /www/wwwroot/career_helper-main
docker compose -f deploy/docker-compose.alternative-ports.yml up -d
```

- 前端：http://localhost:3001
- 后端：http://localhost:8083

### 方案 3：Nginx 反向代理（推荐生产环境）

1. 编辑 `deploy/nginx.conf`，替换 `your-domain.com` 为你的域名

2. 启动服务：

```bash
cd /www/wwwroot/career_helper-main
docker compose -f deploy/docker-compose.nginx.yml up -d
```

3. 访问：
   - 前端：http://your-domain.com
   - 后端 API：http://your-domain.com/api

### 方案 4：仅后端（前端部署在其他地方）

修改 `docker-compose.yml`，只启动 backend：

```bash
docker compose up -d backend
```

## 常见错误解决

### 错误：`bind: address already in use`

**原因**：端口被占用

**解决**：
1. 检查是否有旧容器运行：`docker ps`
2. 停止旧容器：`docker compose down`
3. 检查其他进程：`sudo lsof -i :3000`
4. 或者使用备用端口方案

### 错误：`Cannot connect to the Docker daemon`

**原因**：Docker 服务未运行

**解决**：
```bash
# 启动 Docker 服务
sudo systemctl start docker

# 设置开机自启
sudo systemctl enable docker
```

### 错误：`permission denied`

**原因**：权限不足

**解决**：
```bash
# 使用 sudo
sudo docker compose up -d

# 或将用户加入 docker 组
sudo usermod -aG docker $USER
# 注销后重新登录生效
```

## 生产环境建议

1. **使用 Nginx 反向代理**
   - 统一入口
   - 方便配置 SSL
   - 支持负载均衡

2. **设置重启策略**
   - 已在配置中添加 `restart: unless-stopped`
   - 服务器重启后自动启动容器

3. **配置日志**
   ```bash
   # 查看容器日志
   docker compose logs -f frontend
   docker compose logs -f backend
   ```

4. **监控资源使用**
   ```bash
   # 查看容器资源使用
   docker stats
   ```

5. **定期备份**
   - 备份数据库
   - 备份配置文件
   - 备份上传的文件

## 快速命令参考

```bash
# 启动服务
docker compose up -d

# 停止服务
docker compose down

# 重启服务
docker compose restart

# 查看日志
docker compose logs -f

# 重新构建并启动
docker compose up -d --build

# 查看运行状态
docker compose ps
```
