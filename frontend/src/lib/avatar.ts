/**
 * 头像处理工具
 */

import type { Gender } from '@/types/models';

// 默认头像路径
const DEFAULT_AVATARS = {
  male: '/assets/images/male.png',
  female: '/assets/images/female.png',
  assistant: '/assets/images/assistant.png',
} as const;

/**
 * 获取默认头像路径
 */
export function getDefaultAvatar(
  gender: Gender,
  userGender?: Gender,
  isFriend: boolean = false
): string {
  // 有明确性别
  if (gender === 'male') {
    return DEFAULT_AVATARS.male;
  }
  if (gender === 'female') {
    return DEFAULT_AVATARS.female;
  }

  // 性别未知，根据当前用户推断
  if (userGender && !isFriend) {
    // 自己性别未知，使用相反性别的头像
    return userGender === 'male'
      ? DEFAULT_AVATARS.female
      : DEFAULT_AVATARS.male;
  }

  // 好友性别未知，使用通用头像
  if (isFriend) {
    return DEFAULT_AVATARS.assistant;
  }

  // 默认使用男性头像
  return DEFAULT_AVATARS.male;
}

/**
 * 获取用户头像（带降级）
 */
export function getUserAvatar(
  avatar?: string | null,
  gender?: Gender,
  userGender?: Gender,
  isFriend: boolean = false
): string {
  // 有头像直接返回
  if (avatar && avatar.trim() !== '') {
    return avatar;
  }

  // 无头像返回默认头像
  return getDefaultAvatar(gender || '', userGender, isFriend);
}

/**
 * 验证头像 URL 是否有效
 */
export function validateAvatarUrl(url: string): Promise<boolean> {
  return new Promise((resolve) => {
    const img = new Image();
    img.onload = () => resolve(true);
    img.onerror = () => resolve(false);
    img.src = url;
  });
}

/**
 * 获取头像颜色（用于占位符）
 */
export function getAvatarColor(name: string): string {
  const colors = [
    '#1890ff',
    '#52c41a',
    '#faad14',
    '#f5222d',
    '#722ed1',
    '#eb2f96',
    '#13c2c2',
    '#fa8c16',
  ];

  // 根据名字计算哈希值
  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }

  return colors[Math.abs(hash) % colors.length];
}

/**
 * 获取名字首字母（用于无头像时显示）
 */
export function getAvatarInitial(name: string): string {
  if (!name) return '?';

  // 中文取第一个字
  if (/[\u4e00-\u9fa5]/.test(name)) {
    return name.charAt(0);
  }

  // 英文取首字母
  return name.charAt(0).toUpperCase();
}
