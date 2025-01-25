FROM golang:1.23-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o client .

FROM alpine:3.20

COPY --from=builder /app/client .

CMD ["./client"]