'use client';

import { useState, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { MessagePlugin, Avatar } from 'tdesign-react';
import { UploadIcon, UserIcon, ArrowLeftIcon } from 'tdesign-icons-react';
import { useUserStore } from '@/stores/user';
import { sessionApi } from '@/services/api/session';
import { profileApi } from '@/services/api/profile';
import { uploadApi } from '@/services/api/upload';
import { getUserAvatar } from '@/lib/avatar';
import { ROUTES } from '@/constants/routes';
type Relation = '' | '部门同事' | '直属上级' | '下属' | '合作方';

export default function NewSessionPage() {
  const router = useRouter();
  const { profile } = useUserStore();
  const [friendName, setFriendName] = useState('');
  const [friendRelation, setFriendRelation] = useState<Relation>('');
  const [friendAvatar, setFriendAvatar] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isUploading, setIsUploading] = useState(false);

  const fileInputRef = useRef<HTMLInputElement>(null);

  // Handle Avatar Upload
  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setIsUploading(true);
    try {
      const url = await uploadApi.uploadAvatar(file);
      setFriendAvatar(url);
      MessagePlugin.success('头像上传成功');
    } catch (error: any) {
      MessagePlugin.error(error.message || '头像上传失败');
    } finally {
      setIsUploading(false);
      // Reset input
      if (fileInputRef.current) fileInputRef.current.value = '';
    }
  };

  // Create Session
  const handleSubmit = async () => {
    if (!friendName.trim()) {
      MessagePlugin.error('请输入对方姓名');
      return;
    }

    if (!friendRelation) {
      MessagePlugin.error('请选择对方关系');
      return;
    }

    setIsSubmitting(true);
    try {
      const session = await sessionApi.createSession({
        profile: {
          name: friendName.trim(),
          avatar: friendAvatar || undefined,
        },
      });

      if (session.profileId) {
        try {
          await profileApi.updateProfile({
            id: session.profileId,
            custom: [{ name: '关系', value: friendRelation }],
          });
        } catch (error: any) {
          console.warn('Failed to update relation:', error);
        }
      }

      MessagePlugin.success('创建成功');
      router.push(ROUTES.SESSION_DETAIL(session.sessionId));
    } catch (error: any) {
      MessagePlugin.error(error.message || '创建失败');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Preview Avatar logic
  const previewImage = friendAvatar || getUserAvatar('', '', profile?.gender, true);

  return (
    <div className="h-full flex flex-col items-center justify-center p-4 bg-slate-900 relative overflow-hidden">
        {/* Mobile Header */}
        <div className="lg:hidden absolute top-0 left-0 p-4">
             <button onClick={() => router.back()} className="text-slate-400 hover:text-white">
                 <ArrowLeftIcon />
             </button>
        </div>

        {/* Decorative Background */}
        <div className="absolute top-[-20%] right-[-20%] w-[500px] h-[500px] bg-blue-600/10 rounded-full mix-blend-screen filter blur-[120px] pointer-events-none"></div>
        <div className="absolute bottom-[-20%] left-[-20%] w-[500px] h-[500px] bg-purple-600/10 rounded-full mix-blend-screen filter blur-[120px] pointer-events-none"></div>

      <div className="w-full max-w-lg z-10 animate-fade-in-up">
        <div className="text-center mb-8">
             <h1 className="text-3xl font-bold text-white mb-2">新建对话</h1>
             <p className="text-slate-400 text-sm">填写对方信息，AI 将为您提供更精准的辅助</p>
        </div>

        <div className="bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-8 shadow-2xl">
          {/* Avatar Section */}
          <div className="flex flex-col items-center mb-8">
            <div className="relative group cursor-pointer" onClick={() => fileInputRef.current?.click()}>
                <div className="w-24 h-24 rounded-full p-1 bg-gradient-to-tr from-blue-500 to-purple-500">
                     <Avatar size="100%" image={previewImage} className="bg-slate-800" />
                </div>
                
                {/* Upload Overlay */}
                <div className="absolute inset-0 bg-black/50 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                     <UploadIcon className="text-white" />
                </div>
                
                {/* Loading Spinner */}
                {isUploading && (
                     <div className="absolute inset-0 bg-black/60 rounded-full flex items-center justify-center">
                         <svg className="animate-spin h-8 w-8 text-white" viewBox="0 0 24 24"><circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/><path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/></svg>
                     </div>
                )}
            </div>
            <input 
                ref={fileInputRef} 
                type="file" 
                accept="image/*" 
                className="hidden" 
                onChange={handleFileChange} 
            />
            <button 
                onClick={() => fileInputRef.current?.click()}
                className="mt-3 text-sm text-blue-400 hover:text-blue-300 transition-colors"
            >
                {friendAvatar ? '更换头像' : '上传头像 (可选)'}
            </button>
          </div>

          <div className="space-y-6">
            {/* Name Input */}
            <div className="space-y-2">
                <label className="text-sm font-medium text-slate-300 ml-1">对方姓名</label>
                <div className="relative">
                    <div className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-500">
                        <UserIcon />
                    </div>
                    <input
                        type="text"
                        value={friendName}
                        onChange={(e) => setFriendName(e.target.value)}
                        placeholder="例如：李经理、张主管..."
                        maxLength={20}
                        className="w-full bg-slate-900/50 text-white placeholder-slate-500 border border-white/10 rounded-xl py-3 pl-11 pr-4 outline-none focus:border-blue-500/50 focus:ring-1 focus:ring-blue-500/50 transition-all"
                    />
                </div>
            </div>

            {/* Relation Selection */}
            <div className="space-y-2">
                <label className="text-sm font-medium text-slate-300 ml-1">对方关系</label>
                <div className="grid grid-cols-2 gap-4">
                    {['部门同事', '直属上级', '下属', '合作方'].map((relation) => (
                      <button
                        key={relation}
                        onClick={() => setFriendRelation(relation as Relation)}
                        className={`
                            flex items-center justify-center gap-2 py-3 px-4 rounded-xl border transition-all duration-200
                            ${friendRelation === relation
                                ? 'bg-blue-600/20 border-blue-500 text-white shadow-[0_0_15px_rgba(59,130,246,0.2)]'
                                : 'bg-slate-900/50 border-white/10 text-slate-400 hover:bg-slate-800 hover:border-white/20'
                            }
                        `}
                      >
                        <span className="text-sm">{relation}</span>
                      </button>
                    ))}
                </div>
            </div>

            {/* Submit Button */}
            <button
              onClick={handleSubmit}
              disabled={!friendName.trim() || !friendRelation || isSubmitting}
              className={`
                w-full mt-6 py-3.5 rounded-xl font-semibold text-white transition-all duration-200 shadow-lg
                ${!friendName.trim() || !friendRelation || isSubmitting
                    ? 'bg-slate-700 text-slate-500 cursor-not-allowed shadow-none'
                    : 'bg-gradient-to-r from-blue-600 to-indigo-600 hover:shadow-blue-500/30 hover:scale-[1.02] active:scale-[0.98]'
                }
              `}
            >
              {isSubmitting ? '创建中...' : '开始对话'}
            </button>
          </div>
        </div>
      </div>

      <style jsx global>{`
          @keyframes fade-in-up {
              0% { opacity: 0; transform: translateY(20px); }
              100% { opacity: 1; transform: translateY(0); }
          }
          .animate-fade-in-up {
              animation: fade-in-up 0.6s ease-out forwards;
          }
      `}</style>
    </div>
  );
}
