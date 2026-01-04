# ChatHandy 服务端基础设施层 (pkg)

## 📋 概述

`pkg` 目录包含了 ChatHandy 后端服务的所有基础设施和公共工具包。这些包提供了核心功能支持，包括 AI 服务集成、认证、数据库、文件存储、HTTP 客户端等基础设施组件。

## 🏗️ 包结构图

```
pkg/
├── aiapi/        # AI API 集成（火山引擎 OCR）
├── cbind/        # Connect-Go 与 Gin 的绑定器
├── cfg/          # 配置管理（基于 Viper）
├── db/           # 数据库连接和 GORM 配置
├── fn/           # 函数式编程工具集
├── httpc/        # HTTP 客户端封装
├── idgen/        # 分布式 ID 生成（雪花算法）
├── jwt/          # JWT 认证管理
├── openaic/      # OpenAI 客户端封装
└── ossc/         # 阿里云 OSS 客户端封装
```

## 📦 各包详细说明

### 1. aiapi - AI API 集成包

**功能**: 提供 AI 相关功能的高级封装，目前主要用于图片聊天记录的 OCR 解析。

**核心功能**:
- `ParseImageChat(imageUrl string) ([]string, error)`: 解析图片中的聊天内容
  - 使用火山引擎的 OCR 模型
  - 自动识别聊天记录格式
  - 返回格式化的聊天行列表
- `ParseChatLine(line string) (role string, content string, ok bool)`: 解析单行聊天记录
  - 识别【朋友】和【自己】角色
  - 返回标准化的角色标识（FRIEND/SELF）

**使用示例**:
```go
// 解析聊天截图
lines, err := aiapi.ParseImageChat("https://example.com/chat.png")
// 返回: ["【朋友】你好", "【自己】你好呀"]

// 解析单行
role, content, ok := aiapi.ParseChatLine("【朋友】你好")
// role: "FRIEND", content: "你好", ok: true
```

### 2. cbind - Connect-Go 绑定器

**功能**: 提供 Connect-Go (gRPC) 与 Gin Web 框架的集成绑定。

**核心组件**:
- `Binder`: 将 Connect-Go 的 HTTP Handler 绑定到 Gin 路由
- 自动处理路径模式和 POST 请求映射

**使用示例**:
```go
router := gin.New()
binder := cbind.NewBinder(router.Group("/api"))
binder.Bind("/user.UserService/Login", userServiceHandler)
```

### 3. cfg - 配置管理

**功能**: 基于 Viper 的配置管理，支持 YAML 配置文件。

**核心功能**:
- `Init(file string)`: 初始化配置文件
- `Viper() *viper.Viper`: 获取 Viper 实例
- `UnmarshalKey[T any](key string) T`: 泛型配置解析

**使用示例**:
```go
// 初始化配置
cfg.Init("config.yaml")

// 解析特定配置
dbConfig := cfg.UnmarshalKey[DatabaseConfig]("database")
```

### 4. db - 数据库管理

**功能**: 提供 MySQL 数据库连接管理和 GORM 配置。

**核心特性**:
- 基于 GORM 的 ORM 封装
- 自动 ID 生成（使用雪花算法）
- 单表命名策略
- 支持测试环境的数据库替换

**自动 ID 生成机制**:
- 在创建记录前自动设置 ID
- 支持批量插入
- 兼容 int/uint 类型的 ID 字段

**使用示例**:
```go
// 初始化数据库
err := db.Init("user:pass@tcp(localhost:3306)/chathandy?charset=utf8mb4")

// 获取数据库实例
database := db.GetDB()

// 创建记录（ID 自动生成）
user := &User{Name: "张三"}
database.Create(user) // user.ID 自动填充
```

### 5. fn - 函数式编程工具集

**功能**: 提供泛型函数式编程工具，简化数据转换和处理。

#### 5.1 类型转换工具 (cast.go)
- `Atoi[T constraints.Integer](s string) T`: 字符串转整数
- `Itoa[T constraints.Integer](v T) string`: 整数转字符串
- `CastNumber[A, B](a A) B`: 数字类型转换
- `CastNumbers[A, B](a []A) []B`: 批量数字类型转换

