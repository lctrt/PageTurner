import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': '/src',
    },
  },
  server: {
    host: '0.0.0.0',
    allowedHosts: process.env.ALLOWED_HOSTS.split(','),
    port: 5173,
    watch: {
      usePolling: true,
      interval: 100,
    },
  },
})
