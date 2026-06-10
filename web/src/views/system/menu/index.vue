<template>
  <div class="page-wrap">
    <el-card shadow="never" class="table-card">
      <div class="table-toolbar">
        <el-button v-permission="'menu:add'" type="primary" @click="open()">
          <el-icon><Plus /></el-icon>新增菜单
        </el-button>
      </div>

      <el-table :data="flat" v-loading="loading" row-key="id" default-expand-all stripe border
                :tree-props="{ children: 'children' }">
      <el-table-column prop="title" label="标题" width="180" show-overflow-tooltip />
      <el-table-column prop="type" label="类型" width="80" align="center">
        <template #default="s">
          <el-tag :type="tagType(s.row.type)" size="small" effect="dark">{{ typeLabel(s.row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="icon" label="图标" width="80" align="center" />
      <el-table-column prop="name" label="路由名" width="130" show-overflow-tooltip />
      <el-table-column prop="path" label="路径" width="150" show-overflow-tooltip />
      <el-table-column prop="permission" label="权限码" width="120" show-overflow-tooltip>
        <template #default="s">
          <el-tag v-if="s.row.permission" size="small" type="warning" effect="plain">{{ s.row.permission }}</el-tag>
          <span v-else style="color:#c0c4cc">-</span>
        </template>
      </el-table-column>
      <el-table-column prop="component" label="组件路径" width="200" show-overflow-tooltip />
      <el-table-column prop="sort" label="排序" width="60" align="center" />
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="s">
          <el-button v-permission="'menu:edit'" type="primary" link size="small" @click="open(s.row)">编辑</el-button>
          <el-button v-permission="'menu:del'"  type="danger"  link size="small" @click="del(s.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dlg" :title="form.id ? '编辑菜单' : '新增菜单'" width="580px" destroy-on-close>
      <el-form :model="form" label-width="90px">
        <el-form-item label="类型" required>
          <el-radio-group v-model="form.type">
            <el-radio-button label="catalog">目录</el-radio-button>
            <el-radio-button label="menu">菜单</el-radio-button>
            <el-radio-button label="button">按钮</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="父菜单 ID">
          <el-input-number v-model="form.parent_id" :min="0" placeholder="0 表示顶级" />
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="标题">
              <el-input v-model="form.title" placeholder="显示标题" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="路由名">
              <el-input v-model="form.name" placeholder="Vue Router name" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item v-if="form.type !== 'button'" label="路径">
          <el-input v-model="form.path" placeholder="路由路径，如 /system 或 user" />
        </el-form-item>
        <el-form-item v-if="form.type !== 'button'" label="组件">
          <el-input v-model="form.component" placeholder="Layout 或 system/user/index" />
        </el-form-item>
        <el-form-item v-if="form.type === 'catalog'" label="重定向">
          <el-input v-model="form.redirect" placeholder="如 /system/user" />
        </el-form-item>
        <el-form-item v-if="form.type === 'button'" label="权限码">
          <el-input v-model="form.permission" placeholder="如 user:add" />
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item v-if="form.type !== 'button'" label="图标">
              <el-input v-model="form.icon" placeholder="element-plus 图标名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序">
              <el-input-number v-model="form.sort" :min="0" controls-position="right" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12" v-if="form.type !== 'button'">
            <el-form-item label="隐藏">
              <el-switch v-model="form.hidden" />
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.type === 'menu'">
            <el-form-item label="缓存">
              <el-switch v-model="form.keep_alive" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dlg = false">取消</el-button>
        <el-button type="primary" @click="submit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { menuTree, menuCreate, menuUpdate, menuDelete } from '@/api/system'

const TYPE_MAP = { catalog: '目录', menu: '菜单', button: '按钮' }
const TAG_MAP  = { catalog: '', menu: 'success', button: 'warning' }

const flat = ref([])
const loading = ref(false)
const dlg = ref(false)
const submitting = ref(false)
const form = reactive({ parent_id: 0, type: 'menu', sort: 0, hidden: false, keep_alive: false })

function typeLabel(t) { return TYPE_MAP[t] || t }
function tagType(t)   { return TAG_MAP[t] || '' }

async function load() {
  loading.value = true
  try {
    const { data } = await menuTree()
    flat.value = data.list || []
  } finally { loading.value = false }
}

function open(row) {
  if (row) {
    Object.assign(form, { ...row })
  } else {
    Object.assign(form, {
      id: undefined, parent_id: 0, type: 'menu', sort: 0, hidden: false, keep_alive: false,
      title: '', name: '', path: '', component: '', permission: '', redirect: '', icon: ''
    })
  }
  dlg.value = true
}

async function submit() {
  submitting.value = true
  try {
    if (form.id) {
      await menuUpdate({ ...form })
      ElMessage.success('编辑成功')
    } else {
      await menuCreate({ ...form })
      ElMessage.success('新增成功')
    }
    dlg.value = false
    load()
  } finally { submitting.value = false }
}

async function del(row) {
  await ElMessageBox.confirm('确认删除？子节点将被一并删除', '警告', { type: 'warning' })
  await menuDelete(row.id)
  ElMessage.success('删除成功')
  load()
}

load()
</script>

<style scoped>
.page-wrap { display: flex; flex-direction: column; gap: 12px; }
.table-card { }
.table-toolbar { margin-bottom: 12px; }
</style>
