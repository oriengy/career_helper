'use client';

import { useState, useEffect, useRef } from 'react';
import { useRouter, useParams } from 'next/navigation';
import {
  Button,
  Textarea,
  Message,
  Avatar,
  Popup,
  Divider,
  ImageViewer,
} from 'tdesign-react';
import {
  ArrowLeftIcon,
  SendIcon,
  ImageIcon,
  ChatIcon,
  TranslateIcon,
} from 'tdesign-icons-react';
import ReactMarkdown from 'react-markdown';
import { useUserStore } from '@/stores/user';
import { sessionApi } from '@/services/api/session';
import { messageApi } from '@/services/api/message';
import { translateApi } from '@/services/api/translate';
import { uploadApi } from '@/services/api/upload';
import { getUserAvatar } from '@/lib/avatar';
import { formatRelativeTime } from '@/services/utils/format';
import Loading from '@/components/common/Loading';
import type { ChatSession, Message as MessageType, MessageRole } from '@/types/models';

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
  const [inputText, setInputText] = useState('');
  const [selectedMessageId, setSelectedMessageId] = useState<string | null>(null);
  const [imageViewerVisible, setImageViewerVisible] = useState(false);
  const [currentImageUrl, setCurrentImageUrl] = useState('');

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // 加载会话信息和消息列表
  const loadData = async () => {
    setIsLoading(true);
    try {
      const [sessionData, messagesData] = await Promise.all([
        sessionApi.listSessions(),
        messageApi.listMessages({ sessionId }),
      ]);

      const currentSession = sessionData.find((s) => s.sessionId === sessionId);
      if (!currentSession) {
        Message.error('会话不存在');
        router.back();
        return;
      }

      setSession(currentSession);
      setMessages(messagesData);
    } catch (error: any) {
      Message.error(error.message || '加载失败');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, [sessionId]);

  // 自动滚动到底部
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  // 发送文本消息
  const handleSendMessage = async () => {
    const text = inputText.trim();
    if (!text) {
      Message.error('请输入消息内容');
      return;
    }

    setIsSending(true);
    try {
      await messageApi.createMessage({
        sessionId,
        role: 'FRIEND',
        content: text,
      });

      setInputText('');
      await loadData();
    } catch (error: any) {
      Message.error(error.message || '发送失败');
    } finally {
      setIsSending(false);
    }
  };

  // 发送图片
  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      const url = await uploadApi.uploadChatImage(file);
      await messageApi.createMessage({
        sessionId,
        role: 'FRIEND',
        imageUrl: url,
      });

      await loadData();
      Message.success('图片发送成功');
    } catch (error: any) {
      Message.error(error.message || '图片上传失败');
    }

    // 重置文件输入
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  // 翻译消息
  const handleTranslate = async (messageId: string) => {
    setIsTranslating(true);
    try {
      await translateApi.translateMessage({ sessionId, messageId });
      await loadData();
      Message.success('翻译成功');
    } catch (error: any) {
      Message.error(error.message || '翻译失败');
    } finally {
      setIsTranslating(false);
      setSelectedMessageId(null);
    }
  };

  // AI 咨询
  const handleConsult = async () => {
    const text = inputText.trim();
    if (!text) {
      Message.error('请输入咨询内容');
      return;
    }

    setIsConsulting(true);
    try {
      await messageApi.sendConsultMessage({
        sessionId,
        content: text,
      });

      setInputText('');
      await loadData();
    } catch (error: any) {
      Message.error(error.message || '咨询失败');
    } finally {
      setIsConsulting(false);
    }
  };

  // 删除消息
  const handleDeleteMessage = async (messageId: string) => {
    try {
      await messageApi.deleteMessage(messageId);
      await loadData();
      Message.success('删除成功');
    } catch (error: any) {
      Message.error(error.message || '删除失败');
    } finally {
      setSelectedMessageId(null);
    }
  };

  // 查看图片
  const handleViewImage = (imageUrl: string) => {
    setCurrentImageUrl(imageUrl);
    setImageViewerVisible(true);
  };

  // 获取消息头像
  const getMessageAvatar = (role: MessageRole) => {
    if (role === 'SELF') {
      return getUserAvatar(profile?.avatar, profile?.gender, profile?.gender, false);
    } else if (role === 'FRIEND') {
      return getUserAvatar(
        session?.friendAvatar,
        session?.friendGender,
        profile?.gender,
        true
      );
    } else if (role === 'AI') {
      return '/ai-avatar.png'; // AI 默认头像
    }
    return '';
  };

  // 获取消息气泡样式
  const getMessageBubbleClass = (role: MessageRole, messageType: string) => {
    if (role === 'SELF') {
      return 'bg-primary-500 text-white';
    } else if (messageType === 'TRANSLATE') {
      return 'bg-yellow-50 border border-yellow-200';
    } else if (messageType === 'CONSULT') {
      return 'bg-blue-50 border border-blue-200';
    }
    return 'bg-white border border-gray-200';
  };

  if (isLoading) {
    return <Loading fullScreen text="加载中..." />;
  }

  if (!session) {
    return null;
  }

  return (
    <div className="h-screen flex flex-col bg-gray-50">
      {/* 头部 */}
      <div className="bg-white border-b border-gray-200 flex-shrink-0">
        <div className="max-w-4xl mx-auto px-4 py-4 flex items-center gap-4">
          <Button
            variant="text"
            shape="circle"
            icon={<ArrowLeftIcon />}
            onClick={() => router.back()}
          />
          <Avatar size="40px" image={getMessageAvatar('FRIEND')} />
          <div className="flex-1">
            <h1 className="text-lg font-semibold text-gray-800">{session.friendName}</h1>
          </div>
        </div>
      </div>

      {/* 消息列表 */}
      <div className="flex-1 overflow-y-auto">
        <div className="max-w-4xl mx-auto px-4 py-6 space-y-4">
          {messages.length === 0 ? (
            <div className="text-center py-12">
              <p className="text-gray-500">暂无消息，开始聊天吧</p>
            </div>
          ) : (
            messages.map((msg) => (
              <div
                key={msg.messageId}
                className={`flex gap-3 ${msg.role === 'SELF' ? 'flex-row-reverse' : ''}`}
              >
                {/* 头像 */}
                <Avatar size="40px" image={getMessageAvatar(msg.role)} />

                {/* 消息内容 */}
                <div
                  className={`flex-1 max-w-[70%] ${msg.role === 'SELF' ? 'items-end' : ''}`}
                >
                  {/* 时间 */}
                  <div
                    className={`text-xs text-gray-500 mb-1 ${
                      msg.role === 'SELF' ? 'text-right' : ''
                    }`}
                  >
                    {formatRelativeTime(msg.createdAt)}
                  </div>

                  {/* 消息气泡 */}
                  <Popup
                    visible={selectedMessageId === msg.messageId}
                    onVisibleChange={(visible) => {
                      if (!visible) setSelectedMessageId(null);
                    }}
                    placement={msg.role === 'SELF' ? 'bottom-end' : 'bottom-start'}
                    content={
                      <div className="p-2 space-y-1">
                        {msg.role === 'FRIEND' && msg.messageType === 'HISTORY' && (
                          <Button
                            variant="text"
                            block
                            onClick={() => handleTranslate(msg.messageId)}
                            disabled={isTranslating}
                            icon={<TranslateIcon />}
                          >
                            翻译此消息
                          </Button>
                        )}
                        <Button
                          variant="text"
                          block
                          onClick={() => handleDeleteMessage(msg.messageId)}
                          theme="danger"
                        >
                          删除
                        </Button>
                      </div>
                    }
                  >
                    <div
                      className={`rounded-lg p-3 inline-block cursor-pointer ${getMessageBubbleClass(
                        msg.role,
                        msg.messageType
                      )}`}
                      onClick={() => setSelectedMessageId(msg.messageId)}
                    >
                      {/* 图片消息 */}
                      {msg.imageUrl && (
                        <img
                          src={msg.imageUrl}
                          alt="消息图片"
                          className="max-w-xs rounded cursor-pointer"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleViewImage(msg.imageUrl!);
                          }}
                        />
                      )}

                      {/* 文本消息 */}
                      {msg.content && (
                        <div
                          className={`markdown ${
                            msg.role === 'SELF' ? 'text-white' : 'text-gray-800'
                          }`}
                        >
                          <ReactMarkdown>{msg.content}</ReactMarkdown>
                        </div>
                      )}

                      {/* 消息类型标签 */}
                      {msg.messageType !== 'HISTORY' && (
                        <div className="text-xs mt-2 opacity-70">
                          {msg.messageType === 'TRANSLATE' && '💬 AI 翻译'}
                          {msg.messageType === 'CONSULT' && '🤖 AI 咨询'}
                        </div>
                      )}
                    </div>
                  </Popup>
                </div>
              </div>
            ))
          )}
          <div ref={messagesEndRef} />
        </div>
      </div>

      {/* 输入区域 */}
      <div className="bg-white border-t border-gray-200 flex-shrink-0">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <div className="flex gap-2 items-end">
            {/* 图片上传 */}
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              className="hidden"
              onChange={handleImageUpload}
            />
            <Button
              variant="outline"
              shape="circle"
              icon={<ImageIcon />}
              onClick={() => fileInputRef.current?.click()}
            />

            {/* 文本输入 */}
            <Textarea
              value={inputText}
              onChange={(value) => setInputText(value)}
              placeholder="输入消息..."
              autosize={{ minRows: 1, maxRows: 4 }}
              className="flex-1"
              onKeyDown={(e) => {
                if (e.key === 'Enter' && !e.shiftKey) {
                  e.preventDefault();
                  handleSendMessage();
                }
              }}
            />

            {/* AI 咨询按钮 */}
            <Button
              variant="outline"
              icon={<ChatIcon />}
              onClick={handleConsult}
              loading={isConsulting}
              disabled={!inputText.trim()}
            >
              咨询
            </Button>

            {/* 发送按钮 */}
            <Button
              theme="primary"
              icon={<SendIcon />}
              onClick={handleSendMessage}
              loading={isSending}
              disabled={!inputText.trim()}
            >
              发送
            </Button>
          </div>

          <div className="mt-2 text-xs text-gray-500">
            按 Enter 发送消息，Shift + Enter 换行
          </div>
        </div>
      </div>

      {/* 图片查看器 */}
      <ImageViewer
        images={[currentImageUrl]}
        visible={imageViewerVisible}
        onClose={() => setImageViewerVisible(false)}
      />
    </div>
  );
}
