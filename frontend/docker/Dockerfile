# =============================================================================
# Stage 1: Build
# =============================================================================
FROM node:22-alpine AS builder

# corepack enables pnpm without explicit install
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app

# Copy package manifests first for layer caching
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./

# Install dependencies (frozen lockfile ensures reproducibility)
RUN pnpm install --frozen-lockfile

# Copy the rest of the source
COPY . .

# Build environment variables (override via --build-arg)
ARG VITE_API_BASE_URL=http://localhost:8080/api
ARG VITE_WIDGET_BASE_URL=http://localhost:5173

ENV VITE_API_BASE_URL=${VITE_API_BASE_URL} \
    VITE_WIDGET_BASE_URL=${VITE_WIDGET_BASE_URL}

# Build both the frontend SPA and the widget library
RUN pnpm build

# =============================================================================
# Stage 2: Serve (nginx)
# =============================================================================
FROM nginx:stable-alpine AS production

# Copy custom nginx config
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf

# Copy built artifacts from the builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 5173

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost/ || exit 1

CMD ["nginx", "-g", "daemon off;"]