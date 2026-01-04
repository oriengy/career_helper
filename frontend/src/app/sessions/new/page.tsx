'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Input, Button, Message, Radio, Upload, Avatar } from 'tdesign-react';
import { ArrowLeftIcon, UploadIcon } from 'tdesign-icons-react';
import { useUserStore } from '@/stores/user';
import { sessionApi } from '@/services/api/session';
import { uploadApi } from '@/services/api/upload';
import { getUserAvatar } from '@/lib/avatar';
import { ROUTES } from '@/constants/routes';
import type { Gender } from '@/types/models';
import type { UploadFile } from 'tdesign-react';

export default function NewSessionPage() {
  const router = useRouter();
  const { profile } = useUserStore();
  const [friendName, setFriendName] = useState('');
  const [friendGender, setFriendGender] = useState<Gender>('');
  const [friendAvatar, setFriendAvatar] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);

  // 处理头像上传
  const handleAvatarUpload = async (file: UploadFile) => {
    if (!file.raw) {
      return { status: 'fail', error: '无效的文件' };
    }

    try {
      const url = await uploadApi.uploadAvatar(file.raw, (percent) => {
        setUploadProgress(percent);
      });
      setFriendAvatar(url);
      Message.success('头像上传成功');
      return { status: 'success', url };
    } catch (error: any) {
      Message.error(error.message || '头像上传失败');
      return { status: 'fail', error: error.message };
    }
  };

  // 创建会话
  const handleSubmit = async () => {
    if (!friendName.trim()) {
      Message.error('请输入对方昵称');
      return;
    }

    if (!friendGender) {
      Message.error('请选择对方性别');
      return;
    }

    setIsSubmitting(true);
    try {
      const session = await sessionApi.createSession({
        friendName: friendName.trim(),
        friendGender,
        friendAvatar: friendAvatar || undefined,
      });

      Message.success('创建成功');
      router.push(ROUTES.SESSION_DETAIL(session.sessionId));
    } catch (error: any) {
      Message.error(error.message || '创建失败');
    } finally {
      setIsSubmitting(false);
    }
  };

  // 获取预览头像
  const getPreviewAvatar = () => {
    if (friendAvatar) {
      return friendAvatar;
    }
    return getUserAvatar('', friendGender, profile?.gender, true);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 头部 */}
      <div className="bg-white border-b border-gray-200 sticky top-0 z-10">
        <div className="max-w-2xl mx-auto px-4 py-4 flex items-center gap-4">
          <Button
            variant="text"
            shape="circle"
            icon={<ArrowLeftIcon />}
            onClick={() => router.back()}
          />
          <h1 className="text-xl font-semibold text-gray-800">新建会话</h1>
        </div>
      </div>

      {/* 表单内容 */}
      <div className="max-w-2xl mx-auto px-4 py-6">
        <div className="bg-white rounded-lg shadow-sm p-6 space-y-6">
          {/* 头像上传 */}
          <div className="flex flex-col items-center">
            <div className="mb-4">
              <Avatar size="80px" image={getPreviewAvatar()} />
            </div>
            <Upload
              theme="custom"
              accept="image/*"
              requestMethod={handleAvatarUpload}
              showUploadProgress={false}
              max={1}
            >
              <Button variant="outline" icon={<UploadIcon />}>
                {friendAvatar ? '更换头像' : '上传头像（可选）'}
              </Button>
            </Upload>
            {uploadProgress > 0 && uploadProgress < 100 && (
              <p className="text-sm text-gray-500 mt-2">上传中... {uploadProgress}%</p>
            )}
          </div>

          {/* 昵称输入 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              对方昵称 <span className="text-red-500">*</span>
            </label>
            <Input
              value={friendName}
              onChange={(value) => setFriendName(value)}
              placeholder="请输入对方的昵称"
              maxlength={20}
              clearable
              size="large"
            />
          </div>

          {/* 性别选择 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              对方性别 <span className="text-red-500">*</span>
            </label>
            <Radio.Group
              value={friendGender}
              onChange={(value) => setFriendGender(value as Gender)}
              variant="default-filled"
            >
              <Radio.Button value="male">男性</Radio.Button>
              <Radio.Button value="female">女性</Radio.Button>
            </Radio.Group>
          </div>

          {/* 提示信息 */}
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <p className="text-sm text-blue-800">
              💡 填写对方信息后，AI 将根据性别提供更精准的翻译和建议
            </p>
          </div>

          {/* 提交按钮 */}
          <div className="flex gap-3">
            <Button
              block
              size="large"
              variant="outline"
              onClick={() => router.back()}
              disabled={isSubmitting}
            >
              取消
            </Button>
            <Button
              block
              size="large"
              theme="primary"
              onClick={handleSubmit}
              loading={isSubmitting}
              disabled={!friendName.trim() || !friendGender}
            >
              创建会话
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
