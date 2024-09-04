pull:
	@git pull

build:
	@go build -o bin/beatify-core

run: build
	@./bin/beatify-core

status:
	@sudo systemctl status beatify-core

serve: pull build
	@sudo systemctl restart beatify-core nginx
	@sleep 1
	@$(MAKE) status
