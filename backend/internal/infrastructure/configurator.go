package infrastructure

import (
	"errors"
	"fmt"
	"strings"
)

// Mode 定义应用统一的状态模式
type Mode string

const (
	DebugMode   Mode = "debug"
	ReleaseMode Mode = "release"
	TestMode    Mode = "test"
)

// 预定义错误信息
var (
	ErrInvalidMode = errors.New("invalid application mode")
)

// String 实现 fmt.Stringer 接口
func (m Mode) String() string {
	return string(m)
}

// IsValid 检查当前 Mode 是否合法
func (m Mode) IsValid() bool {
	switch m {
	case DebugMode, ReleaseMode, TestMode:
		return true
	default:
		return false
	}
}

// Configurator 任何需要响应 Mode 变更的组件/框架都需要实现该接口
type Configurator interface {
	// ApplyMode 根据传入的 Mode 设置组件自身的内部状态
	ApplyMode(mode Mode) error
}

// StringToMode 将字符串解析为对应的 Mode 结构
func StringToMode(m string) (Mode, error) {
	mode := Mode(strings.ToLower(strings.TrimSpace(m)))
	if mode.IsValid() {
		return mode, nil
	}
	return "", fmt.Errorf("%w: %q (valid modes: %s, %s, %s)", ErrInvalidMode, m, DebugMode, ReleaseMode, TestMode)
}
