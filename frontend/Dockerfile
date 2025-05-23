FROM node:24-slim AS builder

WORKDIR /app
COPY . /app

RUN npm install -g pnpm

RUN pnpm install && pnpm build

FROM alpine:3.21 AS exporter
WORKDIR /export
COPY --from=builder /app/build ./dist
