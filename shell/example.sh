#!/usr/bin/env sh

echo "启动示例应用..."

# 听、说、唱 服务，互相调用（同步服务）
go run ./examples/services/sing/main.go &
go run ./examples/services/speak/main.go &
go run ./examples/services/listen/main.go &

# 异步消息（事件/订阅）
go run ./examples/services/async/event/main.go &
go run ./examples/services/async/subscriber/main.go &