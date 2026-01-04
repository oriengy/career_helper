/**
 * 路由常量
 */

export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  SESSIONS: '/sessions',
  SESSION_NEW: '/sessions/new',
  SESSION_DETAIL: (id: string) => `/sessions/${id}`,
  PROFILE: '/profile',
  PROFILE_EDIT: '/profile/edit',
  PROFILE_EDIT_FRIEND: (id: string) => `/profile/${id}/edit`,
  GENDER: '/gender',
} as const;

// 不需要认证的路由白名单
export const PUBLIC_ROUTES = [ROUTES.LOGIN];
