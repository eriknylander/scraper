test:
	go test ./...

build: test
	go build -o scraper.exe cmd/scraper/main.go