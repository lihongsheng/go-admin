<template>
  <div>
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
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()
</script>
