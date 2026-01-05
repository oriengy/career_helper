/**
 * 认证状态管理
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { authApi } from '@/services/api/auth';
import { storageHelpers } from '@/services/utils/storage';
import { useUserStore } from './user';

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
        storageHelpers.setToken(token);
      },

      // 清除 token
      clearToken: () => {
        set({ token: null, isLoggedIn: false });
        storageHelpers.removeToken();
      },

      // 检查认证状态
      checkAuth: async () => {
        const { token } = get();
        if (!token) {
          return false;
        }

        try {
          // 调用获取用户信息接口验证 token
          const data = await authApi.getUserProfile();

          // 更新用户信息
          const userStore = useUserStore.getState();
          userStore.setUser(data.user);
          userStore.setProfile(data.profile);

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
          const response = await authApi.phoneLogin({
            phone,
            verificationCode: code,
          });
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
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ token: state.token }),
    }
  )
);
