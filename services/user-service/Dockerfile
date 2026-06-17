FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/user-service ./cmd/main.go

FROM alpine:3.19

COPY --from=builder /app/user-service /app/user-service
COPY migrations /app/migrations

EXPOSE 8080

CMD ["/app/user-service"]


