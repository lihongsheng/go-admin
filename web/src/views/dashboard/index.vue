<template>
  <div>
    <!-- 统计卡片 -->
    <el-row :gutter="16">
      <el-col :xs="24" :sm="12" :lg="6" v-for="c in cards" :key="c.label" style="margin-bottom:16px">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-body">
            <div class="stat-icon" :style="{ background: c.color }">
              <el-icon :size="26" color="#fff"><component :is="c.icon" /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ c.value }}</div>
              <div class="stat-label">{{ c.label }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 欢迎信息 -->
    <el-card shadow="never">
      <template #header><span style="font-weight:600">系统信息</span></template>
      <el-descriptions :column="2" border size="small">
        <el-descriptions-item label="当前用户">{{ userStore.userInfo?.username }}</el-descriptions-item>
        <el-descriptions-item label="昵称">{{ userStore.userInfo?.nickname }}</el-descriptions-item>
        <el-descriptions-item label="角色">
          <el-tag v-for="r in (userStore.userInfo?.roles || [])" :key="r.id" size="small" type="success" effect="plain" style="margin-right:4px">
            {{ r.name }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="系统版本">0.1.0</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { User, Avatar, Grid, Link } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'
import { userList, roleList, menuTree, apiList } from '@/api/system'

const userStore = useUserStore()

const cards = ref([
  { label: '用户数', value: 0, icon: User, color: '#409EFF' },
  { label: '角色数', value: 0, icon: Avatar, color: '#67C23A' },
  { label: '菜单数', value: 0, icon: Grid, color: '#E6A23C' },
  { label: 'API 数', value: 0, icon: Link, color: '#F56C6C' }
])

async function loadStats() {
  try {
    const [users, roles, menus, apis] = await Promise.all([
      userList({ page: 1, size: 1 }),
      roleList(),
      menuTree(),
      apiList({})
    ])
    cards.value[0].value = users.data.total || 0
    cards.value[1].value = (roles.data.list || []).length
    cards.value[2].value = countMenuNodes(menus.data.list || [])
    cards.value[3].value = (apis.data.list || []).length
  } catch (_) { /* ignore */ }
}

function countMenuNodes(nodes) {
  let n = 0
  for (const node of nodes) {
    n++
    if (node.children) n += countMenuNodes(node.children)
  }
  return n
}

loadStats()
</script>

<style scoped>
.stat-card { cursor: pointer; transition: transform 0.2s; }
.stat-card:hover { transform: translateY(-2px); }
.stat-body { display: flex; align-items: center; gap: 16px; }
.stat-icon { width: 52px; height: 52px; border-radius: 10px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.stat-value { font-size: 26px; font-weight: 700; color: #303133; line-height: 1.2; }
html.dark .stat-value { color: #e5eaf3; }
.stat-label { font-size: 13px; color: #909399; margin-top: 2px; }
</style>
