FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o job-poll-queue main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/job-poll-queue /app/job-poll-queue

EXPOSE 3000
EXPOSE 50051

CMD ["./job-poll-queue"]