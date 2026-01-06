/**
 * 文件上传 API
 */

import { apiClient } from './client';
import type { UploadFileResponse } from '@/types/api';

export interface UploadOptions {
  usageType?: 'avatar' | 'chat_image' | 'temp_upload';
  onProgress?: (percent: number) => void;
}

/**
 * 上传文件到服务器
 */
export async function uploadFile(
  file: File,
  options: UploadOptions = {}
): Promise<UploadFileResponse> {
  const { usageType, onProgress } = options;

  // 创建 FormData
  const formData = new FormData();
  formData.append('file', file);

  // 发送请求
  const response = await apiClient.post<{ data: UploadFileResponse }>(
    '/file/wx_upload',
    formData,
    {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      params: usageType ? { usage_type: usageType } : undefined,
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const percent = Math.round(
            (progressEvent.loaded * 100) / progressEvent.total
          );
          onProgress(percent);
        }
      },
    }
  );

  return response.data.data;
}

/**
 * 上传头像（自动压缩）
 */
export async function uploadAvatar(
  file: File,
  onProgress?: (percent: number) => void
): Promise<string> {
  // 1. 压缩图片
  const compressedFile = await compressImage(file, {
    maxWidth: 800,
    maxHeight: 800,
    quality: 0.8,
  });

  // 2. 上传
  const response = await uploadFile(compressedFile, {
    usageType: 'avatar',
    onProgress,
  });

  // 3. 返回公共 URL
  return response;
}

/**
 * 上传聊天图片（自动压缩）
 */
export async function uploadChatImage(
  file: File,
  onProgress?: (percent: number) => void
): Promise<UploadFileResponse> {
  // 1. 压缩图片
  const compressedFile = await compressImage(file, {
    maxWidth: 1200,
    maxHeight: 1200,
    quality: 0.85,
  });

  // 2. 上传
  const response = await uploadFile(compressedFile, {
    usageType: 'chat_image',
    onProgress,
  });

  // 3. 返回公共 URL
  return response;
}

/**
 * 压缩图片
 */
async function compressImage(
  file: File,
  options: {
    maxWidth: number;
    maxHeight: number;
    quality: number;
  }
): Promise<File> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    reader.onload = (e) => {
      const img = new Image();

      img.onload = () => {
        // 计算压缩后的尺寸
        let { width, height } = img;
        const { maxWidth, maxHeight } = options;

        if (width > maxWidth || height > maxHeight) {
          const ratio = Math.min(maxWidth / width, maxHeight / height);
          width *= ratio;
          height *= ratio;
        }

        // 创建 canvas
        const canvas = document.createElement('canvas');
        canvas.width = width;
        canvas.height = height;

        const ctx = canvas.getContext('2d');
        if (!ctx) {
          reject(new Error('无法创建 Canvas 上下文'));
          return;
        }

        // 绘制图片
        ctx.drawImage(img, 0, 0, width, height);

        // 转换为 Blob
        canvas.toBlob(
          (blob) => {
            if (!blob) {
              reject(new Error('图片压缩失败'));
              return;
            }

            // 创建新的 File 对象
            const compressedFile = new File([blob], file.name, {
              type: 'image/jpeg',
              lastModified: Date.now(),
            });

            resolve(compressedFile);
          },
          'image/jpeg',
          options.quality
        );
      };

      img.onerror = () => reject(new Error('图片加载失败'));
      img.src = e.target?.result as string;
    };

    reader.onerror = () => reject(new Error('文件读取失败'));
    reader.readAsDataURL(file);
  });
}

/**
 * 上传 API 对象（统一导出）
 */
export const uploadApi = {
  uploadFile,
  uploadAvatar,
  uploadChatImage,
};
