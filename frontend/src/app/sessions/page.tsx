'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Button, Message, Dialog, Empty, Avatar } from 'tdesign-react';
import { AddIcon, DeleteIcon } from 'tdesign-icons-react';
import { useUserStore } from '@/stores/user';
import { sessionApi } from '@/services/api/session';
import { getUserAvatar } from '@/lib/avatar';
import { formatRelativeTime } from '@/services/utils/format';
import { ROUTES } from '@/constants/routes';
import Loading from '@/components/common/Loading';
import type { ChatSession } from '@/types/models';

export default function SessionsPage() {
  const router = useRouter();
  const { profile } = useUserStore();
  const [sessions, setSessions] = useState<ChatSession[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreating, setIsCreating] = useState(false);

  // 加载会话列表
  const loadSessions = async () => {
    setIsLoading(true);
    try {
      const data = await sessionApi.listSessions();
      setSessions(data);
    } catch (error: any) {
      Message.error(error.message || '加载会话列表失败');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadSessions();
  }, []);

  // 创建新会话
  const handleCreateSession = async () => {
    router.push(ROUTES.SESSION_NEW);
  };

  // 删除会话
  const handleDeleteSession = async (sessionId: string, friendName: string) => {
    const confirmResult = await Dialog.confirm({
      header: '确认删除',
      body: `确定要删除与 ${friendName} 的聊天记录吗？此操作不可恢复。`,
      confirmBtn: '删除',
      cancelBtn: '取消',
      theme: 'warning',
    });

    if (confirmResult) {
      try {
        await sessionApi.deleteSession(sessionId);
        Message.success('删除成功');
        // 刷新列表
        loadSessions();
      } catch (error: any) {
        Message.error(error.message || '删除失败');
      }
    }
  };

  // 进入会话详情
  const handleSessionClick = (sessionId: string) => {
    router.push(ROUTES.SESSION_DETAIL(sessionId));
  };

  if (isLoading) {
    return <Loading fullScreen text="加载中..." />;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 头部 */}
      <div className="bg-white border-b border-gray-200 sticky top-0 z-10">
        <div className="max-w-4xl mx-auto px-4 py-4 flex items-center justify-between">
          <h1 className="text-xl font-semibold text-gray-800">聊天列表</h1>
          <Button
            icon={<AddIcon />}
            theme="primary"
            onClick={handleCreateSession}
            loading={isCreating}
          >
            新建会话
          </Button>
        </div>
      </div>

      {/* 会话列表 */}
      <div className="max-w-4xl mx-auto px-4 py-6">
        {sessions.length === 0 ? (
          <div className="bg-white rounded-lg shadow-sm p-12">
            <Empty description="暂无聊天记录">
              <Button theme="primary" onClick={handleCreateSession} className="mt-4">
                创建第一个会话
              </Button>
            </Empty>
          </div>
        ) : (
          <div className="space-y-2">
            {sessions.map((session) => (
              <div
                key={session.sessionId}
                className="bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow cursor-pointer"
              >
                <div className="p-4 flex items-center gap-4">
                  {/* 头像 */}
                  <div
                    onClick={() => handleSessionClick(session.sessionId)}
                    className="flex-shrink-0"
                  >
                    <Avatar
                      size="large"
                      image={getUserAvatar(
                        session.friendAvatar,
                        session.friendGender,
                        profile?.gender,
                        true
                      )}
                    />
                  </div>

                  {/* 会话信息 */}
                  <div
                    onClick={() => handleSessionClick(session.sessionId)}
                    className="flex-1 min-w-0"
                  >
                    <div className="flex items-center justify-between mb-1">
                      <h3 className="text-base font-semibold text-gray-800 truncate">
                        {session.friendName}
                      </h3>
                      <span className="text-xs text-gray-500 flex-shrink-0 ml-2">
                        {formatRelativeTime(session.updatedAt)}
                      </span>
                    </div>
                    <p className="text-sm text-gray-500 truncate">
                      {session.lastMessage || '暂无消息'}
                    </p>
                  </div>

                  {/* 删除按钮 */}
                  <div className="flex-shrink-0">
                    <Button
                      variant="text"
                      shape="circle"
                      icon={<DeleteIcon />}
                      onClick={(e) => {
                        e.stopPropagation();
                        handleDeleteSession(session.sessionId, session.friendName);
                      }}
                      theme="danger"
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
