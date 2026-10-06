package cli

import (
	"flag"
	"strconv"
)


func assignString(src *string, dst *string) {
	if src != nil {
		*dst = *src
	}
}

func assignInt(src *int, dst *int) {
	if src != nil {
		*dst = *src
	}
}

func parseOriginList(s string) []string {
	if s == "" {
		return nil
	}
	// 直接拆分，与 utils.EnvOriginList 逻辑一致
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			token := trimSpace(s[start:i])
			if token != "" {
				out = append(out, token)
			}
			start = i + 1
		}
	}
	return out
}

func trimSpace(s string) string {
	lo, hi := 0, len(s)
	for lo < hi && (s[lo] == ' ' || s[lo] == '\t') {
		lo++
	}
	for hi > lo && (s[hi-1] == ' ' || s[hi-1] == '\t') {
		hi--
	}
	return s[lo:hi]
}


func stringFlag(name, usage string) *string {
	var v *string
	flag.Func(name, usage, func(s string) error {
		v = &s
		return nil
	})
	return v
}

func intFlag(name, usage string) *int {
	var v *int
	flag.Func(name, usage, func(s string) error {
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		v = &n
		return nil
	})
	return v
}

