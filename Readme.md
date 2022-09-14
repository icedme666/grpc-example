# grpc基础知识
* 类型
    1. simple RPC
    2. 服务器端流式RPC：Server-side streaming RPC，服务器端流式 RPC，显然是单向流，并代指 Server 为 Stream 而 Client 为普通 RPC 请求
    3. 客户端流式RPC：Client-side streaming RPC，单向流，客户端通过流式发起多次 RPC 请求给服务端，服务端发起一次响应给客户端
    4. 双向流式RPC：Bidirectional streaming RPC，双向流，由客户端以流式的方式发起请求，服务端同样以流式的方式响应请求


* Streaming RPC应用场景
    + 大规模数据包
    + 实时场景
* grpc安装：
  ```bash
  go get -u google.golang.org/grpc
  brew install protobuf
  protoc --version
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
  ```  
    + 文档
        - https://developers.google.com/protocol-buffers/docs/reference/go-generated#package

# Client and Server
1. 生成pb.go文件：proto/search.proto
   ```bash
   protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative search.proto
   ```
   + 字段名称从小写下划线转换为大写驼峰模式（字段导出）
   + 生成一组 Getters 方法，能便于处理一些空指针取值的情况
   + ProtoMessage 方法实现 proto.Message 的接口
   + 生成 Rest 方法，便于将 Protobuf 结构体恢复为零值
   + Repeated 转换为切片
2. server：server/simple_server/server.go
3. client：client/simple_client/client.go
4. 访问
   ```bash
   go run server/simple_server/server.go
   go run client/simple_client/client.go
   ```
   
# 服务器端流式RPC
1. 生成pb.go文件：proto/stream.proto
2. server：server/stream_server/server.go
3. client：client/stream_client/client.go
4. 访问
   ```bash
   go run server/stream_server/server.go
   go run client/stream_client/client.go
   ```