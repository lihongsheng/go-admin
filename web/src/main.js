import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

import App from './App.vue'
import router from './router'
import pinia from './store'
import { setupPermission } from './directive/permission'
import { loadPlugins } from './plugin'

// 路由守卫
import './permission'

async function bootstrap() {
  const app = createApp(App)

  app.use(router)
  app.use(pinia)
  app.use(ElementPlus, { locale: zhCn })

  // 注册 v-permission 指令
  setupPermission(app)

  // 加载前端插件（传入 app 实例）
  await loadPlugins(app)

  app.mount('#app')
}

bootstrap()
