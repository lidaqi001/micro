#!/usr/bin/env sh


# 听、说、唱 服务，互相调用（同步服务）

echo "启动示例服务：sing"
go run ./examples/services/sing/main.go &

echo "启动示例服务：speak"
go run ./examples/services/speak/main.go &

echo "启动示例服务：listen"
go run ./examples/services/listen/main.go &

# 异步消息（事件/订阅）

echo "启动示例Event：async_rabbitmq"
go run ./examples/services/async_rabbitmq/event/main.go &

echo "启动示例Subscribe：async_rabbitmq"
go run ./examples/services/async_rabbitmq/subscriber/main.go &