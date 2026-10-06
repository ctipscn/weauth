package cli

// CliArgs 保存解析后的命令行参数。
type CliArgs struct {
	// ConfigPath 配置文件路径（--config），空值表示不覆盖默认配置。
	ConfigPath string

	// Init 是否执行用户初始化（--init）。
	Init bool

	// Username 用户名（配合 --init 使用）。
	Username string

	// Password 密码（配合 --init 使用）。
	Password string

	// FrontendTrialURL 试用模式的前端 URL（--frontend-trial）。
	FrontendTrialURL string

	// Flags 通过 --section-key 参数指定的配置覆盖值。
	Flags ConfigOfArgs
}
