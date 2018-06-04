SHELL := bash

all: xray

godep:
	@go get -u github.com/golang/dep/...

deps: godep
	@dep ensure

xray: deps
	@go build -o xray .

clean:
	@rm -rf xray

install:
	@cp xray /usr/local/bin/

docker:
	@docker build -t xray:latest .
