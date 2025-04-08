<script setup>
defineProps({
  msg: {
    type: String,
    required: true
  }
})
</script>

<template>
  <h1>{{ msg }}</h1>
  <p>
    Recommended IDE setup:
    <a href="https://code.visualstudio.com/" target="_blank">VSCode</a>
    +
    <a href="https://marketplace.visualstudio.com/items?itemName=Vue.volar" target="_blank">Volar</a>
  </p>
  <p>
    <a href="https://vitejs.dev/guide/features.html" target="_blank">Vite Documentation</a> |
    <a href="https://v3.vuejs.org/" target="_blank">Vue 3 Documentation</a>
  </p>
</template>

<style scoped>
a {
  color: #42b883;
}
</style> 