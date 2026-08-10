# 前端体验增强

顶部多标签页与右侧抽屉表单两项布局级增强。

## 顶部多标签页（TagsView）

浏览器式多标签导航，位于头部下方、内容区上方。

### 行为

| 行为 | 说明 |
|---|---|
| 首页常驻 | 第一个 tab 固定为「首页」（`/dashboard`），不可关闭 |
| 自动追加 | 打开任意菜单自动新增 tab（叶子页面，目录容器不生成） |
| 点击切换 | 点击 tab 路由跳转，激活 tab 高亮 |
| 关闭 | 非首页 tab 带关闭按钮；关闭当前页时自动跳转左侧相邻 tab |
| 关闭其他 / 关闭全部 | 标签栏右侧下拉菜单 |
| 横向滚动 | tab 过多时 el-scrollbar 横向滚动 |
| 刷新恢复 | 刷新后当前页 tab 自动恢复 |
| 登出清理 | 切换账号时清空历史 tab |

### 实现

| 文件 | 职责 |
|---|---|
| `web/src/store/modules/tags.js` | Pinia store：`addTab / removeTab / closeOthers / closeAll`，首页 tab 常驻（affix） |
| `web/src/layout/TagsView.vue` | 标签栏组件（滚动条 + 关闭按钮 + 下拉操作 + 深色模式） |
| `web/src/layout/index.vue` | 布局接入 + 登出时 `tagsStore.closeAll()` |

```js
// 路由变化自动追加（immediate 保证刷新后恢复当前页）
watch(() => route.path, () => tagsStore.addTab(route), { immediate: true })

// 关闭 tab：返回应跳转的相邻路径
function handleClose(tag) {
  const nextPath = tagsStore.removeTab(tag.path)
  if (tag.path === route.path && nextPath) router.push(nextPath)
}
```

## 右侧抽屉表单

系统管理页面与资源管理页面的表单统一从居中弹窗改为**右侧抽屉**（el-drawer，direction 默认 rtl）。

| 页面 | 抽屉内容 | 宽度 |
|---|---|---|
| 角色管理 | 新增/编辑角色、角色授权（菜单树） | 480px / 720px |
| 用户管理 | 新增/编辑用户（含头像上传） | 520px |
| 菜单管理 | 新增/编辑菜单（含 API 规则编辑器） | 620px |
| 资源管理 | 新增/编辑/查看资源 | 480px |

统一样式约定：

```html
<el-drawer v-model="dlg" :title="..." size="480px" destroy-on-close>
  <el-form :model="form" label-width="80px">
    <!-- 表单字段 -->
  </el-form>
  <template #footer>
    <div class="drawer-footer">
      <el-button @click="dlg = false">取消</el-button>
      <el-button type="primary" @click="submit" :loading="submitting">确定</el-button>
    </div>
  </template>
</el-drawer>
```

```css
.drawer-footer { display: flex; justify-content: flex-end; gap: 12px; }
```

特点：
- `destroy-on-close`：关闭销毁表单 DOM，每次打开重新初始化
- footer 插槽统一右对齐操作区（取消 / 确定）
- 支持深色模式
