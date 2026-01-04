# Domain 层设计文档

## 🎯 概述

Domain 层是 ChatHandy 应用的核心业务逻辑层，采用领域驱动设计（Domain-Driven Design, DDD）的思想，封装了应用的核心业务规则和领域知识。该层独立于具体的技术实现，专注于表达业务概念和业务规则。

## 🏗️ 领域驱动设计概述

### 设计原则
1. **业务逻辑集中化**：所有核心业务规则都在 domain 层实现，而非分散在 service 或 controller 层
2. **事务一致性**：使用数据库事务确保业务操作的原子性
3. **错误处理标准化**：通过 `exterr` 包提供统一的业务错误定义
4. **领域模型独立性**：领域模型（model 包）独立于外部依赖，仅包含业务属性

### 目录结构
```
domain/
├── chat_session.go      # 聊天会话相关的领域逻辑
├── new_user.go         # 新用户注册和初始化逻辑
├── demo_data.go        # 演示数据初始化逻辑
├── new_user_test.go    # 单元测试
└── exterr/            # 业务错误定义
    └── external_errors.go
```

## 📊 核心领域模型说明

### 1. User（用户）
```go
type User struct {
    gorm.Model
    Name       string  // 用户显示名称
    ImName     string  // IM中的名称
    ExternalId string  // 外部系统ID（如微信OpenID）
    Phone      string  // 手机号
    Avatar     string  // 头像URL
    ProfileID  uint    // 关联的个人资料ID
}
```

**业务规则**：
- 用户可以通过微信（ExternalId）或手机号（Phone）注册
- 每个用户都有一个主 Profile
- 用户名和 IM 名称可以不同，提供灵活性

### 2. Profile（个人资料）
```go
type Profile struct {
    gorm.Model
    UserID          uint         // 所属用户ID
    Name            string       // 姓名
    ImName          string       // IM显示名称
    Avatar          string       // 头像
    Age             int          // 年龄
    Gender          string       // 性别
    Birthday        sql.NullTime // 生日
    BirthLocation   string       // 出生地
    CurrentLocation string       // 当前位置
}
```

**业务规则**：
- 一个用户可以有多个 Profile（代表不同的聊天对象）
- Profile 包含详细的个人信息，用于 AI 分析时提供上下文
- 用户自己也有一个主 Profile

### 3. ChatSession（聊天会话）
```go
type ChatSession struct {
    gorm.Model
    Name      string // 会话名称
    UserID    uint   // 所属用户ID
    ProfileID uint   // 关联的Profile ID
    Avatar    string // 会话头像
}
```

**业务规则**：
- 每个会话关联一个 Profile（聊天对象）
- 会话名称和头像可以自定义，方便用户识别
- 会话是用户级别隔离的

### 4. ConsultMessage（咨询消息）
```go
type ConsultMessage struct {
    gorm.Model
    UserID    uint     // 用户ID
    SessionID uint     // 会话ID
    ParentID  uint     // 父消息ID（用于AI翻译关联）
    ProfileID uint     // 关联的Profile ID
    Role      string   // 角色：SELF/FRIEND/AI
    MsgType   string   // 消息类型：text/translation等
    Content   string   // 消息内容
    Tags      []string // 标签（如demo、translation_to_male等）
    MsgAt     time.Time // 消息时间
}
```

**业务规则**：
- 支持三种角色：自己（SELF）、朋友（FRIEND）、AI助手（AI）
- AI 翻译消息通过 ParentID 关联到原始消息
- Tags 用于标记消息特征，如翻译方向、是否为演示数据等

## 🔧 业务规则和约束

### 1. 用户注册流程
```go
func FindOrRegisterUser(ctx context.Context, user *model.User) error
```
- **查找优先**：先检查用户是否已存在（通过 ExternalId 或 Phone）
- **自动注册**：新用户自动完成注册流程
- **初始化数据**：新用户获得演示数据，帮助理解产品功能
- **事务保证**：整个注册过程在事务中完成，保证数据一致性

### 2. 新用户初始化
```go
func RegisterNewUser(ctx context.Context, tx *gorm.DB, user *model.User) error
```
- **创建主 Profile**：每个用户都有一个代表自己的 Profile
- **复制演示数据**：从 user_id=1 复制演示数据
- **保持关系完整性**：正确映射 Profile、ChatSession 和 Message 的关系
- **演示数据内容**：
  - 3个预设会话（小美、小丽、小芳）
  - 每个会话包含不同场景的对话示例
  - 展示翻译功能的实际效果
  - 包含男女交流的典型场景

