/**
 * 本地存储工具类
 */

import { StorageKey } from '@/types/common';
import type { User, Profile } from '@/types/models';

// 存储选项
interface StorageOptions {
  expires?: number; // 过期时间（毫秒）
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
   */
  set<T>(key: StorageKey | string, value: T, options?: StorageOptions): void {
    try {
      const data: StorageData<T> = {
        value,
        timestamp: Date.now(),
        expires: options?.expires ? Date.now() + options.expires : undefined,
      };

      localStorage.setItem(key, JSON.stringify(data));
    } catch (error) {
      console.error('Storage set error:', error);
      throw new Error('存储失败');
    }
  }

  /**
   * 获取存储项
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
  },
};
