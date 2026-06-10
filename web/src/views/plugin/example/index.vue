<template>
  <el-card shadow="never">
    <template #header><span style="font-weight:600">示例插件 - 笔记管理</span></template>

    <div class="toolbar">
      <el-button v-permission="'example:add'" type="primary" @click="add">
        <el-icon><Plus /></el-icon>新增笔记
      </el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe border style="margin-top:12px">
      <el-table-column prop="id" label="ID" width="80" align="center" />
      <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
      <el-table-column prop="content" label="内容" min-width="300" show-overflow-tooltip />
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="s">
          <el-button v-permission="'example:del'" type="danger" link size="small" @click="del(s.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-empty v-if="!loading && list.length === 0" description="暂无笔记" />
  </el-card>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const list = ref([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const { data } = await request.get('/api/v1/plugin/example/note/list')
    list.value = data.list || []
  } finally { loading.value = false }
}

async function add() {
  try {
    const { value } = await ElMessageBox.prompt('请输入笔记标题', '新增笔记', {
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })
    if (value) {
      await request.post('/api/v1/plugin/example/note', { title: value, content: value })
      ElMessage.success('新增成功')
      load()
    }
  } catch (_) { /* cancel */ }
}

async function del(row) {
  await ElMessageBox.confirm(`确认删除笔记「${row.title}」？`, '提示', { type: 'warning' })
  await request.delete('/api/v1/plugin/example/note/' + row.id)
  ElMessage.success('删除成功')
  load()
}

load()
</script>

<style scoped>
.toolbar { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 12px; }
</style>
