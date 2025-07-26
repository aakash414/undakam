FROM node:22-alpine AS builder
RUN npm install -g pnpm

WORKDIR /app/frontend
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY frontend ./
RUN pnpm run build
RUN pnpm run generate

FROM node:22-alpine
RUN npm install -g pnpm

WORKDIR /app
COPY --from=builder /app/frontend/.output ./

EXPOSE 3000
CMD ["node", "server/index.mjs"]
