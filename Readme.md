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
  + 根证书：根证书（root certificate）是属于根证书颁发机构（CA）的公钥证书。我们可以通过验证 CA 的签名从而信任 CA ，任何人都可以得到 CA 的证书（含公钥），用以验证它所签发的证书（客户端、服务端）
    - 公钥
    - 密钥
  + CSR 是 Cerificate Signing Request 的英文缩写，为证书请求文件。主要作用是 CA 会利用 CSR 文件进行签名使得攻击者无法伪装或篡改原有证书
* 流程
  1. Client 通过请求得到 Server 端的证书
  2. 使用 CA 认证的根证书对 Server 端的证书进行可靠性、有效性等校验
  3. 校验 ServerName 是否可用、有效

# 拦截器
* 功能：每个 RPC 方法的前或后做某些事情
* RPC方法与拦截器的对应：
  + 普通方法：一元拦截器（grpc.UnaryInterceptor）
    - 实现 UnaryServerInterceptor 方法，形参如下：
      - ctx context.Context：请求上下文
      - req interface{}：RPC 方法的请求参数
      - info *UnaryServerInfo：RPC 方法的所有信息
      - handler UnaryHandler：RPC 方法本身
  + 流方法：流拦截器（grpc.StreamInterceptor）
* 多个拦截器：开源项目go-grpc-middleware

# HTTP接口
* 原理：gRPC 的协议是基于 HTTP/2 的，因此应用程序能够在单个 TCP 端口上提供 HTTP/1.1 和 gRPC 接口服务（两种不同的流量）
* 流程
  1. 检测请求协议是否为 HTTP/2
  2. 判断 Content-Type 是否为 application/grpc（gRPC 的默认标识位）
  3. 根据协议的不同转发到不同的服务处理

# 对RPC自定义认证
* PerRPCCredentials
  + 功能： gRPC 默认提供用于自定义认证的接口，作用是将所需的安全认证信息添加到每个 RPC 方法的上下文中
  + 方法：
    - GetRequestMetadata：获取当前请求认证所需的元数据（metadata）
    - RequireTransportSecurity：是否需要基于 TLS 认证进行安全传输

# GRPC Deadlines

# 分布式链路追踪器

# swagger