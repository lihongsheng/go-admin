# ---- Build Stage ----
FROM node:20-alpine AS builder

WORKDIR /app

# 国内构建可取消下行注释，配置 NPM 镜像加速
# RUN npm config set registry https://registry.npmmirror.com

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

RUN npm run build:prod

# ---- Run Stage ----
FROM nginx:alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

ENTRYPOINT ["nginx", "-g", "daemon off;"]
