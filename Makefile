NAME=xray
SOURCE=cmd/$(NAME)/*.go
GOBUILD=go build
DEPEND=github.com/Masterminds/glide
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: format
format:
	@gofmt -s -w $(GOFILES)

# Command to get glide, you need to run it only once
.PHONY: get_glide
get_glide:
	go get -u -v $(DEPEND)
	$(GOPATH)/bin/glide install

# Command to install dependencies using glide
.PHONY: install_dependencies
install_dependencies:
	glide install

# Run tests in verbose mode with race detector and display coverage
.PHONY: test
test:
	go test -v -cover -race $(shell glide novendor)

# Removing artifacts
.PHONY: clean
clean:
	@rm -rf $(NAME) ui.go
	@rm -rf build

.PHONY: static
static: format
	go-bindata -o cmd/xray/ui.go -pkg main ui

# Building linux binaries
.PHONY: build
build: static
	@mkdir -p build		
	@$(GOBUILD) -o build/$(NAME) $(SOURCE)

# Run the application
.PHONY: run
run:
	go run $(SOURCE)
