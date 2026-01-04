'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { ROUTES } from '@/constants/routes';

export default function HomePage() {
  const router = useRouter();

  useEffect(() => {
    // 重定向到会话列表页
    router.replace(ROUTES.SESSIONS);
  }, [router]);

  return <div className="flex items-center justify-center h-screen">加载中...</div>;
}
