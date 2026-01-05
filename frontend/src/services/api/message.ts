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
   * 标准化消息字段，兼容后端返回
   */
  mapMessage(message: import('@/types/models').Message): import('@/types/models').Message {
    const rawMessage = message as any;
    const messageId = message.messageId || message.id || rawMessage.message_id;
    const sessionId = message.sessionId || rawMessage.sessionId || rawMessage.session_id;
    const parentId = message.parentId || rawMessage.parentId || rawMessage.parent_id;
    const msgType = message.msgType || rawMessage.msgType || rawMessage.msg_type || rawMessage.messageType;
    const tags = message.tags || rawMessage.tags || [];
    const content = message.content ?? rawMessage.content;
    const imageUrl = message.imageUrl || (tags.includes('image') ? content : undefined);

    return {
      ...message,
      messageId,
      sessionId,
      parentId,
      msgType,
      tags,
      content,
      imageUrl,
    };
  },

  /**
   * 获取消息列表
   */
  async listMessages(
    params: ListMessagesParams
  ): Promise<import('@/types/models').Message[]> {
    const response = await apiClient.post<ListMessagesResponse>(
      '/message.ChatMessageService/ListChatMessages',
      params
    );
    return (response.data.messages || []).map((msg) => messageApi.mapMessage(msg));
  },

  /**
   * 创建消息
   */
  async createMessage(
    params: CreateMessageParams
  ): Promise<CreateMessageResponse> {
    const tags = params.tags ? [...params.tags] : [];
    const content = params.imageUrl ? params.imageUrl : params.content;
    if (params.imageUrl && !tags.includes('image')) {
      tags.push('image');
    }

    const response = await apiClient.post<CreateMessageResponse>(
      '/message.ChatMessageService/CreateChatMessage',
      {
        messages: [
          {
            sessionId: params.sessionId,
            parentId: params.parentId,
            role: params.role,
            msgType: params.msgType || 'HISTORY',
            content,
            tags,
          },
        ],
      }
    );
    return {
      messages: (response.data.messages || []).map((msg) =>
        messageApi.mapMessage(msg)
      ),
    };
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
      ids: [id],
    });
  },

  /**
   * 图片解析
   */
  async parseImageMessages(sessionId: string, imageUrl: string): Promise<any> {
    const response = await apiClient.post(
      '/message.ChatMessageService/ParseImageMessages',
      {
        sessionId,
        imageUrl,
      }
    );
    return response.data;
  },
};
