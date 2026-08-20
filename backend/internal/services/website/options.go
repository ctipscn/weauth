package website

// CreateSiteOptions 创建站点的可选配置
type CreateSiteOptions struct {
	BatchModeEnabled *bool
	BatchDifficulty  *int
	MinBatchCount    *int
	MaxBatchCount    *int
}