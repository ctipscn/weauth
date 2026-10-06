-- WeAuth 数据库初始化脚本
-- 基于 internal/infrastructure/database/models 自动生成
-- 请使用具有创建表权限的用户执行此脚本

-- ========== 用户管理 ==========

-- 用户表（对应 models.User）
CREATE TABLE IF NOT EXISTS users (
    id          BIGSERIAL    PRIMARY KEY,
    username    VARCHAR(100) DEFAULT '',
    email       VARCHAR(255) DEFAULT '',
    avatar      VARCHAR(500) DEFAULT '',
    role        VARCHAR(20)  DEFAULT 'staff',
    enabled     BOOLEAN      DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- 认证凭证表（对应 models.Auth，与 users 分离，无 FK 约束，由业务逻辑维护一致性）
CREATE TABLE IF NOT EXISTS auth (
    id            BIGSERIAL    PRIMARY KEY,
    user_id       BIGINT       NOT NULL UNIQUE,
    oidc_sub      VARCHAR(255),
    password_hash VARCHAR(255) DEFAULT '',
    login_mode    INT          NOT NULL DEFAULT 1,
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_oidc_sub ON auth(oidc_sub);
CREATE INDEX IF NOT EXISTS idx_auth_user_id ON auth(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_deleted_at ON auth(deleted_at);

-- ========== 站点管理 ==========

-- 站点表（对应 models.Site）
CREATE TABLE IF NOT EXISTS sites (
    id                 BIGSERIAL    PRIMARY KEY,
    user_id            BIGINT       NOT NULL DEFAULT 0,
    name               VARCHAR(100) NOT NULL,
    sitekey            VARCHAR(32)  NOT NULL,
    secret             VARCHAR(64)  NOT NULL,
    base_difficulty    INT          NOT NULL DEFAULT 4,
    max_difficulty     INT          NOT NULL DEFAULT 8,
    challenge_timeout  INT          DEFAULT 120,
    token_timeout      INT          DEFAULT 60,
    enabled            BOOLEAN      DEFAULT TRUE,
    batch_mode_enabled BOOLEAN      DEFAULT FALSE,
    batch_difficulty   INT          DEFAULT 4,
    min_batch_count    INT          DEFAULT 1,
    max_batch_count    INT          DEFAULT 10,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at         TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_sites_sitekey ON sites(sitekey);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sites_secret  ON sites(secret);
CREATE INDEX IF NOT EXISTS idx_sites_deleted_at ON sites(deleted_at);
CREATE INDEX IF NOT EXISTS idx_sites_user_id ON sites(user_id);
CREATE INDEX IF NOT EXISTS idx_sites_user_created ON sites(user_id, created_at DESC);

-- 站点域名白名单表（对应 models.Domain）
-- 注意：GORM 模式名为 "Domain"，表名为 "domains"
CREATE TABLE IF NOT EXISTS domains (
    id         BIGSERIAL    PRIMARY KEY,
    site_id    BIGINT       NOT NULL,
    domain     VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_domains_site_domain UNIQUE (site_id, domain)
);

CREATE INDEX IF NOT EXISTS idx_domains_site_id ON domains(site_id);

-- ========== IP 策略 ==========

-- IP 策略配置表（对应 models.IpPolicy）
-- site_id 为 NULL 表示全局策略
CREATE TABLE IF NOT EXISTS ip_policies (
    id                   BIGSERIAL    PRIMARY KEY,
    site_id              BIGINT       UNIQUE,
    time_window_seconds  INT          DEFAULT 300,
    request_threshold    INT          DEFAULT 10,
    difficulty_increment INT          DEFAULT 1,
    batch_increment      INT          DEFAULT 1,
    failure_weight       DECIMAL(3,2) DEFAULT 2.0,
    enabled              BOOLEAN      DEFAULT TRUE,
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    -- 外键约束已移除，业务层保证一致性
);

-- IP 访问记录表（对应 models.IpRecord）
CREATE TABLE IF NOT EXISTS ip_records (
    id                 BIGSERIAL    PRIMARY KEY,
    site_id            BIGINT       NOT NULL,
    client_ip          VARCHAR(45)  NOT NULL,
    request_count      INT          DEFAULT 0,
    failure_count      INT          DEFAULT 0,
    current_difficulty INT          DEFAULT 4,
    window_start       TIMESTAMPTZ  NOT NULL,
    last_request_at    TIMESTAMPTZ,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_ip_records_site_ip UNIQUE (site_id, client_ip)
);

-- ========== Challenge & Token ==========

-- Challenge 记录表（对应 models.Challenge）
CREATE TABLE IF NOT EXISTS challenges (
    id          UUID         PRIMARY KEY,
    site_id     BIGINT       NOT NULL,
    challenge   VARCHAR(128) NOT NULL,
    difficulty  INT          NOT NULL,
    client_ip   VARCHAR(45)  DEFAULT '',
    origin      VARCHAR(255) DEFAULT '',
    action      VARCHAR(100) DEFAULT '',
    solved      BOOLEAN      DEFAULT FALSE,
    token       VARCHAR(128) DEFAULT '',
    token_used  BOOLEAN      DEFAULT FALSE,
    expires_at  TIMESTAMPTZ  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    solved_at   TIMESTAMPTZ,
    batch_id    UUID,
    batch_index INT          DEFAULT 0,
    batch_total INT          DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_challenges_site_id ON challenges(site_id);
CREATE INDEX IF NOT EXISTS idx_challenges_token ON challenges(token);
CREATE INDEX IF NOT EXISTS idx_challenges_batch_id ON challenges(batch_id);

-- ========== 验证日志 ==========

-- 验证日志表（对应 models.VerificationLog）
-- SiteID/ChallengeID 可为 NULL：日志仅作审计用途，不阻塞操作
CREATE TABLE IF NOT EXISTS verification_logs (
    id           BIGSERIAL    PRIMARY KEY,
    challenge_id UUID,
    site_id      BIGINT,
    client_ip    VARCHAR(45)  DEFAULT '',
    origin       VARCHAR(255) DEFAULT '',
    action       VARCHAR(100) DEFAULT '',
    iterations   BIGINT       DEFAULT 0,
    solve_time_ms INT         DEFAULT 0,
    success      BOOLEAN      DEFAULT FALSE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_verification_logs_site_id ON verification_logs(site_id);
CREATE INDEX IF NOT EXISTS idx_verification_logs_challenge_id ON verification_logs(challenge_id);