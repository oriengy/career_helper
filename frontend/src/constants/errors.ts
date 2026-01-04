/**
 * 错误信息常量
 */

import { ErrorCode } from '@/types/common';

export const ERROR_MESSAGES: Record<ErrorCode, string> = {
  [ErrorCode.INVALID_PHONE]: '手机号格式不正确',
  [ErrorCode.INVALID_CODE]: '验证码错误',
  [ErrorCode.CODE_EXPIRED]: '验证码已过期，请重新获取',
  [ErrorCode.TOO_MANY_REQUESTS]: '请求过于频繁，请稍后再试',
  [ErrorCode.TOKEN_EXPIRED]: '登录已过期，请重新登录',
  [ErrorCode.TOKEN_INVALID]: '登录状态异常，请重新登录',
  [ErrorCode.NETWORK_ERROR]: '网络连接失败，请检查网络',
  [ErrorCode.UNKNOWN_ERROR]: '操作失败，请重试',
};
