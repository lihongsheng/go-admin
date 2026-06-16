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
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="s">
            <el-button v-permission="'user:edit'" type="primary" link size="small" @click="open(s.row)">编辑</el-button>
            <el-button v-permission="'user:reset'" type="warning" link size="small">重置密码</el-button>
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

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dlg" :title="form.id ? '编辑用户' : '新增用户'" width="520px" destroy-on-close>
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
        <el-form-item label="角色">
          <el-select v-model="form.role_ids" multiple placeholder="请选择角色" style="width:100%">
            <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
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
import { ref, reactive } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Search, Plus, Upload } from '@element-plus/icons-vue'
import { userList, userCreate, userUpdate, userDelete, roleList, uploadFile } from '@/api/system'

const list = ref([])
const total = ref(0)
const page = ref(1)
const size = ref(10)
const loading = ref(false)
const dlg = ref(false)
const submitting = ref(false)
const query = reactive({ kw: '' })
const form = reactive({ statusBool: true, role_ids: [], status: 1, phone: '', avatar: '' })
const roles = ref([])
const avatarFile = ref(null)
const avatarInput = ref(null)

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

async function fetchRoles() {
  const { data } = await roleList()
  roles.value = data.list || []
}

async function load() {
  loading.value = true
  try {
    const { data } = await userList({ page: page.value, size: size.value, keyword: query.kw })
    list.value = data.list || []
    total.value = data.total || 0
  } finally { loading.value = false }
}

function open(row) {
  avatarFile.value = null
  if (row) {
    Object.assign(form, { ...row, statusBool: row.status === 1, role_ids: (row.roles || []).map(r => r.id), password: '', avatar: row.avatar || '' })
  } else {
    Object.assign(form, { id: undefined, username: '', password: '', nickname: '', email: '', phone: '', avatar: '', statusBool: true, role_ids: [], status: 1 })
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

fetchRoles()
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
</style>
