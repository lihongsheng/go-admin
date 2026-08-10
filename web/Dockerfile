# ---- Build Stage ----
FROM node:20-alpine AS builder

WORKDIR /app

# 国内构建可取消下行注释，配置 NPM 镜像加速
# RUN npm config set registry https://registry.npmmirror.com

# Mock 模式开关：docker compose 构建 web-mock 时传入 true，
# 前端整体 mock 生效（免后端/免安装/免登录，见 src/api/mock/index.js）
ARG VITE_USE_MOCK=false

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

RUN VITE_USE_MOCK=$VITE_USE_MOCK npm run build:prod

# ---- Run Stage ----
FROM nginx:alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

ENTRYPOINT ["nginx", "-g", "daemon off;"]