#### 5.2 错误处理工具 (error.go)
- `NoErr[T any](v T, err error) T`: 忽略错误，返回值

#### 5.3 函数式编程工具 (functional.go)
- `Map[T, R](arr []T, fn func(T) R) []R`: 映射转换
- `Filter[T](arr []T, fn func(T) bool) []T`: 过滤操作
- `Partial[T1, T2, R](f func(T1, T2) R, arg1 T1) func(T2) R`: 偏函数（固定第一个参数）
- `PartialR[T1, T2, R](f func(T1, T2) R, arg2 T2) func(T1) R`: 偏函数（固定第二个参数）
- `DropLast[P, R1, R2](f func(P) (R1, R2)) func(P) R1`: 丢弃第二个返回值
- `Ptr[T](item T) *T`: 值转指针
- `DerefOr0[T](item *T) T`: 安全解引用

#### 5.4 JSON 工具 (json.go)
- `JsonUnmarshalStr[T any](data string) (T, error)`: 泛型 JSON 解析

**使用示例**:
```go
// Map 示例
ids := []string{"1", "2", "3"}
numbers := fn.Map(ids, fn.Atoi[int])
// numbers: []int{1, 2, 3}

// Filter 示例
positive := fn.Filter(numbers, func(n int) bool { return n > 0 })

// 偏函数示例
add10 := fn.Partial(func(a, b int) int { return a + b }, 10)
result := add10(5) // 15
```

### 6. httpc - HTTP 客户端

**功能**: 提供统一的 HTTP 客户端实例（基于 req/v3）。

**使用示例**:
```go
client := httpc.Client()
resp, err := client.R().Get("https://api.example.com/data")
```

### 7. idgen - ID 生成器

**功能**: 基于雪花算法的分布式 ID 生成器。

**核心特性**:
- 64 位分布式唯一 ID
- 支持多种格式输出（int64、Base36）
- 时间戳解析功能
- 随机节点 ID（0-1023）

**API 接口**:
- `NewID() ID`: 生成新 ID
- `Int() int`: 直接生成 int 类型 ID
- `Base36() string`: 生成 Base36 字符串 ID
- `FromInt(idint int) ID`: 从整数创建 ID
- `FromBase36(idstring string) (ID, error)`: 从 Base36 字符串解析 ID
- `FromTime(t time.Time) int64`: 从时间生成 ID 前缀
- `ToTime(id int64) time.Time`: 从 ID 解析时间

**使用示例**:
```go
// 生成新 ID
id := idgen.NewID()
fmt.Println(id.Int())    // 7123456789012345678
fmt.Println(id.String()) // "1a2b3c4d5e"

// 直接生成
intId := idgen.Int()
strId := idgen.Base36()

// 时间相关
idFromTime := idgen.FromTime(time.Now())
timeFromId := idgen.ToTime(id.Int64())
```

### 8. jwt - JWT 认证管理

**功能**: 提供 JWT 令牌的生成和解析功能。

**核心功能**:
- `GenerateToken(userID string, exp time.Duration) (string, error)`: 生成令牌
- `ParseToken(tokenString string) (userID string, err error)`: 解析令牌
- 使用 HS256 签名算法
- 支持过期时间设置

**使用示例**:
```go
// 初始化
jwt.Init([]byte("your-secret-key"))

// 生成令牌（有效期 7 天）
token, err := jwt.Get().GenerateToken("user123", 7*24*time.Hour)

// 解析令牌
userID, err := jwt.Get().ParseToken(token)
```

### 9. openaic - OpenAI 客户端

**功能**: 封装 OpenAI API 客户端，支持火山引擎等兼容接口。

**配置结构**:
```go
type Config struct {
    ApiKey  string // API 密钥
    BaseURL string // API 基础 URL
    Models  Models // 模型配置
}

type Models struct {
    Chat string // 聊天模型
    Ocr  string // OCR 模型
}
```

