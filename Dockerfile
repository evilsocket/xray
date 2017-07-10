FROM golang:alpine

# Add git and make to alpine
RUN apk add --no-cache git make

# Download and install xray as per instructions
RUN go get github.com/evilsocket/xray && \
    cd $GOPATH/src/github.com/evilsocket/xray/ && \
    make get_glide && \
    make install_dependencies && \
    go get -u github.com/jteeuwen/go-bindata/... && \
    make build

# Default port for xray
EXPOSE 8080

# Settings for run
ENV PATH /go/src/github.com/evilsocket/xray/build/:$PATH

# Build directory
WORKDIR /go/src/github.com/evilsocket/xray/

CMD ["xray"]
