<template>
  <div class="login">
    <el-card class="box">
      <template #header>go-admin 登录</template>
      <el-form :model="form" size="small" label-width="60px">
        <el-form-item label="账号"><el-input v-model="form.username" /></el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" show-password @keyup.enter="submit" />
        </el-form-item>
        <el-form-item label="验证码">
          <div class="captcha-row">
            <el-input v-model="form.captcha_code" maxlength="5" @keyup.enter="submit" />
            <img v-if="captchaB64" :src="captchaB64" class="captcha-img" title="点击刷新" @click="loadCaptcha" />
          </div>
        </el-form-item>
        <el-button type="primary" :loading="loading" style="width:100%" @click="submit">登录</el-button>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { captcha } from '@/api/base'
import { useUserStore } from '@/store/modules/user'

const router = useRouter()
const userStore = useUserStore()

const form = reactive({ username: 'admin', password: '', captcha_id: '', captcha_code: '' })
const captchaB64 = ref('')
const loading = ref(false)

async function loadCaptcha() {
  const { data } = await captcha()
  form.captcha_id = data.captcha_id
  captchaB64.value = data.captcha_b64
  form.captcha_code = ''
}

async function submit() {
  loading.value = true
  try {
    await userStore.doLogin({ ...form })
    router.replace('/')
  } catch (_) {
    loadCaptcha()
  } finally {
    loading.value = false
  }
}

loadCaptcha()
</script>

<style scoped>
.login { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: #f5f7fa; }
.box { width: 380px; }
.captcha-row { display: flex; gap: 8px; align-items: center; }
.captcha-img { height: 36px; border: 1px solid #dcdfe6; border-radius: 4px; cursor: pointer; }
</style>