**使用示例**:
```go
// 初始化
openaic.Init(openaic.Config{
    ApiKey:  "your-api-key",
    BaseURL: "https://api.volcengine.com/v1",
    Models: openaic.Models{
        Chat: "gpt-3.5-turbo",
        Ocr:  "gpt-4-vision",
    },
})

// 获取客户端
client := openaic.Get()
```

### 10. ossc - 阿里云 OSS 客户端

**功能**: 封装阿里云 OSS 存储服务，提供文件上传下载功能。

**核心特性**:
- 双客户端模式（内网/公网）
- 预配置的存储桶访问
- 简化的 API 接口

**配置结构**:
```go
type Cfg struct {
    PublicEndpoint  string // 公网访问端点
    Endpoint        string // 内网访问端点
    AccessKeyId     string // 访问密钥 ID
    AccessKeySecret string // 访问密钥
    UserFileBucket  string // 用户文件存储桶
}
```

**使用示例**:
```go
// 初始化
err := ossc.Init(ossc.Cfg{
    PublicEndpoint:  "oss-cn-beijing.aliyuncs.com",
    Endpoint:        "oss-cn-beijing-internal.aliyuncs.com",
    AccessKeyId:     "your-key-id",
    AccessKeySecret: "your-key-secret",
    UserFileBucket:  "chathandy-files",
})

// 获取客户端
client := ossc.Get()           // 内网客户端
publicClient := ossc.GetPublic() // 公网客户端

// 获取存储桶
bucket := client.UserFileBucket()
```

## 🔧 最佳实践

### 1. 初始化顺序
建议按以下顺序初始化各个包：
```go
// 1. 配置管理
cfg.Init("config.yaml")

// 2. 数据库
db.Init(cfg.Viper().GetString("database.dsn"))

// 3. JWT
jwt.Init([]byte(cfg.Viper().GetString("jwt.secret")))

// 4. AI 客户端
openaic.Init(cfg.UnmarshalKey[openaic.Config]("openai"))

// 5. OSS
ossc.Init(cfg.UnmarshalKey[ossc.Cfg]("oss"))
```

### 2. 错误处理
- 使用 `fn.NoErr` 仅在确定不会出错的场景
- 优先使用泛型工具函数减少类型断言
- 在初始化阶段使用 `lo.Must` 确保配置正确

### 3. ID 生成策略
- 所有实体 ID 使用 `idgen` 生成
- 数据库会自动为新记录生成 ID
- 对外暴露使用 Base36 格式，内部使用 int 格式

### 4. 函数式编程
- 使用 `fn.Map` 进行批量转换
- 使用 `fn.Filter` 进行条件筛选
- 组合使用偏函数简化复杂操作

## 🚨 注意事项

1. **安全性**
   - JWT 密钥必须足够复杂且定期更换
   - OSS 密钥应通过环境变量或密钥管理服务提供
   - 避免在日志中打印敏感信息

2. **性能优化**
   - HTTP 客户端已预配置连接池
   - 数据库连接池由 GORM 自动管理
   - ID 生成器使用本地缓存，性能极高

3. **错误处理**
   - 所有初始化函数都应检查错误
   - 使用结构化日志记录错误上下文
   - 对外接口统一错误格式

4. **配置管理**
   - 敏感配置使用环境变量覆盖
   - 开发/测试/生产环境使用不同配置文件
   - 配置变更需要重启服务

## 📈 扩展指南

### 添加新的基础设施包
1. 在 `pkg/` 下创建新目录
2. 实现初始化函数 `Init()`
3. 提供全局访问器 `Get()`
4. 在主程序初始化流程中调用
5. 更新本文档

### 集成新的 AI 服务
1. 在 `aiapi/` 中添加新的服务文件
2. 实现标准化的接口
3. 在配置中添加相应配置项
4. 提供降级和容错机制

---

**最后更新**: 2025-01-02  
**维护者**: ChatHandy 开发团队