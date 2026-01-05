/**
 * 会话相关 API
 */

import { apiClient } from './client';
import type {
  ListChatSessionsParams,
  ListChatSessionsResponse,
  CreateChatSessionParams,
  CreateChatSessionResponse,
  UpdateChatSessionParams,
  DeleteChatSessionParams,
} from '@/types/api';
import type { ChatSession } from '@/types/models';

const mapChatSession = (session: ChatSession): ChatSession => {
  const sessionId = session.sessionId || session.id;
  const name = session.friendName || session.name;
  const avatar = session.friendAvatar || session.avatar;
  const gender = session.friendGender || session.gender;

  return {
    ...session,
    sessionId,
    friendName: name,
    friendAvatar: avatar,
    friendGender: gender,
  };
};

export const sessionApi = {
  /**
   * 获取会话列表
   */
  async listSessions(
    params: ListChatSessionsParams = {}
  ): Promise<ChatSession[]> {
    const response = await apiClient.post<ListChatSessionsResponse>(
      '/chat.ChatService/ListChatSessions',
      params
    );
    return (response.data.data || []).map(mapChatSession);
  },

  /**
   * 创建会话
   */
  async createSession(
    params: CreateChatSessionParams
  ): Promise<ChatSession> {
    const response = await apiClient.post<CreateChatSessionResponse>(
      '/chat.ChatService/CreateChatSession',
      params
    );
    return mapChatSession(response.data.chatSession);
  },

  /**
   * 更新会话
   */
  async updateSession(params: UpdateChatSessionParams): Promise<ChatSession> {
    const response = await apiClient.post<{ chatSession: ChatSession }>(
      '/chat.ChatService/UpdateChatSession',
      params
    );
    return mapChatSession(response.data.chatSession);
  },

  /**
   * 删除会话
   */
  async deleteSession(sessionId: string): Promise<void> {
    await apiClient.post('/chat.ChatService/DeleteChatSession', { id: sessionId });
  },
};
