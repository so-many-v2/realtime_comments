// Фронт ходит через nginx относительными путями (/api/...) в тот же origin, что и SPA.
// nginx распроксирует на нужный сервис (см. nginx/nginx.conf). Для dev можно задать
// абсолютный URL через VITE_*-переменные.
const POST_API = import.meta.env.VITE_POST_API || ''
const COMMENT_API = import.meta.env.VITE_COMMENT_API || ''
const SSE_API = import.meta.env.VITE_SSE_API || ''

async function request(url, options) {
  const res = await fetch(url, options)
  if (!res.ok) {
    let detail = ''
    try {
      detail = (await res.json()).error || ''
    } catch (_) {}
    throw new Error(detail || `HTTP ${res.status}`)
  }
  return res.status === 204 ? null : res.json()
}

export const api = {
  // ---- posts (post_service) ----
  listPosts({ limit = 20, offset = 0 } = {}) {
    return request(`${POST_API}/api/posts?limit=${limit}&offset=${offset}`)
  },
  getPost(id) {
    return request(`${POST_API}/api/posts/${id}`)
  },
  createPost({ title, text, user_id }) {
    return request(`${POST_API}/api/posts`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title, text, user_id }),
    })
  },

  // ---- comments (comment_service) ----
  // time_from — unix-секунды; 0 = с начала эпохи (все комментарии).
  getComments({ post_id, time_from = 0, limit = 50 }) {
    return request(`${COMMENT_API}/api/comments?post_id=${post_id}&time_from=${time_from}&limit=${limit}`)
  },
  createComment({ post_id, text, user_id }) {
    return request(`${COMMENT_API}/api/comments`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ post_id, text, user_id }),
    })
  },
}

// SSE-стрим новых комментариев поста (connection_service).
export function commentStreamUrl(postId) {
  return `${SSE_API}/api/sse/${postId}`
}