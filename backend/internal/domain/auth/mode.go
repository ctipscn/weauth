package auth

// LoginMode 用户认证登录模式枚举
type LoginMode int

const (
	LoginModeOIDC   LoginMode = iota // 0 - 仅 OIDC 登录
	LoginModePasswd                  // 1 - 仅密码登录
	LoginModeBoth                    // 2 - 支持 OIDC 和密码
)

// String 返回登录模式的字符串表示
func (m LoginMode) String() string {
	switch m {
	case LoginModeOIDC:
		return "oidc"
	case LoginModePasswd:
		return "passwd"
	case LoginModeBoth:
		return "both"
	default:
		return "passwd"
	}
}

// ParseLoginMode 将字符串解析为登录模式枚举
func ParseLoginMode(s string) (LoginMode, error) {
	switch s {
	case "oidc":
		return LoginModeOIDC, nil
	case "passwd":
		return LoginModePasswd, nil
	case "both":
		return LoginModeBoth, nil
	default:
		return LoginModePasswd, nil
	}
}