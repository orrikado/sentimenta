# --- Stage 1: Build Go backend ---
FROM golang:1.24.3-alpine AS backend-builder

WORKDIR /app
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download
COPY ./backend ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# --- Stage 2: Build frontend ---
FROM node:24-slim AS frontend-builder

WORKDIR /app
COPY ./frontend ./
RUN npm install -g pnpm
RUN pnpm install && pnpm build

# --- Stage 3: Final image with Nginx ---
FROM alpine:3.21

# Установим nginx и tzdata
RUN apk add --no-cache nginx tzdata

# Создадим нужные директории
RUN mkdir -p /var/www/html /etc/nginx/conf.d /app

# Копируем backend бинарник
COPY --from=backend-builder /app/main /app/main

# Копируем фронтенд (dist)
COPY --from=frontend-builder /app/build /var/www/html

# Копируем nginx конфиг
COPY ./nginx/nginx.conf /etc/nginx/nginx.conf
COPY ./nginx/default.conf /etc/nginx/conf.d/default.conf

# Подготовка запуска nginx и backend
EXPOSE 80
EXPOSE 8000

CMD sh -c "/app/main & nginx"
