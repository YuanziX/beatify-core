build:
	@go build -o bin/beatify-core

run: build
	@./bin/beatify-core
