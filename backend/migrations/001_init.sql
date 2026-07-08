-- PoW Verify 数据库初始化脚本
-- 请使用具有创建表权限的用户执行此脚本

-- ========== 站点管理 ==========

-- 站点表
CREATE TABLE IF NOT EXISTS sites (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    sitekey VARCHAR(32) UNIQUE NOT NULL,
    secret VARCHAR(64) UNIQUE NOT NULL,
    base_difficulty INT NOT NULL DEFAULT 4,
    max_difficulty INT NOT NULL DEFAULT 8,
    challenge_timeout INT DEFAULT 120,
    token_timeout INT DEFAULT 300,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 站点域名白名单
CREATE TABLE IF NOT EXISTS site_domains (
    id SERIAL PRIMARY KEY,
    site_id INT NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    domain VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(site_id, domain)
);

-- ========== IP 策略 ==========

-- IP 策略配置表
CREATE TABLE IF NOT EXISTS ip_policies (
    id SERIAL PRIMARY KEY,
    site_id INT REFERENCES sites(id) ON DELETE CASCADE,
    time_window_seconds INT DEFAULT 300,
    request_threshold INT DEFAULT 10,
    difficulty_increment INT DEFAULT 1,
    failure_weight DECIMAL(3,2) DEFAULT 2.0,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(site_id)
);

-- IP 访问记录表
CREATE TABLE IF NOT EXISTS ip_records (
    id SERIAL PRIMARY KEY,
    site_id INT NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    client_ip VARCHAR(45) NOT NULL,
    request_count INT DEFAULT 0,
    failure_count INT DEFAULT 0,
    current_difficulty INT DEFAULT 4,
    window_start TIMESTAMP NOT NULL,
    last_request_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(site_id, client_ip)
);

-- ========== Challenge & Token ==========

-- Challenge 记录表
CREATE TABLE IF NOT EXISTS challenges (
    id UUID PRIMARY KEY,
    site_id INT NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    challenge VARCHAR(64) NOT NULL,
    difficulty INT NOT NULL,
    client_ip VARCHAR(45),
    origin VARCHAR(255),
    action VARCHAR(100),
    solved BOOLEAN DEFAULT false,
    token VARCHAR(64),
    token_used BOOLEAN DEFAULT false,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    solved_at TIMESTAMP
);

-- 验证日志表
CREATE TABLE IF NOT EXISTS verification_logs (
    id SERIAL PRIMARY KEY,
    challenge_id UUID REFERENCES challenges(id),
    site_id INT REFERENCES sites(id),
    client_ip VARCHAR(45),
    origin VARCHAR(255),
    action VARCHAR(100),
    iterations BIGINT,
    solve_time_ms INT,
    success BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ========== 管理员用户 ==========

-- 管理员用户表
CREATE TABLE IF NOT EXISTS admin_users (
    id SERIAL PRIMARY KEY,
    oidc_sub VARCHAR(255) UNIQUE,
    username VARCHAR(100),
    email VARCHAR(255),
    avatar VARCHAR(500),
    password_hash VARCHAR(255),
    role VARCHAR(20) DEFAULT 'user',
    enabled BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- ========== 索引 ==========
CREATE INDEX IF NOT EXISTS idx_sites_sitekey ON sites(sitekey);
CREATE INDEX IF NOT EXISTS idx_sites_secret ON sites(secret);
CREATE INDEX IF NOT EXISTS idx_sites_deleted_at ON sites(deleted_at);
CREATE INDEX IF NOT EXISTS idx_site_domains_site_id ON site_domains(site_id);
CREATE INDEX IF NOT EXISTS idx_ip_records_site_ip ON ip_records(site_id, client_ip);
CREATE INDEX IF NOT EXISTS idx_challenges_site_id ON challenges(site_id);
CREATE INDEX IF NOT EXISTS idx_challenges_token ON challenges(token);
CREATE INDEX IF NOT EXISTS idx_verification_logs_site_id ON verification_logs(site_id);
CREATE INDEX IF NOT EXISTS idx_admin_users_oidc_sub ON admin_users(oidc_sub);
CREATE INDEX IF NOT EXISTS idx_admin_users_deleted_at ON admin_users(deleted_at);

-- 授权给 pow_dev 用户 (如果需要)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO pow_dev;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO pow_dev;
