<template>
  <el-card shadow="never">
    <template #header><span style="font-weight:600">已安装插件</span></template>
    <el-table :data="list" v-loading="loading" stripe border>
      <el-table-column prop="name" label="插件名称" min-width="200" />
      <el-table-column prop="version" label="版本" width="120" align="center">
        <template #default="s">
          <el-tag size="small" effect="plain">{{ s.row.version }}</el-tag>
        </template>
      </el-table-column>
    </el-table>
    <el-empty v-if="!loading && list.length === 0" description="暂无已安装插件" />
  </el-card>
</template>

<script setup>
import { ref } from 'vue'
import { pluginList } from '@/api/system'

const list = ref([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const { data } = await pluginList()
    list.value = data.list || []
  } finally { loading.value = false }
}

load()
</script>
