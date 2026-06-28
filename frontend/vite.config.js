import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

const POST_SERVICE = process.env.POST_SERVICE_URL || 'http://localhost:8081'
const COMMENT_SERVICE = process.env.COMMENT_SERVICE_URL || 'http://localhost:8082'
const CONNECTION_SERVICE = process.env.CONNECTION_SERVICE_URL || 'http://localhost:8083'

const stripApi = (p) => p.replace(/^\/api/, '')

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      // /api/posts и /api/post/{id} и /api/post -> post_service
      '/api/post': { target: POST_SERVICE, changeOrigin: true, rewrite: stripApi },
      // /api/comments -> comment_service
      '/api/comments': { target: COMMENT_SERVICE, changeOrigin: true, rewrite: stripApi },
      // /api/sse/{post_id} -> connection_service (SSE-стрим)
      '/api/sse': {
        target: CONNECTION_SERVICE,
        changeOrigin: true,
        rewrite: stripApi,
        // важно для SSE: не буферизовать ответ
        configure: (proxy) => {
          proxy.on('proxyReq', (proxyReq) => proxyReq.setHeader('Accept', 'text/event-stream'))
        },
      },
    },
  },
})