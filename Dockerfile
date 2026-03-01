FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o stock_backend .

FROM alpine

COPY --from=builder /app/stock_backend /stock_backend

EXPOSE 8080

CMD ["./stock_backend"]