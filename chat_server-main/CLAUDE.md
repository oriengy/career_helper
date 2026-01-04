# ChatHandy Server - 妙聊后端服务

**版本**: v1.0.0  
**状态**: MVP 完成，可正式运营  
**最后更新**: 2025-01-02

## 项目概述

ChatHandy Server 是一个基于 Go 语言和 gRPC/Connect 协议构建的 AI 驱动社交辅助后端服务。该项目专为男女之间的沟通场景设计，通过 AI 技术帮助用户理解异性交流中的潜在含义，提升沟通效果。

### 核心功能特性

- **用户认证系统**: 支持微信小程序登录，基于 JWT 的身份认证
- **聊天会话管理**: 创建、查询、更新、删除聊天会话
- **AI 翻译服务**: 基于火山引擎 API 的智能翻译，专门针对男女沟通场景优化
- **消息管理**: 支持好友消息、咨询消息的存储和查询
- **个人资料管理**: 用户和对话对象的个人信息管理
- **文件上传服务**: 集成阿里云 OSS 的文件存储服务
- **图片解析功能**: 支持聊天截图的 OCR 识别和内容导入

### 技术架构

- **框架**: Gin Web Framework
- **RPC 协议**: gRPC + Connect-Go
- **数据库**: MySQL + GORM ORM
- **认证**: JWT Token
- **AI 服务**: 火山引擎 AI API、OpenAI API
- **文件存储**: 阿里云 OSS
- **配置管理**: Viper
- **ID 生成**: 雪花算法 (Snowflake)

## 目录结构

```
chat_server/
├── app_server/                 # 主应用服务器
│   ├── cmd/                    # 命令行入口
│   │   └── api_server/         # API 服务器主程序
│   ├── config.yaml             # 开发环境配置
│   ├── config.prod.yaml        # 生产环境配置
│   ├── domain/                 # 业务领域层 [📖 详细文档](./app_server/domain/CLAUDE.md)
│   │   ├── chat_session.go     # 聊天会话领域逻辑
│   │   ├── new_user.go         # 新用户注册逻辑
│   │   └── exterr/             # 业务错误定义
│   ├── http/                   # HTTP 路由和处理器
│   │   └── file/               # 文件上传处理
│   ├── model/                  # 数据模型 [📖 详细文档](./app_server/model/CLAUDE.md)
│   │   ├── user.go             # 用户模型
│   │   ├── profile.go          # 个人资料模型
│   │   ├── chat.go             # 聊天会话模型
│   │   └── message.go          # 消息模型
│   ├── pkg/                    # 公共包和工具类 [📖 详细文档](./app_server/pkg/CLAUDE.md)
│   │   ├── aiapi/              # AI API 集成（OCR、图片解析）
│   │   ├── cbind/              # Connect-Go 与 Gin 绑定器
│   │   ├── cfg/                # 配置管理（基于 Viper）
│   │   ├── db/                 # 数据库连接（GORM + 自动ID生成）
│   │   ├── fn/                 # 函数式编程工具集
│   │   ├── httpc/              # HTTP 客户端（req/v3）
│   │   ├── idgen/              # ID 生成器（雪花算法）
│   │   ├── jwt/                # JWT 认证管理
│   │   ├── openaic/            # OpenAI/火山引擎客户端
│   │   └── ossc/               # 阿里云 OSS 客户端
│   ├── proto/                  # 生成的 Protocol Buffers 代码
│   ├── service/                # 业务服务层 [📖 详细文档](./app_server/service/CLAUDE.md)
│   │   ├── auth/               # 认证授权服务
│   │   ├── user/               # 用户服务
│   │   ├── chat/               # 会话管理服务
│   │   ├── translate/          # AI 翻译服务
│   │   ├── message/            # 消息管理服务
│   │   ├── profile/            # 个人资料服务
│   │   └── admin/              # 管理后台服务
│   └── scripts/                # 脚本工具
├── proto/                      # Protocol Buffers 定义文件
├── buf.yaml                    # Buf 配置文件
└── buf.gen.yaml                # Buf 代码生成配置
```

## 核心模块详解

### 1. 用户认证 (Authentication)

**位置**: `app_server/service/auth/` 和 `app_server/service/user/`

