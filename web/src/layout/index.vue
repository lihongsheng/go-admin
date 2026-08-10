<template>
  <el-container class="layout">
    <!-- 顶部固定栏（全宽：go-admin 与 管理员 同一行） -->
    <el-header class="header">
      <div class="header-left">
        <div class="logo" @click="router.replace('/')">
          <span>go-admin</span>
        </div>
      </div>

        <div class="header-right">
          <!-- Mock 模式动态开关（适配层双环境切换演示，运行时立即生效） -->
          <el-tooltip :content="mockEnabled ? 'Mock 模式：接口返回模拟数据（点击切换真实后端）' : '真实模式：请求真实后端（点击切换 Mock）'" placement="bottom">
            <div class="mock-toggle">
              <span class="mock-label">Mock</span>
              <el-switch v-model="mockEnabled" size="small" @change="onMockToggle" />
            </div>
          </el-tooltip>

          <!-- 深色模式切换 -->
          <el-tooltip :content="dark.isDark ? '切换亮色' : '切换深色'" placement="bottom">
            <el-button class="dark-toggle" text circle @click="dark.toggle()">
              <el-icon :size="18"><Moon v-if="dark.isDark" /><Sunny v-else /></el-icon>
            </el-button>
          </el-tooltip>

          <!-- 用户菜单 -->
          <el-dropdown trigger="click" @command="onCmd">
            <span class="user-info">
              <el-avatar :size="30" icon="UserFilled" />
              <span class="nickname">{{ userStore.userInfo?.nickname || '管理员' }}</span>
              <el-icon style="vertical-align:middle"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

    <!-- 主体区：菜单栏 + 内容（菜单栏不在顶部展示，位于固定 header 下方） -->
    <el-container class="body-row">
      <!-- 侧边栏菜单 -->
      <el-aside :width="collapsed ? '64px' : '220px'" class="aside">
        <el-menu
          :default-active="route.path"
          :collapse="collapsed"
          :collapse-transition="false"
          router
        >
          <SubTree
            v-for="r in permStore.routes"
            :key="r.path || r.name"
            :item="r"
            base=""
          />
        </el-menu>

        <!-- 导航收起（菜单栏底部） -->
        <div class="aside-collapse" @click="collapsed = !collapsed">
          <el-icon :size="18"><Fold v-if="!collapsed" /><Expand v-else /></el-icon>
          <span v-show="!collapsed" class="aside-collapse-text">收起导航</span>
        </div>
      </el-aside>

      <!-- 右侧主体 -->
      <el-container class="main-container">
        <!-- 顶部多标签页（首页常驻，打开菜单自动追加） -->
        <TagsView />

        <!-- 内容区 -->
        <el-main class="main-content">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import {
  Fold, Expand, ArrowDown, Moon, Sunny, SwitchButton
} from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'
import { usePermissionStore } from '@/store/modules/permission'
import { useTagsStore } from '@/store/modules/tags'
import { useDarkMode } from '@/composables/useDarkMode'
import { getMockEnabled, setMockEnabled } from '@/api/mock'
import router from '@/router'
import SubTree from './SubTree.vue'
import TagsView from './TagsView.vue'

const route = useRoute()
const userStore = useUserStore()
const permStore = usePermissionStore()
const tagsStore = useTagsStore()
const dark = useDarkMode()

const collapsed = ref(false)

// Mock 模式动态开关：切换即持久化并刷新，干净地应用新模式
// 由 el-switch @change 触发（v-model 已更新，直接持久化当前值）
const mockEnabled = ref(getMockEnabled())
function onMockToggle() {
  setMockEnabled(mockEnabled.value)
  window.location.reload()
}

async function onCmd(c) {
  if (c === 'logout') {
    await userStore.doLogout()
    tagsStore.closeAll() // 清空历史 tab，避免切换账号后残留
    router.replace('/login')
  }
}
</script>

<style scoped>
.layout { height: 100%; }

