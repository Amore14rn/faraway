FROM golang:1.20.3 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./app/cmd/server ./cmd/server

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

# multistage build to copy only binary and config
FROM scratch

COPY --from=builder /build/main /

EXPOSE 3333

ENTRYPOINT ["/main"]