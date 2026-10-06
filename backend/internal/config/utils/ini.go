package utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// ParseINI 将 INI 文件解析为 section → key → value 的映射。
// 支持 # 和 ; 注释行。
func ParseINI(path string) (map[string]map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sections := map[string]map[string]string{"": {}}
	current := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			current = strings.TrimSpace(line[1 : len(line)-1])
			if _, ok := sections[current]; !ok {
				sections[current] = map[string]string{}
			}
			continue
		}

		if k, v, ok := strings.Cut(line, "="); ok {
			key := strings.TrimSpace(k)
			val := strings.TrimSpace(v)
			sections[current][key] = val
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sections, nil
}

// ApplyString 如果 kvs 中存在 key 则写入 target。
func ApplyString(kvs map[string]string, key string, target *string) {
	if v, ok := kvs[key]; ok {
		*target = v
	}
}

// ApplyInt 如果 kvs 中存在 key 则解析 int 后写入 target。
func ApplyInt(kvs map[string]string, key string, target *int) {
	if v, ok := kvs[key]; ok {
		if n, err := strconv.Atoi(v); err == nil {
			*target = n
		}
	}
}

// ApplyBool 如果 kvs 中存在 key 则解析 bool 后写入 target。
// 支持 true/false/1/0/yes/no。
func ApplyBool(kvs map[string]string, key string, target *bool) {
	if v, ok := kvs[key]; ok {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "true", "1", "yes":
			*target = true
		case "false", "0", "no":
			*target = false
		}
	}
}

// ApplyOriginList 解析逗号分隔的字符串列表。
func ApplyOriginList(kvs map[string]string, key string, target *[]string) {
	if v, ok := kvs[key]; ok {
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