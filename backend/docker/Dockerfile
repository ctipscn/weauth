# =============================================================================
# WeAuth Backend — Multi-stage Dockerfile (Legacy / Non-BuildKit Compatible)
# =============================================================================

# ---------------------------------------------------------------------------
# Stage 1 — Build the Go binary
# ---------------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

WORKDIR /src

# 1. 优先复制 Go modules 文件利用标准 Docker 层缓存
COPY go.mod go.sum ./
RUN go mod download

# 2. 复制源码并编译
COPY . .
RUN CGO_ENABLED=0 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /usr/local/bin/weauth \
    ./cmd/weauth

# ---------------------------------------------------------------------------
# Stage 2 — Minimal runtime image (Minimal Alpine)
# ---------------------------------------------------------------------------
FROM alpine:3.21

WORKDIR /app

# 安装运行基础组件
RUN apk add --no-cache ca-certificates tzdata

# 创建非 root 用户与必要目录
RUN addgroup -S -g 10001 weauth && \
    adduser -S -u 10001 -G weauth weauth && \
    mkdir -p /app/logs /etc/weauth && \
    chown -R weauth:weauth /app/logs /etc/weauth

# 复制编译产物与配置文件
COPY --from=builder /usr/local/bin/weauth /app/weauth
COPY --chown=weauth:weauth configs/weauth.example.conf /etc/weauth/weauth.conf

# 切换非特权用户
USER weauth:weauth

# 声明端口
EXPOSE 8080 8081

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8081/health || exit 1

ENTRYPOINT [ \
    "/bin/sh", \
    "-c", \
    "/app/weauth --init \
        --username admin \
        --password admin123 \
    && /app/weauth --frontend-trial http://localhost:5173 \
    && /app/weauth" \
]
