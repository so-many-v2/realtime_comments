<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api.js'

const router = useRouter()

const posts = ref([])
const loading = ref(true)
const error = ref('')

const showForm = ref(false)
const form = ref({ title: '', text: '', user_id: 1 })
const submitting = ref(false)
const formError = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    posts.value = (await api.listPosts({ limit: 30 })) || []
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

async function submit() {
  if (!form.value.title.trim() || !form.value.text.trim()) return
  submitting.value = true
  formError.value = ''
  try {
    const { post_id } = await api.createPost({
      title: form.value.title.trim(),
      text: form.value.text.trim(),
      user_id: Number(form.value.user_id) || 1,
    })
    router.push({ name: 'post', params: { id: post_id } })
  } catch (e) {
    formError.value = e.message
  } finally {
    submitting.value = false
  }
}

function excerpt(text, n = 160) {
  return text.length > n ? text.slice(0, n).trimEnd() + '…' : text
}

onMounted(load)
</script>

<template>
  <div class="feed-head">
    <div>
      <h1 class="feed-title">Лента</h1>
      <p class="muted feed-sub">Все посты, новые сверху</p>
    </div>
    <button class="btn btn--primary" @click="showForm = !showForm">
      {{ showForm ? 'Закрыть' : 'Новый пост' }}
    </button>
  </div>

  <!-- форма создания поста -->
  <form v-if="showForm" class="card compose" @submit.prevent="submit">
    <input v-model="form.title" class="input" placeholder="Заголовок" maxlength="256" />
    <textarea v-model="form.text" class="textarea" rows="3" placeholder="Текст поста" />
    <div class="compose__row">
      <label class="dim compose__user">
        user_id
        <input v-model="form.user_id" class="input input--mini" type="number" min="1" />
      </label>
      <span class="spacer" />
      <button class="btn btn--primary" :disabled="submitting">
        {{ submitting ? 'Публикуем…' : 'Опубликовать' }}
      </button>
    </div>
    <p v-if="formError" class="err">{{ formError }}</p>
  </form>

  <!-- состояния -->
  <div v-if="loading" class="list">
    <div v-for="i in 4" :key="i" class="card post-card">
      <div class="skeleton" style="height: 18px; width: 55%" />
      <div class="skeleton" style="height: 13px; width: 100%; margin-top: 12px" />
      <div class="skeleton" style="height: 13px; width: 80%; margin-top: 8px" />
    </div>
  </div>

  <div v-else-if="error" class="card state">
    <p class="err">Не удалось загрузить ленту: {{ error }}</p>
    <button class="btn btn--ghost" @click="load">Повторить</button>
  </div>

  <div v-else-if="posts.length === 0" class="card state">
    <p class="muted">Постов пока нет. Создайте первый.</p>
  </div>

  <!-- лента -->
  <div v-else class="list">
    <RouterLink
      v-for="post in posts"
      :key="post.id"
      :to="{ name: 'post', params: { id: post.id } }"
      class="card post-card"
    >
      <div class="post-card__top">
        <span class="avatar">{{ String(post.author_id).slice(-2) }}</span>
        <span class="dim">автор #{{ post.author_id }}</span>
        <span class="spacer" />
        <span class="dim">#{{ post.id }}</span>
      </div>
      <h2 class="post-card__title">{{ post.title }}</h2>
      <p class="post-card__text muted">{{ excerpt(post.text) }}</p>
      <div class="post-card__cta dim">Открыть и читать комментарии →</div>
    </RouterLink>
  </div>
</template>

<style scoped>
.feed-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  margin-bottom: 20px;
}
.feed-title {
  font-size: 30px;
  font-weight: 800;
  letter-spacing: -0.03em;
  margin: 0;
}
.feed-sub {
  margin: 4px 0 0;
  font-size: 14px;
}

.compose {
  padding: 16px;
  margin-bottom: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.compose__row {
  display: flex;
  align-items: center;
  gap: 12px;
}
.compose__user {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12.5px;
}
.input--mini {
  width: 72px;
  padding: 7px 9px;
}

.list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.post-card {
  display: block;
  padding: 18px;
  transition: border-color 0.15s ease, box-shadow 0.15s ease, transform 0.08s ease;
}
.post-card:hover {
  border-color: var(--line-2);
  box-shadow: var(--shadow);
  transform: translateY(-1px);
}
.post-card__top {
  display: flex;
  align-items: center;
  gap: 9px;
  font-size: 12.5px;
  margin-bottom: 12px;
}
.avatar {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: var(--ink);
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.post-card__title {
  font-size: 18px;
  font-weight: 700;
  letter-spacing: -0.02em;
  margin: 0 0 7px;
}
.post-card__text {
  font-size: 14.5px;
  line-height: 1.5;
  margin: 0;
}
.post-card__cta {
  margin-top: 14px;
  font-size: 13px;
  font-weight: 600;
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