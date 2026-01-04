/**
 * 消息相关 API
 */

import { apiClient } from './client';
import type {
  ListMessagesParams,
  ListMessagesResponse,
  CreateMessageParams,
  CreateMessageResponse,
  SendConsultMessageParams,
  SendConsultMessageResponse,
} from '@/types/api';

export const messageApi = {
  /**
   * 获取消息列表
   */
  async listMessages(
    params: ListMessagesParams
  ): Promise<ListMessagesResponse> {
    const response = await apiClient.post<ListMessagesResponse>(
      '/message.ChatMessageService/ListChatMessages',
      params
    );
    return response.data;
  },

  /**
   * 创建消息
   */
  async createMessage(
    params: CreateMessageParams
  ): Promise<CreateMessageResponse> {
    const response = await apiClient.post<CreateMessageResponse>(
      '/message.ChatMessageService/CreateChatMessage',
      params
    );
    return response.data;
  },

  /**
   * 发送咨询消息
   */
  async sendConsultMessage(
    params: SendConsultMessageParams
  ): Promise<SendConsultMessageResponse> {
    const response = await apiClient.post<SendConsultMessageResponse>(
      '/message.ChatMessageService/SendConsultMessage',
      params
    );
    return response.data;
  },

  /**
   * 删除消息
   */
  async deleteMessage(id: string): Promise<void> {
    await apiClient.post('/message.ChatMessageService/DeleteChatMessage', {
      id,
    });
  },

  /**
   * 图片解析
   */
  async parseImageMessages(sessionId: string, imageUrl: string): Promise<any> {
    const response = await apiClient.post(
      '/message.FriendMessageService/ParseImageMessages',
      {
        sessionId,
        imageUrl,
      }
    );
    return response.data;
  },
};
