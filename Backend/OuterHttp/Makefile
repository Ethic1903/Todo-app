.PHONY: build

build:
	protoc -I proto api/proto/authGens/authRpcService.proto --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go
	protoc -I proto api/proto/notifyGens/notifyRpcService.proto --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go
	protoc -I proto api/proto/taskGens/taskRpcService.proto --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go

.DEfAULT_GOAL := build