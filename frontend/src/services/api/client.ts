/**
 * API 客户端配置
 */

import axios from 'axios';
import { config } from '@/constants/config';
import { storageHelpers } from '@/services/utils/storage';

// 创建 axios 实例
export const apiClient = axios.create({
  baseURL: config.apiBaseUrl,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
    'Connect-Protocol-Version': '1',
  },
});

// 请求拦截器：自动添加 token
apiClient.interceptors.request.use(
  (config) => {
    // 获取 token
    const token = storageHelpers.getToken();

    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    // 添加自定义 headers
    config.headers['X-App-Platform'] = 'web';
    config.headers['X-App-Env'] = process.env.NODE_ENV || 'production';
    config.headers['X-App-Version'] =
      process.env.NEXT_PUBLIC_APP_VERSION || '1.0.0';

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
      storageHelpers.clearUserData();

      // 跳转到登录页
      if (typeof window !== 'undefined') {
        window.location.href = '/login';
      }

      return Promise.reject(error);
    }

    return Promise.reject(error);
  }
);
