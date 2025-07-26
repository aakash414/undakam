FROM node:22-alpine
RUN npm install -g pnpm
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY . .
RUN pnpm run build
RUN pnpm run generate

FROM node:22-alpine
RUN npm install -g pnpm
WORKDIR /app
COPY --from=0 /app/.output ./
EXPOSE 3000
CMD ["node", "server/index.mjs"]