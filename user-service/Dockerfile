FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/user-service/

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main /main

EXPOSE 8082

CMD ["/main"]
