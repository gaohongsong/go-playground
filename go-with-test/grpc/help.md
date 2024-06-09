# quickstart

https://grpc.io/docs/languages/go/quickstart/

# install

- protocol buffer compiler
> https://grpc.io/docs/protoc-installation/
```
brew install protobuf
protoc --version
```
- go plugins for protocol compiler
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

# example
https://github.com/grpc/grpc-go/blob/master/examples/

# render proto file
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  helloworld.proto
