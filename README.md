## Berikut adalah cara penggunaan file proto untuk gRPC:

1. Selalu update repository proto dengan yang terbaru
2. Copy folder proto disini ke dalam folder project go
3. Ubah option go_package ke nama package go yang diinginkan berdasarkan nama folder project go tempat protogen disimpan
4. Hapus folder protogen, agar dapat data terbaru dari protogen
5. Jalankan perintah `protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto`
6. Jalankan perintah `go mod tidy`
7. Jalankan perintah `go run cmd/main.go`

## install Framework & Library

- GoFiber (HTTP Framework) : go get -u github.com/gofiber/fiber/v2
- Protobuf : google.golang.org/protobuf
- Protobuf GRPC : go get -u google.golang.org/grpc
- Viper (Configuration) : go get github.com/spf13/viper
- Logrus (Logging) : go get github.com/sirupsen/logrus
- Go Playground Validator (Validation) : go get github.com/go-playground/validator/v10
- Gorm (ORM) : go get -u gorm.io/gorm || go get -u gorm.io/driver/postgres


### Install Plugin Protocol Buffers Go
```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

### Install Plugin Protocol Buffers Go GRPC

```shell
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Run Generate Proto to Go GRPC

```shell
protoc --go_opt=module={module_name} --go_out=. ./proto/*.proto
protoc --go-grpc_opt=module={module_name} --go-grpc_out=. ./proto/*.proto

# example:
protoc --go_opt=module=grpc-server-1 --go_out=. ./proto/user/*.proto
protoc --go-grpc_opt=module=grpc-server-1 --go-grpc_out=. ./proto/user/*.proto
```

### Install Dependency

```bash
go mod tidy
```

### Run GRPC server

```bash
go run cmd/grpc/main.go
```