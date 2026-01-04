/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**',
      },
    ],
  },
  // 环境变量
  env: {
    NEXT_PUBLIC_APP_NAME: '恋爱翻译官',
  },
}

module.exports = nextConfig
