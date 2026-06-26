<template>
  <div class="gva-container">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-lg font-bold">商户管理</h2>
    </div>

    <!-- 搜索区域 -->
    <div class="search-box bg-white p-4 rounded-md shadow mb-4">
      <el-form :inline="true" :model="searchInfo" class="search-form">
        <el-row :gutter="10">
          <el-col :span="6">
            <el-form-item label="商户编号">
              <el-input
                  v-model="searchInfo.mch_no"
                  placeholder="请输入商户编号"
                  clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="商户名称">
              <el-input
                  v-model="searchInfo.mch_name"
                  placeholder="请输入商户名称"
                  clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="状态">
              <el-select
                  v-model="searchInfo.status"
                  placeholder="请选择状态"
                  clearable
                  class="w-full"
              >
                <el-option label="启用" :value="1" />
                <el-option label="禁用" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="10">
          <el-col :span="24" class="text-right">
            <el-button type="primary" @click="onSubmit">查询</el-button>
            <el-button @click="onReset">重置</el-button>
            <el-button type="primary" @click="addMerchant">新增商户</el-button>
          </el-col>
        </el-row>
      </el-form>
    </div>

    <!-- 表格区域 -->
    <div class="table-container bg-white p-4 rounded-md shadow">
      <el-table
          :data="tableData"
          style="width: 100%"
          row-key="id"
          v-loading="loading"
          border
      >
        <el-table-column prop="mch_no" label="商户编号" min-width="150" show-overflow-tooltip />
        <el-table-column prop="mch_name" label="商户名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="linker" label="联系人" min-width="100" show-overflow-tooltip />
        <el-table-column prop="phone" label="联系电话" min-width="120" show-overflow-tooltip />
        <el-table-column prop="email" label="邮箱" min-width="150" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="address" label="地址" min-width="150" show-overflow-tooltip />
        <el-table-column prop="reason" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button
                size="small"
                type="primary"
                link
                @click="viewMerchant(row)"
            >
              查看
            </el-button>
            <el-button
                size="small"
                type="primary"
                link
                @click="editMerchant(row)"
            >
              编辑
            </el-button>
            <el-button
                size="small"
                :type="row.status === 1 ? 'danger' : 'success'"
                link
                @click="changeStatus(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页区域 -->
      <div class="pagination mt-4">
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
    </div>

    <!-- 抽屉侧边栏 -->
    <el-drawer
        v-model="drawerVisible"
        :title="dialogTitle"
        size="700px"
        direction="rtl"
        destroy-on-close
        :with-header="true"
        :z-index="1000"
    >
      <MerchantForm
          v-model="formData"
          :viewMode="viewMode"
          ref="formRef"
      />
      <template #footer v-if="!isViewMode">
        <div class="flex justify-end gap-2 mt-4">
          <el-button @click="closeDrawer">取消</el-button>
          <el-button type="primary" @click="enterDialog">确认</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
/**
 * 商户管理页面
 * @component MerchantView
 */

import { ref, reactive, onMounted, computed } from 'vue'
import MerchantForm from '../form/merchant-form.vue'
import {
  mchList,
  mchCreate,
  mchUpdate,
  mchDetail,
  mchChangeStatus
} from '@/api/system.js'
import { ElMessage, ElMessageBox } from 'element-plus'

// 响应式数据
const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 搜索条件
const searchInfo = reactive({
  mch_no: '',
  mch_name: '',
  status: null
})

// 侧边栏状态
const drawerVisible = ref(false)
const dialogTitle = ref('')
const formData = ref({})
const formRef = ref(null)
const viewMode = ref('add')   // 'add' | 'edit' | 'view'
const isViewMode = computed(() => viewMode.value === 'view')

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 获取商户列表
const loadList = async () => {
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

// 分页处理
const handleCurrentChange = (val) => {
  page.value = val
  loadList()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  loadList()
}

// 搜索和重置
const onSubmit = () => {
  page.value = 1
  loadList()
}

const onReset = () => {
  searchInfo.mch_no = ''
  searchInfo.mch_name = ''
  searchInfo.status = null
  page.value = 1
  loadList()
}

// 新增商户
const addMerchant = () => {
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

// 编辑商户
const editMerchant = (row) => {
  dialogTitle.value = '编辑商户'
  viewMode.value = 'edit'
  formData.value = { ...row }
  drawerVisible.value = true
}

// 查看商户
const viewMerchant = async (row) => {
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

// 关闭侧边栏
const closeDrawer = () => {
  drawerVisible.value = false
  formData.value = {}
}

// 确认提交
const enterDialog = async () => {
  try {
    // 验证表单
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

// 状态切换
const changeStatus = async (row) => {
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
    const res = await mchChangeStatus({
      id: row.id,
      status: row.status === 1 ? 2 : 1
    })
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

// 初始化
onMounted(() => {
  loadList()
})
</script>

<style scoped>
.gva-container {
  padding: 0 20px;
}

.search-box {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.08);
}

.table-container {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.08);
}

.pagination {
  text-align: right;
}

:deep(.el-drawer__footer) {
  padding: 16px;
  border-top: 1px solid #ebeef5;
  margin-top: 10px;
}
</style>
