package helper

import (
	"github.com/lidaqi001/micro/examples/config"
	"os"
	"strconv"
)

// 获取注册中心地址
func GetRegistryAddress() string {
	addr := os.Getenv("REGISTRY_ADDR")
	if len(addr) == 0 {
		return config.REGISTRY_ADDR
	}
	return addr
}

// 获取注册中心地址
func GetTraceAddress() string {
	addr := os.Getenv("TRACE_ADDR")
	if len(addr) == 0 {
		return config.TRACE_ADDR
	}
	return addr
}

// 获取限流QPS
func GetQPS() (float64, int64) {
	var qpsi int64
	var qpsf float64

	qps := os.Getenv("QPS")

	if len(qps) == 0 {
		// 未设置env，使用默认配置
		qps := config.QPS
		qpsf = float64(qps)
		qpsi = int64(qps)

	} else {
		qpsf, _ = strconv.ParseFloat(qps, 64)
		qpsi, _ = strconv.ParseInt(qps, 10, 64)
	}

	return qpsf, qpsi
}
