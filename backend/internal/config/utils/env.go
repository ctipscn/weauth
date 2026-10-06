package utils

import (
	"os"
	"strconv"
	"strings"
)

// EnvString 如果环境变量 envKey 存在且非空，则写入 target。
func EnvString(envKey string, target *string) {
	if v := os.Getenv(envKey); v != "" {
		*target = v
	}
}

// EnvInt 如果环境变量 envKey 存在且非空，则解析 int 后写入 target。
func EnvInt(envKey string, target *int) {
	if v := os.Getenv(envKey); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			*target = n
		}
	}
}

// EnvBool 从环境变量解析布尔值。
// 支持 true/false/1/0/yes/no。
func EnvBool(envKey string, target *bool) {
	if v := os.Getenv(envKey); v != "" {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "true", "1", "yes":
			*target = true
		case "false", "0", "no":
			*target = false
		}
	}
}

// EnvOriginList 从环境变量解析逗号分隔的字符串列表。
func EnvOriginList(envKey string, target *[]string) {
	if v := os.Getenv(envKey); v != "" {
		parts := strings.Split(v, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				out = append(out, trimmed)
			}
		}
		*target = out
	}
}