buildall:
	go build -o bin/explorer ./example/explorer/main.go
	go build -o bin/sitemapper ./example/sitemapper/main.go
	go build -o bin/tianalyser ./example/tianalyser/main.go

ci:
	go test ./...
