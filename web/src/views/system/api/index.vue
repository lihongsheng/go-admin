<template>
  <div class="page-wrap">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="query" class="search-form">
        <el-form-item label="分组">
          <el-input v-model="query.group" placeholder="按分组过滤" clearable style="width:200px"
                    @keyup.enter="load" @clear="load">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item label="方法">
          <el-select v-model="query.method" placeholder="全部" clearable style="width:120px" @change="load">
            <el-option v-for="m in ['GET','POST','PUT','DELETE','PATCH']" :key="m" :label="m" :value="m" />
          </el-select>
        </el-form-item>
        <el-form-item label="系统">
          <el-select v-model="query.system_type" placeholder="系统类型" style="width:120px" @change="load">
            <el-option v-for="st in systemTypesList" :key="st.SystemType" :label="st.Name" :value="st.SystemType" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="load">
            <el-icon><Search /></el-icon>搜索
          </el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据栏 -->
    <el-card shadow="never" class="table-card">
      <div class="table-toolbar">
        <el-button v-permission="'api:add'" type="primary" @click="open()">
          <el-icon><Plus /></el-icon>新增 API
        </el-button>
      </div>
      <el-table :data="list" v-loading="loading" stripe border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="path" label="接口路径" min-width="220" show-overflow-tooltip />
        <el-table-column prop="method" label="方法" width="100" align="center">
          <template #default="s">
            <el-tag :type="methodTag(s.row.method)" size="small" effect="dark">{{ s.row.method }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="group" label="分组" width="120" align="center">
          <template #default="s">
            <el-tag v-if="s.row.group" size="small" type="info" effect="plain">{{ s.row.group }}</el-tag>
            <span v-else style="color:#c0c4cc">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="desc" label="说明" min-width="180" show-overflow-tooltip />
        <el-table-column label="系统" width="80" align="center">
          <template #default="s">
            <el-tag v-if="s.row.system_type" size="small" type="warning">{{ systemTypeName[s.row.system_type] }}</el-tag>
            <span v-else style="color:#c0c4cc">平台</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="s">
            <el-button v-permission="'api:edit'" type="primary" link size="small" @click="open(s.row)">编辑</el-button>
            <el-button v-permission="'api:del'"  type="danger"  link size="small" @click="del(s.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dlg" :title="form.id ? '编辑 API' : '新增 API'" width="500px" destroy-on-close>
      <el-form :model="form" label-width="70px">
        <el-form-item label="路径" required>
          <el-input v-model="form.path" placeholder="/api/v1/system/user" />
        </el-form-item>
        <el-form-item label="方法" required>
          <el-select v-model="form.method" style="width:100%">
            <el-option v-for="m in ['GET','POST','PUT','DELETE','PATCH']" :key="m" :label="m" :value="m">
              <el-tag :type="methodTag(m)" size="small" effect="dark" style="margin-right:8px">{{ m }}</el-tag>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="分组">
          <el-input v-model="form.group" placeholder="如 user / role / menu" />
        </el-form-item>
        <el-form-item label="系统">
          <el-select v-model="form.system_type" style="width:100%">
            <el-option v-for="st in systemTypesList" :key="st.SystemType" :label="st.Name" :value="st.SystemType" />
          </el-select>
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="form.desc" placeholder="接口功能描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dlg = false">取消</el-button>
        <el-button type="primary" @click="submit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'
import { apiList, apiCreate, apiUpdate, apiDelete } from '@/api/system'
import { systemTypes } from '@/api/base'

const list = ref([])
const loading = ref(false)
const dlg = ref(false)
const submitting = ref(false)
const query = reactive({ group: '', method: '', system_type: 1 }) // 默认选中第一个，加载后覆盖
const form = reactive({ method: 'GET' })
const systemTypesList = ref([])

const systemTypeName = computed(() => {
  const m = {}
  systemTypesList.value.forEach(st => { m[st.SystemType] = st.Name })
  return m
})

async function loadSystemTypes() {
  try {
    const { data } = await systemTypes()
    systemTypesList.value = data || []
    // 默认选中第一个
    if (systemTypesList.value.length > 0) {
      query.system_type = systemTypesList.value[0].SystemType
      load()
    }
  } catch (_) { /* ignore */ }
}

function methodTag(m) {
  const map = { GET: 'success', POST: '', PUT: 'warning', DELETE: 'danger', PATCH: 'info' }
  return map[m] || ''
}

function resetQuery() {
  query.group = ''
  query.method = ''
  load()
}

async function load() {
  loading.value = true
  try {
    const { data } = await apiList({ group: query.group, system_type: query.system_type })
    list.value = data.list || []
  } finally { loading.value = false }
}

function open(row) {
  if (row) { Object.assign(form, { ...row }) }
  else { Object.assign(form, { id: undefined, path: '', method: 'GET', group: '', desc: '', system_type: query.system_type }) }
  dlg.value = true
}

async function submit() {
  submitting.value = true
  try {
    if (form.id) { await apiUpdate({ ...form }); ElMessage.success('编辑成功') }
    else { await apiCreate({ ...form }); ElMessage.success('新增成功') }
    dlg.value = false
    load()
  } finally { submitting.value = false }
}

async function del(row) {
  await ElMessageBox.confirm(`确认删除 API「${row.method} ${row.path}」？`, '提示', { type: 'warning' })
  await apiDelete(row.id)
  ElMessage.success('删除成功')
  load()
}

loadSystemTypes()
</script>

<style scoped>
.page-wrap { display: flex; flex-direction: column; gap: 12px; }
.search-card { }
.search-form { display: flex; align-items: center; flex-wrap: wrap; gap: 4px; }
.search-form .el-form-item { margin-bottom: 0; }
.table-card { }
.table-toolbar { margin-bottom: 12px; }
</style>
