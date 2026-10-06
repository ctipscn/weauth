# WeAuth

WeAuth（微验）是一个基于 SHA-256 工作量证明（Proof-of-Work）的人机验证与反滥用平台。它提供类验证码的用户体验，但无需图片识别或滑块拼图：浏览器执行一次轻量级 PoW 计算，后端验证通过后发放一次性短生命周期 Token。

**体验地址：** https://www.weauth.cn

## 特性

- **工作量证明验证** — 基于 SHA-256 的 PoW 挑战，可动态调节难度，浏览器端 Web Worker 求解。
- **嵌入型 Widget** — 通过 `weauth.js` UMD 脚本或 iframe 嵌入任意第三方站点，支持 `<script>` 标签、JS API 以及 HTML 属性自动渲染三种接入方式。
- **批量验证模式** — 一次发起多个 PoW 挑战，批量求解后一次性提交验证，适合高频场景。
- **域名白名单** — 站点级别域名白名单，支持子域名通配符，防止跨站滥用。
- **智能风控** — 基于请求频次和失败次数的 IP 自适应难度调节，可配置每站点的 IP 策略与惩罚权重。
- **双认证模式** — 支持 OIDC 单点登录与本地密码认证，可混合使用。
- **管理控制台** — 站点管理、密钥轮换、验证统计、域名管理、IP 策略配置。
- **管理后台** — 管理员可管理全部用户与站点。

## 技术栈

| 层       | 技术                                                                 |
|----------|----------------------------------------------------------------------|
| 后端     | Go 1.24, Gin, GORM, PostgreSQL, JWT (golang-jwt/v5)                 |
| 前端 SPA | Vue 3 (Composition API), TypeScript, Pinia, TanStack Vue Query, Vue Router 5, Tailwind CSS |
| Widget   | Vue 3 SFC + Vite UMD 构建，Web Worker PoW 求解                     |
| 部署     | Docker Compose (PostgreSQL + Go 二进制 + Nginx 静态托管)            |

## 项目布局

```text
.
├── backend/                  # Go API 服务
│   ├── cmd/weauth/main.go   # 入口
│   ├── internal/
│   │   ├── cli/              # 命令行参数解析
│   │   ├── config/           # 配置加载（INI + 环境变量 + CLI 覆盖）
│   │   ├── domain/           # 领域层（实体、仓储接口、查询）
│   │   ├── application/      # 应用服务（PoW、认证、站点、OIDC、用户、管理）
│   │   ├── infrastructure/   # 基础设施（数据库、日志、安全、健康检查、Gin 服务器）
│   │   └── interface/        # HTTP 接口（路由、处理器、中间件）
│   ├── configs/              # 配置文件示例
│   ├── docker/               # 多阶段 Dockerfile
│   ├── migrations/           # 数据库初始化 SQL
│   └── scripts/              # 本地开发脚本
├── frontend/                 # Vue 3 前端 SPA + Widget
│   ├── src/
│   │   ├── lib/              # API 客户端、认证存储、工具函数
│   │   ├── pages/            # 页面组件（公开页、认证页、控制台、管理后台）
│   │   └── widget/           # 嵌入型验证 Widget（Vue SFC + UMD 构建）
│   │   └── router/           # Vue Router 路由定义（含导航守卫）
│   └── vite.config.*.ts      # SPA 构建 / Widget UMD 构建
└── compose.yml               # 全栈 Docker Compose 编排
```

## 快速开始

### 开发环境

**前置条件：** Go 1.24+, Node.js 22+, pnpm, PostgreSQL 14+

**1. 启动数据库**

```bash
cd backend && bash scripts/prepare.sh
```

**2. 启动后端**

```bash
cd backend
cp configs/weauth.example.conf configs/weauth.conf
go run ./cmd/weauth --config configs/weauth.conf
```

首次运行自动执行数据库迁移，并创建管理员用户（`admin` / `admin123`）和试用站点。

**3. 启动前端**

```bash
cd frontend
pnpm install
pnpm dev          # SPA + Widget 同时监听
```

前端开发服务器默认运行在 `http://localhost:5173`，API 请求通过 Vite proxy 转发到后端 `http://localhost:8080`。

### Docker Compose

```bash
docker compose up -d
```

启动后：
- 前端：http://localhost:5173
- 后端 API：http://localhost:8080
- 健康检查：http://localhost:8081/health

## 配置

配置按以下优先级加载（后项覆盖前项）：

1. 硬编码默认值
2. `/etc/weauth/weauth.conf`
3. `$HOME/.config/weauth/weauth.conf`
4. `$WEAUTH_HOME/conf/weauth.conf`
5. `--config` 命令行参数指定的文件
6. `WEAUTH_*` 环境变量
7. `--server-port`, `--db-password` 等命令行标志

### 主要配置项

