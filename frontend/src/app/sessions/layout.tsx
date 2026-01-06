'use client';

import { useState, useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { MessagePlugin, DialogPlugin, Avatar } from 'tdesign-react';
import { AddIcon, DeleteIcon } from 'tdesign-icons-react';
import { useUserStore } from '@/stores/user';
import { sessionApi } from '@/services/api/session';
import { getUserAvatar } from '@/lib/avatar';
import { ROUTES } from '@/constants/routes';
import type { ChatSession } from '@/types/models';

export default function SessionsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const { profile } = useUserStore();
  const [sessions, setSessions] = useState<ChatSession[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const loadSessions = async () => {
    try {
      const data = await sessionApi.listSessions();
      setSessions(data);
    } catch (error: any) {
      console.error('Failed to load sessions:', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadSessions();
  }, []);

  useEffect(() => {
    if (pathname.startsWith(ROUTES.SESSIONS)) {
      loadSessions();
    }
  }, [pathname]);

  const handleCreateSession = () => {
    router.push(ROUTES.SESSION_NEW);
    if (window.innerWidth < 1024) setIsMobileMenuOpen(false);
  };

  const handleSessionClick = (sessionId: string) => {
    router.push(ROUTES.SESSION_DETAIL(sessionId));
    if (window.innerWidth < 1024) setIsMobileMenuOpen(false);
  };

  const handleDeleteSession = (e: React.MouseEvent, sessionId: string, friendName: string) => {
    e.stopPropagation();
    const dialog = DialogPlugin.confirm({
      header: 'Confirm delete',
      body: `Delete conversation with ${friendName}?`,
      confirmBtn: 'Delete',
      cancelBtn: 'Cancel',
      theme: 'warning',
      onConfirm: async () => {
        try {
          await sessionApi.deleteSession(sessionId);
          MessagePlugin.success('Deleted.');
          loadSessions();
          if (pathname.includes(sessionId)) {
            router.push(ROUTES.SESSIONS);
          }
        } catch (error: any) {
          MessagePlugin.error('Delete failed.');
        } finally {
          dialog.destroy();
        }
      },
      onClose: () => {
        dialog.destroy();
      },
    });
  };

  return (
    <div className="flex h-screen bg-slate-900 text-slate-100 overflow-hidden font-sans">
      {/* Mobile Header */}
      <div className="lg:hidden fixed top-0 left-0 w-full bg-slate-900 border-b border-slate-800 z-50 flex items-center justify-between p-4">
        <button onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)} className="p-2 text-slate-400 hover:text-white">
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" /></svg>
        </button>
        <span className="font-semibold text-slate-200">Career Helper AI</span>
        <button onClick={handleCreateSession} className="p-2 text-slate-400 hover:text-white">
          <AddIcon />
        </button>
      </div>

      {/* Sidebar Overlay for Mobile */}
      {isMobileMenuOpen && (
        <div className="fixed inset-0 bg-black/50 z-40 lg:hidden" onClick={() => setIsMobileMenuOpen(false)} />
      )}

      {/* Sidebar */}
      <div className={`
        fixed lg:static inset-y-0 left-0 z-50 w-[260px] bg-black bg-opacity-95 flex flex-col transition-transform duration-300 ease-in-out border-r border-white/10
        ${isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
      `}>
        {/* New Chat Button */}
        <div className="p-3 mb-2">
          <button
            onClick={handleCreateSession}
            className="w-full flex items-center gap-3 px-4 py-3 text-sm text-white bg-slate-800/50 hover:bg-slate-800 border border-white/10 rounded-lg transition-colors duration-200 group"
          >
            <div className="p-1 bg-white/10 rounded-full group-hover:bg-white/20 transition-colors">
              <AddIcon size="16px" />
            </div>
            <span>New chat</span>
          </button>
        </div>

        {/* Session List */}
        <div className="flex-1 overflow-y-auto custom-scrollbar px-3 space-y-2">
          <div className="text-xs font-semibold text-slate-500 px-4 py-2">Recent</div>
          {isLoading ? (
            <div className="px-4 py-2 text-sm text-slate-500 animate-pulse">Loading...</div>
          ) : sessions.length === 0 ? (
            <div className="px-4 py-10 text-center">
              <p className="text-sm text-slate-600">No history yet</p>
            </div>
          ) : (
            sessions.map((session) => (
              <div
                key={session.sessionId}
                onClick={() => handleSessionClick(session.sessionId)}
                className={`
                  group relative flex items-center gap-3 px-3 py-3 rounded-lg cursor-pointer transition-colors duration-200
                  ${pathname.includes(session.sessionId) ? 'bg-slate-800 text-white' : 'text-slate-300 hover:bg-slate-900'}
                `}
              >
                <div className="flex-shrink-0">
                  <Avatar
                    size="small"
                    shape="circle"
                    image={getUserAvatar(session.friendAvatar, session.friendGender, profile?.gender, true)}
                    className="opacity-80 group-hover:opacity-100 transition-opacity"
                  />
                </div>
                <div className="flex-1 min-w-0 flex flex-col">
                  <span className="text-sm truncate font-medium">{session.friendName || 'Unknown'}</span>
                  <span className="text-xs text-slate-500 truncate">{session.lastMessage || 'Click to view'}</span>
                </div>

                {/* Delete Button (Visible on Hover) */}
                <button
                  onClick={(e) => handleDeleteSession(e, session.sessionId, session.friendName)}
                  className="absolute right-2 opacity-0 group-hover:opacity-100 p-1.5 text-slate-400 hover:text-red-400 hover:bg-slate-700 rounded transition-all"
                >
                  <DeleteIcon size="14px" />
                </button>
              </div>
            ))
          )}
        </div>

        {/* User Profile Area */}
        <div className="p-4 border-t border-white/10">
          <div className="flex items-center gap-3 px-2 py-2 rounded-lg hover:bg-slate-800 cursor-pointer transition-colors">
            <Avatar
              size="small"
              image={getUserAvatar(profile?.avatar, profile?.gender, profile?.gender, false)}
            />
            <div className="flex-1 min-w-0">
              <div className="text-sm font-medium text-white truncate">{profile?.nickname || 'My account'}</div>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="flex-1 relative flex flex-col h-full overflow-hidden bg-slate-900">
        {children}
      </div>

      <style jsx global>{`
        .custom-scrollbar::-webkit-scrollbar {
          width: 6px;
        }
        .custom-scrollbar::-webkit-scrollbar-track {
          background: transparent;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb {
          background-color: rgba(255, 255, 255, 0.1);
          border-radius: 20px;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb:hover {
          background-color: rgba(255, 255, 255, 0.2);
        }
      `}</style>
    </div>
  );
}
