APP_NAME=generator
GRPC_API_PROTO_PATH:=./api/grpc

PROTOC_GEN_GO_VERSION      := v1.33.0
PROTOC_GEN_GO_GRPC_VERSION := v1.3.0
PROTOC_VERSION             := 3.20.3

## Testing

.PHONY: test-coverprofile
test-coverprofile:
	@go test ./... -count=1 -cover -coverprofile=cover.out

.PHONY: test
test: test-coverprofile # run tests
	@go tool cover -func=cover.out

.PHONY: test-coverage
test-coverage: _test-coverprofile # run tests and show coverage
	@go tool cover -html cover.out

## Grpc generation

GRPC_INSTALL_SOURCE_WIN:=https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-win64.zip
GRPC_INSTALL_SOURCE_LIN:=https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-linux-x86_64.zip
GRPC_INSTALL_FILENAME:=third_party/protoc.zip

.PHONY: install-grpc
install-grpc:
ifeq ($(OS),Windows_NT) # Windows
	@mkdir -p third_party/
	@powershell -Command "Invoke-WebRequest -OutFile ${GRPC_INSTALL_FILENAME} -Uri ${GRPC_INSTALL_SOURCE_WIN}"
	@echo "$(CYAN)Downloaded protoc to $(RESET)";
	@powershell -Command "Expand-Archive -Path ${GRPC_INSTALL_FILENAME} -DestinationPath third_party/protoc -Force"
	@echo "$(CYAN)Unzipped protoc to third_party/protoc$(RESET)";
	@rm ${GRPC_INSTALL_FILENAME}
	@LOCAL_VERSION=`third_party/protoc/bin/protoc.exe --version | cut -d ' ' -f 2`; \
	echo "$(CYAN)Installed protoc version $$LOCAL_VERSION$(RESET)";
else # Linux
	@wget -qO ${GRPC_INSTALL_FILENAME} ${GRPC_INSTALL_SOURCE_LIN}
	@echo "$(CYAN)Downloaded protoc$(RESET)";
	@unzip -qod third_party/protoc ${GRPC_INSTALL_FILENAME}
	@echo "$(CYAN)Unzipped protoc to third_party/protoc$(RESET)";
	@rm -f ${GRPC_INSTALL_FILENAME}
	@LOCAL_VERSION=`third_party/protoc/bin/protoc --version | cut -d ' ' -f 2`; \
	echo "$(CYAN)Installed protoc version $$LOCAL_VERSION$(RESET)";
endif
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	@LOCAL_VERSION=`protoc-gen-go --version | cut -d ' ' -f 2`; \
	echo "$(CYAN)Installed protoc-gen-go version $$LOCAL_VERSION$(RESET)";
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	@LOCAL_VERSION=`protoc-gen-go-grpc --version | cut -d ' ' -f 2`; \
	echo "$(CYAN)Installed protoc-gen-go-grpc version $$LOCAL_VERSION$(RESET)";

GRPC_PKG_DIR:=./pkg/api/grpc
GRPC_API_PROTO_PATH=./api/grpc
GOPATH_BIN := $(shell go env GOPATH)\bin
PROTOC_GEN_GO = $(GOPATH_BIN)\protoc-gen-go
PROTOC_GEN_GO_GRPC = $(GOPATH_BIN)\protoc-gen-go-grpc

.PHONY: grpc-gen
grpc-gen:
	@mkdir -p ${GRPC_PKG_DIR}
	./third_party/protoc/bin/protoc -I=${GRPC_API_PROTO_PATH} \
			--plugin=$(PROTOC_GEN_GO) \
			--plugin=$(PROTOC_GEN_GO_GRPC) \
			--go_out=${GRPC_PKG_DIR} \
			--go_opt=paths=source_relative \
			--go-grpc_out=${GRPC_PKG_DIR} \
			--go-grpc_opt=paths=source_relative \
	${GRPC_API_PROTO_PATH}/${APP_NAME}.proto