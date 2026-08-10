<template>
  <div class="page-wrap">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="query" class="search-form">
        <el-form-item label="资源名称">
          <el-input v-model="query.keyword" placeholder="请输入资源名称" clearable style="width:220px"
                    @keyup.enter="onSearch" @clear="onSearch">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSearch">
            <el-icon><Search /></el-icon>查询
          </el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据栏 -->
    <el-card shadow="never" class="table-card">
      <div class="table-toolbar">
        <el-button v-permission="'resource:del'" type="danger" plain :disabled="selectedIds.length === 0" @click="batchDel">
          <el-icon><Delete /></el-icon>批量删除
        </el-button>
        <el-button v-permission="'resource:add'" type="primary" @click="openAdd">
          <el-icon><Plus /></el-icon>新增资源
        </el-button>
      </div>
      <el-table :data="list" v-loading="loading" stripe border
                @selection-change="onSelectionChange">
        <el-table-column type="selection" width="50" align="center" />
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="资源名称" min-width="180" show-overflow-tooltip />
        <el-table-column label="资源类型" min-width="140">
          <template #default="s">
            <el-tag size="small" effect="plain">{{ s.row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90" align="center">
          <template #default="s">
            <el-tag :type="statusTag(s.row.status)" size="small" effect="dark">
              {{ statusLabel(s.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态更新时间" width="170" show-overflow-tooltip>
          <template #default="s">{{ formatTime(s.row.status_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="s">
            <el-button type="primary" link size="small" @click="openView(s.row)">查看</el-button>
            <el-button v-permission="'resource:edit'" type="primary" link size="small" @click="openEdit(s.row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination background layout="total,sizes,prev,pager,next,jumper"
                       :total="total" v-model:current-page="page" v-model:page-size="size"
                       :page-sizes="[10, 20, 50, 100]" @current-change="load" @size-change="load" />
      </div>
    </el-card>

    <!-- 新增资源（右侧抽屉） -->
    <el-drawer v-model="addDlg" title="新增资源" size="480px" destroy-on-close>
      <el-form ref="addFormRef" :model="addForm" :rules="rules" label-width="90px">
        <el-form-item label="资源名称" prop="name">
          <el-input v-model="addForm.name" placeholder="请输入资源名称" maxlength="128" />
        </el-form-item>
        <el-form-item label="资源类型" prop="type">
          <el-select v-model="addForm.type" placeholder="请选择或输入资源类型" filterable allow-create
                     default-first-option style="width:100%">
            <el-option v-for="t in TYPE_OPTIONS" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="addForm.status">
            <el-radio-button v-for="s in STATUS_OPTIONS" :key="s.value" :value="s.value">{{ s.label }}</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="drawer-footer">
          <el-button @click="addDlg = false">取消</el-button>
          <el-button type="primary" @click="submitAdd" :loading="submitting">确定</el-button>
        </div>
      </template>
    </el-drawer>

    <!-- 查看资源（右侧抽屉） -->
    <el-drawer v-model="viewDlg" title="资源详情" size="480px" destroy-on-close>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="资源名称">{{ viewForm.name }}</el-descriptions-item>
        <el-descriptions-item label="资源类型">
          <el-tag size="small" effect="plain">{{ viewForm.type }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusTag(viewForm.status)" size="small" effect="dark">{{ statusLabel(viewForm.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态更新时间">{{ formatTime(viewForm.status_at) }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatTime(viewForm.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatTime(viewForm.updated_at) }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <div class="drawer-footer">
          <el-button @click="viewDlg = false">关闭</el-button>
        </div>
      </template>
    </el-drawer>

    <!-- 编辑资源（右侧抽屉） -->
    <el-drawer v-model="editDlg" title="编辑资源" size="480px" destroy-on-close>
      <el-form ref="editFormRef" :model="editForm" :rules="rules" label-width="90px">
        <el-form-item label="资源名称" prop="name">
          <el-input v-model="editForm.name" placeholder="请输入资源名称" maxlength="128" />
        </el-form-item>
        <el-form-item label="资源类型" prop="type">
          <el-select v-model="editForm.type" placeholder="请选择或输入资源类型" filterable allow-create
                     default-first-option style="width:100%">
            <el-option v-for="t in TYPE_OPTIONS" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="editForm.status">
            <el-radio-button v-for="s in STATUS_OPTIONS" :key="s.value" :value="s.value">{{ s.label }}</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="drawer-footer">
          <el-button @click="editDlg = false">取消</el-button>
          <el-button type="primary" @click="submitEdit" :loading="submitting">确定</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Delete } from '@element-plus/icons-vue'
import { resourceList, resourceDetail, resourceCreate, resourceUpdate, resourceBatchDelete } from '@/api/resource'

// 资源类型选项（ECS 实例规格族，可手动输入）
const TYPE_OPTIONS = ['通用型 g6', '计算型 c6', '内存型 r6', '大数据型 d6', '突发性能型 t6', '共享型 s6']

// 运行状态（与服务端 model.StatusXxx 保持一致）
const STATUS_OPTIONS = [
  { value: 1, label: '运行中' },
  { value: 2, label: '已停止' },
  { value: 3, label: '已释放' },
]

const STATUS_MAP = { 1: { label: '运行中', tag: 'success' }, 2: { label: '已停止', tag: 'info' }, 3: { label: '已释放', tag: 'danger' } }
const statusLabel = (s) => STATUS_MAP[s]?.label || '-'
const statusTag = (s) => STATUS_MAP[s]?.tag || ''

const rules = {
  name: [{ required: true, message: '请输入资源名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择或输入资源类型', trigger: 'change' }],
}

const list = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const loading = ref(false)
const query = reactive({ keyword: '' })
const selectedIds = ref([])

const addDlg = ref(false)
const addFormRef = ref(null)
const addForm = reactive({ name: '', type: '', status: 1 })
const viewDlg = ref(false)
const viewForm = reactive({ name: '', type: '', status: 0, status_at: '', created_at: '', updated_at: '' })
const editDlg = ref(false)
const editFormRef = ref(null)
const editForm = reactive({ id: 0, name: '', type: '', status: 1 })
const submitting = ref(false)

function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

async function load() {
  loading.value = true
  try {
    const params = { page: page.value, limit: size.value, keyword: query.keyword }
    const { data } = await resourceList(params)
    list.value = data.list || []
    total.value = data.total || 0
  } finally { loading.value = false }
}

function onSearch() {
  page.value = 1
  load()
}

function onReset() {
  query.keyword = ''
  page.value = 1
  load()
}

function openAdd() {
  Object.assign(addForm, { name: '', type: '', status: 1 })
  addDlg.value = true
}

async function submitAdd() {
  try {
    await addFormRef.value.validate()
  } catch (_) { return }
  submitting.value = true
  try {
    await resourceCreate({ ...addForm })
    ElMessage.success('新增成功')
    addDlg.value = false
    load()
  } finally { submitting.value = false }
}

async function openView(row) {
  const { data } = await resourceDetail(row.id)
  Object.assign(viewForm, data)
  viewDlg.value = true
}

function openEdit(row) {
  Object.assign(editForm, { id: row.id, name: row.name, type: row.type, status: row.status })
  editDlg.value = true
}

async function submitEdit() {
  try {
    await editFormRef.value.validate()
  } catch (_) { return }
  submitting.value = true
  try {
    await resourceUpdate(editForm.id, { ...editForm })
    ElMessage.success('编辑成功')
    editDlg.value = false
    load()
  } finally { submitting.value = false }
}

function onSelectionChange(rows) {
  selectedIds.value = rows.map(r => r.id)
}

async function batchDel() {
  const ids = [...selectedIds.value]
  if (ids.length === 0) return
  await ElMessageBox.confirm(`确认删除选中的 ${ids.length} 个资源？`, '提示', { type: 'warning' })
  await resourceBatchDelete(ids)
  ElMessage.success('删除成功')
  load()
}

load()
</script>

<style scoped>
.page-wrap { display: flex; flex-direction: column; gap: 12px; }
.search-card { }
.search-form { display: flex; align-items: center; flex-wrap: wrap; gap: 4px; }
.search-form .el-form-item { margin-bottom: 0; }
.table-card { }
.table-toolbar { display: flex; justify-content: flex-end; margin-bottom: 12px; }
.pagination-wrap { display: flex; justify-content: flex-end; margin-top: 16px; }

/* 抽屉底部操作区 */
.drawer-footer { display: flex; justify-content: flex-end; gap: 12px; }
</style>
