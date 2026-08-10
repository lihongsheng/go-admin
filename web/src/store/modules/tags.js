import { defineStore } from 'pinia'
import { ref } from 'vue'

// 常驻首页 tab（不可关闭）
const AFFIX_TAB = { path: '/dashboard', title: '首页', affix: true }

// 顶部多标签页状态：打开菜单路由时自动追加 tab，首页常驻
export const useTagsStore = defineStore('tags', () => {
  const tabs = ref([{ ...AFFIX_TAB }])

  // 路由 → tab（仅叶子页面：有 title、非 catalog 容器）
  function addTab(route) {
    const title = route.meta?.title
    if (!title || route.meta?.type === 'catalog' || route.path === '/') return
    if (tabs.value.some(t => t.path === route.path)) return
    tabs.value.push({ path: route.path, title, affix: false })
  }

  // 移除指定 tab，返回应跳转的路径（移除激活 tab 时取左侧相邻 tab）
  function removeTab(path) {
    const idx = tabs.value.findIndex(t => t.path === path)
    if (idx === -1 || tabs.value[idx].affix) return null
    const nextPath = tabs.value[Math.max(idx - 1, 0)].path
    tabs.value.splice(idx, 1)
    return nextPath
  }

  // 关闭其他：保留首页 + 当前 tab
  function closeOthers(path) {
    tabs.value = tabs.value.filter(t => t.affix || t.path === path)
    return path
  }

  // 关闭全部：只保留首页
  function closeAll() {
    tabs.value = [{ ...AFFIX_TAB }]
  }

  return { tabs, addTab, removeTab, closeOthers, closeAll }
})
