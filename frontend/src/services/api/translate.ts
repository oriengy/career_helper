/**
 * 翻译相关 API
 */

import { apiClient } from './client';
import type { TranslateParams, TranslateResponse } from '@/types/api';

export const translateApi = {
  /**
   * 翻译消息
   */
  async translateMessage(
    params: TranslateParams
  ): Promise<TranslateResponse> {
    const response = await apiClient.post<TranslateResponse>(
      '/translate.TranslateService/TranslateV2',
      {
        chatSessionId: params.chatSessionId,
        targetMessageId: params.targetMessageId,
      }
    );
    return response.data;
  },
};
