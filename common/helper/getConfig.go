package helper

import (
	"github.com/lidaqi001/micro/examples/config"
	"os"
	"strconv"
)

// 获取注册中心地址
func GetRegistryAddress() string {
	return GetConfig("REGISTRY_ADDR", config.REGISTRY_ADDR)
}

// 获取链路追踪地址
func GetTraceAddress() string {
	return GetConfig("TRACE_ADDR", config.TRACE_ADDR)
}

// 获取限流QPS
func GetQPS() (float64, int64) {
	qps := GetConfig("QPS", config.QPS)
	qpsf, _ := strconv.ParseFloat(qps, 64)
	qpsi, _ := strconv.ParseInt(qps, 10, 64)
	return qpsf, qpsi
}

func GetConfig(envK string, defaultK string) string {
	env := os.Getenv(envK)
	if len(env) == 0 {
		return defaultK
	}
	return env
}
