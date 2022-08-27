NAME=xray
SOURCE=cmd/$(NAME)/*.go

all: cmd/xray/ui.go
	@mkdir -p build		
	go build -o build/$(NAME) $(SOURCE)

test:
	go test -v -cover -race

clean:
	@rm -rf cmd/xray/ui.go
	@rm -rf build

cmd/xray/ui.go:
	go-bindata -o cmd/xray/ui.go -pkg main ui