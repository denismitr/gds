deps:
	go mod tidy
	go mod vendor

test:
	go test ./...