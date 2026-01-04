/**
 * 用户状态管理
 */

import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User, Profile } from '@/types/models';
import { storageHelpers } from '@/services/utils/storage';

interface UserState {
  // 状态
  user: User | null;
  profile: Profile | null;

  // 操作
  setUser: (user: User) => void;
  setProfile: (profile: Profile) => void;
  updateProfile: (partialProfile: Partial<Profile>) => void;
  clearUser: () => void;
}

export const useUserStore = create<UserState>()(
  persist(
    (set, get) => ({
      // 初始状态
      user: null,
      profile: null,

      // 设置用户信息
      setUser: (user) => {
        set({ user });
        storageHelpers.setUserInfo(user);
      },

      // 设置用户资料
      setProfile: (profile) => {
        set({ profile });
        storageHelpers.setUserProfile(profile);
      },

      // 更新部分用户资料
      updateProfile: (partialProfile) => {
        const { profile } = get();
        if (profile) {
          const updatedProfile = { ...profile, ...partialProfile };
          get().setProfile(updatedProfile);
        }
      },

      // 清除用户数据
      clearUser: () => {
        set({ user: null, profile: null });
        storageHelpers.clearUserData();
      },
    }),
    {
      name: 'user-storage',
      partialize: (state) => ({ user: state.user, profile: state.profile }),
    }
  )
);
