FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/file-server ./cmd/server/main.go

FROM alpine:latest
RUN apk add --no-cache postgresql-client
WORKDIR /app
COPY --from=builder /app/file-server .
COPY entrypoint.sh .
RUN chmod +x ./entrypoint.sh
# EXPOSE 3100
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./file-server"]
