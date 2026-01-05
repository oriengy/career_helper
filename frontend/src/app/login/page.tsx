'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { MessagePlugin } from 'tdesign-react';
import { useAuthStore } from '@/stores/auth';
import { useUserStore } from '@/stores/user';
import { authApi } from '@/services/api/auth';
import { validatePhone, validateVerifyCode } from '@/services/utils/validator';
import { ROUTES } from '@/constants/routes';
import { config } from '@/constants/config';

// Icons
const SmartphoneIcon = ({ className }: { className?: string }) => (
  <svg xmlns="http://www.w3.org/2000/svg" className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect width="14" height="20" x="5" y="2" rx="2" ry="2"/><path d="M12 18h.01"/></svg>
);
const ShieldCheckIcon = ({ className }: { className?: string }) => (
  <svg xmlns="http://www.w3.org/2000/svg" className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10"/><path d="m9 12 2 2 4-4"/></svg>
);

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuthStore();
  const { profile } = useUserStore();

  const [phone, setPhone] = useState('');
  const [verifyCode, setVerifyCode] = useState('');
  const [countdown, setCountdown] = useState(0);
  const [isLoggingIn, setIsLoggingIn] = useState(false);
  const [isSendingCode, setIsSendingCode] = useState(false);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

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
    <div className="min-h-screen w-full flex bg-slate-900 overflow-hidden">
      {/* Left Side - Dynamic Visuals (3/5) */}
      <div 
        className={`hidden lg:flex w-3/5 relative flex-col justify-center items-center overflow-hidden transition-all duration-1000 ease-out ${
          mounted ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-20'
        }`}
      >
        {/* Animated Background Layers */}
        <div className="absolute inset-0 bg-slate-900">
          <div className="absolute top-0 left-0 w-full h-full bg-[radial-gradient(ellipse_at_top_left,_var(--tw-gradient-stops))] from-indigo-900/40 via-slate-900 to-slate-900"></div>
          <div className="absolute top-[-20%] left-[-20%] w-[80%] h-[80%] bg-purple-600/20 rounded-full mix-blend-screen filter blur-3xl animate-blob"></div>
          <div className="absolute bottom-[-20%] right-[-20%] w-[80%] h-[80%] bg-blue-600/20 rounded-full mix-blend-screen filter blur-3xl animate-blob animation-delay-2000"></div>
          <div className="absolute top-[40%] left-[30%] w-[40%] h-[40%] bg-indigo-500/20 rounded-full mix-blend-screen filter blur-3xl animate-blob animation-delay-4000"></div>
        </div>

        {/* Content Overlay */}
        <div className="relative z-10 text-center px-12">
          <div className="mb-8 inline-block">
             <div className="w-20 h-20 bg-gradient-to-tr from-blue-400 to-purple-500 rounded-2xl rotate-3 shadow-lg shadow-blue-500/20 flex items-center justify-center">
                <span className="text-4xl">📚</span>
             </div>
          </div>
          <h2 className="text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-100 via-white to-purple-100 mb-6 drop-shadow-sm tracking-tight">
            探索职场无限可能
          </h2>
          <p className="text-lg text-blue-200/60 max-w-lg mx-auto leading-relaxed">
            职宝书 AI 助手为您提供专业的职业规划与咨询服务，<br/>让您的每一步职业发展都充满信心。
          </p>
          
          {/* Decorative Lines */}
          <div className="mt-12 flex justify-center gap-2">
            <div className="w-16 h-1 bg-gradient-to-r from-transparent via-blue-500/50 to-transparent rounded-full"></div>
            <div className="w-8 h-1 bg-gradient-to-r from-transparent via-purple-500/50 to-transparent rounded-full"></div>
          </div>
        </div>
      </div>

      {/* Right Side - Login Form (2/5) */}
      <div 
        className={`w-full lg:w-2/5 relative flex items-center justify-center p-8 bg-slate-950/30 backdrop-blur-sm border-l border-white/5 transition-all duration-1000 ease-out delay-200 ${
          mounted ? 'opacity-100 translate-x-0' : 'opacity-0 translate-x-20'
        }`}
      >
        <div className="w-full max-w-sm">
          {/* Header */}
          <div className="text-center mb-10 lg:text-left">
            <h1 className="text-3xl font-bold text-white mb-2">
              欢迎回来
            </h1>
            <p className="text-gray-400 text-sm">
              请登录您的账号以继续
            </p>
          </div>

          <div className="space-y-6">
            {/* Phone Input */}
            <div className="group">
              <label className="block text-xs font-medium text-gray-400 mb-1 ml-1 uppercase tracking-wider group-focus-within:text-blue-300 transition-colors">
                手机号
              </label>
              <div className="relative flex items-center">
                <div className="absolute left-4 text-gray-400 group-focus-within:text-blue-300 transition-colors">
                  <SmartphoneIcon className="w-5 h-5" />
                </div>
                <input
                  type="tel"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  placeholder="请输入手机号"
                  maxLength={11}
                  className="w-full bg-slate-800/50 text-white placeholder-gray-500 border border-white/10 rounded-xl py-4 pl-12 pr-4 outline-none focus:border-blue-500/50 focus:bg-slate-800 focus:ring-1 focus:ring-blue-500/50 transition-all duration-300"
                />
              </div>
            </div>

            {/* Verify Code Input */}
            <div className="group">
              <label className="block text-xs font-medium text-gray-400 mb-1 ml-1 uppercase tracking-wider group-focus-within:text-blue-300 transition-colors">
                验证码
              </label>
              <div className="relative flex gap-3">
                <div className="relative flex-1">
                  <div className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400 group-focus-within:text-blue-300 transition-colors">
                    <ShieldCheckIcon className="w-5 h-5" />
                  </div>
                  <input
                    type="text"
                    value={verifyCode}
                    onChange={(e) => setVerifyCode(e.target.value)}
                    placeholder="6位数字"
                    maxLength={6}
                    className="w-full bg-slate-800/50 text-white placeholder-gray-500 border border-white/10 rounded-xl py-4 pl-12 pr-4 outline-none focus:border-blue-500/50 focus:bg-slate-800 focus:ring-1 focus:ring-blue-500/50 transition-all duration-300"
                  />
                </div>
                <button
                  onClick={handleSendCode}
                  disabled={!validatePhone(phone) || countdown > 0 || isSendingCode}
                  className="px-6 min-w-[120px] bg-slate-800/50 border border-white/10 text-blue-200 text-sm font-medium rounded-xl hover:bg-slate-700 active:bg-slate-800 disabled:opacity-50 disabled:cursor-not-allowed transition-all backdrop-blur-sm whitespace-nowrap"
                >
                  {isSendingCode ? (
                    <span className="flex items-center justify-center">
                      <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      ...
                    </span>
                  ) : countdown > 0 ? (
                    `${countdown}s`
                  ) : (
                    '获取验证码'
                  )}
                </button>
              </div>
            </div>

            {/* Login Button */}
            <button
              onClick={handleLogin}
              disabled={isLoggingIn || !phone || !verifyCode}
              className="w-full mt-4 relative group overflow-hidden rounded-xl bg-gradient-to-r from-blue-600 to-indigo-600 p-px focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-offset-2 focus:ring-offset-slate-900 disabled:opacity-50 disabled:cursor-not-allowed transition-all hover:shadow-[0_0_20px_-5px_rgba(79,70,229,0.5)] transform hover:scale-[1.02] active:scale-[0.98]"
            >
              <span className="absolute inset-0 w-full h-full bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full group-hover:animate-shimmer" />
              <div className="relative flex items-center justify-center w-full h-full bg-transparent px-4 py-4 text-white font-bold tracking-wide">
                {isLoggingIn ? '登录中...' : '立即登录'}
              </div>
            </button>
          </div>

          {/* Footer */}
          <div className="mt-12 text-center lg:text-left">
            <p className="text-xs text-gray-500">
              登录即表示同意{' '}
              <a href="#" className="text-blue-400 hover:text-blue-300 hover:underline transition-colors">
                用户协议
              </a>{' '}
              和{' '}
              <a href="#" className="text-blue-400 hover:text-blue-300 hover:underline transition-colors">
                隐私政策
              </a>
            </p>
          </div>

          {/* Mock Tip */}
          {config.mockLogin && (
            <div className="mt-6 p-3 bg-yellow-500/10 border border-yellow-500/20 rounded-lg backdrop-blur-md">
              <p className="text-xs text-yellow-200/80 font-mono text-center">
                DEV: Code=Any | Phone=138xxxxxxxx
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}