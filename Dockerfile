FROM golang:1.25-alpine AS builder

WORKDIR /app

# Tencent Cloud / CN networks often cannot reach proxy.golang.org.
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM alpine:latest

COPY --from=builder /app/app /app/app
WORKDIR /app
CMD ["./app"]