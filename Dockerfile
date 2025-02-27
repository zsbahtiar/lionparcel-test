FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o lionparcel-test .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/lionparcel-test .
COPY --from=builder /app/internal/migration ./internal/migration

EXPOSE 8080

CMD ["/bin/sh", "-c", "./lionparcel-test migrate up && ./lionparcel-test start"]
