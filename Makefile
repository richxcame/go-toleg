dev:
	@go run main.go

protoc:
	@echo "Started compiling..."
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    rpc/gotoleg/transaction.proto
	@echo "Done."

build:
	@echo "Started building..."
	@go build -o bin/gotoleg
	@echo "Done."