#### 核心功能
- **微信登录**: 支持微信小程序授权登录，自动获取 openid/unionid
- **手机号登录**: 开发阶段使用魔法验证码 "1234"
- **JWT Token**: 基于 JWT 的无状态认证，有效期 365 天
- **认证拦截器**: Connect-Go 拦截器自动验证 Authorization Header

#### 技术实现
- **拦截器机制**: 使用 Connect-Go 的 UnaryInterceptor 实现全局认证
- **上下文注入**: 将用户 ID 注入到 context 中，供后续服务使用
- **错误处理**: 认证失败返回 `CodeUnauthenticated` 错误

#### 关键文件
- `service/auth/auth.go`: 认证中间件和 Token 解析
- `service/user/wx_login.go`: 微信登录实现
- `pkg/jwt/jwt.go`: JWT Token 生成和验证

### 2. 聊天会话管理 (Chat Sessions)

**位置**: `app_server/service/chat/` 和 `app_server/domain/`

#### 核心功能
- **会话 CRUD**: 创建、查询、更新、删除聊天会话
- **快速创建**: 支持同时创建会话和关联的个人资料
- **分页查询**: 默认每页 20 条，按 ID 倒序排列
- **权限验证**: 确保用户只能操作自己的会话

#### 技术实现
- **事务一致性**: 使用 domain 层的 `CreateChatSession` 保证原子性
- **数据关联**: 自动关联 Profile 信息，优先使用 Profile 的名称和头像
- **软删除**: 基于 GORM 的软删除机制

#### 数据模型
```go
type ChatSession struct {
    gorm.Model
    Name      string  // 会话名称
    UserID    uint    // 用户ID
    ProfileID uint    // 关联的个人资料ID
    Avatar    string  // 会话头像
}
```

### 3. AI 翻译服务 (Translation)

**位置**: `app_server/service/translate/`

#### 核心功能
- **智能翻译**: 专为男女沟通场景设计的 AI 翻译
- **上下文感知**: 自动获取前后 24 小时的对话历史作为上下文
- **双向翻译**: 支持「翻译给男生」和「翻译给女生」两种模式
- **结果存储**: 翻译结果自动保存为咨询消息

#### 技术实现
- **AI 集成**: 主要使用火山引擎 AI API，OpenAI 作为备选
- **时间范围查询**: 利用雪花算法 ID 的时间特性进行高效查询
- **消息关联**: 通过 ParentID 关联原始消息和翻译结果

#### 翻译提示模板
```
你是一个帮助男女之间相互理解彼此话语含义的助手。
翻译时请保持内容简短。

请理解两人的对话上下文：
%s

将这句%s说的话翻译给%s听:
%s
```

### 4. 消息管理 (Messages)

**位置**: `app_server/service/message/` 和 `app_server/model/`

#### 消息体系设计
- **统一模型**: 所有消息都存储在 `ConsultMessage` 表中
- **类型区分**: 通过 `msg_type` 字段区分不同类型的消息
- **角色系统**: 支持 SELF（自己）、FRIEND（好友）、AI（机器人）三种角色

#### 核心功能

**好友消息服务**:
- **消息查询**: 实际查询 ConsultMessage 表中的历史消息
- **批量创建**: 支持批量导入聊天记录
- **图片解析**: OCR 识别聊天截图并自动导入
- **智能去重**: 避免重复导入相同的消息

**咨询消息服务**:
- **灵活查询**: 支持按类型、ID 列表、会话 ID 等多种条件过滤
- **翻译关联**: AI 翻译结果通过 ParentID 关联原始消息
- **标签系统**: 使用 Tags 标记消息特征（demo、翻译方向等）

#### 技术亮点
- **倒序分页**: 查询时倒序，返回前再反转，保证时间顺序正确
- **火山 OCR**: 集成火山引擎 API 进行图片文字识别
- **角色识别**: 自动识别聊天记录中的对话角色

### 5. 个人资料管理 (Profiles)

**位置**: `app_server/service/profile/` 和 `app_server/model/profile.go`

#### 核心功能
- **资料 CRUD**: 创建、查询、更新、删除个人资料
- **头像管理**: 自动处理临时文件迁移到用户专属目录
- **OSS 集成**: 生成 10 年有效期的签名 URL
- **数据验证**: 性别限制为 "male" 或 "female"

