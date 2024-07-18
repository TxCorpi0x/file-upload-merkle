MODULE_NAME=github.com/TxCorpi0x/file-upload-merkle
BINARY_NAME=fxmerkle
TEST_FOLDER=.runtime/files

.PHONY: start-server build-client

start-server:
	docker compose up -d

stop-server:
	docker compose down

build-client:
	go build -o $(BINARY_NAME) $(MODULE_NAME)/

test-upload: build-client
	./$(BINARY_NAME) client upload $(TEST_FOLDER)

test-download: build-client
	./$(BINARY_NAME) client download 1
