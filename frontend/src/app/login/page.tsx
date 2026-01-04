'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Input, Button, MessagePlugin } from 'tdesign-react';
import { useAuthStore } from '@/stores/auth';
import { useUserStore } from '@/stores/user';
import { authApi } from '@/services/api/auth';
import { validatePhone, validateVerifyCode } from '@/services/utils/validator';
import { ROUTES } from '@/constants/routes';
import { config } from '@/constants/config';

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuthStore();
  const { profile } = useUserStore();

  const [phone, setPhone] = useState('');
  const [verifyCode, setVerifyCode] = useState('');
  const [countdown, setCountdown] = useState(0);
  const [isLoggingIn, setIsLoggingIn] = useState(false);
  const [isSendingCode, setIsSendingCode] = useState(false);

  // 倒计时逻辑
  useEffect(() => {
    if (countdown > 0) {
      const timer = setTimeout(() => {
        setCountdown(countdown - 1);
      }, 1000);
      return () => clearTimeout(timer);
    }
  }, [countdown]);

  // 发送验证码
  const handleSendCode = async () => {
    // 验证手机号
    if (!validatePhone(phone)) {
      MessagePlugin.error('请输入正确的手机号');
      return;
    }

    setIsSendingCode(true);
    try {
      await authApi.sendVerifyCode({ phone });
      MessagePlugin.success('验证码已发送');
      setCountdown(60);
    } catch (error: any) {
      MessagePlugin.error(error.message || '发送验证码失败');
    } finally {
      setIsSendingCode(false);
    }
  };

  // 登录处理
  const handleLogin = async () => {
    // 验证输入
    if (!phone || !verifyCode) {
      MessagePlugin.error('请填写完整信息');
      return;
    }

    if (!validatePhone(phone)) {
      MessagePlugin.error('请输入正确的手机号');
      return;
    }

    if (!config.mockLogin && !validateVerifyCode(verifyCode)) {
      MessagePlugin.error('请输入6位验证码');
      return;
    }

    setIsLoggingIn(true);
    try {
      await login(phone, verifyCode);
      MessagePlugin.success('登录成功');

      // 判断是否需要性别选择
      const currentProfile = useUserStore.getState().profile;
      if (!currentProfile?.gender) {
        router.push(ROUTES.GENDER);
      } else {
        router.push(ROUTES.SESSIONS);
      }
    } catch (error: any) {
      MessagePlugin.error(error.message || '登录失败');
    } finally {
      setIsLoggingIn(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100 px-4">
      <div className="w-full max-w-md">
        {/* Logo 和标题 */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-800 mb-2">恋爱翻译官</h1>
          <p className="text-gray-600">AI 驱动的社交辅助应用</p>
        </div>

        {/* 登录表单 */}
        <div className="bg-white rounded-2xl shadow-xl p-8">
          <h2 className="text-xl font-semibold text-gray-800 mb-6">欢迎登录</h2>

          {/* 手机号输入 */}
          <div className="mb-4">
            <Input
              value={phone}
              onChange={(value) => setPhone(value)}
              placeholder="请输入手机号"
              maxlength={11}
              clearable
              size="large"
            />
          </div>

          {/* 验证码输入 */}
          <div className="mb-4">
            <div className="flex gap-2">
              <Input
                value={verifyCode}
                onChange={(value) => setVerifyCode(value)}
                placeholder="请输入验证码"
                maxlength={6}
                clearable
                size="large"
                className="flex-1"
              />
              <Button
                onClick={handleSendCode}
                disabled={!validatePhone(phone) || countdown > 0 || isSendingCode}
                loading={isSendingCode}
                size="large"
                variant="outline"
                className="w-32"
              >
                {countdown > 0 ? `${countdown}秒后重试` : '发送验证码'}
              </Button>
            </div>
          </div>

          {/* 登录按钮 */}
          <Button
            onClick={handleLogin}
            loading={isLoggingIn}
            disabled={!phone || !verifyCode}
            block
            size="large"
            theme="primary"
            className="mb-4"
          >
            登录
          </Button>

          {/* 用户协议 */}
          <p className="text-xs text-gray-500 text-center">
            登录即表示同意
            <a href="#" className="text-primary-500 hover:underline">
              《用户协议》
            </a>
            和
            <a href="#" className="text-primary-500 hover:underline">
              《隐私政策》
            </a>
          </p>
        </div>

        {/* 开发环境提示 */}
        {config.mockLogin && (
          <div className="mt-4 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
            <p className="text-sm text-yellow-800">
              <strong>测试登录提示：</strong>
              <br />
              验证码：任意（mock）
              <br />
              测试手机号：138xxxxxxxx
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