#### 技术实现
- **文件路径管理**: 临时上传 → 用户目录迁移
- **权限控制**: 用户只能操作自己创建的资料
- **自动清理**: 文件迁移后清理临时文件

### 6. 文件服务 (File Service)

**位置**: `app_server/http/file/`

#### 核心功能
- **微信文件上传**: 专门适配微信小程序的文件上传接口
- **OSS 存储**: 集成阿里云 OSS，支持内网/公网双客户端
- **安全验证**: 文件大小限制、类型验证、防恶意上传

#### 技术实现
- **临时存储**: 先保存到临时目录，后续迁移到用户目录
- **路径规范**: `uploads/` 临时目录，`user/{user_id}/` 用户目录
- **Gin 集成**: 使用 Gin 框架处理 multipart/form-data

### 7. 领域驱动设计 (Domain Layer)

**位置**: `app_server/domain/`

#### 设计理念
- **业务逻辑集中化**: 核心业务规则都在 domain 层实现
- **事务一致性**: 使用数据库事务确保操作原子性
- **错误标准化**: 通过 `exterr` 包提供统一的业务错误

#### 核心业务逻辑
- **用户注册流程**: 查找或创建用户，自动初始化演示数据
- **新用户初始化**: 创建主 Profile，复制演示数据
- **会话创建**: Profile 和 ChatSession 在同一事务中创建
- **演示数据管理**: 为新用户提供典型场景的示例对话

## 数据库模型设计

### 技术特性
- **ORM 框架**: GORM v2，支持自动迁移和软删除
- **ID 生成**: 雪花算法生成分布式唯一 ID
- **时间处理**: 统一使用 timestamppb 进行 Proto 转换
- **JSON 序列化**: 复杂类型使用 GORM serializer

### 核心数据表

#### 用户表 (user)
```sql
CREATE TABLE `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL COMMENT '用户姓名',
  `im_name` varchar(255) DEFAULT NULL COMMENT 'IM昵称',
  `external_id` varchar(255) DEFAULT NULL COMMENT '微信OpenID',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `avatar` varchar(500) DEFAULT NULL COMMENT '头像URL',
  `profile_id` bigint unsigned DEFAULT NULL COMMENT '主Profile ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_external_id` (`external_id`),
  KEY `idx_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 个人资料表 (profile)
```sql
CREATE TABLE `profile` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL COMMENT '所属用户ID',
  `name` varchar(255) DEFAULT NULL COMMENT '姓名',
  `im_name` varchar(255) DEFAULT NULL COMMENT 'IM名称',
  `avatar` varchar(500) DEFAULT NULL COMMENT '头像',
  `age` int DEFAULT NULL COMMENT '年龄',
  `gender` varchar(10) DEFAULT NULL COMMENT '性别(male/female)',
  `birthday` datetime DEFAULT NULL COMMENT '生日',
  `birth_location` varchar(255) DEFAULT NULL COMMENT '出生地',
  `current_location` varchar(255) DEFAULT NULL COMMENT '现居地',
  PRIMARY KEY (`id`),
  KEY `idx_profile_user_id` (`user_id`),
  KEY `idx_profile_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 聊天会话表 (chat_session)
```sql
CREATE TABLE `chat_session` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL COMMENT '会话名称',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `profile_id` bigint unsigned NOT NULL COMMENT '关联Profile ID',
  `avatar` varchar(500) DEFAULT NULL COMMENT '会话头像',
  PRIMARY KEY (`id`),
  KEY `idx_chat_session_user_id` (`user_id`),
  KEY `idx_chat_session_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 咨询消息表 (consult_message)
```sql
CREATE TABLE `consult_message` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `session_id` bigint unsigned NOT NULL COMMENT '会话ID',
  `parent_id` bigint unsigned DEFAULT NULL COMMENT '父消息ID',
  `profile_id` bigint unsigned DEFAULT NULL COMMENT 'Profile ID',
  `role` varchar(20) NOT NULL COMMENT '角色(SELF/FRIEND/AI)',
  `msg_type` varchar(20) NOT NULL COMMENT '类型(HISTORY/TRANSLATE)',
  `content` text COMMENT '消息内容',
  `tags` json DEFAULT NULL COMMENT '标签数组',
  `msg_at` datetime NOT NULL COMMENT '消息时间',
  PRIMARY KEY (`id`),
  KEY `idx_consult_message_user_id` (`user_id`),
  KEY `idx_consult_message_session_id` (`session_id`),
  KEY `idx_consult_message_msg_at` (`msg_at`),
  KEY `idx_consult_message_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## API 接口说明

### 用户服务 (UserService)
```
POST /user.UserService/WxUserLogin     # 微信用户登录
POST /user.UserService/GetUserProfile  # 获取用户资料
```

### 聊天服务 (ChatService)
```
POST /chat.ChatService/ListChatSessions    # 查询聊天会话列表
POST /chat.ChatService/CreateChatSession   # 创建聊天会话
POST /chat.ChatService/DeleteChatSession   # 删除聊天会话
POST /chat.ChatService/UpdateChatSession   # 更新聊天会话
```

### 翻译服务 (TranslateService)
```
POST /translate.TranslateService/Translate                    # 通用翻译
POST /translate.TranslateService/TranslateFriendMessage       # 翻译好友消息
```

### 消息服务 (MessageService)
```
POST /message.ConsultMessageService/ListConsultMessages      # 查询咨询消息
POST /message.ConsultMessageService/CreateConsultMessage     # 创建咨询消息
POST /message.FriendMessageService/ListFriendMessages        # 查询好友消息
POST /message.FriendMessageService/CreateFriendMessage       # 创建好友消息
POST /message.FriendMessageService/DeleteFriendMessage       # 删除好友消息
POST /message.FriendMessageService/ParseImageMessages        # 解析图片中的聊天记录
```

### 个人资料服务 (ProfileService)
```
POST /profile.ProfileService/ListProfiles      # 查询个人资料列表
POST /profile.ProfileService/CreateProfile     # 创建个人资料
POST /profile.ProfileService/UpdateProfile     # 更新个人资料
POST /profile.ProfileService/DeleteProfile     # 删除个人资料
```

### 文件服务 (FileService)
```
POST /file/wx_upload    # 微信文件上传
GET  /ping             # 健康检查
```

## 配置文件说明

### 主配置 (config.yaml)
```yaml
# 数据库配置
db:
  dsn: "数据库连接字符串"

