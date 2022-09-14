# TLS证书认证
* 证书生成
  ```bash
  openssl ecparam -genkey -name secp384r1 -out server.key
  openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
  ```
  Common Name为grpc-example，其他为空
* 
# 基于CA的TLS证书认证