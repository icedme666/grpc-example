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