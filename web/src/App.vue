<template>
  <div id="app"><router-view /></div>
</template>
<script setup>
import { useDarkMode } from '@/composables/useDarkMode'
// 初始化深色模式（确保 watchEffect 尽早生效）
useDarkMode()
</script>
<style>
/* 全局基础样式 */
* { margin: 0; padding: 0; box-sizing: border-box; }
html, body, #app { height: 100%; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }

/* 深色模式全局过渡 */
html.dark { color-scheme: dark; }
html.dark body { background: #141414; color: #e5eaf3; }
</style>
