buildall:
	go build -o example/bin/sitemapgenerator ./example/sitemapgenerator/main.go
	go build -o example/bin/tianalyser ./example/tianalyser/main.go

ci:
	go test ./...

crawler:
	go build -o bin/crawler ./cmd/crawler/main.go
