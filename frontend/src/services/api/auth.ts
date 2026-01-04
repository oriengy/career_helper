/**
 * 认证相关 API
 */

import { apiClient } from './client';
import { config } from '@/constants/config';
import { storageHelpers } from '@/services/utils/storage';
import type { User, Profile } from '@/types/models';
import type {
  SendVerifyCodeParams,
  PhoneLoginParams,
  PhoneLoginResponse,
  UserProfileResponse,
} from '@/types/api';

export const authApi = {
  /**
   * 发送验证码
   */
  async sendVerifyCode(params: SendVerifyCodeParams): Promise<void> {
    if (config.mockLogin) {
      return Promise.resolve();
    }
    await apiClient.post('/user.UserService/SendVerifyCode', params);
  },

  /**
   * 手机号登录
   */
  async phoneLogin(params: PhoneLoginParams): Promise<PhoneLoginResponse> {
    if (config.mockLogin) {
      const mockUser: User = {
        id: '10001',
        name: 'Test User',
        phone: params.phone,
        profileId: '20001',
      };
      const mockProfile: Profile = {
        id: '20001',
        userId: '10001',
        name: 'Test User',
        gender: '',
      };

      storageHelpers.setUserInfo(mockUser);
      storageHelpers.setUserProfile(mockProfile);

      return {
        token: 'test-token',
        user: mockUser,
        profile: mockProfile,
      };
    }
    const response = await apiClient.post<PhoneLoginResponse>(
      '/user.UserService/PhoneLogin',
      params
    );
    return response.data;
  },

  /**
   * 获取用户信息
   */
  async getUserProfile(): Promise<UserProfileResponse> {
    if (config.mockLogin) {
      const user = storageHelpers.getUserInfo() || {
        id: '10001',
        name: 'Test User',
        profileId: '20001',
      };
      const profile = storageHelpers.getUserProfile() || {
        id: '20001',
        userId: '10001',
        name: 'Test User',
        gender: '',
      };
      return { user, profile };
    }
    const response = await apiClient.post<UserProfileResponse>(
      '/user.UserService/GetUserProfile',
      {}
    );
    return response.data;
  },

  /**
   * 登出（清理客户端状态）
   */
  async logout(): Promise<void> {
    // 后端可能不需要 logout 接口，主要是清理本地状态
    return Promise.resolve();
  },
};
