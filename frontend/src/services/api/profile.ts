/**
 * 资料相关 API
 */

import { apiClient } from './client';
import type {
  CreateProfileParams,
  UpdateProfileParams,
  GetProfileParams,
  ProfileResponse,
} from '@/types/api';
import type { Profile } from '@/types/models';

export const profileApi = {
  /**
   * 获取资料
   */
  async getProfile(params: GetProfileParams): Promise<Profile> {
    const response = await apiClient.post<ProfileResponse>(
      '/profile.ProfileService/GetProfile',
      params
    );
    return response.data.profile;
  },

  /**
   * 创建资料
   */
  async createProfile(params: CreateProfileParams): Promise<Profile> {
    const response = await apiClient.post<ProfileResponse>(
      '/profile.ProfileService/CreateProfile',
      params
    );
    return response.data.profile;
  },

  /**
   * 更新资料
   */
  async updateProfile(params: UpdateProfileParams): Promise<Profile> {
    const response = await apiClient.post<ProfileResponse>(
      '/profile.ProfileService/UpdateProfile',
      params
    );
    return response.data.profile;
  },

  /**
   * 删除资料
   */
  async deleteProfile(id: string): Promise<void> {
    await apiClient.post('/profile.ProfileService/DeleteProfile', { id });
  },
};
