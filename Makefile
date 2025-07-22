LOCAL_BIN := $(CURDIR)/bin
export PATH := $(LOCAL_BIN):$(PATH)

.PHONY: install-deps generate generate-chat

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4

generate: generate-chat

generate-chat:
	mkdir -p pkg/chat_v1
	protoc -I ./api/proto/chat_v1 \
		--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
		--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
		--validate_out=lang=go:pkg/chat_v1 --validate_opt=paths=source_relative \
		./api/proto/chat_v1/chat.proto