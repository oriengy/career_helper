/**
 * 数据模型类型定义
 */

export type Gender = 'male' | 'female' | '';

export type MessageRole = 'SELF' | 'FRIEND' | 'AI' | 'USER';

export type MessageType = 'HISTORY' | 'TRANSLATE' | 'CONSULT';

// 用户信息
export interface User {
  id: string;
  name?: string;
  imName?: string;
  externalId?: string;
  phone?: string;
  avatar?: string;
  profileId?: string;
  createdAt?: string;
  updatedAt?: string;
}

// 用户资料
export interface Profile {
  id: string;
  userId: string;
  name?: string;
  imName?: string;
  avatar?: string;
  age?: number;
  gender?: Gender;
  birthday?: string;
  birthLocation?: string;
  currentLocation?: string;
  createdAt?: string;
  updatedAt?: string;
}

// 聊天会话
export interface ChatSession {
  id: string;
  name: string;
  userId: string;
  profileId: string;
  avatar?: string;
  gender?: Gender;
  createdAt?: string;
  updatedAt?: string;
  displayAvatar?: string; // 显示用的头像（包含默认头像逻辑）
}

// 消息
export interface Message {
  id: string;
  userId: string;
  sessionId: string;
  parentId?: string;
  profileId?: string;
  role: MessageRole;
  msgType: MessageType;
  content: string;
  tags?: string[];
  msgAt: string;
  createdAt?: string;
  updatedAt?: string;
  // 扩展字段（前端使用）
  translateContent?: string;  // 翻译内容
  translateDirection?: 'male' | 'female';  // 翻译方向
}
