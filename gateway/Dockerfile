FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

COPY public.pem /app/public.pem
COPY publicRef.pem /app/publicRef.pem

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/gateway/

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main /main

COPY --from=builder /app/public.pem /app/public.pem
COPY --from=builder /app/publicRef.pem /app/publicRef.pem

EXPOSE 3211

CMD ["/main"]
