# Build โปรแกรม
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /app/server \
    ./cmd/api


# นำไฟล์ที่ build แล้วไปรัน
FROM alpine:3.22

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/server ./server

USER appuser

EXPOSE 8080

CMD ["./server"]