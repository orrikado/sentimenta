# Базовый образ для сборки
FROM node:24-slim AS builder

WORKDIR /app
COPY package*.json pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm fetch

# Установка зависимостей с кэшированием
RUN pnpm install --frozen-lockfile

# Копирование исходников и сборка
COPY . .
RUN pnpm build

# Финальный образ с минимальным Node.js окружением
FROM node:24-alpine

WORKDIR /app

# Копирование необходимых файлов из builder
COPY --from=builder /app/build /app/
COPY --from=builder /app/package.json /app/ 
COPY --from=builder /app/pnpm-lock.yaml /app

# Установка production-зависимостей
RUN npm install -g pnpm && \
    pnpm install --prod --frozen-lockfile

# Экспорт порта и запуск приложения
EXPOSE 3000
CMD ["node", "/app"]