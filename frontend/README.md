# 恋爱翻译官 - Web 端

AI 驱动的社交辅助应用 Web 版本，帮助你更好地理解异性交流。

## 功能特性

- **手机号验证码登录** - 安全便捷的登录方式
- **性别选择** - 首次登录时选择性别，获得个性化服务
- **会话管理** - 创建、查看、删除聊天会话
- **消息发送** - 支持文本和图片消息
- **AI 翻译** - 一键翻译对方的消息，理解真实意图
- **AI 咨询** - 向 AI 咨询如何回复，获得专业建议
- **Markdown 渲染** - 完美展示 AI 生成的内容
- **响应式设计** - 适配桌面和移动端

## 技术栈

- **框架**: Next.js 14+ (App Router)
- **语言**: TypeScript
- **UI 组件库**: TDesign React
- **样式**: Tailwind CSS
- **状态管理**: Zustand (with persist middleware)
- **HTTP 客户端**: Axios
- **时间处理**: dayjs
- **Markdown 渲染**: react-markdown
- **图标**: TDesign Icons React

## 环境要求

- Node.js >= 18.0.0
- npm >= 9.0.0 或 pnpm >= 8.0.0

## 快速开始

### 1. 安装依赖

```bash
npm install
# 或
pnpm install
```

### 2. 配置环境变量

复制 `.env.local` 文件并根据实际情况修改:

```bash
# 开发环境
NEXT_PUBLIC_API_BASE_URL=https://local.chathandy.com
NEXT_PUBLIC_APP_VERSION=1.0.0
NEXT_PUBLIC_APP_ENV=development
```

生产环境使用 `.env.production`:

```bash
# 生产环境
NEXT_PUBLIC_API_BASE_URL=https://api.chathandy.com
NEXT_PUBLIC_APP_VERSION=1.0.0
NEXT_PUBLIC_APP_ENV=production
```

### 3. 启动开发服务器

```bash
npm run dev
# 或
pnpm dev
```

访问 [http://localhost:3000](http://localhost:3000) 查看应用。

### 4. 构建生产版本

```bash
npm run build
npm run start
# 或
pnpm build
pnpm start
```

## 项目结构

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router 页面
│   │   ├── gender/            # 性别选择页
│   │   ├── login/             # 登录页
│   │   ├── sessions/          # 会话列表和详情
│   │   │   ├── new/          # 新建会话
│   │   │   └── [id]/         # 会话详情 (动态路由)
│   │   ├── globals.css       # 全局样式
│   │   ├── layout.tsx        # 根布局
│   │   └── page.tsx          # 首页 (重定向)
│   │
│   ├── components/            # React 组件
│   │   └── common/           # 通用组件
│   │       └── Loading.tsx   # 加载组件
│   │
│   ├── constants/             # 常量定义
│   │   ├── config.ts         # 应用配置
│   │   ├── errors.ts         # 错误码
│   │   └── routes.ts         # 路由常量
│   │
│   ├── lib/                   # 工具库
│   │   └── avatar.ts         # 头像处理逻辑
│   │
│   ├── services/              # 服务层
│   │   ├── api/              # API 客户端
│   │   │   ├── client.ts    # Axios 实例
│   │   │   ├── auth.ts      # 认证 API
│   │   │   ├── message.ts   # 消息 API
│   │   │   ├── profile.ts   # 用户资料 API
│   │   │   ├── session.ts   # 会话 API
│   │   │   ├── translate.ts # 翻译 API
│   │   │   └── upload.ts    # 上传 API
│   │   │
│   │   └── utils/            # 工具函数
│   │       ├── format.ts    # 格式化函数
│   │       ├── markdown.ts  # Markdown 配置
│   │       ├── storage.ts   # 本地存储封装
│   │       └── validator.ts # 验证函数
│   │
│   ├── stores/                # Zustand 状态管理
│   │   ├── auth.ts           # 认证状态
│   │   ├── ui.ts             # UI 状态
│   │   └── user.ts           # 用户状态
│   │
│   └── types/                 # TypeScript 类型定义
│       ├── api.ts            # API 类型
│       ├── common.ts         # 通用类型
│       └── models.ts         # 数据模型
│
├── public/                    # 静态资源
├── .env.local                # 开发环境变量
├── .env.production           # 生产环境变量
├── next.config.js            # Next.js 配置
├── tailwind.config.js        # Tailwind 配置
├── tsconfig.json             # TypeScript 配置
└── package.json              # 项目依赖

```

## 核心功能说明

### 认证流程

1. 用户输入手机号，点击"发送验证码"
2. 后端发送 6 位验证码到手机
3. 用户输入验证码，点击"登录"
4. 后端验证成功后返回 JWT Token (有效期 365 天)
5. 前端存储 Token 并获取用户信息
6. 如果用户未设置性别，跳转到性别选择页
7. 否则跳转到会话列表页

### 会话管理

- **创建会话**: 输入对方昵称、性别、头像(可选)
- **查看会话**: 显示所有聊天会话，按最后更新时间排序
- **删除会话**: 长按或点击删除按钮删除会话

### 消息功能

- **发送文本**: 输入文本消息并发送
- **发送图片**: 上传图片消息(自动压缩)
- **翻译消息**: 点击对方消息，选择"翻译此消息"
- **AI 咨询**: 输入问题，点击"咨询"按钮获取 AI 建议

### 数据持久化

使用 Zustand persist 中间件，以下数据会持久化到 localStorage:

- **认证 Token**: 自动登录
- **用户信息**: 避免重复请求
- **用户资料**: 性别、头像等

## API 集成

后端 API 基于 gRPC/Connect 协议:

- **认证服务**: `/user.UserService/*`
- **会话服务**: `/chat.ChatService/*`
- **消息服务**: `/message.ChatMessageService/*`
- **翻译服务**: `/translate.TranslateService/*`
- **用户资料**: `/profile.ProfileService/*`

所有请求自动携带:

- `Authorization: Bearer {token}` - JWT Token
- `X-App-Platform: web` - 平台标识
- `Connect-Protocol-Version: 1` - Protocol 版本

## 开发提示

### 开发环境魔法验证码

在开发环境下，可以使用魔法验证码 `1234` 登录任意手机号。

### 图片上传

图片会在客户端自动压缩后上传:

- **头像**: 最大 800x800, 质量 0.8
- **聊天图片**: 最大 1200x1200, 质量 0.85

### 错误处理

- **401 未授权**: 自动清除 Token 并跳转到登录页
- **其他错误**: 显示错误提示信息

### Markdown 支持

AI 返回的内容支持完整的 Markdown 语法:

- 标题 (h1-h3)
- 列表 (有序/无序)
- 引用
- 代码块
- 表格
- 链接

## 常见问题

### Q: 如何修改主题颜色?

A: 编辑 `tailwind.config.js` 中的 `primary` 颜色配置。

### Q: 如何添加新页面?

A: 在 `src/app/` 目录下创建新文件夹和 `page.tsx` 文件。

### Q: 如何调试 API 请求?

A: 检查浏览器开发者工具的 Network 标签，或在 `src/services/api/client.ts` 中添加日志。

### Q: Token 过期怎么办?

A: Token 有效期为 365 天，过期后会自动跳转到登录页重新登录。

## License

Copyright © 2024 恋爱翻译官
