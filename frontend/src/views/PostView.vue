<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { api } from '../api.js'
import { useCommentStream } from '../composables/useCommentStream.js'

const props = defineProps({ id: { type: [String, Number], required: true } })
const postId = Number(props.id)

const post = ref(null)
const loadingPost = ref(true)
const postError = ref('')

const { comments, totalReceived, connected, seed, connect, close } = useCommentStream(postId, {
  maxRendered: 200,
  maxPerFrame: 120,
})

// --- compose ---
const draft = ref('')
const userId = ref(1)
const sending = ref(false)

// --- автоскролл ---
const scrollEl = ref(null)
const stickToBottom = ref(true)

function onScroll() {
  const el = scrollEl.value
  if (!el) return
  // «прилипли» к низу, если до конца меньше 80px
  stickToBottom.value = el.scrollHeight - el.scrollTop - el.clientHeight < 80
}

function scrollToBottom() {
  const el = scrollEl.value
  if (el) el.scrollTop = el.scrollHeight
}

// при каждом новом батче, если пользователь у низа — догоняем
watch(
  () => comments.value.length,
  async () => {
    if (!stickToBottom.value) return
    await nextTick()
    scrollToBottom()
  }
)

const hiddenCount = computed(() => Math.max(0, totalReceived.value - comments.value.length))

async function loadPost() {
  loadingPost.value = true
  postError.value = ''
  try {
    post.value = await api.getPost(postId)
  } catch (e) {
    postError.value = e.message
  } finally {
    loadingPost.value = false
  }
}

async function loadHistory() {
  try {
    const history = await api.getComments({ post_id: postId, time_from: 0, limit: 50 })
    seed(history || [])
    await nextTick()
    scrollToBottom()
  } catch (_) {
    // история не критична — поток всё равно подключится
  }
}

async function send() {
  const text = draft.value.trim()
  if (!text || sending.value) return
  sending.value = true
  try {
    await api.createComment({ post_id: postId, text, user_id: Number(userId.value) || 1 })
    draft.value = ''
    // сам комментарий прилетит обратно по SSE — локально не добавляем,
    // чтобы не задвоить (дедуп по id это всё равно поймает).
  } catch (_) {
    // оставляем черновик, чтобы можно было повторить
  } finally {
    sending.value = false
  }
}

function fmtTime(ts) {
  const d = new Date(ts)
  return Number.isNaN(d.getTime()) ? '' : d.toLocaleTimeString()
}

onMounted(async () => {
  await loadPost()
  await loadHistory()
  connect()
})
</script>

