FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

COPY private.pem /app/private.pem
COPY privateRef.pem /app/privateRef.pem
COPY public.pem /app/public.pem
COPY publicRef.pem /app/publicRef.pem

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/identity-service/

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main /main

COPY --from=builder /app/private.pem /app/private.pem
COPY --from=builder /app/privateRef.pem /app/privateRef.pem
COPY --from=builder /app/public.pem /app/public.pem
COPY --from=builder /app/publicRef.pem /app/publicRef.pem

EXPOSE 8081

CMD ["/main"]
