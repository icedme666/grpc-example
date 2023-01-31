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
* gRPC 中强调 TL;DR（Too long, Don’t read）并建议始终设定截止日期
* 原因：当未设置 Deadlines 时，将采用默认的 DEADLINE_EXCEEDED（这个时间非常大）。如果产生了阻塞等待，就会造成大量正在进行的请求都会被保留，并且所有请求都有可能达到最大超时。这会使服务面临资源耗尽的风险，例如内存，这会增加服务的延迟，或者在最坏的情况下可能导致整个进程崩溃
* 方法
  + context.WithDeadline：会返回最终上下文截止时间。第一个形参为父上下文，第二个形参为调整的截止时间。若父级时间早于子级时间，则以父级时间为准，否则以子级时间为最终截止时间
  + context.WithTimeout：很常见的另外一个方法，是便捷操作。实际上是对于 WithDeadline 的封装
  + status.FromError：返回 GRPCStatus 的具体错误码，若为非法，则直接返回 codes.Unknown
* 逻辑
  + Server：检测超时时间，若Client超时则不再执行逻辑并报错
  + Client：设置截止时间

# 分布式链路追踪器
* 组件：gRPC + Opentracing + Zipkin
  + OpenTracing：通过提供平台无关、厂商无关的API，使得开发人员能够方便的添加（或更换）追踪系统的实现
  + Zipkin：分布式追踪系统，作用是收集解决微服务架构中的延迟问题所需的时序数据。它管理这些数据的收集和查找
    - zipkin.NewHTTPCollector：创建一个 Zipkin HTTP 后端收集器
    - zipkin.NewRecorder：创建一个基于 Zipkin 收集器的记录器
    - zipkin.NewTracer：创建一个 OpenTracing 跟踪器（兼容 Zipkin Tracer）
    - otgrpc.OpenTracingClientInterceptor：返回 grpc.UnaryServerInterceptor，不同点在于该拦截器会在 gRPC Metadata 中查找 OpenTracing SpanContext。如果找到则为该服务的 Span Context 的子节点
    - otgrpc.LogPayloads：设置并返回 Option。作用是让 OpenTracing 在双向方向上记录应用程序的有效载荷（payload）
    - otgrpc.OpenTracingClientInterceptor：返回 grpc.UnaryClientInterceptor。该拦截器的核心功能在于：
      - OpenTracing SpanContext 注入 gRPC Metadata
      - 查看 context.Context 中的上下文关系，若存在父级 Span 则创建一个 ChildOf 引用，得到一个子 Span
* 参考：https://github.com/openzipkin-contrib/zipkin-go-opentracing
# swagger