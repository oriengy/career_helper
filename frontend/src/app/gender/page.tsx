'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button, Message } from 'tdesign-react';
import { useUserStore } from '@/stores/user';
import { profileApi } from '@/services/api/profile';
import { ROUTES } from '@/constants/routes';
import type { Gender } from '@/types/models';

export default function GenderPage() {
  const router = useRouter();
  const { profile, setProfile } = useUserStore();
  const [selectedGender, setSelectedGender] = useState<Gender>('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async () => {
    if (!selectedGender) {
      Message.error('请选择性别');
      return;
    }

    setIsSubmitting(true);
    try {
      if (profile) {
        // 更新已存在的 profile
        const updatedProfile = await profileApi.updateProfile({
          gender: selectedGender,
        });
        setProfile(updatedProfile);
      } else {
        // 创建新的 profile
        const newProfile = await profileApi.createProfile({
          gender: selectedGender,
        });
        setProfile(newProfile);
      }

      Message.success('设置成功');
      router.push(ROUTES.SESSIONS);
    } catch (error: any) {
      Message.error(error.message || '设置失败');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100 px-4">
      <div className="w-full max-w-md">
        {/* 标题 */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-800 mb-2">选择性别</h1>
          <p className="text-gray-600">为了提供更好的服务，请选择您的性别</p>
        </div>

        {/* 性别选择卡片 */}
        <div className="bg-white rounded-2xl shadow-xl p-8">
          <div className="space-y-4 mb-8">
            {/* 男性选项 */}
            <div
              onClick={() => setSelectedGender('male')}
              className={`
                relative cursor-pointer rounded-xl p-6 border-2 transition-all
                ${
                  selectedGender === 'male'
                    ? 'border-primary-500 bg-primary-50'
                    : 'border-gray-200 hover:border-primary-300'
                }
              `}
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <div
                    className={`
                    text-4xl w-16 h-16 rounded-full flex items-center justify-center
                    ${selectedGender === 'male' ? 'bg-primary-100' : 'bg-gray-100'}
                  `}
                  >
                    👨
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold text-gray-800">男性</h3>
                    <p className="text-sm text-gray-500">Male</p>
                  </div>
                </div>
                {selectedGender === 'male' && (
                  <div className="text-primary-500 text-2xl">✓</div>
                )}
              </div>
            </div>

            {/* 女性选项 */}
            <div
              onClick={() => setSelectedGender('female')}
              className={`
                relative cursor-pointer rounded-xl p-6 border-2 transition-all
                ${
                  selectedGender === 'female'
                    ? 'border-primary-500 bg-primary-50'
                    : 'border-gray-200 hover:border-primary-300'
                }
              `}
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <div
                    className={`
                    text-4xl w-16 h-16 rounded-full flex items-center justify-center
                    ${selectedGender === 'female' ? 'bg-primary-100' : 'bg-gray-100'}
                  `}
                  >
                    👩
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold text-gray-800">女性</h3>
                    <p className="text-sm text-gray-500">Female</p>
                  </div>
                </div>
                {selectedGender === 'female' && (
                  <div className="text-primary-500 text-2xl">✓</div>
                )}
              </div>
            </div>
          </div>

          {/* 确认按钮 */}
          <Button
            onClick={handleSubmit}
            loading={isSubmitting}
            disabled={!selectedGender}
            block
            size="large"
            theme="primary"
          >
            确认
          </Button>
        </div>

        {/* 说明文字 */}
        <div className="mt-6 text-center">
          <p className="text-sm text-gray-500">
            性别信息将用于提供个性化的翻译和建议服务
          </p>
        </div>
      </div>
    </div>
  );
}
