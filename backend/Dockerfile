FROM golang:1.25

WORKDIR ${GOPATH}/src/backend/
COPY . .

RUN go fmt ./...

RUN go build -o ./bin/app ./cmd/app \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["./bin/app"]

