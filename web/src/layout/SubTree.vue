<template>
  <template v-if="hidden(item) || item.meta?.type === 'button'" />
  <el-sub-menu v-else-if="hasVisibleChildren(item)" :index="fullPath">
    <template #title>
      <el-icon v-if="item.meta?.icon"><component :is="iconMap[item.meta.icon] || MenuIcon" /></el-icon>
      <span>{{ item.meta?.title }}</span>
    </template>
    <SubTree v-for="c in visibleChildren(item)" :key="c.path || c.name" :item="c" :base="fullPath" />
  </el-sub-menu>
  <el-menu-item v-else :index="fullPath">
    <el-icon v-if="item.meta?.icon"><component :is="iconMap[item.meta.icon] || MenuIcon" /></el-icon>
    <span>{{ item.meta?.title }}</span>
  </el-menu-item>
</template>

<script setup>
import { computed } from 'vue'
import {
  DataBoard, Setting, User, Avatar, Grid, List, Link, Edit, Menu as MenuIcon
} from '@element-plus/icons-vue'

const props = defineProps({
  item: Object,
  base: { type: String, default: '' }
})

// 菜单名 → Element Plus 图标映射
const iconMap = {
  dashboard: DataBoard,
  setting: Setting,
  user: User,
  peoples: Avatar,
  'tree-table': Grid,
  component: List,
  api: Link,
  list: List,
  edit: Edit
}

const fullPath = computed(() => {
  const p = props.item.path || ''
  if (p.startsWith('/')) return p
  return (props.base.endsWith('/') ? props.base : props.base + '/') + p
})

function hidden(i) {
  if (i.meta?.hidden) return true
  if (i.meta?.type === 'button') return true
  return false
}

function hasVisibleChildren(item) {
  if (!item.children || item.children.length === 0) return false
  return item.children.some(c => !hidden({ meta: c.meta }))
}

function visibleChildren(item) {
  if (!item.children) return []
  return item.children.filter(c => !hidden({ meta: c.meta }))
}
</script>
