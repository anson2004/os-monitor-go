# ---- build stage ----
FROM node:22-alpine AS build
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# ---- runtime stage ----
FROM node:22-alpine
WORKDIR /app
ENV NODE_ENV=production \
    NITRO_PORT=3000 \
    NITRO_HOST=0.0.0.0 \
    NUXT_API_URL=http://osmonitor:8080

COPY --from=build /app/.output ./.output

USER node
EXPOSE 3000
ENTRYPOINT ["node", ".output/server/index.mjs"]
