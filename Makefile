GO=go
SERVER_DIR=cmd/server
CLIENT_DIR=cmd/client
TEST_DIR=./tests
PROTO_DIR=api/proto
SERVER_BINARY=netvuln-server
CLIENT_BINARY=netvuln-client

PORT=50051
TARGETS=192.168.100.14
CLIENT_ADDRESS=localhost:50051

run-server:
	$(GO) run $(SERVER_DIR)/main.go

# запуск клиента с параметрами
run-client:
	$(GO) run $(CLIENT_DIR)/main.go -address=$(CLIENT_ADDRESS) -targets=$(TARGETS) -port=80

test: test-client test-service

#запуск тестов клиентской части
test-client:
	$(GO) test -v $(TEST_DIR)/client

#запуск тестов серверной части
test-service:
	$(GO) test -v $(TEST_DIR)/service

# Генерация go-кода из proto-файлов
generate-proto:
	protoc --go_out=plugins=grpc:$(PROTO_DIR) $(PROTO_DIR)/netvuln.proto

clean:
	rm -f $(SERVER_BINARY) $(CLIENT_BINARY)

# сборкa сервера
build-server:
	$(GO) build -o $(SERVER_BINARY) $(SERVER_DIR)/main.go

# сборка клиента
build-client:
	$(GO) build -o $(CLIENT_BINARY) $(CLIENT_DIR)/main.go

build: build-server build-client

lint:
	golangci-lint run


style:
	@echo "Formatting Go code with gofmt..."
	@gofmt -w .
	@echo "Go code formatting complete."


