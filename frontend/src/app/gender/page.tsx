'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { MessagePlugin } from 'tdesign-react';
import { useUserStore } from '@/stores/user';
import { profileApi } from '@/services/api/profile';
import { ROUTES } from '@/constants/routes';
import type { Gender } from '@/types/models';

export default function GenderPage() {
  const router = useRouter();
  const { profile, setProfile } = useUserStore();
  const [selectedGender, setSelectedGender] = useState<Gender>('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const handleSubmit = async () => {
    if (!selectedGender) {
      MessagePlugin.error('请选择性别');
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

      MessagePlugin.success('设置成功');
      router.push(ROUTES.SESSIONS);
    } catch (error: any) {
      MessagePlugin.error(error.message || '设置失败');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen w-full flex items-center justify-center bg-slate-900 relative overflow-hidden">
      {/* Background Decor - Same as Login for consistency */}
      <div className="absolute top-0 left-0 w-full h-full overflow-hidden z-0 pointer-events-none">
        <div className="absolute top-[10%] left-[20%] w-[30%] h-[30%] bg-purple-600/20 rounded-full mix-blend-screen filter blur-3xl animate-blob"></div>
        <div className="absolute top-[20%] right-[20%] w-[30%] h-[30%] bg-blue-600/20 rounded-full mix-blend-screen filter blur-3xl animate-blob animation-delay-2000"></div>
        <div className="absolute -bottom-10 left-[40%] w-[40%] h-[40%] bg-indigo-600/20 rounded-full mix-blend-screen filter blur-3xl animate-blob animation-delay-4000"></div>
      </div>

      <div className={`w-full max-w-md z-10 px-4 transition-all duration-1000 ease-out transform ${
        mounted ? 'opacity-100 translate-y-0 scale-100' : 'opacity-0 translate-y-10 scale-95'
      }`}>
        {/* Header Text */}
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-white mb-2 tracking-wide drop-shadow-md">
            选择性别
          </h1>
          <p className="text-blue-200/60 text-sm font-light">
            AI 助手将根据您的性别提供更加个性化的服务
          </p>
        </div>

        {/* Main Card */}
        <div className="relative bg-white/10 backdrop-blur-xl border border-white/20 rounded-3xl shadow-2xl p-8 overflow-hidden">
           {/* Shine effect */}
           <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-white/30 to-transparent opacity-50"></div>

          <div className="space-y-4 mb-8">
            {/* Male Option */}
            <div
              onClick={() => setSelectedGender('male')}
              className={`
                group relative cursor-pointer rounded-2xl p-4 border transition-all duration-300
                ${
                  selectedGender === 'male'
                    ? 'border-blue-500 bg-blue-600/20 shadow-[0_0_20px_rgba(59,130,246,0.3)]'
                    : 'border-white/10 hover:border-blue-400/50 hover:bg-white/5'
                }
              `}
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-5">
                  <div
                    className={`
                    text-3xl w-14 h-14 rounded-full flex items-center justify-center transition-colors duration-300
                    ${selectedGender === 'male' ? 'bg-blue-500 text-white' : 'bg-white/10 text-gray-400 group-hover:bg-white/20 group-hover:text-blue-300'}
                  `}
                  >
                    👨
                  </div>
                  <div>
                    <h3 className={`text-lg font-semibold transition-colors duration-300 ${selectedGender === 'male' ? 'text-white' : 'text-gray-300 group-hover:text-white'}`}>
                      男性
                    </h3>
                    <p className="text-xs text-gray-500 group-hover:text-gray-400">Male</p>
                  </div>
                </div>
                {/* Checkmark */}
                <div className={`
                    w-6 h-6 rounded-full border-2 flex items-center justify-center transition-all duration-300
                    ${selectedGender === 'male' ? 'border-blue-500 bg-blue-500' : 'border-gray-600 group-hover:border-blue-400/50'}
                `}>
                    {selectedGender === 'male' && (
                        <svg className="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                        </svg>
                    )}
                </div>
              </div>
            </div>

            {/* Female Option */}
            <div
              onClick={() => setSelectedGender('female')}
              className={`
                group relative cursor-pointer rounded-2xl p-4 border transition-all duration-300
                ${
                  selectedGender === 'female'
                    ? 'border-pink-500 bg-pink-600/20 shadow-[0_0_20px_rgba(236,72,153,0.3)]'
                    : 'border-white/10 hover:border-pink-400/50 hover:bg-white/5'
                }
              `}
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-5">
                  <div
                    className={`
                    text-3xl w-14 h-14 rounded-full flex items-center justify-center transition-colors duration-300
                    ${selectedGender === 'female' ? 'bg-pink-500 text-white' : 'bg-white/10 text-gray-400 group-hover:bg-white/20 group-hover:text-pink-300'}
                  `}
                  >
                    👩
                  </div>
                  <div>
                    <h3 className={`text-lg font-semibold transition-colors duration-300 ${selectedGender === 'female' ? 'text-white' : 'text-gray-300 group-hover:text-white'}`}>
                      女性
                    </h3>
                    <p className="text-xs text-gray-500 group-hover:text-gray-400">Female</p>
                  </div>
                </div>
                {/* Checkmark */}
                <div className={`
                    w-6 h-6 rounded-full border-2 flex items-center justify-center transition-all duration-300
                    ${selectedGender === 'female' ? 'border-pink-500 bg-pink-500' : 'border-gray-600 group-hover:border-pink-400/50'}
                `}>
                    {selectedGender === 'female' && (
                        <svg className="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                        </svg>
                    )}
                </div>
              </div>
            </div>
          </div>

          {/* Confirm Button */}
          <button
            onClick={handleSubmit}
            disabled={!selectedGender || isSubmitting}
            className={`
                w-full relative group overflow-hidden rounded-xl p-px focus:outline-none transition-all duration-300
                ${!selectedGender || isSubmitting 
                    ? 'opacity-50 cursor-not-allowed bg-gray-700' 
                    : 'bg-gradient-to-r from-blue-600 to-indigo-600 hover:shadow-[0_0_20px_-5px_rgba(79,70,229,0.5)] transform hover:scale-[1.02] active:scale-[0.98]'
                }
            `}
          >
             {/* Gradient Background for active state */}
             {!(!selectedGender || isSubmitting) && (
                <span className="absolute inset-0 w-full h-full bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full group-hover:animate-shimmer" />
             )}
            
            <div className="relative flex items-center justify-center w-full h-full bg-transparent px-4 py-4 text-white font-bold tracking-wide">
              {isSubmitting ? (
                 <span className="flex items-center">
                    <svg className="animate-spin -ml-1 mr-2 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    提交中...
                 </span>
              ) : (
                '确认开启旅程'
              )}
            </div>
          </button>

          {/* Footer Note */}
          <div className="mt-6 text-center">
            <p className="text-xs text-gray-500">
              * 此信息仅用于 AI 个性化调优，严格保密
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}