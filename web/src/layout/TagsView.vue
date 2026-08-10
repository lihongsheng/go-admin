<template>
  <div class="tags-view">
    <el-scrollbar class="tags-scrollbar">
      <div class="tags-list">
        <div v-for="tag in tagsStore.tabs" :key="tag.path" class="tag-item"
             :class="{ active: tag.path === route.path }" @click="handleClick(tag)">
          <span class="tag-title">{{ tag.title }}</span>
          <el-icon v-if="!tag.affix" class="tag-close" :size="12" @click.stop="handleClose(tag)">
            <Close />
          </el-icon>
        </div>
      </div>
    </el-scrollbar>

    <el-dropdown trigger="click" @command="handleCommand">
      <el-button text class="tags-actions" :icon="ArrowDown" />
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="closeOthers">关闭其他</el-dropdown-item>
          <el-dropdown-item command="closeAll">关闭全部</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup>
import { watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Close, ArrowDown } from '@element-plus/icons-vue'
import { useTagsStore } from '@/store/modules/tags'

const route = useRoute()
const router = useRouter()
const tagsStore = useTagsStore()

// 路由变化时自动追加 tab（刷新后首次挂载也补一次，恢复当前页 tab）
watch(() => route.path, () => tagsStore.addTab(route), { immediate: true })

function handleClick(tag) {
  if (tag.path !== route.path) router.push(tag.path)
}

function handleClose(tag) {
  const nextPath = tagsStore.removeTab(tag.path)
  if (tag.path === route.path && nextPath) router.push(nextPath)
}

function handleCommand(cmd) {
  if (cmd === 'closeOthers') {
    router.push(tagsStore.closeOthers(route.path))
  } else if (cmd === 'closeAll') {
    tagsStore.closeAll()
    router.push('/dashboard')
  }
}
</script>

<style scoped>
.tags-view {
  display: flex;
  align-items: center;
  height: 38px;
  background: #fff;
  border-bottom: 1px solid var(--el-border-color-lighter);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  padding: 0 12px;
  flex-shrink: 0;
}

.tags-scrollbar { flex: 1; }

.tags-list { display: flex; align-items: center; white-space: nowrap; }

.tag-item {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  height: 26px;
  line-height: 26px;
  padding: 0 10px;
  margin-right: 8px;
  font-size: 13px;
  color: #606266;
  background: #fff;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  cursor: pointer;
  user-select: none;
  transition: all 0.2s;
}
.tag-item:hover { color: #409eff; border-color: #409eff; }

.tag-item.active {
  color: #fff;
  background: #409eff;
  border-color: #409eff;
}

.tag-title { max-width: 160px; overflow: hidden; text-overflow: ellipsis; }

.tag-close {
  border-radius: 50%;
  padding: 1px;
  cursor: pointer;
  flex-shrink: 0;
}
.tag-close:hover { background: rgba(0, 0, 0, 0.15); }

.tags-actions { color: #606266; padding: 4px 6px; }

/* 深色模式适配 */
html.dark .tags-view { background: #1d1e1f; border-bottom-color: #363637; }
html.dark .tag-item {
  background: #1d1e1f;
  border-color: #363637;
  color: #a3a6ad;
}
html.dark .tag-item:hover { color: #409eff; border-color: #409eff; }
html.dark .tag-item.active { background: #409eff; border-color: #409eff; color: #fff; }
html.dark .tags-actions { color: #a3a6ad; }
</style>
