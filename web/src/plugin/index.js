// 前端插件加载器：扫描 src/plugin/*/index.js，调用其默认导出 install(app, ctx)
const modules = import.meta.glob('./*/index.js')

export async function loadPlugins(app) {
  for (const path in modules) {
    try {
      const mod = await modules[path]()
      if (mod && typeof mod.default === 'function') {
        mod.default(app, { name: path.split('/')[1] })
      }
    } catch (e) {
      console.warn('[plugin] load failed:', path, e)
    }
  }
}
