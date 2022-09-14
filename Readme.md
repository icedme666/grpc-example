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

# 证书认证
* TLS证书认证
* 基于CA的TLS证书认证

# 拦截器