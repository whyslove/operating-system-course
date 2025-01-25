FROM golang:1.23-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o server .


FROM alpine:3.20

COPY --from=builder /app/server .

CMD ["./server"]