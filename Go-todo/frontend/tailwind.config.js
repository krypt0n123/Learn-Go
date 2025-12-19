/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html" ,// 告诉 Tailwind 也要扫描这个文件
    "./src/**/*.ts"
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}