| 配置节        | 环境变量前缀            | 说明                        |
|---------------|------------------------|-----------------------------|
| `[server]`    | `WEAUTH_SERVER_*`      | HTTP 监听地址、端口、运行模式  |
| `[database]`  | `WEAUTH_DATABASE_*`    | PostgreSQL 连接与连接池       |
| `[jwt]`       | `WEAUTH_JWT_*`         | JWT 签名密钥与过期时间         |
| `[oidc]`      | `WEAUTH_OIDC_*`        | OIDC 提供商集成               |
| `[cors]`      | `WEAUTH_CORS_*`        | 跨域配置                     |
| `[ratelimit]` | `WEAUTH_RATELIMIT_*`   | 站点级请求频率限制            |

**重要：** `JWT_SECRET` 必须在生产环境中设置为强随机值；`OIDC_CLIENT_SECRET` 绝不能提交到版本控制系统。

## Widget 接入

WeAuth 提供三种接入方式，无需安装任何 SDK。

### 方式一：HTML 属性自动渲染（推荐）

```html
<script src="https://www.weauth.cn/weauth.js" async defer></script>

<div class="weauth-widget"
     data-sitekey="your-sitekey"
     data-action="login"
     data-theme="light"
     data-callback="onVerify"></div>

<script>
  function onVerify(token) {
    // token 用于后端验证
  }
</script>
```

### 方式二：JavaScript API

```js
const widgetId = weauth.render('#my-widget', {
  sitekey: 'your-sitekey',
  action: 'login',
  callback: (token) => { /* ... */ },
})

// 重置 / 手动触发 / 移除
weauth.reset(widgetId)
weauth.execute(widgetId)
weauth.remove(widgetId)
```

### 方式三：iframe 直嵌

```html
<iframe id="weauth-iframe"
        src="https://www.weauth.cn/widget.html?sitekey=your-key&widgetId=my-widget"
        width="300" height="65" frameborder="0"></iframe>
```

### 后端验证

```go
// Go 示例
resp, err := http.PostForm("https://api.weauth.cn/api/pow/siteverify", url.Values{
    "secret": {"your-secret"},
    "token":  {token},
})
```

详细接入文档见前端 `/docs` 页面。

## 工作流程

```text
用户浏览器                   WeAuth 后端                   第三方应用
    │                           │                            │
    │  POST /api/pow/challenge  │                            │
    │─── sitekey + origin ──────→│                            │
    │                           │ 校验域名、计算难度、创建挑战 │
    │←──── challenge + difficulty ───│                       │
    │                           │                            │
    │  Web Worker 求解 PoW     │                            │
    │  SHA-256(challenge:nonce) │                            │
    │  满足 difficulty 个前导零  │                            │
    │                           │                            │
    │  POST /api/pow/verify     │                            │
    │─── challengeId + nonce ──→│                            │
    │                           │ 验证哈希、发放 Token        │
    │←──────── token ─────────────│                          │
    │                           │                            │
    │  提交 token + 表单        │                            │
    │─────────────────────────────────────────→              │
    │                           │  POST /api/pow/siteverify  │
    │                           │←──────── token ────────────│
    │                           │──── { success: true } ────→│
```

## API 概览

| 分组         | 路径                             | 说明                | 认证     |
|-------------|----------------------------------|--------------------|----------|
| PoW 公开     | `POST /api/pow/challenge`        | 获取 PoW 挑战       | 无       |
| PoW 公开     | `POST /api/pow/verify`           | 验证 PoW 解        | 无       |
| PoW 公开     | `POST /api/pow/siteverify`       | 第三方站点 Token 验证 | Secret   |
| 认证         | `POST /api/auth/register`        | 注册                | 无       |
| 认证         | `POST /api/auth/login`           | 登录                | 无       |
| 认证         | `GET /api/auth/oidc/authorize`   | OIDC 授权跳转       | 无       |
| 认证         | `GET /api/auth/oidc/callback`    | OIDC 回调           | 无       |
| 认证         | `GET /api/auth/me`               | 当前用户信息         | JWT      |
| 用户         | `GET /api/user/profile`          | 用户主页            | JWT      |
| 控制台       | `GET /api/console/sites`         | 站点列表            | JWT      |
| 控制台       | `POST /api/console/sites`        | 创建站点            | JWT      |
| 控制台       | `GET /api/console/sites/:id`     | 站点详情            | JWT      |
| 控制台       | `GET /api/console/sites/:id/stats`| 验证统计            | JWT      |
| 控制台       | `GET /api/console/sites/:id/ip-policy` | IP 策略       | JWT      |
| 管理后台     | `GET /api/admin/users`           | 用户列表            | Admin    |
| 管理后台     | `GET /api/admin/sites`           | 全部站点列表         | Admin    |

## License

Apache-2.0