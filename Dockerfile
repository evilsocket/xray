FROM golang:alpine as build-stage

RUN apk --no-cache add ca-certificates

WORKDIR /go/src/github.com/evilsocket/xray

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /xray ./cmd/xray/*.go

FROM scratch

COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-stage /xray /xray

EXPOSE 8080

ENTRYPOINT ["/xray"]