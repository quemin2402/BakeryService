FROM golang:1.23-alpine
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .
RUN mkdir -p /app/logs && chmod 755 /app/logs
EXPOSE 8080
ENV GIN_MODE=release
CMD ["./main"]