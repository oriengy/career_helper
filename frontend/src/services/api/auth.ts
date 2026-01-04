/**
 * 认证相关 API
 */

import { apiClient } from './client';
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
    await apiClient.post('/user.UserService/SendVerifyCode', params);
  },

  /**
   * 手机号登录
   */
  async phoneLogin(params: PhoneLoginParams): Promise<PhoneLoginResponse> {
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
