<template>
  <div class="tags-view">
    <!-- 左侧：前进/后退（浏览器式导航） -->
    <div class="tags-nav">
      <el-tooltip content="后退" placement="bottom">
        <el-button class="nav-btn" text :icon="ArrowLeft" @click="router.back()" />
      </el-tooltip>
      <el-tooltip content="前进" placement="bottom">
        <el-button class="nav-btn" text :icon="ArrowRight" @click="router.forward()" />
      </el-tooltip>
    </div>

    <!-- 中间：多标签页 -->
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

    <!-- 右侧：刷新 / 关闭其他 / 关闭全部 -->
    <div class="tags-actions">
      <el-tooltip content="刷新当前页" placement="bottom">
        <el-button class="action-btn" text :icon="Refresh" @click="handleRefresh" />
      </el-tooltip>
      <el-dropdown trigger="click" @command="handleCommand">
        <el-button class="action-btn" text :icon="ArrowDown" />
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="closeOthers">关闭其他</el-dropdown-item>
            <el-dropdown-item command="closeAll">关闭全部</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Close, ArrowLeft, ArrowRight, Refresh, ArrowDown } from '@element-plus/icons-vue'
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

// 刷新当前页：同路由 replace + 时间戳参数，强制组件重新挂载
function handleRefresh() {
  const { path, query } = route
  router.replace({ path, query: { ...query, t: Date.now() } })
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
  height: 48px;
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  padding: 0 12px;
  flex-shrink: 0;
}

/* 前进/后退 */
.tags-nav {
  display: flex;
  align-items: center;
  gap: 2px;
  padding-right: 10px;
  border-right: 1px solid #ebeef5;
  margin-right: 10px;
}
.nav-btn {
  color: #606266;
  padding: 6px 4px;
  border-radius: 6px;
}
.nav-btn:hover { color: #2b7de9; background: #f0f6ff; }

.tags-scrollbar { flex: 1; }

.tags-list { display: flex; align-items: center; white-space: nowrap; }

.tag-item {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  height: 34px;
  line-height: 34px;
  padding: 0 16px;
  margin-right: 10px;
  font-size: 14px;
  color: #4b5563;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 14px;
  cursor: pointer;
  user-select: none;
  transition: all 0.25s ease;
}
.tag-item:hover {
  color: #2b7de9;
  border-color: rgba(43, 125, 233, 0.5);
  background: #f0f6ff;
}

/* 激活 tab：主色渐变胶囊 + 轻投影 */
.tag-item.active {
  color: #fff;
  background: linear-gradient(90deg, #409eff, #6cb3ff);
  border-color: transparent;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.35);
}

.tag-title { max-width: 160px; overflow: hidden; text-overflow: ellipsis; }

.tag-close {
  border-radius: 50%;
  padding: 1px;
  cursor: pointer;
  flex-shrink: 0;
  transition: background 0.2s;
}
.tag-close:hover { background: rgba(0, 0, 0, 0.12); }
.tag-item.active .tag-close:hover { background: rgba(255, 255, 255, 0.35); }

/* 右侧操作 */
.tags-actions {
  display: flex;
  align-items: center;
  gap: 2px;
  padding-left: 10px;
  border-left: 1px solid #ebeef5;
}
.action-btn {
  color: #606266;
  padding: 6px 6px;
  border-radius: 6px;
}
.action-btn:hover { color: #2b7de9; background: #f0f6ff; }

/* 深色模式适配 */
html.dark .tags-view {
  background: #1d1e1f;
  border-bottom-color: #363637;
}
html.dark .tags-nav { border-right-color: #363637; }
html.dark .nav-btn { color: #a3a6ad; }
html.dark .nav-btn:hover { color: #409eff; background: #242526; }
html.dark .tag-item {
  background: #1d1e1f;
  border-color: #363637;
  color: #a3a6ad;
}
html.dark .tag-item:hover { color: #409eff; border-color: #409eff; background: #242526; }
html.dark .tag-item.active { background: #409eff; border-color: #409eff; color: #fff; box-shadow: none; }
html.dark .tags-actions { border-left-color: #363637; }
html.dark .action-btn { color: #a3a6ad; }
html.dark .action-btn:hover { color: #409eff; background: #242526; }
</style>
