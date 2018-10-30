NAME=xray
SOURCE=cmd/$(NAME)/*.go
GOBUILD=go build
DEPEND=github.com/Masterminds/glide
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: deps static
	@mkdir -p build		
	@$(GOBUILD) -o build/$(NAME) $(SOURCE)

format:
	@gofmt -s -w $(GOFILES)

deps: godep
	dep ensure
	go get -u github.com/jteeuwen/go-bindata/...

test:
	go test -v -cover -race $(shell glide novendor)

godep:
	@go get -u github.com/golang/dep/...

clean:
	@rm -rf $(NAME) ui.go
	@rm -rf build

static: format
	go-bindata -o cmd/xray/ui.go -pkg main ui

run:
	go run $(SOURCE)
