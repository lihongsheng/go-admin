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
          <div style="display:flex;gap:6px;width:100%">
            <el-tree-select v-model="form.component" :data="componentTreeData"
                            filterable clearable placeholder="选择组件"
                            :props="{ label: 'label', children: 'children' }"
                            :filter-node-method="filterComponentNode"
                            class="component-tree-select" style="flex:1" />
            <el-input v-model="form.component" placeholder="或手动输入" style="width:160px" />
          </div>
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
              <el-select v-model="form.icon" filterable allow-create clearable
                         placeholder="选择或输入图标名">
                <el-option v-for="ic in iconOptions" :key="ic" :label="ic" :value="ic">
                  <el-icon style="margin-right:6px;vertical-align:middle">
                    <component :is="iconMap[ic]" />
                  </el-icon>
                  <span>{{ ic }}</span>
                </el-option>
              </el-select>
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
import { ref, reactive, computed } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import {
  Plus, Aim, Avatar, Bell, Bicycle, Box, Brush, Calendar, CircleCheck, CircleClose,
  Clock, Collection, CollectionTag, Connection, Coordinate, Cpu, DataAnalysis,
  DataBoard, Delete, Document, Download, Edit, Files, Folder, FullScreen, Goods,
  Grid, Guide, Headset, Histogram, HomeFilled, InfoFilled, Key, Link, List, Lock,
  Location, MagicStick, Medal, Menu, Message, Microphone, Monitor, OfficeBuilding,
  Operation, PieChart, Picture, Postcard, PriceTag, Promotion, Rank, Refresh,
  Search, Sell, SetUp, Setting, Share, Ship, ShoppingCart, Star, Switch,
  SwitchButton, TakeawayBox, Ticket, Timer, Tools, TrendCharts, Trophy, Unlock,
  Upload, User, UserFilled, Van, VideoCamera, Warning
} from '@element-plus/icons-vue'
import { menuTree, menuCreate, menuUpdate, menuDelete } from '@/api/system'

// 图标名(kebab) → 图标组件映射
const iconMap = {
  aim: Aim, avatar: Avatar, bell: Bell, bicycle: Bicycle, box: Box, brush: Brush,
  calendar: Calendar, 'circle-check': CircleCheck, 'circle-close': CircleClose,
  clock: Clock, collection: Collection, 'collection-tag': CollectionTag,
  connection: Connection, coordinate: Coordinate, cpu: Cpu,
  'data-analysis': DataAnalysis, 'data-board': DataBoard, dashboard: DataBoard,
  delete: Delete, document: Document, download: Download, edit: Edit,
  files: Files, folder: Folder, 'full-screen': FullScreen, goods: Goods,
  grid: Grid, guide: Guide, headset: Headset, histogram: Histogram,
  'home-filled': HomeFilled, 'info-filled': InfoFilled, key: Key, link: Link,
  list: List, lock: Lock, location: Location, 'magic-stick': MagicStick,
  medal: Medal, menu: Menu, message: Message, microphone: Microphone,
  monitor: Monitor, 'office-building': OfficeBuilding, operation: Operation,
  'pie-chart': PieChart, picture: Picture, plus: Plus, postcard: Postcard,
  'price-tag': PriceTag, promotion: Promotion, rank: Rank, refresh: Refresh,
  search: Search, sell: Sell, 'set-up': SetUp, setting: Setting, share: Share,
  ship: Ship, 'shopping-cart': ShoppingCart, star: Star, switch: Switch,
  'switch-button': SwitchButton, 'takeaway-box': TakeawayBox, ticket: Ticket,
  timer: Timer, tools: Tools, 'trend-charts': TrendCharts, trophy: Trophy,
  unlock: Unlock, upload: Upload, user: User, 'user-filled': UserFilled,
  van: Van, 'video-camera': VideoCamera, warning: Warning
}

// ── 可用组件路径（自动扫描 src/views 目录，go-vue-admin 风格）──
const viewModules = import.meta.glob('../../**/*.vue')
const componentOptions = Object.keys(viewModules)
  .filter(p => !p.includes('/layout/'))       // 排除 layout 组件（非页面）
  .map(p => {
    // 去除所有 ../ 前缀，再去除 .vue 后缀
    return p.replace(/^(\.\.\/)+/, '').replace(/\.vue$/, '')
  })
  .filter(p => p && !p.includes('..'))         // 最终兜底：排除仍含 .. 的非法路径
  .sort()

// 构建 el-tree-select 树数据：目录可展开，文件可选
const componentTreeData = computed(() => {
  const root = { label: 'views', children: [] }

  for (const path of componentOptions) {
    const parts = path.split('/')
    let node = root
    for (let i = 0; i < parts.length; i++) {
      const isLast = i === parts.length - 1
      let child = node.children.find(c => c.label === parts[i])
      if (!child) {
        child = { label: parts[i] }
        if (isLast) {
          child.value = path  // 叶子节点：可选中
        } else {
          child.children = []  // 目录节点：可展开
        }
        node.children.push(child)
      } else if (isLast) {
        // 同名文件已作为目录存在（如 views/plugin 是目录，但也可有 views/plugin 文件）
        // 不过正常文件系统不会出现，保留目录优先
        if (!child.children) child.children = []
      }
      node = child
    }
  }

  // 递归清理空 children 的目录节点（避免 el-tree-select 显示 ... 占位符）
  function clean(node) {
    if (node.children) {
      node.children = node.children.filter(c => {
        if (c.children && c.children.length === 0) return false
        return true
      })
      if (node.children.length === 0) delete node.children
      else node.children.forEach(clean)
    }
  }
  clean(root)

  return [root]
})

// 树节点筛选：匹配标签或值中包含输入文本的节点（目录和文件都匹配）
function filterComponentNode(value, data) {
  if (!value) return true
  const v = value.toLowerCase()
  const text = (data.label || '').toLowerCase()
  const val = (data.value || '').toLowerCase()
  return text.includes(v) || val.includes(v)
}

// ── 常用图标 ──
const iconOptions = [
  'dashboard', 'setting', 'user', 'user-filled', 'avatar', 'list', 'grid',
  'menu', 'edit', 'delete', 'plus', 'search', 'refresh', 'upload', 'download',
  'link', 'share', 'message', 'bell', 'clock', 'timer', 'calendar', 'star',
  'folder', 'document', 'files', 'data-board', 'data-analysis', 'pie-chart',
  'histogram', 'trend-charts', 'office-building', 'home-filled', 'location',
  'monitor', 'cpu', 'connection', 'set-up', 'operation', 'guide', 'tools',
  'key', 'lock', 'unlock', 'warning', 'info-filled', 'circle-check',
  'circle-close', 'promotion', 'sell', 'ticket', 'collection', 'collection-tag',
  'price-tag', 'shopping-cart', 'goods', 'postcard', 'picture', 'video-camera',
  'headset', 'microphone', 'switch', 'switch-button', 'full-screen', 'rank',
  'aim', 'medal', 'trophy', 'magic-stick', 'brush', 'coordinate',
  'box', 'takeaway-box', 'ship', 'bicycle', 'truck', 'van',
]

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

/* 组件树选择器：让树选择器宽度与表单项一致 */
.component-tree-select { width: 100%; }
</style>
