import { usePermissionStore } from '@/store/modules/permission'

export function setupPermission(app) {
  app.directive('permission', {
    mounted(el, binding) {
      const need = binding.value
      if (!need) return
      const permStore = usePermissionStore()
      if (!permStore.btns || !permStore.btns.includes(need)) {
        el.parentNode && el.parentNode.removeChild(el)
      }
    }
  })
}
