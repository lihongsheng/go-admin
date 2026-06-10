// 示例前端插件：在 Vue 3 中注册全局组件或提供其他功能
export default function install(app, ctx) {
  console.log('[plugin] example loaded:', ctx.name)
}
