<template>
  <el-form
    :model="localFormData"
    :rules="rules"
    ref="formRef"
    label-width="100px"
    @submit.prevent
  >
    <el-form-item label="商户名称" prop="mch_name">
      <el-input v-model="localFormData.mch_name" placeholder="请输入商户名称" :disabled="isReadOnly" />
    </el-form-item>
    <el-form-item label="联系人" prop="linker">
      <el-input v-model="localFormData.linker" placeholder="请输入联系人" :disabled="isReadOnly" />
    </el-form-item>
    <el-form-item label="联系电话" prop="phone">
      <el-input v-model="localFormData.phone" placeholder="请输入联系电话" :disabled="isReadOnly" />
    </el-form-item>
    <el-form-item label="邮箱" prop="email">
      <el-input v-model="localFormData.email" placeholder="请输入邮箱" :disabled="isReadOnly" />
    </el-form-item>
    <el-form-item label="地址" prop="address">
      <el-input v-model="localFormData.address" placeholder="请输入地址" type="textarea" :rows="2" :disabled="isReadOnly" />
    </el-form-item>
    <el-form-item label="状态" prop="status">
      <el-radio-group v-model="localFormData.status" :disabled="isReadOnly">
        <el-radio :label="1">启用</el-radio>
        <el-radio :label="2">禁用</el-radio>
      </el-radio-group>
    </el-form-item>
    <el-form-item label="备注" prop="reason">
      <el-input v-model="localFormData.reason" placeholder="请输入备注" type="textarea" :rows="3" :disabled="isReadOnly" />
    </el-form-item>
  </el-form>
</template>

<script setup>
/**
 * 商户表单组件
 * @component MerchantForm
 */
import { ref, reactive, watch, computed } from 'vue'

const props = defineProps({
  modelValue: { type: Object, default: () => ({}) },
  viewMode: { type: String, default: 'add' }
})

const emit = defineEmits(['update:modelValue'])

const localFormData = reactive({
  id: 0,
  mch_name: '',
  linker: '',
  phone: '',
  email: '',
  address: '',
  status: 1,
  reason: ''
})

const isReadOnly = computed(() => props.viewMode === 'view')
const formRef = ref(null)

const rules = reactive({
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
    { min: 1, max: 200, message: '地址长度应在1-200之间', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
})

watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal) {
      Object.assign(localFormData, {
        ...newVal,
        id: newVal.id ? Number(newVal.id) : 0,
        status: newVal.status !== undefined ? Number(newVal.status) : 1,
        address: newVal.address || '',
        reason: newVal.reason || ''
      })
    }
  },
  { immediate: true, deep: true }
)

watch(
  () => localFormData,
  (newVal) => {
    emit('update:modelValue', {
      ...newVal,
      id: Number(newVal.id) || 0
    })
  },
  { deep: true }
)

async function validate() {
  if (!formRef.value) return false
  try {
    localFormData.id = Number(localFormData.id) || 0
    await formRef.value.validate()
    return true
  } catch {
    return false
  }
}

function resetForm() {
  if (formRef.value) formRef.value.resetFields()
  Object.assign(localFormData, {
    id: 0, mch_name: '', linker: '', phone: '', email: '',
    address: '', status: 1, reason: ''
  })
}

defineExpose({
  validate,
  resetForm,
  getFormData: () => ({
    ...localFormData,
    id: Number(localFormData.id) || 0
  })
})
</script>
