package iprecord

import "time"

// IpRecord IP 访问记录实体
type IpRecord struct {
	ID                uint      `json:"id"`
	SiteID            uint      `json:"site_id"`
	ClientIP          string    `json:"client_ip"`
	RequestCount      int       `json:"request_count"`
	FailureCount      int       `json:"failure_count"`
	CurrentDifficulty int       `json:"current_difficulty"`
	WindowStart       time.Time `json:"window_start"`
	LastRequestAt     time.Time `json:"last_request_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// IsWindowExpired 检查时间窗口是否过期
func (r *IpRecord) IsWindowExpired(windowSeconds int) bool {
	return time.Now().After(r.WindowStart.Add(time.Duration(windowSeconds) * time.Second))
}

// ResetWindow 重置统计窗口
func (r *IpRecord) ResetWindow() {
	r.WindowStart = time.Now()
	r.RequestCount = 0
	r.FailureCount = 0
}

// IncrementRequest 增加请求计数
func (r *IpRecord) IncrementRequest() {
	r.RequestCount++
	r.LastRequestAt = time.Now()
}

// IncrementFailure 增加失败计数
func (r *IpRecord) IncrementFailure() {
	r.FailureCount++
	r.LastRequestAt = time.Now()
}

// NewIpRecord 创建新的 IP 记录
func NewIpRecord(siteID uint, clientIP string) *IpRecord {
	return &IpRecord{
		SiteID:      siteID,
		ClientIP:    clientIP,
		WindowStart: time.Now(),
	}
}