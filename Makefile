ci:
	go test ./...

crawler:
	go build -o bin/crawler ./cmd/crawler/main.go
