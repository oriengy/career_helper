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

export const sessionApi = {
  /**
   * 获取会话列表
   */
  async listSessions(
    params: ListChatSessionsParams
  ): Promise<ListChatSessionsResponse> {
    const response = await apiClient.post<ListChatSessionsResponse>(
      '/chat.ChatService/ListChatSessions',
      params
    );
    return response.data;
  },

  /**
   * 创建会话
   */
  async createSession(
    params: CreateChatSessionParams
  ): Promise<CreateChatSessionResponse> {
    const response = await apiClient.post<CreateChatSessionResponse>(
      '/chat.ChatService/CreateChatSession',
      params
    );
    return response.data;
  },

  /**
   * 更新会话
   */
  async updateSession(params: UpdateChatSessionParams): Promise<ChatSession> {
    const response = await apiClient.post<{ chatSession: ChatSession }>(
      '/chat.ChatService/UpdateChatSession',
      params
    );
    return response.data.chatSession;
  },

  /**
   * 删除会话
   */
  async deleteSession(params: DeleteChatSessionParams): Promise<void> {
    await apiClient.post('/chat.ChatService/DeleteChatSession', params);
  },
};
