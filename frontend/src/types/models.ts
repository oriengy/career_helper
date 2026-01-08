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
  custom?: ProfileProperty[];
  birthday?: string;
  birthLocation?: string;
  currentLocation?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface ProfileProperty {
  name: string;
  value: string;
}

// 聊天会话
export interface ChatSession {
  id: string;
  sessionId: string; // 会话ID（与 id 相同，兼容字段）
  name: string;
  friendName: string; // 好友名称（兼容字段）
  friendGender?: Gender; // 好友性别（兼容字段）
  friendAvatar?: string; // 好友头像（兼容字段）
  userId: string;
  profileId: string;
  avatar?: string;
  gender?: Gender;
  createdAt?: string;
  updatedAt?: string;
  lastMessage?: string; // 最后一条消息
  displayAvatar?: string; // 显示用的头像（包含默认头像逻辑）
}

// 消息
export interface Message {
  id: string;
  messageId: string; // 消息ID（与 id 相同，兼容字段）
  userId: string;
  sessionId: string;
  parentId?: string;
  profileId?: string;
  role: MessageRole;
  msgType: MessageType;
  content?: string;
  imageUrl?: string; // 图片URL
  tags?: string[];
  msgAt: string;
  createdAt?: string;
  updatedAt?: string;
  // 扩展字段（前端使用）
  translateContent?: string;  // 翻译内容
  translateDirection?: 'male' | 'female';  // 翻译方向
}
