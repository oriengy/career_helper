'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from 'tdesign-react';
import { AddIcon, ChatIcon } from 'tdesign-icons-react';
import { ROUTES } from '@/constants/routes';

export default function SessionsIndexPage() {
  const router = useRouter();

  const suggestions = [
    { title: '模拟面试', desc: '帮我进行一场模拟前端开发面试' },
    { title: '简历优化', desc: '分析我的简历并给出改进建议' },
    { title: '职业规划', desc: '我现在迷茫，帮我规划职业路径' },
    { title: '职场沟通', desc: '如何委婉地拒绝同事的请求？' },
  ];

  return (
    <div className="flex-1 flex flex-col items-center justify-center p-4 text-center h-full relative overflow-hidden">
        {/* Background blobs for subtle effect */}
        <div className="absolute top-[20%] left-[20%] w-[300px] h-[300px] bg-purple-500/10 rounded-full mix-blend-screen filter blur-[100px] animate-blob pointer-events-none"></div>
        <div className="absolute bottom-[20%] right-[20%] w-[300px] h-[300px] bg-blue-500/10 rounded-full mix-blend-screen filter blur-[100px] animate-blob animation-delay-2000 pointer-events-none"></div>

      <div className="z-10 max-w-2xl w-full animate-fade-in-up">
        <div className="mb-10 flex flex-col items-center">
            <div className="w-20 h-20 bg-gradient-to-tr from-blue-500 to-indigo-600 rounded-2xl flex items-center justify-center mb-6 shadow-lg shadow-blue-500/20 rotate-3 transform hover:rotate-6 transition-transform duration-300">
                <ChatIcon size="40px" className="text-white" />
            </div>
            <h1 className="text-4xl font-bold text-white mb-3 tracking-tight">
                职宝书 AI
            </h1>
            <p className="text-lg text-slate-400">
                您的全天候智能职场助手
            </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 w-full px-4">
          {suggestions.map((item, index) => (
            <button
              key={index}
              onClick={() => router.push(ROUTES.SESSION_NEW)} // In a real app, this might pre-fill the prompt
              className="group text-left p-4 rounded-xl border border-white/10 bg-white/5 hover:bg-white/10 transition-all duration-200 hover:border-blue-500/30 hover:shadow-lg hover:-translate-y-0.5"
            >
              <h3 className="font-medium text-slate-200 mb-1 group-hover:text-blue-300 transition-colors">{item.title}</h3>
              <p className="text-sm text-slate-500 group-hover:text-slate-400 transition-colors">{item.desc}</p>
            </button>
          ))}
        </div>

        <div className="mt-10">
             <Button 
                theme="primary" 
                size="large" 
                icon={<AddIcon />} 
                onClick={() => router.push(ROUTES.SESSION_NEW)}
                className="!bg-gradient-to-r !from-blue-600 !to-indigo-600 !border-none !rounded-full !px-8 !h-12 !text-lg !font-medium hover:!shadow-lg hover:!shadow-blue-500/30 transition-all"
            >
                开始新对话
             </Button>
        </div>
      </div>
      
      <style jsx global>{`
          @keyframes fade-in-up {
              0% { opacity: 0; transform: translateY(20px); }
              100% { opacity: 1; transform: translateY(0); }
          }
          .animate-fade-in-up {
              animation: fade-in-up 0.8s ease-out forwards;
          }
      `}</style>
    </div>
  );
}