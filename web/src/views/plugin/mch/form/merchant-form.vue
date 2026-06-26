<!-- web/src/plugin/payment/form/merchant-form.vue -->
<template>
<el-form
    :model="localFormData"
    :rules="rules"
    ref="formRef"
    label-width="120px"
    @submit.prevent
    class="merchant-form"
>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="商户名称" prop="mch_name">
          <el-input
              v-model="localFormData.mch_name"
              placeholder="请输入商户名称"
              :disabled="isReadOnly"
          />
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item label="联系人" prop="linker">
          <el-input
              v-model="localFormData.linker"
              placeholder="请输入联系人"
              :disabled="isReadOnly"
          />
        </el-form-item>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="联系电话" prop="phone">
          <el-input
              v-model="localFormData.phone"
              placeholder="请输入联系电话"
              :disabled="isReadOnly"
          />
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item label="邮箱" prop="email">
          <el-input
              v-model="localFormData.email"
              placeholder="请输入邮箱"
              :disabled="isReadOnly"
          />
        </el-form-item>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="地址" prop="address" >
          <el-input
              v-model="localFormData.address"
              placeholder="请输入地址"
              type="textarea"
              :rows="2"
              :disabled="isReadOnly"
          />
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="localFormData.status" :disabled="isReadOnly">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="24">
        <el-form-item label="备注" prop="reason">
          <el-input
              v-model="localFormData.reason"
              placeholder="请输入备注"
              type="textarea"
              :rows="3"
              :disabled="isReadOnly"
          />
        </el-form-item>
      </el-col>
    </el-row>
  </el-form>
</template>

<script setup>
/**
 * 商户表单组件
 * @component MerchantForm
 * @description 商户信息表单组件，确保id为int类型提交
 */

import {ref, reactive, watch, nextTick, computed} from 'vue'

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({})
  },
  viewMode: {
    type: String,
    default: 'view'
  }
})

const emit = defineEmits(['update:modelValue', 'submit', 'validate-success', 'validate-fail'])

// 修复1：改用reactive创建响应式对象，确保嵌套属性响应式正常
const localFormData = reactive({
  id: 0, // 明确声明为int类型，默认值0
  mch_name: '',
  linker: '',
  phone: '',
  email: '',
  address: '',
  status: 1, // 默认启用
  reason: ''
})
// 计算属性：是否为只读模式
// 计算属性：是否为只读模式
const isReadOnly = computed(() => props.viewMode === 'view')
// 表单引用
const formRef = ref(null)

// 表单验证规则（补充address/reason规则，添加id类型校验）
const rules = reactive({
  id: [
    { type: 'number', message: 'ID必须为数字类型', trigger: 'blur' }
  ],
  mch_name: [
    { required: true, message: '请输入商户名称', trigger: 'blur' },
    { min: 1, max: 100, message: '商户名称长度应在1-100之间', trigger: 'blur' }
  ],
  linker: [
    { required: true, message: '请输入联系人', trigger: 'blur' },
    { min: 1, max: 50, message: '联系人长度应在1-50之间', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$|^(\d{3,4}-?)?\d{7,8}$/, message: '请输入正确的手机号或固定电话格式', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  address: [
    { required: false, message: '请输入地址', trigger: 'blur' },
    { min: 1, max: 200, message: '地址长度应在1-200之间', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' },
    { type: 'number', message: '状态值必须为数字', trigger: 'change' }
  ]
})

// 修复2：监听父组件传入值，强制转换id为int类型
watch(
    () => props.modelValue,
    (newVal) => {
      if (newVal) {
        // 深度合并并强制转换id为int
        Object.assign(localFormData, {
          ...newVal,
          id: newVal.id ? Number(newVal.id) : 0, // 关键：确保id是int类型
          status: newVal.status !== undefined ? Number(newVal.status) : 1,
          address: newVal.address || '',
          reason: newVal.reason || ''
        })
      }
    },
    { immediate: true, deep: true }
)

// 修复3：监听本地数据变化，同步时确保id为int
watch(
    () => localFormData,
    (newVal) => {
      // 提交前强制转换id为int
      const submitData = {
        ...newVal,
        id: Number(newVal.id) || 0 // 确保id始终是int类型
      }
      emit('update:modelValue', submitData)
    },
    { deep: true }
)

// 修复4：完善验证逻辑，确保id类型正确
const validate = async () => {
  if (!formRef.value) {
    emit('validate-fail', new Error('表单引用不存在'))
    return false
  }
  try {
    // 验证前先强制转换id为int
    localFormData.id = Number(localFormData.id) || 0
    await formRef.value.validate()
    emit('validate-success')
    return true
  } catch (error) {
    emit('validate-fail', error)
    return false
  }
}

// 修复5：提交处理，确保id为int类型
const handleSubmit = async () => {
  const isValid = await validate()
  if (isValid) {
    // 最终提交数据，确保id是int
    const submitData = {
      ...localFormData,
      id: Number(localFormData.id) || 0
    }
    emit('submit', submitData)
  }
}

// 重置表单方法
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  Object.assign(localFormData, {
    id: 0,
    mch_name: '',
    linker: '',
    phone: '',
    email: '',
    address: '',
    status: 1,
    reason: ''
  })
}

// 暴露方法给父组件
defineExpose({
  validate,
  resetForm,
  handleSubmit,
  // 暴露获取格式化数据的方法，确保id为int
  getFormData: () => ({
    ...localFormData,
    id: Number(localFormData.id) || 0
  })
})

</script>

<style scoped>
.merchant-form {
  height: 100%;
  overflow-y: auto;  /* 侧边栏内容过长时滚动 */
  padding-right: 10px;
}

:deep(.el-form-item) {
  margin-bottom: 20px;
}

/* 适配textarea在侧边栏的样式 */
:deep(.el-textarea__inner) {
  resize: vertical;
  min-height: 60px;
}
</style>