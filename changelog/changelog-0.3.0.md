# 拦截器实现：
* 实现：server/simple_server/server.go
  + logging：RPC 方法的入参出参的日志输出
  + recover：RPC 方法的异常保护和日志输出
* 访问
  1. 验证logging
  ```bash
  go run server/interceptor_server/server.go
  go run client/ca_client/client.go
  ```
  2. 验证recover
  在 RPC 方法中人为地制造运行时错误，再重复启动 server/client.go，得到报错

# HTTP接口
* 实现
  + 封装获取证书的方法：pkg/gtls
  + Server：server/simple_http_server/server.go
  + Client：client/simple_http_client/client.go
* 验证
  ```bash
  go run server/simple_http_server/server.go
  go run client/simple_http_client/client.go
  curl https://127.0.0.1:9003 -k  #或者请求postman
  ```

# 自定义认证
* 实现
  + Server：server/simple_token_server/server.go
    - 调用 metadata.FromIncomingContext 从上下文中获取 metadata，再在不同的 RPC 方法中进行认证检查
  + Client：client/simple_token_client/client.go
    - 重点实现 type PerRPCCredentials interface，关注：
      - struct Auth：GetRequestMetadata、RequireTransportSecurity
      - grpc.WithPerRPCCredentials
* 验证
  ```bash
  go run server/simple_token_server/server.go
  go run client/simple_token_client/client.go
  curl https://127.0.0.1:9003 -k  #或者请求postman
  ```
* 进阶：将type PerRPCCredentials interface 做成一个拦截器