src = $(shell find . -name '*.go')
templates = $(shell find . -name '*.html')

server: $(src) $(templates)
	./tailwindcss -i ./assets/input.css -o ./assets/index.css --minify
	go mod tidy
	go build -o ./server cmd/server/main.go

prod: $(src) $(templates)
	mkdir build
	mkdir build/logs
	cp -R cert build
	cp -R assets build
	cp -R views build
	cp .env build
	./tailwindcss -i ./assets/input.css -o build/assets/index.css --minify
	go mod tidy
	go build -o build/server cmd/server/main.go

test: $(src)
	go mod tidy
	go build -o test cmd/test/main.go

.PHONY: clean
clean:
	rm -rf build
