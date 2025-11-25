FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/db-seeder ./cmd/seeder/main.go

FROM alpine:latest
RUN apk add --no-cache postgresql-client
WORKDIR /app
COPY --from=builder /app/api-server .
COPY --from=builder /app/db-seeder .
COPY entrypoint.sh .
RUN chmod +x ./entrypoint.sh
# EXPOSE 3100
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./api-server"]
