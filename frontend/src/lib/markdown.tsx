/**
 * Markdown 渲染配置
 */

import React from 'react';

/**
 * Markdown 组件配置
 */
export const markdownComponents: any = {
  // 代码块
  code({ node, inline, className, children, ...props }: any) {
    const match = /language-(\w+)/.exec(className || '');

    if (!inline && match) {
      // 代码块 - 简单版本，不使用 Prism 以避免依赖问题
      return (
        <pre className={className} {...props}>
          <code>{children}</code>
        </pre>
      );
    }

    // 行内代码
    return (
      <code className={className} {...props}>
        {children}
      </code>
    );
  },

  // 链接（新窗口打开）
  a({ node, children, href, ...props }: any) {
    return (
      <a href={href} target="_blank" rel="noopener noreferrer" {...props}>
        {children}
      </a>
    );
  },

  // 图片（懒加载）
  img({ node, src, alt, ...props }: any) {
    return (
      <img
        src={src}
        alt={alt}
        loading="lazy"
        style={{ maxWidth: '100%', height: 'auto' }}
        {...props}
      />
    );
  },
};

/**
 * Markdown 渲染选项
 */
export const markdownOptions = {
  components: markdownComponents,
};
