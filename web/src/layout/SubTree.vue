<template>
  <template v-if="hidden(item) || item.meta?.type === 'button'" />
  <el-sub-menu v-else-if="hasVisibleChildren(item)" :index="fullPath">
    <template #title>
      <el-icon v-if="item.meta?.icon"><component :is="iconMap[item.meta.icon] || MenuIcon" /></el-icon>
      <span>{{ item.meta?.title }}</span>
    </template>
    <SubTree v-for="c in visibleChildren(item)" :key="c.path || c.name" :item="c" :base="fullPath" />
  </el-sub-menu>
  <el-menu-item v-else :index="fullPath">
    <el-icon v-if="item.meta?.icon"><component :is="iconMap[item.meta.icon] || MenuIcon" /></el-icon>
    <span>{{ item.meta?.title }}</span>
  </el-menu-item>
</template>

<script setup>
import { computed } from 'vue'
import {
  DataBoard, Setting, User, UserFilled, Avatar, Grid, List, Link, Edit,
  Menu as MenuIcon, Plus, Search, Delete, Upload, Download, Star, StarFilled,
  Folder, Document, Files, DataAnalysis, PieChart, Histogram, TrendCharts,
  OfficeBuilding, HomeFilled, Location, Monitor, Cpu, Connection, SetUp,
  Operation, Guide, Tools, Key, Lock, Unlock, Warning, InfoFilled,
  CircleCheck, CircleClose, Promotion, Sell, Ticket, Collection,
  CollectionTag, PriceTag, ShoppingCart, Goods, Postcard, Picture, VideoCamera,
  Headset, Microphone, Switch, SwitchButton, FullScreen, Rank, Aim, Medal,
  Trophy, MagicStick, Brush, Coordinate, Box, TakeawayBox, Ship, Bicycle,
  Van, Clock, Timer, Calendar, Bell, Message, Share, Refresh
} from '@element-plus/icons-vue'

const props = defineProps({
  item: Object,
  base: { type: String, default: '' }
})

// 菜单名 → Element Plus 图标映射
const iconMap = {
  dashboard: DataBoard,
  setting: Setting,
  user: User,
  'user-filled': UserFilled,
  peoples: Avatar,
  'tree-table': Grid,
  component: List,
  api: Link,
  list: List,
  edit: Edit,
  plus: Plus,
  search: Search,
  delete: Delete,
  upload: Upload,
  download: Download,
  star: Star,
  'star-filled': StarFilled,
  folder: Folder,
  document: Document,
  files: Files,
  'data-analysis': DataAnalysis,
  'pie-chart': PieChart,
  histogram: Histogram,
  'trend-charts': TrendCharts,
  'office-building': OfficeBuilding,
  'home-filled': HomeFilled,
  location: Location,
  monitor: Monitor,
  cpu: Cpu,
  connection: Connection,
  'set-up': SetUp,
  operation: Operation,
  guide: Guide,
  tools: Tools,
  key: Key,
  lock: Lock,
  unlock: Unlock,
  warning: Warning,
  'info-filled': InfoFilled,
  'circle-check': CircleCheck,
  'circle-close': CircleClose,
  promotion: Promotion,
  sell: Sell,
  ticket: Ticket,
  collection: Collection,
  'collection-tag': CollectionTag,
  'price-tag': PriceTag,
  'shopping-cart': ShoppingCart,
  goods: Goods,
  postcard: Postcard,
  picture: Picture,
  'video-camera': VideoCamera,
  headset: Headset,
  microphone: Microphone,
  switch: Switch,
  'switch-button': SwitchButton,
  'full-screen': FullScreen,
  rank: Rank,
  aim: Aim,
  medal: Medal,
  trophy: Trophy,
  'magic-stick': MagicStick,
  brush: Brush,
  coordinate: Coordinate,
  box: Box,
  'takeaway-box': TakeawayBox,
  ship: Ship,
  bicycle: Bicycle,
  van: Van,
  clock: Clock,
  timer: Timer,
  calendar: Calendar,
  bell: Bell,
  message: Message,
  share: Share,
  refresh: Refresh,
  grid: Grid,
  menu: MenuIcon,
  link: Link,
  avatar: Avatar,
}

const fullPath = computed(() => {
  const p = props.item.path || ''
  if (p.startsWith('/')) return p
  return (props.base.endsWith('/') ? props.base : props.base + '/') + p
})

function hidden(i) {
  if (i.meta?.hidden) return true
  if (i.meta?.type === 'button') return true
  return false
}

function hasVisibleChildren(item) {
  if (!item.children || item.children.length === 0) return false
  return item.children.some(c => !hidden({ meta: c.meta }))
}

function visibleChildren(item) {
  if (!item.children) return []
  return item.children.filter(c => !hidden({ meta: c.meta }))
}
</script>
