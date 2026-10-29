# Multi-stage Go build for Railway
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Build the application
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Final stage - minimal runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