# JWT 配置
jwt:
  secret: "JWT密钥"

# 微信配置
wechat:
  mp:
    translator:
      app_id: "微信小程序AppID"
      app_secret: "微信小程序AppSecret"
      prompt: "翻译提示模板"

# AI 服务配置
ai:
  volces:
    api_key: "火山引擎API密钥"
    base_url: "API基础URL"
    models:
      chat: "聊天模型"
      ocr: "OCR模型"

# 阿里云 OSS 配置
aliyun:
  oss:
    access_key_id: "AccessKey ID"
    access_key_secret: "AccessKey Secret"
    endpoint: "OSS终端节点"
    user_file_bucket: "用户文件存储桶"

# 服务器配置
server:
  address: ":8080"
```

## 开发和部署指南

### 环境要求
- Go 1.23.6+
- MySQL 5.7+
- 火山引擎 AI API 账号
- 阿里云 OSS 账号
- 微信小程序开发者账号

### 本地开发
1. **克隆项目并安装依赖**:
   ```bash
   cd chat_server/app_server
   go mod tidy
   ```

2. **配置数据库**:
   - 创建 MySQL 数据库
   - 更新 `config.yaml` 中的数据库连接信息

3. **生成 Protocol Buffers 代码**:
   ```bash
   # 在项目根目录执行
   buf generate
   ```

4. **启动服务**:
   ```bash
   # 启动 API 服务器
   go run cmd/api_server/main.go -c config.yaml
   ```

### 生产部署
1. **构建可执行文件**:
   ```bash
   # 构建 API 服务器
   go build -o output/api_server cmd/api_server/main.go
   ```

2. **使用生产配置**:
   ```bash
   ./output/api_server -c config.prod.yaml
   ```

3. **Docker 部署** (推荐):
   ```dockerfile
   FROM golang:1.23-alpine AS builder
   WORKDIR /app
   COPY . .
   RUN go build -o api_server cmd/api_server/main.go
   
   FROM alpine:latest
   RUN apk --no-cache add ca-certificates
   WORKDIR /root/
   COPY --from=builder /app/api_server .
   COPY --from=builder /app/config.prod.yaml .
   CMD ["./api_server", "-c", "config.prod.yaml"]
   ```

### 服务监控
- **健康检查**: `GET /ping` 接口返回 "pong"
- **日志记录**: 使用 Go 标准库 `slog` 进行结构化日志记录
- **错误处理**: 统一的错误码和错误信息返回

## 安全考虑

1. **认证安全**:
   - JWT Token 过期机制
   - 安全的密钥管理
   - 请求头验证

2. **数据安全**:
   - 数据库连接加密
   - 敏感信息不记录日志
   - API 参数验证

3. **访问控制**:
   - 用户只能访问自己的数据
   - 基于用户ID的数据隔离

## 扩展性设计

1. **微服务架构**: 支持拆分为独立的微服务
2. **水平扩展**: 无状态设计，支持负载均衡
3. **缓存层**: 可集成 Redis 等缓存系统
4. **消息队列**: 支持异步处理和削峰填谷

## 最佳实践

### 代码组织原则
1. **分层架构**: 严格遵循 Service → Domain → Model 的调用关系
2. **依赖注入**: 使用接口而非具体实现，便于测试和扩展
3. **错误处理**: 使用 Connect-Go 的标准错误码，统一错误格式
4. **日志规范**: 使用结构化日志，包含足够的上下文信息

### 开发规范
1. **Proto 优先**: API 定义从 Proto 文件开始，确保前后端一致
2. **事务管理**: 复杂操作使用 Domain 层的事务封装
3. **ID 生成**: 统一使用雪花算法，确保分布式唯一性
4. **时间处理**: 数据库使用 UTC 时间，前端显示时转换时区

### 性能优化
1. **数据库查询**: 使用索引、预加载、批量操作
2. **并发控制**: 合理使用 goroutine，避免资源竞争
3. **缓存策略**: 热点数据缓存在内存或 Redis
4. **连接池**: HTTP 客户端和数据库连接池合理配置

### 安全实践
1. **认证授权**: 所有 API 都需要 JWT 认证（除登录外）
2. **数据隔离**: 通过 UserID 实现租户级别的数据隔离
3. **参数验证**: 服务层进行完整的参数验证
4. **敏感信息**: 不在日志中记录用户隐私数据

## 故障排查

### 常见问题
1. **数据库连接失败**: 检查配置文件中的 DSN 设置
2. **JWT 验证失败**: 确认密钥配置和 Token 格式
3. **AI API 调用失败**: 检查网络连接和 API 密钥
4. **文件上传失败**: 验证 OSS 配置和权限设置

### 日志级别
- `INFO`: 正常业务流程
- `ERROR`: 错误信息和异常处理
- `DEBUG`: 详细的调试信息

## 项目状态

### 当前版本 (v1.0.0)
- ✅ 完整的用户认证系统
- ✅ 会话和消息管理功能
- ✅ AI 翻译核心功能
- ✅ 文件上传和存储
- ✅ 图片 OCR 解析
- ✅ 生产环境就绪

### 技术特色
- **高性能**: 基于 Go 语言的高并发处理能力
- **现代化架构**: gRPC + Connect 协议，支持 HTTP/2
- **安全可靠**: JWT 认证 + 数据隔离
- **易于扩展**: 模块化设计，便于功能扩展

## 深入阅读

为了更深入了解各个模块的技术细节，请参考以下专题文档：

- **[服务层架构文档](./app_server/service/CLAUDE.md)**: 详细介绍各个服务的实现细节、API 设计和最佳实践
- **[领域层设计文档](./app_server/domain/CLAUDE.md)**: 深入探讨 DDD 设计理念、业务规则实现和事务管理
- **[基础设施层文档](./app_server/pkg/CLAUDE.md)**: 全面介绍各个基础包的功能、使用方法和扩展指南
- **[数据模型层文档](./app_server/model/CLAUDE.md)**: 详细说明数据库设计、ORM 使用和 Proto 转换机制

---

**项目状态**: MVP 完成，具备正式运营的技术基础  
**维护团队**: ChatHandy 开发团队  
**技术栈**: Go + gRPC + MySQL + 火山引擎 AI  
**版本**: v1.0.0  
**最后更新**: 2025-01-02