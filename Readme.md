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
  
# HTTP接口

# 对RPC自定义认证

# GRPC Deadlines

# 分布式链路追踪器

# swagger