### 3. 聊天会话创建
```go
func CreateChatSession(ctx context.Context, userID uint, profile *model.Profile) (*model.ChatSession, error)
```
- **Profile 优先创建**：如果提供了 Profile，先创建 Profile
- **自动关联**：会话自动关联到创建的 Profile
- **事务原子性**：Profile 和 ChatSession 在同一事务中创建

### 4. 演示数据管理
```go
func InitDemoData(ctx context.Context) error
```
- **数据隔离**：演示数据仅属于 user_id=1
- **场景覆盖**：包含初次聊天、暧昧期、情感危机等典型场景
- **翻译示例**：展示 AI 翻译功能的实际效果

## 🚨 错误处理策略

### 错误类型定义
```go
type ExternalError struct {
    Code  int32  // 错误码
    Msg   string // 英文错误信息
    CnMsg string // 中文错误信息
}
```

### 预定义错误
- `ErrInternal (10000)`：内部错误，通常是数据库或系统级错误
- `ErrParamMissing (20001)`：必填参数缺失
- `ErrParamInvalid (20002)`：参数格式或值无效

### 错误处理原则
1. **分层处理**：domain 层定义业务错误，service 层转换为 API 错误
2. **日志记录**：使用 slog 记录详细错误信息，方便调试
3. **用户友好**：提供中文错误信息，提升用户体验
4. **事务回滚**：错误发生时自动回滚事务，保证数据一致性

## 💼 领域服务的实现

### 1. 事务管理模式
```go
err = db.GetDB().Transaction(func(tx *gorm.DB) error {
    // 业务逻辑
    return nil
})
```
- 使用 GORM 的事务回调，自动处理提交和回滚
- 支持 panic 恢复，防止程序崩溃
- 统一的错误处理和日志记录

### 2. 数据复制策略
演示数据复制采用批量操作和 ID 映射：
- **批量创建**：减少数据库交互次数
- **ID 映射表**：维护新旧 ID 的对应关系
- **关系重建**：正确处理 ParentID 等外键关系
- **智能去重**：避免重复复制演示数据
- **标签管理**：为演示数据添加 "demo" 标签

### 3. 查询优化
- 使用预加载避免 N+1 查询问题
- 合理使用索引提升查询性能
- 批量操作减少数据库往返

## 🔄 与其他层的交互方式

### 1. 与 Model 层
- Domain 层依赖 Model 定义的数据结构
- Model 提供与 Protobuf 的转换方法
- Domain 负责维护模型间的业务关系

### 2. 与 Service 层
- Service 层调用 Domain 层的业务方法
- Domain 层返回领域对象或业务错误
- Service 层负责将结果转换为 API 响应

### 3. 与数据库层
- 通过 `pkg/db` 包获取数据库连接
- 使用 GORM 进行 ORM 操作
- 事务管理确保数据一致性

## 🧪 测试策略

### 单元测试设计
1. **使用内存数据库**：SQLite 内存模式，快速且隔离
2. **完整场景覆盖**：测试正常流程和异常情况
3. **数据验证**：验证业务规则是否正确执行

### 测试重点
- 新用户注册和数据初始化
- 演示数据复制的正确性
- ParentID 映射的准确性
- 事务回滚的有效性

## 📈 性能考虑

1. **批量操作**：演示数据复制使用批量插入
2. **延迟加载**：仅在需要时加载关联数据
3. **事务范围**：合理控制事务大小，避免长事务
4. **错误快速失败**：参数验证提前进行，减少无效操作

## 🔐 安全考虑

1. **用户隔离**：所有查询都包含 UserID 条件
2. **参数验证**：防止 SQL 注入和无效输入
3. **权限检查**：确保用户只能访问自己的数据
4. **敏感信息**：日志中避免记录用户隐私数据

## 🚀 未来优化方向

1. **领域事件**：引入领域事件机制，解耦业务流程
2. **聚合根**：明确聚合根边界，提升模型一致性
3. **仓储模式**：抽象数据访问层，提高可测试性
4. **缓存策略**：对热点数据进行缓存，提升性能

---

**最后更新**: 2025-01-07  
**维护者**: ChatHandy 开发团队