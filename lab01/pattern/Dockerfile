FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /build
COPY go.mod ./
RUN go mod tidy && go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build -o app ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache sqlite-libs ca-certificates

WORKDIR /app
COPY --from=builder /build/app ./app
COPY static/ ./static/

RUN mkdir -p /data

EXPOSE 8082
CMD ["./app"]
