<template>
  <div class="install">
    <el-card class="box" shadow="always">
      <template #header>
        <span>go-admin 安装向导</span>
      </template>
      <el-steps :active="step" align-center finish-status="success">
        <el-step title="环境检测" />
        <el-step title="数据库配置" />
        <el-step title="初始化" />
      </el-steps>

      <!-- 1. 环境检测 -->
      <div v-if="step === 0" class="pane">
        <el-alert v-if="status" :type="status.installed ? 'warning' : 'info'" show-icon
                  :title="status.installed ? '系统已安装' : (status.db_connected ? 'DB 可连接，未初始化' : '尚未配置 DB 或无法连接')" />
        <el-button type="primary" :loading="loading" @click="step = 1">下一步</el-button>
      </div>

      <!-- 2. DB 配置 -->
      <div v-else-if="step === 1" class="pane">
        <el-form :model="db" label-width="100px" size="small">
          <el-form-item label="驱动">
            <el-select v-model="db.driver" @change="onDriver">
              <el-option label="MySQL" value="mysql" />
              <el-option label="PostgreSQL" value="postgres" />
              <el-option label="SQLite" value="sqlite" />
            </el-select>
          </el-form-item>
          <template v-if="db.driver !== 'sqlite'">
            <el-form-item label="主机"><el-input v-model="db.host" /></el-form-item>
            <el-form-item label="端口"><el-input-number v-model="db.port" :min="1" :max="65535" /></el-form-item>
            <el-form-item label="用户名"><el-input v-model="db.username" /></el-form-item>
            <el-form-item label="密码"><el-input v-model="db.password" show-password /></el-form-item>
            <el-form-item label="数据库"><el-input v-model="db.database" /></el-form-item>
            <el-form-item label="charset"><el-input v-model="db.charset" /></el-form-item>
          </template>
          <template v-else>
            <el-form-item label="文件路径"><el-input v-model="db.path" placeholder="go-admin.db" /></el-form-item>
          </template>
        </el-form>
        <el-button @click="step = 0">上一步</el-button>
        <el-button :loading="loading" @click="onCheck">测试连接</el-button>
        <el-button v-if="db.driver !== 'sqlite'" :loading="loading" @click="onCreateDB">连接并自动建库</el-button>
        <el-button type="primary" :loading="loading" :disabled="!canNext" @click="step = 2">下一步</el-button>
        <el-alert v-if="checkResult" style="margin-top:12px"
                  :type="checkResult.db_connected ? 'success' : 'error'" show-icon
                  :title="checkResult.db_connected
                            ? ('连接成功' + (checkResult.created ? '（已自动创建数据库）' : ''))
                            : ('连接失败：' + (checkResult.reason || ''))" />
      </div>

      <!-- 3. 管理员 + 执行 -->
      <div v-else class="pane">
        <el-form :model="admin" label-width="100px" size="small">
          <el-form-item label="用户名"><el-input v-model="admin.username" /></el-form-item>
          <el-form-item label="密码"><el-input v-model="admin.password" show-password /></el-form-item>
          <el-form-item label="昵称"><el-input v-model="admin.nickname" /></el-form-item>
          <el-form-item label="Email"><el-input v-model="admin.email" /></el-form-item>
        </el-form>
        <el-button @click="step = 1">上一步</el-button>
        <el-button type="primary" :loading="loading" @click="onInstall">开始安装</el-button>

        <el-table v-if="steps.length" :data="steps" size="small" style="margin-top:12px">
          <el-table-column prop="table" label="表/步骤" />
          <el-table-column prop="state" label="状态">
            <template #default="s">
              <el-tag :type="tagType(s.row.state)" size="small">{{ s.row.state }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="err" label="错误" />
        </el-table>

        <el-result v-if="done" icon="success" title="安装完成"
                   sub-title="3 秒后跳转登录页">
          <template #extra>
            <el-button type="primary" @click="router.replace('/login')">立即登录</el-button>
          </template>
        </el-result>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getInstallStatus, checkDB, doInstall } from '@/api/install'
import { refreshInstallStatus } from '@/permission'

const router = useRouter()

const step = ref(0)
const loading = ref(false)
const status = ref(null)
const checkResult = ref(null)
const db = reactive({ driver: 'mysql', host: '', port: 3306, username: 'root', password: '', database: 'go_admin', charset: 'utf8mb4', path: '' })
const admin = reactive({ username: 'admin', password: '', nickname: '管理员', email: '' })
const steps = ref([])
const done = ref(false)

const canNext = computed(() => checkResult.value && checkResult.value.db_connected && !checkResult.value.installed)

function tagType(s) {
  return s === 'done' ? 'success'
       : s === 'failed' ? 'danger'
       : s === 'migrating' || s === 'seeding' ? '' : 'info'
}

function onDriver() {
  if (db.driver === 'postgres') db.port = 5432
  if (db.driver === 'mysql') db.port = 3306
}

async function onCheck() {
  loading.value = true
  try { const r = await checkDB({ ...db }); checkResult.value = r.data }
  finally { loading.value = false }
}

async function onCreateDB() {
  loading.value = true
  try {
    const r = await checkDB({ ...db, create_if_missing: true })
    checkResult.value = r.data
    if (r.data.created) ElMessage.success('数据库不存在，已自动创建')
  } finally { loading.value = false }
}

async function onInstall() {
  loading.value = true
  steps.value = []
  try {
    const r = await doInstall({ db: { ...db }, admin: { ...admin } })
    steps.value = r.data.steps || []
    if (r.data.ok) {
      if (r.data.created_db) ElMessage.success('数据库不存在，已自动创建')
      done.value = true
      await refreshInstallStatus()
      setTimeout(() => router.replace('/login'), 1500)
    } else {
      ElMessage.error(r.data.error || '安装失败')
    }
  } finally { loading.value = false }
}

// 加载初始状态
getInstallStatus().then(r => status.value = r.data).catch(() => {})
</script>

<style scoped>
.install { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: #f5f7fa; }
.box { width: 720px; }
.pane { margin-top: 24px; }
</style>
