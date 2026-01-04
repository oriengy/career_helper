/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#07c160',
          50: '#e6f9f0',
          100: '#c2f2dc',
          200: '#9aebc7',
          300: '#6ee4b1',
          400: '#45dd9c',
          500: '#07c160',
          600: '#06a855',
          700: '#058f4a',
          800: '#04763e',
          900: '#035d32',
        },
      },
    },
  },
  plugins: [],
}
