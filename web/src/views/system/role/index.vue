<template>
  <div class="page-wrap">
    <el-card shadow="never" class="table-card">
      <div class="table-toolbar">
        <template v-if="isPlatformUser">
          <el-select v-model="currentSysType" placeholder="系统类型" style="width:140px;margin-right:12px" @change="load">
            <el-option v-for="t in SYS_TYPES" :key="t.system_type" :label="t.name" :value="t.system_type" />
          </el-select>
        </template>
        <el-button v-permission="'role:add'" type="primary" @click="open()">
          <el-icon><Plus /></el-icon>新增角色
        </el-button>
      </div>

      <el-table :data="list" v-loading="loading" stripe border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column prop="name" label="角色名称" min-width="140" />
      <el-table-column prop="default_router" label="默认首页" min-width="140">
        <template #default="s">
          <el-tag v-if="s.row.default_router" type="warning" size="small" effect="light">{{ s.row.default_router }}</el-tag>
          <span v-else class="text-muted">/dashboard（默认）</span>
        </template>
      </el-table-column>
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
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="角色名称" />
        </el-form-item>
        <el-form-item label="默认首页">
          <el-input v-model="form.default_router" placeholder="默认首页路由，如 /dashboard" />
          <div class="form-tip">登录后跳转的默认页面，为空则使用 /dashboard</div>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="备注信息" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.statusBool" active-text="启用" inactive-text="禁用" />
        </el-form-item>
        <template v-if="isPlatformUser">
          <el-form-item label="系统类型">
            <el-select v-model="form.system_type" placeholder="选择系统类型" @change="onSysTypeChange" style="width:100%">
              <el-option v-for="t in SYS_TYPES" :key="t.system_type" :label="t.name" :value="t.system_type" />
            </el-select>
          </el-form-item>
          <el-form-item label="所属商户" v-if="form.system_type === 1">
            <el-select v-model="form.mch_id" placeholder="输入商户名称搜索" filterable remote
                       :remote-method="searchMch" :loading="mchLoading" style="width:100%">
              <el-option v-for="m in merchants" :key="m.id" :label="m.mch_name" :value="m.id" />
            </el-select>
          </el-form-item>
        </template>
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
          <!-- 当前首页提示 -->
          <el-alert v-if="authForm.default_router" type="warning" :closable="false" show-icon style="margin-bottom:12px">
            <template #title>
              当前默认首页：<b>{{ authForm.default_router }}</b>
              <el-button type="danger" link size="small" @click="clearDefaultRouter">清除</el-button>
            </template>
          </el-alert>
          <el-alert v-else type="info" :closable="false" show-icon style="margin-bottom:12px">
            <template #title>未设置默认首页，登录后默认跳转 /dashboard</template>
          </el-alert>

          <el-tree ref="menuTreeRef" :data="menuTreeData" show-checkbox node-key="id"
                   :props="{ label: (d) => d.title || d.name, children: 'children' }"
                   :default-checked-keys="authForm.menu_ids" default-expand-all highlight-current>
            <template #default="{ node, data }">
              <span class="tree-node">
                <span>{{ node.label }}</span>
                <el-tag v-if="data.type" size="small" :type="typeTag(data.type)" class="tree-type-tag">{{ typeLabel(data.type) }}</el-tag>
                <template v-if="data.type === 'menu'">
                  <el-tag v-if="isDefaultRouter(node)" type="warning" size="small" effect="dark" class="home-tag">首页</el-tag>
                  <el-button v-else type="warning" link size="small" class="set-home-btn" @click.stop="setDefaultRouter(node)">设为首页</el-button>
                </template>
              </span>
            </template>
          </el-tree>
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
import { ref, reactive, computed } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'
import { roleList, roleCreate, roleUpdate, roleDelete, roleAuth, roleAuthDetail, roleSetDefaultRouter, menuTree, apiList, mchList } from '@/api/system'

// 系统类型常量（与服务端 enum.SystemType 保持一致）
const SYS_TYPES = [
  { system_type: 0, name: '平台' },
  { system_type: 1, name: '商户' },
  { system_type: 2, name: '代理' },
]

const list = ref([])
const loading = ref(false)
const dlg = ref(false)
const submitting = ref(false)
const form = reactive({ statusBool: true, status: 1, default_router: '', system_type: 0, mch_id: null })
const currentSysType = ref(0)
const authDlg = ref(false)
const authTab = ref('menu')
const authForm = reactive({ role_id: 0, menu_ids: [], api_ids: [], default_router: '' })
const menuTreeData = ref([])
const apiGroups = ref([])
const menuTreeRef = ref(null)

const userStore = useUserStore()
const isPlatformUser = computed(() => userStore.userInfo?.system_type === 0)
const merchants = ref([])
const mchLoading = ref(false)

