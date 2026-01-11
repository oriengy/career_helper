/**
 * 通用类型定义
 */

// 存储键名枚举
export enum StorageKey {
  ACCESS_TOKEN = 'access_token',
  USER_INFO = 'user_info',
  USER_PROFILE = 'user_profile',
  THEME = 'theme',
  LANGUAGE = 'language',
}

// 错误码
export enum ErrorCode {
  INVALID_PHONE = 'INVALID_PHONE',
  INVALID_CODE = 'INVALID_CODE',
  CODE_EXPIRED = 'CODE_EXPIRED',
  TOO_MANY_REQUESTS = 'TOO_MANY_REQUESTS',
  TOKEN_EXPIRED = 'TOKEN_EXPIRED',
  TOKEN_INVALID = 'TOKEN_INVALID',
  NETWORK_ERROR = 'NETWORK_ERROR',
  UNKNOWN_ERROR = 'UNKNOWN_ERROR',
}

// 路由路径
export enum RoutePath {
  HOME = '/',
  LOGIN = '/login',
  SESSIONS = '/sessions',
  SESSION_DETAIL = '/sessions/[id]',
  SESSION_NEW = '/sessions/new',
  PROFILE = '/profile',
  PROFILE_EDIT = '/profile/edit',
  PROFILE_EDIT_FRIEND = '/profile/[id]/edit',
}
