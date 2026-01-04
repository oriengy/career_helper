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
        messageId: params.messageId,
        direction: params.direction,
      }
    );
    return response.data;
  },
};