function onSysTypeChange(val) {
  form.mch_id = null
  if (val === 1) searchMch('')
}

async function searchMch(keyword) {
  mchLoading.value = true
  try {
    const params = { page: 1, limit: 10 }
    if (keyword) params.mch_name = keyword
    const { data } = await mchList(params)
    merchants.value = data.list || []
  } catch (_) {
    merchants.value = []
  } finally {
    mchLoading.value = false
  }
}

function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

function methodTag(m) {
  const map = { GET: 'success', POST: '', PUT: 'warning', DELETE: 'danger', PATCH: 'info' }
  return map[m] || ''
}

function typeTag(t) {
  const map = { catalog: '', menu: 'success', button: 'info' }
  return map[t] || ''
}

function typeLabel(t) {
  const map = { catalog: '目录', menu: '菜单', button: '按钮' }
  return map[t] || t
}

// 从 tree node 向上遍历，拼接完整路由路径（如 /system/user）
function buildFullPath(node) {
  const segs = []
  let cur = node
  while (cur) {
    const d = cur.data
    if (d && d.path && d.type === 'catalog') {
      segs.unshift(d.path.replace(/^\//, ''))
    } else if (d && d.path && d.type === 'menu') {
      segs.push(d.path.replace(/^\//, ''))
    }
    cur = cur.parent
  }
  return '/' + segs.filter(Boolean).join('/')
}

function isDefaultRouter(node) {
  const full = buildFullPath(node)
  return full === authForm.default_router
}

async function setDefaultRouter(node) {
  const full = buildFullPath(node)
  try {
    await roleSetDefaultRouter(authForm.role_id, { default_router: full })
    authForm.default_router = full
    ElMessage.success(`已设置默认首页：${full}`)
  } catch (e) {
    ElMessage.error('设置默认首页失败')
  }
}

async function clearDefaultRouter() {
  try {
    await roleSetDefaultRouter(authForm.role_id, { default_router: '' })
    authForm.default_router = ''
    ElMessage.info('已清除默认首页，将使用 /dashboard')
  } catch (e) {
    ElMessage.error('清除默认首页失败')
  }
}

async function load() {
  loading.value = true
  try {
    const params = {}
    if (isPlatformUser.value) params.system_type = currentSysType.value
    const { data } = await roleList(params)
    list.value = data.list || []
  } finally { loading.value = false }
}

function open(row) {
  if (row) {
    Object.assign(form, { ...row, statusBool: row.status === 1, default_router: row.default_router || '' })
    // 编辑商户角色时加载商户列表
    if (isPlatformUser.value && row.system_type === 1) {
      searchMch('')
    }
  } else {
    const defaultSysType = isPlatformUser.value ? currentSysType.value : 0
    Object.assign(form, { id: undefined, name: '', remark: '', default_router: '', statusBool: true, status: 1, system_type: defaultSysType, mch_id: null })
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
    const detail = await roleAuthDetail(row.id).catch(() => ({ data: { menu_ids: [], api_ids: [], default_router: '', system_type: 0 } }))
    const sysType = detail.data.system_type || 0
    const params = sysType > 0 ? { system_type: sysType } : {}
    const [m, a] = await Promise.all([
      menuTree(params),
      apiList(params),
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
    authForm.default_router = detail.data.default_router || ''
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
    await roleAuth({
      role_id: authForm.role_id,
      menu_ids: authForm.menu_ids,
      api_ids: authForm.api_ids
    })
    authDlg.value = false
    ElMessage.success('授权已保存')
    load()
  } finally { submitting.value = false }
}

load()
</script>

<style scoped>
.page-wrap { display: flex; flex-direction: column; gap: 12px; }
.table-card { }
.table-toolbar { margin-bottom: 12px; }
.text-muted { color: #909399; font-size: 12px; }
.form-tip { color: #909399; font-size: 12px; margin-top: 4px; }
.api-group { margin-bottom: 16px; }
.api-group-title { font-weight: 600; font-size: 13px; color: #606266; margin-bottom: 8px; padding-bottom: 4px; border-bottom: 1px solid #ebeef5; }
.api-item { display: flex; align-items: center; gap: 8px; margin: 4px 0 4px 8px; }
.api-path { font-size: 12px; color: #909399; font-family: monospace; }

/* 菜单树节点样式 */
.tree-node {
  display: inline-flex; align-items: center; gap: 6px;
  width: 100%;
}
.tree-type-tag {
  flex-shrink: 0;
}
.set-home-btn {
  flex-shrink: 0; font-size: 11px; padding: 0 4px;
}
.home-tag {
  flex-shrink: 0;
}
</style>
