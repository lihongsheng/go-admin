<template>
  <el-container class="layout">
    <!-- 侧边栏 -->
    <el-aside :width="collapsed ? '64px' : '220px'" class="aside">
      <div class="logo" @click="router.replace('/')">
        <span v-show="!collapsed">go-admin</span>
        <span v-show="collapsed" style="font-size:18px">GA</span>
      </div>
      <el-menu
        :default-active="route.path"
        :collapse="collapsed"
        :collapse-transition="false"
        router
        background-color="transparent"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
      >
        <SubTree
          v-for="r in permStore.routes"
          :key="r.path || r.name"
          :item="r"
          base=""
        />
      </el-menu>
    </el-aside>

    <!-- 右侧主体 -->
    <el-container class="main-container">
      <!-- 头部 -->
      <el-header class="header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="collapsed = !collapsed" :size="20">
            <Fold v-if="!collapsed" /><Expand v-else />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="route.meta?.title && route.path !== '/'">
              {{ route.meta.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>

        <div class="header-right">
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

      <!-- 内容区 -->
      <el-main class="main-content">
        <router-view />
      </el-main>
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
import { useDarkMode } from '@/composables/useDarkMode'
import router from '@/router'
import SubTree from './SubTree.vue'

const route = useRoute()
const userStore = useUserStore()
const permStore = usePermissionStore()
const dark = useDarkMode()

const collapsed = ref(false)

async function onCmd(c) {
  if (c === 'logout') {
    await userStore.doLogout()
    router.replace('/login')
  }
}
</script>

<style scoped>
.layout { height: 100%; }

/* ===== 侧边栏 ===== */
.aside {
  background: #304156;
  transition: width 0.28s;
  overflow: hidden;
}
.logo {
  height: 56px; line-height: 56px; color: #fff; text-align: center;
  font-weight: 600; font-size: 18px; cursor: pointer; user-select: none;
  border-bottom: 1px solid rgba(255,255,255,0.08);
  letter-spacing: 1px;
}
.aside .el-menu { border-right: none; }

/* ===== 头部 ===== */
.main-container { flex-direction: column; height: 100%; }
.header {
  height: 50px !important;
  background: #fff;
  display: flex; align-items: center; justify-content: space-between;
  padding: 0 16px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.08);
  z-index: 10;
}
.header-left { display: flex; align-items: center; gap: 12px; }
.collapse-btn { cursor: pointer; color: #606266; }
.collapse-btn:hover { color: #409EFF; }
.header-right { display: flex; align-items: center; gap: 8px; }
.dark-toggle { color: #606266; }
.dark-toggle:hover { color: #409EFF; }
.user-info {
  display: flex; align-items: center; gap: 6px; cursor: pointer;
  color: #303133; font-size: 13px;
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
html.dark .collapse-btn,
html.dark .dark-toggle { color: #a3a6ad; }
html.dark .collapse-btn:hover,
html.dark .dark-toggle:hover { color: #409EFF; }
html.dark .user-info { color: #e5eaf3; }
html.dark .aside { background: #1a1a1a; }
html.dark .main-content { background: #141414; }
</style>
