'use client';

import { useState, useEffect, useRef } from 'react';
import { useRouter, useParams } from 'next/navigation';
import {
  Button,
  MessagePlugin,
  Avatar,
  ImageViewer,
} from 'tdesign-react';
import {
  SendIcon,
  ImageIcon,
  ChatIcon,
  TranslateIcon,
  DeleteIcon,
  ArrowLeftIcon
} from 'tdesign-icons-react';
import ReactMarkdown from 'react-markdown';
import { useUserStore } from '@/stores/user';
import { sessionApi } from '@/services/api/session';
import { messageApi } from '@/services/api/message';
import { translateApi } from '@/services/api/translate';
import { uploadApi } from '@/services/api/upload';
import { getUserAvatar } from '@/lib/avatar';
import { formatRelativeTime } from '@/services/utils/format';
import { ROUTES } from '@/constants/routes';
import Loading from '@/components/common/Loading';
import type { ChatSession, Message as MessageType, MessageRole } from '@/types/models';

// Icons for message actions
const CopyIcon = () => <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>;

export default function SessionDetailPage() {
  const router = useRouter();
  const params = useParams();
  const sessionId = params.id as string;
  const { profile } = useUserStore();

  const [session, setSession] = useState<ChatSession | null>(null);
  const [messages, setMessages] = useState<MessageType[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSending, setIsSending] = useState(false);
  const [isTranslating, setIsTranslating] = useState(false);
  const [isConsulting, setIsConsulting] = useState(false);
  const [translatingMessageId, setTranslatingMessageId] = useState<string | null>(null);
  
  // Input states
  const [inputText, setInputText] = useState('');
  const [pendingFile, setPendingFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string>('');

  const [selectedMessageId, setSelectedMessageId] = useState<string | null>(null);
  const [imageViewerVisible, setImageViewerVisible] = useState(false);
  const [currentImageUrl, setCurrentImageUrl] = useState('');

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const textAreaRef = useRef<HTMLTextAreaElement>(null);

  // Load Data
  const loadData = async () => {
    if (messages.length === 0) setIsLoading(true);
    try {
      const [sessionData, messagesData] = await Promise.all([
        sessionApi.listSessions(),
        messageApi.listMessages({ sessionId }),
      ]);

      const currentSession = sessionData.find((s) => s.sessionId === sessionId);
      if (!currentSession) {
        MessagePlugin.error('会话不存在');
        router.push(ROUTES.SESSIONS);
        return;
      }

      setSession(currentSession);
      setMessages(messagesData);
    } catch (error: any) {
      MessagePlugin.error(error.message || '加载失败');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadData();
    return () => {
      if (previewUrl) URL.revokeObjectURL(previewUrl);
    };
  }, [sessionId]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, pendingFile]);

  // Handle Send
  const handleSendMessage = async () => {
    const text = inputText.trim();
    if (!text && !pendingFile) {
      MessagePlugin.error('Please enter a message or select an image.');
      return;
    }

    setIsSending(true);

    try {
      // 1. Upload and parse image if exists
      if (pendingFile) {
        try {
            const uploadResult = await uploadApi.uploadChatImage(pendingFile);
            let imageKey = uploadResult.url || uploadResult.publicUrl;
            if (imageKey && imageKey.startsWith('http')) {
              try {
                imageKey = decodeURIComponent(new URL(imageKey).pathname).replace(/^\/+/, '');
              } catch {}
            }
            if (!imageKey) {
              throw new Error('Missing image key');
            }
            await messageApi.parseImageMessages(sessionId, imageKey);
        } catch (error: any) {
            console.error('Image upload failed', error);
            MessagePlugin.error('Image send failed.');
        }
      }

      // 2. Send consult if exists
      if (text) {
          await messageApi.sendConsultMessage({ sessionId, content: text });
      }

      // Reset states
      setInputText('');
      setPendingFile(null);
      setPreviewUrl('');
      if (textAreaRef.current) {
          textAreaRef.current.style.height = 'auto';
          textAreaRef.current.style.minHeight = '40px';
      }

      await loadData();
    } catch (error: any) {
      MessagePlugin.error(error.message || 'Send failed.');
    } finally {
      setIsSending(false);
    }
  };

  // Handle Image Selection (Button)
  const handleImageSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setFilePreview(file);
    if (fileInputRef.current) fileInputRef.current.value = '';
  };

  // Handle Paste
  const handlePaste = (e: React.ClipboardEvent) => {
    const items = e.clipboardData.items;
    for (let i = 0; i < items.length; i++) {
      if (items[i].type.indexOf('image') !== -1) {
        e.preventDefault();
        const file = items[i].getAsFile();
        if (file) setFilePreview(file);
        return;
      }
    }
  };

  const setFilePreview = (file: File) => {
      if (previewUrl) URL.revokeObjectURL(previewUrl);
      const url = URL.createObjectURL(file);
      setPendingFile(file);
      setPreviewUrl(url);
  };

  const clearPreview = () => {
      setPendingFile(null);
      setPreviewUrl('');
      if (fileInputRef.current) fileInputRef.current.value = '';
  };

  // Handle Translate
  const handleTranslate = async (messageId: string) => {
    setIsTranslating(true);
    setTranslatingMessageId(messageId);
    try {
      await translateApi.translateMessage({
        chatSessionId: sessionId,
        targetMessageId: messageId,
      });
      await loadData();
      MessagePlugin.success('????');
    } catch (error: any) {
      MessagePlugin.error(error.message || '????');
    } finally {
      setIsTranslating(false);
      setTranslatingMessageId(null);
      setSelectedMessageId(null);
    }
  };

  // Handle Consult
  const handleConsult = async () => {
    const text = inputText.trim();
    if (!text) {
        MessagePlugin.error('请输入咨询内容');
        return;
    }
    setIsConsulting(true);
    try {
      await messageApi.sendConsultMessage({ sessionId, content: text });
      setInputText('');
      await loadData();
    } catch (error: any) {
      MessagePlugin.error(error.message || '咨询失败');
    } finally {
      setIsConsulting(false);
    }
  };

  // Handle Delete
  const handleDeleteMessage = async (messageId: string) => {
    try {
      await messageApi.deleteMessage(messageId);
      setMessages(prev => prev.filter(m => m.messageId !== messageId));
      MessagePlugin.success('删除成功');
    } catch (error: any) {
      MessagePlugin.error(error.message || '删除失败');
    } finally {
      setSelectedMessageId(null);
    }
  };

  const handleCopy = (content: string) => {
      navigator.clipboard.writeText(content);
      MessagePlugin.success('已复制');
      setSelectedMessageId(null);
  };

  // Avatars
  const getMessageAvatar = (role: MessageRole) => {
    if (role === 'SELF') {
      return getUserAvatar(profile?.avatar, profile?.gender, profile?.gender, false);
    } else if (role === 'FRIEND') {
      return getUserAvatar(session?.friendAvatar, session?.friendGender, profile?.gender, true);
    } else if (role === 'AI') {
      return 'https://tdesign.gtimg.com/site/avatar.jpg'; // Placeholder for AI
    }
    return '';
  };

  if (isLoading && messages.length === 0) {
    return (
        <div className="h-full flex items-center justify-center bg-slate-900">
            <Loading text="正在连接..." />
        </div>
    );
  }

  if (!session) return null;

  const historyMessages = messages.filter((msg) => msg.msgType !== 'CONSULT');
  const consultMessages = messages.filter((msg) => msg.msgType === 'CONSULT');

  return (
    <div className="h-full flex flex-col bg-slate-900 relative">
      {/* Mobile Header - Only visible on small screens */}
      <div className="lg:hidden absolute top-0 left-0 right-0 z-20 bg-slate-900/80 backdrop-blur-md border-b border-white/5 px-4 py-3 flex items-center gap-3">
          <button onClick={() => router.push(ROUTES.SESSIONS)} className="text-slate-400 hover:text-white">
            <ArrowLeftIcon />
          </button>
          <span className="text-white font-medium">{session.friendName}</span>
      </div>

      {/* Messages Area */}
      <div className="flex-1 overflow-y-auto custom-scrollbar pt-16 lg:pt-4 pb-4 px-4 scroll-smooth">
        <div className="max-w-3xl mx-auto space-y-6">
          {historyMessages.length === 0 && consultMessages.length === 0 ? (
            <div className="text-center py-20 opacity-50 animate-fade-in-up">
                <div className="w-16 h-16 bg-white/5 rounded-full mx-auto flex items-center justify-center mb-4">
                     <ChatIcon size="32px" className="text-slate-400" />
                </div>
                <h2 className="text-xl font-medium text-slate-200 mb-2">开始与 {session.friendName} 的对话</h2>
                <p className="text-sm text-slate-500">这里是你们的专属空间</p>
            </div>
          ) : (
            historyMessages.map((msg, index) => {
               const isUser = msg.role === 'SELF' || msg.role === 'USER';
               return (
                <div
                    key={msg.messageId || index}
                    className={`flex gap-4 animate-message-slide-in ${isUser ? 'flex-row-reverse' : 'flex-row'}`}
                    style={{ animationDelay: `${index * 50}ms` }}
                >
                    {/* Avatar */}
                    <div className="flex-shrink-0 mt-1">
                        <Avatar 
                            size="40px" 
                            shape="circle"
                            image={getMessageAvatar(msg.role)} 
                            className="border border-white/10"
                        />
                    </div>

                    {/* Content */}
                    <div className={`flex flex-col max-w-[85%] lg:max-w-[75%] ${isUser ? 'items-end' : 'items-start'}`}>
                        <div className="flex items-center gap-2 mb-1 px-1">
                             <span className="text-xs text-slate-500 font-medium">
                                {isUser ? 'You' : (msg.role === 'AI' ? 'AI Assistant' : session.friendName)}
                             </span>
                             <span className="text-xs text-slate-600">
                                {formatRelativeTime(msg.createdAt)}
                             </span>
                        </div>

                        <div className="relative group">
                            {/* Message Bubble */}
                            <div 
                                className={`
                                    rounded-2xl px-4 py-3 text-sm leading-relaxed shadow-sm relative z-10
                                    ${isUser 
                                        ? 'bg-blue-600 text-white rounded-tr-sm' 
                                        : 'bg-slate-800 text-slate-200 border border-white/5 rounded-tl-sm'
                                    }
                                    ${msg.msgType === 'TRANSLATE' ? '!bg-amber-900/30 !border-amber-700/50 !text-amber-100' : ''}
                                    ${msg.msgType === 'CONSULT' ? '!bg-indigo-900/30 !border-indigo-700/50 !text-indigo-100' : ''}
                                `}
                            >
                                {msg.imageUrl ? (
                                    <img
                                        src={msg.imageUrl}
                                        alt="Uploaded"
                                        className="max-w-full rounded-lg cursor-zoom-in hover:opacity-90 transition-opacity"
                                        style={{ maxHeight: '300px' }}
                                        onClick={() => { setCurrentImageUrl(msg.imageUrl!); setImageViewerVisible(true); }}
                                    />
                                ) : (
                                    <div className="markdown prose prose-invert prose-sm max-w-none">
                                        <ReactMarkdown>{msg.content || ''}</ReactMarkdown>
                                    </div>
                                )}
                            </div>

                            {msg.translateContent && (
                                <div className={`mt-2 text-xs rounded-lg border px-3 py-2 ${isUser ? 'bg-slate-700/50 border-white/10 text-slate-200' : 'bg-amber-900/20 border-amber-700/40 text-amber-100'}`}>
                                  <div className="markdown prose prose-invert prose-sm max-w-none">
                                    <ReactMarkdown>{msg.translateContent}</ReactMarkdown>
                                  </div>
                                </div>
                            )}

                            {/* Action Tools (Visible on Hover) */}
                            <div className={`
                                absolute top-0 ${isUser ? '-left-24' : '-right-24'} h-full flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity px-2
                            `}>
                                {msg.content && (
                                    <button 
                                        onClick={() => handleCopy(msg.content!)}
                                        className="p-1.5 text-slate-400 hover:text-white hover:bg-white/10 rounded" 
                                        title="复制"
                                    >
                                        <CopyIcon />
                                    </button>
                                )}
                                <button 
                                    onClick={() => handleDeleteMessage(msg.messageId)}
                                    className="p-1.5 text-slate-400 hover:text-red-400 hover:bg-white/10 rounded"
                                    title="删除"
                                >
                                    <DeleteIcon size="14px" />
                                </button>
                                {msg.role === 'FRIEND' && msg.msgType === 'HISTORY' && (
                                     <button 
                                        onClick={() => handleTranslate(msg.messageId)}
                                        className={`p-1.5 text-slate-400 hover:text-amber-400 hover:bg-white/10 rounded ${translatingMessageId === msg.messageId ? 'animate-pulse' : ''}`}
                                        title="翻译"
                                     >
                                        <TranslateIcon size="14px" />
                                     </button>
                                )}
                            </div>
                        </div>
                        
                        {/* Type Labels */}
                        {msg.msgType !== 'HISTORY' && msg.msgType !== 'TEXT' && (
                            <span className="text-[10px] mt-1 text-slate-500 uppercase tracking-wider px-1">
                                {msg.msgType === 'TRANSLATE' && 'AI 翻译'}
                                {msg.msgType === 'CONSULT' && 'AI 咨询'}
                            </span>
                        )}
                    </div>
                </div>
               );
            })
          )}
          {consultMessages.length > 0 && (
            <div className="pt-4 border-t border-white/5 space-y-3">
              {consultMessages.map((msg, index) => {
                const isUser = msg.role === 'USER' || msg.role === 'SELF';
                return (
                  <div
                    key={msg.messageId || `consult-${index}`}
                    className={`flex ${isUser ? 'justify-end' : 'justify-start'}`}
                  >
                    <div
                      className={`
                        max-w-[85%] lg:max-w-[75%] text-sm leading-relaxed
                        ${isUser ? 'bg-slate-200 text-slate-900' : 'text-slate-200'}
                        ${isUser ? 'rounded-2xl rounded-tr-sm px-4 py-2' : 'px-1'}
                      `}
                    >
                      <div className="markdown prose prose-invert prose-sm max-w-none">
                        <ReactMarkdown>{msg.content || ''}</ReactMarkdown>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
          <div ref={messagesEndRef} className="h-4" />
        </div>
      </div>

      {/* Input Area */}
      <div className="flex-shrink-0 bg-gradient-to-t from-slate-900 via-slate-900 to-transparent pt-10 pb-6">
        <div className="max-w-3xl mx-auto px-4">
             {/* Action Buttons Floating above input */}
             <div className="flex justify-center gap-3 mb-3">
                 {inputText.trim().length > 0 && (
                     <button
                        onClick={handleConsult}
                        disabled={isConsulting}
                        className={`
                            flex items-center gap-2 px-4 py-2 rounded-full text-xs font-medium backdrop-blur-md border transition-all
                            ${isConsulting 
                                ? 'bg-indigo-500/20 border-indigo-500/50 text-indigo-300' 
                                : 'bg-slate-800/80 border-white/10 text-slate-300 hover:bg-indigo-600 hover:text-white hover:border-indigo-500'
                            }
                        `}
                     >
                        {isConsulting ? (
                            <span className="animate-spin mr-1">⟳</span>
                        ) : (
                            <ChatIcon />
                        )}
                        AI 深度咨询
                     </button>
                 )}
             </div>

             {/* Main Input Box */}
             <div className="relative bg-slate-800 rounded-2xl border border-white/10 shadow-lg focus-within:border-white/20 focus-within:ring-1 focus-within:ring-white/20 transition-all overflow-hidden">
                <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/*"
                    className="hidden"
                    onChange={handleImageSelect}
                />

                {/* Image Preview Area */}
                {pendingFile && (
                    <div className="px-4 pt-4 pb-2">
                        <div className="relative inline-block group">
                            <img 
                                src={previewUrl} 
                                alt="Preview" 
                                className="h-20 w-auto rounded-lg border border-white/10 object-cover"
                            />
                            <button 
                                onClick={clearPreview}
                                className="absolute -top-2 -right-2 bg-slate-700 text-white rounded-full p-1 hover:bg-red-500 transition-colors shadow-md"
                            >
                                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                            </button>
                        </div>
                    </div>
                )}
                
                <div className="flex items-end p-3 gap-2">
                     <button 
                        onClick={() => fileInputRef.current?.click()}
                        className="p-2 text-slate-400 hover:text-white hover:bg-white/10 rounded-xl transition-colors mb-0.5"
                        title="上传图片"
                     >
                        <ImageIcon size="20px" />
                     </button>
                     
                     <textarea
                        ref={textAreaRef}
                        rows={1}
                        value={inputText}
                        onChange={(e) => {
                            setInputText(e.target.value);
                            e.target.style.height = 'auto';
                            e.target.style.height = `${Math.min(e.target.scrollHeight, 120)}px`;
                        }}
                        onPaste={handlePaste}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter' && !e.shiftKey) {
                                e.preventDefault();
                                handleSendMessage();
                            }
                        }}
                        placeholder="输入消息，可粘贴图片..."
                        className="flex-1 bg-transparent text-white placeholder-slate-500 resize-none outline-none max-h-[120px] py-2 text-sm leading-6 custom-scrollbar"
                        style={{ minHeight: '40px' }}
                     />
                     
                     <button
                        onClick={handleSendMessage}
                        disabled={(!inputText.trim() && !pendingFile) || isSending}
                        className={`
                            p-2 rounded-xl transition-all duration-200 mb-0.5
                            ${(inputText.trim() || pendingFile)
                                ? 'bg-blue-600 text-white hover:bg-blue-500 shadow-lg shadow-blue-900/50' 
                                : 'bg-slate-700 text-slate-500 cursor-not-allowed'
                            }
                        `}
                     >
                        {isSending ? (
                             <svg className="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                             </svg>
                        ) : (
                            <SendIcon size="20px" />
                        )}
                     </button>
                </div>
             </div>
             
             <p className="text-center text-xs text-slate-600 mt-3">
                 职宝书 AI 可能会产生错误信息。请核对重要信息。
             </p>
        </div>
      </div>

      <ImageViewer
        images={currentImageUrl ? [currentImageUrl] : []}
        visible={imageViewerVisible && Boolean(currentImageUrl)}
        trigger={<span className="hidden" />}
        closeOnOverlay
        closeOnEscKeydown
        closeBtn
        onClose={() => {
          setImageViewerVisible(false);
          setCurrentImageUrl('');
        }}
        index={0}
      />

      <style jsx global>{`
          @keyframes message-slide-in {
              from { opacity: 0; transform: translateY(10px); }
              to { opacity: 1; transform: translateY(0); }
          }
          .animate-message-slide-in {
              animation: message-slide-in 0.4s cubic-bezier(0.2, 0.8, 0.2, 1) forwards;
          }
      `}</style>
    </div>
  );
}
