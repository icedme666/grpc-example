# TLS证书认证
* 证书生成
  ```bash
  openssl ecparam -genkey -name secp384r1 -out server.key
  openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650 -config server.conf -extensions SAN
  ```
  Common Name为grpc-example，其他为空
* 实现
  + server/tls_server/server.go
  + client/tls_client/client.go

# 基于CA的TLS证书认证
* 证书生成
  ```bash
  # 生成key
  openssl genrsa -out ca.key 2048
  # 生成密钥
  openssl req -new -x509 -days 7200 -key ca.key -out ca.pem
  # Server
  # 生成csr
  openssl req -new -key server.key -out server.csr -config server.conf -extensions SAN
  # 基于CA签发
  openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in server.csr -out server.pem -extensions SAN
  # Client
  # 生成key
  openssl ecparam -genkey -name secp384r1 -out client.key
  # 生成csr
  openssl req -new -key client.key -out client.csr
  # 基于CA签发
  openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.pem
  ```