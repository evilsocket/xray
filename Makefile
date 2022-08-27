NAME=xray
SOURCE=cmd/$(NAME)/*.go

all: static
	@mkdir -p build		
	go build -o build/$(NAME) $(SOURCE)

test:
	go test -v -cover -race

clean:
	@rm -rf cmd/xray/ui.go
	@rm -rf build

static:
	go-bindata -o cmd/xray/ui.go -pkg main ui