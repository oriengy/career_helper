import type { Metadata } from 'next';
import './globals.css';
import 'tdesign-react/es/style/index.css';

export const metadata: Metadata = {
  title: '恋爱翻译官 - AI 驱动的社交辅助应用',
  description: 'AI 驱动的社交辅助应用，帮助你更好地理解异性交流',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <body className="antialiased">{children}</body>
    </html>
  );
}
