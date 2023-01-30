# TLS证书认证
* 证书生成
  ```bash
  openssl ecparam -genkey -name secp384r1 -out conf/tls/server.key
  openssl req -new -x509 -sha256 -key conf/tls/server.key -out conf/tls/server.pem -days 3650 -extensions req_ext -config conf/server.conf
  ```
  Common Name为grpc-example，其他为空
* 实现
  + server/tls_server/server.go
  + client/tls_client/client.go
* 访问
   ```bash
   go run server/tls_server/server.go
   go run client/tls_client/client.go
   ```

# 基于CA的TLS证书认证
  
  go 1.15 版本开始废弃 CommonName，需要开启SAN扩展
* 证书生成
  ```bash
  # 1.生成ca秘钥，得到ca.key
  openssl genrsa -out conf/ca/ca.key 4096
  # 2.生成ca证书签发请求，得到ca.csr
  openssl req -new -sha256 -out conf/ca/ca.csr -key conf/ca/ca.key -config conf/ca/ca.conf
  # 3.生成ca根证书，得到ca.crt
  openssl x509 -req -days 3650 -in conf/ca/ca.csr -signkey conf/ca/ca.key -out conf/ca/ca.crt
  ```
* Server
  ```bash
  # 1. 生成私钥，得到server.key
  openssl genrsa -out conf/ca/server/server.key 2048
  # 2. 生成证书签发请求，得到server.csr
  openssl req -new -sha256 -out conf/ca/server/server.csr -key conf/ca/server/server.key -config conf/server.conf
  # 3. 用CA证书生成终端用户证书，得到server.crt
  openssl x509 -req -sha256 -CA conf/ca/ca.crt -CAkey conf/ca/ca.key -CAcreateserial -days 365 -in conf/ca/server/server.csr -out conf/ca/server/server.crt -extensions req_ext -extfile conf/server.conf
  ```
* Client
  ```bash
  # 1. 生成私钥，得到client.key
  openssl genrsa -out conf/ca/client/client.key 2048 
  # 2. 生成证书签发请求，得到client.csr
  openssl req -new -key conf/ca/client/client.key -out conf/ca/client/client.csr -config conf/server.conf
  # 3. 用CA证书生成客户端证书，得到client.crt
  openssl x509 -req -sha256 -CA conf/ca/ca.crt -CAkey conf/ca/ca.key -CAcreateserial -days 365  -in conf/ca/client/client.csr -out conf/ca/client/client.crt
  ```
* 实现
  + server/ca_server/server.go
  + client/ca_client/client.go
* 访问
   ```bash
   go run server/ca_server/server.go
   go run client/ca_client/client.go
   ```
# 参考
* http://www.yaotu.net/biancheng/73602.html
