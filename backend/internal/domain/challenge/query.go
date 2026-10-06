package challenge

// ChallengeQuery 挑战查询过滤条件
type ChallengeQuery struct {
	SiteID   uint
	Solved   *bool
	ClientIP string
}