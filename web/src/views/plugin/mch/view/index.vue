<template>
  <div class="page-wrap">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchInfo" class="search-form">
        <el-form-item label="商户编号">
          <el-input v-model="searchInfo.mch_no" placeholder="请输入商户编号" clearable style="width:160px"
                    @keyup.enter="onSubmit" @clear="onSubmit" />
        </el-form-item>
        <el-form-item label="商户名称">
          <el-input v-model="searchInfo.mch_name" placeholder="请输入商户名称" clearable style="width:160px"
                    @keyup.enter="onSubmit" @clear="onSubmit" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="请选择状态" clearable style="width:120px" @change="onSubmit">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据栏 -->
    <el-card shadow="never" class="table-card">
      <div class="table-toolbar">
        <el-button v-permission="'mch:add'" type="primary" @click="addMerchant">
          <el-icon><Plus /></el-icon>新增商户
        </el-button>
      </div>
      <el-table :data="tableData" v-loading="loading" stripe border row-key="id">
        <el-table-column prop="mch_no" label="商户编号" min-width="150" show-overflow-tooltip />
        <el-table-column prop="mch_name" label="商户名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="linker" label="联系人" min-width="100" show-overflow-tooltip />
        <el-table-column prop="phone" label="联系电话" min-width="120" show-overflow-tooltip />
        <el-table-column prop="email" label="邮箱" min-width="150" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small" effect="dark">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="address" label="地址" min-width="150" show-overflow-tooltip />
        <el-table-column prop="reason" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="更新时间" width="170">
          <template #default="{ row }">
            {{ formatTime(row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button v-permission="'mch:view'" type="primary" link size="small" @click="viewMerchant(row)">
              查看
            </el-button>
            <el-button v-permission="'mch:edit'" type="primary" link size="small" @click="editMerchant(row)">
              编辑
            </el-button>
            <el-button
              v-permission="'mch:status'"
              :type="row.status === 1 ? 'danger' : 'success'"
              link size="small"
              @click="changeStatus(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          background
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 20, 30, 50]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </el-card>

    <!-- 新增/编辑/查看弹窗 -->
    <el-dialog v-model="drawerVisible" :title="dialogTitle" width="700px" destroy-on-close>
      <MerchantForm v-model="formData" :viewMode="viewMode" ref="formRef" />
      <template #footer v-if="!isViewMode">
        <el-button @click="closeDrawer">取消</el-button>
        <el-button type="primary" @click="enterDialog">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
/**
 * 商户管理页面
 * @component MerchantView
 */
import { ref, reactive, onMounted, computed } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import MerchantForm from '../form/merchant-form.vue'
import { mchList, mchCreate, mchUpdate, mchDetail, mchChangeStatus } from '@/api/system.js'

const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const searchInfo = reactive({
  mch_no: '',
  mch_name: '',
  status: null
})

const drawerVisible = ref(false)
const dialogTitle = ref('')
const formData = ref({})
const formRef = ref(null)
const viewMode = ref('add')
const isViewMode = computed(() => viewMode.value === 'view')

function formatTime(t) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

async function loadList() {
  loading.value = true
  try {
    const params = {
      mch_no: searchInfo.mch_no || undefined,
      mch_name: searchInfo.mch_name || undefined,
      status: searchInfo.status || 0,
      page: page.value,
      limit: pageSize.value
    }
    const res = await mchList(params)
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    } else {
      ElMessage.error(res.msg || '获取商户列表失败')
    }
  } catch (error) {
    console.error('获取商户列表失败:', error)
    ElMessage.error('获取商户列表失败')
  } finally {
    loading.value = false
  }
}

function handleCurrentChange(val) {
  page.value = val
  loadList()
}

function handleSizeChange(val) {
  pageSize.value = val
  page.value = 1
  loadList()
}

function onSubmit() {
  page.value = 1
  loadList()
}

function onReset() {
  searchInfo.mch_no = ''
  searchInfo.mch_name = ''
  searchInfo.status = null
  page.value = 1
  loadList()
}

function addMerchant() {
  dialogTitle.value = '新增商户'
  viewMode.value = 'add'
  formData.value = {
    id: 0,
    mch_name: '',
    linker: '',
    phone: '',
    email: '',
    status: 1,
    address: '',
    reason: ''
  }
  drawerVisible.value = true
}

function editMerchant(row) {
  dialogTitle.value = '编辑商户'
  viewMode.value = 'edit'
  formData.value = { ...row }
  drawerVisible.value = true
}

async function viewMerchant(row) {
  try {
    const res = await mchDetail(row.id)
    if (res.code === 0) {
      viewMode.value = 'view'
      dialogTitle.value = '查看商户'
      formData.value = { ...res.data }
      drawerVisible.value = true
    } else {
      ElMessage.error(res.msg || '获取商户详情失败')
    }
  } catch (error) {
    console.error('获取商户详情失败:', error)
    ElMessage.error('获取商户详情失败')
  }
}

function closeDrawer() {
  drawerVisible.value = false
  formData.value = {}
}

async function enterDialog() {
  try {
    if (formRef.value) {
      const isValid = await formRef.value.validate()
      if (!isValid) return
    }
    const data = formRef.value ? formRef.value.getFormData() : formData.value
    const isEdit = viewMode.value === 'edit'
    const res = isEdit ? await mchUpdate(data) : await mchCreate(data)
    if (res.code === 0) {
      ElMessage.success(isEdit ? '编辑商户成功' : '新增商户成功')
      drawerVisible.value = false
      loadList()
    } else {
      ElMessage.error(res.msg || '操作失败')
    }
  } catch (error) {
    console.error('提交商户信息失败:', error)
    ElMessage.error('提交商户信息失败')
  }
}

async function changeStatus(row) {
  try {
    await ElMessageBox.confirm(
      `确定要${row.status === 1 ? '禁用' : '启用'}此商户吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: row.status === 1 ? 'warning' : 'info'
      }
    )
    const res = await mchChangeStatus({ id: row.id, status: row.status === 1 ? 2 : 1 })
    if (res.code === 0) {
      ElMessage.success(`商户已${row.status === 1 ? '禁用' : '启用'}`)
      loadList()
    } else {
      ElMessage.error(res.msg || '状态修改失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('修改商户状态失败:', error)
      ElMessage.error('修改商户状态失败')
    }
  }
}

onMounted(() => {
  loadList()
})
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
