# Deadlines
* 实现
  + Server：server/deadline_http_server/server.go
  + Client：client/deadline_http_client/client.go
* 验证
  1. 
  ```bash
  go run server/deadline_http_server/server.go
  go run client/deadline_http_client/client.go  #超时报错
  ```
  2. 去掉服务端循环的逻辑，client再次请求则成功

# 分布式链路追踪
* 实现：
  1. 搭建Zipkin
    ```bash
    docker run -d -p 9411:9411 openzipkin/zipkin
    ```
  2. 初始化Zipkin
     + Server：包含收集器、记录器、跟踪器，再利用拦截器在 Server 端实现 SpanContext、Payload 的双向读取和管理
       - server/simple_zipkin_server/server.go
     + Client：实现拦截器
       - client/simple_zipkin_client/client.go 
* 验证
  1. 启动服务，并访问
  ```bash
    go run server/simple_zipkin_server/server.go
    go run client/simple_zipkin_client/client.go  #超时报错
  ```
  2. 查看http://127.0.0.1:9411/zipkin，点击查询