/* ===== 顶部固定栏（全宽：go-admin + 管理员同行） ===== */
.header {
  height: 50px !important;
  background: #fff;
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 16px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  z-index: 10;
  flex-shrink: 0;
}
.header-left { display: flex; align-items: center; gap: 12px; }
.header-right { display: flex; align-items: center; gap: 8px; }

/* 品牌区 */
.logo {
  height: 50px; line-height: 50px;
  font-weight: 700; font-size: 18px; cursor: pointer; user-select: none;
  letter-spacing: 1px;
  color: #2b7de9;
  padding-right: 16px;
}

/* ===== 主体区（菜单栏在 header 下方） ===== */
.body-row { flex: 1; min-height: 0; }

/* ===== 侧边栏（白色背景） ===== */
.aside {
  background: #fff;
  transition: width 0.28s;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #ebeef5;
}
.aside .el-menu { flex: 1; overflow-y: auto; }

.aside .el-menu {
  border-right: none;
  --el-menu-bg-color: transparent;
  --el-menu-text-color: #1f1f1f;   /* 菜单黑色字 */
  --el-menu-active-color: #2b7de9;
  --el-menu-hover-bg-color: #f5f7fa;
  --el-menu-item-height: 44px;
}
/* 激活菜单：主色浅底圆角，主色文字 */
.aside .el-menu-item.is-active {
  background: rgba(64, 158, 255, 0.1);
  border-radius: 10px;
  font-weight: 600;
}
.aside .el-menu-item { border-radius: 10px; margin: 2px 8px; }
.aside .el-sub-menu__title { border-radius: 10px; margin: 2px 8px; }

/* 导航收起（菜单栏底部） */
.aside-collapse {
  display: flex; align-items: center; justify-content: center; gap: 6px;
  height: 44px; flex-shrink: 0;
  color: #606266; font-size: 13px; cursor: pointer; user-select: none;
  border-top: 1px solid #ebeef5;
  background: #fff;
  transition: all 0.25s;
}
.aside-collapse:hover {
  color: #2b7de9;
  background: #f5f7fa;
}

/* ===== 右侧主体（tab + 内容） ===== */
.main-container { flex-direction: column; height: 100%; }

/* Mock 动态开关 */
.mock-toggle {
  display: flex; align-items: center; gap: 6px;
  padding: 0 8px; cursor: pointer; user-select: none;
}
.mock-label { font-size: 12px; color: #606266; font-weight: 500; }
html.dark .mock-label { color: #a3a6ad; }

.dark-toggle { color: #606266; }
.dark-toggle:hover { color: #2b7de9; }
.user-info {
  display: flex; align-items: center; gap: 6px; cursor: pointer;
  color: #303133; font-size: 13px; font-weight: 500;
}
.nickname { max-width: 80px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* ===== 内容区 ===== */
.main-content {
  background: #f0f2f5;
  min-height: 0;
  padding: 16px;
}

/* ===== 深色模式适配 ===== */
html.dark .header {
  background: #1d1e1f;
  box-shadow: 0 1px 4px rgba(0,0,0,0.3);
}
html.dark .dark-toggle { color: #a3a6ad; }
html.dark .dark-toggle:hover { color: #409EFF; }
html.dark .aside-collapse { color: #a3a6ad; border-top-color: rgba(255, 255, 255, 0.08); background: #1a1a1a; }
html.dark .aside-collapse:hover { color: #409EFF; }
html.dark .user-info { color: #e5eaf3; }
html.dark .aside { background: #1a1a1a; border-right-color: #1a1a1a; }
/* 深色模式：菜单文字反白 */
html.dark .aside .el-menu {
  --el-menu-text-color: #cfd3dc;
  --el-menu-active-color: #409eff;
  --el-menu-hover-bg-color: #242526;
}
html.dark .aside .el-menu-item.is-active { background: rgba(64, 158, 255, 0.2); }
html.dark .logo { color: #409EFF; }
html.dark .main-content { background: #141414; }
</style>
