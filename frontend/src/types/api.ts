/**
 * API 相关类型定义
 */

import type { User, Profile, ChatSession, Message } from './models';

// 通用响应类型
export interface ApiResponse<T = any> {
  code?: number;
  message?: string;
  data?: T;
}

// 分页请求参数
export interface PageRequest {
  pageToken?: string;
  pageSize?: number;
}

// 分页响应
export interface PageResponse<T> {
  data: T[];
  nextPageToken?: string;
  total?: number;
}

// 登录相关
export interface SendVerifyCodeParams {
  phone: string;
}

export interface PhoneLoginParams {
  phone: string;
  verificationCode: string;
}

export interface PhoneLoginResponse {
  token: string;
  user: User;
  profile: Profile;
}

export interface UserProfileResponse {
  user: User;
  profile: Profile;
}

// 会话相关
export interface ListChatSessionsParams extends PageRequest {}

export interface ListChatSessionsResponse extends PageResponse<ChatSession> {}

export interface CreateChatSessionParams {
  profile: {
    name: string;
    gender?: string;
    avatar?: string;
  };
}

export interface CreateChatSessionResponse {
  chatSession: ChatSession;
  profile: Profile;
}

export interface UpdateChatSessionParams {
  id: string;
  name?: string;
  avatar?: string;
}

export interface DeleteChatSessionParams {
  id: string;
}

// 消息相关
export interface ListMessagesParams extends PageRequest {
  sessionId: string;
  msgType?: string;
  ids?: string[];
}

export interface ListMessagesResponse extends PageResponse<Message> {}

export interface CreateMessageParams {
  sessionId: string;
  role: string;
  content: string;
  msgAt?: string;
}

export interface CreateMessageResponse {
  message: Message;
}

export interface SendConsultMessageParams {
  sessionId: string;
  content: string;
}

export interface SendConsultMessageResponse {
  message: Message;
}

// 翻译相关
export interface TranslateParams {
  messageId: string;
  direction: 'male' | 'female';
}

export interface TranslateResponse {
  message: Message;
}

// 资料相关
export interface CreateProfileParams {
  name?: string;
  imName?: string;
  avatar?: string;
  gender?: string;
  age?: number;
  birthday?: string;
  birthLocation?: string;
  currentLocation?: string;
}

export interface UpdateProfileParams extends CreateProfileParams {
  id: string;
}

export interface GetProfileParams {
  id: string;
}

export interface ProfileResponse {
  profile: Profile;
}

// 文件上传
export interface UploadFileResponse {
  publicUrl: string;
  privateUrl?: string;
  fileId?: string;
  size?: number;
}

// 配置
export interface ConfigItem {
  key: string;
  value: string;
}

export interface GetConfigParams {
  keys: string[];
  app: string;
  platform: string;
  env: string;
  version: string;
}

export interface GetConfigResponse {
  configs: ConfigItem[];
}
