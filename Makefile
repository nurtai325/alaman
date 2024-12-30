src = $(shell find . -name '*.go')
templates = $(shell find . -name '*.html')

server: $(src) $(templates)
	go mod tidy
	go build -o server cmd/server/main.go

test: $(src)
	go mod tidy
	go build -o test cmd/test/main.go

.PHONY: clean
clean:
	rm server test build
