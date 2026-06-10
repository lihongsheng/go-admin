<template>
  <div class="page-wrap">
    <el-card shadow="never" class="table-card">
      <div class="table-toolbar">
        <el-button v-permission="'role:add'" type="primary" @click="open()">
          <el-icon><Plus /></el-icon>新增角色
        </el-button>
      </div>

      <el-table :data="list" v-loading="loading" stripe border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column prop="name" label="角色名称" min-width="140" />
      <el-table-column prop="code" label="角色标识" min-width="140" />
      <el-table-column prop="remark" label="备注" min-width="160" show-overflow-tooltip />
      <el-table-column label="状态" width="80" align="center">
        <template #default="s">
          <el-tag :type="s.row.status === 1 ? 'success' : 'info'" size="small" effect="dark">
            {{ s.row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="170">
        <template #default="s">{{ formatTime(s.row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="240" fixed="right">
        <template #default="s">
          <el-button v-permission="'role:auth'" type="success" link size="small" @click="auth(s.row)">授权</el-button>
          <el-button v-permission="'role:edit'" type="primary" link size="small" @click="open(s.row)">编辑</el-button>
          <el-button v-permission="'role:del'"  type="danger"  link size="small" @click="del(s.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新增/编辑角色 -->
    <el-dialog v-model="dlg" :title="form.id ? '编辑角色' : '新增角色'" width="480px" destroy-on-close>
      <el-form :model="form" label-width="70px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="角色名称" />
        </el-form-item>
        <el-form-item label="标识" required>
          <el-input v-model="form.code" placeholder="角色标识（英文）" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="备注信息" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.statusBool" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dlg = false">取消</el-button>
        <el-button type="primary" @click="submit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 角色授权 -->
    <el-dialog v-model="authDlg" title="角色授权" width="680px" destroy-on-close>
      <el-tabs v-model="authTab" type="border-card">
        <el-tab-pane label="菜单权限" name="menu">
          <el-tree ref="menuTreeRef" :data="menuTreeData" show-checkbox node-key="id"
                   :props="{ label: (d) => d.title || d.name, children: 'children' }"
                   :default-checked-keys="authForm.menu_ids" default-expand-all highlight-current />
        </el-tab-pane>
        <el-tab-pane label="API 权限" name="api">
          <el-checkbox-group v-model="authForm.api_ids">
            <div v-for="g in apiGroups" :key="g.name" class="api-group">
              <div class="api-group-title">{{ g.name }}</div>
              <el-checkbox v-for="a in g.items" :key="a.id" :label="a.id" class="api-item">
                <el-tag size="small" :type="methodTag(a.method)">{{ a.method }}</el-tag>
                <span class="api-path">{{ a.path }}</span>
              </el-checkbox>
            </div>
          </el-checkbox-group>
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="authDlg = false">取消</el-button>
        <el-button type="primary" @click="submitAuth" :loading="submitting">保存授权</el-button>
      </template>
    </el-dialog>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { roleList, roleCreate, roleUpdate, roleDelete, roleAuth, roleAuthDetail, menuTree, apiList } from '@/api/system'

const list = ref([])
const loading = ref(false)
const dlg = ref(false)
const submitting = ref(false)
const form = reactive({ statusBool: true, status: 1 })
const authDlg = ref(false)
const authTab = ref('menu')
const authForm = reactive({ role_id: 0, menu_ids: [], api_ids: [] })
const menuTreeData = ref([])
const apiGroups = ref([])
const menuTreeRef = ref(null)

function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

function methodTag(m) {
  const map = { GET: 'success', POST: '', PUT: 'warning', DELETE: 'danger', PATCH: 'info' }
  return map[m] || ''
}

async function load() {
  loading.value = true
  try {
    const { data } = await roleList()
    list.value = data.list || []
  } finally { loading.value = false }
}

function open(row) {
  if (row) {
    Object.assign(form, { ...row, statusBool: row.status === 1 })
  } else {
    Object.assign(form, { id: undefined, name: '', code: '', remark: '', statusBool: true, status: 1 })
  }
  dlg.value = true
}

async function submit() {
  submitting.value = true
  try {
    const body = { ...form, status: form.statusBool ? 1 : 0 }
    delete body.statusBool
    if (form.id) {
      await roleUpdate(body)
      ElMessage.success('编辑成功')
    } else {
      await roleCreate(body)
      ElMessage.success('新增成功')
    }
    dlg.value = false
    load()
  } finally { submitting.value = false }
}

async function del(row) {
  await ElMessageBox.confirm(`确认删除角色「${row.name}」？`, '提示', { type: 'warning' })
  await roleDelete(row.id)
  ElMessage.success('删除成功')
  load()
}

// 提取所有节点 id，构建 id -> 是否为父节点 的映射
function collectParentIds(tree, parentSet = new Set()) {
  for (const node of tree) {
    if (node.children && node.children.length > 0) {
      parentSet.add(node.id)
      collectParentIds(node.children, parentSet)
    }
  }
  return parentSet
}

async function auth(row) {
  try {
    const [m, a, detail] = await Promise.all([
      menuTree(),
      apiList(),
      roleAuthDetail(row.id).catch(() => ({ data: { menu_ids: [], api_ids: [] } }))
    ])
    menuTreeData.value = m.data.list || []
    const map = {}
    ;(a.data.list || []).forEach(x => {
      const key = x.group || '默认'
      if (!map[key]) map[key] = []
      map[key].push(x)
    })
    apiGroups.value = Object.entries(map).map(([name, items]) => ({ name, items }))
    authForm.role_id = row.id
    // 关键修复：el-tree 在非 check-strictly 模式下，勾选父节点会自动勾选全部子节点。
    // 后端为保持菜单树完整性会保存所有父级 ID，因此渲染时必须过滤掉父级 ID，
    // 只把叶子节点的 ID 设为默认勾选，让 el-tree 自行推断父级的全选/半选状态。
    const parentSet = collectParentIds(menuTreeData.value)
    const allIds = detail.data.menu_ids || []
    authForm.menu_ids = allIds.filter(id => !parentSet.has(id))
    authForm.api_ids = detail.data.api_ids || []
    authDlg.value = true
  } catch (e) {
    ElMessage.error('加载授权数据失败')
  }
}

async function submitAuth() {
  submitting.value = true
  try {
    // 关键修复：合并 getCheckedKeys()（完全选中）和 getHalfCheckedKeys()（半选中的父节点）
    // 否则当取消某个子菜单时，半选状态的父级菜单 ID 不会被提交，导致父级权限丢失
    const checked = menuTreeRef.value.getCheckedKeys() || []
    const halfChecked = menuTreeRef.value.getHalfCheckedKeys() || []
    authForm.menu_ids = [...checked, ...halfChecked]
    await roleAuth({ ...authForm })
    authDlg.value = false
    ElMessage.success('授权已保存')
  } finally { submitting.value = false }
}

load()
</script>

<style scoped>
.page-wrap { display: flex; flex-direction: column; gap: 12px; }
.table-card { }
.table-toolbar { margin-bottom: 12px; }
.api-group { margin-bottom: 16px; }
.api-group-title { font-weight: 600; font-size: 13px; color: #606266; margin-bottom: 8px; padding-bottom: 4px; border-bottom: 1px solid #ebeef5; }
.api-item { display: flex; align-items: center; gap: 8px; margin: 4px 0 4px 8px; }
.api-path { font-size: 12px; color: #909399; font-family: monospace; }
</style>
