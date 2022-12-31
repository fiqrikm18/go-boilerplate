CWD=$(shell pwd)

.PHONY: certs
certs:
	openssl genrsa -out ./certs/app.rsa 2048 && \
	openssl rsa -in ./certs/app.rsa -pubout > ./certs/app.rsa.pub

test:
	export APP_MODE=test && \
	go test ./... -v

dev:
	export APP_CONFIG_PATH=$(CWD)/config && \
	go run main.go