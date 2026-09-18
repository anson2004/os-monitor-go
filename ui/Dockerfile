# ---- build stage ----
FROM node:22-alpine AS build
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run generate

# ---- runtime stage ----
FROM nginx:1.27-alpine-slim

# Default upstream; override with -e API_URL=... (docker-compose sets it).
ENV API_URL=http://osmonitor:8080

COPY nginx/default.conf.template /etc/nginx/templates/default.conf.template
COPY --from=build /app/.output/public /usr/share/nginx/html

EXPOSE 80
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://localhost/health || exit 1
