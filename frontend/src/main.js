import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'

import './styles.css'
import App from './App.vue'
import FeedView from './views/FeedView.vue'
import PostView from './views/PostView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'feed', component: FeedView },
    { path: '/posts/:id', name: 'post', component: PostView, props: true },
  ],
})

createApp(App).use(router).mount('#app')