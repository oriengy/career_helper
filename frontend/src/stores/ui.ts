/**
 * UI 状态管理
 */

import { create } from 'zustand';

interface UIState {
  // 状态
  isLoading: boolean;
  isMobileMenuOpen: boolean;

  // 操作
  setLoading: (loading: boolean) => void;
  toggleMobileMenu: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  // 初始状态
  isLoading: false,
  isMobileMenuOpen: false,

  // 操作
  setLoading: (loading) => set({ isLoading: loading }),
  toggleMobileMenu: () =>
    set((state) => ({ isMobileMenuOpen: !state.isMobileMenuOpen })),
}));
