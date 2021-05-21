package jaeger

import (
	"errors"
	"strings"
)

func getTraceIp(addr string) (string, error) {
	split := strings.Split(addr, ":")
	if len(split) == 2 {
		return split[0], nil
	}
	return "", errors.New("IP解析失败")
}

func empty(ip string) bool {
	return len(ip) == 0
}
