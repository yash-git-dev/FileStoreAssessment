FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum .env ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]
