.PHONY: test

test:
	go test ./...

lint:
	golangci-lint run

dev: test lint

linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o dist/tender-linux-amd64 ./cmd/tender
	cd ./dist/ && cp ./tender-linux-amd64 ./tender && tar cfz tender-linux-amd64.tar.gz ./tender

linux_arm64:
	GOOS=linux GOARCH=arm64 go build -o dist/tender-linux-arm64 ./cmd/tender
	cd ./dist/ && cp ./tender-linux-arm64 ./tender && tar cfz tender-linux-arm64.tar.gz ./tender

darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -o dist/tender-darwin-amd64 ./cmd/tender
	cd ./dist/ && cp ./tender-darwin-amd64 ./tender && tar cfz tender-darwin-amd64.tar.gz ./tender

darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build -o dist/tender-darwin-arm64 ./cmd/tender
	cd ./dist/ && cp ./tender-darwin-arm64 ./tender && tar cfz tender-darwin-arm64.tar.gz ./tender

all: linux_amd64 linux_arm64 darwin_amd64 darwin_arm64

clean:
	rm ./dist/tender-*
