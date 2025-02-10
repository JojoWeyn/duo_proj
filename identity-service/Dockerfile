FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/identity-service/

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main /main

EXPOSE 8081

CMD ["/main"]
