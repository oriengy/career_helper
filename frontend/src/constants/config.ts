/**
 * 应用配置常量
 */

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

  // Token 有效期
  tokenExpiration: 365 * 24 * 60 * 60 * 1000, // 365天
} as const;

// 环境判断
export const isDev = config.appEnv === 'development';
export const isProd = config.appEnv === 'production';
