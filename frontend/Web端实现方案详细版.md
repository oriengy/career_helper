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

**当前进度**: ✅ 第 2 章节（登录与鉴权模块）已完成，等待用户确认后继续第 3 章节（工具函数层）。

