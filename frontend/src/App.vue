<script setup>
import { RouterView, RouterLink } from 'vue-router'
</script>

<template>
  <header class="topbar">
    <div class="container topbar__inner">
      <RouterLink to="/" class="brand">Posts</RouterLink>
      <span class="spacer" />
      <span class="pill"><span class="dot dot--live" /> realtime</span>
    </div>
  </header>

  <main class="container" style="padding-top: 24px; padding-bottom: 80px">
    <RouterView />
  </main>
</template>