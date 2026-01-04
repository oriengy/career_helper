# Web 端重构实现方案（详细版）

## 📋 目录
- [1. 项目架构概述](#1-项目架构概述)
- [2. 登录与鉴权模块](#2-登录与鉴权模块)
- [3. 工具函数层](#3-工具函数层)
- [4. 会话列表页面](#4-会话列表页面)
- [5. 聊天详情页面](#5-聊天详情页面)
- [6. 资料编辑页面](#6-资料编辑页面)
- [7. 性别选择流程](#7-性别选择流程)
- [8. 状态管理方案](#8-状态管理方案)
- [9. 错误处理与边界情况](#9-错误处理与边界情况)
- [10. 部署与配置](#10-部署与配置)

---

## 1. 项目架构概述

### 1.1 整体技术栈

**前端框架**: Next.js 14+ (App Router)
**UI 组件库**: TDesign React
**状态管理**: Zustand
**HTTP 客户端**: Axios
**样式方案**: CSS Modules + Tailwind CSS
**Markdown 渲染**: react-markdown
**时间处理**: dayjs
**表单管理**: React Hook Form
**类型检查**: TypeScript

### 1.2 项目目录结构

```
career_helper/frontend/
├── src/
│   ├── app/                          # Next.js App Router
│   │   ├── layout.tsx                # 根布局
│   │   ├── page.tsx                  # 首页（重定向到 /sessions）
│   │   ├── sessions/                 # 会话模块
│   │   │   ├── page.tsx              # 会话列表页
│   │   │   ├── new/page.tsx          # 创建会话页
│   │   │   └── [id]/page.tsx         # 聊天详情页
│   │   ├── profile/                  # 资料模块
│   │   │   ├── page.tsx              # 个人资料页
│   │   │   ├── edit/page.tsx         # 编辑自己资料
│   │   │   └── [id]/edit/page.tsx    # 编辑好友资料
│   │   ├── login/                    # 登录页
│   │   │   └── page.tsx
│   │   └── gender/                   # 性别选择页
│   │       └── page.tsx
│   │
│   ├── components/                   # UI 组件
│   │   ├── layout/                   # 布局组件
│   │   │   ├── Header.tsx
│   │   │   ├── Footer.tsx
│   │   │   └── MainLayout.tsx
│   │   ├── sessions/                 # 会话相关组件
│   │   │   ├── SessionList.tsx
│   │   │   ├── SessionCard.tsx
│   │   │   └── CreateSessionModal.tsx
│   │   ├── chat/                     # 聊天相关组件
│   │   │   ├── MessageList.tsx
│   │   │   ├── MessageItem.tsx
│   │   │   ├── MessageInput.tsx
│   │   │   ├── TranslateResult.tsx
│   │   │   └── ImageUpload.tsx
│   │   ├── profile/                  # 资料相关组件
│   │   │   ├── ProfileForm.tsx
│   │   │   ├── AvatarUpload.tsx
│   │   │   └── GenderSelector.tsx
│   │   └── common/                   # 通用组件
│   │       ├── Loading.tsx
│   │       ├── ErrorBoundary.tsx
│   │       ├── MarkdownRenderer.tsx
│   │       └── Avatar.tsx
│   │
│   ├── services/                     # 业务逻辑服务层
│   │   ├── api/                      # API 层
│   │   │   ├── client.ts             # Axios 客户端配置
│   │   │   ├── auth.ts               # 认证相关 API
│   │   │   ├── session.ts            # 会话相关 API
│   │   │   ├── message.ts            # 消息相关 API
│   │   │   ├── translate.ts          # 翻译相关 API
│   │   │   ├── profile.ts            # 资料相关 API
│   │   │   └── upload.ts             # 文件上传 API
│   │   └── utils/                    # 业务工具函数
│   │       ├── auth.ts               # 认证工具
│   │       ├── storage.ts            # 本地存储工具
│   │       └── format.ts             # 格式化工具
│   │
│   ├── stores/                       # Zustand 状态管理
│   │   ├── auth.ts                   # 认证状态
│   │   ├── user.ts                   # 用户状态
│   │   ├── session.ts                # 会话状态
│   │   └── ui.ts                     # UI 状态
│   │
│   ├── types/                        # TypeScript 类型定义
│   │   ├── api.ts                    # API 响应类型
│   │   ├── models.ts                 # 数据模型类型
│   │   └── common.ts                 # 通用类型
│   │
│   ├── hooks/                        # 自定义 Hooks
│   │   ├── useAuth.ts                # 认证相关 Hook
│   │   ├── useSession.ts             # 会话相关 Hook
│   │   ├── useMessage.ts             # 消息相关 Hook
│   │   └── useInfiniteScroll.ts      # 无限滚动 Hook
│   │
│   ├── lib/                          # 第三方库配置
│   │   ├── axios.ts                  # Axios 配置
│   │   └── markdown.ts               # Markdown 配置
│   │
│   ├── constants/                    # 常量定义
│   │   ├── api.ts                    # API 常量
│   │   ├── routes.ts                 # 路由常量
│   │   └── config.ts                 # 配置常量
│   │
│   └── styles/                       # 全局样式
│       ├── globals.css               # 全局 CSS
│       └── variables.css             # CSS 变量
│
├── public/                           # 静态资源
│   └── assets/
│       └── images/
│           ├── male.png              # 男性默认头像
│           ├── female.png            # 女性默认头像
│           └── assistant.png         # 通用默认头像
│
├── .env.local                        # 环境变量（开发）
├── .env.production                   # 环境变量（生产）
├── next.config.js                    # Next.js 配置
├── tailwind.config.js                # Tailwind 配置
├── tsconfig.json                     # TypeScript 配置
└── package.json                      # 项目依赖
```

### 1.3 技术架构图

```
┌─────────────────────────────────────────────────────┐
│                   用户界面层                         │
│            (Next.js Pages & Components)              │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                  状态管理层                          │
│                   (Zustand Stores)                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐          │
│  │ Auth     │  │ User     │  │ Session  │          │
│  └──────────┘  └──────────┘  └──────────┘          │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                  服务层                              │
│              (API Services)                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐          │
│  │ Auth API │  │ Chat API │  │ File API │          │
│  └──────────┘  └──────────┘  └──────────┘          │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                 HTTP 客户端                          │
│            (Axios with Interceptors)                 │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│                   后端服务                           │
│              (Go Backend Server)                     │
└─────────────────────────────────────────────────────┘
```

---

## 2. 登录与鉴权模块

### 2.1 模块概述

本模块负责用户身份认证、Token 管理和登录状态维护，是整个应用的安全基础。

**核心职责**:
- 用户登录（手机号登录为主）
- Token 的获取、存储和刷新
- 登录状态检查和维护
- 401 错误自动重试
- 自动登录保持

### 2.2 与小程序的差异对比

| 功能点 | 小程序实现 | Web 端实现 | 差异说明 |
|--------|-----------|-----------|---------|
| 登录方式 | wx.login() 获取 code | 手机号 + 验证码 | 小程序使用微信授权，Web 使用手机号 |
| Token 存储 | wx.setStorageSync() | localStorage + sessionStorage | 存储 API 不同，Web 需考虑跨标签页同步 |
| 自动登录 | 启动时自动检查 token | 页面加载时检查 token | 机制相同，触发时机略有不同 |
| 401 处理 | 清除 token 后重新 wx.login | 清除 token 后跳转登录页 | Web 需要明确的登录页面 |
| 登录锁 | 使用 Promise 防重复登录 | 使用 Zustand 状态 + Promise | Web 使用状态管理更清晰 |

### 2.3 详细技术方案

#### 2.3.1 认证流程图

```
┌─────────────┐
│ 用户访问页面 │
└──────┬──────┘
       │
       ▼
┌──────────────────┐
│ 检查本地 token    │
│ (localStorage)   │
└──────┬───────────┘
       │
       ├─────────────┐
       │             │
   token 存在     token 不存在
       │             │
       ▼             ▼
┌──────────────┐  ┌─────────────┐
│ 验证 token    │  │ 跳转登录页   │
│ 有效性       │  └─────────────┘
└──────┬───────┘
       │
       ├─────────────┐
       │             │
    有效          无效
       │             │
       ▼             ▼
┌──────────────┐  ┌──────────────┐
│ 加载用户数据  │  │ 清除 token   │
└──────────────┘  │ 跳转登录页   │
                  └──────────────┘
```

#### 2.3.2 登录页面实现方案

**页面路径**: `src/app/login/page.tsx`

**功能需求**:
1. 手机号输入（带格式校验）
2. 验证码输入（6位数字）
3. 发送验证码按钮（60秒倒计时）
4. 登录按钮
5. 登录中状态展示
6. 错误提示

**UI 结构**:
```tsx
<LoginPage>
  <Logo />
  <Title>欢迎使用恋爱翻译官</Title>

  <Form>
    <PhoneInput
      placeholder="请输入手机号"
      validation={phoneValidator}
    />

    <VerifyCodeGroup>
      <VerifyCodeInput
        placeholder="请输入验证码"
        maxLength={6}
      />
      <SendCodeButton
        disabled={!phoneValid || countdown > 0}
        onClick={sendVerifyCode}
      >
        {countdown > 0 ? `${countdown}秒后重试` : '发送验证码'}
      </SendCodeButton>
    </VerifyCodeGroup>

    <LoginButton
      loading={isLoggingIn}
      disabled={!canSubmit}
      onClick={handleLogin}
    >
      登录
    </LoginButton>
  </Form>

  <Agreement>
    登录即表示同意《用户协议》和《隐私政策》
  </Agreement>
</LoginPage>
```

**状态管理**:
```typescript
interface LoginPageState {
  phone: string;              // 手机号
  verifyCode: string;         // 验证码
  countdown: number;          // 倒计时秒数
  isLoggingIn: boolean;       // 登录中
  isSendingCode: boolean;     // 发送验证码中
  error: string | null;       // 错误信息
}
```

**关键逻辑**:
```typescript
// 发送验证码
const sendVerifyCode = async () => {
  // 1. 验证手机号格式
  if (!validatePhone(phone)) {
    showError('请输入正确的手机号');
    return;
  }

  // 2. 调用发送验证码 API
  try {
    setIsSendingCode(true);
    await authApi.sendVerifyCode({ phone });

    // 3. 开始60秒倒计时
    setCountdown(60);
    const timer = setInterval(() => {
      setCountdown(prev => {
        if (prev <= 1) {
          clearInterval(timer);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    showSuccess('验证码已发送');
  } catch (error) {
    showError('发送验证码失败，请重试');
  } finally {
    setIsSendingCode(false);
  }
};

// 登录处理
const handleLogin = async () => {
  // 1. 验证输入
  if (!phone || !verifyCode) {
    showError('请填写完整信息');
    return;
  }

  // 2. 调用登录 API
  try {
    setIsLoggingIn(true);
    const response = await authApi.phoneLogin({
      phone,
      code: verifyCode
    });

    // 3. 保存 token
    const { token, user, profile } = response;
    authStore.setToken(token);
    userStore.setUser(user);
    userStore.setProfile(profile);

    // 4. 判断是否需要性别选择
    if (!profile.gender) {
      router.push('/gender');
    } else {
      router.push('/sessions');
    }
  } catch (error) {
    showError(error.message || '登录失败');
  } finally {
    setIsLoggingIn(false);
  }
};
```

#### 2.3.3 Auth API 服务实现

**文件位置**: `src/services/api/auth.ts`

**API 接口定义**:
```typescript
interface AuthApi {
  // 发送验证码
  sendVerifyCode(params: SendVerifyCodeParams): Promise<void>;

  // 手机号登录
  phoneLogin(params: PhoneLoginParams): Promise<PhoneLoginResponse>;

  // 获取用户信息（验证 token 有效性）
  getUserProfile(): Promise<UserProfileResponse>;

  // 登出
  logout(): Promise<void>;
}

// 请求参数类型
interface SendVerifyCodeParams {
  phone: string;
}

interface PhoneLoginParams {
  phone: string;
  code: string;
}

// 响应类型
interface PhoneLoginResponse {
  token: string;
  user: User;
  profile: Profile;
}

interface UserProfileResponse {
  user: User;
  profile: Profile;
}
```

**实现代码**:
```typescript
import { apiClient } from './client';

export const authApi = {
  // 发送验证码
  async sendVerifyCode(params: SendVerifyCodeParams) {
    // 注意：这个接口可能需要后端新增
    await apiClient.post('/user.UserService/SendVerifyCode', params);
  },

  // 手机号登录
  async phoneLogin(params: PhoneLoginParams) {
    const response = await apiClient.post<PhoneLoginResponse>(
      '/user.UserService/PhoneLogin',
      params
    );
    return response.data;
  },

  // 获取用户信息
  async getUserProfile() {
    const response = await apiClient.post<UserProfileResponse>(
      '/user.UserService/GetUserProfile',
      {}
    );
    return response.data;
  },

  // 登出（清理客户端状态）
  async logout() {
    // 后端可能不需要 logout 接口，主要是清理本地状态
    return Promise.resolve();
  }
};
```

#### 2.3.4 Auth Store 状态管理

**文件位置**: `src/stores/auth.ts`

**状态定义**:
```typescript
interface AuthState {
  // 状态
  token: string | null;
  isLoggedIn: boolean;
  isLoggingIn: boolean;

  // 操作
  setToken: (token: string) => void;
  clearToken: () => void;
  checkAuth: () => Promise<boolean>;
  login: (phone: string, code: string) => Promise<void>;
  logout: () => Promise<void>;
}
```

**完整实现**:
```typescript
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { authApi } from '@/services/api/auth';

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // 初始状态
      token: null,
      isLoggedIn: false,
      isLoggingIn: false,

      // 设置 token
      setToken: (token) => {
        set({ token, isLoggedIn: true });
        // 同时存储到 localStorage
        localStorage.setItem('access_token', token);
      },

      // 清除 token
      clearToken: () => {
        set({ token: null, isLoggedIn: false });
        localStorage.removeItem('access_token');
      },

      // 检查认证状态
      checkAuth: async () => {
        const { token } = get();
        if (!token) {
          return false;
        }

        try {
          // 调用获取用户信息接口验证 token
          await authApi.getUserProfile();
          set({ isLoggedIn: true });
          return true;
        } catch (error) {
          // token 无效，清除
          get().clearToken();
          return false;
        }
      },

      // 登录
      login: async (phone, code) => {
        set({ isLoggingIn: true });
        try {
          const response = await authApi.phoneLogin({ phone, code });
          get().setToken(response.token);

          // 更新用户信息到 userStore
          const userStore = useUserStore.getState();
          userStore.setUser(response.user);
          userStore.setProfile(response.profile);
        } finally {
          set({ isLoggingIn: false });
        }
      },

      // 登出
      logout: async () => {
        await authApi.logout();
        get().clearToken();

        // 清除用户信息
        const userStore = useUserStore.getState();
        userStore.clearUser();
      }
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ token: state.token })
    }
  )
);
```

#### 2.3.5 Axios 拦截器配置

**文件位置**: `src/services/api/client.ts`

**请求拦截器**:
```typescript
import axios from 'axios';
import { useAuthStore } from '@/stores/auth';

// 创建 axios 实例
export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
    'Connect-Protocol-Version': '1'
  }
});

// 请求拦截器：自动添加 token
apiClient.interceptors.request.use(
  (config) => {
    // 获取 token
    const token = useAuthStore.getState().token;

    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // 添加自定义 headers
    config.headers['X-App-Platform'] = 'web';
    config.headers['X-App-Env'] = process.env.NODE_ENV;
    config.headers['X-App-Version'] = process.env.NEXT_PUBLIC_APP_VERSION || '1.0.0';

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器：处理 401 错误
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config;

    // 401 错误：token 失效
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      // 清除 token
      const authStore = useAuthStore.getState();
      authStore.clearToken();

      // 跳转到登录页
      if (typeof window !== 'undefined') {
        window.location.href = '/login';
      }

      return Promise.reject(error);
    }

    return Promise.reject(error);
  }
);
```

#### 2.3.6 路由守卫实现

**文件位置**: `src/components/layout/AuthGuard.tsx`

**功能需求**:
- 检查用户是否已登录
- 未登录自动跳转登录页
- 白名单路由（登录页、注册页等）不需要认证

**实现代码**:
```typescript
'use client';

import { useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAuthStore } from '@/stores/auth';
import Loading from '@/components/common/Loading';

// 不需要认证的路由白名单
const PUBLIC_ROUTES = ['/login', '/register'];

interface AuthGuardProps {
  children: React.ReactNode;
}

export default function AuthGuard({ children }: AuthGuardProps) {
  const router = useRouter();
  const pathname = usePathname();
  const { isLoggedIn, checkAuth } = useAuthStore();
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    const verify = async () => {
      // 白名单路由直接放行
      if (PUBLIC_ROUTES.includes(pathname)) {
        setIsChecking(false);
        return;
      }

      // 检查登录状态
      const isValid = await checkAuth();

      if (!isValid) {
        // 未登录，跳转登录页
        router.push(`/login?redirect=${encodeURIComponent(pathname)}`);
      } else {
        setIsChecking(false);
      }
    };

    verify();
  }, [pathname]);

  // 检查中显示 loading
  if (isChecking) {
    return <Loading fullScreen />;
  }

  return <>{children}</>;
}
```

**在根布局中使用**:
```typescript
// src/app/layout.tsx
import AuthGuard from '@/components/layout/AuthGuard';

export default function RootLayout({ children }) {
  return (
    <html lang="zh-CN">
      <body>
        <AuthGuard>
          {children}
        </AuthGuard>
      </body>
    </html>
  );
}
```

#### 2.3.7 自动登录保持

**使用场景**:
- 用户刷新页面
- 用户关闭后重新打开浏览器
- 多标签页同步登录状态

**实现方案**:

**1. 本地存储策略**:
```typescript
// src/services/utils/storage.ts

export const storage = {
  // Token 存储（持久化）
  setToken(token: string) {
    localStorage.setItem('access_token', token);
  },

  getToken(): string | null {
    return localStorage.getItem('access_token');
  },

  removeToken() {
    localStorage.removeItem('access_token');
  },

  // 用户信息存储（持久化）
  setUser(user: User) {
    localStorage.setItem('user_info', JSON.stringify(user));
  },

  getUser(): User | null {
    const data = localStorage.getItem('user_info');
    return data ? JSON.parse(data) : null;
  },

  // Profile 存储（持久化）
  setProfile(profile: Profile) {
    localStorage.setItem('user_profile', JSON.stringify(profile));
  },

  getProfile(): Profile | null {
    const data = localStorage.getItem('user_profile');
    return data ? JSON.parse(data) : null;
  },

  // 清除所有数据
  clear() {
    localStorage.removeItem('access_token');
    localStorage.removeItem('user_info');
    localStorage.removeItem('user_profile');
  }
};
```

**2. 多标签页同步**:
```typescript
// src/hooks/useStorageSync.ts

export function useStorageSync() {
  const { setToken, clearToken } = useAuthStore();

  useEffect(() => {
    // 监听 storage 变化（其他标签页修改）
    const handleStorageChange = (e: StorageEvent) => {
      if (e.key === 'access_token') {
        if (e.newValue) {
          setToken(e.newValue);
        } else {
          clearToken();
        }
      }
    };

    window.addEventListener('storage', handleStorageChange);

    return () => {
      window.removeEventListener('storage', handleStorageChange);
    };
  }, []);
}
```

#### 2.3.8 登录状态持久化

**使用 Zustand persist 中间件**:
```typescript
// 已在 authStore 中配置
persist(
  (set, get) => ({
    // store 实现
  }),
  {
    name: 'auth-storage',          // localStorage key
    partialize: (state) => ({       // 只持久化这些字段
      token: state.token
    })
  }
)
```

### 2.4 错误处理方案

#### 2.4.1 错误类型定义

```typescript
// src/types/errors.ts

export enum AuthErrorCode {
  INVALID_PHONE = 'INVALID_PHONE',
  INVALID_CODE = 'INVALID_CODE',
  CODE_EXPIRED = 'CODE_EXPIRED',
  TOO_MANY_REQUESTS = 'TOO_MANY_REQUESTS',
  TOKEN_EXPIRED = 'TOKEN_EXPIRED',
  TOKEN_INVALID = 'TOKEN_INVALID',
  NETWORK_ERROR = 'NETWORK_ERROR',
  UNKNOWN_ERROR = 'UNKNOWN_ERROR'
}

export class AuthError extends Error {
  code: AuthErrorCode;

  constructor(code: AuthErrorCode, message: string) {
    super(message);
    this.code = code;
    this.name = 'AuthError';
  }
}
```

#### 2.4.2 错误提示映射

```typescript
// src/constants/errors.ts

export const AUTH_ERROR_MESSAGES: Record<AuthErrorCode, string> = {
  [AuthErrorCode.INVALID_PHONE]: '手机号格式不正确',
  [AuthErrorCode.INVALID_CODE]: '验证码错误',
  [AuthErrorCode.CODE_EXPIRED]: '验证码已过期，请重新获取',
  [AuthErrorCode.TOO_MANY_REQUESTS]: '请求过于频繁，请稍后再试',
  [AuthErrorCode.TOKEN_EXPIRED]: '登录已过期，请重新登录',
  [AuthErrorCode.TOKEN_INVALID]: '登录状态异常，请重新登录',
  [AuthErrorCode.NETWORK_ERROR]: '网络连接失败，请检查网络',
  [AuthErrorCode.UNKNOWN_ERROR]: '操作失败，请重试'
};
```

### 2.5 测试方案

#### 2.5.1 单元测试

```typescript
// __tests__/auth.test.ts

describe('Auth Store', () => {
  it('should set token correctly', () => {
    const { setToken, token } = useAuthStore.getState();
    setToken('test-token');
    expect(token).toBe('test-token');
  });

  it('should clear token correctly', () => {
    const { setToken, clearToken, token } = useAuthStore.getState();
    setToken('test-token');
    clearToken();
    expect(token).toBeNull();
  });
});
```

#### 2.5.2 集成测试

```typescript
// __tests__/login.test.tsx

describe('Login Page', () => {
  it('should login successfully', async () => {
    render(<LoginPage />);

    // 输入手机号
    const phoneInput = screen.getByPlaceholderText('请输入手机号');
    fireEvent.change(phoneInput, { target: { value: '13800138000' } });

    // 发送验证码
    const sendButton = screen.getByText('发送验证码');
    fireEvent.click(sendButton);

    // 输入验证码
    const codeInput = screen.getByPlaceholderText('请输入验证码');
    fireEvent.change(codeInput, { target: { value: '123456' } });

    // 点击登录
    const loginButton = screen.getByText('登录');
    fireEvent.click(loginButton);

    // 验证跳转
    await waitFor(() => {
      expect(window.location.pathname).toBe('/sessions');
    });
  });
});
```

---

## ✅ 第一模块完成检查清单

- [x] 登录流程图
- [x] 登录页面 UI 设计
- [x] 手机号验证码登录逻辑
- [x] Auth API 接口定义
- [x] Auth Store 状态管理
- [x] Axios 拦截器配置
- [x] 路由守卫实现
- [x] 自动登录保持
- [x] Token 持久化方案
- [x] 多标签页同步
- [x] 401 错误处理
- [x] 错误类型定义
- [x] 错误提示映射
- [x] 测试方案

---

## 3. 工具函数层

### 3.1 模块概述

工具函数层是应用的基础设施层，提供可复用的通用功能。Web 端需要重新实现小程序的所有工具函数，并适配浏览器环境。

**核心模块**:
- Request 请求封装（已在第 2 章实现）
- Upload 文件上传
- Storage 本地存储
- Avatar 头像处理
- Markdown 渲染
- Config 配置管理
- Format 格式化工具
- Validator 验证工具

### 3.2 与小程序工具层对比

| 小程序工具 | Web 端对应 | 主要差异 | 迁移难度 |
|-----------|-----------|---------|---------|
| utils/auth.js | services/api/auth.ts | 登录方式不同（微信 vs 手机号） | ⭐⭐⭐ |
| utils/request.js | services/api/client.ts | wx.request vs Axios | ⭐⭐ |
| utils/upload.js | services/api/upload.ts | wx.uploadFile vs FormData | ⭐⭐ |
| utils/user.js | services/utils/storage.ts | wx.storage vs localStorage | ⭐ |
| utils/avatar.js | lib/avatar.ts | 逻辑基本一致 | ⭐ |
| utils/markdown.js | lib/markdown.ts | marked 同样可用 | ⭐ |
| utils/config.js | constants/config.ts | 环境检测方式不同 | ⭐⭐ |
| utils/configManager.js | services/api/config.ts | 逻辑基本一致 | ⭐ |

**难度说明**: ⭐ 简单 | ⭐⭐ 中等 | ⭐⭐⭐ 复杂

### 3.3 文件上传工具 (Upload)

#### 3.3.1 小程序实现分析

**小程序代码**（utils/upload.js）:
```javascript
// 使用 wx.uploadFile
wx.uploadFile({
    url: `${API_BASE_URL}/file/wx_upload`,
    filePath: filePath,
    name: 'file',
    formData: { usage_type: 'avatar' },
    header: { Authorization: `Bearer ${token}` }
});
```

**关键特性**:
- 自动添加 Authorization header
- 支持额外的 formData
- 401 错误自动重新登录并重试
- 返回 JSON 自动解析

#### 3.3.2 Web 端实现方案

**文件位置**: `src/services/api/upload.ts`

**完整实现**:
```typescript
import { apiClient } from './client';
import type { UploadFileResponse } from '@/types/api';

export interface UploadOptions {
  usageType?: 'avatar' | 'chat_image' | 'temp_upload';
  onProgress?: (percent: number) => void;
}

/**
 * 上传文件到服务器
 * @param file - File 对象
 * @param options - 上传选项
 */
export async function uploadFile(
  file: File,
  options: UploadOptions = {}
): Promise<UploadFileResponse> {
  const { usageType, onProgress } = options;

  // 创建 FormData
  const formData = new FormData();
  formData.append('file', file);

  if (usageType) {
    formData.append('usage_type', usageType);
  }

  // 发送请求
  const response = await apiClient.post<UploadFileResponse>(
    '/file/wx_upload',
    formData,
    {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const percent = Math.round(
            (progressEvent.loaded * 100) / progressEvent.total
          );
          onProgress(percent);
        }
      }
    }
  );

  return response.data;
}

/**
 * 上传头像（自动压缩）
 * @param file - 图片 File 对象
 * @param onProgress - 进度回调
 */
export async function uploadAvatar(
  file: File,
  onProgress?: (percent: number) => void
): Promise<string> {
  // 1. 压缩图片
  const compressedFile = await compressImage(file, {
    maxWidth: 800,
    maxHeight: 800,
    quality: 0.8
  });

  // 2. 上传
  const response = await uploadFile(compressedFile, {
    usageType: 'avatar',
    onProgress
  });

  // 3. 返回公共 URL
  return response.publicUrl;
}

/**
 * 上传聊天图片
 * @param file - 图片 File 对象
 * @param onProgress - 进度回调
 */
export async function uploadChatImage(
  file: File,
  onProgress?: (percent: number) => void
): Promise<string> {
  const response = await uploadFile(file, {
    usageType: 'chat_image',
    onProgress
  });

  return response.publicUrl;
}

/**
 * 压缩图片
 * @param file - 原始图片
 * @param options - 压缩选项
 */
async function compressImage(
  file: File,
  options: {
    maxWidth: number;
    maxHeight: number;
    quality: number;
  }
): Promise<File> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    reader.onload = (e) => {
      const img = new Image();

      img.onload = () => {
        // 计算压缩后的尺寸
        let { width, height } = img;
        const { maxWidth, maxHeight } = options;

        if (width > maxWidth || height > maxHeight) {
          const ratio = Math.min(maxWidth / width, maxHeight / height);
          width *= ratio;
          height *= ratio;
        }

        // 创建 canvas
        const canvas = document.createElement('canvas');
        canvas.width = width;
        canvas.height = height;

        const ctx = canvas.getContext('2d');
        if (!ctx) {
          reject(new Error('无法创建 Canvas 上下文'));
          return;
        }

        // 绘制图片
        ctx.drawImage(img, 0, 0, width, height);

        // 转换为 Blob
        canvas.toBlob(
          (blob) => {
            if (!blob) {
              reject(new Error('图片压缩失败'));
              return;
            }

            // 创建新的 File 对象
            const compressedFile = new File([blob], file.name, {
              type: 'image/jpeg',
              lastModified: Date.now()
            });

            resolve(compressedFile);
          },
          'image/jpeg',
          options.quality
        );
      };

      img.onerror = () => reject(new Error('图片加载失败'));
      img.src = e.target?.result as string;
    };

    reader.onerror = () => reject(new Error('文件读取失败'));
    reader.readAsDataURL(file);
  });
}
```

**类型定义**:
```typescript
// src/types/api.ts

export interface UploadFileResponse {
  publicUrl: string;      // 公开访问 URL
  privateUrl?: string;    // 私有 URL（如果有）
  fileId?: string;        // 文件 ID
  size?: number;          // 文件大小
}
```

#### 3.3.3 上传组件封装

**文件位置**: `src/components/common/FileUpload.tsx`

```typescript
'use client';

import { useState } from 'react';
import { Upload, Message } from 'tdesign-react';
import { uploadFile } from '@/services/api/upload';
import type { UploadOptions } from '@/services/api/upload';

interface FileUploadProps {
  accept?: string;
  maxSize?: number;          // MB
  usageType?: UploadOptions['usageType'];
  onSuccess?: (url: string) => void;
  onError?: (error: Error) => void;
  children?: React.ReactNode;
}

export default function FileUpload({
  accept = 'image/*',
  maxSize = 10,
  usageType,
  onSuccess,
  onError,
  children
}: FileUploadProps) {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);

  const handleUpload = async (file: File) => {
    // 文件大小校验
    if (file.size > maxSize * 1024 * 1024) {
      Message.error(`文件大小不能超过 ${maxSize}MB`);
      onError?.(new Error('文件过大'));
      return false;
    }

    setUploading(true);
    setProgress(0);

    try {
      const response = await uploadFile(file, {
        usageType,
        onProgress: setProgress
      });

      Message.success('上传成功');
      onSuccess?.(response.publicUrl);
    } catch (error) {
      Message.error('上传失败');
      onError?.(error as Error);
    } finally {
      setUploading(false);
      setProgress(0);
    }

    return false; // 阻止默认上传行为
  };

  return (
    <Upload
      accept={accept}
      beforeUpload={handleUpload}
      showUploadProgress={uploading}
      disabled={uploading}
    >
      {children || (
        <div className="upload-trigger">
          {uploading ? `上传中 ${progress}%` : '点击上传'}
        </div>
      )}
    </Upload>
  );
}
```

### 3.4 本地存储工具 (Storage)

#### 3.4.1 小程序实现分析

**小程序代码**（utils/user.js）:
```javascript
// 使用 wx.getStorageSync 和 wx.setStorageSync
function getUserInfo() {
  return wx.getStorageSync('user_info') || null;
}

function setUserInfo(userInfo) {
  wx.setStorageSync('user_info', userInfo);
}
```

**特点**:
- 同步操作
- 自动序列化/反序列化
- 容量限制 10MB

#### 3.4.2 Web 端实现方案

**文件位置**: `src/services/utils/storage.ts`

**完整实现**:
```typescript
/**
 * 本地存储工具类
 * 提供类型安全的存储操作
 */

// 存储键名枚举
export enum StorageKey {
  ACCESS_TOKEN = 'access_token',
  USER_INFO = 'user_info',
  USER_PROFILE = 'user_profile',
  THEME = 'theme',
  LANGUAGE = 'language'
}

// 存储选项
interface StorageOptions {
  expires?: number;  // 过期时间（毫秒）
}

// 存储数据包装
interface StorageData<T> {
  value: T;
  expires?: number;
  timestamp: number;
}

class Storage {
  /**
   * 设置存储项
   * @param key - 存储键
   * @param value - 存储值
   * @param options - 选项
   */
  set<T>(key: StorageKey | string, value: T, options?: StorageOptions): void {
    try {
      const data: StorageData<T> = {
        value,
        timestamp: Date.now(),
        expires: options?.expires
          ? Date.now() + options.expires
          : undefined
      };

      localStorage.setItem(key, JSON.stringify(data));
    } catch (error) {
      console.error('Storage set error:', error);
      throw new Error('存储失败');
    }
  }

  /**
   * 获取存储项
   * @param key - 存储键
   * @returns 存储值或 null
   */
  get<T>(key: StorageKey | string): T | null {
    try {
      const item = localStorage.getItem(key);
      if (!item) return null;

      const data: StorageData<T> = JSON.parse(item);

      // 检查是否过期
      if (data.expires && Date.now() > data.expires) {
        this.remove(key);
        return null;
      }

      return data.value;
    } catch (error) {
      console.error('Storage get error:', error);
      return null;
    }
  }

  /**
   * 移除存储项
   * @param key - 存储键
   */
  remove(key: StorageKey | string): void {
    try {
      localStorage.removeItem(key);
    } catch (error) {
      console.error('Storage remove error:', error);
    }
  }

  /**
   * 清空所有存储
   */
  clear(): void {
    try {
      localStorage.clear();
    } catch (error) {
      console.error('Storage clear error:', error);
    }
  }

  /**
   * 检查键是否存在
   * @param key - 存储键
   */
  has(key: StorageKey | string): boolean {
    return localStorage.getItem(key) !== null;
  }

  /**
   * 获取所有键
   */
  keys(): string[] {
    return Object.keys(localStorage);
  }

  /**
   * 获取存储大小（估算，字节）
   */
  size(): number {
    let total = 0;
    for (const key in localStorage) {
      if (localStorage.hasOwnProperty(key)) {
        total += localStorage[key].length + key.length;
      }
    }
    return total;
  }
}

// 导出单例
export const storage = new Storage();

// 便捷方法（类型安全）
export const storageHelpers = {
  // Token
  getToken(): string | null {
    return storage.get<string>(StorageKey.ACCESS_TOKEN);
  },

  setToken(token: string): void {
    storage.set(StorageKey.ACCESS_TOKEN, token);
  },

  removeToken(): void {
    storage.remove(StorageKey.ACCESS_TOKEN);
  },

  // 用户信息
  getUserInfo(): User | null {
    return storage.get<User>(StorageKey.USER_INFO);
  },

  setUserInfo(user: User): void {
    storage.set(StorageKey.USER_INFO, user);
  },

  // 用户资料
  getUserProfile(): Profile | null {
    return storage.get<Profile>(StorageKey.USER_PROFILE);
  },

  setUserProfile(profile: Profile): void {
    storage.set(StorageKey.USER_PROFILE, profile);
  },

  // 清除用户数据
  clearUserData(): void {
    storage.remove(StorageKey.ACCESS_TOKEN);
    storage.remove(StorageKey.USER_INFO);
    storage.remove(StorageKey.USER_PROFILE);
  }
};
```

**SessionStorage 版本**（用于临时数据）:
```typescript
// src/services/utils/sessionStorage.ts

class SessionStorage {
  set<T>(key: string, value: T): void {
    sessionStorage.setItem(key, JSON.stringify(value));
  }

  get<T>(key: string): T | null {
    const item = sessionStorage.getItem(key);
    return item ? JSON.parse(item) : null;
  }

  remove(key: string): void {
    sessionStorage.removeItem(key);
  }

  clear(): void {
    sessionStorage.clear();
  }
}

export const sessionStore = new SessionStorage();
```

### 3.5 头像处理工具 (Avatar)

#### 3.5.1 小程序实现分析

**小程序代码**（utils/avatar.js）:
```javascript
// 根据性别返回默认头像
function getDefaultAvatar(gender, userGender, isFriend = false) {
  if (gender === 'male') {
    return '/assets/images/male.png';
  } else if (gender === 'female') {
    return '/assets/images/female.png';
  }

  // 好友性别未知时，使用 assistant.png
  if (isFriend) {
    return '/assets/images/assistant.png';
  }

  return '/assets/images/male.png';
}
```

**逻辑要点**:
1. 优先使用用户上传的头像
2. 无头像时根据性别显示默认头像
3. 好友性别未知时显示通用头像
4. 头像加载失败时自动降级

#### 3.5.2 Web 端实现方案

**文件位置**: `src/lib/avatar.ts`

**完整实现**:
```typescript
/**
 * 头像处理工具
 */

export type Gender = 'male' | 'female' | '';

// 默认头像路径
const DEFAULT_AVATARS = {
  male: '/assets/images/male.png',
  female: '/assets/images/female.png',
  assistant: '/assets/images/assistant.png'
} as const;

/**
 * 获取默认头像路径
 * @param gender - 用户性别
 * @param userGender - 当前用户性别（用于推断好友性别）
 * @param isFriend - 是否为好友头像
 */
export function getDefaultAvatar(
  gender: Gender,
  userGender?: Gender,
  isFriend: boolean = false
): string {
  // 有明确性别
  if (gender === 'male') {
    return DEFAULT_AVATARS.male;
  }
  if (gender === 'female') {
    return DEFAULT_AVATARS.female;
  }

  // 性别未知，根据当前用户推断
  if (userGender && !isFriend) {
    // 自己性别未知，使用相反性别的头像
    return userGender === 'male'
      ? DEFAULT_AVATARS.female
      : DEFAULT_AVATARS.male;
  }

  // 好友性别未知，使用通用头像
  if (isFriend) {
    return DEFAULT_AVATARS.assistant;
  }

  // 默认使用男性头像
  return DEFAULT_AVATARS.male;
}

/**
 * 获取用户头像（带降级）
 * @param avatar - 用户头像 URL
 * @param gender - 用户性别
 * @param userGender - 当前用户性别
 * @param isFriend - 是否为好友
 */
export function getUserAvatar(
  avatar?: string | null,
  gender?: Gender,
  userGender?: Gender,
  isFriend: boolean = false
): string {
  // 有头像直接返回
  if (avatar && avatar.trim() !== '') {
    return avatar;
  }

  // 无头像返回默认头像
  return getDefaultAvatar(gender || '', userGender, isFriend);
}

/**
 * 验证头像 URL 是否有效
 * @param url - 头像 URL
 */
export function validateAvatarUrl(url: string): Promise<boolean> {
  return new Promise((resolve) => {
    const img = new Image();
    img.onload = () => resolve(true);
    img.onerror = () => resolve(false);
    img.src = url;
  });
}

/**
 * 获取头像颜色（用于占位符）
 * @param name - 用户名
 */
export function getAvatarColor(name: string): string {
  const colors = [
    '#1890ff', '#52c41a', '#faad14', '#f5222d',
    '#722ed1', '#eb2f96', '#13c2c2', '#fa8c16'
  ];

  // 根据名字计算哈希值
  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }

  return colors[Math.abs(hash) % colors.length];
}

/**
 * 获取名字首字母（用于无头像时显示）
 * @param name - 用户名
 */
export function getAvatarInitial(name: string): string {
  if (!name) return '?';

  // 中文取第一个字
  if (/[\u4e00-\u9fa5]/.test(name)) {
    return name.charAt(0);
  }

  // 英文取首字母
  return name.charAt(0).toUpperCase();
}
```

**Avatar 组件封装**:
```typescript
// src/components/common/Avatar.tsx

'use client';

import { useState } from 'react';
import { Avatar as TAvatar } from 'tdesign-react';
import {
  getUserAvatar,
  getAvatarColor,
  getAvatarInitial
} from '@/lib/avatar';
import type { Gender } from '@/lib/avatar';

interface AvatarProps {
  src?: string | null;
  name?: string;
  gender?: Gender;
  userGender?: Gender;
  isFriend?: boolean;
  size?: 'small' | 'medium' | 'large' | number;
  className?: string;
}

export default function Avatar({
  src,
  name = '',
  gender,
  userGender,
  isFriend = false,
  size = 'medium',
  className
}: AvatarProps) {
  const [error, setError] = useState(false);

  // 获取头像 URL
  const avatarUrl = error
    ? getUserAvatar(null, gender, userGender, isFriend)
    : getUserAvatar(src, gender, userGender, isFriend);

  // 头像加载失败处理
  const handleError = () => {
    setError(true);
  };

  // 如果没有图片，显示名字首字母
  if (!avatarUrl || error) {
    const initial = getAvatarInitial(name);
    const bgColor = getAvatarColor(name);

    return (
      <TAvatar
        size={size}
        className={className}
        style={{ backgroundColor: bgColor }}
      >
        {initial}
      </TAvatar>
    );
  }

  return (
    <TAvatar
      image={avatarUrl}
      size={size}
      className={className}
      onError={handleError}
    />
  );
}
```

### 3.6 Markdown 渲染工具

#### 3.6.1 小程序实现分析

**小程序代码**（utils/markdown.js）:
- 使用 `marked` 库解析 Markdown
- 自定义渲染器添加 inline styles
- 支持 GFM（GitHub Flavored Markdown）
- 输出带样式的 HTML 供 rich-text 组件使用

**样式特点**:
- 标题分级（h1-h6）不同字号
- 代码块背景色
- 引用块左边框
- 列表缩进

#### 3.6.2 Web 端实现方案

**文件位置**: `src/lib/markdown.ts`

**使用 react-markdown**:
```typescript
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';

/**
 * Markdown 组件配置
 */
export const markdownComponents = {
  // 代码块
  code({ node, inline, className, children, ...props }) {
    const match = /language-(\w+)/.exec(className || '');

    return !inline && match ? (
      <SyntaxHighlighter
        style={vscDarkPlus}
        language={match[1]}
        PreTag="div"
        {...props}
      >
        {String(children).replace(/\n$/, '')}
      </SyntaxHighlighter>
    ) : (
      <code className={className} {...props}>
        {children}
      </code>
    );
  },

  // 链接（新窗口打开）
  a({ node, children, href, ...props }) {
    return (
      <a
        href={href}
        target="_blank"
        rel="noopener noreferrer"
        {...props}
      >
        {children}
      </a>
    );
  },

  // 图片（懒加载）
  img({ node, src, alt, ...props }) {
    return (
      <img
        src={src}
        alt={alt}
        loading="lazy"
        style={{ maxWidth: '100%', height: 'auto' }}
        {...props}
      />
    );
  }
};

/**
 * Markdown 渲染配置
 */
export const markdownOptions = {
  remarkPlugins: [remarkGfm],  // GitHub 风格 Markdown
  components: markdownComponents
};
```

**Markdown 渲染组件**:
```typescript
// src/components/common/MarkdownRenderer.tsx

'use client';

import ReactMarkdown from 'react-markdown';
import { markdownOptions } from '@/lib/markdown';
import styles from './MarkdownRenderer.module.css';

interface MarkdownRendererProps {
  content: string;
  className?: string;
}

export default function MarkdownRenderer({
  content,
  className
}: MarkdownRendererProps) {
  return (
    <div className={`${styles.markdown} ${className || ''}`}>
      <ReactMarkdown {...markdownOptions}>
        {content}
      </ReactMarkdown>
    </div>
  );
}
```

**样式文件**:
```css
/* src/components/common/MarkdownRenderer.module.css */

.markdown {
  font-size: 16px;
  line-height: 1.75;
  color: #404040;
}

/* 标题 */
.markdown h1 {
  font-size: 22px;
  font-weight: 600;
  margin: 16px 0 8px 0;
}

.markdown h2 {
  font-size: 20px;
  font-weight: 600;
  margin: 16px 0 8px 0;
}

.markdown h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 16px 0 8px 0;
}

/* 段落 */
.markdown p {
  margin: 8px 0;
}

/* 列表 */
.markdown ul,
.markdown ol {
  margin: 8px 0;
  padding-left: 24px;
}

.markdown li + li {
  margin-top: 4px;
}

/* 引用块 */
.markdown blockquote {
  border-left: 2px solid #a3a3a3;
  padding-left: 16px;
  margin: 12px 0;
  color: #737373;
}

/* 代码 */
.markdown code {
  background-color: #f5f5f5;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
}

.markdown pre {
  background-color: #1e1e1e;
  padding: 12px;
  border-radius: 4px;
  overflow-x: auto;
  margin: 12px 0;
}

.markdown pre code {
  background: none;
  padding: 0;
  color: #d4d4d4;
}

/* 水平线 */
.markdown hr {
  height: 1px;
  margin: 12px 0;
  background-color: #e5e5e5;
  border: none;
}

/* 链接 */
.markdown a {
  color: #3b82f6;
  text-decoration: none;
}

.markdown a:hover {
  text-decoration: underline;
}

/* 表格 */
.markdown table {
  border-collapse: collapse;
  width: 100%;
  margin: 12px 0;
}

.markdown th,
.markdown td {
  border: 1px solid #e5e5e5;
  padding: 8px 12px;
  text-align: left;
}

.markdown th {
  background-color: #f5f5f5;
  font-weight: 600;
}
```

### 3.7 配置管理工具 (Config)

#### 3.7.1 小程序实现分析

**小程序代码**（utils/config.js）:
```javascript
// 环境识别
const accountInfo = wx.getAccountInfoSync();
const ENV_VERSION = accountInfo.miniProgram.envVersion;

// develop -> 开发环境
// trial/release -> 生产环境

const API_BASE_URL = ENV_VERSION === 'develop'
  ? 'https://local.chathandy.com'
  : 'https://mp.chathandy.com';
```

**配置管理器**（utils/configManager.js）:
- 从服务端拉取动态配置
- 本地缓存 30 分钟
- 批量获取配置
- 自动刷新过期配置

#### 3.7.2 Web 端实现方案

**环境配置文件**:

```bash
# .env.local (开发环境)
NEXT_PUBLIC_API_BASE_URL=https://local.chathandy.com
NEXT_PUBLIC_APP_VERSION=1.0.0
NEXT_PUBLIC_APP_ENV=development
```

```bash
# .env.production (生产环境)
NEXT_PUBLIC_API_BASE_URL=https://mp.chathandy.com
NEXT_PUBLIC_APP_VERSION=1.0.0
NEXT_PUBLIC_APP_ENV=production
```

**配置常量**:
```typescript
// src/constants/config.ts

export const config = {
  // API 配置
  apiBaseUrl: process.env.NEXT_PUBLIC_API_BASE_URL || '',
  appVersion: process.env.NEXT_PUBLIC_APP_VERSION || '1.0.0',
  appEnv: process.env.NEXT_PUBLIC_APP_ENV || 'production',

  // 应用配置
  appName: '恋爱翻译官',
  appDescription: 'AI 驱动的社交辅助应用',

  // 分页配置
  pageSize: 10,

  // 文件上传配置
  maxFileSize: 10 * 1024 * 1024, // 10MB
  allowedImageTypes: ['image/jpeg', 'image/png', 'image/gif', 'image/webp'],

  // 缓存配置
  cacheExpiration: 30 * 60 * 1000, // 30分钟
} as const;

// 环境判断
export const isDev = config.appEnv === 'development';
export const isProd = config.appEnv === 'production';
```

**动态配置管理**:
```typescript
// src/services/api/config.ts

interface ConfigItem {
  key: string;
  value: string;
}

interface ConfigResponse {
  configs: ConfigItem[];
}

/**
 * 配置管理器
 */
class ConfigManager {
  private cache: Map<string, { value: string; timestamp: number }>;
  private cacheExpiration: number;

  constructor() {
    this.cache = new Map();
    this.cacheExpiration = config.cacheExpiration;
  }

  /**
   * 获取单个配置
   * @param key - 配置键
   * @param defaultValue - 默认值
   */
  async get(key: string, defaultValue: string = ''): Promise<string> {
    const result = await this.getMulti({ [key]: defaultValue });
    return result[key];
  }

  /**
   * 批量获取配置
   * @param keysWithDefaults - 配置键和默认值的对象
   */
  async getMulti(
    keysWithDefaults: Record<string, string>
  ): Promise<Record<string, string>> {
    const keys = Object.keys(keysWithDefaults);
    const result: Record<string, string> = {};

    // 检查是否需要刷新
    const needRefresh = keys.some((key) => this.isExpired(key));

    if (needRefresh) {
      await this.refresh(keys);
    }

    // 获取配置值
    for (const key of keys) {
      const cached = this.cache.get(key);
      result[key] = cached?.value || keysWithDefaults[key];
    }

    return result;
  }

  /**
   * 刷新配置
   * @param keys - 需要刷新的配置键
   */
  private async refresh(keys: string[]): Promise<void> {
    try {
      const response = await apiClient.post<ConfigResponse>(
        '/config.ConfigService/GetConfig',
        {
          keys,
          app: 'chathandy',
          platform: 'web',
          env: config.appEnv,
          version: config.appVersion
        }
      );

      const timestamp = Date.now();

      // 更新缓存
      for (const item of response.data.configs) {
        this.cache.set(item.key, {
          value: item.value,
          timestamp
        });
      }

      // 对于没有返回的 key，设置为空并记录时间戳
      for (const key of keys) {
        if (!response.data.configs.find((c) => c.key === key)) {
          this.cache.set(key, {
            value: '',
            timestamp
          });
        }
      }
    } catch (error) {
      console.error('刷新配置失败:', error);
    }
  }

  /**
   * 检查配置是否过期
   * @param key - 配置键
   */
  private isExpired(key: string): boolean {
    const cached = this.cache.get(key);
    if (!cached) return true;

    return Date.now() - cached.timestamp > this.cacheExpiration;
  }

  /**
   * 清除缓存
   */
  clearCache(): void {
    this.cache.clear();
  }
}

// 导出单例
export const configManager = new ConfigManager();
```

### 3.8 格式化工具 (Format)

**文件位置**: `src/services/utils/format.ts`

```typescript
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';

// 配置 dayjs
dayjs.extend(relativeTime);
dayjs.locale('zh-cn');

/**
 * 格式化时间
 * @param date - 日期
 * @param format - 格式
 */
export function formatDate(
  date: string | Date | number,
  format: string = 'YYYY-MM-DD HH:mm:ss'
): string {
  return dayjs(date).format(format);
}

/**
 * 格式化相对时间（刚刚、1分钟前等）
 * @param date - 日期
 */
export function formatRelativeTime(date: string | Date | number): string {
  const now = dayjs();
  const target = dayjs(date);

  const diffMinutes = now.diff(target, 'minute');
  const diffHours = now.diff(target, 'hour');
  const diffDays = now.diff(target, 'day');

  if (diffMinutes < 1) return '刚刚';
  if (diffMinutes < 60) return `${diffMinutes}分钟前`;
  if (diffHours < 24) return `${diffHours}小时前`;
  if (diffDays < 7) return `${diffDays}天前`;

  return formatDate(date, 'MM-DD HH:mm');
}

/**
 * 格式化文件大小
 * @param bytes - 字节数
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B';

  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`;
}

/**
 * 格式化手机号（隐藏中间4位）
 * @param phone - 手机号
 */
export function formatPhone(phone: string): string {
  if (!phone || phone.length !== 11) return phone;
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2');
}

/**
 * 格式化数字（千分位）
 * @param num - 数字
 */
export function formatNumber(num: number): string {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
}
```

### 3.9 验证工具 (Validator)

**文件位置**: `src/services/utils/validator.ts`

```typescript
/**
 * 验证手机号
 * @param phone - 手机号
 */
export function validatePhone(phone: string): boolean {
  const phoneRegex = /^1[3-9]\d{9}$/;
  return phoneRegex.test(phone);
}

/**
 * 验证邮箱
 * @param email - 邮箱
 */
export function validateEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

/**
 * 验证验证码（6位数字）
 * @param code - 验证码
 */
export function validateVerifyCode(code: string): boolean {
  const codeRegex = /^\d{6}$/;
  return codeRegex.test(code);
}

/**
 * 验证密码强度
 * @param password - 密码
 * @returns 强度等级 0-4
 */
export function validatePasswordStrength(password: string): number {
  if (password.length < 6) return 0;

  let strength = 0;

  // 长度
  if (password.length >= 8) strength++;
  if (password.length >= 12) strength++;

  // 包含小写字母
  if (/[a-z]/.test(password)) strength++;

  // 包含大写字母
  if (/[A-Z]/.test(password)) strength++;

  // 包含数字
  if (/\d/.test(password)) strength++;

  // 包含特殊字符
  if (/[!@#$%^&*(),.?":{}|<>]/.test(password)) strength++;

  return Math.min(strength, 4);
}

/**
 * 验证 URL
 * @param url - URL
 */
export function validateUrl(url: string): boolean {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
}

/**
 * 验证身份证号
 * @param idCard - 身份证号
 */
export function validateIdCard(idCard: string): boolean {
  const idCardRegex = /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/;
  return idCardRegex.test(idCard);
}
```

---

## ✅ 第三模块完成检查清单

- [x] 文件上传工具（upload）
  - [x] 基础上传功能
  - [x] 图片压缩
  - [x] 进度回调
  - [x] 头像专用上传
  - [x] Upload 组件封装
- [x] 本地存储工具（storage）
  - [x] 类型安全的存储操作
  - [x] 过期时间支持
  - [x] 便捷方法封装
  - [x] SessionStorage 支持
- [x] 头像处理工具（avatar）
  - [x] 默认头像逻辑
  - [x] 性别判断
  - [x] 头像组件封装
  - [x] 占位符支持
- [x] Markdown 渲染
  - [x] react-markdown 配置
  - [x] 代码高亮
  - [x] 样式定制
  - [x] 组件封装
- [x] 配置管理（config）
  - [x] 环境配置
  - [x] 静态配置
  - [x] 动态配置管理
  - [x] 缓存机制
- [x] 格式化工具（format）
  - [x] 时间格式化
  - [x] 相对时间
  - [x] 文件大小
  - [x] 手机号脱敏
- [x] 验证工具（validator）
  - [x] 手机号验证
  - [x] 邮箱验证
  - [x] 验证码验证
  - [x] URL 验证

---

**当前进度**: ✅ 第 3 章节（工具函数层）已完成，等待用户确认后继续第 4 章节（会话列表页面）。

