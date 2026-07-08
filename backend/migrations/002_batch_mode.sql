-- 批量 PoW 验证模式迁移脚本
-- 添加批量模式相关字段

-- ========== sites 表添加批量模式配置 ==========

-- 是否启用批量模式（新站点默认启用）
ALTER TABLE sites ADD COLUMN IF NOT EXISTS batch_mode_enabled BOOLEAN DEFAULT true;

-- 批量模式固定难度
ALTER TABLE sites ADD COLUMN IF NOT EXISTS batch_difficulty INT DEFAULT 4;

-- 最小批次数
ALTER TABLE sites ADD COLUMN IF NOT EXISTS min_batch_count INT DEFAULT 1;

-- 最大批次数
ALTER TABLE sites ADD COLUMN IF NOT EXISTS max_batch_count INT DEFAULT 10;

-- ========== challenges 表添加批量模式字段 ==========

-- 批次ID，同批次的 challenge 共享
ALTER TABLE challenges ADD COLUMN IF NOT EXISTS batch_id UUID;

-- 在批次中的序号
ALTER TABLE challenges ADD COLUMN IF NOT EXISTS batch_index INT DEFAULT 0;

-- 批次总数
ALTER TABLE challenges ADD COLUMN IF NOT EXISTS batch_total INT DEFAULT 1;

-- ========== ip_policies 表添加批量模式策略字段 ==========

-- 批次增量（批量模式专用）
ALTER TABLE ip_policies ADD COLUMN IF NOT EXISTS batch_increment INT DEFAULT 1;

-- ========== 索引 ==========

-- 批次ID索引，用于批量验证时快速查询同批次的 challenges
CREATE INDEX IF NOT EXISTS idx_challenges_batch_id ON challenges(batch_id);

-- ========== 注释 ==========
COMMENT ON COLUMN sites.batch_mode_enabled IS '是否启用批量模式：固定难度，动态数量';
COMMENT ON COLUMN sites.batch_difficulty IS '批量模式的固定难度（建议4或5）';
COMMENT ON COLUMN sites.min_batch_count IS '批量模式最小批次数';
COMMENT ON COLUMN sites.max_batch_count IS '批量模式最大批次数';
COMMENT ON COLUMN challenges.batch_id IS '批次ID，同批次的challenges共享此ID';
COMMENT ON COLUMN challenges.batch_index IS '在批次中的序号（从0开始）';
COMMENT ON COLUMN challenges.batch_total IS '批次总数';
COMMENT ON COLUMN ip_policies.batch_increment IS '批量模式：每超过阈值增加的批次数';
