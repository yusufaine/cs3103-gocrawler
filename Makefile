build-crawler:
	go build -o bin/crawler ./cmd/crawler/main.go
ci:
	go test ./...
