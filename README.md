# micro

> 基于[asim/go-micro](https://github.com/asim/go-micro)微服务框架(Thanks [asim](https://github.com/asim))

> 进一步封装了一些功能组件，包括对应的示例

实现了如下基础服务：

- 服务治理，
    - 服务限流
    - 服务注册发现
    - 重试，降级，熔断
        - 基于hystrix
- 链路追踪
    - 基于jaeger
    - 同步/异步
- 同步服务
    - 客户端
    - 服务端
- 异步服务
    - 事件/订阅
- 对应示例
    - 同步服务
    - 异步服务
    - 多服务互相调用