<template>
  <RouterLink to="/" class="back dim">← к ленте</RouterLink>

  <!-- пост -->
  <article v-if="loadingPost" class="card post">
    <div class="skeleton" style="height: 22px; width: 60%" />
    <div class="skeleton" style="height: 14px; width: 100%; margin-top: 14px" />
    <div class="skeleton" style="height: 14px; width: 90%; margin-top: 8px" />
  </article>

  <div v-else-if="postError" class="card state">
    <p class="err">{{ postError }}</p>
    <button class="btn btn--ghost" @click="loadPost">Повторить</button>
  </div>

  <article v-else class="card post">
    <div class="post__meta">
      <span class="avatar">{{ String(post.author_id).slice(-2) }}</span>
      <span class="dim">автор #{{ post.author_id }}</span>
    </div>
    <h1 class="post__title">{{ post.title }}</h1>
    <p class="post__text">{{ post.text }}</p>
  </article>

  <!-- комментарии -->
  <section class="card comments">
    <header class="comments__head">
      <strong>Комментарии</strong>
      <span class="spacer" />
      <span class="pill" :class="connected ? 'pill--live' : ''">
        <span class="dot" :class="connected ? 'dot--live' : ''" />
        {{ connected ? 'в эфире' : 'подключение…' }}
      </span>
      <span class="pill">{{ totalReceived }} всего</span>
    </header>

    <div v-if="hiddenCount > 0" class="comments__trim dim">
      показаны последние {{ comments.length }} · скрыто ранних: {{ hiddenCount }}
    </div>

    <div ref="scrollEl" class="comments__list" @scroll.passive="onScroll">
      <p v-if="comments.length === 0" class="comments__empty dim">
        Пока пусто. Новые комментарии появятся здесь в реальном времени.
      </p>

      <div v-for="c in comments" :key="c.id" class="comment">
        <span class="avatar avatar--sm">{{ String(c.author_id).slice(-2) }}</span>
        <div class="comment__body">
          <div class="comment__meta dim">
            <span>#{{ c.author_id }}</span>
            <span v-if="fmtTime(c.created)">· {{ fmtTime(c.created) }}</span>
          </div>
          <div class="comment__text">{{ c.text }}</div>
        </div>
      </div>
    </div>

    <button v-if="!stickToBottom" class="jump btn btn--ghost" @click="(stickToBottom = true), scrollToBottom()">
      ↓ к новым
    </button>

    <form class="comments__compose" @submit.prevent="send">
      <input v-model.number="userId" class="input input--mini" type="number" min="1" title="user_id" />
      <input
        v-model="draft"
        class="input"
        placeholder="Написать комментарий…"
        maxlength="2000"
        @keydown.enter.exact.prevent="send"
      />
      <button class="btn btn--primary" :disabled="sending || !draft.trim()">Отправить</button>
    </form>
  </section>
</template>

<style scoped>
.back {
  display: inline-block;
  font-size: 13.5px;
  font-weight: 600;
  margin-bottom: 16px;
}
.post {
  padding: 22px;
  margin-bottom: 16px;
}
.post__meta {
  display: flex;
  align-items: center;
  gap: 9px;
  font-size: 12.5px;
  margin-bottom: 12px;
}
.post__title {
  font-size: 24px;
  font-weight: 800;
  letter-spacing: -0.03em;
  margin: 0 0 10px;
}
.post__text {
  font-size: 15px;
  line-height: 1.6;
  color: var(--ink-2);
  margin: 0;
  white-space: pre-wrap;
}

.comments {
  padding: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  position: relative;
}
.comments__head {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 18px;
  border-bottom: 1px solid var(--line);
}
.comments__trim {
  font-size: 12px;
  padding: 8px 18px;
  background: var(--surface-2);
  border-bottom: 1px solid var(--line);
}
.comments__list {
  height: 52vh;
  min-height: 320px;
  overflow-y: auto;
  padding: 12px 18px;
  scroll-behavior: auto; /* при потоке плавный scroll-behavior только мешает */
}
.comments__empty {
  text-align: center;
  padding: 40px 0;
  font-size: 14px;
}
.comment {
  display: flex;
  gap: 11px;
  padding: 9px 0;
  border-bottom: 1px solid var(--line);
}
.comment:last-child {
  border-bottom: none;
}
.avatar--sm {
  width: 30px;
  height: 30px;
  font-size: 11px;
  flex: 0 0 auto;
}
.comment__body {
  min-width: 0;
}
.comment__meta {
  display: flex;
  gap: 6px;
  font-size: 11.5px;
  margin-bottom: 2px;
}
.comment__text {
  font-size: 14.5px;
  line-height: 1.45;
  word-break: break-word;
}

.jump {
  position: absolute;
  right: 18px;
  bottom: 84px;
  padding: 8px 14px;
  font-size: 13px;
  box-shadow: var(--shadow);
}

.comments__compose {
  display: flex;
  gap: 8px;
  padding: 14px 18px;
  border-top: 1px solid var(--line);
  background: var(--surface-2);
}
.input--mini {
  width: 66px;
  flex: 0 0 auto;
  padding: 10px 8px;
  text-align: center;
}
.state {
  padding: 28px;
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 14px;
  align-items: center;
}
.err {
  color: var(--danger);
  font-size: 13.5px;
  margin: 0;
}
</style>