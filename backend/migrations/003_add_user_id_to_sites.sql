-- 为 sites 表添加 user_id 字段以支持多用户隔离

-- 添加 user_id 字段（先允许为空）
ALTER TABLE sites ADD COLUMN IF NOT EXISTS user_id INT;

-- 为现有站点设置默认所有者（假设 ID 为 1 的用户，如果不存在则创建）
-- 注意：在实际迁移时，你可能需要手动处理现有数据
DO $$
BEGIN
    -- 如果有站点但没有 user_id，为它们分配一个默认用户
    IF EXISTS (SELECT 1 FROM sites WHERE user_id IS NULL) THEN
        -- 确保至少有一个管理员用户
        INSERT INTO admin_users (username, email, password_hash, role, enabled, created_at, updated_at)
        VALUES ('admin', 'admin@example.com', '', 'admin', true, NOW(), NOW())
        ON CONFLICT (email) DO NOTHING;
        
        -- 更新所有没有 user_id 的站点
        UPDATE sites 
        SET user_id = (SELECT id FROM admin_users ORDER BY id LIMIT 1)
        WHERE user_id IS NULL;
    END IF;
END $$;

-- 将字段设置为 NOT NULL
ALTER TABLE sites ALTER COLUMN user_id SET NOT NULL;

-- 添加外键约束
ALTER TABLE sites 
ADD CONSTRAINT fk_sites_owner 
FOREIGN KEY (user_id) REFERENCES admin_users(id) ON DELETE CASCADE;

-- 添加索引
CREATE INDEX IF NOT EXISTS idx_sites_user_id ON sites(user_id);

-- 添加复合索引以优化按用户查询
CREATE INDEX IF NOT EXISTS idx_sites_user_created ON sites(user_id, created_at DESC);
