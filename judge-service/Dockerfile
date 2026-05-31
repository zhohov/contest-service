FROM golang:1.25

WORKDIR ${GOPATH}/src/judge-service/
COPY . .

RUN go fmt ./...

RUN go build -o ./bin/app ./cmd/app \
    && go clean -cache -modcache

CMD ["./bin/app"]

