<template>
  <div class="page-wrap">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="query" class="search-form">
        <el-form-item label="关键字">
          <el-input v-model="query.kw" placeholder="用户名 / 昵称" clearable style="width:220px"
                    @keyup.enter="load" @clear="load">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
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
        <el-button v-permission="'user:add'" type="primary" @click="open()">
          <el-icon><Plus /></el-icon>新增用户
        </el-button>
      </div>
      <el-table :data="list" v-loading="loading" stripe border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="nickname" label="昵称" min-width="120" />
        <el-table-column prop="email" label="邮箱" min-width="160" show-overflow-tooltip />
        <el-table-column label="角色" min-width="180">
          <template #default="s">
            <el-tag v-for="r in s.row.roles" :key="r.id" size="small" type="success" effect="plain" style="margin-right:4px">
              {{ r.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80" align="center">
          <template #default="s">
            <el-tag :type="s.row.status === 1 ? 'success' : 'info'" size="small" effect="dark">
              {{ s.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="170" show-overflow-tooltip>
          <template #default="s">{{ formatTime(s.row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="s">
            <el-button v-permission="'user:edit'" type="primary" link size="small" @click="open(s.row)">编辑</el-button>
            <el-button v-permission="'user:del'"  type="danger"  link size="small" @click="del(s.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination background layout="total,sizes,prev,pager,next,jumper"
                       :total="total" v-model:current-page="page" v-model:page-size="size"
                       :page-sizes="[10, 20, 50, 100]" @current-change="load" @size-change="load" />
      </div>
    </el-card>

    <!-- 新增/编辑用户（右侧抽屉） -->
    <el-drawer v-model="dlg" :title="form.id ? '编辑用户' : '新增用户'" size="520px" destroy-on-close>
      <el-form ref="formRef" :model="form" label-width="70px">
        <el-form-item label="用户名" required>
          <el-input v-model="form.username" :disabled="!!form.id" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" :required="!form.id">
          <el-input v-model="form.password" show-password :placeholder="form.id ? '留空不修改' : '请输入密码'" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="form.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="头像">
          <div style="display:flex;align-items:center;gap:12px">
            <el-avatar :size="48" :src="form.avatar" style="flex-shrink:0">
              <el-icon :size="24"><Upload /></el-icon>
            </el-avatar>
            <input ref="avatarInput" type="file" accept="image/*" style="display:none" @change="onAvatarPick" />
            <el-button size="small" type="primary" plain @click="$refs.avatarInput.click()">
              <el-icon><Upload /></el-icon>选择图片
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
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
                       :remote-method="searchMch" :loading="mchLoading" style="width:100%"
                       @change="onMchChange">
              <el-option v-for="m in merchants" :key="m.id" :label="m.mch_name" :value="m.id" />
            </el-select>
          </el-form-item>
        </template>
        <el-form-item label="角色">
          <el-select v-model="form.role_ids" multiple placeholder="请选择角色" style="width:100%">
            <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="drawer-footer">
          <el-button @click="dlg = false">取消</el-button>
          <el-button type="primary" @click="submit" :loading="submitting">确定</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Search, Plus, Upload } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/modules/user'
import { userList, userCreate, userUpdate, userDelete, roleList, uploadFile, mchList } from '@/api/system'

// 系统类型常量（与服务端 enum.SystemType 保持一致）
const SYS_TYPES = [
  { system_type: 0, name: '平台' },
  { system_type: 1, name: '商户' },
  { system_type: 2, name: '代理' },
]

const list = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const loading = ref(false)
const dlg = ref(false)
const submitting = ref(false)
const query = reactive({ kw: '' })
const form = reactive({ statusBool: true, role_ids: [], status: 1, phone: '', avatar: '', system_type: 0, mch_id: null })
const roles = ref([])
const avatarFile = ref(null)
const avatarInput = ref(null)

const userStore = useUserStore()
const isPlatformUser = computed(() => userStore.userInfo?.system_type === 0)
const isMerchantUser = computed(() => userStore.userInfo?.system_type === 1)
const merchants = ref([])
const mchLoading = ref(false)

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

function onSysTypeChange(val) {
  form.mch_id = null
  form.role_ids = []
  // 按选择的系统类型加载角色
  const params = { system_type: val }
  fetchRoles(params)
  if (val === 1) searchMch('')
}

function onMchChange(val) {
  if (!val) return
  form.role_ids = []
  // 按选中的商户加载角色
  fetchRoles({ system_type: 1, mch_id: val })
}

function onAvatarPick(e) {
  const file = e.target.files[0]
  if (!file) return
  avatarFile.value = file
  // 本地预览
  form.avatar = URL.createObjectURL(file)
  e.target.value = ''
  // 立即上传
  uploadFile(file).then(res => {
    form.avatar = res.data.url
    avatarFile.value = null
  }).catch(() => {
    ElMessage.error('头像上传失败')
    form.avatar = ''
    avatarFile.value = null
  })
}

function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

function resetQuery() {
  query.kw = ''
  page.value = 1
  load()
}

async function fetchRoles(params) {
  const { data } = await roleList(params)
  roles.value = data.list || []
}

async function load() {
  loading.value = true
  try {
    const params = { page: page.value, size: size.value, keyword: query.kw }
    const { data } = await userList(params)
    list.value = data.list || []
    total.value = data.total || 0
  } finally { loading.value = false }
}

function open(row) {
  avatarFile.value = null
  if (row) {
    Object.assign(form, { ...row, statusBool: row.status === 1, role_ids: (row.roles || []).map(r => r.id), password: '', avatar: row.avatar || '' })
    // 编辑时按角色所属系统类型加载角色列表
    const roleParams = {}
    if (row.system_type > 0) roleParams.system_type = row.system_type
    if (row.mch_id > 0) roleParams.mch_id = row.mch_id
    fetchRoles(roleParams)
    // 如果是商户用户，加载商户列表供选择
    if (isPlatformUser.value && row.system_type === 1 && row.mch_id > 0) {
      searchMch('')
    }
  } else {
    Object.assign(form, { id: undefined, username: '', password: '', nickname: '', email: '', phone: '', avatar: '', statusBool: true, role_ids: [], status: 1, system_type: 0, mch_id: null })
  }
  dlg.value = true
}

async function submit() {
  submitting.value = true
  try {
    const body = { ...form, status: form.statusBool ? 1 : 0 }
    delete body.statusBool
    if (body.id) { await userUpdate(body); ElMessage.success('编辑成功') }
    else { await userCreate(body); ElMessage.success('新增成功') }
    dlg.value = false
    load()
  } finally { submitting.value = false }
}

async function del(row) {
  await ElMessageBox.confirm(`确认删除用户「${row.username}」？`, '提示', { type: 'warning' })
  await userDelete(row.id)
  ElMessage.success('删除成功')
  load()
}

// 初始化：商户用户自动加载其商户的角色列表
if (isMerchantUser.value) {
  fetchRoles()
} else {
  fetchRoles({ system_type: 0 })
}
load()
</script>

<style scoped>
.page-wrap { display: flex; flex-direction: column; gap: 12px; }
.search-card { }
.search-form { display: flex; align-items: center; flex-wrap: wrap; gap: 4px; }
.search-form .el-form-item { margin-bottom: 0; }
.table-card { }
.table-toolbar { margin-bottom: 12px; }
.pagination-wrap { display: flex; justify-content: flex-end; margin-top: 16px; }

/* 抽屉底部操作区 */
.drawer-footer { display: flex; justify-content: flex-end; gap: 12px; }
</style>
