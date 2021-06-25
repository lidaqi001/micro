# micro

> 基于[asim/go-micro](https://github.com/asim/go-micro)微服务框架(Thanks [asim](https://github.com/asim))

> 进一步封装了一些功能组件，包括对应的示例

- 实现了如下基础服务：

    - 服务治理，
        - 服务限流
        - 服务注册发现
        - 重试，降级，熔断
            - 基于hystrix
            
    - 链路追踪
        - 基于jaeger
        - 同步/异步
        
    - 配置中心
        - [viper](https://github.com/spf13/viper)
            - 基于etcd远程配置
            - 动态更新配置
        
    - 同步服务
        - 客户端
        - 服务端
        
    - 异步服务
        - 事件/订阅
        
    - 对应示例
        - 同步服务
        - 异步服务
        - 多服务互相调用
        
- env
```env
// 限流设置QPS
QPS=100

// 注册中心地址（etcd）
REGISTRY_ADDR=127.0.0.1:2379

// Jaeger Agent地址（6831端口）
TRACE_ADDR=127.0.0.1:6831

// 
jiagnsujiang 
```

- 快速运行

    - 安装依赖包
    ```
  go mod tidy
    ```
    - 运行网关
    ```
  sh shell/api.sh
    ```
    - 运行示例服务
        - 同步服务：sing、speak、listen
            - sing 调用 speak
            - sing 调用 listen
        - 异步服务：async_event、async_subscriber
            - async_event 创建事件
            - async_subscriber 创建订阅
    ```
  sh shell/example.sh
    ```
  
  - 具体路由查看网关定义
  
  ```
  文件地址：examples/api/api.go
  ```  